// Package userbalance_test tests the user balance service. Covers:
//   - Retrieving pair balances by pair ID and pair name
//   - Fetching all user balances with USDT-equivalent calculations
//   - Withdraw/deposit data retrieval for BTC and multi-network coins (USDT ERC20/TRC20)
//   - Balance lookup by coin ID and bulk balance queries for multiple users
//   - Generating balances and wallet addresses for new users
//   - Upserting wallet balances (existing and new records)
//   - Admin balance updates with negative, over-limit, and under-limit amounts
//   - Auto-exchange coin configuration (add/remove with validation)
//
// Test data: mock repositories, currency service, wallet service, price generator,
// and go-sqlmock for GORM database interactions with flexible query matching.
package userbalance_test

import (
	"database/sql"
	"exchange-go/internal/currency"
	"exchange-go/internal/mocks"
	"exchange-go/internal/user"
	"exchange-go/internal/userbalance"
	"net/http"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// userBalanceQueryMatcher is a sqlmock.QueryMatcher that accepts any SQL query,
// allowing tests to focus on service logic rather than exact SQL statements.
type userBalanceQueryMatcher struct {
}

// Match satisfies the sqlmock.QueryMatcher interface by unconditionally matching
// any expected SQL against any actual SQL.
func (userBalanceQueryMatcher) Match(expectedSQL, actualSQL string) error {
	return nil
}

func TestService_GetPairBalances_ByPairID(t *testing.T) {
	qm := userBalanceQueryMatcher{}
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
	db, err := gorm.Open(dialector, &gorm.Config{SkipDefaultTransaction: true})
	dbMock.MatchExpectationsInOrder(false)

	cs := new(mocks.CurrencyService)

	pair := currency.Pair{
		ID:                 1,
		Name:               "BTC-USDT",
		BasisCoinID:        1,
		DependentCoinID:    2,
		MinimumOrderAmount: sql.NullString{String: "10", Valid: true},
		MakerFee:           0.9,
		TakerFee:           0.9,
	}

	cs.On("GetPairByID", int64(1)).Once().Return(pair, nil)

	usdtCoin := currency.Coin{
		ID:   1,
		Code: "USDT",
		Name: "Tether",
	}

	btcCoin := currency.Coin{
		ID:   2,
		Code: "BTC",
		Name: "Bitcoin",
	}
	repo := new(mocks.UserBalanceRepository)
	ubs := []userbalance.UserBalance{
		{
			Coin:         usdtCoin,
			CoinID:       1,
			Amount:       "1000.00",
			FrozenAmount: "100.00",
		},
		{
			Coin:         btcCoin,
			CoinID:       2,
			Amount:       "1.00",
			FrozenAmount: "0.1",
		},
	}
	repo.On("GetUserBalancesForCoins", 1, []int64{1, 2}).Once().Return(ubs)
	pg := new(mocks.PriceGenerator)
	permissionManager := new(mocks.UserPermissionManager)
	walletService := new(mocks.WalletService)
	userService := new(mocks.UserService)
	userWalletBalanceRepository := new(mocks.UserWalletBalanceRepository)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	ubService := userbalance.NewBalanceService(db, repo, cs, pg, permissionManager, walletService, userService,
		userWalletBalanceRepository, configs, logger)

	u := &user.User{
		ID: 1,
	}

	params := userbalance.GetPairBalancesParams{
		PairID:   1,
		PairName: "",
	}
	res, statusCode := ubService.GetPairBalances(u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	resp, ok := res.Data.(userbalance.GetPairBalanceResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}

	assert.Equal(t, "BTC-USDT", resp.PairData.Name)
	assert.Equal(t, int64(1), resp.PairData.ID)
	assert.Equal(t, "10", resp.PairData.MinimumOrderAmount)

	assert.Equal(t, 0.9, resp.Fee["makerFee"])
	assert.Equal(t, 0.9, resp.Fee["takerFee"])

	usdtBalance := resp.PairBalances[0]
	assert.Equal(t, int64(1), usdtBalance.CoinID)
	assert.Equal(t, "Tether", usdtBalance.CoinName)
	assert.Equal(t, "USDT", usdtBalance.CoinCode)
	assert.Equal(t, "900.00000000", usdtBalance.Balance)

	btcBalance := resp.PairBalances[1]
	assert.Equal(t, int64(2), btcBalance.CoinID)
	assert.Equal(t, "Bitcoin", btcBalance.CoinName)
	assert.Equal(t, "BTC", btcBalance.CoinCode)
	assert.Equal(t, "0.90000000", btcBalance.Balance)

	repo.AssertExpectations(t)
	cs.AssertExpectations(t)
}

func TestService_GetPairBalances_ByPairName(t *testing.T) {
	qm := userBalanceQueryMatcher{}
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
	db, err := gorm.Open(dialector, &gorm.Config{SkipDefaultTransaction: true})
	dbMock.MatchExpectationsInOrder(false)

	cs := new(mocks.CurrencyService)
	pair := currency.Pair{
		ID:                 1,
		Name:               "BTC-USDT",
		BasisCoinID:        1,
		DependentCoinID:    2,
		MinimumOrderAmount: sql.NullString{String: "10", Valid: true},
		MakerFee:           0.9,
		TakerFee:           0.9,
	}

	cs.On("GetPairByName", "BTC-USDT").Once().Return(pair, nil)

	usdtCoin := currency.Coin{
		ID:   1,
		Code: "USDT",
		Name: "Tether",
	}

	btcCoin := currency.Coin{
		ID:   2,
		Code: "BTC",
		Name: "Bitcoin",
	}
	repo := new(mocks.UserBalanceRepository)
	ubs := []userbalance.UserBalance{
		{
			Coin:         usdtCoin,
			CoinID:       1,
			Amount:       "1000.00",
			FrozenAmount: "100.00",
		},
		{
			Coin:         btcCoin,
			CoinID:       2,
			Amount:       "1.00",
			FrozenAmount: "0.1",
		},
	}
	repo.On("GetUserBalancesForCoins", 1, []int64{1, 2}).Once().Return(ubs)
	pg := new(mocks.PriceGenerator)
	permissionManager := new(mocks.UserPermissionManager)
	walletService := new(mocks.WalletService)
	userService := new(mocks.UserService)
	userWalletBalanceRepository := new(mocks.UserWalletBalanceRepository)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	ubService := userbalance.NewBalanceService(db, repo, cs, pg, permissionManager, walletService, userService,
		userWalletBalanceRepository, configs, logger)

	u := &user.User{
		ID: 1,
	}

	params := userbalance.GetPairBalancesParams{
		PairID:   1,
		PairName: "BTC-USDT",
	}
	res, statusCode := ubService.GetPairBalances(u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	resp, ok := res.Data.(userbalance.GetPairBalanceResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}

	assert.Equal(t, "BTC-USDT", resp.PairData.Name)
	assert.Equal(t, int64(1), resp.PairData.ID)
	assert.Equal(t, "10", resp.PairData.MinimumOrderAmount)

	assert.Equal(t, 0.9, resp.Fee["makerFee"])
	assert.Equal(t, 0.9, resp.Fee["takerFee"])

	usdtBalance := resp.PairBalances[0]
	assert.Equal(t, int64(1), usdtBalance.CoinID)
	assert.Equal(t, "Tether", usdtBalance.CoinName)
	assert.Equal(t, "USDT", usdtBalance.CoinCode)
	assert.Equal(t, "900.00000000", usdtBalance.Balance)

	btcBalance := resp.PairBalances[1]
	assert.Equal(t, int64(2), btcBalance.CoinID)
	assert.Equal(t, "Bitcoin", btcBalance.CoinName)
	assert.Equal(t, "BTC", btcBalance.CoinCode)
	assert.Equal(t, "0.90000000", btcBalance.Balance)

	repo.AssertExpectations(t)
	cs.AssertExpectations(t)
}

func TestService_GetAllBalances(t *testing.T) {

	qm := userBalanceQueryMatcher{}
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
	db, err := gorm.Open(dialector, &gorm.Config{SkipDefaultTransaction: true})
	dbMock.MatchExpectationsInOrder(false)

	cs := new(mocks.CurrencyService)
	usdtCoin := currency.Coin{
		ID:              1,
		Code:            "USDT",
		Name:            "Tether",
		SubUnit:         8,
		MinimumWithdraw: "100",
	}

	btcCoin := currency.Coin{
		ID:              2,
		Code:            "BTC",
		Name:            "Bitcoin",
		SubUnit:         8,
		MinimumWithdraw: "0.01",
	}

	ethCoin := currency.Coin{
		ID:              3,
		Code:            "ETH",
		Name:            "Ethereum",
		SubUnit:         8,
		MinimumWithdraw: "0.01",
	}

	repo := new(mocks.UserBalanceRepository)
	ubs := []userbalance.UserBalance{
		{
			ID:           1,
			UserID:       1,
			CoinID:       1,
			Coin:         usdtCoin,
			Address:      sql.NullString{String: "USDTAddress", Valid: true},
			Amount:       "1000.00",
			FrozenAmount: "100.00",
		},
		{
			ID:               2,
			UserID:           1,
			CoinID:           2,
			Coin:             btcCoin,
			FrozenBalance:    "",
			Address:          sql.NullString{String: "BTCAddress", Valid: true},
			Amount:           "1.00",
			FrozenAmount:     "0.1",
			AutoExchangeCoin: sql.NullString{String: "ETH", Valid: true},
		},
		{
			ID:            3,
			UserID:        1,
			CoinID:        3,
			Coin:          ethCoin,
			FrozenBalance: "",
			Address:       sql.NullString{String: "ETHAddress", Valid: true},
			Amount:        "1.5",
			FrozenAmount:  "0.1",
		},
	}
	repo.On("GetUserAllBalances", 1, mock.Anything).Once().Return(ubs)
	pg := new(mocks.PriceGenerator)
	//for usdt
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "USDT", "1000.00").Once().Return("1000.00", nil)
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "USDT", "900.00000000").Once().Return("900.000000", nil)
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "USDT", "100.00").Once().Return("100.00", nil)
	pg.On("GetAmountBasedOnBTC", mock.Anything, "USDT", "1000.00").Once().Return("0.0200000", nil)
	pg.On("GetAmountBasedOnBTC", mock.Anything, "USDT", "900.00000000").Once().Return("0.0018000", nil)
	pg.On("GetAmountBasedOnBTC", mock.Anything, "USDT", "100.00").Once().Return("0.0020000", nil)
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "USDT", "1.0").Once().Return("1.0", nil)

	//for btc
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "BTC", "1.00").Once().Return("50000", nil)
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "BTC", "0.90000000").Once().Return("45000", nil)
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "BTC", "0.1").Once().Return("5000", nil)
	pg.On("GetAmountBasedOnBTC", mock.Anything, "BTC", "1.00").Once().Return("1", nil)
	pg.On("GetAmountBasedOnBTC", mock.Anything, "BTC", "0.90000000").Once().Return("0.9", nil)
	pg.On("GetAmountBasedOnBTC", mock.Anything, "BTC", "0.1").Once().Return("0.1", nil)
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "BTC", "1.0").Once().Return("50000", nil)

	//for eth
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "ETH", "1.5").Once().Return("3000", nil)
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "ETH", "1.40000000").Once().Return("2800", nil)
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "ETH", "0.1").Once().Return("200", nil)
	pg.On("GetAmountBasedOnBTC", mock.Anything, "ETH", "1.5").Once().Return("0.06", nil)
	pg.On("GetAmountBasedOnBTC", mock.Anything, "ETH", "1.40000000").Once().Return("0.056", nil)
	pg.On("GetAmountBasedOnBTC", mock.Anything, "ETH", "0.1").Once().Return("0.004", nil)
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "ETH", "1.0").Once().Return("2000", nil)

	permissionManager := new(mocks.UserPermissionManager)
	walletService := new(mocks.WalletService)
	userService := new(mocks.UserService)
	userWalletBalanceRepository := new(mocks.UserWalletBalanceRepository)
	configs := new(mocks.Configs)
	configs.On("GetImagePath").Times(3).Return("http://127.0.0.1/")
	logger := new(mocks.Logger)

	ubService := userbalance.NewBalanceService(db, repo, cs, pg, permissionManager, walletService, userService,
		userWalletBalanceRepository, configs, logger)

	u := &user.User{
		ID:     1,
		Status: user.StatusVerified,
	}

	params := userbalance.GetAllBalancesParams{
		CoinCode: "",
		Name:     "",
		Sort:     "desc",
	}

	res, statusCode := ubService.GetAllBalances(u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	resp, ok := res.Data.(userbalance.GetAllBalancesResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}

	assert.Equal(t, "54000.00000000", resp.TotalSum)
	assert.Equal(t, "48700.00000000", resp.AvailableSum)
	assert.Equal(t, "5300.00000000", resp.InOrderSum)
	assert.Equal(t, "1.08000000", resp.BtcTotalSum)
	assert.Equal(t, "0.95780000", resp.BtcAvailableSum)
	assert.Equal(t, "0.10600000", resp.BtcInOrderSum)
	assert.Equal(t, "0", resp.MinimumOfSmallBalances)

	allBalances := resp.Balances
	//how we are sure this is btc because we sort desc and the btc balance is higher than others
	btcBalance := allBalances[0]
	assert.Equal(t, "1.00000000", btcBalance.TotalAmount)
	assert.Equal(t, "0.90000000", btcBalance.AvailableAmount)
	assert.Equal(t, "0.10000000", btcBalance.InOrderAmount)
	assert.Equal(t, "50000", btcBalance.EquivalentTotalAmount)
	assert.Equal(t, "45000", btcBalance.EquivalentAvailableAmount)
	assert.Equal(t, "5000", btcBalance.EquivalentInOrderAmount)
	assert.Equal(t, "BTC", btcBalance.CoinCode)
	assert.Equal(t, "Bitcoin", btcBalance.CoinName)
	assert.Equal(t, "50000", btcBalance.Price)
	assert.Equal(t, "0", btcBalance.Fee)
	assert.Equal(t, 8, btcBalance.SubUnit)
	assert.Equal(t, "0.01", btcBalance.MinimumWithdraw)
	assert.Equal(t, "1", btcBalance.BtcTotalEquivalentAmount)
	assert.Equal(t, "0.9", btcBalance.BtcAvailableEquivalentAmount)
	assert.Equal(t, "0.1", btcBalance.BtcInOrderEquivalentAmount)
	assert.Equal(t, "BTCAddress", btcBalance.Address)
	assert.Equal(t, "ETH", btcBalance.AutoExchangeCode)

	ethBalance := allBalances[1]
	assert.Equal(t, "1.50000000", ethBalance.TotalAmount)
	assert.Equal(t, "1.40000000", ethBalance.AvailableAmount)
	assert.Equal(t, "0.10000000", ethBalance.InOrderAmount)
	assert.Equal(t, "3000", ethBalance.EquivalentTotalAmount)
	assert.Equal(t, "2800", ethBalance.EquivalentAvailableAmount)
	assert.Equal(t, "200", ethBalance.EquivalentInOrderAmount)
	assert.Equal(t, "ETH", ethBalance.CoinCode)
	assert.Equal(t, "Ethereum", ethBalance.CoinName)
	assert.Equal(t, "2000", ethBalance.Price)
	assert.Equal(t, "0", ethBalance.Fee)
	assert.Equal(t, 8, ethBalance.SubUnit)
	assert.Equal(t, "0.01", ethBalance.MinimumWithdraw)
	assert.Equal(t, "0.06", ethBalance.BtcTotalEquivalentAmount)
	assert.Equal(t, "0.056", ethBalance.BtcAvailableEquivalentAmount)
	assert.Equal(t, "0.004", ethBalance.BtcInOrderEquivalentAmount)
	assert.Equal(t, "ETHAddress", ethBalance.Address)

	usdtBalance := allBalances[2]
	assert.Equal(t, "1000.00000000", usdtBalance.TotalAmount)
	assert.Equal(t, "900.00000000", usdtBalance.AvailableAmount)
	assert.Equal(t, "100.00000000", usdtBalance.InOrderAmount)
	assert.Equal(t, "1000.00", usdtBalance.EquivalentTotalAmount)
	assert.Equal(t, "900.000000", usdtBalance.EquivalentAvailableAmount)
	assert.Equal(t, "100.00", usdtBalance.EquivalentInOrderAmount)
	assert.Equal(t, "USDT", usdtBalance.CoinCode)
	assert.Equal(t, "Tether", usdtBalance.CoinName)
	assert.Equal(t, "1.0", usdtBalance.Price)
	assert.Equal(t, "0", usdtBalance.Fee)
	assert.Equal(t, 8, usdtBalance.SubUnit)
	assert.Equal(t, "100", usdtBalance.MinimumWithdraw)
	assert.Equal(t, "0.0200000", usdtBalance.BtcTotalEquivalentAmount)
	assert.Equal(t, "0.0018000", usdtBalance.BtcAvailableEquivalentAmount)
	assert.Equal(t, "0.0020000", usdtBalance.BtcInOrderEquivalentAmount)
	assert.Equal(t, "USDTAddress", usdtBalance.Address)

	repo.AssertExpectations(t)
	pg.AssertExpectations(t)
	configs.AssertExpectations(t)
}

func TestService_GetWithdrawDepositData_BTC_WithoutUserBalance(t *testing.T) {
	qm := userBalanceQueryMatcher{}
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
	db, err := gorm.Open(dialector, &gorm.Config{SkipDefaultTransaction: true})
	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectExec("INSERT INTO user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))

	btcCoin := currency.Coin{
		ID:               2,
		Code:             "BTC",
		Name:             "Bitcoin",
		SubUnit:          8,
		MinimumWithdraw:  "0.01",
		WithdrawalFee:    sql.NullFloat64{Float64: 0.01, Valid: true},
		SupportsWithdraw: sql.NullBool{Bool: true, Valid: true},
		SupportsDeposit:  sql.NullBool{Bool: true, Valid: true},
		DepositComments: sql.NullString{
			String: "test1|test2",
			Valid:  true,
		},
		WithdrawComments: sql.NullString{
			String: "",
			Valid:  false,
		},
	}

	cs := new(mocks.CurrencyService)
	cs.On("GetCoinByCode", "BTC").Once().Return(btcCoin, nil)

	repo := new(mocks.UserBalanceRepository)
	repo.On("GetBalanceOfUserByCoinID", 1, int64(2), mock.Anything).Once().Return(gorm.ErrRecordNotFound)

	pg := new(mocks.PriceGenerator)
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "BTC", "0.0").Once().Return("0", nil)
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "BTC", "0.0").Once().Return("0", nil)
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "BTC", "0.00000000").Once().Return("0", nil)
	pg.On("GetAmountBasedOnBTC", mock.Anything, "BTC", "0.0").Once().Return("0", nil)
	pg.On("GetAmountBasedOnBTC", mock.Anything, "BTC", "0.0").Once().Return("0", nil)
	pg.On("GetAmountBasedOnBTC", mock.Anything, "BTC", "0.00000000").Once().Return("0", nil)
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "BTC", "1.0").Once().Return("50000", nil)

	permissionManager := new(mocks.UserPermissionManager)
	permissionManager.On("IsPermissionGrantedToUserFor", mock.Anything, user.PermissionDeposit).Once().Return(true)
	permissionManager.On("IsPermissionGrantedToUserFor", mock.Anything, user.PermissionWithdraw).Once().Return(true)

	walletService := new(mocks.WalletService)
	walletService.On("GetAddressForUser", "BTC", mock.Anything).Once().Return("BTCAddress", nil)
	userService := new(mocks.UserService)
	userWalletBalanceRepository := new(mocks.UserWalletBalanceRepository)
	configs := new(mocks.Configs)
	configs.On("GetImagePath").Once().Return("http://127.0.0.1/")
	logger := new(mocks.Logger)

	ubService := userbalance.NewBalanceService(db, repo, cs, pg, permissionManager, walletService, userService,
		userWalletBalanceRepository, configs, logger)

	u := &user.User{
		ID:     1,
		Status: user.StatusVerified,
	}
	params := userbalance.GetWithdrawDepositParams{
		Coin: "BTC",
	}

	res, statusCode := ubService.GetWithdrawDepositData(u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	resp, ok := res.Data.(userbalance.GetWithdrawDepositDataResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}

	assert.Equal(t, "BTCAddress", resp.WalletAddress)
	assert.Equal(t, true, resp.SupportsWithdraw)
	assert.Equal(t, true, resp.SupportsDeposit)
	assert.Equal(t, true, resp.HasDepositPermission)
	assert.Equal(t, true, resp.HasWithdrawPermission)
	assert.Equal(t, "Bitcoin(BTC)", resp.CompletedNetworkName)

	networkConfigs := resp.NetworksConfigs
	assert.Equal(t, 1, len(networkConfigs))

	assert.Equal(t, true, networkConfigs[0].SupportsDeposit)
	assert.Equal(t, true, networkConfigs[0].SupportsWithdraw)
	assert.Equal(t, "BTCAddress", networkConfigs[0].Address)

	btcBalance := resp.Balance
	assert.Equal(t, "0.00000000", btcBalance.TotalAmount)
	assert.Equal(t, "0.00000000", btcBalance.AvailableAmount)
	assert.Equal(t, "0.00000000", btcBalance.InOrderAmount)
	assert.Equal(t, "0", btcBalance.EquivalentTotalAmount)
	assert.Equal(t, "0", btcBalance.EquivalentAvailableAmount)
	assert.Equal(t, "0", btcBalance.EquivalentInOrderAmount)
	assert.Equal(t, "BTC", btcBalance.CoinCode)
	assert.Equal(t, "Bitcoin", btcBalance.CoinName)
	assert.Equal(t, "50000", btcBalance.Price)
	assert.Equal(t, "0.01000000", btcBalance.Fee)
	assert.Equal(t, 8, btcBalance.SubUnit)
	assert.Equal(t, "0.01", btcBalance.MinimumWithdraw)
	assert.Equal(t, "0", btcBalance.BtcTotalEquivalentAmount)
	assert.Equal(t, "0", btcBalance.BtcAvailableEquivalentAmount)
	assert.Equal(t, "0", btcBalance.BtcInOrderEquivalentAmount)
	assert.Equal(t, "BTCAddress", btcBalance.Address)

	assert.Equal(t, []string{"test1", "test2"}, resp.DepositComments)
	assert.Equal(t, []string{}, resp.WithdrawComments)

	repo.AssertExpectations(t)
	cs.AssertExpectations(t)
	pg.AssertExpectations(t)
	permissionManager.AssertExpectations(t)
	walletService.AssertExpectations(t)
	configs.AssertExpectations(t)
}

func TestService_GetWithdrawDepositData_BTC_WithoutAddress(t *testing.T) {
	qm := userBalanceQueryMatcher{}
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
	db, err := gorm.Open(dialector, &gorm.Config{SkipDefaultTransaction: true})
	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))

	btcCoin := currency.Coin{
		ID:               2,
		Code:             "BTC",
		Name:             "Bitcoin",
		SubUnit:          8,
		MinimumWithdraw:  "0.01",
		WithdrawalFee:    sql.NullFloat64{Float64: 0.01, Valid: true},
		SupportsWithdraw: sql.NullBool{Bool: true, Valid: true},
		SupportsDeposit:  sql.NullBool{Bool: true, Valid: true},
		DepositComments: sql.NullString{
			String: "test1|test2",
			Valid:  true,
		},
		WithdrawComments: sql.NullString{
			String: "",
			Valid:  false,
		},
	}

	cs := new(mocks.CurrencyService)
	cs.On("GetCoinByCode", "BTC").Once().Return(btcCoin, nil)

	repo := new(mocks.UserBalanceRepository)
	repo.On("GetBalanceOfUserByCoinID", 1, int64(2), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub := args.Get(2).(*userbalance.UserBalance)
		ub.ID = 1
		ub.UserID = 1
		ub.Amount = "1.00"
		ub.FrozenAmount = "0.1"
		ub.Coin = btcCoin
		ub.CoinID = 2
	})

	pg := new(mocks.PriceGenerator)
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "BTC", "1.00").Once().Return("50000", nil)
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "BTC", "0.90000000").Once().Return("45000", nil)
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "BTC", "0.1").Once().Return("5000", nil)
	pg.On("GetAmountBasedOnBTC", mock.Anything, "BTC", "1.00").Once().Return("1", nil)
	pg.On("GetAmountBasedOnBTC", mock.Anything, "BTC", "0.90000000").Once().Return("0.9", nil)
	pg.On("GetAmountBasedOnBTC", mock.Anything, "BTC", "0.1").Once().Return("0.1", nil)
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "BTC", "1.0").Once().Return("50000", nil)

	permissionManager := new(mocks.UserPermissionManager)
	permissionManager.On("IsPermissionGrantedToUserFor", mock.Anything, user.PermissionDeposit).Once().Return(true)
	permissionManager.On("IsPermissionGrantedToUserFor", mock.Anything, user.PermissionWithdraw).Once().Return(true)

	walletService := new(mocks.WalletService)
	walletService.On("GetAddressForUser", "BTC", mock.Anything).Once().Return("BTCAddress", nil)
	userService := new(mocks.UserService)
	userWalletBalanceRepository := new(mocks.UserWalletBalanceRepository)
	configs := new(mocks.Configs)
	configs.On("GetImagePath").Once().Return("http://127.0.0.1/")
	logger := new(mocks.Logger)

	ubService := userbalance.NewBalanceService(db, repo, cs, pg, permissionManager, walletService, userService,
		userWalletBalanceRepository, configs, logger)

	u := &user.User{
		ID:     1,
		Status: user.StatusVerified,
	}

	params := userbalance.GetWithdrawDepositParams{
		Coin: "BTC",
	}

	res, statusCode := ubService.GetWithdrawDepositData(u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	resp, ok := res.Data.(userbalance.GetWithdrawDepositDataResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}

	assert.Equal(t, "BTCAddress", resp.WalletAddress)
	assert.Equal(t, true, resp.SupportsWithdraw)
	assert.Equal(t, true, resp.SupportsDeposit)
	assert.Equal(t, true, resp.HasDepositPermission)
	assert.Equal(t, true, resp.HasWithdrawPermission)
	assert.Equal(t, "Bitcoin(BTC)", resp.CompletedNetworkName)

	networkConfigs := resp.NetworksConfigs
	assert.Equal(t, 1, len(networkConfigs))

	assert.Equal(t, true, networkConfigs[0].SupportsDeposit)
	assert.Equal(t, true, networkConfigs[0].SupportsWithdraw)
	assert.Equal(t, "BTCAddress", networkConfigs[0].Address)

	btcBalance := resp.Balance
	assert.Equal(t, "1.00000000", btcBalance.TotalAmount)
	assert.Equal(t, "0.90000000", btcBalance.AvailableAmount)
	assert.Equal(t, "0.10000000", btcBalance.InOrderAmount)
	assert.Equal(t, "50000", btcBalance.EquivalentTotalAmount)
	assert.Equal(t, "45000", btcBalance.EquivalentAvailableAmount)
	assert.Equal(t, "5000", btcBalance.EquivalentInOrderAmount)
	assert.Equal(t, "BTC", btcBalance.CoinCode)
	assert.Equal(t, "Bitcoin", btcBalance.CoinName)
	assert.Equal(t, "50000", btcBalance.Price)
	assert.Equal(t, "0.01000000", btcBalance.Fee)
	assert.Equal(t, 8, btcBalance.SubUnit)
	assert.Equal(t, "0.01", btcBalance.MinimumWithdraw)
	assert.Equal(t, "1", btcBalance.BtcTotalEquivalentAmount)
	assert.Equal(t, "0.9", btcBalance.BtcAvailableEquivalentAmount)
	assert.Equal(t, "0.1", btcBalance.BtcInOrderEquivalentAmount)
	assert.Equal(t, "BTCAddress", btcBalance.Address)

	assert.Equal(t, []string{"test1", "test2"}, resp.DepositComments)
	assert.Equal(t, []string{}, resp.WithdrawComments)

	repo.AssertExpectations(t)
	cs.AssertExpectations(t)
	pg.AssertExpectations(t)
	permissionManager.AssertExpectations(t)
	walletService.AssertExpectations(t)
	configs.AssertExpectations(t)
}

func TestService_GetWithdrawDepositData_BTC_WithAddress(t *testing.T) {
	qm := userBalanceQueryMatcher{}
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
	db, err := gorm.Open(dialector, &gorm.Config{SkipDefaultTransaction: true})
	dbMock.MatchExpectationsInOrder(false)

	btcCoin := currency.Coin{
		ID:               2,
		Code:             "BTC",
		Name:             "Bitcoin",
		SubUnit:          8,
		MinimumWithdraw:  "0.01",
		WithdrawalFee:    sql.NullFloat64{Float64: 0.01, Valid: true},
		SupportsWithdraw: sql.NullBool{Bool: true, Valid: true},
		SupportsDeposit:  sql.NullBool{Bool: true, Valid: true},
		DepositComments: sql.NullString{
			String: "test1|test2",
			Valid:  true,
		},
		WithdrawComments: sql.NullString{
			String: "",
			Valid:  false,
		},
	}

	cs := new(mocks.CurrencyService)
	cs.On("GetCoinByCode", "BTC").Once().Return(btcCoin, nil)

	repo := new(mocks.UserBalanceRepository)
	repo.On("GetBalanceOfUserByCoinID", 1, int64(2), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub := args.Get(2).(*userbalance.UserBalance)
		ub.ID = 1
		ub.UserID = 1
		ub.Amount = "1.00"
		ub.FrozenAmount = "0.1"
		ub.Coin = btcCoin
		ub.CoinID = 2
		ub.Address = sql.NullString{String: "BTCAddress", Valid: true}
	})

	pg := new(mocks.PriceGenerator)
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "BTC", "1.00").Once().Return("50000", nil)
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "BTC", "0.90000000").Once().Return("45000", nil)
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "BTC", "0.1").Once().Return("5000", nil)
	pg.On("GetAmountBasedOnBTC", mock.Anything, "BTC", "1.00").Once().Return("1", nil)
	pg.On("GetAmountBasedOnBTC", mock.Anything, "BTC", "0.90000000").Once().Return("0.9", nil)
	pg.On("GetAmountBasedOnBTC", mock.Anything, "BTC", "0.1").Once().Return("0.1", nil)
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "BTC", "1.0").Once().Return("50000", nil)

	permissionManager := new(mocks.UserPermissionManager)
	permissionManager.On("IsPermissionGrantedToUserFor", mock.Anything, user.PermissionDeposit).Once().Return(true)
	permissionManager.On("IsPermissionGrantedToUserFor", mock.Anything, user.PermissionWithdraw).Once().Return(true)

	walletService := new(mocks.WalletService)
	userService := new(mocks.UserService)
	userWalletBalanceRepository := new(mocks.UserWalletBalanceRepository)
	configs := new(mocks.Configs)
	configs.On("GetImagePath").Once().Return("http://127.0.0.1/")
	logger := new(mocks.Logger)

	ubService := userbalance.NewBalanceService(db, repo, cs, pg, permissionManager, walletService, userService,
		userWalletBalanceRepository, configs, logger)

	u := &user.User{
		ID:     1,
		Status: user.StatusVerified,
	}

	params := userbalance.GetWithdrawDepositParams{
		Coin: "BTC",
	}

	res, statusCode := ubService.GetWithdrawDepositData(u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	resp, ok := res.Data.(userbalance.GetWithdrawDepositDataResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}

	assert.Equal(t, "BTCAddress", resp.WalletAddress)
	assert.Equal(t, true, resp.SupportsWithdraw)
	assert.Equal(t, true, resp.SupportsDeposit)
	assert.Equal(t, true, resp.HasDepositPermission)
	assert.Equal(t, true, resp.HasWithdrawPermission)
	assert.Equal(t, "Bitcoin(BTC)", resp.CompletedNetworkName)

	networkConfigs := resp.NetworksConfigs
	assert.Equal(t, 1, len(networkConfigs))

	assert.Equal(t, true, networkConfigs[0].SupportsDeposit)
	assert.Equal(t, true, networkConfigs[0].SupportsWithdraw)
	assert.Equal(t, "BTCAddress", networkConfigs[0].Address)

	btcBalance := resp.Balance
	assert.Equal(t, "1.00000000", btcBalance.TotalAmount)
	assert.Equal(t, "0.90000000", btcBalance.AvailableAmount)
	assert.Equal(t, "0.10000000", btcBalance.InOrderAmount)
	assert.Equal(t, "50000", btcBalance.EquivalentTotalAmount)
	assert.Equal(t, "45000", btcBalance.EquivalentAvailableAmount)
	assert.Equal(t, "5000", btcBalance.EquivalentInOrderAmount)
	assert.Equal(t, "BTC", btcBalance.CoinCode)
	assert.Equal(t, "Bitcoin", btcBalance.CoinName)
	assert.Equal(t, "50000", btcBalance.Price)
	assert.Equal(t, "0.01000000", btcBalance.Fee)
	assert.Equal(t, 8, btcBalance.SubUnit)
	assert.Equal(t, "0.01", btcBalance.MinimumWithdraw)
	assert.Equal(t, "1", btcBalance.BtcTotalEquivalentAmount)
	assert.Equal(t, "0.9", btcBalance.BtcAvailableEquivalentAmount)
	assert.Equal(t, "0.1", btcBalance.BtcInOrderEquivalentAmount)
	assert.Equal(t, "BTCAddress", btcBalance.Address)

	assert.Equal(t, []string{"test1", "test2"}, resp.DepositComments)
	assert.Equal(t, []string{}, resp.WithdrawComments)

	repo.AssertExpectations(t)
	cs.AssertExpectations(t)
	pg.AssertExpectations(t)
	permissionManager.AssertExpectations(t)
	configs.AssertExpectations(t)
}

func TestService_GetWithdrawDepositData_USDT_WithOutTRC20Address_WithTrxAddress(t *testing.T) {
	qm := userBalanceQueryMatcher{}
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
	db, err := gorm.Open(dialector, &gorm.Config{SkipDefaultTransaction: true})
	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))

	otherNetworksConfigs := `[{"code":"TRX","supportsWithdraw":true,"supportsDeposit":true,"completedNetworkName":"Tron (TRX)","fee":"2.5"}]`

	usdtCoin := currency.Coin{
		ID:                             1,
		Code:                           "USDT",
		Name:                           "Tether",
		SubUnit:                        8,
		MinimumWithdraw:                "100.0",
		CompletedNetworkName:           sql.NullString{String: "Ethereum (ETH)", Valid: true},
		WithdrawalFee:                  sql.NullFloat64{Float64: 10, Valid: true},
		SupportsWithdraw:               sql.NullBool{Bool: true, Valid: true},
		SupportsDeposit:                sql.NullBool{Bool: true, Valid: true},
		OtherBlockchainNetworksConfigs: sql.NullString{String: otherNetworksConfigs, Valid: true},
		DepositComments: sql.NullString{
			String: "test1|test2",
			Valid:  true,
		},
		WithdrawComments: sql.NullString{
			String: "",
			Valid:  false,
		},
	}

	trxCoin := currency.Coin{
		ID:                   3,
		Code:                 "TRX",
		Name:                 "Tron",
		SubUnit:              8,
		MinimumWithdraw:      "100.0",
		CompletedNetworkName: sql.NullString{String: "Tron (Trx)", Valid: true},
		WithdrawalFee:        sql.NullFloat64{Float64: 10, Valid: true},
		SupportsWithdraw:     sql.NullBool{Bool: true, Valid: true},
		SupportsDeposit:      sql.NullBool{Bool: true, Valid: true},
	}

	cs := new(mocks.CurrencyService)
	cs.On("GetCoinByCode", "USDT").Once().Return(usdtCoin, nil)
	cs.On("GetCoinByCode", "TRX").Once().Return(trxCoin, nil)

	repo := new(mocks.UserBalanceRepository)
	repo.On("GetBalanceOfUserByCoinID", 1, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub := args.Get(2).(*userbalance.UserBalance)
		ub.ID = 1
		ub.UserID = 1
		ub.Amount = "1000"
		ub.FrozenAmount = "100"
		ub.Coin = usdtCoin
		ub.CoinID = 1
		ub.Address = sql.NullString{String: "USDTAddress", Valid: true}
	})

	repo.On("GetBalanceOfUserByCoinID", 1, int64(3), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub := args.Get(2).(*userbalance.UserBalance)
		ub.ID = 1
		ub.UserID = 1
		ub.Amount = "1000"
		ub.FrozenAmount = "100"
		ub.Coin = trxCoin
		ub.CoinID = 3
		ub.Address = sql.NullString{String: "TRXAddress", Valid: true}
	})

	pg := new(mocks.PriceGenerator)
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "USDT", "1000").Once().Return("1000", nil)
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "USDT", "900.00000000").Once().Return("900", nil)
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "USDT", "100").Once().Return("100", nil)
	pg.On("GetAmountBasedOnBTC", mock.Anything, "USDT", "1000").Once().Return("0.02", nil)
	pg.On("GetAmountBasedOnBTC", mock.Anything, "USDT", "900.00000000").Once().Return("0.0018", nil)
	pg.On("GetAmountBasedOnBTC", mock.Anything, "USDT", "100").Once().Return("0.002", nil)
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "USDT", "1.0").Once().Return("1", nil)

	permissionManager := new(mocks.UserPermissionManager)
	permissionManager.On("IsPermissionGrantedToUserFor", mock.Anything, user.PermissionDeposit).Once().Return(true)
	permissionManager.On("IsPermissionGrantedToUserFor", mock.Anything, user.PermissionWithdraw).Once().Return(true)

	walletService := new(mocks.WalletService)
	userService := new(mocks.UserService)
	userWalletBalanceRepository := new(mocks.UserWalletBalanceRepository)
	configs := new(mocks.Configs)
	configs.On("GetImagePath").Once().Return("http://127.0.0.1/")
	logger := new(mocks.Logger)

	ubService := userbalance.NewBalanceService(db, repo, cs, pg, permissionManager, walletService, userService,
		userWalletBalanceRepository, configs, logger)

	u := &user.User{
		ID:     1,
		Status: user.StatusVerified,
	}

	params := userbalance.GetWithdrawDepositParams{
		Coin: "USDT",
	}

	res, statusCode := ubService.GetWithdrawDepositData(u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	resp, ok := res.Data.(userbalance.GetWithdrawDepositDataResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}

	assert.Equal(t, "USDTAddress", resp.WalletAddress)
	assert.Equal(t, true, resp.SupportsWithdraw)
	assert.Equal(t, true, resp.SupportsDeposit)
	assert.Equal(t, true, resp.HasDepositPermission)
	assert.Equal(t, true, resp.HasWithdrawPermission)
	assert.Equal(t, "Ethereum (ETH)", resp.CompletedNetworkName)

	networkConfigs := resp.NetworksConfigs
	assert.Equal(t, 2, len(networkConfigs))
	assert.Equal(t, true, networkConfigs[0].SupportsDeposit)
	assert.Equal(t, true, networkConfigs[0].SupportsWithdraw)
	assert.Equal(t, "10.00000000", networkConfigs[0].Fee)
	assert.Equal(t, "USDTAddress", networkConfigs[0].Address)

	assert.Equal(t, true, networkConfigs[1].SupportsDeposit)
	assert.Equal(t, true, networkConfigs[1].SupportsWithdraw)
	assert.Equal(t, "2.5", networkConfigs[1].Fee)
	assert.Equal(t, "TRXAddress", networkConfigs[1].Address)

	otherConfigs := resp.OtherNetworksConfigs
	assert.Equal(t, 1, len(otherConfigs))
	assert.Equal(t, true, otherConfigs[0].SupportsDeposit)
	assert.Equal(t, true, otherConfigs[0].SupportsWithdraw)
	assert.Equal(t, "2.5", otherConfigs[0].Fee)
	assert.Equal(t, "TRXAddress", otherConfigs[0].Address)

	btcBalance := resp.Balance
	assert.Equal(t, "1000.00000000", btcBalance.TotalAmount)
	assert.Equal(t, "900.00000000", btcBalance.AvailableAmount)
	assert.Equal(t, "100.00000000", btcBalance.InOrderAmount)
	assert.Equal(t, "1000", btcBalance.EquivalentTotalAmount)
	assert.Equal(t, "900", btcBalance.EquivalentAvailableAmount)
	assert.Equal(t, "100", btcBalance.EquivalentInOrderAmount)
	assert.Equal(t, "USDT", btcBalance.CoinCode)
	assert.Equal(t, "Tether", btcBalance.CoinName)
	assert.Equal(t, "1", btcBalance.Price)
	assert.Equal(t, "10.00000000", btcBalance.Fee)
	assert.Equal(t, 8, btcBalance.SubUnit)
	assert.Equal(t, "100.0", btcBalance.MinimumWithdraw)
	assert.Equal(t, "0.02", btcBalance.BtcTotalEquivalentAmount)
	assert.Equal(t, "0.0018", btcBalance.BtcAvailableEquivalentAmount)
	assert.Equal(t, "0.002", btcBalance.BtcInOrderEquivalentAmount)
	assert.Equal(t, "USDTAddress", btcBalance.Address)

	assert.Equal(t, []string{"test1", "test2"}, resp.DepositComments)
	assert.Equal(t, []string{}, resp.WithdrawComments)

	repo.AssertExpectations(t)
	cs.AssertExpectations(t)
	pg.AssertExpectations(t)
	permissionManager.AssertExpectations(t)
	configs.AssertExpectations(t)
}

func TestService_GetWithdrawDepositData_USDT_WithOutERC20_And_TRC20Address(t *testing.T) {
	qm := userBalanceQueryMatcher{}
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
	db, err := gorm.Open(dialector, &gorm.Config{SkipDefaultTransaction: true})
	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))

	otherNetworksConfigs := `[{"code":"TRX","supportsWithdraw":true,"supportsDeposit":true,"completedNetworkName":"Tron (TRX)","fee":"2.5"}]`
	usdtCoin := currency.Coin{
		ID:                             1,
		Code:                           "USDT",
		Name:                           "Tether",
		SubUnit:                        8,
		MinimumWithdraw:                "100.0",
		BlockchainNetwork:              sql.NullString{String: "ETH", Valid: true},
		CompletedNetworkName:           sql.NullString{String: "Ethereum (ETH)", Valid: true},
		WithdrawalFee:                  sql.NullFloat64{Float64: 10, Valid: true},
		SupportsWithdraw:               sql.NullBool{Bool: true, Valid: true},
		SupportsDeposit:                sql.NullBool{Bool: true, Valid: true},
		OtherBlockchainNetworksConfigs: sql.NullString{String: otherNetworksConfigs, Valid: true},
		DepositComments: sql.NullString{
			String: "test1|test2",
			Valid:  true,
		},
		WithdrawComments: sql.NullString{
			String: "",
			Valid:  false,
		},
	}

	ethCoin := currency.Coin{
		ID:                   2,
		Code:                 "ETH",
		Name:                 "Ethereum",
		SubUnit:              8,
		MinimumWithdraw:      "0.1",
		CompletedNetworkName: sql.NullString{String: "Ethereum (ETH)", Valid: true},
		WithdrawalFee:        sql.NullFloat64{Float64: 0.01, Valid: true},
		SupportsWithdraw:     sql.NullBool{Bool: true, Valid: true},
		SupportsDeposit:      sql.NullBool{Bool: true, Valid: true},
	}

	trxCoin := currency.Coin{
		ID:                   3,
		Code:                 "TRX",
		Name:                 "Tron",
		SubUnit:              8,
		MinimumWithdraw:      "100.0",
		CompletedNetworkName: sql.NullString{String: "Tron (Trx)", Valid: true},
		WithdrawalFee:        sql.NullFloat64{Float64: 10, Valid: true},
		SupportsWithdraw:     sql.NullBool{Bool: true, Valid: true},
		SupportsDeposit:      sql.NullBool{Bool: true, Valid: true},
	}

	cs := new(mocks.CurrencyService)
	cs.On("GetCoinByCode", "USDT").Once().Return(usdtCoin, nil)
	cs.On("GetCoinByCode", "ETH").Once().Return(ethCoin, nil)
	cs.On("GetCoinByCode", "TRX").Once().Return(trxCoin, nil)

	repo := new(mocks.UserBalanceRepository)
	repo.On("GetBalanceOfUserByCoinID", 1, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub := args.Get(2).(*userbalance.UserBalance)
		ub.ID = 1
		ub.UserID = 1
		ub.Amount = "1000"
		ub.FrozenAmount = "100"
		ub.Coin = usdtCoin
		ub.CoinID = 1
	})

	repo.On("GetBalanceOfUserByCoinID", 1, int64(2), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub := args.Get(2).(*userbalance.UserBalance)
		ub.ID = 2
		ub.UserID = 1
		ub.Amount = "0"
		ub.FrozenAmount = "0"
		ub.Coin = ethCoin
		ub.CoinID = 2
	})

	repo.On("GetBalanceOfUserByCoinID", 1, int64(3), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub := args.Get(2).(*userbalance.UserBalance)
		ub.ID = 3
		ub.UserID = 1
		ub.Amount = "1000"
		ub.FrozenAmount = "100"
		ub.Coin = trxCoin
		ub.CoinID = 3
	})

	pg := new(mocks.PriceGenerator)
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "USDT", "1000").Once().Return("1000", nil)
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "USDT", "900.00000000").Once().Return("900", nil)
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "USDT", "100").Once().Return("100", nil)
	pg.On("GetAmountBasedOnBTC", mock.Anything, "USDT", "1000").Once().Return("0.02", nil)
	pg.On("GetAmountBasedOnBTC", mock.Anything, "USDT", "900.00000000").Once().Return("0.0018", nil)
	pg.On("GetAmountBasedOnBTC", mock.Anything, "USDT", "100").Once().Return("0.002", nil)
	pg.On("GetAmountBasedOnUSDT", mock.Anything, "USDT", "1.0").Once().Return("1", nil)

	permissionManager := new(mocks.UserPermissionManager)
	permissionManager.On("IsPermissionGrantedToUserFor", mock.Anything, user.PermissionDeposit).Once().Return(true)
	permissionManager.On("IsPermissionGrantedToUserFor", mock.Anything, user.PermissionWithdraw).Once().Return(true)

	walletService := new(mocks.WalletService)
	walletService.On("GetAddressForUser", "ETH", mock.Anything).Once().Return("ETHAddress", nil)
	walletService.On("GetAddressForUser", "TRX", mock.Anything).Once().Return("TRXAddress", nil)

	userService := new(mocks.UserService)
	userWalletBalanceRepository := new(mocks.UserWalletBalanceRepository)
	configs := new(mocks.Configs)
	configs.On("GetImagePath").Once().Return("http://127.0.0.1/")
	logger := new(mocks.Logger)

	ubService := userbalance.NewBalanceService(db, repo, cs, pg, permissionManager, walletService, userService,
		userWalletBalanceRepository, configs, logger)

	u := &user.User{
		ID:     1,
		Status: user.StatusVerified,
	}

	params := userbalance.GetWithdrawDepositParams{
		Coin: "USDT",
	}

	res, statusCode := ubService.GetWithdrawDepositData(u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	resp, ok := res.Data.(userbalance.GetWithdrawDepositDataResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}

	assert.Equal(t, "ETHAddress", resp.WalletAddress)
	assert.Equal(t, true, resp.SupportsWithdraw)
	assert.Equal(t, true, resp.SupportsDeposit)
	assert.Equal(t, true, resp.HasDepositPermission)
	assert.Equal(t, true, resp.HasWithdrawPermission)
	assert.Equal(t, "Ethereum (ETH)", resp.CompletedNetworkName)

	networkConfigs := resp.NetworksConfigs
	assert.Equal(t, 2, len(networkConfigs))
	assert.Equal(t, true, networkConfigs[0].SupportsDeposit)
	assert.Equal(t, true, networkConfigs[0].SupportsWithdraw)
	assert.Equal(t, "10.00000000", networkConfigs[0].Fee)
	assert.Equal(t, "ETHAddress", networkConfigs[0].Address)

	assert.Equal(t, true, networkConfigs[1].SupportsDeposit)
	assert.Equal(t, true, networkConfigs[1].SupportsWithdraw)
	assert.Equal(t, "2.5", networkConfigs[1].Fee)
	assert.Equal(t, "TRXAddress", networkConfigs[1].Address)

	otherConfigs := resp.OtherNetworksConfigs
	assert.Equal(t, 1, len(otherConfigs))
	assert.Equal(t, true, otherConfigs[0].SupportsDeposit)
	assert.Equal(t, true, otherConfigs[0].SupportsWithdraw)
	assert.Equal(t, "2.5", otherConfigs[0].Fee)
	assert.Equal(t, "TRXAddress", otherConfigs[0].Address)

	btcBalance := resp.Balance
	assert.Equal(t, "1000.00000000", btcBalance.TotalAmount)
	assert.Equal(t, "900.00000000", btcBalance.AvailableAmount)
	assert.Equal(t, "100.00000000", btcBalance.InOrderAmount)
	assert.Equal(t, "1000", btcBalance.EquivalentTotalAmount)
	assert.Equal(t, "900", btcBalance.EquivalentAvailableAmount)
	assert.Equal(t, "100", btcBalance.EquivalentInOrderAmount)
	assert.Equal(t, "USDT", btcBalance.CoinCode)
	assert.Equal(t, "Tether", btcBalance.CoinName)
	assert.Equal(t, "1", btcBalance.Price)
	assert.Equal(t, "10.00000000", btcBalance.Fee)
	assert.Equal(t, 8, btcBalance.SubUnit)
	assert.Equal(t, "100.0", btcBalance.MinimumWithdraw)
	assert.Equal(t, "0.02", btcBalance.BtcTotalEquivalentAmount)
	assert.Equal(t, "0.0018", btcBalance.BtcAvailableEquivalentAmount)
	assert.Equal(t, "0.002", btcBalance.BtcInOrderEquivalentAmount)
	assert.Equal(t, "ETHAddress", btcBalance.Address)

	assert.Equal(t, []string{"test1", "test2"}, resp.DepositComments)
	assert.Equal(t, []string{}, resp.WithdrawComments)

	repo.AssertExpectations(t)
	cs.AssertExpectations(t)
	pg.AssertExpectations(t)
	walletService.AssertExpectations(t)
	permissionManager.AssertExpectations(t)
	configs.AssertExpectations(t)
}

func TestService_GetBalanceOfUserByCoinID(t *testing.T) {
	qm := userBalanceQueryMatcher{}
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
	db, err := gorm.Open(dialector, &gorm.Config{SkipDefaultTransaction: true})
	dbMock.MatchExpectationsInOrder(false)

	cs := new(mocks.CurrencyService)

	ethCoin := currency.Coin{
		ID:               2,
		Code:             "ETH",
		Name:             "Ethereum",
		SubUnit:          8,
		MinimumWithdraw:  "0.1",
		WithdrawalFee:    sql.NullFloat64{Float64: 0.01, Valid: true},
		SupportsWithdraw: sql.NullBool{Bool: true, Valid: true},
		SupportsDeposit:  sql.NullBool{Bool: true, Valid: true},
	}
	repo := new(mocks.UserBalanceRepository)
	repo.On("GetBalanceOfUserByCoinID", 1, int64(2), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub := args.Get(2).(*userbalance.UserBalance)
		ub.ID = 1
		ub.UserID = 1
		ub.Amount = "1"
		ub.FrozenAmount = "0.1"
		ub.Coin = ethCoin
		ub.CoinID = 2
	})
	pg := new(mocks.PriceGenerator)
	permissionManager := new(mocks.UserPermissionManager)
	walletService := new(mocks.WalletService)
	userService := new(mocks.UserService)
	userWalletBalanceRepository := new(mocks.UserWalletBalanceRepository)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	ubService := userbalance.NewBalanceService(db, repo, cs, pg, permissionManager, walletService, userService,
		userWalletBalanceRepository, configs, logger)

	ub := userbalance.UserBalance{}
	err = ubService.GetBalanceOfUserByCoinID(1, int64(2), &ub)
	assert.Nil(t, err)
	repo.AssertExpectations(t)
	assert.Equal(t, int64(1), ub.ID)
	assert.Equal(t, "1", ub.Amount)
	assert.Equal(t, "0.1", ub.FrozenAmount)
	assert.Equal(t, int64(2), ub.CoinID)
}

func TestService_GetBalancesOfUsersForCoins(t *testing.T) {
	qm := userBalanceQueryMatcher{}
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
	db, err := gorm.Open(dialector, &gorm.Config{SkipDefaultTransaction: true})
	dbMock.MatchExpectationsInOrder(false)

	cs := new(mocks.CurrencyService)

	usdtCoin := currency.Coin{
		ID:               1,
		Code:             "USDT",
		Name:             "Tether",
		SubUnit:          8,
		MinimumWithdraw:  "100",
		WithdrawalFee:    sql.NullFloat64{Float64: 0.01, Valid: true},
		SupportsWithdraw: sql.NullBool{Bool: true, Valid: true},
		SupportsDeposit:  sql.NullBool{Bool: true, Valid: true},
	}

	btcCoin := currency.Coin{
		ID:               2,
		Code:             "BTC",
		Name:             "Bitcoin",
		SubUnit:          8,
		MinimumWithdraw:  "0.001",
		WithdrawalFee:    sql.NullFloat64{Float64: 0.01, Valid: true},
		SupportsWithdraw: sql.NullBool{Bool: true, Valid: true},
		SupportsDeposit:  sql.NullBool{Bool: true, Valid: true},
	}

	ethCoin := currency.Coin{
		ID:               3,
		Code:             "ETH",
		Name:             "Ethereum",
		SubUnit:          8,
		MinimumWithdraw:  "0.1",
		WithdrawalFee:    sql.NullFloat64{Float64: 0.01, Valid: true},
		SupportsWithdraw: sql.NullBool{Bool: true, Valid: true},
		SupportsDeposit:  sql.NullBool{Bool: true, Valid: true},
	}

	ubs := []userbalance.UserBalance{
		{
			ID:           1,
			UserID:       1,
			CoinID:       1,
			Coin:         usdtCoin,
			Address:      sql.NullString{String: "USDTAddress", Valid: true},
			Amount:       "1000.00",
			FrozenAmount: "100.00",
		},
		{
			ID:            2,
			UserID:        1,
			CoinID:        2,
			Coin:          btcCoin,
			FrozenBalance: "",
			Address:       sql.NullString{String: "BTCAddress", Valid: true},
			Amount:        "1.00",
			FrozenAmount:  "0.1",
		},
		{
			ID:            3,
			UserID:        1,
			CoinID:        3,
			Coin:          ethCoin,
			FrozenBalance: "",
			Address:       sql.NullString{String: "ETHAddress", Valid: true},
			Amount:        "1.5",
			FrozenAmount:  "0.1",
		},

		{
			ID:           1,
			UserID:       2,
			CoinID:       1,
			Coin:         usdtCoin,
			Address:      sql.NullString{String: "USDTAddress", Valid: true},
			Amount:       "1000.00",
			FrozenAmount: "100.00",
		},
		{
			ID:            2,
			UserID:        2,
			CoinID:        2,
			Coin:          btcCoin,
			FrozenBalance: "",
			Address:       sql.NullString{String: "BTCAddress", Valid: true},
			Amount:        "1.00",
			FrozenAmount:  "0.1",
		},
		{
			ID:            3,
			UserID:        2,
			CoinID:        3,
			Coin:          ethCoin,
			FrozenBalance: "",
			Address:       sql.NullString{String: "ETHAddress", Valid: true},
			Amount:        "1.5",
			FrozenAmount:  "0.1",
		},
	}

	repo := new(mocks.UserBalanceRepository)
	repo.On("GetBalancesOfUsersForCoins", []int{1, 2}, []int64{1, 2, 3}).Once().Return(ubs)
	pg := new(mocks.PriceGenerator)
	permissionManager := new(mocks.UserPermissionManager)
	walletService := new(mocks.WalletService)
	userService := new(mocks.UserService)
	userWalletBalanceRepository := new(mocks.UserWalletBalanceRepository)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	ubService := userbalance.NewBalanceService(db, repo, cs, pg, permissionManager, walletService, userService,
		userWalletBalanceRepository, configs, logger)

	userBalances := ubService.GetBalancesOfUsersForCoins([]int{1, 2}, []int64{1, 2, 3})
	assert.Equal(t, 6, len(userBalances))

	user1UsdtBalance := userBalances[0]
	assert.Equal(t, 1, user1UsdtBalance.UserID)
	assert.Equal(t, int64(1), user1UsdtBalance.CoinID)
	assert.Equal(t, "1000.00", user1UsdtBalance.Amount)
	assert.Equal(t, "100.00", user1UsdtBalance.FrozenAmount)

	user1BtcBalance := userBalances[1]
	assert.Equal(t, 1, user1BtcBalance.UserID)
	assert.Equal(t, int64(2), user1BtcBalance.CoinID)
	assert.Equal(t, "1.00", user1BtcBalance.Amount)
	assert.Equal(t, "0.1", user1BtcBalance.FrozenAmount)

	user1EthBalance := userBalances[2]
	assert.Equal(t, 1, user1EthBalance.UserID)
	assert.Equal(t, int64(3), user1EthBalance.CoinID)
	assert.Equal(t, "1.5", user1EthBalance.Amount)
	assert.Equal(t, "0.1", user1EthBalance.FrozenAmount)

	user2UsdtBalance := userBalances[3]
	assert.Equal(t, 2, user2UsdtBalance.UserID)
	assert.Equal(t, int64(1), user2UsdtBalance.CoinID)
	assert.Equal(t, "1000.00", user2UsdtBalance.Amount)
	assert.Equal(t, "100.00", user2UsdtBalance.FrozenAmount)

	user2BtcBalance := userBalances[4]
	assert.Equal(t, 2, user2BtcBalance.UserID)
	assert.Equal(t, int64(2), user2BtcBalance.CoinID)
	assert.Equal(t, "1.00", user2BtcBalance.Amount)
	assert.Equal(t, "0.1", user2BtcBalance.FrozenAmount)

	user2EthBalance := userBalances[5]
	assert.Equal(t, 2, user2EthBalance.UserID)
	assert.Equal(t, int64(3), user2EthBalance.CoinID)
	assert.Equal(t, "1.5", user2EthBalance.Amount)
	assert.Equal(t, "0.1", user2EthBalance.FrozenAmount)
}

func TestService_GenerateBalancesAndAddressForUser(t *testing.T) {
	qm := userBalanceQueryMatcher{}
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
	db, err := gorm.Open(dialector, &gorm.Config{SkipDefaultTransaction: true})
	dbMock.MatchExpectationsInOrder(false)

	dbMock.ExpectExec("INSERT INTO user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("Update user_balances").WillReturnResult(sqlmock.NewResult(1, 1))

	cs := new(mocks.CurrencyService)

	coins := []currency.Coin{
		{
			ID:   1,
			Name: "Bitcoin",
			Code: "BTC",
		},
	}

	cs.On("GetActiveCoins").Once().Return(coins)

	repo := new(mocks.UserBalanceRepository)
	pg := new(mocks.PriceGenerator)
	permissionManager := new(mocks.UserPermissionManager)
	walletService := new(mocks.WalletService)
	walletService.On("GetAddressForUser", "BTC", "").Once().Return("", nil)
	userService := new(mocks.UserService)
	userWalletBalanceRepository := new(mocks.UserWalletBalanceRepository)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	ubService := userbalance.NewBalanceService(db, repo, cs, pg, permissionManager, walletService, userService,
		userWalletBalanceRepository, configs, logger)

	u := user.User{
		ID: 1,
	}

	ubService.GenerateBalancesAndAddressForUser(u)
	repo.AssertExpectations(t)
	cs.AssertExpectations(t)
	walletService.AssertExpectations(t)

}

func TestService_GenerateSingleUserBalanceForCoin(t *testing.T) {
	qm := userBalanceQueryMatcher{}
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
	db, err := gorm.Open(dialector, &gorm.Config{SkipDefaultTransaction: true})
	dbMock.MatchExpectationsInOrder(false)

	dbMock.ExpectExec("INSERT INTO user_balances").WillReturnResult(sqlmock.NewResult(1, 1))

	cs := new(mocks.CurrencyService)

	repo := new(mocks.UserBalanceRepository)
	pg := new(mocks.PriceGenerator)
	permissionManager := new(mocks.UserPermissionManager)
	walletService := new(mocks.WalletService)
	userService := new(mocks.UserService)
	userWalletBalanceRepository := new(mocks.UserWalletBalanceRepository)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	ubService := userbalance.NewBalanceService(db, repo, cs, pg, permissionManager, walletService, userService,
		userWalletBalanceRepository, configs, logger)

	u := user.User{
		ID: 1,
	}

	coin := currency.Coin{ID: 1, Name: "Bitcoin", Code: "BTC"}
	ub, err := ubService.GenerateSingleUserBalanceForCoin(u, coin)
	assert.Nil(t, err)
	assert.Equal(t, int64(1), ub.ID)
	repo.AssertExpectations(t)
	cs.AssertExpectations(t)

}

func TestService_GenerateAddress(t *testing.T) {
	qm := userBalanceQueryMatcher{}
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
	db, err := gorm.Open(dialector, &gorm.Config{SkipDefaultTransaction: true})
	dbMock.MatchExpectationsInOrder(false)

	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))

	cs := new(mocks.CurrencyService)

	repo := new(mocks.UserBalanceRepository)
	pg := new(mocks.PriceGenerator)
	permissionManager := new(mocks.UserPermissionManager)
	walletService := new(mocks.WalletService)
	walletService.On("GetAddressForUser", "BTC", "").Once().Return("BtcAddress1", nil)
	userService := new(mocks.UserService)
	userWalletBalanceRepository := new(mocks.UserWalletBalanceRepository)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	ubService := userbalance.NewBalanceService(db, repo, cs, pg, permissionManager, walletService, userService,
		userWalletBalanceRepository, configs, logger)

	u := user.User{
		ID: 1,
	}
	ub := userbalance.UserBalance{
		ID: 1,
	}

	coin := currency.Coin{ID: 1, Name: "Bitcoin", Code: "BTC"}
	address, err := ubService.GenerateAddress(ub, u, coin)
	assert.Nil(t, err)
	assert.Equal(t, "BtcAddress1", address)
	repo.AssertExpectations(t)
	cs.AssertExpectations(t)
	walletService.AssertExpectations(t)
}

func TestService_UpsertUserWalletBalance_AlreadyExists(t *testing.T) {
	qm := userBalanceQueryMatcher{}
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
	db, err := gorm.Open(dialector, &gorm.Config{SkipDefaultTransaction: true})

	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectExec("UPDATE user_wallet_balance").WillReturnResult(sqlmock.NewResult(2, 1))

	cs := new(mocks.CurrencyService)
	repo := new(mocks.UserBalanceRepository)
	pg := new(mocks.PriceGenerator)
	permissionManager := new(mocks.UserPermissionManager)
	walletService := new(mocks.WalletService)
	walletService.On("GetAddressBalance", "USDT", "ETH", "someAddress", true).Once().Return("0.1", nil)
	userService := new(mocks.UserService)
	userWalletBalanceRepository := new(mocks.UserWalletBalanceRepository)
	uwb := &userbalance.UserWalletBalance{}
	userWalletBalanceRepository.On("FindUserWalletBalance", 21, int64(1), int64(2), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		uwb = args.Get(3).(*userbalance.UserWalletBalance)
		uwb.ID = 1
	})
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	ubService := userbalance.NewBalanceService(db, repo, cs, pg, permissionManager, walletService, userService,
		userWalletBalanceRepository, configs, logger)

	params := userbalance.UpsertUserWalletBalancesParams{
		UserID:             21,
		CoinID:             1,
		CoinCode:           "USDT",
		BlockchainCoinID:   2,
		BlockchainCoinCode: "ETH",
		Address:            "someAddress",
	}

	err = ubService.UpsertUserWalletBalance(params)
	assert.Nil(t, err)

	repo.AssertExpectations(t)
}

func TestService_UpsertUserWalletBalance_AlreadyDoesNotExist(t *testing.T) {
	qm := userBalanceQueryMatcher{}
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
	db, err := gorm.Open(dialector, &gorm.Config{SkipDefaultTransaction: true})

	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectExec("INSERT INTO user_wallet_balance").WillReturnResult(sqlmock.NewResult(2, 1))

	cs := new(mocks.CurrencyService)
	repo := new(mocks.UserBalanceRepository)
	pg := new(mocks.PriceGenerator)
	permissionManager := new(mocks.UserPermissionManager)
	walletService := new(mocks.WalletService)
	walletService.On("GetAddressBalance", "USDT", "ETH", "someAddress", true).Once().Return("0.1", nil)
	userService := new(mocks.UserService)
	userWalletBalanceRepository := new(mocks.UserWalletBalanceRepository)
	userWalletBalanceRepository.On("FindUserWalletBalance", 21, int64(1), int64(2), mock.Anything).Once().Return(gorm.ErrRecordNotFound)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	ubService := userbalance.NewBalanceService(db, repo, cs, pg, permissionManager, walletService, userService,
		userWalletBalanceRepository, configs, logger)

	params := userbalance.UpsertUserWalletBalancesParams{
		UserID:             21,
		CoinID:             1,
		CoinCode:           "USDT",
		BlockchainCoinID:   2,
		BlockchainCoinCode: "ETH",
		Address:            "someAddress",
	}

	err = ubService.UpsertUserWalletBalance(params)
	assert.Nil(t, err)

	repo.AssertExpectations(t)
}

func TestService_UpdateUserBalanceFromAdmin_NegativeAmount(t *testing.T) {
	db := &gorm.DB{}
	cs := new(mocks.CurrencyService)
	repo := new(mocks.UserBalanceRepository)
	pg := new(mocks.PriceGenerator)
	permissionManager := new(mocks.UserPermissionManager)
	walletService := new(mocks.WalletService)
	walletService.On("GetAddressForUser", "BTC", "").Once().Return("BtcAddress1", nil)
	userService := new(mocks.UserService)
	userWalletBalanceRepository := new(mocks.UserWalletBalanceRepository)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	ubService := userbalance.NewBalanceService(db, repo, cs, pg, permissionManager, walletService, userService,
		userWalletBalanceRepository, configs, logger)

	u := &user.User{
		ID: 1,
	}
	params := userbalance.UpdateUserBalanceFromAdminParams{
		ID:     1,
		Amount: "-1.8",
	}

	res, statusCode := ubService.UpdateUserBalanceFromAdmin(u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "amount is not correct", res.Message)
}

func TestService_UpdateUserBalanceFromAdmin_MoreThanCurrentBalanceAmount(t *testing.T) {
	qm := userBalanceQueryMatcher{}
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
	db, err := gorm.Open(dialector, &gorm.Config{SkipDefaultTransaction: true})

	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(2, 1))
	dbMock.ExpectCommit()

	cs := new(mocks.CurrencyService)
	repo := new(mocks.UserBalanceRepository)
	ub := &userbalance.UserBalance{}
	repo.On("GetUserBalanceByIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub = args.Get(2).(*userbalance.UserBalance)
		ub.ID = 1
		ub.UserID = 1
		ub.Amount = "1"
		ub.FrozenAmount = "0.1"
		ub.Coin = currency.Coin{
			ID:   1,
			Code: "ETH",
		}
		ub.CoinID = 2
	})

	pg := new(mocks.PriceGenerator)
	permissionManager := new(mocks.UserPermissionManager)
	walletService := new(mocks.WalletService)
	userService := new(mocks.UserService)
	userWalletBalanceRepository := new(mocks.UserWalletBalanceRepository)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	ubService := userbalance.NewBalanceService(db, repo, cs, pg, permissionManager, walletService, userService,
		userWalletBalanceRepository, configs, logger)

	u := &user.User{
		ID: 1,
	}
	params := userbalance.UpdateUserBalanceFromAdminParams{
		ID:     1,
		Amount: "1.8",
	}

	res, statusCode := ubService.UpdateUserBalanceFromAdmin(u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	assert.Equal(t, "1.80000000", ub.Amount)

	repo.AssertExpectations(t)
}

func TestService_UpdateUserBalanceFromAdmin_LessThanCurrentBalanceAmount(t *testing.T) {
	qm := userBalanceQueryMatcher{}
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
	db, err := gorm.Open(dialector, &gorm.Config{SkipDefaultTransaction: true})

	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO transactions").WillReturnResult(sqlmock.NewResult(2, 1))
	dbMock.ExpectCommit()

	cs := new(mocks.CurrencyService)
	repo := new(mocks.UserBalanceRepository)
	ub := &userbalance.UserBalance{}
	repo.On("GetUserBalanceByIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub = args.Get(2).(*userbalance.UserBalance)
		ub.ID = 1
		ub.UserID = 1
		ub.Amount = "1"
		ub.FrozenAmount = "0.1"
		ub.Coin = currency.Coin{
			ID:   1,
			Code: "ETH",
		}
		ub.CoinID = 2
	})

	pg := new(mocks.PriceGenerator)
	permissionManager := new(mocks.UserPermissionManager)
	walletService := new(mocks.WalletService)
	userService := new(mocks.UserService)
	userWalletBalanceRepository := new(mocks.UserWalletBalanceRepository)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	ubService := userbalance.NewBalanceService(db, repo, cs, pg, permissionManager, walletService, userService,
		userWalletBalanceRepository, configs, logger)

	u := &user.User{
		ID: 1,
	}
	params := userbalance.UpdateUserBalanceFromAdminParams{
		ID:     1,
		Amount: "0.8",
	}

	res, statusCode := ubService.UpdateUserBalanceFromAdmin(u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	assert.Equal(t, "0.80000000", ub.Amount)

	repo.AssertExpectations(t)
}

func TestService_SetAutoExchangeCoin_CodeNotFound(t *testing.T) {
	db := &gorm.DB{}
	cs := new(mocks.CurrencyService)
	cs.On("GetCoinByCode", "NOTEXISTING").Once().Return(currency.Coin{}, gorm.ErrRecordNotFound)
	repo := new(mocks.UserBalanceRepository)
	pg := new(mocks.PriceGenerator)
	permissionManager := new(mocks.UserPermissionManager)
	walletService := new(mocks.WalletService)
	userService := new(mocks.UserService)
	userWalletBalanceRepository := new(mocks.UserWalletBalanceRepository)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	ubService := userbalance.NewBalanceService(db, repo, cs, pg, permissionManager, walletService, userService,
		userWalletBalanceRepository, configs, logger)

	u := &user.User{
		ID: 1,
	}
	params := userbalance.SetAutoExchangeCoinParams{
		Code:             "notExisting",
		AutoExchangeCode: "BTC",
		Mode:             "add",
	}

	res, statusCode := ubService.SetAutoExchangeCoin(u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "code not found", res.Message)

	cs.AssertExpectations(t)
}

func TestService_SetAutoExchangeCoin_ExchangeCodeNotFound(t *testing.T) {
	db := &gorm.DB{}
	cs := new(mocks.CurrencyService)
	cs.On("GetCoinByCode", "BTC").Once().Return(currency.Coin{ID: 1}, nil)
	cs.On("GetCoinByCode", "NOTEXISTING").Once().Return(currency.Coin{}, gorm.ErrRecordNotFound)
	repo := new(mocks.UserBalanceRepository)
	pg := new(mocks.PriceGenerator)
	permissionManager := new(mocks.UserPermissionManager)
	walletService := new(mocks.WalletService)
	userService := new(mocks.UserService)
	userWalletBalanceRepository := new(mocks.UserWalletBalanceRepository)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	ubService := userbalance.NewBalanceService(db, repo, cs, pg, permissionManager, walletService, userService,
		userWalletBalanceRepository, configs, logger)

	u := &user.User{
		ID: 1,
	}
	params := userbalance.SetAutoExchangeCoinParams{
		Code:             "BTC",
		AutoExchangeCode: "notExisting",
		Mode:             "add",
	}

	res, statusCode := ubService.SetAutoExchangeCoin(u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "auto exchange code not found", res.Message)

	cs.AssertExpectations(t)
}

func TestService_SetAutoExchangeCoin_Add_PairNotFound(t *testing.T) {
	db := &gorm.DB{}
	cs := new(mocks.CurrencyService)
	cs.On("GetCoinByCode", "DAI").Once().Return(currency.Coin{ID: 1}, nil)
	cs.On("GetCoinByCode", "GRS").Once().Return(currency.Coin{ID: 3}, nil)
	pairs := []currency.Pair{
		{
			ID:              1,
			Name:            "BTC-USDT",
			BasisCoinID:     1,
			DependentCoinID: 2,
		},
	}
	cs.On("GetActivePairCurrenciesList").Once().Return(pairs)
	repo := new(mocks.UserBalanceRepository)
	ub := &userbalance.UserBalance{}
	repo.On("GetBalanceOfUserByCoinID", 1, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub = args.Get(2).(*userbalance.UserBalance)
		ub.ID = 1
		ub.AutoExchangeCoin = sql.NullString{String: "", Valid: false}
	})
	pg := new(mocks.PriceGenerator)
	permissionManager := new(mocks.UserPermissionManager)
	walletService := new(mocks.WalletService)
	userService := new(mocks.UserService)
	userWalletBalanceRepository := new(mocks.UserWalletBalanceRepository)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	ubService := userbalance.NewBalanceService(db, repo, cs, pg, permissionManager, walletService, userService,
		userWalletBalanceRepository, configs, logger)

	u := &user.User{
		ID: 1,
	}
	params := userbalance.SetAutoExchangeCoinParams{
		Code:             "DAI",
		AutoExchangeCode: "GRS",
		Mode:             "add",
	}

	res, statusCode := ubService.SetAutoExchangeCoin(u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "pair not found with these two coins", res.Message)

	cs.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestService_SetAutoExchangeCoin_Add_Successful(t *testing.T) {
	qm := userBalanceQueryMatcher{}
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
	db, err := gorm.Open(dialector, &gorm.Config{SkipDefaultTransaction: true})

	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))

	cs := new(mocks.CurrencyService)
	cs.On("GetCoinByCode", "BTC").Once().Return(currency.Coin{ID: 2}, nil)
	cs.On("GetCoinByCode", "USDT").Once().Return(currency.Coin{ID: 1}, nil)
	pairs := []currency.Pair{
		{
			ID:              1,
			Name:            "BTC-USDT",
			BasisCoinID:     1,
			DependentCoinID: 2,
		},
	}
	cs.On("GetActivePairCurrenciesList").Once().Return(pairs)
	repo := new(mocks.UserBalanceRepository)
	ub := &userbalance.UserBalance{}
	repo.On("GetBalanceOfUserByCoinID", 1, int64(2), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub = args.Get(2).(*userbalance.UserBalance)
		ub.ID = 1
		ub.AutoExchangeCoin = sql.NullString{String: "", Valid: false}
	})
	pg := new(mocks.PriceGenerator)
	permissionManager := new(mocks.UserPermissionManager)
	walletService := new(mocks.WalletService)
	userService := new(mocks.UserService)
	userWalletBalanceRepository := new(mocks.UserWalletBalanceRepository)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	ubService := userbalance.NewBalanceService(db, repo, cs, pg, permissionManager, walletService, userService,
		userWalletBalanceRepository, configs, logger)

	u := &user.User{
		ID: 1,
	}
	params := userbalance.SetAutoExchangeCoinParams{
		Code:             "BTC",
		AutoExchangeCode: "USDT",
		Mode:             "add",
	}

	res, statusCode := ubService.SetAutoExchangeCoin(u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	assert.Equal(t, true, ub.AutoExchangeCoin.Valid)
	assert.Equal(t, "USDT", ub.AutoExchangeCoin.String)

	cs.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestService_SetAutoExchangeCoin_Remove_Successful(t *testing.T) {
	qm := userBalanceQueryMatcher{}
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
	db, err := gorm.Open(dialector, &gorm.Config{SkipDefaultTransaction: true})

	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))

	cs := new(mocks.CurrencyService)
	cs.On("GetCoinByCode", "BTC").Once().Return(currency.Coin{ID: 2}, nil)
	cs.On("GetCoinByCode", "USDT").Once().Return(currency.Coin{ID: 1}, nil)
	repo := new(mocks.UserBalanceRepository)
	ub := &userbalance.UserBalance{}
	repo.On("GetBalanceOfUserByCoinID", 1, int64(2), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub = args.Get(2).(*userbalance.UserBalance)
		ub.ID = 1
		ub.AutoExchangeCoin = sql.NullString{String: "USDT", Valid: true}
	})
	pg := new(mocks.PriceGenerator)
	permissionManager := new(mocks.UserPermissionManager)
	walletService := new(mocks.WalletService)
	userService := new(mocks.UserService)
	userWalletBalanceRepository := new(mocks.UserWalletBalanceRepository)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	ubService := userbalance.NewBalanceService(db, repo, cs, pg, permissionManager, walletService, userService,
		userWalletBalanceRepository, configs, logger)

	u := &user.User{
		ID: 1,
	}
	params := userbalance.SetAutoExchangeCoinParams{
		Code:             "BTC",
		AutoExchangeCode: "USDT",
		Mode:             "delete",
	}

	res, statusCode := ubService.SetAutoExchangeCoin(u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)

	assert.Equal(t, false, ub.AutoExchangeCoin.Valid)
	assert.Equal(t, "", ub.AutoExchangeCoin.String)

	cs.AssertExpectations(t)
	repo.AssertExpectations(t)
}
