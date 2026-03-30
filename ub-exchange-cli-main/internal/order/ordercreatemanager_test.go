// Package order_test tests the OrderCreateManager. Covers:
//   - Validation: missing or negative amount, missing or negative price for LIMIT orders
//   - Validation: negative stop point price for stop-limit orders
//   - Validation: order amount below minimum pair threshold
//   - Balance check: insufficient user balance to place order
//   - Level check: user exchange level too low for the requested order volume
//   - Successful creation of LIMIT, MARKET, and stop-limit orders with correct field calculations
//
// Test data: sqlmock MySQL DB, mocked user balance service, user level service,
// price generator, and currency pair/coin fixtures for BTC-USDT.
package order_test

import (
	"database/sql"
	"exchange-go/internal/currency"
	"exchange-go/internal/mocks"
	"exchange-go/internal/order"
	"exchange-go/internal/user"
	"exchange-go/internal/userbalance"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// orderCreateQueryMatcher is a sqlmock.QueryMatcher that accepts any SQL query,
// allowing tests to focus on order creation logic rather than exact SQL strings.
type orderCreateQueryMatcher struct {
}

// Match always returns nil, effectively disabling SQL query matching so that
// any expected SQL statement is considered a match against any actual SQL.
func (orderCreateQueryMatcher) Match(expectedSQL, actualSQL string) error {
	return nil
}

func TestCreateManager_CreateOrder_NoAmountProvidedOrNegativeAmount(t *testing.T) {
	qm := orderCreateQueryMatcher{}
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

	userBalanceService := new(mocks.UserBalanceService)
	userLevelService := new(mocks.UserLevelService)
	priceGenerator := new(mocks.PriceGenerator)

	ocm := order.NewOrderCreateManager(db, userBalanceService, userLevelService, priceGenerator)
	uai := order.UserAgentInfo{
		Device:  "web",
		IP:      "127.0.0.1",
		Browser: "chrome",
	}
	data := order.CreateRequiredData{
		User:           &user.User{},
		Pair:           &currency.Pair{},
		Amount:         "",
		OrderType:      "BUY",
		ExchangeType:   "LIMIT",
		Price:          "5000",
		UserAgentInfo:  uai,
		StopPointPrice: "",
		CurrentPrice:   "50000",
		IsInstant:      true,
	}

	_, err = ocm.CreateOrder(data)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "amount is not valid")

	//negative amount

	data.Amount = "-1.5"
	_, err = ocm.CreateOrder(data)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "amount is not valid")

	userBalanceService.AssertExpectations(t)
}

func TestCreateManager_CreateOrder_Limit_NoPriceProvidedOrNegativePrice(t *testing.T) {
	qm := orderCreateQueryMatcher{}
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

	userBalanceService := new(mocks.UserBalanceService)
	userLevelService := new(mocks.UserLevelService)
	priceGenerator := new(mocks.PriceGenerator)

	ocm := order.NewOrderCreateManager(db, userBalanceService, userLevelService, priceGenerator)
	uai := order.UserAgentInfo{
		Device:  "web",
		IP:      "127.0.0.1",
		Browser: "chrome",
	}
	data := order.CreateRequiredData{
		User:           &user.User{},
		Pair:           &currency.Pair{},
		Amount:         "1.0",
		OrderType:      "BUY",
		ExchangeType:   "LIMIT",
		Price:          "",
		UserAgentInfo:  uai,
		StopPointPrice: "",
		CurrentPrice:   "50000",
		IsInstant:      true,
	}

	_, err = ocm.CreateOrder(data)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "price is not valid")

	//negative price
	data.Price = "-50000"
	_, err = ocm.CreateOrder(data)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "price is not valid")

}

func TestCreateManager_CreateOrder_StopLimit_NegativeStopPointPriceProvided(t *testing.T) {
	qm := orderCreateQueryMatcher{}
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

	userBalanceService := new(mocks.UserBalanceService)
	userLevelService := new(mocks.UserLevelService)
	priceGenerator := new(mocks.PriceGenerator)

	ocm := order.NewOrderCreateManager(db, userBalanceService, userLevelService, priceGenerator)
	uai := order.UserAgentInfo{
		Device:  "web",
		IP:      "127.0.0.1",
		Browser: "chrome",
	}
	data := order.CreateRequiredData{
		User:           &user.User{},
		Pair:           &currency.Pair{},
		Amount:         "1.0",
		OrderType:      "BUY",
		ExchangeType:   "LIMIT",
		Price:          "50000",
		UserAgentInfo:  uai,
		StopPointPrice: "-52000",
		CurrentPrice:   "50000",
		IsInstant:      true,
	}

	_, err = ocm.CreateOrder(data)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "stop point price is not valid")
}

func TestCreateManager_CreateOrder_LessThanMinimumOrderAmount(t *testing.T) {
	qm := orderCreateQueryMatcher{}
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

	userBalanceService := new(mocks.UserBalanceService)
	userLevelService := new(mocks.UserLevelService)
	priceGenerator := new(mocks.PriceGenerator)

	ocm := order.NewOrderCreateManager(db, userBalanceService, userLevelService, priceGenerator)
	uai := order.UserAgentInfo{
		Device:  "web",
		IP:      "127.0.0.1",
		Browser: "chrome",
	}

	basisCoin := &currency.Coin{
		ID:   1,
		Code: "USDT",
	}

	pair := &currency.Pair{
		ID:                 1,
		Name:               "BTC-USDT",
		MinimumOrderAmount: sql.NullString{String: "10", Valid: true},
		BasisCoin:          *basisCoin,
	}
	data := order.CreateRequiredData{
		User:           &user.User{},
		Pair:           pair,
		Amount:         "0.00001",
		OrderType:      "BUY",
		ExchangeType:   "LIMIT",
		Price:          "50000",
		UserAgentInfo:  uai,
		StopPointPrice: "",
		CurrentPrice:   "50000",
		IsInstant:      true,
	}

	_, err = ocm.CreateOrder(data)
	assert.NotNil(t, err)
	assert.Equal(t, "the minimum order amount must be more than 10 USDT", err.Error())
}

func TestCreateManager_CreateOrder_UserBalanceIsNotEnough(t *testing.T) {
	qm := orderCreateQueryMatcher{}
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
	dbMock.ExpectBegin()
	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectRollback()

	userBalanceService := new(mocks.UserBalanceService)
	userBalanceService.On("GetBalanceOfUserByCoinUsingTx", mock.Anything, 1, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub := args.Get(3).(*userbalance.UserBalance)
		ub.ID = 1
		ub.UserID = 1
		ub.Amount = "60000"
		ub.FrozenAmount = "30000"
	})
	userLevelService := new(mocks.UserLevelService)
	ul := user.Level{
		ID:                        1,
		ExchangeVolumeLimitAmount: "1",
		ExchangeNumberLimit:       100,
	}
	userLevelService.On("GetLevelByID", int64(1)).Once().Return(ul, nil)
	priceGenerator := new(mocks.PriceGenerator)
	priceGenerator.On("GetAmountBasedOnBTC", mock.Anything, "BTC", "1.00000000").Once().Return("1.0", nil)
	ocm := order.NewOrderCreateManager(db, userBalanceService, userLevelService, priceGenerator)
	uai := order.UserAgentInfo{
		Device:  "web",
		IP:      "127.0.0.1",
		Browser: "chrome",
	}
	basisCoin := currency.Coin{
		ID:   1,
		Name: "USDT",
		Code: "USDT",
	}
	dependentCoin := currency.Coin{
		ID:   2,
		Name: "BTC",
		Code: "BTC",
	}

	pair := &currency.Pair{
		ID:                 1,
		Name:               "BTC-USDT",
		MinimumOrderAmount: sql.NullString{String: "10", Valid: true},
		BasisCoin:          basisCoin,
		DependentCoin:      dependentCoin,
	}

	user := &user.User{ID: 1, UserLevelID: 1}
	data := order.CreateRequiredData{
		User:           user,
		Pair:           pair,
		Amount:         "1.0",
		OrderType:      "BUY",
		ExchangeType:   "LIMIT",
		Price:          "50000",
		UserAgentInfo:  uai,
		StopPointPrice: "",
		CurrentPrice:   "50000",
		IsInstant:      true,
	}

	_, err = ocm.CreateOrder(data)
	assert.NotNil(t, err)
	assert.ErrorContains(t, err, "user balance is not enough")

	userBalanceService.AssertExpectations(t)

}

func TestCreateManager_CreateOrder_UserLevelDoesNotAllows(t *testing.T) {
	qm := orderCreateQueryMatcher{}
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

	userBalanceService := new(mocks.UserBalanceService)
	ul := user.Level{
		ID:                        1,
		ExchangeVolumeLimitAmount: "1",
		ExchangeNumberLimit:       100,
	}
	userLevelService := new(mocks.UserLevelService)
	userLevelService.On("GetLevelByID", int64(1)).Once().Return(ul, nil)

	priceGenerator := new(mocks.PriceGenerator)
	priceGenerator.On("GetAmountBasedOnBTC", mock.Anything, "BTC", "1.00000000").Once().Return("1.0", nil)

	ocm := order.NewOrderCreateManager(db, userBalanceService, userLevelService, priceGenerator)
	uai := order.UserAgentInfo{
		Device:  "web",
		IP:      "127.0.0.1",
		Browser: "chrome",
	}
	basisCoin := currency.Coin{
		ID:   1,
		Code: "USDT",
	}
	dependentCoin := currency.Coin{
		ID:   2,
		Code: "BTC",
	}

	pair := &currency.Pair{
		ID:                 1,
		Name:               "BTC-USDT",
		MinimumOrderAmount: sql.NullString{String: "10", Valid: true},
		BasisCoin:          basisCoin,
		DependentCoin:      dependentCoin,
	}

	u := &user.User{
		ID:                   1,
		UserLevelID:          1,
		Kyc:                  user.KycLevelMinimum,
		ExchangeVolumeAmount: "1.0", //in btc
		ExchangeNumber:       20,
	}
	data := order.CreateRequiredData{
		User:           u,
		Pair:           pair,
		Amount:         "1.0",
		OrderType:      "BUY",
		ExchangeType:   "LIMIT",
		Price:          "50000",
		UserAgentInfo:  uai,
		StopPointPrice: "",
		CurrentPrice:   "50000",
		IsInstant:      true,
	}

	_, err = ocm.CreateOrder(data)
	assert.NotNil(t, err)
	assert.Equal(t, "your user level is low to place this order. please verify your identity to boost up your level", err.Error())

	userLevelService.AssertExpectations(t)
	priceGenerator.AssertExpectations(t)
}

func TestCreateManager_CreateOrder_Limit_Successful(t *testing.T) {
	qm := orderCreateQueryMatcher{}
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
	dbMock.ExpectExec("INSERT INTO orders").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO orders_extra_info").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE orders").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit()

	userBalanceService := new(mocks.UserBalanceService)
	userBalanceService.On("GetBalanceOfUserByCoinUsingTx", mock.Anything, 1, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub := args.Get(3).(*userbalance.UserBalance)
		ub.ID = 1
		ub.UserID = 1
		ub.Amount = "60000"
		ub.FrozenAmount = "0"
	})

	ul := user.Level{
		ID:                        1,
		ExchangeVolumeLimitAmount: "1",
		ExchangeNumberLimit:       100,
	}
	userLevelService := new(mocks.UserLevelService)
	userLevelService.On("GetLevelByID", int64(1)).Once().Return(ul, nil)

	priceGenerator := new(mocks.PriceGenerator)
	priceGenerator.On("GetAmountBasedOnBTC", mock.Anything, "BTC", "1.00000000").Once().Return("1.0", nil)

	ocm := order.NewOrderCreateManager(db, userBalanceService, userLevelService, priceGenerator)
	uai := order.UserAgentInfo{
		Device:  "web",
		IP:      "127.0.0.1",
		Browser: "chrome",
	}
	basisCoin := currency.Coin{
		ID:   1,
		Code: "USDT",
	}
	dependentCoin := currency.Coin{
		ID:   2,
		Code: "BTC",
	}

	pair := &currency.Pair{
		ID:                 1,
		Name:               "BTC-USDT",
		MinimumOrderAmount: sql.NullString{String: "10", Valid: true},
		BasisCoin:          basisCoin,
		DependentCoin:      dependentCoin,
	}

	u := &user.User{
		ID:                   1,
		UserLevelID:          1,
		Kyc:                  user.KycLevelMinimum,
		ExchangeVolumeAmount: "0.0", //in btc
		ExchangeNumber:       20,
	}
	data := order.CreateRequiredData{
		User:           u,
		Pair:           pair,
		Amount:         "1.0",
		OrderType:      "BUY",
		ExchangeType:   "LIMIT",
		Price:          "50000",
		UserAgentInfo:  uai,
		StopPointPrice: "",
		CurrentPrice:   "50000",
		IsInstant:      true,
	}

	o, err := ocm.CreateOrder(data)
	assert.Nil(t, err)
	assert.Equal(t, "50000.00000000", o.Price.String)
	assert.Equal(t, "1.00000000", o.DemandedAmount.String)
	assert.Equal(t, "50000.00000000", o.PayedByAmount.String)
	assert.Equal(t, 1, o.UserID)
	assert.Equal(t, "BUY", o.Type)
	assert.Equal(t, "LIMIT", o.ExchangeType)

	userBalanceService.AssertExpectations(t)
	userLevelService.AssertExpectations(t)
	priceGenerator.AssertExpectations(t)
}

func TestCreateManager_CreateOrder_Market_Successful(t *testing.T) {
	qm := orderCreateQueryMatcher{}
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
	dbMock.ExpectExec("INSERT INTO orders").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO orders_extra_info").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE orders").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit()

	userBalanceService := new(mocks.UserBalanceService)
	userBalanceService.On("GetBalanceOfUserByCoinUsingTx", mock.Anything, 1, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub := args.Get(3).(*userbalance.UserBalance)
		ub.ID = 1
		ub.UserID = 1
		ub.Amount = "60000"
		ub.FrozenAmount = "0"
	})

	ul := user.Level{
		ID:                        1,
		ExchangeVolumeLimitAmount: "1",
		ExchangeNumberLimit:       100,
	}
	userLevelService := new(mocks.UserLevelService)
	userLevelService.On("GetLevelByID", int64(1)).Once().Return(ul, nil)

	priceGenerator := new(mocks.PriceGenerator)
	priceGenerator.On("GetAmountBasedOnBTC", mock.Anything, "BTC", "0.00200000").Once().Return("1.0", nil)

	ocm := order.NewOrderCreateManager(db, userBalanceService, userLevelService, priceGenerator)
	uai := order.UserAgentInfo{
		Device:  "web",
		IP:      "127.0.0.1",
		Browser: "chrome",
	}
	basisCoin := currency.Coin{
		ID:   1,
		Code: "USDT",
	}
	dependentCoin := currency.Coin{
		ID:   2,
		Code: "BTC",
	}

	pair := &currency.Pair{
		ID:                 1,
		Name:               "BTC-USDT",
		MinimumOrderAmount: sql.NullString{String: "10", Valid: true},
		BasisCoin:          basisCoin,
		DependentCoin:      dependentCoin,
	}

	u := &user.User{
		ID:                   1,
		UserLevelID:          1,
		Kyc:                  user.KycLevelMinimum,
		ExchangeVolumeAmount: "0.0", //in btc
		ExchangeNumber:       20,
	}
	data := order.CreateRequiredData{
		User:           u,
		Pair:           pair,
		Amount:         "100",
		OrderType:      "BUY",
		ExchangeType:   "MARKET",
		Price:          "",
		UserAgentInfo:  uai,
		StopPointPrice: "",
		CurrentPrice:   "50000",
		IsInstant:      true,
	}

	o, err := ocm.CreateOrder(data)
	assert.Nil(t, err)
	assert.Equal(t, "", o.Price.String)
	assert.Equal(t, "0.00200000", o.DemandedAmount.String)
	assert.Equal(t, "100.00000000", o.PayedByAmount.String)
	assert.Equal(t, 1, o.UserID)
	assert.Equal(t, "BUY", o.Type)
	assert.Equal(t, "MARKET", o.ExchangeType)

	userBalanceService.AssertExpectations(t)
	userLevelService.AssertExpectations(t)
	priceGenerator.AssertExpectations(t)

}

func TestCreateManager_CreateOrder_StopOrder_Successful(t *testing.T) {
	qm := orderCreateQueryMatcher{}
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
	dbMock.ExpectExec("INSERT INTO orders").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("INSERT INTO orders_extra_info").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE orders").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit()

	userBalanceService := new(mocks.UserBalanceService)
	userBalanceService.On("GetBalanceOfUserByCoinUsingTx", mock.Anything, 1, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub := args.Get(3).(*userbalance.UserBalance)
		ub.ID = 1
		ub.UserID = 1
		ub.Amount = "60000"
		ub.FrozenAmount = "0"
	})

	ul := user.Level{
		ID:                        1,
		ExchangeVolumeLimitAmount: "1",
		ExchangeNumberLimit:       100,
	}
	userLevelService := new(mocks.UserLevelService)
	userLevelService.On("GetLevelByID", int64(1)).Once().Return(ul, nil)

	priceGenerator := new(mocks.PriceGenerator)
	priceGenerator.On("GetAmountBasedOnBTC", mock.Anything, "BTC", "1.00000000").Once().Return("1.0", nil)

	ocm := order.NewOrderCreateManager(db, userBalanceService, userLevelService, priceGenerator)
	uai := order.UserAgentInfo{
		Device:  "web",
		IP:      "127.0.0.1",
		Browser: "chrome",
	}
	basisCoin := currency.Coin{
		ID:   1,
		Code: "USDT",
	}
	dependentCoin := currency.Coin{
		ID:   2,
		Code: "BTC",
	}

	pair := &currency.Pair{
		ID:                 1,
		Name:               "BTC-USDT",
		MinimumOrderAmount: sql.NullString{String: "10", Valid: true},
		BasisCoin:          basisCoin,
		DependentCoin:      dependentCoin,
	}

	u := &user.User{
		ID:                   1,
		UserLevelID:          1,
		Kyc:                  user.KycLevelMinimum,
		ExchangeVolumeAmount: "0.0", //in btc
		ExchangeNumber:       20,
	}
	data := order.CreateRequiredData{
		User:           u,
		Pair:           pair,
		Amount:         "1.0",
		OrderType:      "BUY",
		ExchangeType:   "LIMIT",
		Price:          "51000",
		UserAgentInfo:  uai,
		StopPointPrice: "50000",
		CurrentPrice:   "50000",
		IsInstant:      true,
	}

	o, err := ocm.CreateOrder(data)
	assert.Nil(t, err)
	assert.Equal(t, "51000.00000000", o.Price.String)
	assert.Equal(t, "50000.00000000", o.StopPointPrice.String)
	assert.Equal(t, "1.00000000", o.DemandedAmount.String)
	assert.Equal(t, "51000.00000000", o.PayedByAmount.String)
	assert.Equal(t, 1, o.UserID)
	assert.Equal(t, "BUY", o.Type)
	assert.Equal(t, "LIMIT", o.ExchangeType)

	userBalanceService.AssertExpectations(t)
	userLevelService.AssertExpectations(t)
	priceGenerator.AssertExpectations(t)

}
