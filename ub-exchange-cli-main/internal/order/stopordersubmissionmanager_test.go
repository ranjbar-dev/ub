// Package order_test tests the StopOrderSubmissionManager. Covers:
//   - Submit: fetches triggered stop orders from Redis, updates DB status, and fires order creation events
//   - SubmitOrderInDb: updates a single stop order's current market price in the database
//
// Test data: sqlmock MySQL DB, mocked order repository, Redis manager with stop order
// queue data, events handler, and BTC-USDT pair with price transitions.
package order_test

import (
	"context"
	"exchange-go/internal/mocks"
	"exchange-go/internal/order"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// submissionQueryMatcher is a sqlmock.QueryMatcher that accepts any SQL query,
// allowing tests to focus on submission logic rather than exact SQL strings.
type submissionQueryMatcher struct {
}

// Match always returns nil, effectively disabling SQL query matching so that
// any expected SQL statement is considered a match against any actual SQL.
func (submissionQueryMatcher) Match(expectedSQL, actualSQL string) error {
	return nil
}

func TestStopOrderSubmissionManager_Submit(t *testing.T) {
	qm := submissionQueryMatcher{}
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
	//dbMock.ExpectBegin()
	dbMock.ExpectExec("UPDATE orders").WillReturnResult(sqlmock.NewResult(1, 1))

	orderRepo := new(mocks.OrderRepository)
	orders := []order.Order{
		{
			ID:     1,
			Status: order.StatusOpen,
		},
	}
	orderRepo.On("GetOrdersByIds", mock.Anything).Once().Return(orders)
	liveDataService := new(mocks.LiveData)
	orderRedisManager := new(mocks.OrderRedisManager)
	inRedisOrders := []redis.Z{
		{
			Score:  50000,
			Member: "1",
		},
	}
	orderRedisManager.On("GetStopOrdersFromQueue", mock.Anything, "BTC-USDT", "51000", "50000", false).Once().Return(inRedisOrders, nil)
	orderRedisManager.On("RemoveStopOrderFromQueue", mock.Anything, mock.Anything, "BTC-USDT").Once().Return(nil)
	eh := new(mocks.EventsHandler)
	eh.On("HandleOrderCreation", mock.Anything, true).Once().Return()
	logger := new(mocks.Logger)

	sm := order.NewStopOrderSubmissionManager(db, orderRepo, liveDataService, orderRedisManager, eh, logger)
	ctx := context.Background()
	pairName := "BTC-USDT"
	price := "50000"
	formerPrice := "51000"
	sm.Submit(ctx, pairName, price, formerPrice)
	time.Sleep(20 * time.Millisecond)
	orderRepo.AssertExpectations(t)
	orderRedisManager.AssertExpectations(t)
	eh.AssertExpectations(t)
}

func TestStopOrderSubmissionManager_SubmitOrderInDb(t *testing.T) {
	qm := submissionQueryMatcher{}
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
	//dbMock.ExpectBegin()
	dbMock.ExpectExec("UPDATE orders").WillReturnResult(sqlmock.NewResult(1, 1))

	orderRepo := new(mocks.OrderRepository)

	liveDataService := new(mocks.LiveData)
	orderRedisManager := new(mocks.OrderRedisManager)
	eh := new(mocks.EventsHandler)
	logger := new(mocks.Logger)

	sm := order.NewStopOrderSubmissionManager(db, orderRepo, liveDataService, orderRedisManager, eh, logger)
	ctx := context.Background()
	o := &order.Order{
		ID: 1,
	}

	price := "50000.00000000"
	err = sm.SubmitOrderInDb(ctx, o, price)

	assert.Nil(t, err)
}
