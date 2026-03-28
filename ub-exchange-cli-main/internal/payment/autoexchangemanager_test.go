// Package payment_test tests the AutoExchangeManager which automatically
// converts deposited funds into a user's preferred coin. Covers:
//   - Skipping auto-exchange when disabled, payment is not a deposit, or
//     payment status is not completed
//   - Successful sell auto-exchange (deposit coin is the dependent coin of the pair)
//   - Successful buy auto-exchange (deposit coin is the basis coin of the pair)
//   - Exception error handling when price retrieval fails (failure type: exception)
//   - Logical error handling when order creation fails with a validation error
//     (failure type: logical, with reason persisted in ExtraInfo)
//
// Test data: mocked PaymentRepository, OrderCreateManager, EventsHandler,
// UserService, CurrencyService, PriceGenerator, and Logger; go-sqlmock for
// GORM database interactions with the crypto_payment_extra_info table.
package payment_test

import (
	"database/sql"
	"exchange-go/internal/currency"
	"exchange-go/internal/mocks"
	"exchange-go/internal/order"
	"exchange-go/internal/payment"
	"exchange-go/internal/platform"
	"exchange-go/internal/user"
	"exchange-go/internal/userbalance"
	"fmt"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestAutoExchangeManager_ShouldNotAutoExchange(t *testing.T) {
	db := &gorm.DB{}
	paymentRepository := new(mocks.PaymentRepository)
	orderCreateManager := new(mocks.OrderCreateManager)
	orderEventsHandler := new(mocks.EventsHandler)
	userService := new(mocks.UserService)
	currencyService := new(mocks.CurrencyService)
	priceGenerator := new(mocks.PriceGenerator)
	logger := new(mocks.Logger)
	autoExchangeManger := payment.NewAutoExchangeManger(db, paymentRepository, orderCreateManager, orderEventsHandler, userService,
		currencyService, priceGenerator, logger)

	// autoexchange is not enabled
	p := &payment.Payment{}
	ub := &userbalance.UserBalance{
		AutoExchangeCoin: sql.NullString{String: "", Valid: false},
	}
	autoExchangeManger.AutoExchange(p, ub)

	// payment is not deposit
	p = &payment.Payment{
		Type:   payment.TypeWithdraw,
		Status: payment.StatusCompleted,
	}
	ub = &userbalance.UserBalance{
		AutoExchangeCoin: sql.NullString{String: "BTC", Valid: true},
	}
	autoExchangeManger.AutoExchange(p, ub)

	// payment status is not completed
	p = &payment.Payment{
		Type:   payment.TypeDeposit,
		Status: payment.StatusCreated,
	}
	ub = &userbalance.UserBalance{
		AutoExchangeCoin: sql.NullString{String: "BTC", Valid: true},
	}
	autoExchangeManger.AutoExchange(p, ub)
}

func TestAutoExchangeManager_Successful_Sell(t *testing.T) {
	user := user.User{
		ID: 1,
	}
	pair := currency.Pair{
		ID:            1,
		Name:          "BTC-USDT",
		BasisCoin:     currency.Coin{Code: "USDT"},
		DependentCoin: currency.Coin{Code: "BTC"},
	}
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
	db, err := gorm.Open(dialector, &gorm.Config{SkipDefaultTransaction: true})
	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectExec("UPDATE crypto_payment_extra_info").WillReturnResult(sqlmock.NewResult(1, 1))
	paymentRepository := new(mocks.PaymentRepository)
	extraInfo := &payment.ExtraInfo{}
	paymentRepository.On("GetExtraInfoByPaymentID", int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		extraInfo = args.Get(1).(*payment.ExtraInfo)
		extraInfo.PaymentID = 1
		extraInfo.IP = sql.NullString{String: "127.0.0.1", Valid: true}
	})
	orderCreateManager := new(mocks.OrderCreateManager)
	requiredData := order.CreateRequiredData{
		User:         &user,
		Pair:         &pair,
		Amount:       "0.1",
		OrderType:    order.TypeSell,
		ExchangeType: order.ExchangeTypeMarket,
		Price:        "",
		UserAgentInfo: order.UserAgentInfo{
			Device:  "",
			IP:      "127.0.0.1",
			Browser: "",
		},
		StopPointPrice: "",
		CurrentPrice:   "50000.0",
		IsInstant:      true,
		IsFastExchange: false,
		IsAutoExchange: true,
	}
	o := &order.Order{}
	orderCreateManager.On("CreateOrder", requiredData).Once().Return(o, nil)
	orderEventsHandler := new(mocks.EventsHandler)
	orderEventsHandler.On("HandleOrderCreation", mock.Anything, false).Once().Return()
	userService := new(mocks.UserService)
	userService.On("GetUserByID", 1).Once().Return(user, nil)
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetActivePairCurrenciesList").Once().Return([]currency.Pair{pair})
	priceGenerator := new(mocks.PriceGenerator)
	priceGenerator.On("GetPrice", mock.Anything, "BTC-USDT").Once().Return("50000.0", nil)
	logger := new(mocks.Logger)
	autoExchangeManger := payment.NewAutoExchangeManger(db, paymentRepository, orderCreateManager, orderEventsHandler, userService,
		currencyService, priceGenerator, logger)
	p := &payment.Payment{
		ID:     1,
		Type:   payment.TypeDeposit,
		Status: payment.StatusCompleted,
		Amount: sql.NullString{String: "0.1", Valid: true},
		UserID: 1,
	}
	ub := &userbalance.UserBalance{
		BalanceCoin:      "BTC",
		AutoExchangeCoin: sql.NullString{String: "USDT", Valid: true},
	}
	autoExchangeManger.AutoExchange(p, ub)

	time.Sleep(20 * time.Millisecond)
	paymentRepository.AssertExpectations(t)
	orderCreateManager.AssertExpectations(t)
	orderEventsHandler.AssertExpectations(t)
	userService.AssertExpectations(t)
	currencyService.AssertExpectations(t)
	priceGenerator.AssertExpectations(t)
}

func TestAutoExchangeManager_Successful_Buy(t *testing.T) {
	user := user.User{
		ID: 1,
	}
	pair := currency.Pair{
		ID:            1,
		Name:          "BTC-USDT",
		BasisCoin:     currency.Coin{Code: "USDT"},
		DependentCoin: currency.Coin{Code: "BTC"},
	}
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
	db, err := gorm.Open(dialector, &gorm.Config{SkipDefaultTransaction: true})
	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectExec("UPDATE crypto_payment_extra_info").WillReturnResult(sqlmock.NewResult(1, 1))
	paymentRepository := new(mocks.PaymentRepository)
	extraInfo := &payment.ExtraInfo{}
	paymentRepository.On("GetExtraInfoByPaymentID", int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		extraInfo = args.Get(1).(*payment.ExtraInfo)
		extraInfo.PaymentID = 1
		extraInfo.IP = sql.NullString{String: "127.0.0.1", Valid: true}
	})
	orderCreateManager := new(mocks.OrderCreateManager)
	requiredData := order.CreateRequiredData{
		User:         &user,
		Pair:         &pair,
		Amount:       "100.0",
		OrderType:    order.TypeBuy,
		ExchangeType: order.ExchangeTypeMarket,
		Price:        "",
		UserAgentInfo: order.UserAgentInfo{
			Device:  "",
			IP:      "127.0.0.1",
			Browser: "",
		},
		StopPointPrice: "",
		CurrentPrice:   "50000.0",
		IsInstant:      true,
		IsFastExchange: false,
		IsAutoExchange: true,
	}
	o := &order.Order{}
	orderCreateManager.On("CreateOrder", requiredData).Once().Return(o, nil)
	orderEventsHandler := new(mocks.EventsHandler)
	orderEventsHandler.On("HandleOrderCreation", mock.Anything, false).Once().Return()
	userService := new(mocks.UserService)
	userService.On("GetUserByID", 1).Once().Return(user, nil)
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetActivePairCurrenciesList").Once().Return([]currency.Pair{pair})
	priceGenerator := new(mocks.PriceGenerator)
	priceGenerator.On("GetPrice", mock.Anything, "BTC-USDT").Once().Return("50000.0", nil)
	logger := new(mocks.Logger)
	autoExchangeManger := payment.NewAutoExchangeManger(db, paymentRepository, orderCreateManager, orderEventsHandler, userService,
		currencyService, priceGenerator, logger)
	p := &payment.Payment{
		ID:     1,
		Type:   payment.TypeDeposit,
		Status: payment.StatusCompleted,
		Amount: sql.NullString{String: "100.0", Valid: true},
		UserID: 1,
	}
	ub := &userbalance.UserBalance{
		BalanceCoin:      "USDT",
		AutoExchangeCoin: sql.NullString{String: "BTC", Valid: true},
	}
	autoExchangeManger.AutoExchange(p, ub)
	time.Sleep(20 * time.Millisecond)
	paymentRepository.AssertExpectations(t)
	orderCreateManager.AssertExpectations(t)
	orderEventsHandler.AssertExpectations(t)
	userService.AssertExpectations(t)
	currencyService.AssertExpectations(t)
	priceGenerator.AssertExpectations(t)
}

func TestAutoExchangeManager_WithExceptionError(t *testing.T) {
	user := user.User{
		ID: 1,
	}
	pair := currency.Pair{
		ID:            1,
		Name:          "BTC-USDT",
		BasisCoin:     currency.Coin{Code: "USDT"},
		DependentCoin: currency.Coin{Code: "BTC"},
	}
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
	db, err := gorm.Open(dialector, &gorm.Config{SkipDefaultTransaction: true})
	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectExec("UPDATE crypto_payment_extra_info").WillReturnResult(sqlmock.NewResult(1, 1))
	paymentRepository := new(mocks.PaymentRepository)
	extraInfo := &payment.ExtraInfo{}
	paymentRepository.On("GetExtraInfoByPaymentID", int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		extraInfo = args.Get(1).(*payment.ExtraInfo)
		extraInfo.ID = 1
		extraInfo.PaymentID = 1
		extraInfo.IP = sql.NullString{String: "127.0.0.1", Valid: true}
	})
	orderCreateManager := new(mocks.OrderCreateManager)
	orderEventsHandler := new(mocks.EventsHandler)
	userService := new(mocks.UserService)
	userService.On("GetUserByID", 1).Once().Return(user, nil)
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetActivePairCurrenciesList").Once().Return([]currency.Pair{pair})
	priceGenerator := new(mocks.PriceGenerator)
	priceGenerator.On("GetPrice", mock.Anything, "BTC-USDT").Once().Return("", fmt.Errorf("price not found"))
	logger := new(mocks.Logger)
	logger.On("Error2", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return()

	autoExchangeManger := payment.NewAutoExchangeManger(db, paymentRepository, orderCreateManager, orderEventsHandler, userService,
		currencyService, priceGenerator, logger)
	p := &payment.Payment{
		ID:     1,
		Type:   payment.TypeDeposit,
		Status: payment.StatusCompleted,
		Amount: sql.NullString{String: "100.0", Valid: true},
		UserID: 1,
	}
	ub := &userbalance.UserBalance{
		BalanceCoin:      "USDT",
		AutoExchangeCoin: sql.NullString{String: "BTC", Valid: true},
	}
	autoExchangeManger.AutoExchange(p, ub)

	assert.Equal(t, false, extraInfo.AutoExchangeFailureReason.Valid)
	assert.Equal(t, "", extraInfo.AutoExchangeFailureReason.String)
	assert.Equal(t, true, extraInfo.AutoExchangeFailureType.Valid)
	assert.Equal(t, payment.FailureTypeException, extraInfo.AutoExchangeFailureType.String)
	paymentRepository.AssertExpectations(t)
	userService.AssertExpectations(t)
	currencyService.AssertExpectations(t)
	priceGenerator.AssertExpectations(t)

}

func TestAutoExchangeManager_WithLogicalError(t *testing.T) {
	user := user.User{
		ID: 1,
	}
	pair := currency.Pair{
		ID:            1,
		Name:          "BTC-USDT",
		BasisCoin:     currency.Coin{Code: "USDT"},
		DependentCoin: currency.Coin{Code: "BTC"},
	}
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
	db, err := gorm.Open(dialector, &gorm.Config{SkipDefaultTransaction: true})
	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectExec("UPDATE crypto_payment_extra_info").WillReturnResult(sqlmock.NewResult(1, 1))
	paymentRepository := new(mocks.PaymentRepository)
	extraInfo := &payment.ExtraInfo{}
	paymentRepository.On("GetExtraInfoByPaymentID", int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		extraInfo = args.Get(1).(*payment.ExtraInfo)
		extraInfo.ID = 1
		extraInfo.PaymentID = 1
		extraInfo.IP = sql.NullString{String: "127.0.0.1", Valid: true}
	})
	orderCreateManager := new(mocks.OrderCreateManager)
	requiredData := order.CreateRequiredData{
		User:         &user,
		Pair:         &pair,
		Amount:       "100.0",
		OrderType:    order.TypeBuy,
		ExchangeType: order.ExchangeTypeMarket,
		Price:        "",
		UserAgentInfo: order.UserAgentInfo{
			Device:  "",
			IP:      "127.0.0.1",
			Browser: "",
		},
		StopPointPrice: "",
		CurrentPrice:   "50000.0",
		IsInstant:      true,
		IsFastExchange: false,
		IsAutoExchange: true,
	}
	o := &order.Order{}
	orderCreateManager.On("CreateOrder", requiredData).Once().Return(o, platform.OrderCreateValidationError{Err: fmt.Errorf("user level does not allowed")})
	orderEventsHandler := new(mocks.EventsHandler)
	userService := new(mocks.UserService)
	userService.On("GetUserByID", 1).Once().Return(user, nil)
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetActivePairCurrenciesList").Once().Return([]currency.Pair{pair})
	priceGenerator := new(mocks.PriceGenerator)
	priceGenerator.On("GetPrice", mock.Anything, "BTC-USDT").Once().Return("50000.0", nil)
	logger := new(mocks.Logger)
	logger.On("Error2", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return()

	autoExchangeManger := payment.NewAutoExchangeManger(db, paymentRepository, orderCreateManager, orderEventsHandler, userService,
		currencyService, priceGenerator, logger)
	p := &payment.Payment{
		ID:     1,
		Type:   payment.TypeDeposit,
		Status: payment.StatusCompleted,
		Amount: sql.NullString{String: "100.0", Valid: true},
		UserID: 1,
	}
	ub := &userbalance.UserBalance{
		BalanceCoin:      "USDT",
		AutoExchangeCoin: sql.NullString{String: "BTC", Valid: true},
	}
	autoExchangeManger.AutoExchange(p, ub)
	assert.Equal(t, true, extraInfo.AutoExchangeFailureReason.Valid)
	assert.Equal(t, "user level does not allowed", extraInfo.AutoExchangeFailureReason.String)
	assert.Equal(t, true, extraInfo.AutoExchangeFailureType.Valid)
	assert.Equal(t, payment.FailureTypeLogical, extraInfo.AutoExchangeFailureType.String)

	paymentRepository.AssertExpectations(t)
	orderCreateManager.AssertExpectations(t)
	userService.AssertExpectations(t)
	currencyService.AssertExpectations(t)
	priceGenerator.AssertExpectations(t)
}
