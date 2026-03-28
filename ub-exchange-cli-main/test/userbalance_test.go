package test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"exchange-go/internal/api"
	"exchange-go/internal/di"
	"exchange-go/internal/response"
	"exchange-go/internal/transaction"
	"exchange-go/internal/user"
	"exchange-go/internal/userbalance"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UserBalanceTests struct {
	*suite.Suite
	httpServer      http.Handler
	adminHTTPServer http.Handler
	db              *gorm.DB
	redisClient     *redis.Client
	userActor       *userActor
	adminUserActor  *userActor
}

func (t *UserBalanceTests) SetupSuite() {
	container := getContainer()
	t.httpServer = container.Get(di.HTTPServer).(api.HTTPServer).GetEngine()
	t.adminHTTPServer = container.Get(di.HTTPServer).(api.HTTPServer).GetAdminEngine()
	t.db = getDb()
	t.redisClient = getRedis()
	t.userActor = getUserActor()
	t.adminUserActor = getAdminUserActor()
}

func (t *UserBalanceTests) SetupTest() {
	t.db.Where("user_id = ?", t.userActor.ID).Delete(userbalance.UserBalance{})
	t.db.Where("user_id = ?", t.userActor.ID).Delete(user.UsersPermissions{})
}

func (t *UserBalanceTests) TearDownTest() {
	t.db.Where("user_id = ?", t.userActor.ID).Delete(userbalance.UserBalance{})
}

func (t *UserBalanceTests) TearDownSuite() {

}

func (t *UserBalanceTests) TestGetPairBalances() {
	usdtUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "10000.00",
		FrozenAmount:  "1000.00",
	}

	btcUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        2, //for btc from currency seed
		FrozenBalance: "",
		BalanceCoin:   "BTC",
		Status:        userbalance.StatusEnabled,
		Amount:        "0.1",
		FrozenAmount:  "0",
	}
	t.db.Create(usdtUb)
	t.db.Create(btcUb)

	queryParams := url.Values{}
	queryParams.Set("pair_currency_name", "BTC-USDT")

	paramsString := queryParams.Encode()

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/user-balance/pair-balance?"+paramsString, nil)
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    userbalance.GetPairBalanceResponse
	}{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)
	pairBalances := result.Data.PairBalances

	var usdtBalance userbalance.PartialBalance
	var btcBalance userbalance.PartialBalance
	if pairBalances[0].CoinCode == "USDT" {
		usdtBalance = pairBalances[0]
		btcBalance = pairBalances[1]
	} else {
		usdtBalance = pairBalances[1]
		btcBalance = pairBalances[0]
	}

	assert.Equal(t.T(), "Tether", usdtBalance.CoinName)
	assert.Equal(t.T(), int64(1), usdtBalance.CoinID)
	assert.Equal(t.T(), "USDT", usdtBalance.CoinCode)
	assert.Equal(t.T(), "9000.00000000", usdtBalance.Balance)

	assert.Equal(t.T(), "Bitcoin", btcBalance.CoinName)
	assert.Equal(t.T(), int64(2), btcBalance.CoinID)
	assert.Equal(t.T(), "BTC", btcBalance.CoinCode)
	assert.Equal(t.T(), "0.10000000", btcBalance.Balance)

	fee := result.Data.Fee

	assert.Equal(t.T(), float64(0.2), fee["makerFee"])
	assert.Equal(t.T(), float64(0.3), fee["takerFee"])

	pairData := result.Data.PairData
	assert.Equal(t.T(), int64(1), pairData.ID)
	assert.Equal(t.T(), "10", pairData.MinimumOrderAmount)
	assert.Equal(t.T(), "BTC-USDT", pairData.Name)

}

func (t *UserBalanceTests) TestGetAllBalances() {
	ctx := context.Background()
	t.redisClient.HMSet(ctx, "live_data:pair_currency:BTC-USDT", "price", "50000")
	t.redisClient.HMSet(ctx, "live_data:pair_currency:ETH-USDT", "price", "2000")
	t.redisClient.HMSet(ctx, "live_data:pair_currency:ETH-BTC", "price", "0.04")
	t.redisClient.HMSet(ctx, "live_data:pair_currency:GRS-BTC", "price", "0.00002")
	t.redisClient.HMSet(ctx, "live_data:pair_currency:USDT-DAI", "price", "1")
	t.redisClient.HMSet(ctx, "live_data:pair_currency:BTC-DAI", "price", "50000")

	usdtUb := &userbalance.UserBalance{
		UserID:       t.userActor.ID,
		CoinID:       1, //for USDT from currency seed
		BalanceCoin:  "USDT",
		Status:       userbalance.StatusEnabled,
		Amount:       "10000.00",
		FrozenAmount: "1000.00",
		Address:      sql.NullString{String: "USDTAddress", Valid: true},
	}

	btcUb := &userbalance.UserBalance{
		UserID:           t.userActor.ID,
		CoinID:           2, //for BTC from currency seed
		BalanceCoin:      "BTC",
		Status:           userbalance.StatusEnabled,
		Amount:           "0.1",
		FrozenAmount:     "0",
		Address:          sql.NullString{String: "BTCAddress", Valid: true},
		AutoExchangeCoin: sql.NullString{String: "ETH", Valid: true},
	}

	ethUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        3, //for ETH from currency seed
		FrozenBalance: "",
		BalanceCoin:   "BTC",
		Status:        userbalance.StatusEnabled,
		Amount:        "1.0",
		FrozenAmount:  "0.2",
		Address:       sql.NullString{String: "ETHAddress", Valid: true},
	}

	grsUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        4, //for GRS from currency seed
		FrozenBalance: "",
		BalanceCoin:   "GRS",
		Status:        userbalance.StatusEnabled,
		Amount:        "30.0",
		FrozenAmount:  "0",
		Address:       sql.NullString{String: "GRSAddress", Valid: true},
	}

	daiUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        5, //for DAI from currency seed
		FrozenBalance: "",
		BalanceCoin:   "DAI",
		Status:        userbalance.StatusEnabled,
		Amount:        "3000.0",
		FrozenAmount:  "500.00",
		Address:       sql.NullString{String: "DAIAddress", Valid: true},
	}

	t.db.Create(usdtUb)
	t.db.Create(btcUb)
	t.db.Create(ethUb)
	t.db.Create(grsUb)
	t.db.Create(daiUb)

	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/user-balance/balance", nil)
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    userbalance.GetAllBalancesResponse
	}{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)

	assert.Equal(t.T(), "20030.00000000", result.Data.TotalSum)
	assert.Equal(t.T(), "18130.00000000", result.Data.AvailableSum)
	assert.Equal(t.T(), "1900.00000000", result.Data.InOrderSum)
	assert.Equal(t.T(), "0.40060000", result.Data.BtcTotalSum)
	assert.Equal(t.T(), "0.36260000", result.Data.BtcAvailableSum)
	assert.Equal(t.T(), "0.03800000", result.Data.BtcInOrderSum)
	balances := result.Data.Balances
	for _, b := range balances {
		switch b.CoinCode {
		case "USDT":
			assert.Equal(t.T(), "10000.00000000", b.TotalAmount)
			assert.Equal(t.T(), "9000.00000000", b.AvailableAmount)
			assert.Equal(t.T(), "1000.00000000", b.InOrderAmount)
			assert.Equal(t.T(), "10000.00", b.EquivalentTotalAmount)
			assert.Equal(t.T(), "9000.00000000", b.EquivalentAvailableAmount)
			assert.Equal(t.T(), "1000.00", b.EquivalentInOrderAmount)
			assert.Equal(t.T(), "Tether", b.CoinName)
			assert.Equal(t.T(), "1.0", b.Price)
			assert.Equal(t.T(), "0", b.Fee)
			assert.Equal(t.T(), 6, b.SubUnit)
			assert.Equal(t.T(), "10", b.MinimumWithdraw)
			assert.Equal(t.T(), "0.20000000", b.BtcTotalEquivalentAmount)
			assert.Equal(t.T(), "0.18000000", b.BtcAvailableEquivalentAmount)
			assert.Equal(t.T(), "0.02000000", b.BtcInOrderEquivalentAmount)
			assert.Equal(t.T(), "USDTAddress", b.Address)
		case "BTC":
			assert.Equal(t.T(), "0.10000000", b.TotalAmount)
			assert.Equal(t.T(), "0.10000000", b.AvailableAmount)
			assert.Equal(t.T(), "0.00000000", b.InOrderAmount)
			assert.Equal(t.T(), "5000.00000000", b.EquivalentTotalAmount)
			assert.Equal(t.T(), "5000.00000000", b.EquivalentAvailableAmount)
			assert.Equal(t.T(), "0.00000000", b.EquivalentInOrderAmount)
			assert.Equal(t.T(), "Bitcoin", b.CoinName)
			assert.Equal(t.T(), "50000.00000000", b.Price)
			assert.Equal(t.T(), "0", b.Fee)
			assert.Equal(t.T(), 8, b.SubUnit)
			assert.Equal(t.T(), "0.001", b.MinimumWithdraw)
			assert.Equal(t.T(), "0.1", b.BtcTotalEquivalentAmount)
			assert.Equal(t.T(), "0.10000000", b.BtcAvailableEquivalentAmount)
			assert.Equal(t.T(), "0", b.BtcInOrderEquivalentAmount)
			assert.Equal(t.T(), "BTCAddress", b.Address)
			assert.Equal(t.T(), "ETH", b.AutoExchangeCode)
		case "ETH":
			assert.Equal(t.T(), "1.00000000", b.TotalAmount)
			assert.Equal(t.T(), "0.80000000", b.AvailableAmount)
			assert.Equal(t.T(), "0.20000000", b.InOrderAmount)
			assert.Equal(t.T(), "2000.00000000", b.EquivalentTotalAmount)
			assert.Equal(t.T(), "1600.00000000", b.EquivalentAvailableAmount)
			assert.Equal(t.T(), "400.00000000", b.EquivalentInOrderAmount)
			assert.Equal(t.T(), "Ethereum", b.CoinName)
			assert.Equal(t.T(), "2000.00000000", b.Price)
			assert.Equal(t.T(), "0", b.Fee)
			assert.Equal(t.T(), 8, b.SubUnit)
			assert.Equal(t.T(), "0.001", b.MinimumWithdraw)
			assert.Equal(t.T(), "0.04000000", b.BtcTotalEquivalentAmount)
			assert.Equal(t.T(), "0.03200000", b.BtcAvailableEquivalentAmount)
			assert.Equal(t.T(), "0.00800000", b.BtcInOrderEquivalentAmount)
			assert.Equal(t.T(), "ETHAddress", b.Address)
		case "GRS":
			assert.Equal(t.T(), "30.00000000", b.TotalAmount)
			assert.Equal(t.T(), "30.00000000", b.AvailableAmount)
			assert.Equal(t.T(), "0.00000000", b.InOrderAmount)
			assert.Equal(t.T(), "30.00000000", b.EquivalentTotalAmount)
			assert.Equal(t.T(), "30.00000000", b.EquivalentAvailableAmount)
			assert.Equal(t.T(), "0.00000000", b.EquivalentInOrderAmount)
			assert.Equal(t.T(), "Groestlcoin", b.CoinName)
			assert.Equal(t.T(), "1.00000000", b.Price)
			assert.Equal(t.T(), "0", b.Fee)
			assert.Equal(t.T(), 8, b.SubUnit)
			assert.Equal(t.T(), "10.0", b.MinimumWithdraw)
			assert.Equal(t.T(), "0.00060000", b.BtcTotalEquivalentAmount)
			assert.Equal(t.T(), "0.00060000", b.BtcAvailableEquivalentAmount)
			assert.Equal(t.T(), "0.00000000", b.BtcInOrderEquivalentAmount)
			assert.Equal(t.T(), "GRSAddress", b.Address)
		case "DAI":
			assert.Equal(t.T(), "3000.00000000", b.TotalAmount)
			assert.Equal(t.T(), "2500.00000000", b.AvailableAmount)
			assert.Equal(t.T(), "500.00000000", b.InOrderAmount)
			assert.Equal(t.T(), "3000.00000000", b.EquivalentTotalAmount)
			assert.Equal(t.T(), "2500.00000000", b.EquivalentAvailableAmount)
			assert.Equal(t.T(), "500.00000000", b.EquivalentInOrderAmount)
			assert.Equal(t.T(), "Dai", b.CoinName)
			assert.Equal(t.T(), "1.00000000", b.Price)
			assert.Equal(t.T(), "0", b.Fee)
			assert.Equal(t.T(), 6, b.SubUnit)
			assert.Equal(t.T(), "10", b.MinimumWithdraw)
			assert.Equal(t.T(), "0.06000000", b.BtcTotalEquivalentAmount)
			assert.Equal(t.T(), "0.05000000", b.BtcAvailableEquivalentAmount)
			assert.Equal(t.T(), "0.01000000", b.BtcInOrderEquivalentAmount)
			assert.Equal(t.T(), "DAIAddress", b.Address)
		default:
			t.Fail("we should not be in default case")
		}
	}

}

func (t *UserBalanceTests) TestGetWithdrawDepositData_BTC_NoFormerUserBalance() {
	up1 := user.UsersPermissions{
		UserID:           t.userActor.ID,
		UserPermissionID: 2, //see the userPermissionSeed id for withdraw is 3
	}
	up2 := user.UsersPermissions{
		UserID:           t.userActor.ID,
		UserPermissionID: 1, //see the userPermissionSeed id for deposit is 3
	}
	t.db.Create(&up1)
	t.db.Create(&up2)

	ctx := context.Background()
	t.redisClient.HMSet(ctx, "live_data:pair_currency:BTC-USDT", "price", "50000")

	queryParams := url.Values{}
	queryParams.Set("code", "BTC")
	paramsString := queryParams.Encode()
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/user-balance/withdraw-deposit?"+paramsString, nil)
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    userbalance.GetWithdrawDepositDataResponse
	}{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)

	assert.Equal(t.T(), true, result.Data.SupportsWithdraw)
	assert.Equal(t.T(), true, result.Data.SupportsDeposit)
	assert.Equal(t.T(), true, result.Data.HasWithdrawPermission)
	assert.Equal(t.T(), true, result.Data.HasDepositPermission)
	assert.Equal(t.T(), "BTCAddress", result.Data.WalletAddress)
	assert.Equal(t.T(), 1, len(result.Data.NetworksConfigs))
	assert.Equal(t.T(), 0, len(result.Data.OtherNetworksConfigs))

	assert.Equal(t.T(), true, result.Data.NetworksConfigs[0].SupportsDeposit)
	assert.Equal(t.T(), true, result.Data.NetworksConfigs[0].SupportsWithdraw)
	assert.Equal(t.T(), "BTC", result.Data.NetworksConfigs[0].Coin)
	assert.Equal(t.T(), "BTCAddress", result.Data.NetworksConfigs[0].Address)

	balance := result.Data.Balance
	assert.Equal(t.T(), "0.00000000", balance.TotalAmount)
	assert.Equal(t.T(), "0.00000000", balance.AvailableAmount)
	assert.Equal(t.T(), "0.00000000", balance.InOrderAmount)
	assert.Equal(t.T(), "0.00000000", balance.EquivalentTotalAmount)
	assert.Equal(t.T(), "0.00000000", balance.EquivalentAvailableAmount)
	assert.Equal(t.T(), "0.00000000", balance.EquivalentInOrderAmount)
	assert.Equal(t.T(), "BTC", balance.CoinCode)
	assert.Equal(t.T(), "Bitcoin", balance.CoinName)
	assert.Equal(t.T(), "50000.00000000", balance.Price)
	assert.Equal(t.T(), "0.00010000", balance.Fee)
	assert.Equal(t.T(), 8, balance.SubUnit)
	assert.Equal(t.T(), "0.001", balance.MinimumWithdraw)
	assert.Equal(t.T(), "0.0", balance.BtcTotalEquivalentAmount)
	assert.Equal(t.T(), "0.00000000", balance.BtcAvailableEquivalentAmount)
	assert.Equal(t.T(), "0.0", balance.BtcInOrderEquivalentAmount)
	assert.Equal(t.T(), "BTCAddress", balance.Address)

	assert.IsType(t.T(), []string{}, result.Data.DepositComments)
	assert.IsType(t.T(), []string{}, result.Data.WithdrawComments)

	//checking if user balance is inserted in db
	btcUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{UserID: t.userActor.ID, CoinID: 2}).Find(btcUb).Error
	assert.Nil(t.T(), err)
	assert.Greater(t.T(), btcUb.ID, int64(0))
	assert.Equal(t.T(), "BTCAddress", btcUb.Address.String)
	assert.Equal(t.T(), "0.0", btcUb.Amount)
	assert.Equal(t.T(), "0.0", btcUb.FrozenAmount)

}

func (t *UserBalanceTests) TestGetWithdrawDepositData_BTC_NoFormerAddress() {
	up1 := user.UsersPermissions{
		UserID:           t.userActor.ID,
		UserPermissionID: 2, //see the userPermissionSeed id for withdraw is 3
	}
	up2 := user.UsersPermissions{
		UserID:           t.userActor.ID,
		UserPermissionID: 1, //see the userPermissionSeed id for deposit is 3
	}
	t.db.Create(&up1)
	t.db.Create(&up2)

	btcUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        2, //for btc from currency seed
		FrozenBalance: "",
		BalanceCoin:   "BTC",
		Status:        userbalance.StatusEnabled,
		Amount:        "0",
		FrozenAmount:  "0",
	}
	t.db.Create(btcUb)

	ctx := context.Background()
	t.redisClient.HMSet(ctx, "live_data:pair_currency:BTC-USDT", "price", "50000")

	queryParams := url.Values{}
	queryParams.Set("code", "BTC")
	paramsString := queryParams.Encode()
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/user-balance/withdraw-deposit?"+paramsString, nil)
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    userbalance.GetWithdrawDepositDataResponse
	}{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)

	assert.Equal(t.T(), true, result.Data.SupportsWithdraw)
	assert.Equal(t.T(), true, result.Data.SupportsDeposit)
	assert.Equal(t.T(), true, result.Data.HasWithdrawPermission)
	assert.Equal(t.T(), true, result.Data.HasDepositPermission)
	assert.Equal(t.T(), "BTCAddress", result.Data.WalletAddress)
	assert.Equal(t.T(), 1, len(result.Data.NetworksConfigs))
	assert.Equal(t.T(), 0, len(result.Data.OtherNetworksConfigs))

	assert.Equal(t.T(), true, result.Data.NetworksConfigs[0].SupportsDeposit)
	assert.Equal(t.T(), true, result.Data.NetworksConfigs[0].SupportsWithdraw)
	assert.Equal(t.T(), "BTC", result.Data.NetworksConfigs[0].Coin)
	assert.Equal(t.T(), "BTCAddress", result.Data.NetworksConfigs[0].Address)

	balanceResult := result.Data.Balance
	assert.Equal(t.T(), "0.00000000", balanceResult.TotalAmount)
	assert.Equal(t.T(), "0.00000000", balanceResult.AvailableAmount)
	assert.Equal(t.T(), "0.00000000", balanceResult.InOrderAmount)
	assert.Equal(t.T(), "0.00000000", balanceResult.EquivalentTotalAmount)
	assert.Equal(t.T(), "0.00000000", balanceResult.EquivalentAvailableAmount)
	assert.Equal(t.T(), "0.00000000", balanceResult.EquivalentInOrderAmount)
	assert.Equal(t.T(), "BTC", balanceResult.CoinCode)
	assert.Equal(t.T(), "Bitcoin", balanceResult.CoinName)
	assert.Equal(t.T(), "50000.00000000", balanceResult.Price)
	assert.Equal(t.T(), "0.00010000", balanceResult.Fee)
	assert.Equal(t.T(), 8, balanceResult.SubUnit)
	assert.Equal(t.T(), "0.001", balanceResult.MinimumWithdraw)
	assert.Equal(t.T(), "0", balanceResult.BtcTotalEquivalentAmount)
	assert.Equal(t.T(), "0.00000000", balanceResult.BtcAvailableEquivalentAmount)
	assert.Equal(t.T(), "0", balanceResult.BtcInOrderEquivalentAmount)
	assert.Equal(t.T(), "BTCAddress", balanceResult.Address)

	assert.IsType(t.T(), []string{}, result.Data.DepositComments)
	assert.IsType(t.T(), []string{}, result.Data.WithdrawComments)

	//checking if user balance address is updated in db
	updatedUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: btcUb.ID}).Find(updatedUb).Error
	assert.Nil(t.T(), err)
	assert.Greater(t.T(), btcUb.ID, int64(0))
	assert.Equal(t.T(), "BTCAddress", updatedUb.Address.String)

}

func (t *UserBalanceTests) TestGetWithdrawDepositData_USDT_NoFormerBalance() {
	up1 := user.UsersPermissions{
		UserID:           t.userActor.ID,
		UserPermissionID: 2, //see the userPermissionSeed id for withdraw is 3
	}
	up2 := user.UsersPermissions{
		UserID:           t.userActor.ID,
		UserPermissionID: 1, //see the userPermissionSeed id for deposit is 3
	}

	t.db.Create(&up1)
	t.db.Create(&up2)

	ctx := context.Background()
	t.redisClient.HMSet(ctx, "live_data:pair_currency:BTC-USDT", "price", "50000")

	queryParams := url.Values{}
	queryParams.Set("code", "USDT")
	paramsString := queryParams.Encode()
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/user-balance/withdraw-deposit?"+paramsString, nil)
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    userbalance.GetWithdrawDepositDataResponse
	}{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)

	assert.Equal(t.T(), true, result.Data.SupportsWithdraw)
	assert.Equal(t.T(), true, result.Data.SupportsDeposit)
	assert.Equal(t.T(), true, result.Data.HasWithdrawPermission)
	assert.Equal(t.T(), true, result.Data.HasDepositPermission)
	assert.Equal(t.T(), "ETHAddress", result.Data.WalletAddress)
	assert.Equal(t.T(), 2, len(result.Data.NetworksConfigs))
	assert.Equal(t.T(), 1, len(result.Data.OtherNetworksConfigs))

	assert.Equal(t.T(), true, result.Data.NetworksConfigs[0].SupportsDeposit)
	assert.Equal(t.T(), true, result.Data.NetworksConfigs[0].SupportsWithdraw)
	assert.Equal(t.T(), "ETH", result.Data.NetworksConfigs[0].Coin)
	assert.Equal(t.T(), "2.00000000", result.Data.NetworksConfigs[0].Fee)
	assert.Equal(t.T(), "ETHAddress", result.Data.NetworksConfigs[0].Address)

	assert.Equal(t.T(), true, result.Data.NetworksConfigs[1].SupportsDeposit)
	assert.Equal(t.T(), true, result.Data.NetworksConfigs[1].SupportsWithdraw)
	assert.Equal(t.T(), "TRX", result.Data.NetworksConfigs[1].Coin)
	assert.Equal(t.T(), "2.5", result.Data.NetworksConfigs[1].Fee)
	assert.Equal(t.T(), "TRXAddress", result.Data.NetworksConfigs[1].Address)

	assert.Equal(t.T(), true, result.Data.OtherNetworksConfigs[0].SupportsDeposit)
	assert.Equal(t.T(), true, result.Data.OtherNetworksConfigs[0].SupportsWithdraw)
	assert.Equal(t.T(), "TRX", result.Data.OtherNetworksConfigs[0].Coin)
	assert.Equal(t.T(), "2.5", result.Data.OtherNetworksConfigs[0].Fee)
	assert.Equal(t.T(), "TRXAddress", result.Data.OtherNetworksConfigs[0].Address)

	balanceResult := result.Data.Balance
	assert.Equal(t.T(), "0.00000000", balanceResult.TotalAmount)
	assert.Equal(t.T(), "0.00000000", balanceResult.AvailableAmount)
	assert.Equal(t.T(), "0.00000000", balanceResult.InOrderAmount)
	assert.Equal(t.T(), "0.0", balanceResult.EquivalentTotalAmount)
	assert.Equal(t.T(), "0.00000000", balanceResult.EquivalentAvailableAmount)
	assert.Equal(t.T(), "0.0", balanceResult.EquivalentInOrderAmount)
	assert.Equal(t.T(), "USDT", balanceResult.CoinCode)
	assert.Equal(t.T(), "Tether", balanceResult.CoinName)
	assert.Equal(t.T(), "1.0", balanceResult.Price)
	assert.Equal(t.T(), "2.00000000", balanceResult.Fee)
	assert.Equal(t.T(), 6, balanceResult.SubUnit)
	assert.Equal(t.T(), "10", balanceResult.MinimumWithdraw)
	assert.Equal(t.T(), "0.00000000", balanceResult.BtcTotalEquivalentAmount)
	assert.Equal(t.T(), "0.00000000", balanceResult.BtcAvailableEquivalentAmount)
	assert.Equal(t.T(), "0.00000000", balanceResult.BtcInOrderEquivalentAmount)
	assert.Equal(t.T(), "ETHAddress", balanceResult.Address)

	//checking if user balance is inserted in db  for usdt
	usdtUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{UserID: t.userActor.ID, CoinID: 1}).Find(usdtUb).Error
	assert.Nil(t.T(), err)
	assert.Greater(t.T(), usdtUb.ID, int64(0))
	assert.Equal(t.T(), "ETHAddress", usdtUb.Address.String)
	assert.Equal(t.T(), "[{\"code\":\"TRX\",\"address\":\"TRXAddress\"}]", usdtUb.OtherAddresses.String)
	assert.Equal(t.T(), "0.0", usdtUb.Amount)
	assert.Equal(t.T(), "0.0", usdtUb.FrozenAmount)

	//checking if user balance is inserted in db for eth
	ethUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{UserID: t.userActor.ID, CoinID: 3}).Find(ethUb).Error
	assert.Nil(t.T(), err)
	assert.Greater(t.T(), ethUb.ID, int64(0))
	assert.Equal(t.T(), "ETHAddress", ethUb.Address.String)
	assert.Equal(t.T(), "", ethUb.OtherAddresses.String)
	assert.Equal(t.T(), "0.0", ethUb.Amount)
	assert.Equal(t.T(), "0.0", ethUb.FrozenAmount)

	//checking if user balance is inserted in db for trx
	trxUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{UserID: t.userActor.ID, CoinID: 6}).Find(trxUb).Error
	assert.Nil(t.T(), err)
	assert.Greater(t.T(), trxUb.ID, int64(0))
	assert.Equal(t.T(), "TRXAddress", trxUb.Address.String)
	assert.Equal(t.T(), "", trxUb.OtherAddresses.String)
	assert.Equal(t.T(), "0.0", trxUb.Amount)
	assert.Equal(t.T(), "0.0", trxUb.FrozenAmount)

	assert.IsType(t.T(), []string{}, result.Data.DepositComments)
	assert.IsType(t.T(), []string{}, result.Data.WithdrawComments)

}

func (t *UserBalanceTests) TestGetWithdrawDepositData_USDT_NoFormerAddressAndOtherAddresses() {
	up1 := user.UsersPermissions{
		UserID:           t.userActor.ID,
		UserPermissionID: 2, //see the userPermissionSeed id for withdraw is 3
	}
	up2 := user.UsersPermissions{
		UserID:           t.userActor.ID,
		UserPermissionID: 1, //see the userPermissionSeed id for deposit is 3
	}
	t.db.Create(&up1)
	t.db.Create(&up2)

	usdtUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "0",
		FrozenAmount:  "0",
	}
	t.db.Create(usdtUb)

	ctx := context.Background()
	t.redisClient.HMSet(ctx, "live_data:pair_currency:BTC-USDT", "price", "50000")

	queryParams := url.Values{}
	queryParams.Set("code", "USDT")
	paramsString := queryParams.Encode()
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/user-balance/withdraw-deposit?"+paramsString, nil)
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    userbalance.GetWithdrawDepositDataResponse
	}{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)

	assert.Equal(t.T(), true, result.Data.SupportsWithdraw)
	assert.Equal(t.T(), true, result.Data.SupportsDeposit)
	assert.Equal(t.T(), true, result.Data.HasWithdrawPermission)
	assert.Equal(t.T(), true, result.Data.HasDepositPermission)
	assert.Equal(t.T(), "ETHAddress", result.Data.WalletAddress)
	assert.Equal(t.T(), 2, len(result.Data.NetworksConfigs))
	assert.Equal(t.T(), 1, len(result.Data.OtherNetworksConfigs))

	assert.Equal(t.T(), true, result.Data.NetworksConfigs[0].SupportsDeposit)
	assert.Equal(t.T(), true, result.Data.NetworksConfigs[0].SupportsWithdraw)
	assert.Equal(t.T(), "ETH", result.Data.NetworksConfigs[0].Coin)
	assert.Equal(t.T(), "2.00000000", result.Data.NetworksConfigs[0].Fee)
	assert.Equal(t.T(), "ETHAddress", result.Data.NetworksConfigs[0].Address)

	assert.Equal(t.T(), true, result.Data.NetworksConfigs[1].SupportsDeposit)
	assert.Equal(t.T(), true, result.Data.NetworksConfigs[1].SupportsWithdraw)
	assert.Equal(t.T(), "TRX", result.Data.NetworksConfigs[1].Coin)
	assert.Equal(t.T(), "2.5", result.Data.NetworksConfigs[1].Fee)
	assert.Equal(t.T(), "TRXAddress", result.Data.NetworksConfigs[1].Address)

	assert.Equal(t.T(), true, result.Data.OtherNetworksConfigs[0].SupportsDeposit)
	assert.Equal(t.T(), true, result.Data.OtherNetworksConfigs[0].SupportsWithdraw)
	assert.Equal(t.T(), "TRX", result.Data.OtherNetworksConfigs[0].Coin)
	assert.Equal(t.T(), "2.5", result.Data.OtherNetworksConfigs[0].Fee)
	assert.Equal(t.T(), "TRXAddress", result.Data.OtherNetworksConfigs[0].Address)

	balanceResult := result.Data.Balance
	assert.Equal(t.T(), "0.00000000", balanceResult.TotalAmount)
	assert.Equal(t.T(), "0.00000000", balanceResult.AvailableAmount)
	assert.Equal(t.T(), "0.00000000", balanceResult.InOrderAmount)
	assert.Equal(t.T(), "0", balanceResult.EquivalentTotalAmount)
	assert.Equal(t.T(), "0.00000000", balanceResult.EquivalentAvailableAmount)
	assert.Equal(t.T(), "0", balanceResult.EquivalentInOrderAmount)
	assert.Equal(t.T(), "USDT", balanceResult.CoinCode)
	assert.Equal(t.T(), "Tether", balanceResult.CoinName)
	assert.Equal(t.T(), "1.0", balanceResult.Price)
	assert.Equal(t.T(), "2.00000000", balanceResult.Fee)
	assert.Equal(t.T(), 6, balanceResult.SubUnit)
	assert.Equal(t.T(), "10", balanceResult.MinimumWithdraw)
	assert.Equal(t.T(), "0.00000000", balanceResult.BtcTotalEquivalentAmount)
	assert.Equal(t.T(), "0.00000000", balanceResult.BtcAvailableEquivalentAmount)
	assert.Equal(t.T(), "0.00000000", balanceResult.BtcInOrderEquivalentAmount)
	assert.Equal(t.T(), "ETHAddress", balanceResult.Address)

	//checking if user balance is updated
	updatingUsdtUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).Find(updatingUsdtUb).Error
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), "ETHAddress", updatingUsdtUb.Address.String)
	assert.Equal(t.T(), "0", updatingUsdtUb.Amount)
	assert.Equal(t.T(), "[{\"code\":\"TRX\",\"address\":\"TRXAddress\"}]", updatingUsdtUb.OtherAddresses.String)
	assert.Equal(t.T(), "0", updatingUsdtUb.FrozenAmount)

	//checking if user balance is inserted in db for eth
	ethUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{UserID: t.userActor.ID, CoinID: 3}).Find(ethUb).Error
	assert.Nil(t.T(), err)
	assert.Greater(t.T(), ethUb.ID, int64(0))
	assert.Equal(t.T(), "ETHAddress", ethUb.Address.String)
	assert.Equal(t.T(), "", ethUb.OtherAddresses.String)
	assert.Equal(t.T(), "0.0", ethUb.Amount)
	assert.Equal(t.T(), "0.0", ethUb.FrozenAmount)

	//checking if user balance is inserted in db for trx
	trxUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{UserID: t.userActor.ID, CoinID: 6}).Find(trxUb).Error
	assert.Nil(t.T(), err)
	assert.Greater(t.T(), trxUb.ID, int64(0))
	assert.Equal(t.T(), "TRXAddress", trxUb.Address.String)
	assert.Equal(t.T(), "", trxUb.OtherAddresses.String)
	assert.Equal(t.T(), "0.0", trxUb.Amount)
	assert.Equal(t.T(), "0.0", trxUb.FrozenAmount)

	assert.IsType(t.T(), []string{}, result.Data.DepositComments)
	assert.IsType(t.T(), []string{}, result.Data.WithdrawComments)
}

func (t *UserBalanceTests) TestGetWithdrawDepositData_USDT_FormerAddressForEthAndTrx() {
	up1 := user.UsersPermissions{
		UserID:           t.userActor.ID,
		UserPermissionID: 2, //see the userPermissionSeed id for withdraw is 3
	}
	up2 := user.UsersPermissions{
		UserID:           t.userActor.ID,
		UserPermissionID: 1, //see the userPermissionSeed id for deposit is 3
	}

	t.db.Create(&up1)
	t.db.Create(&up2)

	usdtUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "0",
		FrozenAmount:  "0",
	}

	ethUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        3, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "ETH",
		Status:        userbalance.StatusEnabled,
		Amount:        "0",
		FrozenAmount:  "0",
		Address:       sql.NullString{String: "ETHAddress1", Valid: true},
	}

	trxUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        6, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "0",
		FrozenAmount:  "0",
		Address:       sql.NullString{String: "TRXAddress1", Valid: true},
	}

	t.db.Create(usdtUb)
	t.db.Create(ethUb)
	t.db.Create(trxUb)

	ctx := context.Background()
	t.redisClient.HMSet(ctx, "live_data:pair_currency:BTC-USDT", "price", "50000")

	queryParams := url.Values{}
	queryParams.Set("code", "USDT")
	paramsString := queryParams.Encode()
	res := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/user-balance/withdraw-deposit?"+paramsString, nil)
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	result := struct {
		Status  bool
		Message string
		Data    userbalance.GetWithdrawDepositDataResponse
	}{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), http.StatusOK, res.Code)

	assert.Equal(t.T(), true, result.Data.SupportsWithdraw)
	assert.Equal(t.T(), true, result.Data.SupportsDeposit)
	assert.Equal(t.T(), true, result.Data.HasWithdrawPermission)
	assert.Equal(t.T(), true, result.Data.HasDepositPermission)
	assert.Equal(t.T(), "ETHAddress1", result.Data.WalletAddress)
	assert.Equal(t.T(), 2, len(result.Data.NetworksConfigs))
	assert.Equal(t.T(), 1, len(result.Data.OtherNetworksConfigs))

	assert.Equal(t.T(), true, result.Data.NetworksConfigs[0].SupportsDeposit)
	assert.Equal(t.T(), true, result.Data.NetworksConfigs[0].SupportsWithdraw)
	assert.Equal(t.T(), "ETH", result.Data.NetworksConfigs[0].Coin)
	assert.Equal(t.T(), "2.00000000", result.Data.NetworksConfigs[0].Fee)
	assert.Equal(t.T(), "ETHAddress1", result.Data.NetworksConfigs[0].Address)

	assert.Equal(t.T(), true, result.Data.NetworksConfigs[1].SupportsDeposit)
	assert.Equal(t.T(), true, result.Data.NetworksConfigs[1].SupportsWithdraw)
	assert.Equal(t.T(), "TRX", result.Data.NetworksConfigs[1].Coin)
	assert.Equal(t.T(), "2.5", result.Data.NetworksConfigs[1].Fee)
	assert.Equal(t.T(), "TRXAddress1", result.Data.NetworksConfigs[1].Address)

	assert.Equal(t.T(), true, result.Data.OtherNetworksConfigs[0].SupportsDeposit)
	assert.Equal(t.T(), true, result.Data.OtherNetworksConfigs[0].SupportsWithdraw)
	assert.Equal(t.T(), "TRX", result.Data.OtherNetworksConfigs[0].Coin)
	assert.Equal(t.T(), "2.5", result.Data.OtherNetworksConfigs[0].Fee)
	assert.Equal(t.T(), "TRXAddress1", result.Data.OtherNetworksConfigs[0].Address)

	balanceResult := result.Data.Balance
	assert.Equal(t.T(), "0.00000000", balanceResult.TotalAmount)
	assert.Equal(t.T(), "0.00000000", balanceResult.AvailableAmount)
	assert.Equal(t.T(), "0.00000000", balanceResult.InOrderAmount)
	assert.Equal(t.T(), "0", balanceResult.EquivalentTotalAmount)
	assert.Equal(t.T(), "0.00000000", balanceResult.EquivalentAvailableAmount)
	assert.Equal(t.T(), "0", balanceResult.EquivalentInOrderAmount)
	assert.Equal(t.T(), "USDT", balanceResult.CoinCode)
	assert.Equal(t.T(), "Tether", balanceResult.CoinName)
	assert.Equal(t.T(), "1.0", balanceResult.Price)
	assert.Equal(t.T(), "2.00000000", balanceResult.Fee)
	assert.Equal(t.T(), 6, balanceResult.SubUnit)
	assert.Equal(t.T(), "10", balanceResult.MinimumWithdraw)
	assert.Equal(t.T(), "0.00000000", balanceResult.BtcTotalEquivalentAmount)
	assert.Equal(t.T(), "0.00000000", balanceResult.BtcAvailableEquivalentAmount)
	assert.Equal(t.T(), "0.00000000", balanceResult.BtcInOrderEquivalentAmount)
	assert.Equal(t.T(), "ETHAddress1", balanceResult.Address)

	//checking if user balance is updated
	updatingUsdtUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).Find(updatingUsdtUb).Error
	assert.Nil(t.T(), err)
	assert.Equal(t.T(), "ETHAddress1", updatingUsdtUb.Address.String)
	assert.Equal(t.T(), "0", updatingUsdtUb.Amount)
	assert.Equal(t.T(), "[{\"code\":\"TRX\",\"address\":\"TRXAddress1\"}]", updatingUsdtUb.OtherAddresses.String)
	assert.Equal(t.T(), "0", updatingUsdtUb.FrozenAmount)

	assert.IsType(t.T(), []string{}, result.Data.DepositComments)
	assert.IsType(t.T(), []string{}, result.Data.WithdrawComments)
}

func (t *UserBalanceTests) TestUpdateUserBalanceForAdmin_MoreThanCurrentBalance() {
	user := getNewUserActor()
	usdtUb := &userbalance.UserBalance{
		UserID:        user.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "100.00",
		FrozenAmount:  "50.00",
		Address:       sql.NullString{String: "user1USDTAddress", Valid: true},
	}
	err := t.db.Create(usdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}

	res := httptest.NewRecorder()
	data := fmt.Sprintf(`{
		"id" : %d,
		"amount" : "110.00"
	}`, usdtUb.ID)
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/user-balance/update", bytes.NewReader(body))
	token := "Bearer " + t.adminUserActor.Token
	req.Header.Set("Authorization", token)
	t.adminHTTPServer.ServeHTTP(res, req)

	assert.Equal(t.T(), http.StatusOK, res.Code)

	result := response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.T().Error(err)
	}
	updatedUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(updatedUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "110.00000000", updatedUb.Amount)

	tx := &transaction.Transaction{}
	err = t.db.Where(transaction.Transaction{Type: "ADMIN_ADDITION", UserID: user.ID}).First(tx).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "10.00000000", tx.Amount.String)
	assert.Equal(t.T(), int64(1), tx.CoinID)
	assert.Equal(t.T(), "USDT", tx.CoinName)
}

func (t *UserBalanceTests) TestUpdateUserBalanceForAdmin_LessThanCurrentBalance() {
	user := getNewUserActor()
	usdtUb := &userbalance.UserBalance{
		UserID:        user.ID,
		CoinID:        1, //for usdt from currency seed
		FrozenBalance: "",
		BalanceCoin:   "USDT",
		Status:        userbalance.StatusEnabled,
		Amount:        "100.00",
		FrozenAmount:  "50.00",
		Address:       sql.NullString{String: "user1USDTAddress", Valid: true},
	}
	err := t.db.Create(usdtUb).Error
	if err != nil {
		t.Fail(err.Error())
	}

	res := httptest.NewRecorder()
	data := fmt.Sprintf(`{
		"id" : %d,
		"amount" : "80.00"
	}`, usdtUb.ID)
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/user-balance/update", bytes.NewReader(body))
	token := "Bearer " + t.adminUserActor.Token
	req.Header.Set("Authorization", token)
	t.adminHTTPServer.ServeHTTP(res, req)

	assert.Equal(t.T(), http.StatusOK, res.Code)

	result := response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.T().Error(err)
	}
	updatedUb := &userbalance.UserBalance{}
	err = t.db.Where(userbalance.UserBalance{ID: usdtUb.ID}).First(updatedUb).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "80.00000000", updatedUb.Amount)

	tx := &transaction.Transaction{}
	err = t.db.Where(transaction.Transaction{Type: "ADMIN_REDUCTION", UserID: user.ID}).First(tx).Error
	if err != nil {
		t.Fail(err.Error())
	}
	assert.Equal(t.T(), "20.00000000", tx.Amount.String)
	assert.Equal(t.T(), int64(1), tx.CoinID)
	assert.Equal(t.T(), "USDT", tx.CoinName)
}

func (t *UserBalanceTests) TestSetAutoExchange_CodeAndExchangeCodeAndPairNotFound() {
	//api call with code not exist
	res := httptest.NewRecorder()
	data := `{
		"code" : "dsaa",
		"auto_exchange_code" : "BTC",
		"mode" : "add"
	}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/user-balance/auto-exchange", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)

	result := response.APIResponse{}
	err := json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.T().Error(err)
	}
	assert.Equal(t.T(), false, result.Status)
	assert.Equal(t.T(), "code not found", result.Message)

	//api call with exchange code not exist
	res = httptest.NewRecorder()
	data = `{
		"code" : "BTC",
		"auto_exchange_code" : "dsa",
		"mode" : "add"
	}`
	body = []byte(data)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/user-balance/auto-exchange", bytes.NewReader(body))
	token = "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)

	result = response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.T().Error(err)
	}
	assert.Equal(t.T(), false, result.Status)
	assert.Equal(t.T(), "auto exchange code not found", result.Message)

	//api call with the two coins that we do not have pair for

	trxUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        6, //TRX from currency seed
		FrozenBalance: "",
		BalanceCoin:   "",
		Status:        "",
		Amount:        "",
		FrozenAmount:  "",
	}
	err = t.db.Create(trxUb).Error
	if err != nil {
		t.Fail(err.Error())
	}

	//in currency seed we do not have pair for TRX and DAI
	res = httptest.NewRecorder()
	data = `{
		"code" : "TRX",
		"auto_exchange_code" : "DAI",
		"mode" : "add"
	}`
	body = []byte(data)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/user-balance/auto-exchange", bytes.NewReader(body))
	token = "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	assert.Equal(t.T(), http.StatusUnprocessableEntity, res.Code)

	result = response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.T().Error(err)
	}
	assert.Equal(t.T(), false, result.Status)
	assert.Equal(t.T(), "pair not found with these two coins", result.Message)
}

func (t *UserBalanceTests) TestSetAutoExchange_Add_Then_Delete_BothSuccessful() {
	usdtUb := &userbalance.UserBalance{
		UserID:        t.userActor.ID,
		CoinID:        1, //USDT from currency seed
		FrozenBalance: "",
		BalanceCoin:   "",
		Status:        "",
		Amount:        "",
		FrozenAmount:  "",
	}
	err := t.db.Create(usdtUb).Error
	if err != nil {
		t.Fail(err.Error())

	}

	res := httptest.NewRecorder()
	data := `{
		"code" : "USDT",
		"auto_exchange_code" : "BTC",
		"mode" : "add"
	}`
	body := []byte(data)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/user-balance/auto-exchange", bytes.NewReader(body))
	token := "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	assert.Equal(t.T(), http.StatusOK, res.Code)

	result := response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.T().Error(err)
	}
	assert.Equal(t.T(), true, result.Status)
	assert.Equal(t.T(), "", result.Message)

	updatedUSDTUb := &userbalance.UserBalance{}
	err = t.db.Where("id = ?", usdtUb.ID).First(updatedUSDTUb).Error
	if err != nil {
		t.T().Error(err)
	}

	assert.Equal(t.T(), true, updatedUSDTUb.AutoExchangeCoin.Valid)
	assert.Equal(t.T(), "BTC", updatedUSDTUb.AutoExchangeCoin.String)

	//calling the api with delete
	res = httptest.NewRecorder()
	data = `{
		"code" : "USDT",
		"auto_exchange_code" : "BTC",
		"mode" : "delete"
	}`
	body = []byte(data)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/user-balance/auto-exchange", bytes.NewReader(body))
	token = "Bearer " + t.userActor.Token
	req.Header.Set("Authorization", token)
	t.httpServer.ServeHTTP(res, req)

	assert.Equal(t.T(), http.StatusOK, res.Code)

	result = response.APIResponse{}
	err = json.Unmarshal(res.Body.Bytes(), &result)
	if err != nil {
		t.T().Error(err)
	}
	assert.Equal(t.T(), true, result.Status)
	assert.Equal(t.T(), "", result.Message)

	updatedUSDTUb = &userbalance.UserBalance{}
	err = t.db.Where("id = ?", usdtUb.ID).First(updatedUSDTUb).Error
	if err != nil {
		t.T().Error(err)
	}

	assert.Equal(t.T(), false, updatedUSDTUb.AutoExchangeCoin.Valid)
	assert.Equal(t.T(), "", updatedUSDTUb.AutoExchangeCoin.String)
}

func TestUserBalance(t *testing.T) {
	suite.Run(t, &UserBalanceTests{
		Suite: new(suite.Suite),
	})
}
