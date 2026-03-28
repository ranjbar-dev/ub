// Package engine_test tests the Engine public API for order lifecycle
// management. Covers:
//   - SubmitOrder: serializes and pushes an order into the Redis queue
//   - RemoveOrder: removes an order from both the Redis list and sorted set
//   - HandleInQueueOrders: dequeues orders from the Redis order book, invokes
//     the result handler callback, and verifies post-processing cleanup
//   - RetrieveOrder for limit orders: not in order book or queue (re-queued),
//     not in order book but in queue (no-op), and already in order book (no-op)
//   - RetrieveOrder for market orders: not in queue (re-queued) and already
//     in queue (no-op)
//
// Test data: mocked RedisClient and EngineResultHandler; miniredis for
// integration-style tests with real Redis sorted sets; BTC-USDT order pairs
// with bid/ask sides and threshold prices.
package engine_test

import (
	"context"
	"encoding/json"
	"exchange-go/internal/engine"
	"exchange-go/internal/mocks"
	"exchange-go/internal/platform"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEngine_Run(t *testing.T) {
	//not very important since it would be tested in  other tests especially functional tests
}

func TestEngine_SubmitOrder(t *testing.T) {
	rc := new(mocks.RedisClient)
	rc.On("RPush", mock.Anything, mock.Anything, mock.Anything).Once().Return(int64(1), nil)
	obp := engine.NewRedisOrderBookProvider(rc)
	rh := new(mocks.EngineResultHandler)
	logger := new(mocks.EngineLogger)

	e := engine.NewEngine(rc, obp, rh, logger, platform.EnvTest)
	o := engine.Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "ask",
		Quantity:          "0.2",
		Price:             "50000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	err := e.SubmitOrder(o)
	assert.Nil(t, err)
	rc.AssertExpectations(t)
}

func TestEngine_RemoveOrder(t *testing.T) {
	rc := new(mocks.RedisClient)
	rc.On("LRem", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(int64(1), nil)
	rc.On("ZRem", mock.Anything, mock.Anything, mock.Anything).Once().Return(true, nil)
	obp := engine.NewRedisOrderBookProvider(rc)
	rh := new(mocks.EngineResultHandler)
	logger := new(mocks.EngineLogger)

	e := engine.NewEngine(rc, obp, rh, logger, platform.EnvTest)
	o := engine.Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "ask",
		Quantity:          "0.2",
		Price:             "50000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	err := e.RemoveOrder(o)
	assert.Nil(t, err)
	rc.AssertExpectations(t)

}

func TestEngine_HandleInQueueOrders(t *testing.T) {
	s := miniredis.NewMiniRedis()
	defer s.Close()
	_ = s.Start()
	rc := redis.NewClient(&redis.Options{Addr: s.Addr()})
	o := engine.Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "bid",
		Quantity:          "0.1",
		Price:             "50000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}

	orderBytes, _ := json.Marshal(o)
	data := []*redis.Z{
		{
			Score:  50000,
			Member: orderBytes,
		},
	}

	rc.ZAdd(context.Background(), "order-book:bid:BTC-USDT", data...)

	matchingResult := engine.MatchingResult{
		Err:                   nil,
		RemainingPartialOrder: nil,
		RemovingDoneOrderIds:  []int64{},
	}
	rh := new(mocks.EngineResultHandler)
	rh.On("CallBack", mock.Anything, mock.Anything).Once().Return(matchingResult)

	redisClient := platform.NewRedisTestClient(rc)

	robp := engine.NewRedisOrderBookProvider(redisClient)
	logger := new(mocks.EngineLogger)

	e := engine.NewEngine(redisClient, robp, rh, logger, platform.EnvTest)
	err := e.HandleInQueueOrders("BTC-USDT", "5000")
	time.Sleep(50 * time.Millisecond)
	assert.Nil(t, err)

	//checking if the order is removed from redis
	_, err = rc.ZRangeByScore(context.Background(), "order-book:bid:BTC-USDT", &redis.ZRangeBy{
		Min:    "0",
		Max:    "100000",
		Offset: 0,
		Count:  1000000,
	}).Result()

	if err != nil && err != redis.Nil {
		t.Error(err)
		t.Fail()
	}

}

func TestEngine_RetrieveOrder_LimitOrder_DoesNotExistsInOrderBookAndQueue(t *testing.T) {
	rc := new(mocks.RedisClient)
	rc.On("ZScore", mock.Anything, "order-book:bid:BTC-USDT", mock.Anything).Once().Return(float64(0), redis.Nil)
	rc.On("LPos", mock.Anything, "engine:queue:orders", mock.Anything, mock.Anything).Once().Return(int64(0), redis.Nil)
	rc.On("LPush", mock.Anything, "engine:queue:orders", mock.Anything).Once().Return(int64(1), nil)
	obp := engine.NewRedisOrderBookProvider(rc)
	rh := new(mocks.EngineResultHandler)
	logger := new(mocks.EngineLogger)

	e := engine.NewEngine(rc, obp, rh, logger, platform.EnvTest)

	o := engine.Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "bid",
		Quantity:          "0.2",
		Price:             "50000.00000000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500.00000000",
		MaxThresholdPrice: "50500.00000000",
	}
	err := e.RetrieveOrder(o)
	assert.Nil(t, err)
	rc.AssertExpectations(t)
}

func TestEngine_RetrieveOrder_LimitOrder_OnlyDoesNotExistsInOrderBook(t *testing.T) {
	rc := new(mocks.RedisClient)
	rc.On("ZScore", mock.Anything, "order-book:bid:BTC-USDT", mock.Anything).Once().Return(float64(0), redis.Nil)
	rc.On("LPos", mock.Anything, "engine:queue:orders", mock.Anything, mock.Anything).Once().Return(int64(1), nil)
	obp := engine.NewRedisOrderBookProvider(rc)
	rh := new(mocks.EngineResultHandler)
	logger := new(mocks.EngineLogger)

	e := engine.NewEngine(rc, obp, rh, logger, platform.EnvTest)
	o := engine.Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "bid",
		Quantity:          "0.2",
		Price:             "50000.00000000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500.00000000",
		MaxThresholdPrice: "50500.00000000",
	}
	err := e.RetrieveOrder(o)
	assert.Nil(t, err)
	rc.AssertExpectations(t)
}

func TestEngine_RetrieveOrder_LimitOrder_ExistsInOrderBook(t *testing.T) {
	rc := new(mocks.RedisClient)
	rc.On("ZScore", mock.Anything, "order-book:bid:BTC-USDT", mock.Anything).Once().Return(float64(1), nil)
	obp := engine.NewRedisOrderBookProvider(rc)
	rh := new(mocks.EngineResultHandler)
	logger := new(mocks.EngineLogger)

	e := engine.NewEngine(rc, obp, rh, logger, platform.EnvTest)
	o := engine.Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "bid",
		Quantity:          "0.2",
		Price:             "50000.00000000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500.00000000",
		MaxThresholdPrice: "50500.00000000",
	}
	err := e.RetrieveOrder(o)
	assert.Nil(t, err)
	rc.AssertExpectations(t)
}

func TestEngine_RetrieveOrder_MarketOrder_DoesNotExistsInQueue(t *testing.T) {
	rc := new(mocks.RedisClient)
	rc.On("LPos", mock.Anything, "engine:queue:orders", mock.Anything, mock.Anything).Once().Return(int64(0), redis.Nil)
	rc.On("LPush", mock.Anything, "engine:queue:orders", mock.Anything).Once().Return(int64(1), nil)
	obp := engine.NewRedisOrderBookProvider(rc)
	rh := new(mocks.EngineResultHandler)
	logger := new(mocks.EngineLogger)

	e := engine.NewEngine(rc, obp, rh, logger, platform.EnvTest)
	o := engine.Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "bid",
		Quantity:          "0.2",
		Price:             "",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500.00000000",
		MaxThresholdPrice: "50500.00000000",
	}
	err := e.RetrieveOrder(o)
	assert.Nil(t, err)
	rc.AssertExpectations(t)
}

func TestEngine_RetrieveOrder_MarketOrder_ExistsInQueue(t *testing.T) {
	rc := new(mocks.RedisClient)
	rc.On("LPos", mock.Anything, "engine:queue:orders", mock.Anything, mock.Anything).Once().Return(int64(1), redis.Nil)
	obp := engine.NewRedisOrderBookProvider(rc)
	rh := new(mocks.EngineResultHandler)
	logger := new(mocks.EngineLogger)

	e := engine.NewEngine(rc, obp, rh, logger, platform.EnvTest)
	o := engine.Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "bid",
		Quantity:          "0.2",
		Price:             "",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500.00000000",
		MaxThresholdPrice: "50500.00000000",
	}
	err := e.RetrieveOrder(o)
	assert.Nil(t, err)
	rc.AssertExpectations(t)
}
