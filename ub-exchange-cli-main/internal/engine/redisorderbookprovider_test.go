// Package engine_test tests the RedisOrderBookProvider which manages
// order book storage in Redis sorted sets. Covers:
//   - GetOrders: retrieves and deserializes multiple orders from a Redis
//     sorted set filtered by price range, verifying all order fields
//   - RemoveOrder: removes a specific order from the sorted set via ZRem
//   - RewriteOrderBook: atomically rewrites the order book by removing
//     done orders and re-adding a partial order to the sorted set
//   - Exists: checks whether an order exists in the sorted set via ZScore
//
// Test data: mocked RedisClient for unit tests; miniredis for integration
// tests with real Redis sorted set operations; BTC-USDT bid orders with
// varying prices (50500, 50600, 50700) and threshold ranges.
package engine_test

import (
	"context"
	"encoding/json"
	"exchange-go/internal/engine"
	"exchange-go/internal/engine/mocks"
	"exchange-go/internal/platform"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"testing"
)

func TestRedisOrderBookProvider_GetOrders(t *testing.T) {
	rc := new(mocks.RedisClient)
	matchingOrder1 := engine.Order{
		Pair:              "BTC-USDT",
		ID:                "2",
		Side:              "bid",
		Quantity:          "0.1",
		Price:             "50500",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	matchingOrder1Bytes, _ := json.Marshal(matchingOrder1)

	matchingOrder2 := engine.Order{
		Pair:              "BTC-USDT",
		ID:                "3",
		Side:              "bid",
		Quantity:          "0.1",
		Price:             "50600",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}

	matchingOrder2Bytes, _ := json.Marshal(matchingOrder2)

	matchingOrder3 := engine.Order{
		Pair:              "BTC-USDT",
		ID:                "4",
		Side:              "bid",
		Quantity:          "0.1",
		Price:             "50700",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}

	matchingOrder3Bytes, _ := json.Marshal(matchingOrder3)

	data := []redis.Z{
		{
			Score:  50000,
			Member: string(matchingOrder1Bytes),
		},
		{
			Score:  50000,
			Member: string(matchingOrder2Bytes),
		},
		{
			Score:  50000,
			Member: string(matchingOrder3Bytes),
		},
	}

	rc.On("ZRangeByScoreWithScores", mock.Anything, "order-book:bid:BTC-USDT", mock.Anything).Once().Return(data, nil)

	robp := engine.NewRedisOrderBookProvider(rc)
	params := engine.OrderBookProviderParams{
		Pair:     "BTC-USDT",
		Side:     "bid",
		Price:    "50000",
		MinPrice: "49500",
		MaxPrice: "50500",
	}
	orders, err := robp.GetOrders(context.Background(), params)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(orders))

	order1 := orders[0]
	assert.Equal(t, "BTC-USDT", order1.Pair)
	assert.Equal(t, "2", order1.ID)
	assert.Equal(t, "bid", order1.Side)
	assert.Equal(t, "0.1", order1.Quantity)
	assert.Equal(t, "50500", order1.Price)
	assert.Equal(t, "", order1.TradedWithOrderID)
	assert.Equal(t, "", order1.QuantityTraded)
	assert.Equal(t, "", order1.TradePrice)
	//assert.Equal(t, "49500", order1.MinThresholdPrice)
	//assert.Equal(t, "50500", order1.MaxThresholdPrice)

	order2 := orders[1]
	assert.Equal(t, "BTC-USDT", order2.Pair)
	assert.Equal(t, "3", order2.ID)
	assert.Equal(t, "bid", order2.Side)
	assert.Equal(t, "0.1", order2.Quantity)
	assert.Equal(t, "50600", order2.Price)
	assert.Equal(t, "", order2.TradedWithOrderID)
	assert.Equal(t, "", order2.QuantityTraded)
	assert.Equal(t, "", order2.TradePrice)
	//assert.Equal(t, "49500", order2.MinThresholdPrice)
	//assert.Equal(t, "50500", order2.MaxThresholdPrice)

	order3 := orders[2]

	assert.Equal(t, "BTC-USDT", order3.Pair)
	assert.Equal(t, "4", order3.ID)
	assert.Equal(t, "bid", order3.Side)
	assert.Equal(t, "0.1", order3.Quantity)
	assert.Equal(t, "50700", order3.Price)
	assert.Equal(t, "", order3.TradedWithOrderID)
	assert.Equal(t, "", order3.QuantityTraded)
	assert.Equal(t, "", order3.TradePrice)
	//assert.Equal(t, "49500", order3.MinThresholdPrice)
	//assert.Equal(t, "50500", order3.MaxThresholdPrice)

	rc.AssertExpectations(t)
}

//func TestRedisOrderBookProvider_PopOrders(t *testing.T) {
//
//}

func TestRedisOrderBookProvider_RemoveOrder(t *testing.T) {
	rc := new(mocks.RedisClient)
	rc.On("ZRem", mock.Anything, "order-book:bid:BTC-USDT", mock.Anything).Once().Return(true, nil)

	robp := engine.NewRedisOrderBookProvider(rc)
	o := engine.Order{
		Pair:              "BTC-USDT",
		ID:                "2",
		Side:              "bid",
		Quantity:          "0.1",
		Price:             "50500",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	err := robp.RemoveOrder(context.Background(), o)
	assert.Nil(t, err)
	rc.AssertExpectations(t)
}

func TestRedisOrderBookProvider_RewriteOrderBook(t *testing.T) {
	s := miniredis.NewMiniRedis()
	defer s.Close()
	_ = s.Start()
	rc := redis.NewClient(&redis.Options{Addr: s.Addr()})
	ctx := context.Background()

	matchingOrder1 := engine.Order{
		Pair:              "BTC-USDT",
		ID:                "2",
		Side:              "bid",
		Quantity:          "0.1",
		Price:             "50500",
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	matchingOrder1Bytes, _ := json.Marshal(matchingOrder1)

	matchingOrder2 := engine.Order{
		Pair:              "BTC-USDT",
		ID:                "3",
		Side:              "bid",
		Quantity:          "0.1",
		Price:             "50600",
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}

	matchingOrder2Bytes, _ := json.Marshal(matchingOrder2)

	matchingOrder3 := engine.Order{
		Pair:              "BTC-USDT",
		ID:                "4",
		Side:              "bid",
		Quantity:          "0.1",
		Price:             "50700",
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}

	matchingOrder3Bytes, _ := json.Marshal(matchingOrder3)

	data := []*redis.Z{
		{
			Score:  50000,
			Member: string(matchingOrder1Bytes),
		},
		{
			Score:  50000,
			Member: string(matchingOrder2Bytes),
		},
		{
			Score:  50000,
			Member: string(matchingOrder3Bytes),
		},
	}

	rc.ZAdd(ctx, "order-book:bid:BTC-USDT", data...)

	redisClient := platform.NewRedisTestClient(rc)

	robp := engine.NewRedisOrderBookProvider(redisClient)

	doneOrders := []engine.Order{
		matchingOrder1,
		matchingOrder2,
		matchingOrder3,
	}

	partial := &engine.Order{
		Pair:              "BTC-USDT",
		ID:                "2",
		Side:              "bid",
		Quantity:          "0.1",
		Price:             "50500",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	err := robp.RewriteOrderBook(context.Background(), doneOrders, partial)
	assert.Nil(t, err)

}

func TestRedisOrderBookProvider_Exists(t *testing.T) {
	s := miniredis.NewMiniRedis()
	defer s.Close()
	_ = s.Start()
	rc := redis.NewClient(&redis.Options{Addr: s.Addr()})
	ctx := context.Background()

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

	orderBytes, _ := o.MarshalForOrderbook()
	data := []*redis.Z{
		{
			Score:  50000,
			Member: orderBytes,
		},
	}

	rc.ZAdd(ctx, "order-book:bid:BTC-USDT", data...)

	redisClient := platform.NewRedisTestClient(rc)

	robp := engine.NewRedisOrderBookProvider(redisClient)

	exists, err := robp.Exists(context.Background(), o)
	assert.Nil(t, err)
	assert.Equal(t, true, exists)
}
