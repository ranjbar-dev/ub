// Package order_test tests the order Service. Covers:
//   - CreateOrder: wrong pair ID, unverified user, missing permission, no price in Redis, creation error, successful creation
//   - CancelOrder: non-existent order, order not belonging to user, non-open status, action not allowed, successful cancellation, stop order cancellation
//   - GetOpenOrders: retrieval with pagination and pair filtering
//   - GetOrdersHistory: retrieval with pagination and pair filtering
//   - GetTradesHistory: retrieval with pagination and pair filtering
//   - FulfillOrder: non-existent order, non-open status, successful admin fulfillment
//   - GetOrderDetail: retrieval of a single order
//
// Test data: sqlmock MySQL DB, mocked order repository, create manager, events handler,
// currency service, price generator, user balance service, Redis manager, and configs.
package order_test

import (
	"database/sql"
	"exchange-go/internal/currency"
	"exchange-go/internal/mocks"
	"exchange-go/internal/order"
	"exchange-go/internal/platform"
	"exchange-go/internal/user"
	"exchange-go/internal/userbalance"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestService_CreateOrder_WrongPairId(t *testing.T) {
	db := &gorm.DB{}
	orderRepo := new(mocks.OrderRepository)
	orderCreateManager := new(mocks.OrderCreateManager)
	eh := new(mocks.EventsHandler)
	cs := new(mocks.CurrencyService)
	cs.On("GetPairByID", int64(0)).Once().Return(currency.Pair{}, gorm.ErrRecordNotFound)
	pg := new(mocks.PriceGenerator)
	userBalanceService := new(mocks.UserBalanceService)
	redisManager := new(mocks.OrderRedisManager)
	userConfigService := new(mocks.UserConfigService)
	userPermissionManager := new(mocks.UserPermissionManager)
	adminOrderManager := new(mocks.AdminOrderManager)
	engineCommunicator := new(mocks.EngineCommunicator)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	orderService := order.NewOrderService(db, orderRepo, orderCreateManager, eh, cs, pg, userBalanceService, redisManager,
		userConfigService, userPermissionManager, adminOrderManager, engineCommunicator, configs, logger)

	u := user.User{
		Status: user.StatusVerified,
	}

	params := order.CreateOrderParams{
		Type:           "buy",
		ExchangeType:   "market",
		Amount:         "20.0",
		PairID:         0,
		Price:          "",
		StopPointPrice: "",
	}

	res, statusCode := orderService.CreateOrder(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)

	assert.Equal(t, false, res.Status)
	assert.NotEqual(t, "", res.Message)
	cs.AssertExpectations(t)

}

func TestService_CreateOrder_NotVerifiedUser(t *testing.T) {
	db := &gorm.DB{}
	orderRepo := new(mocks.OrderRepository)
	orderCreateManager := new(mocks.OrderCreateManager)
	eh := new(mocks.EventsHandler)
	cs := new(mocks.CurrencyService)
	pg := new(mocks.PriceGenerator)
	userBalanceService := new(mocks.UserBalanceService)
	redisManager := new(mocks.OrderRedisManager)
	userConfigService := new(mocks.UserConfigService)
	userPermissionManager := new(mocks.UserPermissionManager)
	adminOrderManager := new(mocks.AdminOrderManager)
	engineCommunicator := new(mocks.EngineCommunicator)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	orderService := order.NewOrderService(db, orderRepo, orderCreateManager, eh, cs, pg, userBalanceService, redisManager,
		userConfigService, userPermissionManager, adminOrderManager, engineCommunicator, configs, logger)

	u := user.User{
		Status: user.StatusRegistered,
	}

	params := order.CreateOrderParams{
		Type:           "buy",
		ExchangeType:   "market",
		Amount:         "20.0",
		PairID:         0,
		Price:          "",
		StopPointPrice: "",
	}

	res, statusCode := orderService.CreateOrder(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.NotEqual(t, "", res.Message)
}

func TestService_CreateOrder_WithoutUserPermission(t *testing.T) {
	db := &gorm.DB{}
	orderRepo := new(mocks.OrderRepository)
	orderCreateManager := new(mocks.OrderCreateManager)
	eh := new(mocks.EventsHandler)
	cs := new(mocks.CurrencyService)
	cs.On("GetPairByID", int64(1)).Once().Return(currency.Pair{ID: 1}, nil)
	pg := new(mocks.PriceGenerator)
	userBalanceService := new(mocks.UserBalanceService)
	redisManager := new(mocks.OrderRedisManager)
	userConfigService := new(mocks.UserConfigService)
	userConfigService.On("GetUserConfig", mock.Anything).Return(user.Config{}, gorm.ErrRecordNotFound)
	userPermissionManager := new(mocks.UserPermissionManager)
	userPermissionManager.On("IsPermissionGrantedToUserFor", mock.Anything, mock.Anything).Once().Return(false)
	adminOrderManager := new(mocks.AdminOrderManager)
	engineCommunicator := new(mocks.EngineCommunicator)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	orderService := order.NewOrderService(db, orderRepo, orderCreateManager, eh, cs, pg, userBalanceService, redisManager,
		userConfigService, userPermissionManager, adminOrderManager, engineCommunicator, configs, logger)

	u := user.User{
		Status: user.StatusVerified,
	}

	params := order.CreateOrderParams{
		Type:           "buy",
		ExchangeType:   "market",
		Amount:         "20.0",
		PairID:         1,
		Price:          "",
		StopPointPrice: "",
	}

	res, statusCode := orderService.CreateOrder(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.NotEqual(t, "", res.Message)
	cs.AssertExpectations(t)
	userConfigService.AssertExpectations(t)
	userPermissionManager.AssertExpectations(t)
}

func TestService_CreateOrder_NoPriceInRedis(t *testing.T) {
	db := &gorm.DB{}
	orderRepo := new(mocks.OrderRepository)
	orderCreateManager := new(mocks.OrderCreateManager)
	eh := new(mocks.EventsHandler)
	cs := new(mocks.CurrencyService)
	cs.On("GetPairByID", int64(1)).Once().Return(currency.Pair{ID: 1}, nil)
	pg := new(mocks.PriceGenerator)
	pg.On("GetPrice", mock.Anything, mock.Anything).Once().Return("", fmt.Errorf("price does not exist"))
	userBalanceService := new(mocks.UserBalanceService)
	redisManager := new(mocks.OrderRedisManager)
	userConfigService := new(mocks.UserConfigService)
	userConfigService.On("GetUserConfig", mock.Anything).Return(user.Config{}, gorm.ErrRecordNotFound)
	userPermissionManager := new(mocks.UserPermissionManager)
	userPermissionManager.On("IsPermissionGrantedToUserFor", mock.Anything, mock.Anything).Once().Return(true)
	adminOrderManager := new(mocks.AdminOrderManager)
	engineCommunicator := new(mocks.EngineCommunicator)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	logger.On("Error2", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return()

	orderService := order.NewOrderService(db, orderRepo, orderCreateManager, eh, cs, pg, userBalanceService, redisManager,
		userConfigService, userPermissionManager, adminOrderManager, engineCommunicator, configs, logger)

	u := user.User{
		Status: user.StatusVerified,
	}

	params := order.CreateOrderParams{
		Type:           "buy",
		ExchangeType:   "market",
		Amount:         "20.0",
		PairID:         1,
		Price:          "",
		StopPointPrice: "",
	}

	res, statusCode := orderService.CreateOrder(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.NotEqual(t, "", res.Message)
	cs.AssertExpectations(t)
	userConfigService.AssertExpectations(t)
	pg.AssertExpectations(t)
	userPermissionManager.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestService_CreateOrder_ErrorInCreateOrder(t *testing.T) {
	db := &gorm.DB{}
	orderRepo := new(mocks.OrderRepository)
	orderCreateManager := new(mocks.OrderCreateManager)
	err := platform.OrderCreateValidationError{Err: fmt.Errorf("balance is not enough")}
	orderCreateManager.On("CreateOrder", mock.Anything).Once().Return(nil, err)
	eh := new(mocks.EventsHandler)
	cs := new(mocks.CurrencyService)
	cs.On("GetPairByID", int64(1)).Once().Return(currency.Pair{ID: 1}, nil)
	pg := new(mocks.PriceGenerator)
	pg.On("GetPrice", mock.Anything, mock.Anything).Once().Return("50000", nil)
	userBalanceService := new(mocks.UserBalanceService)
	redisManager := new(mocks.OrderRedisManager)
	userConfigService := new(mocks.UserConfigService)
	userConfigService.On("GetUserConfig", mock.Anything).Return(user.Config{}, gorm.ErrRecordNotFound)
	userPermissionManager := new(mocks.UserPermissionManager)
	userPermissionManager.On("IsPermissionGrantedToUserFor", mock.Anything, mock.Anything).Once().Return(true)
	adminOrderManager := new(mocks.AdminOrderManager)
	engineCommunicator := new(mocks.EngineCommunicator)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	logger.On("Error2", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return()

	orderService := order.NewOrderService(db, orderRepo, orderCreateManager, eh, cs, pg, userBalanceService, redisManager,
		userConfigService, userPermissionManager, adminOrderManager, engineCommunicator, configs, logger)

	u := user.User{
		Status: user.StatusVerified,
	}

	params := order.CreateOrderParams{
		Type:           "buy",
		ExchangeType:   "market",
		Amount:         "20.0",
		PairID:         1,
		Price:          "",
		StopPointPrice: "",
	}

	res, statusCode := orderService.CreateOrder(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.NotEqual(t, "", res.Message)
	cs.AssertExpectations(t)
	userConfigService.AssertExpectations(t)
	pg.AssertExpectations(t)
	userPermissionManager.AssertExpectations(t)
	orderCreateManager.AssertExpectations(t)
	logger.AssertExpectations(t)
}

func TestService_CreateOrder_successful(t *testing.T) {
	db := &gorm.DB{}
	orderRepo := new(mocks.OrderRepository)
	orderCreateManager := new(mocks.OrderCreateManager)
	o := &order.Order{}
	orderCreateManager.On("CreateOrder", mock.Anything).Once().Return(o, nil)
	eh := new(mocks.EventsHandler)
	eh.On("HandleOrderCreation", mock.Anything, mock.Anything).Once().Return()
	cs := new(mocks.CurrencyService)
	cs.On("GetPairByID", int64(1)).Once().Return(currency.Pair{ID: 1}, nil)
	pg := new(mocks.PriceGenerator)
	pg.On("GetPrice", mock.Anything, mock.Anything).Once().Return("50000", nil)
	userBalanceService := new(mocks.UserBalanceService)
	redisManager := new(mocks.OrderRedisManager)
	userConfigService := new(mocks.UserConfigService)
	userConfigService.On("GetUserConfig", mock.Anything).Return(user.Config{}, gorm.ErrRecordNotFound)
	userPermissionManager := new(mocks.UserPermissionManager)
	userPermissionManager.On("IsPermissionGrantedToUserFor", mock.Anything, mock.Anything).Once().Return(true)
	adminOrderManager := new(mocks.AdminOrderManager)
	engineCommunicator := new(mocks.EngineCommunicator)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	orderService := order.NewOrderService(db, orderRepo, orderCreateManager, eh, cs, pg, userBalanceService, redisManager,
		userConfigService, userPermissionManager, adminOrderManager, engineCommunicator, configs, logger)

	u := user.User{
		Status: user.StatusVerified,
	}

	params := order.CreateOrderParams{
		Type:           "buy",
		ExchangeType:   "market",
		Amount:         "20.0",
		PairID:         1,
		Price:          "",
		StopPointPrice: "",
	}

	res, statusCode := orderService.CreateOrder(&u, params)
	time.Sleep(20 * time.Millisecond)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)
	cs.AssertExpectations(t)
	userConfigService.AssertExpectations(t)
	pg.AssertExpectations(t)
	userPermissionManager.AssertExpectations(t)
	orderCreateManager.AssertExpectations(t)
	eh.AssertExpectations(t)
}

func TestService_CancelOrder_OrderDoesNotExist(t *testing.T) {
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

	orderRepo := new(mocks.OrderRepository)
	orderRepo.On("GetOrderByIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(gorm.ErrRecordNotFound)
	orderCreateManager := new(mocks.OrderCreateManager)

	eh := new(mocks.EventsHandler)
	cs := new(mocks.CurrencyService)
	pg := new(mocks.PriceGenerator)
	userBalanceService := new(mocks.UserBalanceService)
	redisManager := new(mocks.OrderRedisManager)
	userConfigService := new(mocks.UserConfigService)
	userPermissionManager := new(mocks.UserPermissionManager)
	adminOrderManager := new(mocks.AdminOrderManager)
	engineCommunicator := new(mocks.EngineCommunicator)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	orderService := order.NewOrderService(db, orderRepo, orderCreateManager, eh, cs, pg, userBalanceService, redisManager,
		userConfigService, userPermissionManager, adminOrderManager, engineCommunicator, configs, logger)

	u := user.User{
		Status: user.StatusVerified,
	}

	params := order.CancelOrderParams{
		ID: 1,
	}

	res, statusCode := orderService.CancelOrder(&u, params)

	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.NotEqual(t, "", res.Message)
	orderRepo.AssertExpectations(t)
}

func TestService_CancelOrder_OrderDoesNotBelongToUser(t *testing.T) {
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

	orderRepo := new(mocks.OrderRepository)
	orderRepo.On("GetOrderByIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		o := args.Get(2).(*order.Order)
		o.ID = 1
		o.UserID = 21
	})

	orderCreateManager := new(mocks.OrderCreateManager)

	eh := new(mocks.EventsHandler)
	cs := new(mocks.CurrencyService)
	pg := new(mocks.PriceGenerator)
	userBalanceService := new(mocks.UserBalanceService)
	redisManager := new(mocks.OrderRedisManager)
	userConfigService := new(mocks.UserConfigService)
	userPermissionManager := new(mocks.UserPermissionManager)
	adminOrderManager := new(mocks.AdminOrderManager)
	engineCommunicator := new(mocks.EngineCommunicator)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	orderService := order.NewOrderService(db, orderRepo, orderCreateManager, eh, cs, pg, userBalanceService, redisManager,
		userConfigService, userPermissionManager, adminOrderManager, engineCommunicator, configs, logger)

	u := user.User{
		ID: 22,
	}

	params := order.CancelOrderParams{
		ID: 1,
	}

	res, statusCode := orderService.CancelOrder(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.NotEqual(t, "", res.Message)
	orderRepo.AssertExpectations(t)
}

func TestService_CancelOrder_OrderStatusIsNotOpen(t *testing.T) {
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

	orderRepo := new(mocks.OrderRepository)
	//o := &order.Order{}
	orderRepo.On("GetOrderByIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		o := args.Get(2).(*order.Order)
		o.ID = 1
		o.UserID = 21
		o.Status = "FILLED"
	})

	orderCreateManager := new(mocks.OrderCreateManager)

	eh := new(mocks.EventsHandler)
	cs := new(mocks.CurrencyService)
	pg := new(mocks.PriceGenerator)
	userBalanceService := new(mocks.UserBalanceService)
	redisManager := new(mocks.OrderRedisManager)
	userConfigService := new(mocks.UserConfigService)
	userPermissionManager := new(mocks.UserPermissionManager)
	adminOrderManager := new(mocks.AdminOrderManager)
	engineCommunicator := new(mocks.EngineCommunicator)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	orderService := order.NewOrderService(db, orderRepo, orderCreateManager, eh, cs, pg, userBalanceService, redisManager,
		userConfigService, userPermissionManager, adminOrderManager, engineCommunicator, configs, logger)

	u := user.User{
		ID: 21,
	}

	params := order.CancelOrderParams{
		ID: 1,
	}

	res, statusCode := orderService.CancelOrder(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.NotEqual(t, "", res.Message)
	orderRepo.AssertExpectations(t)
}

func TestService_CancelOrder_ActionNotAllowed(t *testing.T) {
	oldIsActionAllowed := currency.IsActionAllowed
	defer func() { order.IsActionAllowed = oldIsActionAllowed }()
	order.IsActionAllowed = func(pair currency.Pair, action string) bool {
		return false
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
	db, err := gorm.Open(dialector, &gorm.Config{})
	dbMock.MatchExpectationsInOrder(false)
	dbMock.ExpectBegin()
	dbMock.ExpectRollback()

	orderRepo := new(mocks.OrderRepository)
	//o := &order.Order{}
	orderRepo.On("GetOrderByIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		o := args.Get(2).(*order.Order)
		o.ID = 1
		o.UserID = 21
		o.Status = "OPEN"
		o.PairID = 1
	})

	orderCreateManager := new(mocks.OrderCreateManager)

	eh := new(mocks.EventsHandler)
	cs := new(mocks.CurrencyService)
	cs.On("GetPairByID", int64(1)).Once().Return(currency.Pair{}, nil)
	pg := new(mocks.PriceGenerator)
	userBalanceService := new(mocks.UserBalanceService)
	redisManager := new(mocks.OrderRedisManager)
	userConfigService := new(mocks.UserConfigService)
	userPermissionManager := new(mocks.UserPermissionManager)
	adminOrderManager := new(mocks.AdminOrderManager)
	engineCommunicator := new(mocks.EngineCommunicator)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	orderService := order.NewOrderService(db, orderRepo, orderCreateManager, eh, cs, pg, userBalanceService, redisManager,
		userConfigService, userPermissionManager, adminOrderManager, engineCommunicator, configs, logger)

	u := user.User{
		ID: 21,
	}

	params := order.CancelOrderParams{
		ID: 1,
	}

	res, statusCode := orderService.CancelOrder(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.NotEqual(t, "", res.Message)
	orderRepo.AssertExpectations(t)
}

type queryMatcher struct {
}

func (queryMatcher) Match(expectedSQL, actualSQL string) error {
	return nil
}

func TestService_CancelOrder_Successful(t *testing.T) {
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
	dbMock.ExpectExec("UPDATE orders").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit()

	orderRepo := new(mocks.OrderRepository)
	//o := &order.Order{}
	orderRepo.On("GetOrderByIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		o := args.Get(2).(*order.Order)
		o.ID = 1
		o.UserID = 21
		o.Status = "OPEN"
		o.PairID = 1
	})

	orderCreateManager := new(mocks.OrderCreateManager)

	eh := new(mocks.EventsHandler)
	eh.On("HandleOrderCancellation", mock.Anything).Once().Return()
	cs := new(mocks.CurrencyService)
	cs.On("GetPairByID", int64(1)).Once().Return(currency.Pair{}, nil)
	pg := new(mocks.PriceGenerator)
	userBalanceService := new(mocks.UserBalanceService)
	userBalanceService.On("GetBalanceOfUserByCoinUsingTx", mock.Anything, 21, mock.Anything, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub := args.Get(3).(*userbalance.UserBalance)
		ub.ID = 1
		ub.UserID = 21
	})

	redisManager := new(mocks.OrderRedisManager)
	userConfigService := new(mocks.UserConfigService)
	userPermissionManager := new(mocks.UserPermissionManager)
	adminOrderManager := new(mocks.AdminOrderManager)
	engineCommunicator := new(mocks.EngineCommunicator)
	engineCommunicator.On("RemoveOrder", mock.Anything).Once().Return(nil)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	orderService := order.NewOrderService(db, orderRepo, orderCreateManager, eh, cs, pg, userBalanceService, redisManager,
		userConfigService, userPermissionManager, adminOrderManager, engineCommunicator, configs, logger)

	u := user.User{
		ID: 21,
	}

	params := order.CancelOrderParams{
		ID: 1,
	}

	res, statusCode := orderService.CancelOrder(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)
	orderRepo.AssertExpectations(t)
	userBalanceService.AssertExpectations(t)
	eh.AssertExpectations(t)
}

func TestService_CancelOrder_stopOrderSuccessful(t *testing.T) {
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
	dbMock.ExpectExec("UPDATE orders").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectExec("UPDATE user_balances").WillReturnResult(sqlmock.NewResult(1, 1))
	dbMock.ExpectCommit()

	orderRepo := new(mocks.OrderRepository)
	//o := &order.Order{}
	orderRepo.On("GetOrderByIDUsingTx", mock.Anything, int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		o := args.Get(2).(*order.Order)
		o.ID = 1
		o.UserID = 21
		o.Status = "OPEN"
		o.StopPointPrice = sql.NullString{String: "54000", Valid: true}
		o.PairID = 1
	})

	orderCreateManager := new(mocks.OrderCreateManager)

	eh := new(mocks.EventsHandler)
	eh.On("HandleOrderCancellation", mock.Anything).Once().Return()

	cs := new(mocks.CurrencyService)
	cs.On("GetPairByID", int64(1)).Once().Return(currency.Pair{}, nil)
	pg := new(mocks.PriceGenerator)
	userBalanceService := new(mocks.UserBalanceService)
	userBalanceService.On("GetBalanceOfUserByCoinUsingTx", mock.Anything, 21, mock.Anything, mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		ub := args.Get(3).(*userbalance.UserBalance)
		ub.ID = 1
		ub.UserID = 21
	})

	redisManager := new(mocks.OrderRedisManager)
	redisManager.On("RemoveStopOrderFromQueue", mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)
	userConfigService := new(mocks.UserConfigService)
	userPermissionManager := new(mocks.UserPermissionManager)
	adminOrderManager := new(mocks.AdminOrderManager)
	engineCommunicator := new(mocks.EngineCommunicator)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	orderService := order.NewOrderService(db, orderRepo, orderCreateManager, eh, cs, pg, userBalanceService, redisManager,
		userConfigService, userPermissionManager, adminOrderManager, engineCommunicator, configs, logger)

	u := user.User{
		ID: 21,
	}

	params := order.CancelOrderParams{
		ID: 1,
	}

	res, statusCode := orderService.CancelOrder(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)
	orderRepo.AssertExpectations(t)
	userBalanceService.AssertExpectations(t)
	redisManager.AssertExpectations(t)
}

func TestService_GetOpenOrders(t *testing.T) {
	db := &gorm.DB{}
	orderRepo := new(mocks.OrderRepository)
	btcUsdtPair := currency.Pair{
		ID:   1,
		Name: "BTC-USDT",
	}
	orders := []order.Order{
		{
			ID:             1,
			Type:           "SELL",
			ExchangeType:   "LIMIT",
			Pair:           btcUsdtPair,
			PayedByAmount:  sql.NullString{String: "0.1", Valid: true},
			DemandedAmount: sql.NullString{String: "5000", Valid: true},
			Price:          sql.NullString{String: "50000", Valid: true},
			Path:           sql.NullString{String: "1,", Valid: true},
			CreatedAt:      time.Now().Add(5 * time.Second),
		},
		{
			ID:             2,
			Type:           "BUY",
			ExchangeType:   "LIMIT",
			Pair:           btcUsdtPair,
			Price:          sql.NullString{String: "52000", Valid: true},
			PayedByAmount:  sql.NullString{String: "5200", Valid: true},
			DemandedAmount: sql.NullString{String: "0.1", Valid: true},
			StopPointPrice: sql.NullString{String: "50000", Valid: true},
			Path:           sql.NullString{String: "2,", Valid: true},
			CreatedAt:      time.Now().Add(3 * time.Second),
		},
		{
			ID:           4,
			Type:         "BUY",
			ExchangeType: "LIMIT",
			Status:       "OPEN",

			Pair:           btcUsdtPair,
			PayedByAmount:  sql.NullString{String: "2500", Valid: true},
			DemandedAmount: sql.NullString{String: "0.05", Valid: true},
			Price:          sql.NullString{String: "50000", Valid: true},
			Path:           sql.NullString{String: "3,4,", Valid: true},
			CreatedAt:      time.Now(),
		},
	}

	orderRepo.On("GetUserOpenOrders", 21, int64(1)).Once().Return(orders)

	allAncestorOrders := []order.Order{
		{
			ID:                  3,
			Type:                "BUY",
			ExchangeType:        "LIMIT",
			Status:              "FILLED",
			Pair:                btcUsdtPair,
			PayedByAmount:       sql.NullString{String: "5000", Valid: true},
			DemandedAmount:      sql.NullString{String: "0.1", Valid: true},
			FinalPayedByAmount:  sql.NullString{String: "2500", Valid: true},
			FinalDemandedAmount: sql.NullString{String: "0.05", Valid: true},
			Price:               sql.NullString{String: "50000", Valid: true},
			Path:                sql.NullString{String: "3,", Valid: true},
			CreatedAt:           time.Now(),
		},
	}
	orderRepo.On("GetOrdersAncestors", mock.Anything).Once().Return(allAncestorOrders)

	orderCreateManager := new(mocks.OrderCreateManager)
	eh := new(mocks.EventsHandler)
	cs := new(mocks.CurrencyService)
	pg := new(mocks.PriceGenerator)
	userBalanceService := new(mocks.UserBalanceService)
	redisManager := new(mocks.OrderRedisManager)
	userConfigService := new(mocks.UserConfigService)
	userPermissionManager := new(mocks.UserPermissionManager)
	adminOrderManager := new(mocks.AdminOrderManager)
	engineCommunicator := new(mocks.EngineCommunicator)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	orderService := order.NewOrderService(db, orderRepo, orderCreateManager, eh, cs, pg, userBalanceService, redisManager,
		userConfigService, userPermissionManager, adminOrderManager, engineCommunicator, configs, logger)

	u := user.User{
		ID: 21,
	}

	params := order.GetOpenOrdersParams{
		PairID: 1,
	}

	res, statusCode := orderService.GetOpenOrders(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)
	openOrders, ok := res.Data.([]order.OpenOrdersResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}
	assert.Equal(t, 3, len(openOrders))

	order1 := openOrders[0]
	assert.Equal(t, int64(1), order1.ID)
	assert.Equal(t, "sell", order1.Side)
	assert.Equal(t, "BTC-USDT", order1.Pair)
	assert.Equal(t, "limit", order1.OrderType)
	assert.Equal(t, "order", order1.MainType)
	assert.Equal(t, "", order1.TriggerCondition)
	assert.Equal(t, "50000", order1.Price)
	assert.Equal(t, "0.1", order1.Amount)
	assert.Equal(t, "5000.00000000", order1.Total)
	assert.Equal(t, "0.00 %", order1.Executed)

	order2 := openOrders[1]
	assert.Equal(t, int64(2), order2.ID)
	assert.Equal(t, "buy", order2.Side)
	assert.Equal(t, "BTC-USDT", order2.Pair)
	assert.Equal(t, "limit", order2.OrderType)
	assert.Equal(t, "stopOrder", order2.MainType)
	assert.Equal(t, ">= 50000", order2.TriggerCondition)
	assert.Equal(t, "52000", order2.Price)
	assert.Equal(t, "0.1", order2.Amount)
	assert.Equal(t, "5200.00000000", order2.Total)
	assert.Equal(t, "0.00 %", order2.Executed)

	order3 := openOrders[2]

	assert.Equal(t, int64(4), order3.ID)
	assert.Equal(t, "buy", order3.Side)
	assert.Equal(t, "BTC-USDT", order3.Pair)
	assert.Equal(t, "limit", order3.OrderType)
	assert.Equal(t, "order", order3.MainType)
	assert.Equal(t, "", order3.TriggerCondition)
	assert.Equal(t, "50000", order3.Price)
	assert.Equal(t, "0.1", order3.Amount)
	assert.Equal(t, "5000.00000000", order3.Total)
	assert.Equal(t, "50.00 %", order3.Executed)
	orderRepo.AssertExpectations(t)
}

func TestService_GetOrdersHistory(t *testing.T) {
	db := &gorm.DB{}
	orderRepo := new(mocks.OrderRepository)
	btcUsdtPair := currency.Pair{
		ID:   1,
		Name: "BTC-USDT",
	}

	fields := []order.HistoryNeededField{
		{
			OrderID:   1,
			Path:      "1,",
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		},
		{
			OrderID:   2,
			Path:      "2,",
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		},
		{
			OrderID:   4,
			Path:      "3,4,",
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		},
		{
			OrderID:   6,
			Path:      "5,6,",
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		},
		{
			OrderID:   7,
			Path:      "7,",
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		},
		{
			OrderID:   8,
			Path:      "8,",
			CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		},
	}

	orderRepo.On("GetLeafOrders", mock.Anything).Once().Return(fields)
	allOrders := []order.Order{
		{
			ID:                  1,
			Type:                "SELL",
			ExchangeType:        "LIMIT",
			Status:              "FILLED",
			Pair:                btcUsdtPair,
			PayedByAmount:       sql.NullString{String: "0.1", Valid: true},
			DemandedAmount:      sql.NullString{String: "5000", Valid: true},
			FinalPayedByAmount:  sql.NullString{String: "0.1", Valid: true},
			FinalDemandedAmount: sql.NullString{String: "5000", Valid: true},
			Price:               sql.NullString{String: "50000", Valid: true},
			TradePrice:          sql.NullString{String: "50000", Valid: true},
			Path:                sql.NullString{String: "1,", Valid: true},
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		},
		{
			ID:                  2,
			Type:                "BUY",
			ExchangeType:        "LIMIT",
			Status:              "FILLED",
			Pair:                btcUsdtPair,
			PayedByAmount:       sql.NullString{String: "5200", Valid: true},
			DemandedAmount:      sql.NullString{String: "0.1", Valid: true},
			FinalPayedByAmount:  sql.NullString{String: "5200", Valid: true},
			FinalDemandedAmount: sql.NullString{String: "0.1", Valid: true},
			StopPointPrice:      sql.NullString{String: "50000", Valid: true},
			Price:               sql.NullString{String: "52000", Valid: true},
			TradePrice:          sql.NullString{String: "52000", Valid: true},
			Path:                sql.NullString{String: "2,", Valid: true},
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		},
		{
			ID:                  3,
			Type:                "BUY",
			ExchangeType:        "LIMIT",
			Status:              "FILLED",
			Pair:                btcUsdtPair,
			PayedByAmount:       sql.NullString{String: "5000", Valid: true},
			DemandedAmount:      sql.NullString{String: "0.1", Valid: true},
			FinalPayedByAmount:  sql.NullString{String: "2500", Valid: true},
			FinalDemandedAmount: sql.NullString{String: "0.05", Valid: true},
			Price:               sql.NullString{String: "50000", Valid: true},
			TradePrice:          sql.NullString{String: "50000", Valid: true},
			Path:                sql.NullString{String: "3,", Valid: true},
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		},
		{
			ID:             4,
			Type:           "BUY",
			ExchangeType:   "LIMIT",
			Status:         "CANCELED",
			Pair:           btcUsdtPair,
			PayedByAmount:  sql.NullString{String: "2500", Valid: true},
			DemandedAmount: sql.NullString{String: "0.05", Valid: true},
			//FinalPayedByAmount:  sql.NullString{String: "2500", Valid: true},
			//FinalDemandedAmount: sql.NullString{String: "0.05", Valid: true},
			Price:     sql.NullString{String: "50000", Valid: true},
			Path:      sql.NullString{String: "3,4", Valid: true},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},

		{
			ID:                  5,
			Type:                "SELL",
			ExchangeType:        "LIMIT",
			Status:              "FILLED",
			Pair:                btcUsdtPair,
			PayedByAmount:       sql.NullString{String: "0.1", Valid: true},
			DemandedAmount:      sql.NullString{String: "5000", Valid: true},
			FinalPayedByAmount:  sql.NullString{String: "0.05", Valid: true},
			FinalDemandedAmount: sql.NullString{String: "2500", Valid: true},
			Price:               sql.NullString{String: "50000", Valid: true},
			TradePrice:          sql.NullString{String: "50000", Valid: true},
			Path:                sql.NullString{String: "5,", Valid: true},
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		},
		{
			ID:                  6,
			Type:                "SELL",
			ExchangeType:        "LIMIT",
			Status:              "FILLED",
			Pair:                btcUsdtPair,
			PayedByAmount:       sql.NullString{String: "0.05", Valid: true},
			DemandedAmount:      sql.NullString{String: "2500", Valid: true},
			FinalPayedByAmount:  sql.NullString{String: "0.05", Valid: true},
			FinalDemandedAmount: sql.NullString{String: "2500", Valid: true},
			Price:               sql.NullString{String: "50000", Valid: true},
			TradePrice:          sql.NullString{String: "50000", Valid: true},
			Path:                sql.NullString{String: "5,6,", Valid: true},
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		},
		{
			ID:                  7,
			Type:                "SELL",
			ExchangeType:        "MARKET",
			Status:              "FILLED",
			Pair:                btcUsdtPair,
			PayedByAmount:       sql.NullString{String: "0.1", Valid: true},
			FinalPayedByAmount:  sql.NullString{String: "0.1", Valid: true},
			FinalDemandedAmount: sql.NullString{String: "5000", Valid: true},
			TradePrice:          sql.NullString{String: "50000", Valid: true},
			Path:                sql.NullString{String: "7,", Valid: true},
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
			IsFastExchange:      true,
		},
		{
			ID:                  8,
			Type:                "BUY",
			ExchangeType:        "MARKET",
			Status:              "FILLED",
			Pair:                btcUsdtPair,
			PayedByAmount:       sql.NullString{String: "5000", Valid: true},
			DemandedAmount:      sql.NullString{String: "0.1", Valid: true},
			FinalPayedByAmount:  sql.NullString{String: "5000", Valid: true},
			FinalDemandedAmount: sql.NullString{String: "0.1", Valid: true},
			TradePrice:          sql.NullString{String: "50000", Valid: true},
			Path:                sql.NullString{String: "8,", Valid: true},
			CreatedAt:           time.Now(),
			UpdatedAt:           time.Now(),
		},
	}
	orderRepo.On("GetOrdersByIds", mock.Anything).Once().Return(allOrders)

	orderCreateManager := new(mocks.OrderCreateManager)
	eh := new(mocks.EventsHandler)
	cs := new(mocks.CurrencyService)
	pg := new(mocks.PriceGenerator)
	userBalanceService := new(mocks.UserBalanceService)
	redisManager := new(mocks.OrderRedisManager)
	userConfigService := new(mocks.UserConfigService)
	userPermissionManager := new(mocks.UserPermissionManager)
	adminOrderManager := new(mocks.AdminOrderManager)
	engineCommunicator := new(mocks.EngineCommunicator)
	configs := new(mocks.Configs)
	configs.On("GetEnv").Once().Return(platform.EnvTest)
	logger := new(mocks.Logger)

	orderService := order.NewOrderService(db, orderRepo, orderCreateManager, eh, cs, pg, userBalanceService, redisManager,
		userConfigService, userPermissionManager, adminOrderManager, engineCommunicator, configs, logger)

	u := user.User{
		ID: 21,
	}

	params := order.GetOrdersHistoryParams{
		PairID: 1,
	}

	res, statusCode := orderService.GetOrdersHistory(&u, params, true)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)
	ordersHistory, ok := res.Data.([]order.HistoryResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}
	assert.Equal(t, 6, len(ordersHistory))

	order1 := ordersHistory[0]
	assert.Equal(t, int64(1), order1.ID)
	assert.Equal(t, "sell", order1.Side)
	assert.Equal(t, "BTC-USDT", order1.Pair)
	assert.Equal(t, "limit", order1.OrderType)
	assert.Equal(t, "order", order1.MainType)
	assert.Equal(t, "", order1.TriggerCondition)
	assert.Equal(t, "50000", order1.Price)
	assert.Equal(t, "50000.00000000", order1.AveragePrice)
	assert.Equal(t, "0.1 BTC", order1.Amount)
	assert.Equal(t, "5000 USDT", order1.Total)
	assert.Equal(t, "100.00 %", order1.Executed)
	assert.Equal(t, "filled", order1.Status)
	assert.NotEmpty(t, order1.CreatedAt)
	assert.NotEmpty(t, order1.UpdatedAt)

	order2 := ordersHistory[1]
	assert.Equal(t, int64(2), order2.ID)
	assert.Equal(t, "buy", order2.Side)
	assert.Equal(t, "BTC-USDT", order2.Pair)
	assert.Equal(t, "limit", order2.OrderType)
	assert.Equal(t, "stopOrder", order2.MainType)
	assert.Equal(t, ">= 50000", order2.TriggerCondition)
	assert.Equal(t, "52000", order2.Price)
	assert.Equal(t, "52000.00000000", order2.AveragePrice)
	assert.Equal(t, "5200 USDT", order2.Amount)
	assert.Equal(t, "0.1 BTC", order2.Total)
	assert.Equal(t, "100.00 %", order2.Executed)
	assert.Equal(t, "filled", order2.Status)
	assert.NotEmpty(t, order2.CreatedAt)
	assert.NotEmpty(t, order2.UpdatedAt)

	order3 := ordersHistory[2]
	assert.Equal(t, int64(3), order3.ID)
	assert.Equal(t, "buy", order3.Side)
	assert.Equal(t, "BTC-USDT", order3.Pair)
	assert.Equal(t, "limit", order3.OrderType)
	assert.Equal(t, "order", order3.MainType)
	assert.Equal(t, "", order3.TriggerCondition)
	assert.Equal(t, "50000", order3.Price)
	assert.Equal(t, "50000.00000000", order3.AveragePrice)
	assert.Equal(t, "5000 USDT", order3.Amount)
	assert.Equal(t, "0.05 BTC", order3.Total)
	assert.Equal(t, "50.00 %", order3.Executed)
	assert.Equal(t, "canceled", order3.Status)
	assert.NotEmpty(t, order3.CreatedAt)
	assert.NotEmpty(t, order3.UpdatedAt)

	order4 := ordersHistory[3]
	assert.Equal(t, int64(5), order4.ID)
	assert.Equal(t, "sell", order4.Side)
	assert.Equal(t, "BTC-USDT", order4.Pair)
	assert.Equal(t, "limit", order4.OrderType)
	assert.Equal(t, "order", order4.MainType)
	assert.Equal(t, "", order4.TriggerCondition)
	assert.Equal(t, "50000", order4.Price)
	assert.Equal(t, "50000.00000000", order4.AveragePrice)
	assert.Equal(t, "0.1 BTC", order4.Amount)
	assert.Equal(t, "5000 USDT", order4.Total)
	assert.Equal(t, "100.00 %", order4.Executed)
	assert.Equal(t, "filled", order4.Status)
	assert.NotEmpty(t, order4.CreatedAt)
	assert.NotEmpty(t, order4.UpdatedAt)

	order5 := ordersHistory[4]
	assert.Equal(t, int64(7), order5.ID)
	assert.Equal(t, "sell", order5.Side)
	assert.Equal(t, "BTC-USDT", order5.Pair)
	assert.Equal(t, "fast exchange", order5.OrderType)
	assert.Equal(t, "order", order5.MainType)
	assert.Equal(t, "", order5.TriggerCondition)
	assert.Equal(t, "", order5.Price)
	assert.Equal(t, "50000.00000000", order5.AveragePrice)
	assert.Equal(t, "0.1 BTC", order5.Amount)
	assert.Equal(t, "5000 USDT", order5.Total)
	assert.Equal(t, "100.00 %", order5.Executed)
	assert.Equal(t, "filled", order5.Status)
	assert.Equal(t, "filled", order5.Status)
	assert.NotEmpty(t, order5.CreatedAt)
	assert.NotEmpty(t, order5.UpdatedAt)

	order6 := ordersHistory[5]
	assert.Equal(t, int64(8), order6.ID)
	assert.Equal(t, "buy", order6.Side)
	assert.Equal(t, "BTC-USDT", order6.Pair)
	assert.Equal(t, "market", order6.OrderType)
	assert.Equal(t, "order", order6.MainType)
	assert.Equal(t, "", order6.TriggerCondition)
	assert.Equal(t, "", order6.Price)
	assert.Equal(t, "50000.00000000", order6.AveragePrice)
	assert.Equal(t, "5000 USDT", order6.Amount)
	assert.Equal(t, "0.1 BTC", order6.Total)
	assert.Equal(t, "100.00 %", order6.Executed)
	assert.Equal(t, "filled", order6.Status)
	assert.NotEmpty(t, order6.CreatedAt)
	assert.NotEmpty(t, order6.UpdatedAt)

	orderRepo.AssertExpectations(t)
}

func TestService_GetTradesHistory(t *testing.T) {

	db := &gorm.DB{}
	orderRepo := new(mocks.OrderRepository)
	btcUsdtPair := currency.Pair{
		ID:   1,
		Name: "BTC-USDT",
	}

	orders := []order.Order{
		{
			ID:                  1,
			Type:                "SELL",
			ExchangeType:        "LIMIT",
			Status:              "FILLED",
			Pair:                btcUsdtPair,
			PayedByAmount:       sql.NullString{String: "0.1", Valid: true},
			DemandedAmount:      sql.NullString{String: "5000", Valid: true},
			FinalPayedByAmount:  sql.NullString{String: "0.1", Valid: true},
			FinalDemandedAmount: sql.NullString{String: "5000", Valid: true},
			Price:               sql.NullString{String: "50000", Valid: true},
			TradePrice:          sql.NullString{String: "50000", Valid: true},
			FeePercentage:       sql.NullFloat64{Float64: 0.1, Valid: true},
			Path:                sql.NullString{String: "1,", Valid: true},
			CreatedAt:           time.Now(),
		},
		{
			ID:                  2,
			Type:                "BUY",
			ExchangeType:        "LIMIT",
			Status:              "FILLED",
			Pair:                btcUsdtPair,
			PayedByAmount:       sql.NullString{String: "5000", Valid: true},
			DemandedAmount:      sql.NullString{String: "0.1", Valid: true},
			FinalPayedByAmount:  sql.NullString{String: "5000", Valid: true},
			FinalDemandedAmount: sql.NullString{String: "0.1", Valid: true},
			Price:               sql.NullString{String: "50000", Valid: true},
			TradePrice:          sql.NullString{String: "50000", Valid: true},
			FeePercentage:       sql.NullFloat64{Float64: 0.1, Valid: true},
			Path:                sql.NullString{String: "2,", Valid: true},
			CreatedAt:           time.Now(),
		},
		{
			ID:                  3,
			Type:                "SELL",
			ExchangeType:        "MARKET",
			Status:              "FILLED",
			Pair:                btcUsdtPair,
			PayedByAmount:       sql.NullString{String: "0.1", Valid: true},
			DemandedAmount:      sql.NullString{String: "5000", Valid: true},
			FinalPayedByAmount:  sql.NullString{String: "0.1", Valid: true},
			FinalDemandedAmount: sql.NullString{String: "5000", Valid: true},
			TradePrice:          sql.NullString{String: "50000", Valid: true},
			FeePercentage:       sql.NullFloat64{Float64: 0.1, Valid: true},
			Path:                sql.NullString{String: "1,", Valid: true},
			CreatedAt:           time.Now(),
		},
		{
			ID:                  4,
			Type:                "BUY",
			ExchangeType:        "MARKET",
			Status:              "FILLED",
			Pair:                btcUsdtPair,
			PayedByAmount:       sql.NullString{String: "5000", Valid: true},
			DemandedAmount:      sql.NullString{String: "0.1", Valid: true},
			FinalPayedByAmount:  sql.NullString{String: "5000", Valid: true},
			FinalDemandedAmount: sql.NullString{String: "0.1", Valid: true},
			TradePrice:          sql.NullString{String: "50000", Valid: true},
			FeePercentage:       sql.NullFloat64{Float64: 0.1, Valid: true},
			Path:                sql.NullString{String: "2,", Valid: true},
			CreatedAt:           time.Now(),
		},
		{
			ID:                  5,
			Type:                "BUY",
			ExchangeType:        "LIMIT",
			Status:              "FILLED",
			Pair:                btcUsdtPair,
			PayedByAmount:       sql.NullString{String: "5000", Valid: true},
			DemandedAmount:      sql.NullString{String: "0.1", Valid: true},
			FinalPayedByAmount:  sql.NullString{String: "2500", Valid: true},
			FinalDemandedAmount: sql.NullString{String: "0.05", Valid: true},
			TradePrice:          sql.NullString{String: "50000", Valid: true},
			Price:               sql.NullString{String: "50000", Valid: true},
			FeePercentage:       sql.NullFloat64{Float64: 0.1, Valid: true},
			Path:                sql.NullString{String: "2,", Valid: true},
			CreatedAt:           time.Now(),
		},
	}
	orderRepo.On("GetUserTradedOrders", mock.Anything).Once().Return(orders)

	orderCreateManager := new(mocks.OrderCreateManager)
	eh := new(mocks.EventsHandler)
	cs := new(mocks.CurrencyService)
	pg := new(mocks.PriceGenerator)
	userBalanceService := new(mocks.UserBalanceService)
	redisManager := new(mocks.OrderRedisManager)
	userConfigService := new(mocks.UserConfigService)
	userPermissionManager := new(mocks.UserPermissionManager)
	adminOrderManager := new(mocks.AdminOrderManager)
	engineCommunicator := new(mocks.EngineCommunicator)
	configs := new(mocks.Configs)
	configs.On("GetEnv").Once().Return(platform.EnvTest)
	logger := new(mocks.Logger)

	orderService := order.NewOrderService(db, orderRepo, orderCreateManager, eh, cs, pg, userBalanceService, redisManager,
		userConfigService, userPermissionManager, adminOrderManager, engineCommunicator, configs, logger)

	u := user.User{
		ID: 21,
	}

	params := order.GetTradesHistoryParams{
		PairID: 1,
	}

	res, statusCode := orderService.GetTradesHistory(&u, params, true)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)
	tradesHistory, ok := res.Data.([]order.TradeHistoryResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}
	assert.Equal(t, 5, len(tradesHistory))

	order1 := tradesHistory[0]
	assert.Equal(t, int64(1), order1.ID)
	assert.Equal(t, "BTC-USDT", order1.Pair)
	assert.Equal(t, "sell", order1.OrderType)
	assert.Equal(t, "50000", order1.Price)
	assert.Equal(t, "5000 USDT", order1.Amount)
	assert.Equal(t, "5000 USDT", order1.Total)
	assert.Equal(t, "500 USDT", order1.Fee)
	assert.Equal(t, "0.1 BTC", order1.Executed)

	order2 := tradesHistory[1]
	assert.Equal(t, int64(2), order2.ID)
	assert.Equal(t, "BTC-USDT", order2.Pair)
	assert.Equal(t, "buy", order2.OrderType)
	assert.Equal(t, "50000", order2.Price)
	assert.Equal(t, "5000 USDT", order2.Amount)
	assert.Equal(t, "5000 USDT", order2.Total)
	assert.Equal(t, "0.01 BTC", order2.Fee)
	assert.Equal(t, "0.1 BTC", order2.Executed)

	order3 := tradesHistory[2]
	assert.Equal(t, int64(3), order3.ID)
	assert.Equal(t, "BTC-USDT", order3.Pair)
	assert.Equal(t, "sell", order3.OrderType)
	assert.Equal(t, "50000", order3.Price)
	assert.Equal(t, "5000 USDT", order3.Amount)
	assert.Equal(t, "5000 USDT", order3.Total)
	assert.Equal(t, "500 USDT", order3.Fee)
	assert.Equal(t, "0.1 BTC", order3.Executed)

	order4 := tradesHistory[3]
	assert.Equal(t, int64(4), order4.ID)
	assert.Equal(t, "BTC-USDT", order4.Pair)
	assert.Equal(t, "buy", order4.OrderType)
	assert.Equal(t, "50000", order4.Price)
	assert.Equal(t, "5000 USDT", order4.Amount)
	assert.Equal(t, "5000 USDT", order4.Total)
	assert.Equal(t, "0.01 BTC", order4.Fee)
	assert.Equal(t, "0.1 BTC", order4.Executed)

	order5 := tradesHistory[4]
	assert.Equal(t, int64(5), order5.ID)
	assert.Equal(t, "BTC-USDT", order5.Pair)
	assert.Equal(t, "buy", order5.OrderType)
	assert.Equal(t, "50000", order5.Price)
	assert.Equal(t, "2500 USDT", order5.Amount)
	assert.Equal(t, "2500 USDT", order5.Total)
	assert.Equal(t, "0.005 BTC", order5.Fee)
	assert.Equal(t, "0.05 BTC", order5.Executed)

	orderRepo.AssertExpectations(t)

}

func TestService_FulfillOrder_OrderDoesNotExist(t *testing.T) {
	db := &gorm.DB{}
	orderRepo := new(mocks.OrderRepository)
	orderRepo.On("GetOrderByID", int64(1), mock.Anything).Once().Return(gorm.ErrRecordNotFound)
	orderCreateManager := new(mocks.OrderCreateManager)

	eh := new(mocks.EventsHandler)
	cs := new(mocks.CurrencyService)
	pg := new(mocks.PriceGenerator)
	userBalanceService := new(mocks.UserBalanceService)
	redisManager := new(mocks.OrderRedisManager)
	userConfigService := new(mocks.UserConfigService)
	userPermissionManager := new(mocks.UserPermissionManager)
	adminOrderManager := new(mocks.AdminOrderManager)
	engineCommunicator := new(mocks.EngineCommunicator)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	orderService := order.NewOrderService(db, orderRepo, orderCreateManager, eh, cs, pg, userBalanceService, redisManager,
		userConfigService, userPermissionManager, adminOrderManager, engineCommunicator, configs, logger)

	u := user.User{
		Status: user.StatusVerified,
	}

	params := order.FulfillOrderParams{
		ID: 1,
	}

	res, statusCode := orderService.FulfillOrder(&u, params)

	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "order not found", res.Message)
	orderRepo.AssertExpectations(t)
}

func TestService_FulfillOrder_StatusIsNotOpen(t *testing.T) {
	db := &gorm.DB{}
	orderRepo := new(mocks.OrderRepository)
	//o := &order.Order{}
	orderRepo.On("GetOrderByID", int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		o := args.Get(1).(*order.Order)
		o.ID = 1
		o.Status = "FILLED"
	})

	orderCreateManager := new(mocks.OrderCreateManager)

	eh := new(mocks.EventsHandler)
	cs := new(mocks.CurrencyService)
	pg := new(mocks.PriceGenerator)
	userBalanceService := new(mocks.UserBalanceService)
	redisManager := new(mocks.OrderRedisManager)
	userConfigService := new(mocks.UserConfigService)
	userPermissionManager := new(mocks.UserPermissionManager)
	adminOrderManager := new(mocks.AdminOrderManager)
	engineCommunicator := new(mocks.EngineCommunicator)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	orderService := order.NewOrderService(db, orderRepo, orderCreateManager, eh, cs, pg, userBalanceService, redisManager,
		userConfigService, userPermissionManager, adminOrderManager, engineCommunicator, configs, logger)

	u := user.User{
		ID: 22,
	}

	params := order.FulfillOrderParams{
		ID: 1,
	}

	res, statusCode := orderService.FulfillOrder(&u, params)
	assert.Equal(t, http.StatusUnprocessableEntity, statusCode)
	assert.Equal(t, false, res.Status)
	assert.Equal(t, "order status is not open", res.Message)
	orderRepo.AssertExpectations(t)

}

func TestService_FulfillOrder_Successful(t *testing.T) {
	db := &gorm.DB{}
	orderRepo := new(mocks.OrderRepository)
	//o := &order.Order{}
	orderRepo.On("GetOrderByID", int64(1), mock.Anything).Once().Return(nil).Run(func(args mock.Arguments) {
		o := args.Get(1).(*order.Order)
		o.ID = 1
		o.Status = "OPEN"
	})

	orderCreateManager := new(mocks.OrderCreateManager)

	eh := new(mocks.EventsHandler)
	cs := new(mocks.CurrencyService)
	pg := new(mocks.PriceGenerator)
	userBalanceService := new(mocks.UserBalanceService)
	redisManager := new(mocks.OrderRedisManager)
	userConfigService := new(mocks.UserConfigService)
	userPermissionManager := new(mocks.UserPermissionManager)
	adminOrderManager := new(mocks.AdminOrderManager)
	adminOrderManager.On("TryToFulfillOrder", mock.Anything).Once().Return(nil)
	engineCommunicator := new(mocks.EngineCommunicator)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)

	orderService := order.NewOrderService(db, orderRepo, orderCreateManager, eh, cs, pg, userBalanceService, redisManager,
		userConfigService, userPermissionManager, adminOrderManager, engineCommunicator, configs, logger)

	u := user.User{
		ID: 22,
	}

	params := order.FulfillOrderParams{
		ID: 1,
	}

	res, statusCode := orderService.FulfillOrder(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)
	orderRepo.AssertExpectations(t)
	adminOrderManager.AssertExpectations(t)

}

func TestService_GetOrderDetail(t *testing.T) {
	db := &gorm.DB{}
	orderRepo := new(mocks.OrderRepository)
	btcUsdtPair := currency.Pair{
		ID:   1,
		Name: "BTC-USDT",
	}

	orders := []order.Order{
		{
			ID:                  1,
			Type:                "BUY",
			ExchangeType:        "LIMIT",
			Status:              "FILLED",
			Pair:                btcUsdtPair,
			PayedByAmount:       sql.NullString{String: "5000", Valid: true},
			DemandedAmount:      sql.NullString{String: "0.1", Valid: true},
			FinalPayedByAmount:  sql.NullString{String: "2500", Valid: true},
			FinalDemandedAmount: sql.NullString{String: "0.05", Valid: true},
			Price:               sql.NullString{String: "50000", Valid: true},
			TradePrice:          sql.NullString{String: "50000", Valid: true},
			Path:                sql.NullString{String: "1,", Valid: true},
			FeePercentage:       sql.NullFloat64{Float64: 0.01, Valid: true},
			CreatedAt:           time.Now(),
		},
		{
			ID:                  2,
			Type:                "BUY",
			ExchangeType:        "LIMIT",
			Status:              "FILLED",
			Pair:                btcUsdtPair,
			PayedByAmount:       sql.NullString{String: "2500", Valid: true},
			DemandedAmount:      sql.NullString{String: "0.05", Valid: true},
			FinalPayedByAmount:  sql.NullString{String: "2500", Valid: true},
			FinalDemandedAmount: sql.NullString{String: "0.05", Valid: true},
			Price:               sql.NullString{String: "50000", Valid: true},
			TradePrice:          sql.NullString{String: "50000", Valid: true},
			Path:                sql.NullString{String: "1,2,", Valid: true},
			FeePercentage:       sql.NullFloat64{Float64: 0.01, Valid: true},
			CreatedAt:           time.Now(),
		},
	}
	orderRepo.On("GetUserOrderDetailsByID", int64(1), 21).Once().Return(orders)

	orderCreateManager := new(mocks.OrderCreateManager)
	eh := new(mocks.EventsHandler)
	cs := new(mocks.CurrencyService)
	pg := new(mocks.PriceGenerator)
	userBalanceService := new(mocks.UserBalanceService)
	redisManager := new(mocks.OrderRedisManager)
	userConfigService := new(mocks.UserConfigService)
	userPermissionManager := new(mocks.UserPermissionManager)
	adminOrderManager := new(mocks.AdminOrderManager)
	engineCommunicator := new(mocks.EngineCommunicator)
	configs := new(mocks.Configs)
	configs.On("GetEnv").Once().Return(platform.EnvTest)
	logger := new(mocks.Logger)

	orderService := order.NewOrderService(db, orderRepo, orderCreateManager, eh, cs, pg, userBalanceService, redisManager,
		userConfigService, userPermissionManager, adminOrderManager, engineCommunicator, configs, logger)

	u := user.User{
		ID: 21,
	}

	params := order.GetOrderDetailParams{
		ID: 1,
	}

	res, statusCode := orderService.GetOrderDetail(&u, params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)
	details, ok := res.Data.([]order.DetailResponse)
	if !ok {
		t.Error("can not cast response to struct")
	}
	assert.Equal(t, 2, len(details))

	order1 := details[0]
	assert.Equal(t, "BTC-USDT", order1.Pair)
	assert.Equal(t, "buy", order1.Type)
	assert.Equal(t, "50000", order1.Price)
	assert.Equal(t, "2500 USDT", order1.Amount)
	assert.Equal(t, "0.0005 BTC", order1.Fee)
	assert.Equal(t, "0.05 BTC", order1.Executed)

	order2 := details[1]
	assert.Equal(t, "BTC-USDT", order2.Pair)
	assert.Equal(t, "buy", order2.Type)
	assert.Equal(t, "50000", order2.Price)
	assert.Equal(t, "2500 USDT", order2.Amount)
	assert.Equal(t, "0.0005 BTC", order2.Fee)
	assert.Equal(t, "0.05 BTC", order2.Executed)

	orderRepo.AssertExpectations(t)
}
