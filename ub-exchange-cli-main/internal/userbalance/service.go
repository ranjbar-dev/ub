package userbalance

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"exchange-go/internal/currency"
	"exchange-go/internal/platform"
	"exchange-go/internal/response"
	"exchange-go/internal/transaction"
	"exchange-go/internal/user"
	"exchange-go/internal/wallet"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/shopspring/decimal"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	StatusEnabled  = "ENABLED"
	StatusDisabled = "DISABLED"
)

type GetPairBalancesParams struct {
	PairID   int64  `form:"pair_currency_id"`
	PairName string `form:"pair_currency_name"`
}

type GetAllBalancesParams struct {
	Sort     string `form:"sort"`
	Name     string `form:"name"`
	CoinCode string `form:"code"`
}

type GetWithdrawDepositParams struct {
	Coin   string `form:"code" binding:"required"`
	Device string `form:"device"`
	Period string `form:"period"`
}

type UpdateUserBalanceFromAdminParams struct {
	ID     int64  `json:"id"`
	Amount string `json:"amount"`
}

type UpsertUserWalletBalancesParams struct {
	UserID             int
	CoinID             int64
	CoinCode           string
	BlockchainCoinID   int64
	BlockchainCoinCode string
	Address            string
}

type SetAutoExchangeCoinParams struct {
	Code             string `json:"code"`
	AutoExchangeCode string `json:"auto_exchange_code"`
	Mode             string `json:"mode"`
}

// Service provides the public API and internal operations for user balance management,
// including balance queries, address generation, auto-exchange configuration, and
// admin balance updates.
type Service interface {
	// GetPairBalances returns the user's balances for both coins in a trading pair.
	GetPairBalances(u *user.User, params GetPairBalancesParams) (apiResponse response.APIResponse, statusCode int)
	// GetAllBalances returns all of the user's coin balances with optional sorting and filtering.
	GetAllBalances(u *user.User, params GetAllBalancesParams) (apiResponse response.APIResponse, statusCode int)
	// GetWithdrawDepositData returns deposit address and withdrawal information for a specific coin.
	GetWithdrawDepositData(u *user.User, params GetWithdrawDepositParams) (apiResponse response.APIResponse, statusCode int)
	// SetAutoExchangeCoin configures the user's preferred auto-exchange target coin for deposits.
	SetAutoExchangeCoin(u *user.User, params SetAutoExchangeCoinParams) (apiResponse response.APIResponse, statusCode int)
	// GetBalanceOfUserByCoinID retrieves a user's balance for a specific coin.
	GetBalanceOfUserByCoinID(userID int, coinID int64, balance *UserBalance) error
	// GetBalanceOfUserByCoinUsingTx retrieves a user's balance within a database transaction.
	GetBalanceOfUserByCoinUsingTx(tx *gorm.DB, userID int, coinID int64, balance *UserBalance) error
	// GetBalancesOfUsersForCoins returns balances for multiple users and coins.
	GetBalancesOfUsersForCoins(userIds []int, coinIds []int64) []UserBalance
	// GetBalancesOfUsersForCoinsUsingTx returns balances for multiple users/coins within a transaction.
	GetBalancesOfUsersForCoinsUsingTx(tx *gorm.DB, userIds []int, coinIds []int64) []UserBalance
	// GenerateBalancesAndAddressForUser creates balance records and deposit addresses
	// for all active coins for a newly registered user.
	GenerateBalancesAndAddressForUser(u user.User)
	// GenerateSingleUserBalanceForCoin creates a balance record for one coin for the user.
	GenerateSingleUserBalanceForCoin(u user.User, coin currency.Coin) (ub UserBalance, err error)
	// GenerateAddress requests a new deposit address from the wallet service for the given balance.
	GenerateAddress(ub UserBalance, u user.User, coin currency.Coin) (address string, err error)
	// GetUserBalanceByCoinAndAddressUsingTx looks up a balance by coin and address within a transaction.
	GetUserBalanceByCoinAndAddressUsingTx(tx *gorm.DB, coinID int64, address string, ub *UserBalance) error
	// UpsertUserWalletBalance creates or updates the per-network wallet balance for a user.
	UpsertUserWalletBalance(params UpsertUserWalletBalancesParams) error

	//for admin
	// UpdateUserBalanceFromAdmin allows an admin to manually adjust a user's balance.
	UpdateUserBalanceFromAdmin(u *user.User, params UpdateUserBalanceFromAdminParams) (apiResponse response.APIResponse, statusCode int)
}

type service struct {
	db                          *gorm.DB
	repo                        Repository
	currencyService             currency.Service
	priceGenerator              currency.PriceGenerator
	permissionManager           user.PermissionManager
	walletService               wallet.Service
	configs                     platform.Configs
	logger                      platform.Logger
	userService                 user.Service
	userWalletBalanceRepository UserWalletBalanceRepository
}

type PartialBalance struct {
	CoinID   int64  `json:"currencyId"`
	CoinCode string `json:"currencyCode"`
	CoinName string `json:"currencyName"`
	Balance  string `json:"balance"`
}

type PairData struct {
	ID                 int64  `json:"id"`
	Name               string `json:"name"`
	MinimumOrderAmount string `json:"minimumOrderAmount"`
}
type GetPairBalanceResponse struct {
	PairBalances [2]PartialBalance  `json:"pairBalances"`
	PairData     PairData           `json:"pairData"`
	Fee          map[string]float64 `json:"fee"`
	Sum          string             `json:"sum"`
}

type AllBalancesFilters struct {
	CoinName string
	CoinCode string
}

type GetAllBalancesResponse struct {
	Balances               []BalancesResponse `json:"balances"`
	TotalSum               string             `json:"totalSum"`
	AvailableSum           string             `json:"availableSum"`
	InOrderSum             string             `json:"inOrderSum"`
	BtcTotalSum            string             `json:"btcTotalSum"`
	BtcAvailableSum        string             `json:"btcAvailableSum"`
	BtcInOrderSum          string             `json:"btcInOrderSum"`
	MinimumOfSmallBalances string             `json:"minimumOfSmallBalances"`
}

type BalancesResponse struct {
	TotalAmount                  string `json:"totalAmount"`
	AvailableAmount              string `json:"availableAmount"`
	InOrderAmount                string `json:"inOrderAmount"`
	EquivalentTotalAmount        string `json:"equivalentTotalAmount"`
	EquivalentAvailableAmount    string `json:"equivalentAvailableAmount"`
	EquivalentInOrderAmount      string `json:"equivalentInOrderAmount"`
	CoinCode                     string `json:"code"`
	CoinName                     string `json:"name"`
	Price                        string `json:"price"`
	Fee                          string `json:"fee"`
	SubUnit                      int    `json:"subUnit"`
	MinimumWithdraw              string `json:"minimumWithdraw"`
	BtcTotalEquivalentAmount     string `json:"btcTotalEquivalentAmount"`
	BtcAvailableEquivalentAmount string `json:"btcAvailableEquivalentAmount"`
	BtcInOrderEquivalentAmount   string `json:"btcInOrderEquivalentAmount"`
	Image                        string `json:"image"`
	BackgroundImage              string `json:"backgroundImage"`
	Address                      string `json:"address"`
	AutoExchangeCode             string `json:"autoExchangeCode"`
}

type NetworkConfig struct {
	Coin             string `json:"code"`
	SupportsWithdraw bool   `json:"supportsWithdraw"`
	SupportsDeposit  bool   `json:"supportsDeposit"`
	NetworkName      string `json:"completedNetworkName"`
	Address          string `json:"address"`
	Fee              string `json:"fee"`
}

type GetWithdrawDepositDataResponse struct {
	WalletAddress         string           `json:"walletAddress"`
	Balance               BalancesResponse `json:"balance"`
	SupportsWithdraw      bool             `json:"supportsWithdraw"`
	MainNetwork           string           `json:"mainNetwork"`
	CompletedNetworkName  string           `json:"completedNetworkName"`
	SupportsDeposit       bool             `json:"supportsDeposit"`
	HasDepositPermission  bool             `json:"isDepositPermissionGranted"`
	HasWithdrawPermission bool             `json:"isWithdrawPermissionGranted"`
	OtherNetworksConfigs  []NetworkConfig  `json:"otherNetworksConfigsAndAddresses"`
	NetworksConfigs       []NetworkConfig  `json:"networksConfigsAndAddresses"`
	DepositComments       []string         `json:"depositComments"`
	WithdrawComments      []string         `json:"withdrawComments"`
}

type addressAndPermissions struct {
	address               string
	hasDepositPermission  bool
	hasWithdrawPermission bool
}

func (s *service) GetPairBalances(u *user.User, params GetPairBalancesParams) (apiResponse response.APIResponse, statusCode int) {
	var pair currency.Pair
	var err error
	pairName := strings.ToUpper(strings.Trim(params.PairName, ""))
	if pairName != "" {
		pair, err = s.currencyService.GetPairByName(pairName)
	} else {
		if params.PairID > 0 {
			pair, err = s.currencyService.GetPairByID(params.PairID)
		}
	}

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("can not get pair by id or name", err,
			zap.String("service", "userBalanceService"),
			zap.String("method", "GetPairBalances"),
			zap.String("pairName", params.PairName),
			zap.Int64("pairID", params.PairID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || pair.ID == 0 {
		return response.Error("pair not found", http.StatusUnprocessableEntity, nil)
	}

	coinIds := []int64{
		pair.BasisCoinID,
		pair.DependentCoinID,
	}

	userBalances := s.repo.GetUserBalancesForCoins(u.ID, coinIds)

	var pairBalances [2]PartialBalance
	for _, ub := range userBalances {

		balanceDecimal, _ := decimal.NewFromString(ub.Amount)
		frozenBalanceDecimal, _ := decimal.NewFromString(ub.FrozenAmount)
		remainingBalance := balanceDecimal.Sub(frozenBalanceDecimal).StringFixed(8)
		pb := PartialBalance{
			CoinID:   ub.CoinID,
			CoinCode: ub.Coin.Code,
			CoinName: ub.Coin.Name,
			Balance:  remainingBalance,
		}

		//we put basis on first place and dependent on second place
		if pair.BasisCoinID == ub.CoinID {
			pairBalances[0] = pb
		} else {
			pairBalances[1] = pb
		}
	}

	pd := PairData{
		ID:                 pair.ID,
		Name:               pair.Name,
		MinimumOrderAmount: pair.MinimumOrderAmount.String,
	}

	fee := map[string]float64{
		"makerFee": pair.MakerFee,
		"takerFee": pair.TakerFee,
	}

	res := GetPairBalanceResponse{
		PairBalances: pairBalances,
		PairData:     pd,
		Fee:          fee,
		Sum:          "0.0", //wont be used just to have same structure as v1
	}

	return response.Success(res, "")

}

func (s *service) GetAllBalances(u *user.User, params GetAllBalancesParams) (apiResponse response.APIResponse, statusCode int) {
	ctx := context.Background()
	sortStrategy := strings.ToLower(strings.Trim(params.Sort, ""))
	filters := AllBalancesFilters{
		CoinName: strings.Trim(params.Name, ""),
		CoinCode: strings.Trim(params.CoinCode, ""),
	}

	userBalances := s.repo.GetUserAllBalances(u.ID, filters)

	totalSumDecimal := decimal.NewFromFloat(0)
	inOrderSumDecimal := decimal.NewFromFloat(0)
	availableSumDecimal := decimal.NewFromFloat(0)
	btcTotalSumDecimal := decimal.NewFromFloat(0)
	btcAvailableSumDecimal := decimal.NewFromFloat(0)
	btcInOrderSumDecimal := decimal.NewFromFloat(0)

	balances := make([]BalancesResponse, 0)

	for _, ub := range userBalances {
		r := s.GetAllBalanceResponse(ctx, ub, "")
		//we do not show the address if user did not verify his/her email
		if u.Status != user.StatusVerified {
			r.Address = ""
		}
		balances = append(balances, r)
		totalDecimal, _ := decimal.NewFromString(r.EquivalentTotalAmount)
		inOrderDecimal, _ := decimal.NewFromString(r.EquivalentInOrderAmount)
		availableDecimal, _ := decimal.NewFromString(r.EquivalentAvailableAmount)
		btcTotalDecimal, _ := decimal.NewFromString(r.BtcTotalEquivalentAmount)
		btcInOrderDecimal, _ := decimal.NewFromString(r.BtcInOrderEquivalentAmount)
		btcAvailableDecimal, _ := decimal.NewFromString(r.BtcAvailableEquivalentAmount)
		totalSumDecimal = totalSumDecimal.Add(totalDecimal)
		inOrderSumDecimal = inOrderSumDecimal.Add(inOrderDecimal)
		availableSumDecimal = availableSumDecimal.Add(availableDecimal)
		btcTotalSumDecimal = btcTotalSumDecimal.Add(btcTotalDecimal)
		btcInOrderSumDecimal = btcInOrderSumDecimal.Add(btcInOrderDecimal)
		btcAvailableSumDecimal = btcAvailableSumDecimal.Add(btcAvailableDecimal)
	}

	if sortStrategy == "" {
		sortStrategy = "desc"
	}

	if sortStrategy == "desc" || sortStrategy == "asc" {
		sort.Slice(balances, func(i, j int) bool {
			first, _ := strconv.ParseFloat(balances[i].EquivalentTotalAmount, 64)
			second, _ := strconv.ParseFloat(balances[j].EquivalentTotalAmount, 64)
			if sortStrategy == "desc" {
				return first > second
			}
			return first < second
		})
	} else {

	}

	finalRes := GetAllBalancesResponse{
		Balances:               balances,
		TotalSum:               totalSumDecimal.StringFixed(8),
		AvailableSum:           availableSumDecimal.StringFixed(8),
		InOrderSum:             inOrderSumDecimal.StringFixed(8),
		BtcTotalSum:            btcTotalSumDecimal.StringFixed(8),
		BtcAvailableSum:        btcAvailableSumDecimal.StringFixed(8),
		BtcInOrderSum:          btcInOrderSumDecimal.StringFixed(8),
		MinimumOfSmallBalances: "0", //todo handle this later
	}
	return response.Success(finalRes, "")

}

func (s *service) GetAllBalanceResponse(ctx context.Context, ub UserBalance, fee string) BalancesResponse {
	coin := ub.Coin
	coinCode := coin.Code

	ubAmountDecimal, _ := decimal.NewFromString(ub.Amount)
	ubFrozenAmountDecimal, _ := decimal.NewFromString(ub.FrozenAmount)
	ubAvailableDecimal := ubAmountDecimal.Sub(ubFrozenAmountDecimal)

	equivalentTotalAmount, _ := s.priceGenerator.GetAmountBasedOnUSDT(ctx, coinCode, ub.Amount)
	equivalentAvailableAmount, _ := s.priceGenerator.GetAmountBasedOnUSDT(ctx, coinCode, ubAvailableDecimal.StringFixed(8))
	equivalentInOrderAmount, _ := s.priceGenerator.GetAmountBasedOnUSDT(ctx, coinCode, ub.FrozenAmount)

	btcTotalEquivalentAmount, _ := s.priceGenerator.GetAmountBasedOnBTC(ctx, coinCode, ub.Amount)
	btcAvailableEquivalentAmount, _ := s.priceGenerator.GetAmountBasedOnBTC(ctx, coinCode, ubAvailableDecimal.StringFixed(8))
	btcInOrderEquivalentAmount, _ := s.priceGenerator.GetAmountBasedOnBTC(ctx, coinCode, ub.FrozenAmount)

	unitPrice, _ := s.priceGenerator.GetAmountBasedOnUSDT(ctx, coinCode, "1.0")
	if fee == "" {
		fee = "0"
	}
	imagePath := s.configs.GetImagePath()
	image := imagePath + coin.Image
	backgroundImage := imagePath + coin.Image //todo handle this later

	res := BalancesResponse{
		TotalAmount:                  ubAmountDecimal.StringFixed(8),
		AvailableAmount:              ubAvailableDecimal.StringFixed(8),
		InOrderAmount:                ubFrozenAmountDecimal.StringFixed(8),
		EquivalentTotalAmount:        equivalentTotalAmount,
		EquivalentAvailableAmount:    equivalentAvailableAmount,
		EquivalentInOrderAmount:      equivalentInOrderAmount,
		CoinCode:                     coin.Code,
		CoinName:                     coin.Name,
		Price:                        unitPrice,
		Fee:                          fee,
		SubUnit:                      coin.SubUnit,
		MinimumWithdraw:              coin.MinimumWithdraw,
		BtcTotalEquivalentAmount:     btcTotalEquivalentAmount,
		BtcAvailableEquivalentAmount: btcAvailableEquivalentAmount,
		BtcInOrderEquivalentAmount:   btcInOrderEquivalentAmount,
		Image:                        image,
		BackgroundImage:              backgroundImage,
		Address:                      ub.Address.String,
		AutoExchangeCode:             ub.AutoExchangeCoin.String,
	}

	return res

}

func (s *service) GetWithdrawDepositData(u *user.User, params GetWithdrawDepositParams) (apiResponse response.APIResponse, statusCode int) {

	//this API will return results for only users with verified email. so we check email verification first
	if u.Status == user.StatusRegistered {
		return response.Error("please verify your email first", http.StatusForbidden, nil)
	}

	coinCode := strings.ToUpper(strings.Trim(params.Coin, ""))
	var coin currency.Coin
	if coinCode == "" {
		return response.Error("coin not found", http.StatusUnprocessableEntity, nil)
	}
	coin, err := s.currencyService.GetCoinByCode(coinCode)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("can not get coin by code", err,
			zap.String("service", "userBalanceService"),
			zap.String("method", "GetWithdrawDepositData"),
			zap.String("coinCode", coinCode),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || coin.ID == 0 {
		return response.Error("coin not found", http.StatusUnprocessableEntity, nil)
	}

	ub := UserBalance{}
	err = s.GetBalanceOfUserByCoinID(u.ID, coin.ID, &ub)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("can not get user balance ", err,
			zap.String("service", "userBalanceService"),
			zap.String("method", "GetWithdrawDepositData"),
			zap.Int64("coinID", coin.ID),
			zap.Int("userID", u.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || ub.ID == 0 {
		ub, err = s.GenerateSingleUserBalanceForCoin(*u, coin)
		if err != nil {
			s.logger.Error2("can not generate single user balance ", err,
				zap.String("service", "userBalanceService"),
				zap.String("method", "GetWithdrawDepositData"),
				zap.Int64("coinID", coin.ID),
				zap.Int("userID", u.ID),
			)
			return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
		}
	}

	ub.Coin = coin
	addressAndPermissions, err := s.getAddressAndPermissions(ub, *u, coin)

	if err != nil {
		s.logger.Error2("can not get address and permission ", err,
			zap.String("service", "userBalanceService"),
			zap.String("method", "GetWithdrawDepositData"),
			zap.Int64("coinID", coin.ID),
			zap.Int("userID", u.ID),
			zap.Int64("userBalanceID", ub.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	mainNetwork := coinCode
	if coin.BlockchainNetwork.Valid {
		mainNetwork = coin.BlockchainNetwork.String
	}

	completedNetworkName := ""
	if coin.CompletedNetworkName.Valid {
		completedNetworkName = coin.CompletedNetworkName.String
	} else {
		completedNetworkName = coin.Name + "(" + coin.Code + ")"
	}

	if !ub.Address.Valid && addressAndPermissions.address != "" {
		ub.Address = sql.NullString{String: addressAndPermissions.address, Valid: true}
	}

	withdrawFee := coin.WithdrawalFee.Float64
	withdrawFeeString := strconv.FormatFloat(withdrawFee, 'f', 8, 64)

	withdrawSupported := false
	depositSupported := false
	if coin.SupportsWithdraw.Valid {
		withdrawSupported = coin.SupportsWithdraw.Bool
	}

	if coin.SupportsDeposit.Valid {
		depositSupported = coin.SupportsDeposit.Bool
	}

	ctx := context.Background()
	balanceData := s.GetAllBalanceResponse(ctx, ub, withdrawFeeString)
	networkConfigs, otherNetworkConfigs, err := s.getNetworkConfigs(ub, *u, coin)

	walletAddress := addressAndPermissions.address
	if strings.HasPrefix(walletAddress, "bitcoincash:") {
		walletAddress = strings.Replace(walletAddress, "bitcoincash:", "", -1)
	}

	var depositComments = make([]string, 0)
	if coin.DepositComments.Valid {
		depositComments = strings.Split(coin.DepositComments.String, "|")
	}
	var withdrawComments = make([]string, 0)
	if coin.WithdrawComments.Valid {
		withdrawComments = strings.Split(coin.WithdrawComments.String, "|")
	}

	res := GetWithdrawDepositDataResponse{
		WalletAddress:         walletAddress,
		Balance:               balanceData,
		SupportsWithdraw:      withdrawSupported,
		MainNetwork:           mainNetwork,
		CompletedNetworkName:  completedNetworkName,
		SupportsDeposit:       depositSupported,
		HasDepositPermission:  addressAndPermissions.hasDepositPermission,
		HasWithdrawPermission: addressAndPermissions.hasWithdrawPermission,
		OtherNetworksConfigs:  otherNetworkConfigs,
		NetworksConfigs:       networkConfigs,
		DepositComments:       depositComments,
		WithdrawComments:      withdrawComments,
	}

	return response.Success(res, "")
}

func (s *service) getNetworkConfigs(ub UserBalance, u user.User, coin currency.Coin) (all []NetworkConfig, others []NetworkConfig, err error) {
	coinCode := coin.Code
	if coin.BlockchainNetwork.Valid {
		coinCode = coin.BlockchainNetwork.String
	}

	withdrawFee := coin.WithdrawalFee.Float64
	withdrawFeeString := strconv.FormatFloat(withdrawFee, 'f', 8, 64)
	nc := NetworkConfig{
		Coin:             coinCode,
		SupportsWithdraw: coin.SupportsWithdraw.Bool,
		SupportsDeposit:  coin.SupportsWithdraw.Bool,
		NetworkName:      coin.CompletedNetworkName.String,
		Fee:              withdrawFeeString,
	}

	if nc.SupportsDeposit {
		nc.Address = ub.Address.String
		if strings.HasPrefix(nc.Address, "bitcoincash:") {
			nc.Address = strings.Replace(nc.Address, "bitcoincash:", "", -1)
		}
	}

	//we do not show the address if user did not verify his/her email
	if u.Status != user.StatusVerified {
		nc.Address = ""
	}
	all = append(all, nc)

	if !coin.OtherBlockchainNetworksConfigs.Valid {
		return all, others, nil
	}

	configs, _ := coin.GetOtherBlockchainNetworksConfigs()
	if len(configs) < 1 {
		return all, others, nil
	}

	var otherAddresses []OtherAddress
	if ub.OtherAddresses.Valid {
		otherAddresses, _ = ub.GetOtherAddresses()
	}

	for _, config := range configs {
		address := ""
		nc := NetworkConfig{
			Coin:             config.Code,
			SupportsWithdraw: config.SupportsWithdraw,
			SupportsDeposit:  config.SupportsWithdraw,
			NetworkName:      config.CompletedNetworkName,
			Fee:              config.Fee,
		}

		if nc.SupportsDeposit {
			for _, oa := range otherAddresses {
				if oa.Code == config.Code {
					address = oa.Address
				}
			}
			if address == "" {
				address, err = s.generateAddressForOtherNetwork(ub, u, config.Code)
				if err != nil {
					return all, others, err
				}
			}
			nc.Address = address
		}

		//we do not show the address if user did not verify his/her email
		if u.Status != user.StatusVerified {
			nc.Address = ""
		}

		all = append(all, nc)
		others = append(others, nc)
	}

	return all, others, nil

}

func (s *service) generateAddressForOtherNetwork(ub UserBalance, user user.User, network string) (address string, err error) {
	parentCoin, err := s.currencyService.GetCoinByCode(strings.ToUpper(network))
	if err != nil {
		return "", err
	}
	parentUb := UserBalance{}
	err = s.GetBalanceOfUserByCoinID(user.ID, parentCoin.ID, &parentUb)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return "", err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || parentUb.ID == 0 {
		parentUb, err = s.GenerateSingleUserBalanceForCoin(user, parentCoin)
		if err != nil {
			return "", err
		}
	}

	if parentUb.Address.Valid {
		address := parentUb.Address.String
		err := s.saveAddressInOtherAddresses(ub, parentCoin.Code, address)
		return parentUb.Address.String, err
	}

	address, err = s.walletService.GetAddressForUser(parentCoin.Code, user.PrivateChannelName)
	if err != nil {
		return address, err
	}

	parentUpdatingUb := &UserBalance{
		ID:      parentUb.ID,
		Address: sql.NullString{String: address, Valid: true},
	}

	err = s.db.Model(parentUpdatingUb).Updates(parentUpdatingUb).Error
	if err != nil {
		return address, err
	}

	otherAddresses, _ := ub.GetOtherAddresses()
	newOtherAddress := OtherAddress{
		Code:    parentCoin.Code,
		Address: address,
	}

	found := false
	for i, o := range otherAddresses {
		if o.Code == newOtherAddress.Code {
			otherAddresses[i].Address = newOtherAddress.Address
			found = true
		}

	}

	if !found {
		otherAddresses = append(otherAddresses, newOtherAddress)
	}

	otherAddressesBytes, err := json.Marshal(otherAddresses)
	updatingUb := &UserBalance{
		ID:             ub.ID,
		OtherAddresses: sql.NullString{String: string(otherAddressesBytes), Valid: true},
	}
	err = s.db.Model(updatingUb).Updates(updatingUb).Error

	return address, err

}

func (s *service) saveAddressInOtherAddresses(ub UserBalance, code string, address string) error {
	otherAddresses, _ := ub.GetOtherAddresses()
	newOtherAddress := OtherAddress{
		Code:    code,
		Address: address,
	}

	found := false
	for i, o := range otherAddresses {
		if o.Code == newOtherAddress.Code {
			otherAddresses[i].Address = newOtherAddress.Address
			found = true
		}

	}

	if !found {
		otherAddresses = append(otherAddresses, newOtherAddress)
	}

	otherAddressesBytes, err := json.Marshal(otherAddresses)
	updatingUb := &UserBalance{
		ID:             ub.ID,
		OtherAddresses: sql.NullString{String: string(otherAddressesBytes), Valid: true},
	}
	err = s.db.Model(updatingUb).Updates(updatingUb).Error

	return err
}

func (s *service) getAddressAndPermissions(ub UserBalance, u user.User, coin currency.Coin) (addressAndPermissions, error) {
	aap := addressAndPermissions{}
	//we do not show the address if user did not verify his/her email
	if u.Status != user.StatusVerified {
		return aap, nil
	}
	if s.permissionManager.IsPermissionGrantedToUserFor(u, user.PermissionDeposit) {
		aap.hasDepositPermission = true
		if ub.Address.Valid {
			aap.address = ub.Address.String
		} else {
			address, err := s.GenerateAddress(ub, u, coin)
			if err != nil {
				return aap, err
			}
			aap.address = address
		}
	}

	if s.permissionManager.IsPermissionGrantedToUserFor(u, user.PermissionWithdraw) {
		aap.hasWithdrawPermission = true
	}

	return aap, nil

}
func (s *service) SetAutoExchangeCoin(u *user.User, params SetAutoExchangeCoinParams) (apiResponse response.APIResponse, statusCode int) {
	code := strings.ToUpper(strings.Trim(params.Code, ""))
	autoExchangeCode := strings.ToUpper(strings.Trim(params.AutoExchangeCode, ""))
	mode := strings.Trim(params.Mode, "")
	coin, err := s.currencyService.GetCoinByCode(code)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("can not get coin by code", err,
			zap.String("service", "userBalanceService"),
			zap.String("method", "SetAutoExchangeCoin"),
			zap.String("code", code),
			zap.Int("userId", u.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return response.Error("code not found", http.StatusUnprocessableEntity, nil)
	}

	autoExchangeCoin, err := s.currencyService.GetCoinByCode(autoExchangeCode)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("can not get coin by code", err,
			zap.String("service", "userBalanceService"),
			zap.String("method", "SetAutoExchangeCoin"),
			zap.String("code", code),
			zap.Int("userId", u.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return response.Error("auto exchange code not found", http.StatusUnprocessableEntity, nil)
	}
	ub := &UserBalance{}
	err = s.repo.GetBalanceOfUserByCoinID(u.ID, coin.ID, ub)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		s.logger.Error2("can not get userbalance  by coinID", err,
			zap.String("service", "userBalanceService"),
			zap.String("method", "SetAutoExchangeCoin"),
			zap.String("code", code),
			zap.Int("userId", u.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return response.Error("user balance not found", http.StatusUnprocessableEntity, nil)
	}
	if mode == "delete" {
		if ub.AutoExchangeCoin.Valid && ub.AutoExchangeCoin.String != "" {
			ub.AutoExchangeCoin = sql.NullString{String: "", Valid: false}
			err := s.db.Omit(clause.Associations).Save(ub).Error
			if err != nil {
				s.logger.Error2("can not save userbalance", err,
					zap.String("service", "userBalanceService"),
					zap.String("method", "SetAutoExchangeCoin"),
					zap.String("code", code),
					zap.Int64("userBalanceId", ub.ID),
					zap.Int("userId", u.ID),
					zap.String("mode", mode),
				)
				return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
			}
		}
		return response.Success(nil, "")
	}
	//here means the mod is add
	if coin.ID == autoExchangeCoin.ID {
		return response.Error("the code and exchange code can not be the same", http.StatusUnprocessableEntity, nil)
	}
	allPairs := s.currencyService.GetActivePairCurrenciesList()
	pairFound := false
	for _, pair := range allPairs {
		if pair.BasisCoinID == coin.ID && pair.DependentCoinID == autoExchangeCoin.ID {
			pairFound = true
			break
		}
		if pair.BasisCoinID == autoExchangeCoin.ID && pair.DependentCoinID == coin.ID {
			pairFound = true
			break
		}
	}
	if !pairFound {
		return response.Error("pair not found with these two coins", http.StatusUnprocessableEntity, nil)
	}
	ub.AutoExchangeCoin = sql.NullString{String: autoExchangeCode, Valid: true}
	err = s.db.Omit(clause.Associations).Save(ub).Error
	if err != nil {
		s.logger.Error2("can not save userbalance", err,
			zap.String("service", "userBalanceService"),
			zap.String("method", "SetAutoExchangeCoin"),
			zap.String("code", code),
			zap.Int64("userBalanceId", ub.ID),
			zap.Int("userId", u.ID),
			zap.String("mode", mode),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}
	return response.Success(nil, "")
}

func (s *service) GenerateAddress(ub UserBalance, u user.User, coin currency.Coin) (address string, err error) {
	if coin.BlockchainNetwork.Valid {
		parentCoin, err := s.currencyService.GetCoinByCode(coin.BlockchainNetwork.String)
		if err != nil {
			return address, err
		}
		parentUserBalance := UserBalance{}
		err = s.GetBalanceOfUserByCoinID(u.ID, parentCoin.ID, &parentUserBalance)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return address, err
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			parentUserBalance, err = s.GenerateSingleUserBalanceForCoin(u, parentCoin)
			if err != nil {
				return address, err
			}
		}

		if parentUserBalance.Address.Valid {
			updatingUb := &UserBalance{
				ID:      ub.ID,
				Address: sql.NullString{String: parentUserBalance.Address.String, Valid: true},
			}
			err := s.db.Model(updatingUb).Updates(updatingUb).Error
			address = parentUserBalance.Address.String
			return address, err
		}
		//here we should generate address for parent coin and child coin
		address, err := s.walletService.GetAddressForUser(parentCoin.Code, u.PrivateChannelName)
		if err != nil {
			return address, err
		}
		parentUpdatingUb := &UserBalance{
			ID:      parentUserBalance.ID,
			Address: sql.NullString{String: address, Valid: true},
		}

		err = s.db.Model(parentUpdatingUb).Updates(parentUpdatingUb).Error
		if err != nil {
			return address, err
		}

		childUpdatingUb := &UserBalance{
			ID:      ub.ID,
			Address: sql.NullString{String: address, Valid: true},
		}
		err = s.db.Model(childUpdatingUb).Updates(childUpdatingUb).Error
		if err != nil {
			return address, err
		}
		return address, nil
	}

	//the coin does not have parent so we just generate address for it
	address, err = s.walletService.GetAddressForUser(coin.Code, u.PrivateChannelName)
	if err != nil {
		return address, err
	}
	updatingUb := &UserBalance{
		ID:      ub.ID,
		Address: sql.NullString{String: address, Valid: true},
	}
	err = s.db.Model(updatingUb).Updates(updatingUb).Error
	if err != nil {
		return address, err
	}
	return address, nil
}

func (s *service) GenerateSingleUserBalanceForCoin(u user.User, coin currency.Coin) (ub UserBalance, err error) {
	ub = UserBalance{
		UserID:        u.ID,
		CoinID:        coin.ID,
		FrozenBalance: "0.0",
		BalanceCoin:   coin.Code,
		Status:        StatusEnabled,
		Address:       sql.NullString{String: "", Valid: false},
		Amount:        "0.0",
		FrozenAmount:  "0.0",
	}

	err = s.db.Omit(clause.Associations).Create(&ub).Error
	ub.Coin = coin
	return ub, err
}

func (s *service) GetBalanceOfUserByCoinID(userID int, coinID int64, ba *UserBalance) error {
	return s.repo.GetBalanceOfUserByCoinID(userID, coinID, ba)
}

func (s *service) GetBalanceOfUserByCoinUsingTx(tx *gorm.DB, userID int, coinID int64, ba *UserBalance) error {
	return s.repo.GetBalanceOfUserByCoinIDUsingTx(tx, userID, coinID, ba)
}

func (s *service) GetBalancesOfUsersForCoins(userIds []int, coinIds []int64) []UserBalance {
	return s.repo.GetBalancesOfUsersForCoins(userIds, coinIds)
}

func (s *service) GetBalancesOfUsersForCoinsUsingTx(tx *gorm.DB, userIds []int, coinIds []int64) []UserBalance {
	return s.repo.GetBalancesOfUsersForCoinsUsingTx(tx, userIds, coinIds)
}

func (s *service) GenerateBalancesAndAddressForUser(u user.User) {
	coins := s.currencyService.GetActiveCoins()
	for _, coin := range coins {
		ub, err := s.GenerateSingleUserBalanceForCoin(u, coin)
		if err != nil {
			s.logger.Error2("can not generate single user balance", err,
				zap.String("service", "userBalanceService"),
				zap.String("method", "GenerateBalancesAndAddressForUser"),
				zap.Int64("coinID", coin.ID),
				zap.Int("userID", u.ID),
			)
			continue
		}
		_, err = s.GenerateAddress(ub, u, coin)
		if err != nil {
			s.logger.Error2("can not generate address", err,
				zap.String("service", "userBalanceService"),
				zap.String("method", "GenerateBalancesAndAddressForUser"),
				zap.Int64("coinID", coin.ID),
				zap.Int("userID", u.ID),
				zap.Int64("userBalanceID", ub.ID),
			)
		}

	}
}

func (s *service) GetUserBalanceByCoinAndAddressUsingTx(tx *gorm.DB, coinID int64, address string, ub *UserBalance) error {
	return s.repo.GetUserBalanceByCoinAndAddressUsingTx(tx, coinID, address, ub)
}

func (s *service) UpdateUserBalanceFromAdmin(u *user.User, params UpdateUserBalanceFromAdminParams) (apiResponse response.APIResponse, statusCode int) {
	amountDecimal, err := decimal.NewFromString(params.Amount)
	if err != nil || amountDecimal.IsNegative() {
		return response.Error("amount is not correct", http.StatusUnprocessableEntity, nil)
	}
	tx := s.db.Begin()
	err = tx.Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("error starting transaction", err,
			zap.String("service", "userBalanceService"),
			zap.String("method", "UpdateUserBalanceFromAdmin"),
			zap.Int64("userBalanceId", params.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}
	ub := &UserBalance{}
	err = s.repo.GetUserBalanceByIDUsingTx(tx, params.ID, ub)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		tx.Rollback()
		s.logger.Error2("error getting user balance from db", err,
			zap.String("service", "userBalanceService"),
			zap.String("method", "UpdateUserBalanceFromAdmin"),
			zap.Int64("userBalanceId", params.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)

	}
	if errors.Is(err, gorm.ErrRecordNotFound) || ub.ID == 0 {
		tx.Rollback()
		return response.Error("deposit not found", http.StatusUnprocessableEntity, nil)
	}
	userBalanceAmountDecimal, err := decimal.NewFromString(ub.Amount)
	if err != nil {
		tx.Rollback()
		return response.Error("can not get user balance amount", http.StatusUnprocessableEntity, nil)
	}
	diffAmountDecimal := amountDecimal.Sub(userBalanceAmountDecimal)
	if diffAmountDecimal.IsZero() {
		tx.Rollback()
		return response.Success(nil, "")
	}
	ub.Amount = amountDecimal.StringFixed(8)
	err = tx.Omit(clause.Associations).Save(ub).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("error updating user balance", err,
			zap.String("service", "userBalanceService"),
			zap.String("method", "UpdateUserBalanceFromAdmin"),
			zap.Int64("userBalanceId", params.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}
	transactionType := transaction.TypeAdminAddition

	if diffAmountDecimal.IsNegative() {
		transactionType = transaction.TypeAdminReduction
	}

	coinCode := ub.Coin.Code
	if coinCode == "" {
		tx.Rollback()
		s.logger.Error2("coin in user balance is not loaded", err,
			zap.String("service", "userBalanceService"),
			zap.String("method", "UpdateUserBalanceFromAdmin"),
			zap.Int64("userBalanceId", params.ID),
		)
		return response.Error("coin in user balance is not loaded", http.StatusUnprocessableEntity, nil)

	}
	adminTransaction := &transaction.Transaction{
		UserID:    ub.UserID,
		CoinID:    ub.CoinID,
		OrderID:   sql.NullInt64{Int64: 0, Valid: false},
		Type:      transactionType,
		Amount:    sql.NullString{String: diffAmountDecimal.Abs().StringFixed(8), Valid: true},
		CoinName:  coinCode,
		PaymentID: sql.NullInt64{Int64: 0, Valid: false},
	}
	err = tx.Omit(clause.Associations).Save(adminTransaction).Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("error saving withdraw transaction", err,
			zap.String("service", "paymentService"),
			zap.String("method", "UpdateUserBalanceFromAdmin"),
			zap.Int64("paymentId", params.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}

	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		s.logger.Error2("error committing db transaction", err,
			zap.String("service", "paymentService"),
			zap.String("method", "UpdateUserBalanceFromAdmin"),
			zap.Int64("paymentId", params.ID),
		)
		return response.Error("something went wrong", http.StatusUnprocessableEntity, nil)
	}
	res := make(map[string]interface{})
	return response.Success(res, "")
}

func (s *service) UpsertUserWalletBalance(params UpsertUserWalletBalancesParams) error {
	balance, err := s.walletService.GetAddressBalance(params.CoinCode, params.BlockchainCoinCode, params.Address, true)
	if err != nil {
		return err
	}

	userWalletBalance := &UserWalletBalance{}
	err = s.userWalletBalanceRepository.FindUserWalletBalance(params.UserID, params.CoinID, params.BlockchainCoinID, userWalletBalance)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if errors.Is(err, gorm.ErrRecordNotFound) || userWalletBalance.ID == 0 {
		//insert new
		userWalletBalance.UserID = params.UserID
		userWalletBalance.CoinID = params.CoinID
		userWalletBalance.NetworkCoinID = sql.NullInt64{
			Int64: params.BlockchainCoinID,
			Valid: true,
		}
		userWalletBalance.Balance = balance
		err = s.db.Save(userWalletBalance).Error
		if err != nil {
			return err
		}
	} else {
		//update
		userWalletBalance.Balance = balance
		err := s.db.Save(userWalletBalance).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func NewBalanceService(db *gorm.DB, repo Repository, currencyService currency.Service, pg currency.PriceGenerator,
	permissionManager user.PermissionManager, walletService wallet.Service, userService user.Service,
	userWalletBalanceRepository UserWalletBalanceRepository, configs platform.Configs, logger platform.Logger) Service {
	return &service{
		db:                          db,
		repo:                        repo,
		currencyService:             currencyService,
		priceGenerator:              pg,
		permissionManager:           permissionManager,
		walletService:               walletService,
		userService:                 userService,
		userWalletBalanceRepository: userWalletBalanceRepository,
		configs:                     configs,
		logger:                      logger,
	}
}
