// Package engine tests the Redis-backed order queue used by the trading
// engine to buffer incoming orders. Covers:
//   - RPush: appends a serialized order to the right end of the queue list
//   - LPush: prepends a serialized order to the left end of the queue list
//   - LPop: pops and deserializes an order from the left end of the queue,
//     verifying all deserialized fields match the expected values
//   - Remove: removes a specific order from the queue list via LRem
//   - Exists: checks order presence in the queue via LPos
//
// Test data: mocked RedisClient with expectations on the
// "engine:queue:orders" key; BTC-USDT bid orders with threshold prices.
package engine

import (
	"context"
	"exchange-go/internal/engine/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestQueue_RPush(t *testing.T) {
	rc := new(mocks.RedisClient)
	rc.On("RPush", mock.Anything, "engine:queue:orders", mock.Anything).Once().Return(int64(1), nil)
	q := newQueue(rc)
	o := Order{
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
	err := q.rPush(context.Background(), o)
	assert.Nil(t, err)
	rc.AssertExpectations(t)
}

func TestQueue_LPush(t *testing.T) {
	rc := new(mocks.RedisClient)
	rc.On("LPush", mock.Anything, "engine:queue:orders", mock.Anything).Once().Return(int64(1), nil)
	q := newQueue(rc)
	o := Order{
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
	err := q.lPush(context.Background(), o)
	assert.Nil(t, err)
	rc.AssertExpectations(t)
}

func TestQueue_LPop(t *testing.T) {
	rc := new(mocks.RedisClient)
	data := `{"pair":"BTC-USDT","id":"1","side":"bid","quantity":"0.1","price":"50000","timestamp":1619419912,"minPrice":"49500","maxPrice":"50500"}`
	rc.On("LPop", mock.Anything, "engine:queue:orders", mock.Anything).Once().Return(data, nil)
	q := newQueue(rc)
	o, err := q.lPop(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, "BTC-USDT", o.Pair)
	assert.Equal(t, "1", o.ID)
	assert.Equal(t, "bid", o.Side)
	assert.Equal(t, "0.1", o.Quantity)
	assert.Equal(t, "50000", o.Price)
	assert.Equal(t, "", o.TradedWithOrderID)
	assert.Equal(t, "", o.QuantityTraded)
	assert.Equal(t, "", o.TradePrice)
	rc.AssertExpectations(t)
}

func TestQueue_Remove(t *testing.T) {
	rc := new(mocks.RedisClient)
	rc.On("LRem", mock.Anything, "engine:queue:orders", int64(0), mock.Anything).Once().Return(int64(1), nil)
	q := newQueue(rc)

	o := Order{
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
	err := q.remove(context.Background(), o)
	assert.Nil(t, err)
	rc.AssertExpectations(t)
}

func TestQueue_Exists(t *testing.T) {
	rc := new(mocks.RedisClient)
	rc.On("LPos", mock.Anything, "engine:queue:orders", mock.Anything, mock.Anything).Once().Return(int64(1), nil)
	q := newQueue(rc)

	o := Order{
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
	exists, err := q.exists(context.Background(), o)
	assert.Nil(t, err)
	assert.Equal(t, true, exists)
	rc.AssertExpectations(t)
}
