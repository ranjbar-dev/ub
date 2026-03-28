package payment

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"exchange-go/internal/communication"
	"exchange-go/internal/currency"
	"exchange-go/internal/externalexchange"
	"exchange-go/internal/platform"
	"exchange-go/internal/response"
	"exchange-go/internal/transaction"
	"exchange-go/internal/user"
	"exchange-go/internal/userbalance"
	"exchange-go/internal/userwithdrawaddress"
	"exchange-go/internal/wallet"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	TypeWithdraw = "WITHDRAW"
	TypeDeposit  = "DEPOSIT"

	StatusCreated      = "CREATED"
	StatusPending      = "PENDING"
	StatusCompleted    = "COMPLETED"
	StatusInProgress   = "IN_PROGRESS"
	StatusFailed       = "FAILED"
	StatusCanceled     = "CANCELED"
	StatusUserCanceled = "USER_CANCELED"
	StatusRejected     = "REJECTED"

	AdminStatusPending  = "PENDING"
	AdminStatusRecheck  = "RECHECK"
	AdminStatusApproved = "APPROVED"

	WithdrawTypeExternalExchange = "EXTERNAL_EXCHANGE"
	WithdrawTypeHotWallet        = "HOT_WALLET"
)

type GetPaymentsParams struct {
	Type      string `form:"type"`
	Coin      string `form:"code"`
	StartDate string `form:"start_date"`
	EndDate   string `form:"end_date"`
	Page      int64  `form:"page"`
	PageSize  int    `form:"page_size"`
}

type PreWithdrawParams struct {
	Coin    string `json:"code" binding:"required"`
	Amount  string `json:"amount" binding:"required"`
	Address string `json:"address" binding:"required"`
	Network string `json:"network"`
}

type WithdrawParams struct {
	Coin      string `json:"code" binding:"required"`
	Amount    string `json:"amount" binding:"required"`
	Address   string `json:"address" binding:"required"`
	Label     string `json:"label"`
	TwoFaCode string `json:"2fa_code"`
	EmailCode string `json:"email_code"`
	Network   string `json:"network"`
	IP        string
}

type GetPaymentDetailParams struct {
	ID int64 `form:"id" binding:"required,gt=0"`
}

type GetPaymentFilters struct {
	UserID    int
	Type      string
	Coin      string
	StartDate string
	EndDate   string
	Page      int64
	PageSize  int
}

type GetPaymentsResponse struct {
	ID                 int64  `json:"id"`
	Status             string `json:"status"`
	Type               string `json:"type"`
	Amount             string `json:"amount"`
	Coin               string `json:"code"`
	CreatedAt          string `json:"createdAt"`
	Address            string `json:"address"`
	AddressExplorerURL string `json:"addressExplorerUrl"`
	TxID               string `json:"txId"`
	TxExplorerURL      string `json:"txIdExplorerUrl"`
}

type DetailResponse struct {
	Address            string `json:"address"`
	AddressExplorerURL string `json:"addressExplorerUrl"`
	TxID               string `json:"txId"`
	TxExplorerURL      string `json:"txIdExplorerUrl"`
	RejectionReason    string `json:"rejectionReason"`
}

type DetailQueryFields struct {
	ID              int64
	Code            string
	Network         string
	Address         string
	UserID          int
	TxID            string
	RejectionReason string
}

type PreWithdrawResponse struct {
	Need2fa       bool `json:"need2fa"`
	NeedEmailCode bool `json:"needEmailCode"`
}

type userWithdrawValidationParams struct {
	toAddress            string
	twoFaCode            string
	amount               string
	emailCode            string
	shouldCheckEmailCode bool
	shouldCheck2fa       bool
}

type ExternalWithdrawalUpdateDataNeeded struct {
	PaymentID                  int64
	PaymentExtraInfoID         int64
	UpdatedAt                  time.Time
	ExternalExchangeWithdrawID string
}

type UpdatePaymentInExternalExchangeParams struct {
	PaymentID          int64
	PaymentExtraInfoID int64
	TxID               string
	Status             string
	Data               string
}

type CancelWithdrawParams struct {
	ID int64 `json:"id" binding:"required,gt=0"`
}
type TransactionMeta struct {
	InternalTransferID string `json:"internal_transfer_id" binding:"omitempty"`
}
type WalletCallBackParams struct {
	Code        string `json:"code" binding:"required"`
	Amount      string `json:"amount" binding:"required"`
	Type        string `json:"type" binding:"required"`
	Status      string `json:"status" binding:"required"`
	FromAddress string `json:"from_address" binding:"required"`
	ToAddress   string `json:"to_address" binding:"required"`
	TxID        string `json:"tx_id" binding:"required"`
	Meta        string `json:"meta" binding:"omitempty"`
	Network     string `json:"network"`
}

type WalletCallbackMetadata struct {
	InternalTransferID string `json:"internal_transfer_id"`
}

type paymentPushPayload struct {
	ID     int64  `json:"id"`
	Status string `json:"status"`
}

type UpdateWithdrawParams struct {
	ID              int64  `json:"id" binding:"required"`
	Status          string `json:"status"`
	AdminStatus     string `json:"admin_status"`
	Fee             string `json:"fee"`
	NetworkFee      string `json:"network_fee"`
	AutoTransfer    *bool  `json:"auto_transfer"`
	RejectionReason string `json:"rejection_reason"`
	WithdrawType    string `json:"withdraw_type"`
}

type UpdateDepositParams struct {
	ID            int64  `json:"id" binding:"required"`
	Status        string `json:"status"`
	Amount        string `json:"amount"`
	FromAddress   string `json:"from_address"`
	ToAddress     string `json:"to_address"`
	TxID          string `json:"tx_id"`
	ShouldDeposit *bool  `json:"should_deposit"`
}

// Service provides the public API for payment operations including deposit and
// withdrawal listing, withdrawal initiation and cancellation, and admin callbacks.
type Service interface {
	// GetPayments returns a paginated list of the user's deposit and withdrawal payments.
	GetPayments(u *user.User, params GetPaymentsParams) (apiResponse response.APIResponse, statusCode int)
	// PreWithdraw initiates the pre-withdrawal flow, sending an email confirmation code.
	PreWithdraw(u *user.User, params PreWithdrawParams) (apiResponse response.APIResponse, statusCode int)
	// Withdraw executes a confirmed withdrawal after email code verification.
	Withdraw(u *user.User, params WithdrawParams) (apiResponse response.APIResponse, statusCode int)
	// Detail returns detailed information for a single payment.
	Detail(u *user.User, params GetPaymentDetailParams) (apiResponse response.APIResponse, statusCode int)
	// GetInProgressWithdrawalsInExternalExchange returns withdrawals pending on the external exchange.
	GetInProgressWithdrawalsInExternalExchange() []ExternalWithdrawalUpdateDataNeeded
	// UpdatePaymentInExternalExchange updates a payment's status based on external exchange data.
	UpdatePaymentInExternalExchange(params UpdatePaymentInExternalExchangeParams)
	// CancelWithdraw cancels a pending withdrawal request.
	CancelWithdraw(u *user.User, params CancelWithdrawParams) (apiResponse response.APIResponse, statusCode int)

	//for admin
	// HandleWalletCallBack processes deposit/withdrawal callbacks from the wallet service.
	HandleWalletCallBack(u *user.User, params WalletCallBackParams) (apiResponse response.APIResponse, statusCode int)
	// UpdateWithdraw allows an admin to approve or reject a pending withdrawal.
	UpdateWithdraw(u *user.User, params UpdateWithdrawParams) (apiResponse response.APIResponse, statusCode int)
	// UpdateDeposit allows an admin to manually update a deposit record.
	UpdateDeposit(u *user.User, params UpdateDepositParams) (apiResponse response.APIResponse, statusCode int)
}

type service struct {
	db                               *gorm.DB
	paymentRepository                Repository
	currencyService                  currency.Service
	walletService                    wallet.Service
	userConfigService                user.ConfigService
	twoFaManager                     user.TwoFaManager
	withdrawEmailConfirmationManager WithdrawEmailConfirmationManager
	permissionManager                user.PermissionManager
	userService                      user.Service
	userBalanceService               userbalance.Service
	communicationService             communication.Service
	priceGenerator                   currency.PriceGenerator
	userWithdrawAddressService       userwithdrawaddress.Service
	internalTransferService          InternalTransferService
	externalExchangeService          externalexchange.Service
	autoExchangeManager              AutoExchangeManager
	mqttManager                      communication.MqttManager
	configs                          platform.Configs
	logger                           platform.Logger
}

func (s *service) GetPayments(u *user.User, params GetPaymentsParams) (apiResponse response.APIResponse, statusCode int) {
	result := make([]GetPaymentsResponse, 0)
	filters := s.getFiltersForPayments(params)
	filters.UserID = u.ID

	payments := s.paymentRepository.GetUserPayments(filters)
	for _, p := range payments {
		status := p.Status
		if p.Type == TypeWithdraw && status == StatusCreated { //we do not want to show the user the created status
			status = StatusPending
		}

		if p.Type == TypeDeposit && status == StatusCreated { //we do not want to show the user the created status
			status = StatusInProgress
		}

		r := GetPaymentsResponse{
			ID:                 p.ID,
			Status:             removeUnderline(strings.ToLower(status)),
			Type:               strings.ToLower(p.Type),
			Amount:             p.Amount.String,
			Coin:               p.Coin.Code,
			CreatedAt:          p.CreatedAt.Format("2006-01-02 15:04:05"),
			Address:            p.ToAddress.String,
			AddressExplorerURL: wallet.GetAddressExplorer(p.Coin.Code, p.BlockchainNetwork.String, p.ToAddress.String),
			TxID:               p.TxID.String,
			TxExplorerURL:      wallet.GetTxExplorer(p.Coin.Code, p.BlockchainNetwork.String, p.TxID.String),
		}

		result = append(result, r)
	}

	data := map[string][]GetPaymentsResponse{
		"payments": result,
	}

	return response.Success(data, "")

}

func (s *service) getFiltersForPayments(params GetPaymentsParams) GetPaymentFilters {
	var filters GetPaymentFilters
	filters.Type = strings.ToUpper(params.Type)
	filters.Coin = strings.ToUpper(params.Coin)
	if params.Page >= 0 {
		filters.Page = params.Page
	}

	filters.PageSize = params.PageSize

	if filters.PageSize == 0 {
		filters.PageSize = 20
	}

	if filters.PageSize > 50 {
		filters.PageSize = 50
	}

	if filters.Type != TypeDeposit && filters.Type != TypeWithdraw {
		filters.Type = ""
	}

	_, err := time.Parse("2006-01-02 15:04:05", params.StartDate)
	if err == nil {
		filters.StartDate = params.StartDate
	}

	_, err = time.Parse("2006-01-02 15:04:05", params.EndDate)
	if err == nil {
		filters.EndDate = params.EndDate
	}

	return filters

}

func (s *service) Detail(u *user.User, params GetPaymentDetailParams) (apiResponse response.APIResponse, statusCode int) {
	id := params.ID
	detailFields, err := s.paymentRepository.GetPaymentDetailByID(id)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("can not get payment", err,
			zap.String("service", "paymentService"),
			zap.String("method", "Detail"),
			zap.Int64("paymentID", id),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) || detailFields.ID == 0 || detailFields.UserID != u.ID {
		return response.Error("payment not found", http.StatusUnprocessableEntity, nil)
	}
	res := DetailResponse{
		Address:            detailFields.Address,
		AddressExplorerURL: wallet.GetAddressExplorer(detailFields.Code, detailFields.Network, detailFields.Address),
		TxID:               detailFields.TxID,
		TxExplorerURL:      wallet.GetTxExplorer(detailFields.Code, detailFields.Network, detailFields.TxID),
		RejectionReason:    detailFields.RejectionReason,
	}
	return response.Success(res, "")

}

func (s *service) PreWithdraw(u *user.User, params PreWithdrawParams) (apiResponse response.APIResponse, statusCode int) {
	coinCode := strings.ToUpper(strings.Trim(params.Coin, ""))
	network := strings.ToUpper(strings.Trim(params.Network, ""))
	var coin currency.Coin

	if coinCode == "" {
		return response.Error("coin not found", http.StatusUnprocessableEntity, nil)
	}
	coin, err := s.currencyService.GetCoinByCode(coinCode)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("can not get coin by code", err,
			zap.String("service", "paymentService"),
			zap.String("method", "PreWithdraw"),
			zap.String("code", coinCode),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || coin.ID == 0 {
		return response.Error("coin not found", http.StatusUnprocessableEntity, nil)
	}

	if network != "" {
		// check if network exists
		parentCoin, err := s.currencyService.GetCoinByCode(network)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error2("can not get coin by code for parent coin", err,
				zap.String("service", "paymentService"),
				zap.String("method", "PreWithdraw"),
				zap.String("code", network),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}

		if errors.Is(err, gorm.ErrRecordNotFound) || parentCoin.ID == 0 {
			return response.Error("network not found", http.StatusUnprocessableEntity, nil)
		}
	}

	if network == "" && coin.BlockchainNetwork.Valid {
		network = strings.ToUpper(coin.BlockchainNetwork.String)
	}

	isValid, err := s.walletService.IsAddressValid(coin.Code, params.Address, network)

	if err != nil {
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if !isValid {
		return response.Error("address is not valid", http.StatusUnprocessableEntity, nil)
	}

	amountDecimal, err := decimal.NewFromString(params.Amount)
	if err != nil || !amountDecimal.IsPositive() {
		return response.Error("amount is not correct", http.StatusUnprocessableEntity, nil)
	}

	validationParams := userWithdrawValidationParams{
		toAddress:            params.Address,
		amount:               params.Amount,
		shouldCheckEmailCode: false,
		shouldCheck2fa:       false,
	}

	var uc *user.Config
	userConfig, err := s.userConfigService.GetUserConfig(u.ID)
	if err == nil {
		uc = &userConfig
	}

	err = s.validateIfUserCanWithdraw(u, uc, coin, validationParams)
	if err != nil {
		return response.Error(err.Error(), http.StatusUnprocessableEntity, nil)
	}

	//check user balance
	ub := &userbalance.UserBalance{}
	err = s.userBalanceService.GetBalanceOfUserByCoinID(u.ID, coin.ID, ub)
	if err != nil {
		s.logger.Error2("can not get balance of user by coin ID", err,
			zap.String("service", "paymentService"),
			zap.String("method", "PreWithdraw"),
			zap.Int64("codeID", coin.ID),
			zap.Int("userID", u.ID),
		)
		return response.Error("can not get userBalance right now", http.StatusUnprocessableEntity, nil)
	}
	userBalanceAmountDecimal, _ := decimal.NewFromString(ub.Amount)
	userBalanceFrozenAmountDecimal, _ := decimal.NewFromString(ub.FrozenAmount)
	userBalanceRemainingAmountDecimal := userBalanceAmountDecimal.Sub(userBalanceFrozenAmountDecimal)
	if userBalanceRemainingAmountDecimal.LessThan(amountDecimal) {
		return response.Error("user balance is not enough to withdraw this much", http.StatusUnprocessableEntity, nil)
	}
	//end of check user balance

	res := PreWithdrawResponse{
		Need2fa:       false,
		NeedEmailCode: true,
	}

	if u.IsTwoFaEnabled {
		res.Need2fa = true
	}

	if uc != nil {
		res.Need2fa = uc.IsTwoFaVerificationForWithdrawEnabled
		//res.NeedEmailCode = uc.IsEmailVerificationForWithdrawEnabled
	}

	if res.NeedEmailCode {
		isAllowed, err := s.withdrawEmailConfirmationManager.IsAllowedToSendEmail(*u, coin.Code, params.Amount, params.Address)
		if err != nil {
			s.logger.Error2("error checking isAllowedToSendEmail", err,
				zap.String("service", "paymentService"),
				zap.String("method", "PreWithdraw"),
				zap.String("codeCode", coin.Code),
				zap.Int("userID", u.ID),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}
		if isAllowed {
			err := s.withdrawEmailConfirmationManager.CreateAndSendWithdrawEmailConfirmationCode(*u, coin.Code, params.Amount, params.Address)
			if err != nil {
				s.logger.Error2("can not create and send withdraw emailConfirmation code", err,
					zap.String("service", "paymentService"),
					zap.String("method", "PreWithdraw"),
					zap.String("codeCode", coin.Code),
					zap.String("amount", params.Amount),
					zap.String("address", params.Address),
					zap.Int("userID", u.ID),
				)
				return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
			}
		} else {
			return response.Error("only one email per minute can be sent", http.StatusUnprocessableEntity, nil)
		}
	}

	return response.Success(res, "")

}

func (s *service) validateIfUserCanWithdraw(u *user.User, uc *user.Config, coin currency.Coin, params userWithdrawValidationParams) error {

	if coin.SupportsWithdraw.Valid && !coin.SupportsWithdraw.Bool {
		return fmt.Errorf("withdrawal is not supported now")
	}

	if u.Status != user.StatusVerified {
		return fmt.Errorf("user account is not verified")
	}

	if u.Google2faDisabledAt.Valid {
		disabledAtTimestamp := u.Google2faDisabledAt.Time.Unix()
		now := time.Now().Unix()
		if now-disabledAtTimestamp < 24*60*60 {
			return fmt.Errorf("for the security reasons, after disabling/enabling 2fa the withdraw request is not allowed for 24 hours")
		}
	}

	//todo currently we have no user config setting in front, ignore following if, unless the feature is implemented in front
	if uc != nil && uc.ID > 0 {
		if params.shouldCheck2fa {
			if uc.IsTwoFaVerificationForWithdrawEnabled {
				if params.twoFaCode == "" {
					return fmt.Errorf("2fa code is not provided")
				}

				if !s.twoFaManager.CheckCode(*u, params.twoFaCode) {
					return fmt.Errorf("2fa code is not correct")
				}
			}
		}

		if uc.IsReadOnly {
			return fmt.Errorf("this account is in read only mode")
		}

		//check whitelist
		if uc.IsWhiteListEnabled {
			addresses := s.userWithdrawAddressService.GetUserWithdrawAddressesByAddress(u, coin, params.toAddress)
			if len(addresses) < 1 {
				return fmt.Errorf("this address is not in white list")
			}
		}

		//canUserWithdrawToThisAddress
		if params.shouldCheckEmailCode {
			if uc.IsEmailVerificationForWithdrawEnabled {
				if params.emailCode == "" {
					return fmt.Errorf("email code is not provided")
				}

				isCorrect, err := s.withdrawEmailConfirmationManager.CheckCode(*u, params.emailCode)
				if err != nil {
					s.logger.Error2("can not check email confirmation code", err,
						zap.String("service", "paymentService"),
						zap.String("method", "validateIfUserCanWithdraw"),
						zap.String("codeCode", coin.Code),
						zap.Int("userID", u.ID),
					)
					return fmt.Errorf("something went wrong")
				}
				if !isCorrect {
					return fmt.Errorf("email confirmation code is not correct")
				}

			}
		}
	}

	if !s.permissionManager.IsPermissionGrantedToUserFor(*u, user.PermissionWithdraw) {
		return fmt.Errorf("withdraw permission is not granted")
	}

	minimumWithdrawDecimal, _ := decimal.NewFromString(coin.MinimumWithdraw)
	amountDecimal, _ := decimal.NewFromString(params.amount)

	if amountDecimal.LessThan(minimumWithdrawDecimal) {
		return fmt.Errorf("minimum withdraw is: " + coin.MinimumWithdraw)
	}

	maximumWithdrawDecimal, _ := decimal.NewFromString(coin.MaximumWithdraw)

	if amountDecimal.GreaterThan(maximumWithdrawDecimal) {
		return fmt.Errorf("maximum withdraw is: " + coin.MaximumWithdraw)
	}

	//always checking for email code no matter what
	if params.shouldCheckEmailCode {
		isCorrect, err := s.withdrawEmailConfirmationManager.CheckCode(*u, params.emailCode)
		if err != nil {
			s.logger.Error2("can not check email confirmation code", err,
				zap.String("service", "paymentService"),
				zap.String("method", "validateIfUserCanWithdraw"),
				zap.String("codeCode", coin.Code),
				zap.Int("userID", u.ID),
			)
			return fmt.Errorf("something went wrong")
		}
		if !isCorrect {
			return fmt.Errorf("email confirmation code is not correct")
		}
	}

	if u.IsTwoFaEnabled {
		if params.shouldCheck2fa {
			if params.twoFaCode == "" {
				return fmt.Errorf("2fa code is not provided")
			}

			if !s.twoFaManager.CheckCode(*u, params.twoFaCode) {
				return fmt.Errorf("2fa code is not correct")
			}

		}
	}

	return nil

}

func (s *service) Withdraw(u *user.User, params WithdrawParams) (apiResponse response.APIResponse, statusCode int) {
	coinCode := strings.ToUpper(strings.Trim(params.Coin, ""))
	network := strings.ToUpper(strings.Trim(params.Network, ""))
	label := strings.Trim(params.Label, "")
	var coin currency.Coin

	if coinCode == "" {
		return response.Error("coin not found", http.StatusUnprocessableEntity, nil)
	}
	coin, err := s.currencyService.GetCoinByCode(coinCode)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("can not get coin by code", err,
			zap.String("service", "paymentService"),
			zap.String("method", "Withdraw"),
			zap.String("codeCode", coinCode),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) || coin.ID == 0 {
		return response.Error("coin not found", http.StatusUnprocessableEntity, nil)
	}

	// check if network exists
	if network != "" {
		parentCoin, err := s.currencyService.GetCoinByCode(network)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			s.logger.Error2("can not get coin by code for parent coin", err,
				zap.String("service", "paymentService"),
				zap.String("method", "Withdraw"),
				zap.String("code", network),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}

		if errors.Is(err, gorm.ErrRecordNotFound) || parentCoin.ID == 0 {
			return response.Error("network not found", http.StatusUnprocessableEntity, nil)
		}
	}
	if network == "" && coin.BlockchainNetwork.Valid {
		network = strings.ToUpper(coin.BlockchainNetwork.String)
	}

	isValid, err := s.walletService.IsAddressValid(coin.Code, params.Address, network)
	if err != nil {
		s.logger.Error2("can not check if address is valid", err,
			zap.String("service", "paymentService"),
			zap.String("method", "Withdraw"),
			zap.String("code", network),
		)
		return response.Error("can not check if address is valid", http.StatusUnprocessableEntity, nil)
	}

	if !isValid {
		return response.Error("address is not valid", http.StatusUnprocessableEntity, nil)
	}

	amountDecimal, err := decimal.NewFromString(params.Amount)
	if err != nil || !amountDecimal.IsPositive() {
		return response.Error("amount is not correct", http.StatusUnprocessableEntity, nil)
	}

	//trim email code
	params.EmailCode = strings.Trim(params.EmailCode, "")

	validationParams := userWithdrawValidationParams{
		toAddress:            params.Address,
		twoFaCode:            params.TwoFaCode,
		amount:               params.Amount,
		emailCode:            params.EmailCode,
		shouldCheckEmailCode: true,
		shouldCheck2fa:       true,
	}
	var uc *user.Config
	userConfig, err := s.userConfigService.GetUserConfig(u.ID)
	if err == nil {
		uc = &userConfig
	}

	err = s.validateIfUserCanWithdraw(u, uc, coin, validationParams)
	if err != nil {
		return response.Error(err.Error(), http.StatusUnprocessableEntity, nil)
	}

	if label != "" {
		go func() {
			userWithdrawAddressParams := userwithdrawaddress.CreateAddressParams{
				Coin:    coin.Code,
				Label:   label,
				Address: params.Address,
				Network: network,
			}
			_, err := s.userWithdrawAddressService.SaveNewAddress(u, coin, userWithdrawAddressParams)
			if err != nil {
				//we only log here we do not want stop the flow since this part is not crucial
				s.logger.Error2("can not save new address", err,
					zap.String("service", "paymentService"),
					zap.String("method", "Withdraw"),
					zap.Int64("codeID", coin.ID),
					zap.Int("userID", u.ID),
				)
			}
		}()
	}

	tx := s.db.Begin()
	err = tx.Error
	if err != nil {
		s.logger.Error2("error in starting transaction", err,
			zap.String("service", "paymentService"),
			zap.String("method", "Withdraw"),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	//check user balance
	ub := &userbalance.UserBalance{}
	err = s.userBalanceService.GetBalanceOfUserByCoinUsingTx(tx, u.ID, coin.ID, ub)
	if err != nil {
		tx.Rollback()
		s.logger.Error2("can not get user balance by coin id", err,
			zap.String("service", "paymentService"),
			zap.String("method", "Withdraw"),
			zap.Int64("codeID", coin.ID),
			zap.Int("userID", u.ID),
		)
		return response.Error("can not get userBalance right now", http.StatusUnprocessableEntity, nil)
	}
	userBalanceAmountDecimal, _ := decimal.NewFromString(ub.Amount)
	userBalanceFrozenAmountDecimal, _ := decimal.NewFromString(ub.FrozenAmount)
	userBalanceRemainingAmountDecimal := userBalanceAmountDecimal.Sub(userBalanceFrozenAmountDecimal)
	if userBalanceRemainingAmountDecimal.LessThan(amountDecimal) {
		tx.Rollback()
		return response.Error("user balance is not enough to withdraw this much", http.StatusUnprocessableEntity, nil)
	}
	//end of check user balance

	blockchainNetwork := sql.NullString{String: "", Valid: false}
	if network != "" {
		blockchainNetwork = sql.NullString{String: network, Valid: true}
	}
	withdrawFee := coin.WithdrawalFee.Float64
	withdrawFeeString := strconv.FormatFloat(withdrawFee, 'f', 8, 64)
	if network != coin.BlockchainNetwork.String {
		otherNetworksConfigs, err := coin.GetOtherBlockchainNetworksConfigs()
		if err != nil {
			s.logger.Error2("can not get other network configs", err,
				zap.String("service", "paymentService"),
				zap.String("method", "Withdraw"),
				zap.Int64("codeID", coin.ID),
				zap.String("configs", coin.OtherBlockchainNetworksConfigs.String),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}
		for _, nc := range otherNetworksConfigs {
			if nc.Code == network {
				withdrawFeeString = nc.Fee
				break
			}

		}

	}

	//create payment
	p := &Payment{
		UserID:            u.ID,
		CoinID:            coin.ID,
		Type:              TypeWithdraw,
		Status:            StatusCreated,
		AdminStatus:       sql.NullString{String: AdminStatusPending, Valid: true},
		Code:              coin.Code,
		ToAddress:         sql.NullString{String: params.Address, Valid: true},
		BlockchainNetwork: blockchainNetwork,
		Amount:            sql.NullString{String: amountDecimal.StringFixed(8), Valid: true},
		FeeAmount:         sql.NullString{String: withdrawFeeString, Valid: true},
	}

	err = tx.Omit(clause.Associations).Save(p).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("can not save payment", err,
			zap.String("service", "paymentService"),
			zap.String("method", "Withdraw"),
			zap.Int64("codeID", coin.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	ctx := context.Background()

	//create extra payment info
	extraInfo := &ExtraInfo{
		PaymentID: p.ID,
		IP:        sql.NullString{String: params.IP, Valid: true},
	}

	//setting price, we do this to have a history of prices in user every withdraw and deposit
	btcPrice, err := s.priceGenerator.GetBTCUSDTPrice(ctx)
	if err == nil {
		extraInfo.BtcPrice = sql.NullString{String: btcPrice, Valid: true}
	}

	coinPriceBasedOnBtc, err := s.priceGenerator.GetAmountBasedOnBTC(ctx, coin.Code, "1.0")
	if err == nil {
		extraInfo.Price = sql.NullString{String: coinPriceBasedOnBtc, Valid: true}
	}
	err = tx.Omit(clause.Associations).Save(extraInfo).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("can not save extraInfo", err,
			zap.String("service", "paymentService"),
			zap.String("method", "Withdraw"),
			zap.Int64("codeID", coin.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	//freeze balance
	frozenDecimal, err := decimal.NewFromString(ub.FrozenAmount)
	if err != nil {
		tx.Rollback()
		s.logger.Error2("can not convert frozenAmount to Decimal", err,
			zap.String("service", "paymentService"),
			zap.String("method", "Withdraw"),
			zap.String("FrozenAmount", ub.FrozenAmount),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}
	newFrozenAmount := frozenDecimal.Add(amountDecimal).StringFixed(8)
	ub.FrozenAmount = newFrozenAmount
	err = tx.Omit(clause.Associations).Save(ub).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("can not save user balance", err,
			zap.String("service", "paymentService"),
			zap.String("method", "Withdraw"),
			zap.Int64("userBalanceID", ub.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("can not commit transaction", err,
			zap.String("service", "paymentService"),
			zap.String("method", "Withdraw"),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	//remove confirmation code from redis
	go func() {
		if params.EmailCode != "" {
			err := s.withdrawEmailConfirmationManager.RemoveConfirmationCodeFromRedis(*u, coin.Code, params.Amount, params.Address)
			if err != nil {
				s.logger.Warn("can not remove confirmation code from redis",
					zap.Error(err),
					zap.String("service", "paymentService"),
					zap.String("method", "Withdraw"),
				)
			}

		}
	}()

	go s.notifyUserPaymentStatusUpdate(*u, *p, "")

	result := make([]GetPaymentsResponse, 0)
	status := p.Status
	if p.Type == TypeWithdraw && status == StatusCreated { //we do not want to show the user the created status
		status = StatusPending
	}

	r := GetPaymentsResponse{
		ID:                 p.ID,
		Status:             removeUnderline(strings.ToLower(status)),
		Type:               strings.ToLower(p.Type),
		Amount:             p.Amount.String,
		Coin:               p.Code,
		CreatedAt:          p.CreatedAt.Format("2006-01-02 15:04:05"),
		Address:            p.ToAddress.String,
		AddressExplorerURL: wallet.GetAddressExplorer(coin.Code, p.BlockchainNetwork.String, p.ToAddress.String),
		TxID:               p.TxID.String,
		TxExplorerURL:      wallet.GetTxExplorer(coin.Code, p.BlockchainNetwork.String, p.TxID.String),
	}
	result = append(result, r)

	data := map[string][]GetPaymentsResponse{
		"payments": result,
	}

	return response.Success(data, "")
}

func (s *service) notifyUserPaymentStatusUpdate(u user.User, payment Payment, RejectionReason string) {
	if u.Email == "" {
		//the user is not loaded we get it from database
		var err error
		u, err = s.userService.GetUserByID(u.ID)
		if err != nil {
			s.logger.Error2("can not get user by id", err,
				zap.String("service", "paymentService"),
				zap.String("method", "notifyUserPaymentStatusUpdate"),
				zap.Int("userID", u.ID),
			)
			return
		}

	}
	profile, err := s.userService.GetUserProfile(u)
	if err != nil {
		s.logger.Error2("can not get userProfile", err,
			zap.String("service", "paymentService"),
			zap.String("method", "notifyUserPaymentStatusUpdate"),
			zap.Int("userID", u.ID),
		)
		return
	}
	params := communication.CryptoPaymentStatusUpdateEmailParams{
		FullName:        profile.GetFullName(),
		Type:            payment.Type,
		Status:          payment.Status,
		CurrencyCode:    payment.Code,
		Amount:          payment.Amount.String,
		ToAddress:       payment.ToAddress.String,
		RejectionReason: RejectionReason,
	}
	cu := communication.CommunicatingUser{
		Email: u.Email,
		Phone: "",
	}
	s.communicationService.SendCryptoPaymentStatusUpdateEmail(cu, params)
}

func removeUnderline(text string) string {
	return strings.Replace(text, "_", " ", -1)
}

func (s *service) GetInProgressWithdrawalsInExternalExchange() []ExternalWithdrawalUpdateDataNeeded {
	return s.paymentRepository.GetInProgressWithdrawalsInExternalExchange()
}

func (s *service) UpdatePaymentInExternalExchange(params UpdatePaymentInExternalExchangeParams) {
	status := strings.ToUpper(params.Status)
	if status == StatusInProgress {
		return
	}
	tx := s.db.Begin()
	err := tx.Error
	if err != nil {
		s.logger.Error2("error in starting transaction", err,
			zap.String("service", "paymentService"),
			zap.String("method", "UpdatePaymentInExternalExchange"),
		)
		return
	}
	p := &Payment{}
	err = s.paymentRepository.GetPaymentByIDUsingTx(tx, params.PaymentID, p)
	if err != nil {
		tx.Rollback()
		s.logger.Error2("can not get payment by id", err,
			zap.String("service", "paymentService"),
			zap.String("method", "UpdatePaymentInExternalExchange"),
			zap.Int64("paymentID", params.PaymentID),
		)
		return
	}
	p.Status = status
	if p.Status == StatusCompleted {
		p.TxID = sql.NullString{String: params.TxID, Valid: true}
	}

	err = tx.Omit(clause.Associations).Save(p).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("can not save payment", err,
			zap.String("service", "paymentService"),
			zap.String("method", "UpdatePaymentInExternalExchange"),
			zap.Int64("paymentID", params.PaymentID),
		)
		return
	}

	pei := &ExtraInfo{}
	err = s.paymentRepository.GetExtraInfoByPaymentIDUsingTx(tx, params.PaymentID, pei)
	if err != nil {
		tx.Rollback()
		s.logger.Error2("can not get extra info by payment ID", err,
			zap.String("service", "paymentService"),
			zap.String("method", "UpdatePaymentInExternalExchange"),
			zap.Int64("paymentID", params.PaymentID),
		)
		return
	}
	pei.ExternalExchangeWithdrawInfo = sql.NullString{String: params.Data, Valid: true}
	err = tx.Omit(clause.Associations).Save(pei).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("can not save payment extra info", err,
			zap.String("service", "paymentService"),
			zap.String("method", "UpdatePaymentInExternalExchange"),
			zap.Int64("paymentID", params.PaymentID),
		)
		return
	}

	ub := &userbalance.UserBalance{}
	err = s.userBalanceService.GetBalanceOfUserByCoinUsingTx(tx, p.UserID, p.CoinID, ub)
	if err != nil {
		tx.Rollback()
		s.logger.Error2("can not get balance of user by coin id", err,
			zap.String("service", "paymentService"),
			zap.String("method", "UpdatePaymentInExternalExchange"),
			zap.Int("userID", p.UserID),
			zap.Int64("coinID", p.CoinID),
		)
		return
	}

	amountDecimal, err := decimal.NewFromString(p.Amount.String)
	if err != nil {
		tx.Rollback()
		s.logger.Error2("can not convert amount to Decimal", err,
			zap.String("service", "paymentService"),
			zap.String("method", "UpdatePaymentInExternalExchange"),
			zap.String("amount", p.Amount.String),
			zap.Int64("coinID", p.CoinID),
		)
		return
	}

	userBalanceAmountDecimal, err := decimal.NewFromString(ub.Amount)
	if err != nil {
		tx.Rollback()
		s.logger.Error2("can not convert userbalance amount to Decimal", err,
			zap.String("service", "paymentService"),
			zap.String("method", "UpdatePaymentInExternalExchange"),
			zap.String("amount", ub.Amount),
			zap.Int64("coinID", p.CoinID),
			zap.Int64("ubID", ub.ID),
		)
		return
	}
	userBalanceFrozenAmountDecimal, err := decimal.NewFromString(ub.FrozenAmount)
	if err != nil {
		tx.Rollback()
		s.logger.Error2("can not convert userbalance frozen amount to Decimal", err,
			zap.String("service", "paymentService"),
			zap.String("method", "UpdatePaymentInExternalExchange"),
			zap.String("amount", ub.FrozenAmount),
			zap.Int64("coinID", p.CoinID),
			zap.Int64("ubID", ub.ID),
		)
		return
	}

	switch p.Status {
	case StatusCompleted:
		newFrozenAmountDecimal := userBalanceFrozenAmountDecimal.Sub(amountDecimal)
		if newFrozenAmountDecimal.IsNegative() {
			tx.Rollback()
			err := fmt.Errorf("frozen amount is negative")
			s.logger.Error2("frozen balance in negative", err,
				zap.String("service", "paymentService"),
				zap.String("method", "UpdatePaymentInExternalExchange"),
				zap.String("amount", p.Amount.String),
				zap.String("frozenAmount", ub.FrozenAmount),
				zap.Int64("coinID", p.CoinID),
				zap.String("status", p.Status),
			)
			return
		}
		newFrozenAmount := newFrozenAmountDecimal.StringFixed(8)
		ub.FrozenAmount = newFrozenAmount

		newAmountDecimal := userBalanceAmountDecimal.Sub(amountDecimal)
		if newAmountDecimal.IsNegative() {
			tx.Rollback()
			err := fmt.Errorf("amount is negative")
			s.logger.Error2("balance in negative", err,
				zap.String("service", "paymentService"),
				zap.String("method", "UpdatePaymentInExternalExchange"),
				zap.String("amount", p.Amount.String),
				zap.String("frozenAmount", ub.Amount),
				zap.Int64("coinID", p.CoinID),
				zap.String("status", p.Status),
			)
			return
		}
		newAmount := newAmountDecimal.StringFixed(8)
		ub.Amount = newAmount
		err = tx.Omit(clause.Associations).Save(ub).Error
		if err != nil {
			tx.Rollback()
			s.logger.Error2("can not save user balance", err,
				zap.String("service", "paymentService"),
				zap.String("method", "UpdatePaymentInExternalExchange"),
				zap.Int64("userBalanceID", ub.ID),
				zap.String("status", p.Status),
			)
			return
		}

		withdrawTransaction := &transaction.Transaction{
			UserID:    p.UserID,
			CoinID:    p.CoinID,
			Type:      transaction.TypeWithdraw,
			Amount:    sql.NullString{String: p.Amount.String, Valid: true},
			CoinName:  p.Code,
			PaymentID: sql.NullInt64{Int64: p.ID, Valid: true},
		}

		err = tx.Omit(clause.Associations).Create(withdrawTransaction).Error
		if err != nil {
			tx.Rollback()
			s.logger.Error2("can not create withdraw transaction", err,
				zap.String("service", "paymentService"),
				zap.String("method", "UpdatePaymentInExternalExchange"),
				zap.Int64("paymentID", p.ID),
				zap.String("status", p.Status),
			)
			return
		}

		withdrawFeeTransaction := &transaction.Transaction{
			UserID:    p.UserID,
			CoinID:    p.CoinID,
			Type:      transaction.TypeWithdrawFee,
			Amount:    sql.NullString{String: p.FeeAmount.String, Valid: true},
			CoinName:  p.Code,
			PaymentID: sql.NullInt64{Int64: p.ID, Valid: true},
		}

		err = tx.Omit(clause.Associations).Create(withdrawFeeTransaction).Error
		if err != nil {
			tx.Rollback()
			s.logger.Error2("can not create withdraw fee transaction", err,
				zap.String("service", "paymentService"),
				zap.String("method", "UpdatePaymentInExternalExchange"),
				zap.Int64("paymentID", p.ID),
				zap.String("status", p.Status),
			)
			return
		}
	case StatusRejected, StatusCanceled, StatusFailed:
		newFrozenAmountDecimal := userBalanceFrozenAmountDecimal.Sub(amountDecimal)
		if newFrozenAmountDecimal.IsNegative() {
			tx.Rollback()
			err := fmt.Errorf("frozen amount is negative")
			s.logger.Error2("frozen balance in negative", err,
				zap.String("service", "paymentService"),
				zap.String("method", "UpdatePaymentInExternalExchange"),
				zap.String("amount", p.Amount.String),
				zap.String("frozenAmount", ub.FrozenAmount),
				zap.Int64("coinID", p.CoinID),
				zap.String("status", p.Status),
			)
			return
		}
		newFrozenAmount := newFrozenAmountDecimal.StringFixed(8)
		ub.FrozenAmount = newFrozenAmount
		err = tx.Omit(clause.Associations).Save(ub).Error
		if err != nil {
			tx.Rollback()
			s.logger.Error2("can not save user balance", err,
				zap.String("service", "paymentService"),
				zap.String("method", "UpdatePaymentInExternalExchange"),
				zap.Int64("userBalanceID", ub.ID),
				zap.String("status", p.Status),
			)
			return
		}
	default:
		//we should not be here ever
		tx.Rollback()
		s.logger.Error2("can not save user balance", err,
			zap.String("service", "paymentService"),
			zap.String("method", "UpdatePaymentInExternalExchange"),
			zap.Int64("paymentID", p.ID),
			zap.String("status", p.Status),
		)
		return
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("can not commit the transaction", err,
			zap.String("service", "paymentService"),
			zap.String("method", "UpdatePaymentInExternalExchange"),
		)
		return
	}

	u := user.User{
		ID: p.UserID,
	}
	go s.notifyUserPaymentStatusUpdate(u, *p, "")
}

func (s *service) CancelWithdraw(u *user.User, params CancelWithdrawParams) (apiResponse response.APIResponse, statusCode int) {
	tx := s.db.Begin()
	err := tx.Error
	if err != nil {
		s.logger.Error2("error in starting transaction", err,
			zap.String("service", "paymentService"),
			zap.String("method", "CancelWithdraw"),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}
	p := &Payment{}
	err = s.paymentRepository.GetPaymentByIDUsingTx(tx, params.ID, p)

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		s.logger.Error2("can not get payment by Id", err,
			zap.String("service", "paymentService"),
			zap.String("method", "CancelWithdraw"),
			zap.Int64("paymentID", params.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || p.ID == 0 || p.UserID != u.ID {
		tx.Rollback()
		return response.Error("withdraw not found", http.StatusUnprocessableEntity, nil)
	}

	if p.Type != TypeWithdraw {
		tx.Rollback()
		return response.Error("withdraw not found", http.StatusUnprocessableEntity, nil)
	}

	if p.Status != StatusCreated {
		tx.Rollback()
		return response.Error("withdraw can't be cancelled now", http.StatusUnprocessableEntity, nil)
	}

	p.Status = StatusUserCanceled
	err = tx.Omit(clause.Associations).Save(p).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("can not save payement", err,
			zap.String("service", "paymentService"),
			zap.String("method", "CancelWithdraw"),
			zap.Int64("paymentID", params.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	ub := &userbalance.UserBalance{}
	err = s.userBalanceService.GetBalanceOfUserByCoinUsingTx(tx, p.UserID, p.CoinID, ub)
	if err != nil {
		tx.Rollback()
		s.logger.Error2("can not get balance of user", err,
			zap.String("service", "paymentService"),
			zap.String("method", "CancelWithdraw"),
			zap.Int("userID", p.UserID),
			zap.Int64("coinID", p.CoinID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	amountDecimal, err := decimal.NewFromString(p.Amount.String)
	if err != nil {
		tx.Rollback()
		s.logger.Error2("can not convert amount to Decimal", err,
			zap.String("service", "paymentService"),
			zap.String("method", "CancelWithdraw"),
			zap.String("amount", p.Amount.String),
			zap.Int64("paymentID", params.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	userBalanceFrozenAmountDecimal, err := decimal.NewFromString(ub.FrozenAmount)
	if err != nil {
		tx.Rollback()
		s.logger.Error2("can not convert frozen amount to Decimal", err,
			zap.String("service", "paymentService"),
			zap.String("method", "CancelWithdraw"),
			zap.String("frozenAmount", ub.FrozenAmount),
			zap.Int64("userBalanceID", ub.ID),
			zap.Int64("paymentID", params.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	newFrozenAmountDecimal := userBalanceFrozenAmountDecimal.Sub(amountDecimal)
	if newFrozenAmountDecimal.IsNegative() {
		tx.Rollback()
		err := fmt.Errorf("frozen amount is negative")
		s.logger.Error2("can not convert frozen amount to Decimal", err,
			zap.String("service", "paymentService"),
			zap.String("method", "CancelWithdraw"),
			zap.String("amount", p.Amount.String),
			zap.String("frozenAmount", ub.FrozenAmount),
			zap.Int64("userBalanceID", ub.ID),
			zap.Int64("paymentID", params.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	newFrozenAmount := newFrozenAmountDecimal.StringFixed(8)
	ub.FrozenAmount = newFrozenAmount
	err = tx.Omit(clause.Associations).Save(ub).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("can not save user balance", err,
			zap.String("service", "paymentService"),
			zap.String("method", "CancelWithdraw"),
			zap.Int64("userBalanceID", ub.ID),
			zap.Int64("paymentID", params.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("can not commit transaction", err,
			zap.String("service", "paymentService"),
			zap.String("method", "CancelWithdraw"),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	//just to have same response as withdraw created and payments lists
	result := make([]GetPaymentsResponse, 0)
	r := GetPaymentsResponse{
		ID:                 p.ID,
		Status:             removeUnderline(strings.ToLower(p.Status)),
		Type:               strings.ToLower(p.Type),
		Amount:             p.Amount.String,
		Coin:               p.Code,
		CreatedAt:          p.CreatedAt.Format("2006-01-02 15:04:05"),
		Address:            p.ToAddress.String,
		AddressExplorerURL: wallet.GetAddressExplorer(p.Code, p.BlockchainNetwork.String, p.ToAddress.String),
		TxID:               p.TxID.String,
		TxExplorerURL:      wallet.GetTxExplorer(p.Code, p.BlockchainNetwork.String, p.TxID.String),
	}
	result = append(result, r)

	data := map[string][]GetPaymentsResponse{
		"payments": result,
	}

	return response.Success(data, "")
}

func (s *service) HandleWalletCallBack(u *user.User, params WalletCallBackParams) (apiResponse response.APIResponse, statusCode int) {
	params.Code = strings.ToUpper(params.Code)
	params.Network = strings.ToUpper(params.Network)
	params.Type = strings.ToUpper(params.Type)
	params.Status = strings.ToUpper(params.Status)

	//meta is always would be sent from wallet
	if params.Meta != "" {
		meta := &WalletCallbackMetadata{}
		err := json.Unmarshal([]byte(params.Meta), meta)
		if err != nil {
			s.logger.Error2("can not unmarshal meta ", err,
				zap.String("service", "paymentService"),
				zap.String("method", "HandleWalletCallBack"),
				zap.String("coin", params.Code),
				zap.String("type", params.Type),
				zap.String("txID", params.TxID),
				zap.String("meta", params.Meta),
			)
			return response.Error(err.Error(), http.StatusUnprocessableEntity, nil)
		}
		//if InternalTransferID is not empty string we handle it here, if it is then we have regular crypto payment
		if meta.InternalTransferID != "" {
			err = s.handleInternalTransferCallBack(meta.InternalTransferID, params.Status)
			if err != nil {
				s.logger.Error2("can not handle internal trnasfer callback", err,
					zap.String("service", "paymentService"),
					zap.String("method", "HandleWalletCallBack"),
					zap.String("internalTransferId", meta.InternalTransferID),
					zap.String("coin", params.Code),
					zap.String("type", params.Type),
					zap.String("txID", params.TxID),
				)
				return response.Error(err.Error(), http.StatusUnprocessableEntity, nil)
			}
			return response.Success(nil, "")

		}
	}
	coin, err := s.currencyService.GetCoinByCode(params.Code)
	if err != nil || coin.ID == 0 {
		s.logger.Error2("can not get coin from db", err,
			zap.String("service", "paymentService"),
			zap.String("method", "HandleWalletCallBack"),
			zap.String("coin", params.Code),
			zap.String("type", params.Type),
			zap.String("txID", params.TxID),
		)
	}
	if params.Network == "" && coin.BlockchainNetwork.Valid && coin.BlockchainNetwork.String != "" {
		params.Network = strings.ToUpper(coin.BlockchainNetwork.String)
	}
	if params.Type == TypeDeposit {
		err := s.handleDepositCallBack(params, coin)
		if err != nil {
			s.logger.Error2("can not handle deposit callback", err,
				zap.String("service", "paymentService"),
				zap.String("method", "HandleWalletCallBack"),
				zap.String("coin", params.Code),
				zap.String("type", params.Type),
				zap.String("txID", params.TxID),
			)
			return response.Error(err.Error(), http.StatusUnprocessableEntity, nil)
		}
	} else {
		err := s.handleWithdrawCallBack(params, coin)
		if err != nil {
			s.logger.Error2("can not handle withdraw callback", err,
				zap.String("service", "paymentService"),
				zap.String("method", "HandleWalletCallBack"),
				zap.String("coin", params.Code),
				zap.String("type", params.Type),
				zap.String("txID", params.TxID),
			)
			return response.Error(err.Error(), http.StatusUnprocessableEntity, nil)
		}
	}
	return response.Success(nil, "")
}

func (s *service) handleDepositCallBack(params WalletCallBackParams, coin currency.Coin) error {
	tx := s.db.Begin()
	err := tx.Error
	if err != nil {
		return err
	}
	ub := &userbalance.UserBalance{}
	err = s.userBalanceService.GetUserBalanceByCoinAndAddressUsingTx(tx, coin.ID, params.ToAddress, ub)
	if err != nil {
		tx.Rollback()
		return err
	}
	beforeStatus := ""
	p := &Payment{}
	err = s.paymentRepository.GetPaymentByCoinIDAndTxIDAndTypeUsingTx(tx, coin.ID, params.TxID, params.Type, p)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		return err
	}
	if p.ID > 0 {
		// this means the payment already exists and we should update it
		beforeStatus = p.Status
		if p.Status == StatusCompleted {
			tx.Rollback()
			return fmt.Errorf("payment with id %d is already completed", p.ID)
		}

		if params.Status == StatusCreated {
			tx.Rollback()
			return fmt.Errorf("payment with id %d is already exists but the comming status is created", p.ID)
		}

		p.TxID = sql.NullString{String: params.TxID, Valid: true}
		p.Status = params.Status
		err = tx.Omit(clause.Associations).Save(p).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		amountDecimal, err := decimal.NewFromString(params.Amount)
		if err != nil {
			tx.Rollback()
			return err
		}
		balanceAmountDecimal, err := decimal.NewFromString(ub.Amount)
		if err != nil {
			tx.Rollback()
			return err
		}
		finalBalanceAmountDecimal := balanceAmountDecimal.Add(amountDecimal)
		ub.Amount = finalBalanceAmountDecimal.StringFixed(8)
		err = tx.Omit(clause.Associations).Save(ub).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		transaction := &transaction.Transaction{
			UserID:    ub.UserID,
			CoinID:    coin.ID,
			OrderID:   sql.NullInt64{Int64: 0, Valid: false},
			Type:      transaction.TypeDeposit,
			Amount:    sql.NullString{String: amountDecimal.StringFixed(8), Valid: true},
			CoinName:  coin.Code,
			PaymentID: sql.NullInt64{Int64: p.ID, Valid: true},
		}
		err = tx.Omit(clause.Associations).Save(transaction).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	} else {
		// this means this is a new created deposit and we should create new one
		amountDecimal, err := decimal.NewFromString(params.Amount)
		if err != nil {
			tx.Rollback()
			return err
		}
		p = &Payment{
			UserID:            ub.UserID,
			CoinID:            coin.ID,
			Type:              TypeDeposit,
			Status:            params.Status,
			Code:              coin.Code,
			FromAddress:       sql.NullString{String: params.FromAddress, Valid: true},
			ToAddress:         sql.NullString{String: params.ToAddress, Valid: true},
			TxID:              sql.NullString{String: params.TxID, Valid: true},
			BlockchainNetwork: sql.NullString{String: params.Network, Valid: true},
			Amount:            sql.NullString{String: amountDecimal.StringFixed(8), Valid: true},
		}
		err = tx.Omit(clause.Associations).Save(p).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		extraInfo := &ExtraInfo{
			PaymentID: p.ID,
			IP:        sql.NullString{String: "", Valid: false},
		}

		ctx := context.Background()
		//setting price, we do this to have a history of prices in user every withdraw and deposit
		btcPrice, err := s.priceGenerator.GetBTCUSDTPrice(ctx)
		if err == nil {
			extraInfo.BtcPrice = sql.NullString{String: btcPrice, Valid: true}
		}

		coinPriceBasedOnBtc, err := s.priceGenerator.GetAmountBasedOnBTC(ctx, coin.Code, "1.0")
		if err == nil {
			extraInfo.Price = sql.NullString{String: coinPriceBasedOnBtc, Valid: true}
		}
		err = tx.Omit(clause.Associations).Save(extraInfo).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		if params.Status == StatusCompleted {
			balanceAmountDecimal, err := decimal.NewFromString(ub.Amount)
			if err != nil {
				tx.Rollback()
				return err
			}
			finalBalanceAmountDecimal := balanceAmountDecimal.Add(amountDecimal)
			ub.Amount = finalBalanceAmountDecimal.StringFixed(8)
			err = tx.Omit(clause.Associations).Save(ub).Error
			if err != nil {
				tx.Rollback()
				return err
			}
			transaction := &transaction.Transaction{
				UserID:    ub.UserID,
				CoinID:    coin.ID,
				OrderID:   sql.NullInt64{Int64: 0, Valid: false},
				Type:      transaction.TypeDeposit,
				Amount:    sql.NullString{String: amountDecimal.StringFixed(8), Valid: true},
				CoinName:  coin.Code,
				PaymentID: sql.NullInt64{Int64: p.ID, Valid: true},
			}
			err = tx.Omit(clause.Associations).Save(transaction).Error
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return err
	}
	user := user.User{ID: ub.UserID}
	if beforeStatus != p.Status {
		go s.notifyUserPaymentStatusUpdate(user, *p, "")
	}
	go s.publishPaymentToUser(user, *p)

	//only for completed status we try to autoExchange
	if params.Status == StatusCompleted {
		go s.autoExchangeManager.AutoExchange(p, ub)
	}

	return nil
}

func (s *service) publishPaymentToUser(u user.User, p Payment) {
	if u.PrivateChannelName == "" {
		//the user is not loaded we get it from database
		var err error
		u, err = s.userService.GetUserByID(u.ID)
		if err != nil {
			s.logger.Error2("can not get user by id", err,
				zap.String("service", "paymentService"),
				zap.String("method", "publishPaymentToUser"),
				zap.Int("userID", u.ID),
			)
			return
		}
	}
	pushData := paymentPushPayload{
		ID:     p.ID,
		Status: p.Status,
	}
	payload, err := json.Marshal(pushData)
	if err != nil {
		s.logger.Error2("can not marshal payment push load", err,
			zap.String("service", "paymentService"),
			zap.String("method", "publishPaymentToUser"),
			zap.Int("userID", u.ID),
		)
	}
	s.mqttManager.PublishPayment(context.Background(), u.PrivateChannelName, payload)
}

func (s *service) handleWithdrawCallBack(params WalletCallBackParams, coin currency.Coin) error {
	tx := s.db.Begin()
	err := tx.Error
	if err != nil {
		return err
	}
	p := &Payment{}
	err = s.paymentRepository.GetPaymentByCoinIDAndTxIDAndTypeUsingTx(tx, coin.ID, params.TxID, params.Type, p)
	if err != nil {
		tx.Rollback()
		return err
	}
	beforeStatus := p.Status
	if p.Status == StatusCompleted {
		tx.Rollback()
		return fmt.Errorf("payment with id %d is already completed", p.ID)
	}
	if params.Status != StatusCompleted && params.Status != StatusFailed {
		tx.Rollback()
		return fmt.Errorf("withdraw payment with id %d, the comming status is not completed or failed", p.ID)
	}
	p.Status = params.Status

	if params.FromAddress != "" {
		p.FromAddress = sql.NullString{String: params.FromAddress, Valid: true}
	}

	err = tx.Omit(clause.Associations).Save(p).Error
	if err != nil {
		tx.Rollback()
		return err
	}
	ub := &userbalance.UserBalance{}
	err = s.userBalanceService.GetBalanceOfUserByCoinUsingTx(tx, p.UserID, coin.ID, ub)
	if err != nil {
		tx.Rollback()
		return err
	}
	amountDecimal, err := decimal.NewFromString(p.Amount.String)
	if err != nil {
		tx.Rollback()
		return err
	}
	balanceAmountDecimal, err := decimal.NewFromString(ub.Amount)
	if err != nil {
		tx.Rollback()
		return err
	}
	balanceFrozenAmountDecimal, err := decimal.NewFromString(ub.FrozenAmount)
	if err != nil {
		tx.Rollback()
		return err
	}

	if params.Status == StatusCompleted {
		finalBalanceAmountDecimal := balanceAmountDecimal.Sub(amountDecimal)
		finalBalanceFrozenAmountDecimal := balanceFrozenAmountDecimal.Sub(amountDecimal)
		ub.Amount = finalBalanceAmountDecimal.StringFixed(8)
		ub.FrozenAmount = finalBalanceFrozenAmountDecimal.StringFixed(8)
		err = tx.Omit(clause.Associations).Save(ub).Error
		if err != nil {
			tx.Rollback()
			return err
		}
		withdrawTransaction := &transaction.Transaction{
			UserID:    ub.UserID,
			CoinID:    coin.ID,
			OrderID:   sql.NullInt64{Int64: 0, Valid: false},
			Type:      transaction.TypeWithdraw,
			Amount:    sql.NullString{String: amountDecimal.StringFixed(8), Valid: true},
			CoinName:  coin.Code,
			PaymentID: sql.NullInt64{Int64: p.ID, Valid: true},
		}
		err = tx.Omit(clause.Associations).Save(withdrawTransaction).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		feeDecimal, err := decimal.NewFromString(p.FeeAmount.String)
		if err != nil {
			tx.Rollback()
			return err
		}
		feeTransaction := &transaction.Transaction{
			UserID:    ub.UserID,
			CoinID:    coin.ID,
			OrderID:   sql.NullInt64{Int64: 0, Valid: false},
			Type:      transaction.TypeWithdrawFee,
			Amount:    sql.NullString{String: feeDecimal.StringFixed(8), Valid: true},
			CoinName:  coin.Code,
			PaymentID: sql.NullInt64{Int64: p.ID, Valid: true},
		}
		err = tx.Omit(clause.Associations).Save(feeTransaction).Error
		if err != nil {
			tx.Rollback()
			return err
		}

	} else {
		//failed status removing the frozen balance
		finalBalanceFrozenAmountDecimal := balanceFrozenAmountDecimal.Sub(amountDecimal)
		ub.FrozenAmount = finalBalanceFrozenAmountDecimal.StringFixed(8)
		err = tx.Omit(clause.Associations).Save(ub).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		return err
	}
	user := user.User{
		ID: ub.UserID,
	}
	if beforeStatus != p.Status {
		go s.notifyUserPaymentStatusUpdate(user, *p, "")
	}
	go s.publishPaymentToUser(user, *p)
	return nil
}

func (s *service) handleInternalTransferCallBack(internalTransferID, status string) error {
	id, err := strconv.ParseInt(internalTransferID, 10, 64)
	if err != nil {
		return err
	}
	return s.internalTransferService.UpdateStatus(id, status)
}

func (s *service) UpdateWithdraw(u *user.User, params UpdateWithdrawParams) (apiResponse response.APIResponse, statusCode int) {
	autoTransfer := true
	if params.AutoTransfer != nil {
		autoTransfer = *params.AutoTransfer
	}
	params.AdminStatus = strings.ToUpper(params.AdminStatus)
	params.Status = strings.ToUpper(params.Status)
	params.WithdrawType = strings.ToUpper(params.WithdrawType)
	if params.WithdrawType == "" {
		params.WithdrawType = WithdrawTypeHotWallet
	}
	tx := s.db.Begin()
	err := tx.Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("error starting transaction", err,
			zap.String("service", "paymentService"),
			zap.String("method", "UpdateWithdraw"),
			zap.Int64("paymentId", params.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}
	p := &Payment{}
	err = s.paymentRepository.GetPaymentByIDUsingTx(tx, params.ID, p)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		s.logger.Error2("error getting payment from db", err,
			zap.String("service", "paymentService"),
			zap.String("method", "UpdateWithdraw"),
			zap.Int64("paymentId", params.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || p.ID == 0 || p.Type == TypeDeposit {
		tx.Rollback()
		return response.Error("withdraw not found", http.StatusUnprocessableEntity, nil)
	}
	extraInfo := &ExtraInfo{}
	err = s.paymentRepository.GetExtraInfoByPaymentIDUsingTx(tx, p.ID, extraInfo)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		s.logger.Error2("error getting payment from db", err,
			zap.String("service", "paymentService"),
			zap.String("method", "UpdateWithdraw"),
			zap.Int64("paymentId", params.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || extraInfo.ID == 0 {
		tx.Rollback()
		return response.Error("extra info not found", http.StatusUnprocessableEntity, nil)
	}
	err = s.ValidateUpdateWithdraw(*p, *extraInfo, params)
	if err != nil {
		tx.Rollback()
		return response.Error(err.Error(), http.StatusUnprocessableEntity, nil)
	}

	formerStatus := p.Status

	//here updating the payment and extraInfo
	if params.Status != "" {
		p.Status = params.Status
	}
	if params.RejectionReason != "" {
		extraInfo.RejectionReason = sql.NullString{String: params.RejectionReason, Valid: true}
	}
	if params.AutoTransfer != nil {
		extraInfo.AutoTransfer = sql.NullBool{Bool: *params.AutoTransfer, Valid: true}
	}
	p.WithdrawType = sql.NullString{String: params.WithdrawType, Valid: true}
	if params.NetworkFee != "" {
		extraInfo.NetworkFee = sql.NullString{String: params.NetworkFee, Valid: true}
	}
	extraInfo.LastHandledID = sql.NullInt64{Int64: int64(u.ID), Valid: true}

	amountDecimal, _ := decimal.NewFromString(p.Amount.String)
	newFeeAmount := ""
	if params.Fee != "" {
		feeDecimal, err := decimal.NewFromString(params.Fee)
		if err != nil {
			return response.Error(err.Error(), http.StatusUnprocessableEntity, nil)
		}
		amountDecimal, _ := decimal.NewFromString(p.Amount.String)
		if !amountDecimal.Sub(feeDecimal).IsPositive() {
			return response.Error("fee must be less than amount", http.StatusUnprocessableEntity, nil)
		}
		newFeeAmount = feeDecimal.StringFixed(8)
		p.FeeAmount = sql.NullString{String: newFeeAmount, Valid: true}
	}
	if params.AdminStatus != "" {
		p.AdminStatus = sql.NullString{String: params.AdminStatus, Valid: true}
		if err != nil {
			tx.Rollback()
			s.logger.Error2("error updating payment", err,
				zap.String("service", "paymentService"),
				zap.String("method", "UpdateWithdraw"),
				zap.Int64("paymentId", params.ID),
				zap.String("adminStatus", params.AdminStatus),
				zap.String("newFee", params.Fee),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}
	} else if params.Status != "" {
		shouldSaveBalance := false
		ub := &userbalance.UserBalance{}
		err = s.userBalanceService.GetBalanceOfUserByCoinUsingTx(tx, p.UserID, p.CoinID, ub)
		if err != nil {
			tx.Rollback()
			s.logger.Error2("error getting user balance", err,
				zap.String("service", "paymentService"),
				zap.String("method", "UpdateWithdraw"),
				zap.Int64("paymentId", params.ID),
				zap.Int("userId", p.UserID),
				zap.Int64("coinId", p.CoinID),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}
		balanceAmountDecimal, err := decimal.NewFromString(ub.Amount)
		if err != nil {
			tx.Rollback()
			s.logger.Error2("error getting user balance", err,
				zap.String("service", "paymentService"),
				zap.String("method", "UpdateWithdraw"),
				zap.Int64("paymentId", params.ID),
				zap.Int("userId", p.UserID),
				zap.Int64("coinId", p.CoinID),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}
		balanceFrozenAmountDecimal, err := decimal.NewFromString(ub.FrozenAmount)
		if err != nil {
			tx.Rollback()
			s.logger.Error2("error getting user balance", err,
				zap.String("service", "paymentService"),
				zap.String("method", "UpdateWithdraw"),
				zap.Int64("paymentId", params.ID),
				zap.Int("userId", p.UserID),
				zap.Int64("coinId", p.CoinID),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}

		if formerStatus == StatusFailed || formerStatus == StatusRejected || formerStatus == StatusCanceled {
			remainingBalanceDecimal := balanceAmountDecimal.Sub(balanceFrozenAmountDecimal)
			if remainingBalanceDecimal.Sub(amountDecimal).IsNegative() {
				tx.Rollback()
				return response.Error("user has not enough balance to do this action", http.StatusUnprocessableEntity, nil)
			}
		}
		if params.Status == StatusInProgress && autoTransfer {
			finalFee := p.FeeAmount.String
			if newFeeAmount != "" {
				finalFee = newFeeAmount
			}
			finalFeeDecimal, _ := decimal.NewFromString(finalFee)
			amountWithoutFeeDecimal := amountDecimal.Sub(finalFeeDecimal)
			amountWithoutFee := amountWithoutFeeDecimal.StringFixed(8)

			if params.WithdrawType == WithdrawTypeHotWallet {
				txID, err := s.walletService.SendTransaction(p.Code, amountWithoutFee, p.ToAddress.String, p.BlockchainNetwork.String, params.NetworkFee)
				if err != nil {
					tx.Rollback()
					s.logger.Error2("error sending tx to wallet", err,
						zap.String("service", "paymentService"),
						zap.String("method", "UpdateWithdraw"),
						zap.Int64("paymentId", params.ID),
						zap.String("amount", p.Amount.String),
						zap.String("ToAddress", p.ToAddress.String),
						zap.String("network", p.BlockchainNetwork.String),
						zap.String("networkFee", params.NetworkFee),
					)
					return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)

				}
				p.TxID = sql.NullString{String: txID, Valid: true}

			} else {
				withdrawParams := externalexchange.WithdrawParams{
					Coin:      p.Code,
					Amount:    p.Amount.String, //we send amount including fee, because external exchange (here binance) would reduce that fee
					ToAddress: p.ToAddress.String,
					Network:   p.BlockchainNetwork.String,
				}

				withdrawResult, err := s.externalExchangeService.Withdraw(withdrawParams)
				if err != nil {
					tx.Rollback()
					s.logger.Error2("error withdraw from external exchange", err,
						zap.String("service", "paymentService"),
						zap.String("method", "UpdateWithdraw"),
						zap.Int64("paymentId", params.ID),
						zap.String("amount", p.Amount.String),
						zap.String("ToAddress", p.ToAddress.String),
						zap.String("network", p.BlockchainNetwork.String),
						zap.String("networkFee", params.NetworkFee),
					)
					return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
				}
				extraInfo.ExternalExchangeID = sql.NullInt64{Int64: withdrawResult.ExternalExchangeID, Valid: true}
				extraInfo.ExternalExchangeWithdrawID = sql.NullString{String: withdrawResult.ID, Valid: true}
			}
		}

		if params.Status == StatusFailed || params.Status == StatusRejected || params.Status == StatusCanceled {
			shouldSaveBalance = true
			finalBalanceFrozenAmountDecimal := balanceFrozenAmountDecimal.Sub(amountDecimal)
			if finalBalanceFrozenAmountDecimal.IsNegative() {
				tx.Rollback()
				s.logger.Error2("frozen balance would become negative", err,
					zap.String("service", "paymentService"),
					zap.String("method", "UpdateWithdraw"),
					zap.Int64("paymentId", params.ID),
					zap.String("status", params.Status),
				)
				return response.Error("frozen balance would be negative", http.StatusUnprocessableEntity, nil)
			}
			ub.FrozenAmount = finalBalanceFrozenAmountDecimal.StringFixed(8)
		}
		//if autoTransfer is false then it means we pay to user in a way out of blockchain ask mr jabbari the cases that we may have?!
		if params.Status == StatusInProgress && !autoTransfer {
			shouldSaveBalance = true
			//this case means the status is complated
			p.Status = StatusCompleted
			finalBalanceAmountDecimal := balanceAmountDecimal.Sub(amountDecimal)
			if finalBalanceAmountDecimal.IsNegative() {
				tx.Rollback()
				s.logger.Error2("balance would become negative", err,
					zap.String("service", "paymentService"),
					zap.String("method", "UpdateWithdraw"),
					zap.Int64("paymentId", params.ID),
					zap.String("status", params.Status),
				)
				return response.Error("frozen balance would be negative", http.StatusUnprocessableEntity, nil)
			}
			ub.Amount = finalBalanceAmountDecimal.StringFixed(8)
			/*
			  if former status of payment is one of /FAILED/REJECTED/CANCLED and admin wants to update it and
			  set to to inProgress then we do not remove the user frozen balance because at the previus status
			  changed to /FAILED/REJECTED/CANCLED the fronzen balance is removed
			*/
			if formerStatus != StatusFailed && formerStatus != StatusRejected && formerStatus != StatusCanceled {
				finalBalanceFrozenAmountDecimal := balanceFrozenAmountDecimal.Sub(amountDecimal)
				ub.FrozenAmount = finalBalanceFrozenAmountDecimal.StringFixed(8)
			}

			//set transaction
			withdrawTransaction := &transaction.Transaction{
				UserID:    p.UserID,
				CoinID:    p.CoinID,
				OrderID:   sql.NullInt64{Int64: 0, Valid: false},
				Type:      transaction.TypeWithdraw,
				Amount:    sql.NullString{String: amountDecimal.StringFixed(8), Valid: true},
				CoinName:  p.Code,
				PaymentID: sql.NullInt64{Int64: p.ID, Valid: true},
			}
			err = tx.Omit(clause.Associations).Save(withdrawTransaction).Error
			if err != nil {
				tx.Rollback()
				s.logger.Error2("error saving withdraw transaction", err,
					zap.String("service", "paymentService"),
					zap.String("method", "UpdateWithdraw"),
					zap.Int64("paymentId", params.ID),
				)
				return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
			}

			feeDecimal, err := decimal.NewFromString(p.FeeAmount.String)
			if err != nil {
				tx.Rollback()
				s.logger.Error2("error calculating fee", err,
					zap.String("service", "paymentService"),
					zap.String("method", "UpdateWithdraw"),
					zap.Int64("paymentId", params.ID),
				)
				return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
			}
			feeTransaction := &transaction.Transaction{
				UserID:    ub.UserID,
				CoinID:    p.CoinID,
				OrderID:   sql.NullInt64{Int64: 0, Valid: false},
				Type:      transaction.TypeWithdrawFee,
				Amount:    sql.NullString{String: feeDecimal.StringFixed(8), Valid: true},
				CoinName:  p.Code,
				PaymentID: sql.NullInt64{Int64: p.ID, Valid: true},
			}
			err = tx.Omit(clause.Associations).Save(feeTransaction).Error
			if err != nil {
				tx.Rollback()
				s.logger.Error2("error saving fee transaction", err,
					zap.String("service", "paymentService"),
					zap.String("method", "UpdateWithdraw"),
					zap.Int64("paymentId", params.ID),
				)
				return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
			}

		}
		if shouldSaveBalance {
			err = tx.Omit(clause.Associations).Save(ub).Error
			if err != nil {
				tx.Rollback()
				s.logger.Error2("error updating userBalance", err,
					zap.String("service", "paymentService"),
					zap.String("method", "UpdateWithdraw"),
					zap.Int64("paymentId", params.ID),
				)
				return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
			}

		}
	}
	err = tx.Omit(clause.Associations).Save(p).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("error updating payment", err,
			zap.String("service", "paymentService"),
			zap.String("method", "UpdateWithdraw"),
			zap.Int64("paymentId", params.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}
	err = tx.Omit(clause.Associations).Save(extraInfo).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("error updating extraInfo", err,
			zap.String("service", "paymentService"),
			zap.String("method", "UpdateWithdraw"),
			zap.Int64("paymentId", params.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("error saving userbalance", err,
			zap.String("service", "paymentService"),
			zap.String("method", "UpdateWithdraw"),
			zap.Int64("paymentId", params.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if formerStatus != p.Status {
		user := user.User{ID: p.UserID}
		go s.notifyUserPaymentStatusUpdate(user, *p, params.RejectionReason)
		go s.publishPaymentToUser(user, *p)
	}
	res := make(map[string]interface{})
	return response.Success(res, "")
}

func (s *service) ValidateUpdateWithdraw(p Payment, extraInfo ExtraInfo, params UpdateWithdrawParams) error {
	if params.Status == "" && params.AdminStatus == "" {
		return fmt.Errorf("at least one of the status or admin status should be provided")
	}
	if params.Status != "" && params.AdminStatus != "" {
		return fmt.Errorf("only one of the status or admin status can be set")
	}
	formerStatus := p.Status
	formerAdminStatus := p.AdminStatus.String
	if (params.Status != "" && params.Status == formerStatus) || (params.AdminStatus != "" && params.AdminStatus == formerAdminStatus) {
		return fmt.Errorf("provided status is the same as payment status")
	}
	if formerStatus == StatusCompleted || formerStatus == StatusInProgress || formerStatus == StatusUserCanceled {
		return fmt.Errorf("payment status is already %s", formerStatus)
	}
	if params.Status == StatusRejected && strings.Trim(params.RejectionReason, "") == "" {
		return fmt.Errorf("rejection reason must be provided")
	}
	return nil
}

func (s *service) UpdateDeposit(u *user.User, params UpdateDepositParams) (apiResponse response.APIResponse, statusCode int) {
	shouldDeposit := false
	if params.ShouldDeposit != nil {
		shouldDeposit = *params.ShouldDeposit
	}
	params.Status = strings.ToUpper(params.Status)
	tx := s.db.Begin()
	err := tx.Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("error starting transaction", err,
			zap.String("service", "paymentService"),
			zap.String("method", "UpdateDeposit"),
			zap.Int64("paymentId", params.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}
	p := &Payment{}
	err = s.paymentRepository.GetPaymentByIDUsingTx(tx, params.ID, p)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		s.logger.Error2("error getting payment from db", err,
			zap.String("service", "paymentService"),
			zap.String("method", "UpdateDeposit"),
			zap.Int64("paymentId", params.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || p.ID == 0 || p.Type == TypeWithdraw {
		tx.Rollback()
		return response.Error("deposit not found", http.StatusUnprocessableEntity, nil)
	}
	extraInfo := &ExtraInfo{}
	err = s.paymentRepository.GetExtraInfoByPaymentIDUsingTx(tx, p.ID, extraInfo)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		s.logger.Error2("error getting payment from db", err,
			zap.String("service", "paymentService"),
			zap.String("method", "UpdateDeposit"),
			zap.Int64("paymentId", params.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || extraInfo.ID == 0 {
		tx.Rollback()
		return response.Error("extra info not found", http.StatusUnprocessableEntity, nil)
	}

	if shouldDeposit && p.Status == StatusCompleted {
		return response.Error("the deposit is already completed, can not deposit again", http.StatusUnprocessableEntity, nil)
	}
	if params.Status != "" && params.Status != StatusCompleted {
		p.Status = params.Status
	}
	if params.FromAddress != "" {
		p.FromAddress = sql.NullString{String: params.FromAddress, Valid: true}
	}
	if params.ToAddress != "" {
		p.ToAddress = sql.NullString{String: params.ToAddress, Valid: true}
	}
	if params.TxID != "" {
		p.TxID = sql.NullString{String: params.TxID, Valid: true}
	}
	if params.Amount != "" {
		amountDecimal, err := decimal.NewFromString(params.Amount)
		if err != nil {
			return response.Error("the amount is not in right format", http.StatusUnprocessableEntity, nil)
		}
		p.Amount = sql.NullString{String: amountDecimal.StringFixed(8), Valid: true}
	}
	if shouldDeposit {
		//if we should depoist then the status is completed
		p.Status = StatusCompleted
		ub := &userbalance.UserBalance{}
		err = s.userBalanceService.GetBalanceOfUserByCoinUsingTx(tx, p.UserID, p.CoinID, ub)
		if err != nil {
			tx.Rollback()
			s.logger.Error2("error getting user balance", err,
				zap.String("service", "paymentService"),
				zap.String("method", "UpdateDeposit"),
				zap.Int64("paymentId", params.ID),
				zap.Int("userId", p.UserID),
				zap.Int64("coinId", p.CoinID),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}
		balanceAmountDecimal, err := decimal.NewFromString(ub.Amount)
		if err != nil {
			tx.Rollback()
			s.logger.Error2("error getting user balance", err,
				zap.String("service", "paymentService"),
				zap.String("method", "UpdateDeposit"),
				zap.Int64("paymentId", params.ID),
				zap.Int("userId", p.UserID),
				zap.Int64("coinId", p.CoinID),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}

		amountDecimal, err := decimal.NewFromString(p.Amount.String)
		if err != nil {
			tx.Rollback()
			s.logger.Error2("error getting user balance", err,
				zap.String("service", "paymentService"),
				zap.String("method", "UpdateDeposit"),
				zap.Int64("paymentId", params.ID),
				zap.Int("userId", p.UserID),
				zap.Int64("coinId", p.CoinID),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}
		finalBalanceAmountDecimal := balanceAmountDecimal.Add(amountDecimal)
		ub.Amount = finalBalanceAmountDecimal.StringFixed(8)
		err = tx.Omit(clause.Associations).Save(ub).Error
		if err != nil {
			tx.Rollback()
			s.logger.Error2("error updating user balance", err,
				zap.String("service", "paymentService"),
				zap.String("method", "UpdateDeposit"),
				zap.Int64("paymentId", params.ID),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}
		//set transaction
		depositTransaction := &transaction.Transaction{
			UserID:    p.UserID,
			CoinID:    p.CoinID,
			OrderID:   sql.NullInt64{Int64: 0, Valid: false},
			Type:      transaction.TypeDeposit,
			Amount:    sql.NullString{String: amountDecimal.StringFixed(8), Valid: true},
			CoinName:  p.Code,
			PaymentID: sql.NullInt64{Int64: p.ID, Valid: true},
		}
		err = tx.Omit(clause.Associations).Save(depositTransaction).Error
		if err != nil {
			tx.Rollback()
			s.logger.Error2("error saving depositTransaction", err,
				zap.String("service", "paymentService"),
				zap.String("method", "UpdateDeposit"),
				zap.Int64("paymentId", params.ID),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}

	}
	extraInfo.LastHandledID = sql.NullInt64{Int64: int64(u.ID), Valid: true}

	err = tx.Omit(clause.Associations).Save(p).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("error updating payment", err,
			zap.String("service", "paymentService"),
			zap.String("method", "UpdateDeposit"),
			zap.Int64("paymentId", params.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}
	err = tx.Omit(clause.Associations).Save(extraInfo).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("error updating extraInfo", err,
			zap.String("service", "paymentService"),
			zap.String("method", "UpdateDeposit"),
			zap.Int64("paymentId", params.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("error committing db transaction", err,
			zap.String("service", "paymentService"),
			zap.String("method", "UpdateWithdraw"),
			zap.Int64("paymentId", params.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	res := make(map[string]interface{})
	return response.Success(res, "")
}

func NewPaymentService(db *gorm.DB, paymentRepository Repository, currencyService currency.Service, walletService wallet.Service,
	userConfigService user.ConfigService, twoFaManager user.TwoFaManager, withdrawEmailConfirmationManager WithdrawEmailConfirmationManager,
	permissionManager user.PermissionManager, userWithdrawAddressService userwithdrawaddress.Service, userService user.Service,
	userBalanceService userbalance.Service, communicationService communication.Service, priceGenerator currency.PriceGenerator,
	internalTransferService InternalTransferService, externalExchangeService externalexchange.Service,
	autoExchangeManager AutoExchangeManager, mqttManager communication.MqttManager, configs platform.Configs, logger platform.Logger) Service {
	return &service{
		db:                               db,
		paymentRepository:                paymentRepository,
		currencyService:                  currencyService,
		walletService:                    walletService,
		userConfigService:                userConfigService,
		twoFaManager:                     twoFaManager,
		withdrawEmailConfirmationManager: withdrawEmailConfirmationManager,
		permissionManager:                permissionManager,
		userService:                      userService,
		userBalanceService:               userBalanceService,
		userWithdrawAddressService:       userWithdrawAddressService,
		communicationService:             communicationService,
		priceGenerator:                   priceGenerator,
		internalTransferService:          internalTransferService,
		externalExchangeService:          externalExchangeService,
		autoExchangeManager:              autoExchangeManager,
		mqttManager:                      mqttManager,
		configs:                          configs,
		logger:                           logger,
	}
}
