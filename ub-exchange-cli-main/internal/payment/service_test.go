// Package payment_test tests the payment Service covering end-to-end
// payment operations. Covers:
//   - GetPayments retrieval and Detail lookup (including not-found)
//   - PreWithdraw validation: coin not found, invalid network coin, invalid address,
//     wrong amount, unverified account, disabled 2FA, permission denied, min/max
//     withdraw limits, read-only account, whitelist restrictions, insufficient balance,
//     and successful flows with email code and 2FA combinations
//   - Withdraw execution: same validations as PreWithdraw plus wrong 2FA code,
//     wrong email code, and successful withdrawal with various user config states
//   - GetInProgressWithdrawalsInExternalExchange and UpdatePaymentInExternalExchange
//     for completed and failed statuses
//   - CancelWithdraw: payment not found, not a withdraw, wrong user, wrong status,
//     and successful cancellation with balance restoration
//   - HandleWalletCallback: internal transfers, deposit creation and completion,
//     existing deposit updates, and withdraw status callbacks (failed/completed)
//   - UpdateWithdraw admin operations: status/fee/network-fee updates, rejection
//     with balance refund, in-progress with hot wallet, external exchange wallet,
//     and non-auto-transfer flows
//   - UpdateDeposit: conditional deposit processing based on confirmation thresholds
//
// Test data: mocked PaymentRepository, CurrencyService, WalletService,
// UserConfigService, TwoFaManager, WithdrawEmailConfirmationManager,
// UserPermissionManager, UserWithdrawAddressService, UserService,
// UserBalanceService, CommunicationService, PriceGenerator,
// InternalTransferService, ExternalExchangeService, AutoExchangeManager,
// CentrifugoManager, and Configs; go-sqlmock for GORM database interactions.
package payment_test

import (
	"database/sql"
	"exchange-go/internal/currency"
	"exchange-go/internal/externalexchange"
	"exchange-go/internal/mocks"
	"exchange-go/internal/payment"
	"exchange-go/internal/user"
	"exchange-go/internal/userbalance"
	"exchange-go/internal/userwithdrawaddress"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestService_GetPayments(t *testing.T) {
	db := &gorm.DB{}
	coin := currency.Coin{
		ID:   1,
		Name: "",
		Code: "BTC",
	}

	payments := []payment.Payment{
		{
			ID:        1,
			CoinID:    1,
			Coin:      coin,
			Type:      payment.TypeWithdraw,
			Status:    payment.StatusCompleted,
			ToAddress: sql.NullString{String: "btcAddress1", Valid: true},
			Amount:    sql.NullString{String: "0.1", Valid: true},
			TxID:      sql.NullString{String: "btcTxId1", Valid: true},
		},
		{
			ID:        2,
			CoinID:    1,
			Coin:      coin,
			Type:      payment.TypeDeposit,
			Status:    payment.StatusInProgress,
			ToAddress: sql.NullString{String: "btcAddress2", Valid: true},
			Amount:    sql.NullString{String: "0.01", Valid: true},
			TxID:      sql.NullString{String: "btcTxId2", Valid: true},
		},
		{
			ID:        3,
			CoinID:    1,
			Coin:      coin,
			Type:      payment.TypeWithdraw,
			Status:    payment.StatusFailed,
			ToAddress: sql.NullString{String: "btcAddress3", Valid: true},
			Amount:    sql.NullString{String: "0.05", Valid: true},
			TxID:      sql.NullString{String: "btcTxId3", Valid: true},
		},
	}
	paymentRepo := new(mocks.PaymentRepository)
	paymentRepo.On("GetUserPayments", mock.Anything).Once().Return(payments)
	currencyService := new(mocks.CurrencyService)
	walletService := new(mocks.WalletService)
	userConfigService := new(mocks.UserConfigService)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID: 21,
	}

	params := payment.GetPaymentsParams{}

	res, statusCode := paymentService.GetPayments(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	resultData, ok := res.Data.(map[string][]payment.GetPaymentsResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}

	paymentResults, ok := resultData["payments"]
	if !ok {
		t.Error("payments key does not exists in response")
	}

	assert.Equal(t, 3, len(paymentResults))

	p1 := paymentResults[0]
	assert.Equal(t, int64(1), p1.ID)
	assert.Equal(t, "completed", p1.Status)
	assert.Equal(t, "withdraw", p1.Type)
	assert.Equal(t, "0.1", p1.Amount)
	assert.Equal(t, "BTC", p1.Coin)
	assert.Equal(t, "btcAddress1", p1.Address)
	assert.Equal(t, "btcTxId1", p1.TxID)

	p2 := paymentResults[1]
	assert.Equal(t, int64(2), p2.ID)
	assert.Equal(t, "in progress", p2.Status)
	assert.Equal(t, "deposit", p2.Type)
	assert.Equal(t, "0.01", p2.Amount)
	assert.Equal(t, "BTC", p2.Coin)
	assert.Equal(t, "btcAddress2", p2.Address)
	assert.Equal(t, "btcTxId2", p2.TxID)

	p3 := paymentResults[2]

	assert.Equal(t, int64(3), p3.ID)
	assert.Equal(t, "failed", p3.Status)
	assert.Equal(t, "withdraw", p3.Type)
	assert.Equal(t, "0.05", p3.Amount)
	assert.Equal(t, "BTC", p3.Coin)
	assert.Equal(t, "btcAddress3", p3.Address)
	assert.Equal(t, "btcTxId3", p3.TxID)
	paymentRepo.AssertExpectations(t)

}

func TestService_Detail_Fail_NotFound(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	paymentRepo.On("GetPaymentDetailByID", int64(1)).Once().Return(payment.DetailQueryFields{}, gorm.ErrRecordNotFound)
	currencyService := new(mocks.CurrencyService)
	walletService := new(mocks.WalletService)
	userConfigService := new(mocks.UserConfigService)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)
	u := user.User{
		ID: 21,
	}

	params := payment.GetPaymentDetailParams{
		ID: 1,
	}

	res, statusCode := paymentService.Detail(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "payment not found", res.Message)
	paymentRepo.AssertExpectations(t)
}

func TestService_Detail(t *testing.T) {
	db := &gorm.DB{}
	dqf := payment.DetailQueryFields{
		ID:              1,
		Code:            "BTC",
		Network:         "",
		Address:         "btcAddress1",
		UserID:          21,
		TxID:            "btcTxId1",
		RejectionReason: "",
	}
	paymentRepo := new(mocks.PaymentRepository)
	paymentRepo.On("GetPaymentDetailByID", int64(1)).Once().Return(dqf, nil)
	currencyService := new(mocks.CurrencyService)
	walletService := new(mocks.WalletService)
	userConfigService := new(mocks.UserConfigService)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID: 21,
	}

	params := payment.GetPaymentDetailParams{
		ID: 1,
	}

	res, statusCode := paymentService.Detail(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)
	result, ok := res.Data.(payment.DetailResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}

	assert.Equal(t, "btcAddress1", result.Address)
	assert.Equal(t, "btcTxId1", result.TxID)
	assert.Equal(t, "", result.RejectionReason)
	paymentRepo.AssertExpectations(t)

}

func TestService_PreWithdraw_Fail_CoinNotFound(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "BTCQ").Once().Return(currency.Coin{}, gorm.ErrRecordNotFound)
	walletService := new(mocks.WalletService)
	userConfigService := new(mocks.UserConfigService)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID: 21,
	}

	params := payment.PreWithdrawParams{
		Coin:    "BTCQ",
		Amount:  "0.1",
		Address: "btcAddress1",
		Network: "",
	}

	res, statusCode := paymentService.PreWithdraw(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "coin not found", res.Message)
	currencyService.AssertExpectations(t)

}

func TestService_PreWithdraw_Fail_NetworkCoinNotFound(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "USDT").Once().Return(currency.Coin{ID: 1}, nil)
	currencyService.On("GetCoinByCode", "ETHQ").Once().Return(currency.Coin{}, gorm.ErrRecordNotFound)
	walletService := new(mocks.WalletService)
	userConfigService := new(mocks.UserConfigService)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID: 21,
	}

	params := payment.PreWithdrawParams{
		Coin:    "usdt",
		Amount:  "0.1",
		Address: "ethAddress1",
		Network: "ethq",
	}

	res, statusCode := paymentService.PreWithdraw(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "network not found", res.Message)
	currencyService.AssertExpectations(t)

}

func TestService_PreWithdraw_Fail_AddressIsNotValid(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "USDT").Once().Return(currency.Coin{ID: 1, Code: "USDT"}, nil)
	currencyService.On("GetCoinByCode", "ETH").Once().Return(currency.Coin{ID: 3}, nil)
	walletService := new(mocks.WalletService)
	walletService.On("IsAddressValid", "USDT", "ethAddress1", "ETH").Once().Return(false, nil)
	userConfigService := new(mocks.UserConfigService)
	//userConfigService.On("GetUserConfig", int64(21)).Once().Return(user.Config{}, gorm.ErrRecordNotFound)

	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID:     21,
		Status: user.StatusVerified,
	}

	params := payment.PreWithdrawParams{
		Coin:    "usdt",
		Amount:  "0.10",
		Address: "ethAddress1",
		Network: "eth",
	}

	res, statusCode := paymentService.PreWithdraw(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "address is not valid", res.Message)

	currencyService.AssertExpectations(t)
	walletService.AssertExpectations(t)

}

func TestService_PreWithdraw_Fail_WrongAmount(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "USDT").Once().Return(currency.Coin{ID: 1, Code: "USDT"}, nil)
	currencyService.On("GetCoinByCode", "ETH").Once().Return(currency.Coin{ID: 3}, nil)
	walletService := new(mocks.WalletService)
	walletService.On("IsAddressValid", "USDT", "ethAddress1", "ETH").Once().Return(true, nil)
	userConfigService := new(mocks.UserConfigService)
	//userConfigService.On("GetUserConfig", int64(21)).Once().Return(user.Config{}, gorm.ErrRecordNotFound)

	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID:     21,
		Status: user.StatusVerified,
	}

	params := payment.PreWithdrawParams{
		Coin:    "usdt",
		Amount:  "-0.10",
		Address: "ethAddress1",
		Network: "eth",
	}

	res, statusCode := paymentService.PreWithdraw(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "amount is not correct", res.Message)

	currencyService.AssertExpectations(t)
	walletService.AssertExpectations(t)

}

func TestService_PreWithdraw_Fail_UserAccountStatusIsNotVerified(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "USDT").Once().Return(currency.Coin{ID: 1, Code: "USDT"}, nil)
	currencyService.On("GetCoinByCode", "ETH").Once().Return(currency.Coin{ID: 3}, nil)
	walletService := new(mocks.WalletService)
	walletService.On("IsAddressValid", "USDT", "ethAddress1", "ETH").Once().Return(true, nil)
	userConfigService := new(mocks.UserConfigService)
	userConfigService.On("GetUserConfig", 21).Once().Return(user.Config{}, gorm.ErrRecordNotFound)

	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID:     21,
		Status: user.StatusRegistered,
	}

	params := payment.PreWithdrawParams{
		Coin:    "usdt",
		Amount:  "0.10",
		Address: "ethAddress1",
		Network: "eth",
	}

	res, statusCode := paymentService.PreWithdraw(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "user account is not verified", res.Message)

	currencyService.AssertExpectations(t)
	walletService.AssertExpectations(t)
	userConfigService.AssertExpectations(t)

}

func TestService_PreWithdraw_Fail_UserDisabledTwoFa(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "USDT").Once().Return(currency.Coin{ID: 1, Code: "USDT"}, nil)
	currencyService.On("GetCoinByCode", "ETH").Once().Return(currency.Coin{ID: 3}, nil)
	walletService := new(mocks.WalletService)
	walletService.On("IsAddressValid", "USDT", "ethAddress1", "ETH").Once().Return(true, nil)
	userConfigService := new(mocks.UserConfigService)
	userConfigService.On("GetUserConfig", 21).Once().Return(user.Config{}, gorm.ErrRecordNotFound)

	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID:                  21,
		Status:              user.StatusVerified,
		Google2faDisabledAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	params := payment.PreWithdrawParams{
		Coin:    "usdt",
		Amount:  "0.10",
		Address: "ethAddress1",
		Network: "eth",
	}

	res, statusCode := paymentService.PreWithdraw(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "for the security reasons, after disabling/enabling 2fa the withdraw request is not allowed for 24 hours", res.Message)

	currencyService.AssertExpectations(t)
	walletService.AssertExpectations(t)
	userConfigService.AssertExpectations(t)

}

func TestService_PreWithdraw_Fail_PermissionNotGranted(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "USDT").Once().Return(currency.Coin{ID: 1, Code: "USDT"}, nil)
	currencyService.On("GetCoinByCode", "ETH").Once().Return(currency.Coin{ID: 3}, nil)
	walletService := new(mocks.WalletService)
	walletService.On("IsAddressValid", "USDT", "ethAddress1", "ETH").Once().Return(true, nil)
	userConfigService := new(mocks.UserConfigService)
	userConfigService.On("GetUserConfig", 21).Once().Return(user.Config{}, gorm.ErrRecordNotFound)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	permissionManager.On("IsPermissionGrantedToUserFor", mock.Anything, user.PermissionWithdraw).Once().Return(false)

	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID:     21,
		Status: user.StatusVerified,
	}

	params := payment.PreWithdrawParams{
		Coin:    "usdt",
		Amount:  "0.10",
		Address: "ethAddress1",
		Network: "eth",
	}

	res, statusCode := paymentService.PreWithdraw(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "withdraw permission is not granted", res.Message)

	currencyService.AssertExpectations(t)
	walletService.AssertExpectations(t)
	userConfigService.AssertExpectations(t)
	permissionManager.AssertExpectations(t)

}

func TestService_PreWithdraw_Fail_LessThanMinimumWithdraw(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	usdtCoin := currency.Coin{
		ID:              1,
		Code:            "USDT",
		MinimumWithdraw: "10",
	}
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "USDT").Once().Return(usdtCoin, nil)
	currencyService.On("GetCoinByCode", "ETH").Once().Return(currency.Coin{ID: 3}, nil)
	walletService := new(mocks.WalletService)
	walletService.On("IsAddressValid", "USDT", "ethAddress1", "ETH").Once().Return(true, nil)
	userConfigService := new(mocks.UserConfigService)
	userConfigService.On("GetUserConfig", 21).Once().Return(user.Config{}, gorm.ErrRecordNotFound)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	permissionManager.On("IsPermissionGrantedToUserFor", mock.Anything, user.PermissionWithdraw).Once().Return(true)

	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID:     21,
		Status: user.StatusVerified,
	}

	params := payment.PreWithdrawParams{
		Coin:    "usdt",
		Amount:  "0.10",
		Address: "ethAddress1",
		Network: "eth",
	}

	res, statusCode := paymentService.PreWithdraw(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "minimum withdraw is: 10", res.Message)

	currencyService.AssertExpectations(t)
	walletService.AssertExpectations(t)
	userConfigService.AssertExpectations(t)
	permissionManager.AssertExpectations(t)

}

func TestService_PreWithdraw_Fail_MoreThanMaximumWithdraw(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	usdtCoin := currency.Coin{
		ID:              1,
		Code:            "USDT",
		MinimumWithdraw: "10",
		MaximumWithdraw: "1000",
	}
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "USDT").Once().Return(usdtCoin, nil)
	currencyService.On("GetCoinByCode", "ETH").Once().Return(currency.Coin{ID: 3}, nil)
	walletService := new(mocks.WalletService)
	walletService.On("IsAddressValid", "USDT", "ethAddress1", "ETH").Once().Return(true, nil)
	userConfigService := new(mocks.UserConfigService)
	userConfigService.On("GetUserConfig", 21).Once().Return(user.Config{}, gorm.ErrRecordNotFound)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	permissionManager.On("IsPermissionGrantedToUserFor", mock.Anything, user.PermissionWithdraw).Once().Return(true)

	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID:     21,
		Status: user.StatusVerified,
	}

	params := payment.PreWithdrawParams{
		Coin:    "usdt",
		Amount:  "1100.00",
		Address: "ethAddress1",
		Network: "eth",
	}

	res, statusCode := paymentService.PreWithdraw(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "maximum withdraw is: 1000", res.Message)

	currencyService.AssertExpectations(t)
	walletService.AssertExpectations(t)
	userConfigService.AssertExpectations(t)
	permissionManager.AssertExpectations(t)

}

func TestService_PreWithdraw_Fail_AccountIsReadOnly(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	usdtCoin := currency.Coin{
		ID:              1,
		Code:            "USDT",
		MinimumWithdraw: "10",
		MaximumWithdraw: "100",
	}
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "USDT").Once().Return(usdtCoin, nil)
	currencyService.On("GetCoinByCode", "ETH").Once().Return(currency.Coin{ID: 3}, nil)
	walletService := new(mocks.WalletService)
	walletService.On("IsAddressValid", "USDT", "ethAddress1", "ETH").Once().Return(true, nil)

	uc := user.Config{
		ID:         1,
		IsReadOnly: true,
	}
	userConfigService := new(mocks.UserConfigService)
	userConfigService.On("GetUserConfig", 21).Once().Return(uc, nil)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)

	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID:     21,
		Status: user.StatusVerified,
	}

	params := payment.PreWithdrawParams{
		Coin:    "usdt",
		Amount:  "10.00",
		Address: "ethAddress1",
		Network: "eth",
	}

	res, statusCode := paymentService.PreWithdraw(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "this account is in read only mode", res.Message)

	currencyService.AssertExpectations(t)
	walletService.AssertExpectations(t)
	userConfigService.AssertExpectations(t)

}

func TestService_PreWithdraw_Fail_WhitelistEnabled(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	usdtCoin := currency.Coin{
		ID:              1,
		Code:            "USDT",
		MinimumWithdraw: "10",
		MaximumWithdraw: "100",
	}
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "USDT").Once().Return(usdtCoin, nil)
	currencyService.On("GetCoinByCode", "ETH").Once().Return(currency.Coin{ID: 3}, nil)
	walletService := new(mocks.WalletService)
	walletService.On("IsAddressValid", "USDT", "ethAddress1", "ETH").Once().Return(true, nil)

	uc := user.Config{
		ID:                 1,
		IsReadOnly:         false,
		IsWhiteListEnabled: true,
	}
	userConfigService := new(mocks.UserConfigService)
	userConfigService.On("GetUserConfig", 21).Once().Return(uc, nil)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)

	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userWithdrawAddressService.On("GetUserWithdrawAddressesByAddress", mock.Anything, mock.Anything, mock.Anything).Once().Return([]userwithdrawaddress.UserWithdrawAddress{})
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID:     21,
		Status: user.StatusVerified,
	}

	params := payment.PreWithdrawParams{
		Coin:    "usdt",
		Amount:  "10.00",
		Address: "ethAddress1",
		Network: "eth",
	}

	res, statusCode := paymentService.PreWithdraw(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "this address is not in white list", res.Message)

	currencyService.AssertExpectations(t)
	walletService.AssertExpectations(t)
	userConfigService.AssertExpectations(t)
	userWithdrawAddressService.AssertExpectations(t)

}

func TestService_PreWithdraw_Fail_UserBalanceIsNotEnough(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	usdtCoin := currency.Coin{
		ID:              1,
		Code:            "USDT",
		MinimumWithdraw: "10",
		MaximumWithdraw: "10000",
	}
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "USDT").Once().Return(usdtCoin, nil)
	currencyService.On("GetCoinByCode", "ETH").Once().Return(currency.Coin{ID: 3}, nil)
	walletService := new(mocks.WalletService)
	walletService.On("IsAddressValid", "USDT", "ethAddress1", "ETH").Once().Return(true, nil)

	uc := user.Config{
		ID:                 1,
		IsReadOnly:         false,
		IsWhiteListEnabled: false,
	}
	userConfigService := new(mocks.UserConfigService)
	userConfigService.On("GetUserConfig", 21).Once().Return(uc, nil)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	permissionManager.On("IsPermissionGrantedToUserFor", mock.Anything, user.PermissionWithdraw).Once().Return(true)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	userBalanceService.On("GetBalanceOfUserByCoinID", 21, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub := args.Get(2).(*userbalance.UserBalance)
		ub.ID = 1
		ub.UserID = 21
		ub.Amount = "2000"
		ub.FrozenAmount = "1000"
	})

	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID:     21,
		Status: user.StatusVerified,
	}

	params := payment.PreWithdrawParams{
		Coin:    "usdt",
		Amount:  "1100.00",
		Address: "ethAddress1",
		Network: "eth",
	}

	res, statusCode := paymentService.PreWithdraw(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "user balance is not enough to withdraw this much", res.Message)

	currencyService.AssertExpectations(t)
	walletService.AssertExpectations(t)
	userConfigService.AssertExpectations(t)
	permissionManager.AssertExpectations(t)
	userBalanceService.AssertExpectations(t)

}

func TestService_PreWithdraw_Successful_NeedEmailCodeAndTwoFa_UserConfigExists(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	usdtCoin := currency.Coin{
		ID:              1,
		Code:            "USDT",
		MinimumWithdraw: "10",
		MaximumWithdraw: "10000",
	}
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "USDT").Once().Return(usdtCoin, nil)
	currencyService.On("GetCoinByCode", "ETH").Once().Return(currency.Coin{ID: 3}, nil)
	walletService := new(mocks.WalletService)
	walletService.On("IsAddressValid", "USDT", "ethAddress1", "ETH").Once().Return(true, nil)

	uc := user.Config{
		ID:                                    1,
		IsReadOnly:                            false,
		IsWhiteListEnabled:                    false,
		IsEmailVerificationForWithdrawEnabled: true,
		IsTwoFaVerificationForWithdrawEnabled: true,
	}
	userConfigService := new(mocks.UserConfigService)
	userConfigService.On("GetUserConfig", 21).Once().Return(uc, nil)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	withdrawEmailConfirmationManager.On("IsAllowedToSendEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(true, nil)
	withdrawEmailConfirmationManager.On("CreateAndSendWithdrawEmailConfirmationCode", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)
	permissionManager := new(mocks.UserPermissionManager)
	permissionManager.On("IsPermissionGrantedToUserFor", mock.Anything, user.PermissionWithdraw).Once().Return(true)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	userBalanceService.On("GetBalanceOfUserByCoinID", 21, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub := args.Get(2).(*userbalance.UserBalance)
		ub.ID = 1
		ub.UserID = 21
		ub.Amount = "2000"
		ub.FrozenAmount = "1000"
	})

	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	priceGenerator.On("GetBTCUSDTPrice", mock.Anything).Once().Return("35000.0", nil)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	configs.On("GetEnv").Return("test")
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID:     21,
		Status: user.StatusVerified,
		Email:  "test@test.test",
	}

	params := payment.PreWithdrawParams{
		Coin:    "usdt",
		Amount:  "100.00",
		Address: "ethAddress1",
		Network: "eth",
	}

	res, statusCode := paymentService.PreWithdraw(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)
	preWithdrawResults, ok := res.Data.(payment.PreWithdrawResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}

	assert.Equal(t, true, preWithdrawResults.NeedEmailCode)
	assert.Equal(t, true, preWithdrawResults.Need2fa)

	currencyService.AssertExpectations(t)
	walletService.AssertExpectations(t)
	currencyService.AssertExpectations(t)
	walletService.AssertExpectations(t)
	userConfigService.AssertExpectations(t)
	permissionManager.AssertExpectations(t)
	userBalanceService.AssertExpectations(t)
	withdrawEmailConfirmationManager.AssertExpectations(t)
}

func TestService_PreWithdraw_Successful_NeedTwoFaAndEmailCode_UserConfigDoesNotExists(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	usdtCoin := currency.Coin{
		ID:              1,
		Code:            "USDT",
		MinimumWithdraw: "10",
		MaximumWithdraw: "10000",
	}
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "USDT").Once().Return(usdtCoin, nil)
	currencyService.On("GetCoinByCode", "ETH").Once().Return(currency.Coin{ID: 3}, nil)
	walletService := new(mocks.WalletService)
	walletService.On("IsAddressValid", "USDT", "ethAddress1", "ETH").Once().Return(true, nil)

	uc := user.Config{}
	userConfigService := new(mocks.UserConfigService)
	userConfigService.On("GetUserConfig", 21).Once().Return(uc, gorm.ErrRecordNotFound)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	withdrawEmailConfirmationManager.On("IsAllowedToSendEmail", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(true, nil)
	withdrawEmailConfirmationManager.On("CreateAndSendWithdrawEmailConfirmationCode", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)
	permissionManager := new(mocks.UserPermissionManager)
	permissionManager.On("IsPermissionGrantedToUserFor", mock.Anything, user.PermissionWithdraw).Once().Return(true)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	userBalanceService.On("GetBalanceOfUserByCoinID", 21, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub := args.Get(2).(*userbalance.UserBalance)
		ub.ID = 1
		ub.UserID = 21
		ub.Amount = "2000"
		ub.FrozenAmount = "1000"
	})

	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	configs.On("GetEnv").Return("test")
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID:             21,
		Status:         user.StatusVerified,
		IsTwoFaEnabled: true,
	}

	params := payment.PreWithdrawParams{
		Coin:    "usdt",
		Amount:  "100.00",
		Address: "ethAddress1",
		Network: "eth",
	}

	res, statusCode := paymentService.PreWithdraw(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)
	preWithdrawResults, ok := res.Data.(payment.PreWithdrawResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}

	assert.Equal(t, true, preWithdrawResults.NeedEmailCode)
	assert.Equal(t, true, preWithdrawResults.Need2fa)

	currencyService.AssertExpectations(t)
	walletService.AssertExpectations(t)
	currencyService.AssertExpectations(t)
	walletService.AssertExpectations(t)
	userConfigService.AssertExpectations(t)
	permissionManager.AssertExpectations(t)
	userBalanceService.AssertExpectations(t)
}

func TestService_Withdraw_Fail_CoinNotFound(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "BTCQ").Once().Return(currency.Coin{}, gorm.ErrRecordNotFound)
	walletService := new(mocks.WalletService)
	userConfigService := new(mocks.UserConfigService)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID: 21,
	}

	params := payment.WithdrawParams{
		Coin:    "BTCQ",
		Amount:  "0.1",
		Address: "btcAddress1",
		Network: "",
	}

	res, statusCode := paymentService.Withdraw(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "coin not found", res.Message)
	currencyService.AssertExpectations(t)
}

func TestService_Withdraw_Fail_NetworkCoinNotFound(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "USDT").Once().Return(currency.Coin{ID: 1}, nil)
	currencyService.On("GetCoinByCode", "ETHQ").Once().Return(currency.Coin{}, gorm.ErrRecordNotFound)
	walletService := new(mocks.WalletService)
	userConfigService := new(mocks.UserConfigService)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID: 21,
	}

	params := payment.WithdrawParams{
		Coin:    "usdt",
		Amount:  "0.1",
		Address: "ethAddress1",
		Network: "ethq",
	}

	res, statusCode := paymentService.Withdraw(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "network not found", res.Message)
	currencyService.AssertExpectations(t)

}

func TestService_Withdraw_Fail_AddressIsNotValid(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "USDT").Once().Return(currency.Coin{ID: 1, Code: "USDT"}, nil)
	currencyService.On("GetCoinByCode", "ETH").Once().Return(currency.Coin{ID: 3}, nil)
	walletService := new(mocks.WalletService)
	walletService.On("IsAddressValid", "USDT", "ethAddress1", "ETH").Once().Return(false, nil)
	userConfigService := new(mocks.UserConfigService)
	//userConfigService.On("GetUserConfig", int64(21)).Once().Return(user.Config{}, gorm.ErrRecordNotFound)

	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID:     21,
		Status: user.StatusVerified,
	}

	params := payment.WithdrawParams{
		Coin:    "usdt",
		Amount:  "0.10",
		Address: "ethAddress1",
		Network: "eth",
	}

	res, statusCode := paymentService.Withdraw(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "address is not valid", res.Message)

	currencyService.AssertExpectations(t)
	walletService.AssertExpectations(t)

}

func TestService_Withdraw_Fail_WrongAmount(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "USDT").Once().Return(currency.Coin{ID: 1, Code: "USDT"}, nil)
	currencyService.On("GetCoinByCode", "ETH").Once().Return(currency.Coin{ID: 3}, nil)
	walletService := new(mocks.WalletService)
	walletService.On("IsAddressValid", "USDT", "ethAddress1", "ETH").Once().Return(true, nil)
	userConfigService := new(mocks.UserConfigService)
	//userConfigService.On("GetUserConfig", int64(21)).Once().Return(user.Config{}, gorm.ErrRecordNotFound)

	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID:     21,
		Status: user.StatusVerified,
	}

	params := payment.WithdrawParams{
		Coin:    "usdt",
		Amount:  "-0.10",
		Address: "ethAddress1",
		Network: "eth",
	}

	res, statusCode := paymentService.Withdraw(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "amount is not correct", res.Message)

	currencyService.AssertExpectations(t)
	walletService.AssertExpectations(t)

}

func TestService_Withdraw_Fail_UserAccountStatusIsNotVerified(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "USDT").Once().Return(currency.Coin{ID: 1, Code: "USDT"}, nil)
	currencyService.On("GetCoinByCode", "ETH").Once().Return(currency.Coin{ID: 3}, nil)
	walletService := new(mocks.WalletService)
	walletService.On("IsAddressValid", "USDT", "ethAddress1", "ETH").Once().Return(true, nil)
	userConfigService := new(mocks.UserConfigService)
	userConfigService.On("GetUserConfig", 21).Once().Return(user.Config{}, gorm.ErrRecordNotFound)

	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID:     21,
		Status: user.StatusRegistered,
	}

	params := payment.WithdrawParams{
		Coin:    "usdt",
		Amount:  "0.10",
		Address: "ethAddress1",
		Network: "eth",
	}

	res, statusCode := paymentService.Withdraw(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "user account is not verified", res.Message)

	currencyService.AssertExpectations(t)
	walletService.AssertExpectations(t)
	userConfigService.AssertExpectations(t)

}

func TestService_Withdraw_Fail_UserDisabledTwoFa(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "USDT").Once().Return(currency.Coin{ID: 1, Code: "USDT"}, nil)
	currencyService.On("GetCoinByCode", "ETH").Once().Return(currency.Coin{ID: 3}, nil)
	walletService := new(mocks.WalletService)
	walletService.On("IsAddressValid", "USDT", "ethAddress1", "ETH").Once().Return(true, nil)
	userConfigService := new(mocks.UserConfigService)
	userConfigService.On("GetUserConfig", 21).Once().Return(user.Config{}, gorm.ErrRecordNotFound)

	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID:                  21,
		Status:              user.StatusVerified,
		Google2faDisabledAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	params := payment.WithdrawParams{
		Coin:    "usdt",
		Amount:  "0.10",
		Address: "ethAddress1",
		Network: "eth",
	}

	res, statusCode := paymentService.Withdraw(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "for the security reasons, after disabling/enabling 2fa the withdraw request is not allowed for 24 hours", res.Message)

	currencyService.AssertExpectations(t)
	walletService.AssertExpectations(t)
	userConfigService.AssertExpectations(t)

}

func TestService_Withdraw_Fail_PermissionNotGranted(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "USDT").Once().Return(currency.Coin{ID: 1, Code: "USDT"}, nil)
	currencyService.On("GetCoinByCode", "ETH").Once().Return(currency.Coin{ID: 3}, nil)
	walletService := new(mocks.WalletService)
	walletService.On("IsAddressValid", "USDT", "ethAddress1", "ETH").Once().Return(true, nil)
	userConfigService := new(mocks.UserConfigService)
	userConfigService.On("GetUserConfig", 21).Once().Return(user.Config{}, gorm.ErrRecordNotFound)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	permissionManager.On("IsPermissionGrantedToUserFor", mock.Anything, user.PermissionWithdraw).Once().Return(false)

	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID:     21,
		Status: user.StatusVerified,
	}

	params := payment.WithdrawParams{
		Coin:    "usdt",
		Amount:  "0.10",
		Address: "ethAddress1",
		Network: "eth",
	}

	res, statusCode := paymentService.Withdraw(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "withdraw permission is not granted", res.Message)

	currencyService.AssertExpectations(t)
	walletService.AssertExpectations(t)
	userConfigService.AssertExpectations(t)
	permissionManager.AssertExpectations(t)

}

func TestService_Withdraw_Fail_LessThanMinimumWithdraw(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	usdtCoin := currency.Coin{
		ID:              1,
		Code:            "USDT",
		MinimumWithdraw: "10",
	}
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "USDT").Once().Return(usdtCoin, nil)
	currencyService.On("GetCoinByCode", "ETH").Once().Return(currency.Coin{ID: 3}, nil)
	walletService := new(mocks.WalletService)
	walletService.On("IsAddressValid", "USDT", "ethAddress1", "ETH").Once().Return(true, nil)
	userConfigService := new(mocks.UserConfigService)
	userConfigService.On("GetUserConfig", 21).Once().Return(user.Config{}, gorm.ErrRecordNotFound)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	permissionManager.On("IsPermissionGrantedToUserFor", mock.Anything, user.PermissionWithdraw).Once().Return(true)

	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID:     21,
		Status: user.StatusVerified,
	}

	params := payment.WithdrawParams{
		Coin:    "usdt",
		Amount:  "0.10",
		Address: "ethAddress1",
		Network: "eth",
	}

	res, statusCode := paymentService.Withdraw(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "minimum withdraw is: 10", res.Message)

	currencyService.AssertExpectations(t)
	walletService.AssertExpectations(t)
	userConfigService.AssertExpectations(t)
	permissionManager.AssertExpectations(t)

}

func TestService_Withdraw_Fail_MoreThanMaximumWithdraw(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	usdtCoin := currency.Coin{
		ID:              1,
		Code:            "USDT",
		MinimumWithdraw: "10",
		MaximumWithdraw: "1000",
	}
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "USDT").Once().Return(usdtCoin, nil)
	currencyService.On("GetCoinByCode", "ETH").Once().Return(currency.Coin{ID: 3}, nil)
	walletService := new(mocks.WalletService)
	walletService.On("IsAddressValid", "USDT", "ethAddress1", "ETH").Once().Return(true, nil)
	userConfigService := new(mocks.UserConfigService)
	userConfigService.On("GetUserConfig", 21).Once().Return(user.Config{}, gorm.ErrRecordNotFound)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	permissionManager.On("IsPermissionGrantedToUserFor", mock.Anything, user.PermissionWithdraw).Once().Return(true)

	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID:     21,
		Status: user.StatusVerified,
	}

	params := payment.WithdrawParams{
		Coin:    "usdt",
		Amount:  "1100.00",
		Address: "ethAddress1",
		Network: "eth",
	}

	res, statusCode := paymentService.Withdraw(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "maximum withdraw is: 1000", res.Message)

	currencyService.AssertExpectations(t)
	walletService.AssertExpectations(t)
	userConfigService.AssertExpectations(t)
	permissionManager.AssertExpectations(t)

}

func TestService_Withdraw_Fail_AccountIsReadOnly(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	usdtCoin := currency.Coin{
		ID:              1,
		Code:            "USDT",
		MinimumWithdraw: "10",
		MaximumWithdraw: "100",
	}
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "USDT").Once().Return(usdtCoin, nil)
	currencyService.On("GetCoinByCode", "ETH").Once().Return(currency.Coin{ID: 3}, nil)
	walletService := new(mocks.WalletService)
	walletService.On("IsAddressValid", "USDT", "ethAddress1", "ETH").Once().Return(true, nil)

	uc := user.Config{
		ID:         1,
		IsReadOnly: true,
	}
	userConfigService := new(mocks.UserConfigService)
	userConfigService.On("GetUserConfig", 21).Once().Return(uc, nil)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)

	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID:     21,
		Status: user.StatusVerified,
	}

	params := payment.WithdrawParams{
		Coin:    "usdt",
		Amount:  "10.00",
		Address: "ethAddress1",
		Network: "eth",
	}

	res, statusCode := paymentService.Withdraw(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "this account is in read only mode", res.Message)

	currencyService.AssertExpectations(t)
	walletService.AssertExpectations(t)
	userConfigService.AssertExpectations(t)

}

func TestService_Withdraw_Fail_WhitelistEnabled(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	usdtCoin := currency.Coin{
		ID:              1,
		Code:            "USDT",
		MinimumWithdraw: "10",
		MaximumWithdraw: "100",
	}
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "USDT").Once().Return(usdtCoin, nil)
	currencyService.On("GetCoinByCode", "ETH").Once().Return(currency.Coin{ID: 3}, nil)
	walletService := new(mocks.WalletService)
	walletService.On("IsAddressValid", "USDT", "ethAddress1", "ETH").Once().Return(true, nil)

	uc := user.Config{
		ID:                 1,
		IsReadOnly:         false,
		IsWhiteListEnabled: true,
	}
	userConfigService := new(mocks.UserConfigService)
	userConfigService.On("GetUserConfig", 21).Once().Return(uc, nil)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)

	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userWithdrawAddressService.On("GetUserWithdrawAddressesByAddress", mock.Anything, mock.Anything, mock.Anything).Once().Return([]userwithdrawaddress.UserWithdrawAddress{})
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID:     21,
		Status: user.StatusVerified,
	}

	params := payment.WithdrawParams{
		Coin:    "usdt",
		Amount:  "10.00",
		Address: "ethAddress1",
		Network: "eth",
	}

	res, statusCode := paymentService.Withdraw(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "this address is not in white list", res.Message)

	currencyService.AssertExpectations(t)
	walletService.AssertExpectations(t)
	userConfigService.AssertExpectations(t)
	userWithdrawAddressService.AssertExpectations(t)

}

func TestService_Withdraw_Fail_UserBalanceIsNotEnough(t *testing.T) {
	qm := queryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})

	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	dbMock.ExpectRollback()
	paymentRepo := new(mocks.PaymentRepository)
	usdtCoin := currency.Coin{
		ID:              1,
		Code:            "USDT",
		MinimumWithdraw: "10",
		MaximumWithdraw: "10000",
	}
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "USDT").Once().Return(usdtCoin, nil)
	currencyService.On("GetCoinByCode", "ETH").Once().Return(currency.Coin{ID: 3}, nil)
	walletService := new(mocks.WalletService)
	walletService.On("IsAddressValid", "USDT", "ethAddress1", "ETH").Once().Return(true, nil)

	uc := user.Config{
		ID:                 1,
		IsReadOnly:         false,
		IsWhiteListEnabled: false,
	}
	userConfigService := new(mocks.UserConfigService)
	userConfigService.On("GetUserConfig", 21).Once().Return(uc, nil)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	withdrawEmailConfirmationManager.On("CheckCode", mock.Anything, "123456").Once().Return(true, nil)

	permissionManager := new(mocks.UserPermissionManager)
	permissionManager.On("IsPermissionGrantedToUserFor", mock.Anything, user.PermissionWithdraw).Once().Return(true)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	userBalanceService.On("GetBalanceOfUserByCoinUsingTx", mock.Anything, 21, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub := args.Get(3).(*userbalance.UserBalance)
		ub.ID = 1
		ub.UserID = 21
		ub.Amount = "2000"
		ub.FrozenAmount = "1000"
	})

	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID:     21,
		Status: user.StatusVerified,
	}

	params := payment.WithdrawParams{
		Coin:      "usdt",
		Amount:    "1100.00",
		Address:   "ethAddress1",
		Network:   "eth",
		EmailCode: "123456",
	}

	res, statusCode := paymentService.Withdraw(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "user balance is not enough to withdraw this much", res.Message)

	currencyService.AssertExpectations(t)
	walletService.AssertExpectations(t)
	userConfigService.AssertExpectations(t)
	permissionManager.AssertExpectations(t)
	userBalanceService.AssertExpectations(t)

}

func TestService_Withdraw_Fail_WrongTwoFaCode_UserConfigExist(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	usdtCoin := currency.Coin{
		ID:              1,
		Code:            "USDT",
		MinimumWithdraw: "10",
		MaximumWithdraw: "10000",
	}
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "USDT").Once().Return(usdtCoin, nil)
	currencyService.On("GetCoinByCode", "ETH").Once().Return(currency.Coin{ID: 3}, nil)
	walletService := new(mocks.WalletService)
	walletService.On("IsAddressValid", "USDT", "ethAddress1", "ETH").Once().Return(true, nil)

	uc := user.Config{
		ID:                                    1,
		IsReadOnly:                            false,
		IsWhiteListEnabled:                    false,
		IsTwoFaVerificationForWithdrawEnabled: true,
	}
	userConfigService := new(mocks.UserConfigService)
	userConfigService.On("GetUserConfig", 21).Once().Return(uc, nil)
	twoFaManager := new(mocks.TwoFaManager)
	twoFaManager.On("CheckCode", mock.Anything, "123456").Once().Return(false)

	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)

	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID:     21,
		Status: user.StatusVerified,
	}

	params := payment.WithdrawParams{
		Coin:      "usdt",
		Amount:    "1100.00",
		Address:   "ethAddress1",
		Network:   "eth",
		TwoFaCode: "123456",
	}

	res, statusCode := paymentService.Withdraw(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "2fa code is not correct", res.Message)

	currencyService.AssertExpectations(t)
	walletService.AssertExpectations(t)
	userConfigService.AssertExpectations(t)
	twoFaManager.AssertExpectations(t)

}

func TestService_Withdraw_Fail_WrongEmailCode_UserConfigExist(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	usdtCoin := currency.Coin{
		ID:              1,
		Code:            "USDT",
		MinimumWithdraw: "10",
		MaximumWithdraw: "10000",
	}
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "USDT").Once().Return(usdtCoin, nil)
	currencyService.On("GetCoinByCode", "ETH").Once().Return(currency.Coin{ID: 3}, nil)
	walletService := new(mocks.WalletService)
	walletService.On("IsAddressValid", "USDT", "ethAddress1", "ETH").Once().Return(true, nil)

	uc := user.Config{
		ID:                                    1,
		IsReadOnly:                            false,
		IsWhiteListEnabled:                    false,
		IsTwoFaVerificationForWithdrawEnabled: false,
		IsEmailVerificationForWithdrawEnabled: true,
	}
	userConfigService := new(mocks.UserConfigService)
	userConfigService.On("GetUserConfig", 21).Once().Return(uc, nil)
	twoFaManager := new(mocks.TwoFaManager)

	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	withdrawEmailConfirmationManager.On("CheckCode", mock.Anything, "123456").Once().Return(false, nil)

	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)

	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID:     21,
		Status: user.StatusVerified,
	}

	params := payment.WithdrawParams{
		Coin:      "usdt",
		Amount:    "1100.00",
		Address:   "ethAddress1",
		Network:   "eth",
		EmailCode: "123456",
	}

	res, statusCode := paymentService.Withdraw(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "email confirmation code is not correct", res.Message)

	currencyService.AssertExpectations(t)
	walletService.AssertExpectations(t)
	userConfigService.AssertExpectations(t)

}

func TestService_Withdraw_Fail_WrongTwoFaCode_UserConfigDoesNotExist(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	usdtCoin := currency.Coin{
		ID:              1,
		Code:            "USDT",
		MinimumWithdraw: "10",
		MaximumWithdraw: "10000",
	}
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "USDT").Once().Return(usdtCoin, nil)
	currencyService.On("GetCoinByCode", "ETH").Once().Return(currency.Coin{ID: 3}, nil)
	walletService := new(mocks.WalletService)
	walletService.On("IsAddressValid", "USDT", "ethAddress1", "ETH").Once().Return(true, nil)

	uc := user.Config{}
	userConfigService := new(mocks.UserConfigService)
	userConfigService.On("GetUserConfig", 21).Once().Return(uc, gorm.ErrRecordNotFound)
	twoFaManager := new(mocks.TwoFaManager)
	twoFaManager.On("CheckCode", mock.Anything, "123456").Once().Return(false)

	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	withdrawEmailConfirmationManager.On("CheckCode", mock.Anything, "123456").Once().Return(true, nil)
	permissionManager := new(mocks.UserPermissionManager)
	permissionManager.On("IsPermissionGrantedToUserFor", mock.Anything, user.PermissionWithdraw).Once().Return(true)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)

	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID:             21,
		Status:         user.StatusVerified,
		IsTwoFaEnabled: true,
	}

	params := payment.WithdrawParams{
		Coin:      "usdt",
		Amount:    "1100.00",
		Address:   "ethAddress1",
		Network:   "eth",
		TwoFaCode: "123456",
		EmailCode: "123456",
	}

	res, statusCode := paymentService.Withdraw(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "2fa code is not correct", res.Message)

	currencyService.AssertExpectations(t)
	walletService.AssertExpectations(t)
	userConfigService.AssertExpectations(t)
	twoFaManager.AssertExpectations(t)

}

// queryMatcher is a sqlmock.QueryMatcher that accepts any SQL query,
// allowing tests to bypass strict SQL matching when the exact query
// text is not relevant to the scenario under test.
type queryMatcher struct {
}

// Match always returns nil, unconditionally matching any expected SQL
// against any actual SQL executed by the code under test.
func (queryMatcher) Match(expectedSQL, actualSQL string) error {
	return nil
}

func TestService_Withdraw_Successful_TwoFaAndEmailCodeAreCorrect_UserConfigExist(t *testing.T) {
	qm := queryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	dbMock.ExpectExec("INSERT INTO crypto_payments").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("INSERT INTO crypto_payment_extra_info").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(10, 1))
	dbMock.ExpectCommit()

	paymentRepo := new(mocks.PaymentRepository)

	usdtCoin := currency.Coin{
		ID:              1,
		Code:            "USDT",
		MinimumWithdraw: "10",
		MaximumWithdraw: "10000",
	}
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "USDT").Once().Return(usdtCoin, nil)
	currencyService.On("GetCoinByCode", "ETH").Once().Return(currency.Coin{ID: 3}, nil)
	walletService := new(mocks.WalletService)
	walletService.On("IsAddressValid", "USDT", "ethAddress1", "ETH").Once().Return(true, nil)

	uc := user.Config{
		ID:                                    1,
		IsReadOnly:                            false,
		IsWhiteListEnabled:                    false,
		IsTwoFaVerificationForWithdrawEnabled: true,
		IsEmailVerificationForWithdrawEnabled: true,
	}
	userConfigService := new(mocks.UserConfigService)
	userConfigService.On("GetUserConfig", 21).Once().Return(uc, nil)

	twoFaManager := new(mocks.TwoFaManager)
	twoFaManager.On("CheckCode", mock.Anything, "123456").Once().Return(true)

	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	withdrawEmailConfirmationManager.On("CheckCode", mock.Anything, "123456").Twice().Return(true, nil)
	withdrawEmailConfirmationManager.On("RemoveConfirmationCodeFromRedis", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)

	permissionManager := new(mocks.UserPermissionManager)
	permissionManager.On("IsPermissionGrantedToUserFor", mock.Anything, user.PermissionWithdraw).Once().Return(true)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userService.On("GetUserProfile", mock.Anything).Once().Return(user.Profile{}, nil)

	userBalanceService := new(mocks.UserBalanceService)
	userBalanceService.On("GetBalanceOfUserByCoinUsingTx", mock.Anything, 21, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub := args.Get(3).(*userbalance.UserBalance)
		ub.ID = 10
		ub.UserID = 21
		ub.Amount = "20000"
		ub.FrozenAmount = "1000"
	})

	communicationService := new(mocks.CommunicationService)
	communicationService.On("SendCryptoPaymentStatusUpdateEmail", mock.Anything, mock.Anything).Once().Return()

	priceGenerator := new(mocks.PriceGenerator)
	priceGenerator.On("GetBTCUSDTPrice", mock.Anything).Once().Return("35000.0", nil)
	priceGenerator.On("GetAmountBasedOnBTC", mock.Anything, mock.Anything, "1.0").Once().Return("0.00000350", nil)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID:     21,
		Status: user.StatusVerified,
		Email:  "test@test.test",
	}

	params := payment.WithdrawParams{
		Coin:      "usdt",
		Amount:    "2000.00",
		Address:   "ethAddress1",
		Network:   "eth",
		TwoFaCode: "123456",
		EmailCode: "123456",
	}

	res, statusCode := paymentService.Withdraw(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	resultData, ok := res.Data.(map[string][]payment.GetPaymentsResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}

	withdrawResults, ok := resultData["payments"]
	if !ok {
		t.Error("payments key does not exists in response")
	}

	assert.Equal(t, 1, len(withdrawResults))

	p1 := withdrawResults[0]
	assert.Equal(t, int64(12), p1.ID)
	assert.Equal(t, "pending", p1.Status)
	assert.Equal(t, "withdraw", p1.Type)
	assert.Equal(t, "2000.00000000", p1.Amount)
	assert.Equal(t, "USDT", p1.Coin)
	assert.Equal(t, "ethAddress1", p1.Address)
	assert.Equal(t, "", p1.TxID)

	time.Sleep(100 * time.Millisecond)

	currencyService.AssertExpectations(t)
	walletService.AssertExpectations(t)
	userConfigService.AssertExpectations(t)
	permissionManager.AssertExpectations(t)
	userBalanceService.AssertExpectations(t)
	twoFaManager.AssertExpectations(t)
	withdrawEmailConfirmationManager.AssertExpectations(t)
	communicationService.AssertExpectations(t)

}

func TestService_Withdraw_Successful_TwoFaIsCorrect_UserConfigDoesExist(t *testing.T) {
	qm := queryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	dbMock.ExpectExec("INSERT INTO crypto_payments").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("INSERT INTO crypto_payment_extra_info").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(10, 1))
	dbMock.ExpectCommit()

	paymentRepo := new(mocks.PaymentRepository)

	usdtCoin := currency.Coin{
		ID:              1,
		Code:            "USDT",
		MinimumWithdraw: "10",
		MaximumWithdraw: "10000",
	}
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "USDT").Once().Return(usdtCoin, nil)
	currencyService.On("GetCoinByCode", "ETH").Once().Return(currency.Coin{ID: 3}, nil)
	walletService := new(mocks.WalletService)
	walletService.On("IsAddressValid", "USDT", "ethAddress1", "ETH").Once().Return(true, nil)

	uc := user.Config{}
	userConfigService := new(mocks.UserConfigService)
	userConfigService.On("GetUserConfig", 21).Once().Return(uc, gorm.ErrRecordNotFound)

	twoFaManager := new(mocks.TwoFaManager)
	twoFaManager.On("CheckCode", mock.Anything, "123456").Once().Return(true)

	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	withdrawEmailConfirmationManager.On("CheckCode", mock.Anything, "123456").Once().Return(true, nil)
	withdrawEmailConfirmationManager.On("RemoveConfirmationCodeFromRedis", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)

	permissionManager := new(mocks.UserPermissionManager)
	permissionManager.On("IsPermissionGrantedToUserFor", mock.Anything, user.PermissionWithdraw).Once().Return(true)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userService.On("GetUserProfile", mock.Anything).Once().Return(user.Profile{}, nil)

	userBalanceService := new(mocks.UserBalanceService)
	userBalanceService.On("GetBalanceOfUserByCoinUsingTx", mock.Anything, 21, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub := args.Get(3).(*userbalance.UserBalance)
		ub.ID = 10
		ub.UserID = 21
		ub.Amount = "20000"
		ub.FrozenAmount = "1000"
	})

	communicationService := new(mocks.CommunicationService)
	communicationService.On("SendCryptoPaymentStatusUpdateEmail", mock.Anything, mock.Anything).Once().Return()

	priceGenerator := new(mocks.PriceGenerator)
	priceGenerator.On("GetBTCUSDTPrice", mock.Anything).Once().Return("35000.0", nil)
	priceGenerator.On("GetAmountBasedOnBTC", mock.Anything, mock.Anything, "1.0").Once().Return("0.00000350", nil)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID:             21,
		Status:         user.StatusVerified,
		IsTwoFaEnabled: true,
		Email:          "test@test.test",
	}

	params := payment.WithdrawParams{
		Coin:      "usdt",
		Amount:    "2000.00",
		Address:   "ethAddress1",
		Network:   "eth",
		TwoFaCode: "123456",
		EmailCode: "123456",
	}

	res, statusCode := paymentService.Withdraw(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	resultData, ok := res.Data.(map[string][]payment.GetPaymentsResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}

	withdrawResults, ok := resultData["payments"]
	if !ok {
		t.Error("payments key does not exists in response")
	}

	assert.Equal(t, 1, len(withdrawResults))

	p1 := withdrawResults[0]
	assert.Equal(t, int64(12), p1.ID)
	assert.Equal(t, "pending", p1.Status)
	assert.Equal(t, "withdraw", p1.Type)
	assert.Equal(t, "2000.00000000", p1.Amount)
	assert.Equal(t, "USDT", p1.Coin)
	assert.Equal(t, "ethAddress1", p1.Address)
	assert.Equal(t, "", p1.TxID)

	time.Sleep(100 * time.Millisecond)

	currencyService.AssertExpectations(t)
	walletService.AssertExpectations(t)
	userConfigService.AssertExpectations(t)
	permissionManager.AssertExpectations(t)
	userBalanceService.AssertExpectations(t)
	twoFaManager.AssertExpectations(t)
	communicationService.AssertExpectations(t)

}

func TestService_GetInProgressWithdrawalsInExternalExchange(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	data := []payment.ExternalWithdrawalUpdateDataNeeded{
		{
			PaymentID:                  1,
			PaymentExtraInfoID:         0,
			UpdatedAt:                  time.Time{},
			ExternalExchangeWithdrawID: "",
		},
		{
			PaymentID:                  2,
			PaymentExtraInfoID:         0,
			UpdatedAt:                  time.Time{},
			ExternalExchangeWithdrawID: "",
		},
	}
	paymentRepo.On("GetInProgressWithdrawalsInExternalExchange").Once().Return(data)

	currencyService := new(mocks.CurrencyService)
	walletService := new(mocks.WalletService)
	userConfigService := new(mocks.UserConfigService)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	result := paymentService.GetInProgressWithdrawalsInExternalExchange()

	assert.Equal(t, 2, len(result))
	assert.Equal(t, int64(1), result[0].PaymentID)
	assert.Equal(t, int64(2), result[1].PaymentID)
	paymentRepo.AssertExpectations(t)

}

func TestService_UpdatePaymentInExternalExchange_StatusCompleted(t *testing.T) {
	qm := queryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	dbMock.ExpectExec("UPDATE crypto_payments").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("UPDATE crypto_payment_extra_info").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(10, 1))
	dbMock.ExpectExec("INSERT INTO  transactions").WillReturnResult(sqlmock.NewResult(10, 1))
	dbMock.ExpectExec("INSERT INTO  transactions").WillReturnResult(sqlmock.NewResult(10, 1))
	dbMock.ExpectCommit()

	paymentRepo := new(mocks.PaymentRepository)
	paymentRepo.On("GetPaymentByIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		p := args.Get(2).(*payment.Payment)
		p.ID = 1
		p.CoinID = 1
		p.UserID = 1
		p.Amount = sql.NullString{String: "0.50000000", Valid: true}
	})

	paymentRepo.On("GetExtraInfoByPaymentIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		extraInfo := args.Get(2).(*payment.ExtraInfo)
		extraInfo.ID = 1
	})

	currencyService := new(mocks.CurrencyService)
	walletService := new(mocks.WalletService)
	userConfigService := new(mocks.UserConfigService)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userService.On("GetUserProfile", mock.Anything).Once().Return(user.Profile{}, nil)
	userService.On("GetUserByID", 1).Once().Return(user.User{ID: 1}, nil)

	userBalanceService := new(mocks.UserBalanceService)
	userBalanceService.On("GetBalanceOfUserByCoinUsingTx", mock.Anything, 1, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub := args.Get(3).(*userbalance.UserBalance)
		ub.Amount = "1.00000000"
		ub.FrozenAmount = "0.50000000"
	})

	communicationService := new(mocks.CommunicationService)
	communicationService.On("SendCryptoPaymentStatusUpdateEmail", mock.Anything, mock.Anything).Once().Return()
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	params := payment.UpdatePaymentInExternalExchangeParams{
		PaymentID:          1,
		PaymentExtraInfoID: 0,
		TxID:               "txId",
		Status:             payment.StatusCompleted,
		Data:               "",
	}
	paymentService.UpdatePaymentInExternalExchange(params)
	time.Sleep(20 * time.Millisecond)
	paymentRepo.AssertExpectations(t)
	userBalanceService.AssertExpectations(t)
	communicationService.AssertExpectations(t)
}

func TestService_UpdatePaymentInExternalExchange_StatusFailed(t *testing.T) {
	qm := queryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	dbMock.ExpectExec("UPDATE crypto_payments").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("UPDATE crypto_payment_extra_info").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(10, 1))
	dbMock.ExpectCommit()

	paymentRepo := new(mocks.PaymentRepository)
	paymentRepo.On("GetPaymentByIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		p := args.Get(2).(*payment.Payment)
		p.ID = 1
		p.CoinID = 1
		p.UserID = 1
		p.Amount = sql.NullString{String: "0.50000000", Valid: true}
	})

	paymentRepo.On("GetExtraInfoByPaymentIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		extraInfo := args.Get(2).(*payment.ExtraInfo)
		extraInfo.ID = 1
	})

	currencyService := new(mocks.CurrencyService)
	walletService := new(mocks.WalletService)
	userConfigService := new(mocks.UserConfigService)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userService.On("GetUserProfile", mock.Anything).Once().Return(user.Profile{}, nil)
	userService.On("GetUserByID", 1).Once().Return(user.User{ID: 1}, nil)

	userBalanceService := new(mocks.UserBalanceService)
	userBalanceService.On("GetBalanceOfUserByCoinUsingTx", mock.Anything, 1, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub := args.Get(3).(*userbalance.UserBalance)
		ub.Amount = "1.00000000"
		ub.FrozenAmount = "0.50000000"
	})

	communicationService := new(mocks.CommunicationService)
	communicationService.On("SendCryptoPaymentStatusUpdateEmail", mock.Anything, mock.Anything).Once().Return()
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	params := payment.UpdatePaymentInExternalExchangeParams{
		PaymentID:          1,
		PaymentExtraInfoID: 0,
		TxID:               "txId",
		Status:             payment.StatusRejected,
		Data:               "",
	}
	paymentService.UpdatePaymentInExternalExchange(params)
	time.Sleep(20 * time.Millisecond)

	paymentRepo.AssertExpectations(t)
	userBalanceService.AssertExpectations(t)
}

func TestService_CancelWithdraw_PaymentDoesNotExist(t *testing.T) {
	qm := queryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})

	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	dbMock.ExpectRollback()

	paymentRepo := new(mocks.PaymentRepository)
	paymentRepo.On("GetPaymentByIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(gorm.ErrRecordNotFound)

	currencyService := new(mocks.CurrencyService)
	walletService := new(mocks.WalletService)
	userConfigService := new(mocks.UserConfigService)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID: 21,
	}

	params := payment.CancelWithdrawParams{
		ID: 1,
	}
	res, statusCode := paymentService.CancelWithdraw(&u, params)

	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "withdraw not found", res.Message)
	paymentRepo.AssertExpectations(t)
}

func TestService_CancelWithdraw_PaymentIsNotWithdraw(t *testing.T) {
	qm := queryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})

	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	dbMock.ExpectRollback()

	paymentRepo := new(mocks.PaymentRepository)
	paymentRepo.On("GetPaymentByIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		p := args.Get(2).(*payment.Payment)
		p.ID = 1
		p.UserID = 21
		p.Type = payment.TypeDeposit
	})

	currencyService := new(mocks.CurrencyService)
	walletService := new(mocks.WalletService)
	userConfigService := new(mocks.UserConfigService)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID: 21,
	}

	params := payment.CancelWithdrawParams{
		ID: 1,
	}
	res, statusCode := paymentService.CancelWithdraw(&u, params)

	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "withdraw not found", res.Message)
	paymentRepo.AssertExpectations(t)
}

func TestService_CancelWithdraw_PaymentDoesNotBelongToUser(t *testing.T) {
	qm := queryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})

	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	dbMock.ExpectRollback()

	paymentRepo := new(mocks.PaymentRepository)
	paymentRepo.On("GetPaymentByIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		p := args.Get(2).(*payment.Payment)
		p.ID = 1
		p.UserID = 1
	})

	currencyService := new(mocks.CurrencyService)
	walletService := new(mocks.WalletService)
	userConfigService := new(mocks.UserConfigService)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID: 21,
	}

	params := payment.CancelWithdrawParams{
		ID: 1,
	}
	res, statusCode := paymentService.CancelWithdraw(&u, params)

	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "withdraw not found", res.Message)
	paymentRepo.AssertExpectations(t)

}

func TestService_CancelWithdraw_PaymentStatusIsNotCreated(t *testing.T) {
	qm := queryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})

	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	dbMock.ExpectRollback()

	paymentRepo := new(mocks.PaymentRepository)
	paymentRepo.On("GetPaymentByIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		p := args.Get(2).(*payment.Payment)
		p.ID = 1
		p.UserID = 21
		p.Type = payment.TypeWithdraw
		p.Status = payment.StatusInProgress
	})

	currencyService := new(mocks.CurrencyService)
	walletService := new(mocks.WalletService)
	userConfigService := new(mocks.UserConfigService)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID: 21,
	}

	params := payment.CancelWithdrawParams{
		ID: 1,
	}
	res, statusCode := paymentService.CancelWithdraw(&u, params)

	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "withdraw can't be cancelled now", res.Message)
	paymentRepo.AssertExpectations(t)
}

func TestService_CancelWithdraw_Successful(t *testing.T) {
	qm := queryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	dbMock.ExpectExec("UPDATE crpyto_payments").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit()

	paymentRepo := new(mocks.PaymentRepository)
	paymentRepo.On("GetPaymentByIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		p := args.Get(2).(*payment.Payment)
		p.ID = 1
		p.UserID = 21
		p.CoinID = 1
		p.Code = "BTC"
		p.Type = payment.TypeWithdraw
		p.Status = payment.StatusCreated
		p.Amount = sql.NullString{String: "1.00000000", Valid: true}
	})

	currencyService := new(mocks.CurrencyService)
	walletService := new(mocks.WalletService)
	userConfigService := new(mocks.UserConfigService)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)

	ub := &userbalance.UserBalance{}
	userBalanceService.On("GetBalanceOfUserByCoinUsingTx", mock.Anything, 21, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub = args.Get(3).(*userbalance.UserBalance)
		ub.ID = 1
		ub.CoinID = 1
		ub.FrozenAmount = "1.00000000"
	})

	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID: 21,
	}

	params := payment.CancelWithdrawParams{
		ID: 1,
	}
	res, statusCode := paymentService.CancelWithdraw(&u, params)

	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	resultData, ok := res.Data.(map[string][]payment.GetPaymentsResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}

	withdrawResults, ok := resultData["payments"]
	if !ok {
		t.Error("payments key does not exists in response")
	}

	assert.Equal(t, 1, len(withdrawResults))

	p1 := withdrawResults[0]
	assert.Equal(t, int64(1), p1.ID)
	assert.Equal(t, "user canceled", p1.Status)
	assert.Equal(t, "withdraw", p1.Type)
	assert.Equal(t, "1.00000000", p1.Amount)
	assert.Equal(t, "BTC", p1.Coin)
	assert.Equal(t, "", p1.Address)
	assert.Equal(t, "", p1.TxID)

	assert.Equal(t, "0.00000000", ub.FrozenAmount)

	paymentRepo.AssertExpectations(t)
	userBalanceService.AssertExpectations(t)
}

func TestService_HandleWalletCallback_internalTransfer(t *testing.T) {
	db := &gorm.DB{}
	paymentRepo := new(mocks.PaymentRepository)
	currencyService := new(mocks.CurrencyService)
	walletService := new(mocks.WalletService)
	userConfigService := new(mocks.UserConfigService)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	internalTransferService.On("UpdateStatus", int64(1), payment.StatusCompleted).Once().Return(nil)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID: 21,
	}

	params := payment.WalletCallBackParams{
		Code:        "BTC",
		Amount:      "0.1",
		Type:        payment.TypeWithdraw,
		Status:      payment.StatusCompleted,
		FromAddress: "fromAddress",
		ToAddress:   "toAddress",
		TxID:        "txId",
		Meta:        "{\"internal_transfer_id\":\"1\"}",
		Network:     "",
	}

	res, statusCode := paymentService.HandleWalletCallBack(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)
	internalTransferService.AssertExpectations(t)
}

func TestService_HandleWalletCallback_Deposit_AlreadyDoesNotExist_Status_Created(t *testing.T) {
	qm := queryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	dbMock.ExpectExec("INSERT INTO crypto_payments").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("INSERT INTO crypto_payment_extra_info").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectCommit()

	paymentRepo := new(mocks.PaymentRepository)
	paymentRepo.On("GetPaymentByCoinIDAndTxIDAndTypeUsingTx", mock.Anything, int64(1), "txId", payment.TypeDeposit, mock.Anything).Once().Return(nil)
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "BTC").Once().Return(currency.Coin{ID: 1, Code: "BTC"}, nil)
	walletService := new(mocks.WalletService)
	userConfigService := new(mocks.UserConfigService)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userService.On("GetUserByID", 21).Twice().Return(user.User{}, nil)
	userService.On("GetUserProfile", mock.Anything).Once().Return(user.Profile{}, nil)
	userBalanceService := new(mocks.UserBalanceService)
	ub := &userbalance.UserBalance{}
	userBalanceService.On("GetUserBalanceByCoinAndAddressUsingTx", mock.Anything, int64(1), "toAddress", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub = args.Get(3).(*userbalance.UserBalance)
		ub.ID = 1
		ub.UserID = 21
		ub.Amount = "1.0"
		ub.FrozenAmount = "0.5"
	})
	communicationService := new(mocks.CommunicationService)
	communicationService.On("SendCryptoPaymentStatusUpdateEmail", mock.Anything, mock.Anything).Once().Return()
	priceGenerator := new(mocks.PriceGenerator)
	priceGenerator.On("GetBTCUSDTPrice", mock.Anything).Once().Return("35000.0", nil)
	priceGenerator.On("GetAmountBasedOnBTC", mock.Anything, "BTC", "1.0").Once().Return("0.1", nil)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	mqttManager.On("PublishPayment", mock.Anything, mock.Anything, mock.Anything).Once().Return()
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID: 21,
	}

	params := payment.WalletCallBackParams{
		Code:        "BTC",
		Amount:      "0.1",
		Type:        payment.TypeDeposit,
		Status:      payment.StatusCreated,
		FromAddress: "fromAddress",
		ToAddress:   "toAddress",
		TxID:        "txId",
		Network:     "",
	}

	res, statusCode := paymentService.HandleWalletCallBack(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	assert.Equal(t, "1.0", ub.Amount)
	assert.Equal(t, "0.5", ub.FrozenAmount)
	time.Sleep(50 * time.Millisecond)
	paymentRepo.AssertExpectations(t)
	currencyService.AssertExpectations(t)
	userBalanceService.AssertExpectations(t)
	userService.AssertExpectations(t)
	mqttManager.AssertExpectations(t)
	priceGenerator.AssertExpectations(t)
	communicationService.AssertExpectations(t)
}

func TestService_HandleWalletCallback_Deposit_AlreadyDoesNotExist_Status_Completed(t *testing.T) {
	qm := queryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	dbMock.ExpectExec("INSERT INTO crypto_payments").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("INSERT INTO crypto_payment_extra_info").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(10, 1))
	dbMock.ExpectExec("INSERT INTO  transactions").WillReturnResult(sqlmock.NewResult(10, 1))
	dbMock.ExpectCommit()

	paymentRepo := new(mocks.PaymentRepository)
	paymentRepo.On("GetPaymentByCoinIDAndTxIDAndTypeUsingTx", mock.Anything, int64(1), "txId", payment.TypeDeposit, mock.Anything).Once().Return(nil)
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "BTC").Once().Return(currency.Coin{ID: 1, Code: "BTC"}, nil)
	walletService := new(mocks.WalletService)
	userConfigService := new(mocks.UserConfigService)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userService.On("GetUserByID", 21).Twice().Return(user.User{}, nil)
	userService.On("GetUserProfile", mock.Anything).Once().Return(user.Profile{}, nil)
	userBalanceService := new(mocks.UserBalanceService)
	ub := &userbalance.UserBalance{}
	userBalanceService.On("GetUserBalanceByCoinAndAddressUsingTx", mock.Anything, int64(1), "toAddress", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub = args.Get(3).(*userbalance.UserBalance)
		ub.ID = 1
		ub.UserID = 21
		ub.Amount = "1.0"
		ub.FrozenAmount = "0.5"
	})
	communicationService := new(mocks.CommunicationService)
	communicationService.On("SendCryptoPaymentStatusUpdateEmail", mock.Anything, mock.Anything).Once().Return()
	priceGenerator := new(mocks.PriceGenerator)
	priceGenerator.On("GetBTCUSDTPrice", mock.Anything).Once().Return("35000.0", nil)
	priceGenerator.On("GetAmountBasedOnBTC", mock.Anything, "BTC", "1.0").Once().Return("0.1", nil)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	autoExchangeManager.On("AutoExchange", mock.Anything, mock.Anything).Once().Return()
	mqttManager := new(mocks.CentrifugoManager)
	mqttManager.On("PublishPayment", mock.Anything, mock.Anything, mock.Anything).Once().Return()
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID: 21,
	}

	params := payment.WalletCallBackParams{
		Code:        "BTC",
		Amount:      "0.1",
		Type:        payment.TypeDeposit,
		Status:      payment.StatusCompleted,
		FromAddress: "fromAddress",
		ToAddress:   "toAddress",
		TxID:        "txId",
		Network:     "",
	}

	res, statusCode := paymentService.HandleWalletCallBack(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	assert.Equal(t, "1.10000000", ub.Amount)
	assert.Equal(t, "0.5", ub.FrozenAmount)
	time.Sleep(50 * time.Millisecond)
	paymentRepo.AssertExpectations(t)
	currencyService.AssertExpectations(t)
	userBalanceService.AssertExpectations(t)
	userService.AssertExpectations(t)
	mqttManager.AssertExpectations(t)
	priceGenerator.AssertExpectations(t)
	communicationService.AssertExpectations(t)
	autoExchangeManager.AssertExpectations(t)
}

func TestService_HandleWalletCallback_Deposit_AlreadyExist_Status_Completed(t *testing.T) {
	qm := queryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	dbMock.ExpectExec("UPDATE crypto_payments").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(10, 1))
	dbMock.ExpectExec("INSERT INTO  transactions").WillReturnResult(sqlmock.NewResult(10, 1))
	dbMock.ExpectCommit()

	paymentRepo := new(mocks.PaymentRepository)
	p := &payment.Payment{}
	paymentRepo.On("GetPaymentByCoinIDAndTxIDAndTypeUsingTx", mock.Anything, int64(1), "txId", payment.TypeDeposit, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		p = args.Get(4).(*payment.Payment)
		p.ID = 1
		p.UserID = 21
		p.Status = payment.StatusCreated
	})

	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "BTC").Once().Return(currency.Coin{ID: 1, Code: "BTC"}, nil)
	walletService := new(mocks.WalletService)
	userConfigService := new(mocks.UserConfigService)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userService.On("GetUserByID", 21).Twice().Return(user.User{}, nil)
	userService.On("GetUserProfile", mock.Anything).Once().Return(user.Profile{}, nil)
	userBalanceService := new(mocks.UserBalanceService)
	ub := &userbalance.UserBalance{}
	userBalanceService.On("GetUserBalanceByCoinAndAddressUsingTx", mock.Anything, int64(1), "toAddress", mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub = args.Get(3).(*userbalance.UserBalance)
		ub.ID = 1
		ub.UserID = 21
		ub.Amount = "1.0"
		ub.FrozenAmount = "0.5"
	})
	communicationService := new(mocks.CommunicationService)
	communicationService.On("SendCryptoPaymentStatusUpdateEmail", mock.Anything, mock.Anything).Once().Return()
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	autoExchangeManager.On("AutoExchange", mock.Anything, mock.Anything).Once().Return()
	mqttManager := new(mocks.CentrifugoManager)
	mqttManager.On("PublishPayment", mock.Anything, mock.Anything, mock.Anything).Once().Return()
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID: 21,
	}

	params := payment.WalletCallBackParams{
		Code:        "BTC",
		Amount:      "0.1",
		Type:        payment.TypeDeposit,
		Status:      payment.StatusCompleted,
		FromAddress: "fromAddress",
		ToAddress:   "toAddress",
		TxID:        "txId",
		Network:     "",
	}

	res, statusCode := paymentService.HandleWalletCallBack(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	assert.Equal(t, "txId", p.TxID.String)
	assert.Equal(t, payment.StatusCompleted, p.Status)

	assert.Equal(t, "1.10000000", ub.Amount)
	assert.Equal(t, "0.5", ub.FrozenAmount)
	time.Sleep(50 * time.Millisecond)
	paymentRepo.AssertExpectations(t)
	currencyService.AssertExpectations(t)
	userBalanceService.AssertExpectations(t)
	userService.AssertExpectations(t)
	mqttManager.AssertExpectations(t)
	communicationService.AssertExpectations(t)
	autoExchangeManager.AssertExpectations(t)
}

func TestService_HandleWalletCallback_Withdraw_Status_Failed(t *testing.T) {
	qm := queryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	dbMock.ExpectExec("UPDATE crypto_payments").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(10, 1))
	dbMock.ExpectCommit()

	paymentRepo := new(mocks.PaymentRepository)
	p := &payment.Payment{}
	paymentRepo.On("GetPaymentByCoinIDAndTxIDAndTypeUsingTx", mock.Anything, int64(1), "txId", payment.TypeWithdraw, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		p = args.Get(4).(*payment.Payment)
		p.ID = 1
		p.UserID = 21
		p.Status = payment.StatusCreated
		p.FeeAmount = sql.NullString{String: "0.01", Valid: true}
		p.Amount = sql.NullString{String: "0.1", Valid: true}
	})

	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "BTC").Once().Return(currency.Coin{ID: 1, Code: "BTC"}, nil)
	walletService := new(mocks.WalletService)
	userConfigService := new(mocks.UserConfigService)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userService.On("GetUserByID", 21).Twice().Return(user.User{}, nil)
	userService.On("GetUserProfile", mock.Anything).Once().Return(user.Profile{}, nil)
	userBalanceService := new(mocks.UserBalanceService)
	ub := &userbalance.UserBalance{}
	userBalanceService.On("GetBalanceOfUserByCoinUsingTx", mock.Anything, 21, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub = args.Get(3).(*userbalance.UserBalance)
		ub.ID = 1
		ub.UserID = 21
		ub.Amount = "1.00000000"
		ub.FrozenAmount = "0.50000000"
	})
	communicationService := new(mocks.CommunicationService)
	communicationService.On("SendCryptoPaymentStatusUpdateEmail", mock.Anything, mock.Anything).Once().Return()
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	mqttManager.On("PublishPayment", mock.Anything, mock.Anything, mock.Anything).Once().Return()
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID: 21,
	}

	params := payment.WalletCallBackParams{
		Code:        "BTC",
		Amount:      "0.1",
		Type:        payment.TypeWithdraw,
		Status:      payment.StatusFailed,
		FromAddress: "fromAddress",
		ToAddress:   "toAddress",
		TxID:        "txId",
		Network:     "",
	}

	res, statusCode := paymentService.HandleWalletCallBack(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	assert.Equal(t, payment.StatusFailed, p.Status)
	assert.Equal(t, "fromAddress", p.FromAddress.String)

	assert.Equal(t, "1.00000000", ub.Amount)
	assert.Equal(t, "0.40000000", ub.FrozenAmount)
	time.Sleep(50 * time.Millisecond)
	paymentRepo.AssertExpectations(t)
	currencyService.AssertExpectations(t)
	userBalanceService.AssertExpectations(t)
	userService.AssertExpectations(t)
	mqttManager.AssertExpectations(t)
	communicationService.AssertExpectations(t)
}

func TestService_HandleWalletCallback_Withdraw_Status_Completed(t *testing.T) {
	qm := queryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	dbMock.ExpectExec("UPDATE crypto_payments").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(10, 1))
	dbMock.ExpectExec("INSERT INTO  transactions").WillReturnResult(sqlmock.NewResult(10, 1))
	dbMock.ExpectExec("INSERT INTO  transactions").WillReturnResult(sqlmock.NewResult(11, 1))
	dbMock.ExpectCommit()

	paymentRepo := new(mocks.PaymentRepository)
	p := &payment.Payment{}
	paymentRepo.On("GetPaymentByCoinIDAndTxIDAndTypeUsingTx", mock.Anything, int64(1), "txId", payment.TypeWithdraw, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		p = args.Get(4).(*payment.Payment)
		p.ID = 1
		p.UserID = 21
		p.Status = payment.StatusCreated
		p.FeeAmount = sql.NullString{String: "0.01", Valid: true}
		p.Amount = sql.NullString{String: "0.1", Valid: true}
	})

	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetCoinByCode", "BTC").Once().Return(currency.Coin{ID: 1, Code: "BTC"}, nil)
	walletService := new(mocks.WalletService)
	userConfigService := new(mocks.UserConfigService)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userService.On("GetUserByID", 21).Twice().Return(user.User{}, nil)
	userService.On("GetUserProfile", mock.Anything).Once().Return(user.Profile{}, nil)
	userBalanceService := new(mocks.UserBalanceService)
	ub := &userbalance.UserBalance{}
	userBalanceService.On("GetBalanceOfUserByCoinUsingTx", mock.Anything, 21, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub = args.Get(3).(*userbalance.UserBalance)
		ub.ID = 1
		ub.UserID = 21
		ub.Amount = "1.0"
		ub.FrozenAmount = "0.5"
	})
	communicationService := new(mocks.CommunicationService)
	communicationService.On("SendCryptoPaymentStatusUpdateEmail", mock.Anything, mock.Anything).Once().Return()
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	mqttManager.On("PublishPayment", mock.Anything, mock.Anything, mock.Anything).Once().Return()
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID: 21,
	}

	params := payment.WalletCallBackParams{
		Code:        "BTC",
		Amount:      "0.1",
		Type:        payment.TypeWithdraw,
		Status:      payment.StatusCompleted,
		FromAddress: "fromAddress",
		ToAddress:   "toAddress",
		TxID:        "txId",
		Network:     "",
	}

	res, statusCode := paymentService.HandleWalletCallBack(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	assert.Equal(t, payment.StatusCompleted, p.Status)
	assert.Equal(t, "fromAddress", p.FromAddress.String)

	assert.Equal(t, "0.90000000", ub.Amount)
	assert.Equal(t, "0.40000000", ub.FrozenAmount)
	time.Sleep(50 * time.Millisecond)
	paymentRepo.AssertExpectations(t)
	currencyService.AssertExpectations(t)
	userBalanceService.AssertExpectations(t)
	userService.AssertExpectations(t)
	mqttManager.AssertExpectations(t)
	communicationService.AssertExpectations(t)
}

func TestService_UpdateWithdraw_UpdateAdminStatus_And_Fee_And_NetworkFee(t *testing.T) {
	qm := queryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	dbMock.ExpectExec("UPDATE crypto_payments").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("UPDATE crypto_payment_extra_info").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectCommit()

	paymentRepo := new(mocks.PaymentRepository)
	p := &payment.Payment{}
	paymentRepo.On("GetPaymentByIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		p = args.Get(2).(*payment.Payment)
		p.ID = 1
		p.UserID = 21
		p.Status = payment.StatusCreated
		p.FeeAmount = sql.NullString{String: "0.01", Valid: true}
		p.Amount = sql.NullString{String: "0.1", Valid: true}
	})
	ei := &payment.ExtraInfo{}
	paymentRepo.On("GetExtraInfoByPaymentIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ei = args.Get(2).(*payment.ExtraInfo)
		ei.ID = 1
	})

	currencyService := new(mocks.CurrencyService)
	walletService := new(mocks.WalletService)
	userConfigService := new(mocks.UserConfigService)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID: 23,
	}

	params := payment.UpdateWithdrawParams{
		ID:              1,
		Status:          "",
		AdminStatus:     payment.AdminStatusRecheck,
		Fee:             "0.01",
		NetworkFee:      "0.002",
		AutoTransfer:    nil,
		RejectionReason: "",
		WithdrawType:    "",
	}

	res, statusCode := paymentService.UpdateWithdraw(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	assert.Equal(t, payment.AdminStatusRecheck, p.AdminStatus.String)
	assert.Equal(t, true, p.AdminStatus.Valid)
	assert.Equal(t, "0.01000000", p.FeeAmount.String)
	assert.Equal(t, "0.002", ei.NetworkFee.String)
	assert.Equal(t, int64(23), ei.LastHandledID.Int64)

	paymentRepo.AssertExpectations(t)
}

func TestService_UpdateWithdraw_UpdateStatus_Reject(t *testing.T) {
	qm := queryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	dbMock.ExpectExec("UPDATE crypto_payments").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("UPDATE crypto_payment_extra_info").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectCommit()

	paymentRepo := new(mocks.PaymentRepository)
	p := &payment.Payment{}
	paymentRepo.On("GetPaymentByIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		p = args.Get(2).(*payment.Payment)
		p.ID = 1
		p.UserID = 21
		p.CoinID = 1
		p.Status = payment.StatusCreated
		p.FeeAmount = sql.NullString{String: "0.01", Valid: true}
		p.Amount = sql.NullString{String: "0.1", Valid: true}
	})
	ei := &payment.ExtraInfo{}
	paymentRepo.On("GetExtraInfoByPaymentIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ei = args.Get(2).(*payment.ExtraInfo)
		ei.ID = 1
	})

	currencyService := new(mocks.CurrencyService)
	walletService := new(mocks.WalletService)
	userConfigService := new(mocks.UserConfigService)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userService.On("GetUserByID", 21).Twice().Return(user.User{}, nil)
	userService.On("GetUserProfile", mock.Anything).Once().Return(user.Profile{}, nil)
	userBalanceService := new(mocks.UserBalanceService)
	ub := &userbalance.UserBalance{}
	userBalanceService.On("GetBalanceOfUserByCoinUsingTx", mock.Anything, 21, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub = args.Get(3).(*userbalance.UserBalance)
		ub.ID = 1
		ub.UserID = 21
		ub.Amount = "1.0"
		ub.FrozenAmount = "0.1"
	})

	communicationService := new(mocks.CommunicationService)
	communicationService.On("SendCryptoPaymentStatusUpdateEmail", mock.Anything, mock.Anything).Once().Return()
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	mqttManager.On("PublishPayment", mock.Anything, mock.Anything, mock.Anything).Once().Return()
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID: 23,
	}

	params := payment.UpdateWithdrawParams{
		ID:              1,
		Status:          payment.StatusRejected,
		AdminStatus:     "",
		Fee:             "",
		NetworkFee:      "",
		AutoTransfer:    nil,
		RejectionReason: "rejected",
		WithdrawType:    "",
	}

	res, statusCode := paymentService.UpdateWithdraw(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	assert.Equal(t, payment.StatusRejected, p.Status)
	assert.Equal(t, "rejected", ei.RejectionReason.String)
	assert.Equal(t, int64(23), ei.LastHandledID.Int64)

	assert.Equal(t, "0.00000000", ub.FrozenAmount)

	time.Sleep(50 * time.Millisecond)
	paymentRepo.AssertExpectations(t)
	userBalanceService.AssertExpectations(t)
	userService.AssertExpectations(t)
	communicationService.AssertExpectations(t)
	mqttManager.AssertExpectations(t)
}

func TestService_UpdateWithdraw_UpdateStatus_InProgressAndAutoTransfer_UsingHotWallet(t *testing.T) {
	qm := queryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	dbMock.ExpectExec("UPDATE crypto_payments").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("UPDATE crypto_payment_extra_info").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectCommit()

	paymentRepo := new(mocks.PaymentRepository)
	p := &payment.Payment{}
	paymentRepo.On("GetPaymentByIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		p = args.Get(2).(*payment.Payment)
		p.ID = 1
		p.UserID = 21
		p.CoinID = 1
		p.Code = "BTC"
		p.ToAddress = sql.NullString{String: "toAddress", Valid: true}
		p.Status = payment.StatusCreated
		p.FeeAmount = sql.NullString{String: "0.01", Valid: true}
		p.Amount = sql.NullString{String: "0.1", Valid: true}
		p.BlockchainNetwork = sql.NullString{String: "BTC", Valid: true}
	})
	ei := &payment.ExtraInfo{}
	paymentRepo.On("GetExtraInfoByPaymentIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ei = args.Get(2).(*payment.ExtraInfo)
		ei.ID = 1
	})

	currencyService := new(mocks.CurrencyService)
	walletService := new(mocks.WalletService)
	walletService.On("SendTransaction", "BTC", "0.09000000", "toAddress", "BTC", "0.02").Once().Return("txId", nil)
	userConfigService := new(mocks.UserConfigService)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userService.On("GetUserByID", 21).Twice().Return(user.User{}, nil)
	userService.On("GetUserProfile", mock.Anything).Once().Return(user.Profile{}, nil)
	userBalanceService := new(mocks.UserBalanceService)
	ub := &userbalance.UserBalance{}
	userBalanceService.On("GetBalanceOfUserByCoinUsingTx", mock.Anything, 21, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub = args.Get(3).(*userbalance.UserBalance)
		ub.ID = 1
		ub.UserID = 21
		ub.Amount = "1.0"
		ub.FrozenAmount = "0.1"
	})

	communicationService := new(mocks.CommunicationService)
	communicationService.On("SendCryptoPaymentStatusUpdateEmail", mock.Anything, mock.Anything).Once().Return()
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	mqttManager.On("PublishPayment", mock.Anything, mock.Anything, mock.Anything).Once().Return()
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID: 23,
	}

	autoTransfer := true
	params := payment.UpdateWithdrawParams{
		ID:              1,
		Status:          payment.StatusInProgress,
		AdminStatus:     "",
		Fee:             "",
		NetworkFee:      "0.02",
		AutoTransfer:    &autoTransfer,
		RejectionReason: "",
		WithdrawType:    payment.WithdrawTypeHotWallet,
	}

	res, statusCode := paymentService.UpdateWithdraw(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	assert.Equal(t, payment.StatusInProgress, p.Status)
	assert.Equal(t, "txId", p.TxID.String)
	assert.Equal(t, true, p.TxID.Valid)
	assert.Equal(t, "0.1", ub.FrozenAmount)

	assert.Equal(t, int64(23), ei.LastHandledID.Int64)

	time.Sleep(50 * time.Millisecond)
	paymentRepo.AssertExpectations(t)
	userBalanceService.AssertExpectations(t)
	userService.AssertExpectations(t)
	communicationService.AssertExpectations(t)
	mqttManager.AssertExpectations(t)
	walletService.AssertExpectations(t)
}

func TestService_UpdateWithdraw_UpdateStatus_InProgressAndAutoTransfer_UsingExternalExchangeWallet(t *testing.T) {
	qm := queryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	dbMock.ExpectExec("UPDATE crypto_payments").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("UPDATE crypto_payment_extra_info").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectCommit()

	paymentRepo := new(mocks.PaymentRepository)
	p := &payment.Payment{}
	paymentRepo.On("GetPaymentByIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		p = args.Get(2).(*payment.Payment)
		p.ID = 1
		p.UserID = 21
		p.CoinID = 1
		p.Code = "BTC"
		p.ToAddress = sql.NullString{String: "toAddress", Valid: true}
		p.Status = payment.StatusCreated
		p.FeeAmount = sql.NullString{String: "0.01", Valid: true}
		p.Amount = sql.NullString{String: "0.1", Valid: true}
		p.BlockchainNetwork = sql.NullString{String: "BTC", Valid: true}
	})
	ei := &payment.ExtraInfo{}
	paymentRepo.On("GetExtraInfoByPaymentIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ei = args.Get(2).(*payment.ExtraInfo)
		ei.ID = 1
	})

	currencyService := new(mocks.CurrencyService)
	walletService := new(mocks.WalletService)
	userConfigService := new(mocks.UserConfigService)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userService.On("GetUserByID", 21).Twice().Return(user.User{}, nil)
	userService.On("GetUserProfile", mock.Anything).Once().Return(user.Profile{}, nil)
	userBalanceService := new(mocks.UserBalanceService)
	ub := &userbalance.UserBalance{}
	userBalanceService.On("GetBalanceOfUserByCoinUsingTx", mock.Anything, 21, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub = args.Get(3).(*userbalance.UserBalance)
		ub.ID = 1
		ub.UserID = 21
		ub.Amount = "1.0"
		ub.FrozenAmount = "0.1"
	})

	communicationService := new(mocks.CommunicationService)
	communicationService.On("SendCryptoPaymentStatusUpdateEmail", mock.Anything, mock.Anything).Once().Return()
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	withdrawResult := externalexchange.WithdrawResult{
		ID:                 "someID",
		ErrorMessage:       "",
		ExternalExchangeID: 1,
	}
	externalExchangeService.On("Withdraw", mock.Anything).Once().Return(withdrawResult, nil)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	mqttManager.On("PublishPayment", mock.Anything, mock.Anything, mock.Anything).Once().Return()
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID: 23,
	}

	autoTransfer := true
	params := payment.UpdateWithdrawParams{
		ID:              1,
		Status:          payment.StatusInProgress,
		AdminStatus:     "",
		Fee:             "",
		NetworkFee:      "0.02",
		AutoTransfer:    &autoTransfer,
		RejectionReason: "",
		WithdrawType:    payment.WithdrawTypeExternalExchange,
	}

	res, statusCode := paymentService.UpdateWithdraw(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	assert.Equal(t, payment.StatusInProgress, p.Status)
	assert.Equal(t, "", p.TxID.String)
	assert.Equal(t, false, p.TxID.Valid)
	assert.Equal(t, "0.1", ub.FrozenAmount)
	assert.Equal(t, int64(1), ei.ExternalExchangeID.Int64)
	assert.Equal(t, "someID", ei.ExternalExchangeWithdrawID.String)
	assert.Equal(t, int64(23), ei.LastHandledID.Int64)

	time.Sleep(50 * time.Millisecond)
	paymentRepo.AssertExpectations(t)
	userBalanceService.AssertExpectations(t)
	userService.AssertExpectations(t)
	communicationService.AssertExpectations(t)
	mqttManager.AssertExpectations(t)
	externalExchangeService.AssertExpectations(t)
}

func TestService_UpdateWithdraw_UpdateStatus_InProgressAndNotAutoTransfer(t *testing.T) {
	qm := queryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	dbMock.ExpectExec("UPDATE crypto_payments").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("UPDATE crypto_payment_extra_info").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectCommit()

	paymentRepo := new(mocks.PaymentRepository)
	p := &payment.Payment{}
	paymentRepo.On("GetPaymentByIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		p = args.Get(2).(*payment.Payment)
		p.ID = 1
		p.UserID = 21
		p.CoinID = 1
		p.Code = "BTC"
		p.ToAddress = sql.NullString{String: "toAddress", Valid: true}
		p.Status = payment.StatusCreated
		p.FeeAmount = sql.NullString{String: "0.01", Valid: true}
		p.Amount = sql.NullString{String: "0.1", Valid: true}
		p.BlockchainNetwork = sql.NullString{String: "BTC", Valid: true}
	})
	ei := &payment.ExtraInfo{}
	paymentRepo.On("GetExtraInfoByPaymentIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ei = args.Get(2).(*payment.ExtraInfo)
		ei.ID = 1
	})

	currencyService := new(mocks.CurrencyService)
	walletService := new(mocks.WalletService)
	userConfigService := new(mocks.UserConfigService)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userService.On("GetUserByID", 21).Twice().Return(user.User{}, nil)
	userService.On("GetUserProfile", mock.Anything).Once().Return(user.Profile{}, nil)
	userBalanceService := new(mocks.UserBalanceService)
	ub := &userbalance.UserBalance{}
	userBalanceService.On("GetBalanceOfUserByCoinUsingTx", mock.Anything, 21, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub = args.Get(3).(*userbalance.UserBalance)
		ub.ID = 1
		ub.UserID = 21
		ub.Amount = "1.0"
		ub.FrozenAmount = "0.1"
	})

	communicationService := new(mocks.CommunicationService)
	communicationService.On("SendCryptoPaymentStatusUpdateEmail", mock.Anything, mock.Anything).Once().Return()
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	mqttManager.On("PublishPayment", mock.Anything, mock.Anything, mock.Anything).Once().Return()
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID: 23,
	}

	autoTransfer := false
	params := payment.UpdateWithdrawParams{
		ID:              1,
		Status:          payment.StatusInProgress,
		AdminStatus:     "",
		Fee:             "",
		NetworkFee:      "0.02",
		AutoTransfer:    &autoTransfer,
		RejectionReason: "",
		WithdrawType:    payment.WithdrawTypeHotWallet,
	}

	res, statusCode := paymentService.UpdateWithdraw(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	assert.Equal(t, payment.StatusCompleted, p.Status)
	assert.Equal(t, "", p.TxID.String)
	assert.Equal(t, false, p.TxID.Valid)
	assert.Equal(t, "0.90000000", ub.Amount)
	assert.Equal(t, "0.00000000", ub.FrozenAmount)

	assert.Equal(t, int64(23), ei.LastHandledID.Int64)

	time.Sleep(50 * time.Millisecond)
	paymentRepo.AssertExpectations(t)
	userBalanceService.AssertExpectations(t)
	userService.AssertExpectations(t)
	communicationService.AssertExpectations(t)
	mqttManager.AssertExpectations(t)
}

func TestService_UpdateDeposit_ShouldNotDepoist(t *testing.T) {
	qm := queryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	dbMock.ExpectExec("UPDATE crypto_payments").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("UPDATE crypto_payment_extra_info").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectCommit()
	paymentRepo := new(mocks.PaymentRepository)
	p := &payment.Payment{}
	paymentRepo.On("GetPaymentByIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		p = args.Get(2).(*payment.Payment)
		p.ID = 1
		p.UserID = 21
		p.CoinID = 1
		p.Code = "BTC"
		p.ToAddress = sql.NullString{String: "toAddress", Valid: true}
		p.Status = payment.StatusCreated
		p.FeeAmount = sql.NullString{String: "0.01", Valid: true}
		p.Amount = sql.NullString{String: "0.1", Valid: true}
		p.BlockchainNetwork = sql.NullString{String: "BTC", Valid: true}
		p.Type = payment.TypeDeposit
	})
	ei := &payment.ExtraInfo{}
	paymentRepo.On("GetExtraInfoByPaymentIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ei = args.Get(2).(*payment.ExtraInfo)
		ei.ID = 1
	})

	currencyService := new(mocks.CurrencyService)
	walletService := new(mocks.WalletService)
	userConfigService := new(mocks.UserConfigService)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)

	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID: 23,
	}

	params := payment.UpdateDepositParams{
		ID:            1,
		Status:        "",
		Amount:        "100.00",
		FromAddress:   "fromAddress",
		ToAddress:     "toAddress",
		TxID:          "txId",
		ShouldDeposit: nil,
	}

	res, statusCode := paymentService.UpdateDeposit(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	assert.Equal(t, payment.StatusCreated, p.Status)
	assert.Equal(t, true, p.TxID.Valid)
	assert.Equal(t, "txId", p.TxID.String)
	assert.Equal(t, "toAddress", p.ToAddress.String)
	assert.Equal(t, "fromAddress", p.FromAddress.String)

	assert.Equal(t, int64(23), ei.LastHandledID.Int64)
	paymentRepo.AssertExpectations(t)
}

func TestService_UpdateDeposit_ShouldDepoist(t *testing.T) {
	qm := queryMatcher{}
	sqlDb, dbMock, err := sqlmock.New(sqlmock.QueryMatcherOption(qm))
	if err != nil {
		t.Error(err)
	}
	defer sqlDb.Close()

	dialector := mysql.New(mysql.Config{
		DSN:                       "sqlmock_db_0",
		DriverName:                "mysql",
		Conn:                      sqlDb,
		SkipInitializeWithVersion: true,
	})
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Error(err)
	}

	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	dbMock.ExpectExec("UPDATE crypto_payments").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("UPDATE crypto_payment_extra_info").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(12, 1))
	dbMock.ExpectCommit()
	paymentRepo := new(mocks.PaymentRepository)
	p := &payment.Payment{}
	paymentRepo.On("GetPaymentByIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		p = args.Get(2).(*payment.Payment)
		p.ID = 1
		p.UserID = 21
		p.CoinID = 1
		p.Code = "BTC"
		p.ToAddress = sql.NullString{String: "toAddress", Valid: true}
		p.Status = payment.StatusCreated
		p.FeeAmount = sql.NullString{String: "0.01", Valid: true}
		p.Amount = sql.NullString{String: "0.1", Valid: true}
		p.BlockchainNetwork = sql.NullString{String: "BTC", Valid: true}
		p.Type = payment.TypeDeposit
	})
	ei := &payment.ExtraInfo{}
	paymentRepo.On("GetExtraInfoByPaymentIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ei = args.Get(2).(*payment.ExtraInfo)
		ei.ID = 1
	})

	currencyService := new(mocks.CurrencyService)
	walletService := new(mocks.WalletService)
	userConfigService := new(mocks.UserConfigService)
	twoFaManager := new(mocks.TwoFaManager)
	withdrawEmailConfirmationManager := new(mocks.WithdrawEmailConfirmationManager)
	permissionManager := new(mocks.UserPermissionManager)
	userWithdrawAddressService := new(mocks.UserWithdrawAddressService)
	userService := new(mocks.UserService)
	userBalanceService := new(mocks.UserBalanceService)
	ub := &userbalance.UserBalance{}
	userBalanceService.On("GetBalanceOfUserByCoinUsingTx", mock.Anything, 21, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub = args.Get(3).(*userbalance.UserBalance)
		ub.ID = 1
		ub.UserID = 21
		ub.Amount = "1.0"
		ub.FrozenAmount = "0.1"
	})

	communicationService := new(mocks.CommunicationService)
	priceGenerator := new(mocks.PriceGenerator)
	internalTransferService := new(mocks.InternalTransferService)
	externalExchangeService := new(mocks.ExternalExchangeService)
	autoExchangeManager := new(mocks.AutoExchangeManager)
	mqttManager := new(mocks.CentrifugoManager)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	paymentService := payment.NewPaymentService(db, paymentRepo, currencyService, walletService, userConfigService,
		twoFaManager, withdrawEmailConfirmationManager, permissionManager, userWithdrawAddressService, userService,
		userBalanceService, communicationService, priceGenerator, internalTransferService, externalExchangeService,
		autoExchangeManager, mqttManager, configs, logger)

	u := user.User{
		ID: 23,
	}

	shouldDeposit := true
	params := payment.UpdateDepositParams{
		ID:            1,
		Status:        "completed",
		Amount:        "0.1",
		FromAddress:   "fromAddress",
		ToAddress:     "toAddress",
		TxID:          "txId",
		ShouldDeposit: &shouldDeposit,
	}

	res, statusCode := paymentService.UpdateDeposit(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	assert.Equal(t, payment.StatusCompleted, p.Status)
	assert.Equal(t, true, p.TxID.Valid)
	assert.Equal(t, "txId", p.TxID.String)
	assert.Equal(t, "toAddress", p.ToAddress.String)
	assert.Equal(t, "fromAddress", p.FromAddress.String)
	assert.Equal(t, int64(23), ei.LastHandledID.Int64)

	assert.Equal(t, "1.10000000", ub.Amount)

	paymentRepo.AssertExpectations(t)
	userBalanceService.AssertExpectations(t)
}
