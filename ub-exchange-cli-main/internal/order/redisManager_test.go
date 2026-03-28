// Package order_test tests the RedisManager for stop order queue operations. Covers:
//   - AddStopOrderToQueue: adds a stop order to the Redis sorted set with correct key and score
//   - GetStopOrdersFromQueue: retrieves BUY and SELL stop orders within a price range
//   - RemoveStopOrderFromQueue: removes a specific stop order from the Redis sorted set
//   - Exists: checks whether a stop order exists in the queue via ZScore
//
// Test data: mocked RedisClient with ZAdd/ZRangeByScoreWithScores/ZRem/ZScore expectations
// and BTC-USDT BUY order fixtures with stop point prices.
package order_test

import (
	"context"
	"database/sql"
	"exchange-go/internal/currency"
	"exchange-go/internal/mocks"
	"exchange-go/internal/order"
	"testing"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRedisManager_AddStopOrderToQueue(t *testing.T) {
	rc := new(mocks.RedisClient)
	rc.On("ZAdd", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(true, nil)

	rm := order.NewRedisManager(rc)
	pair := currency.Pair{
		ID:   1,
		Name: "BTC-USDT",
	}
	o := order.Order{
		Type:           "BUY",
		Pair:           pair,
		StopPointPrice: sql.NullString{String: "50000", Valid: true},
	}
	ctx := context.Background()
	err := rm.AddStopOrderToQueue(ctx, o)
	assert.Nil(t, err)
	rc.AssertExpectations(t)
}

func TestRedisManager_GetStopOrdersFromQueue(t *testing.T) {
	rc := new(mocks.RedisClient)
	data := []redis.Z{
		{},
		{},
	}
	rc.On("ZRangeByScoreWithScores", mock.Anything, "queue:stop:order:BUY:BTC-USDT", mock.Anything).Once().Return(data, nil)
	rc.On("ZRangeByScoreWithScores", mock.Anything, "queue:stop:order:SELL:BTC-USDT", mock.Anything).Once().Return(data, nil)

	rm := order.NewRedisManager(rc)

	ctx := context.Background()
	pairName := "BTC-USDT"
	formerPrice := "48000"
	price := "50000"
	res, err := rm.GetStopOrdersFromQueue(ctx, pairName, formerPrice, price, false)
	assert.Nil(t, err)
	assert.Equal(t, 4, len(res))
	rc.AssertExpectations(t)

}

func TestRedisManager_RemoveStopOrderFromQueue(t *testing.T) {
	rc := new(mocks.RedisClient)
	rc.On("ZRem", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(true, nil)

	rm := order.NewRedisManager(rc)
	pair := currency.Pair{
		ID:   1,
		Name: "BTC-USDT",
	}
	o := order.Order{
		Type:           "BUY",
		Pair:           pair,
		StopPointPrice: sql.NullString{String: "50000", Valid: true},
	}
	ctx := context.Background()
	err := rm.RemoveStopOrderFromQueue(ctx, o, "BTC-USDT")
	assert.Nil(t, err)
	rc.AssertExpectations(t)
}

func TestRedisManager_Exists(t *testing.T) {
	rc := new(mocks.RedisClient)
	rc.On("ZScore", mock.Anything, "queue:stop:order:BUY:BTC-USDT", mock.Anything).Once().Return(float64(1), nil)

	rm := order.NewRedisManager(rc)
	pair := currency.Pair{
		ID:   1,
		Name: "BTC-USDT",
	}
	o := order.Order{
		Type:           "BUY",
		Pair:           pair,
		StopPointPrice: sql.NullString{String: "50000", Valid: true},
	}
	ctx := context.Background()
	exists, err := rm.Exists(ctx, o)
	assert.Nil(t, err)
	assert.Equal(t, true, exists)
	rc.AssertExpectations(t)
}
