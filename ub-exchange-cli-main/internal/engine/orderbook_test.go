// Package engine tests the order book matching logic for the trading engine.
// Covers:
//   - Limit ask orders: single match at same price, three matches at same price,
//     three matches with remaining partial (taker and maker), different prices
//     with remaining partial, and no-match (partial only)
//   - Market ask orders: same scenarios as limit asks but without price constraints
//   - Limit bid orders: single match at same price, three matches at same price,
//     three matches with remaining partial (taker and maker), different prices
//     with remaining partial, and no-match (partial only)
//   - Market bid orders: same scenarios as limit bids but without price constraints
//   - removeOrder: verifies order removal from the order book via Redis ZRem
//   - orderExists: verifies existence check against the Redis sorted set
//
// Test data: mocked RedisClient with ZRangeByScoreWithScores and ZRem
// expectations; orders use BTC-USDT pair with decimal quantities and prices.
package engine

import (
	"context"
	"encoding/json"
	"exchange-go/internal/engine/mocks"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestOrderBook_processOrder_Limit_Ask_MatchedWithSingleOrder_SamePrice_WithoutRemainingPartial(t *testing.T) {
	rc := new(mocks.RedisClient)
	matchingOrder := Order{
		Pair:              "BTC-USDT",
		ID:                "2",
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

	matchingOrderBytes, _ := json.Marshal(matchingOrder)
	data := []redis.Z{
		{
			Score:  1,
			Member: string(matchingOrderBytes),
		},
	}
	rc.On("ZRangeByScoreWithScores", mock.Anything, "order-book:bid:BTC-USDT", mock.Anything).Once().Return(data, nil)
	obp := NewRedisOrderBookProvider(rc)

	ob := newOrderBook("BTC-USDT", obp)
	o := Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "ask",
		Quantity:          "0.1",
		Price:             "50000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	doneOrders, partial, err := ob.processOrder(o)
	assert.Nil(t, err)
	assert.Nil(t, partial)
	assert.Equal(t, 2, len(doneOrders))

	currentOrder := doneOrders[0]
	assert.Equal(t, "1", currentOrder.ID)
	assert.Equal(t, "ask", currentOrder.Side)
	assert.Equal(t, "50000", currentOrder.TradePrice)
	assert.Equal(t, "2", currentOrder.TradedWithOrderID)
	assert.Equal(t, "0.1", currentOrder.QuantityTraded)

	matchingDoneOrder := doneOrders[1]
	assert.Equal(t, "2", matchingDoneOrder.ID)
	assert.Equal(t, "bid", matchingDoneOrder.Side)
	assert.Equal(t, "50000", matchingDoneOrder.TradePrice)
	assert.Equal(t, "1", matchingDoneOrder.TradedWithOrderID)
	assert.Equal(t, "0.1", matchingDoneOrder.QuantityTraded)
	rc.AssertExpectations(t)
}

func TestOrderBook_processOrder_Market_Ask_MatchedWithSingleOrder_SamePrice_WithoutRemainingPartial(t *testing.T) {
	rc := new(mocks.RedisClient)
	matchingOrder := Order{
		Pair:              "BTC-USDT",
		ID:                "2",
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

	matchingOrderBytes, _ := json.Marshal(matchingOrder)
	data := []redis.Z{
		{
			Score:  1,
			Member: string(matchingOrderBytes),
		},
	}
	rc.On("ZRangeByScoreWithScores", mock.Anything, "order-book:bid:BTC-USDT", mock.Anything).Once().Return(data, nil)
	obp := NewRedisOrderBookProvider(rc)

	ob := newOrderBook("BTC-USDT", obp)
	o := Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "ask",
		Quantity:          "0.1",
		Price:             "",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	doneOrders, partial, err := ob.processOrder(o)
	assert.Nil(t, err)
	assert.Nil(t, partial)
	assert.Equal(t, 2, len(doneOrders))

	currentOrder := doneOrders[0]
	assert.Equal(t, "1", currentOrder.ID)
	assert.Equal(t, "ask", currentOrder.Side)
	assert.Equal(t, "50000", currentOrder.TradePrice)
	assert.Equal(t, "2", currentOrder.TradedWithOrderID)
	assert.Equal(t, "0.1", currentOrder.QuantityTraded)

	matchingDoneOrder := doneOrders[1]
	assert.Equal(t, "2", matchingDoneOrder.ID)
	assert.Equal(t, "bid", matchingDoneOrder.Side)
	assert.Equal(t, "50000", matchingDoneOrder.TradePrice)
	assert.Equal(t, "1", matchingDoneOrder.TradedWithOrderID)
	assert.Equal(t, "0.1", matchingDoneOrder.QuantityTraded)

	rc.AssertExpectations(t)
}

func TestOrderBook_processOrder_Limit_Ask_MatchedWithThreeOrder_SamePrice_WithoutRemainingPartial(t *testing.T) {
	rc := new(mocks.RedisClient)
	matchingOrder1 := Order{
		Pair:              "BTC-USDT",
		ID:                "2",
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
	matchingOrder1Bytes, _ := json.Marshal(matchingOrder1)

	matchingOrder2 := Order{
		Pair:              "BTC-USDT",
		ID:                "3",
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

	matchingOrder2Bytes, _ := json.Marshal(matchingOrder2)

	matchingOrder3 := Order{
		Pair:              "BTC-USDT",
		ID:                "4",
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
	obp := NewRedisOrderBookProvider(rc)

	ob := newOrderBook("BTC-USDT", obp)
	o := Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "ask",
		Quantity:          "0.3",
		Price:             "50000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	doneOrders, partial, err := ob.processOrder(o)
	assert.Nil(t, err)
	assert.Nil(t, partial)
	assert.Equal(t, 6, len(doneOrders))

	currentOrder := doneOrders[0]
	assert.Equal(t, "1", currentOrder.ID)
	assert.Equal(t, "ask", currentOrder.Side)
	assert.Equal(t, "50000", currentOrder.TradePrice)
	assert.Equal(t, "2", currentOrder.TradedWithOrderID)
	assert.Equal(t, "0.1", currentOrder.QuantityTraded)
	assert.Equal(t, "0.3", currentOrder.Quantity)

	doneOrder1 := doneOrders[1]
	assert.Equal(t, "2", doneOrder1.ID)
	assert.Equal(t, "bid", doneOrder1.Side)
	assert.Equal(t, "50000", doneOrder1.TradePrice)
	assert.Equal(t, "1", doneOrder1.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder1.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder1.Quantity)

	doneOrder2 := doneOrders[2]
	assert.Equal(t, "1", doneOrder2.ID)
	assert.Equal(t, "ask", doneOrder2.Side)
	assert.Equal(t, "50000", doneOrder2.TradePrice)
	assert.Equal(t, "3", doneOrder2.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder2.QuantityTraded)
	assert.Equal(t, "0.2", doneOrder2.Quantity)

	doneOrder3 := doneOrders[3]
	assert.Equal(t, "3", doneOrder3.ID)
	assert.Equal(t, "bid", doneOrder3.Side)
	assert.Equal(t, "50000", doneOrder3.TradePrice)
	assert.Equal(t, "1", doneOrder3.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder3.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder3.Quantity)

	doneOrder4 := doneOrders[4]
	assert.Equal(t, "1", doneOrder4.ID)
	assert.Equal(t, "ask", doneOrder4.Side)
	assert.Equal(t, "50000", doneOrder4.TradePrice)
	assert.Equal(t, "4", doneOrder4.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder4.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder4.Quantity)

	doneOrder5 := doneOrders[5]
	assert.Equal(t, "4", doneOrder5.ID)
	assert.Equal(t, "bid", doneOrder5.Side)
	assert.Equal(t, "50000", doneOrder5.TradePrice)
	assert.Equal(t, "1", doneOrder5.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder5.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder5.Quantity)

	rc.AssertExpectations(t)
}

func TestOrderBook_processOrder_Market_Ask_MatchedWithThreeOrder_SamePrice_WithoutRemainingPartial(t *testing.T) {
	rc := new(mocks.RedisClient)
	matchingOrder1 := Order{
		Pair:              "BTC-USDT",
		ID:                "2",
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
	matchingOrder1Bytes, _ := json.Marshal(matchingOrder1)

	matchingOrder2 := Order{
		Pair:              "BTC-USDT",
		ID:                "3",
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

	matchingOrder2Bytes, _ := json.Marshal(matchingOrder2)

	matchingOrder3 := Order{
		Pair:              "BTC-USDT",
		ID:                "4",
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
	obp := NewRedisOrderBookProvider(rc)

	ob := newOrderBook("BTC-USDT", obp)
	o := Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "ask",
		Quantity:          "0.3",
		Price:             "",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	doneOrders, partial, err := ob.processOrder(o)
	assert.Nil(t, err)
	assert.Nil(t, partial)
	assert.Equal(t, 6, len(doneOrders))

	currentOrder := doneOrders[0]
	assert.Equal(t, "1", currentOrder.ID)
	assert.Equal(t, "ask", currentOrder.Side)
	assert.Equal(t, "50000", currentOrder.TradePrice)
	assert.Equal(t, "2", currentOrder.TradedWithOrderID)
	assert.Equal(t, "0.1", currentOrder.QuantityTraded)
	assert.Equal(t, "0.3", currentOrder.Quantity)

	doneOrder1 := doneOrders[1]
	assert.Equal(t, "2", doneOrder1.ID)
	assert.Equal(t, "bid", doneOrder1.Side)
	assert.Equal(t, "50000", doneOrder1.TradePrice)
	assert.Equal(t, "1", doneOrder1.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder1.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder1.Quantity)

	doneOrder2 := doneOrders[2]
	assert.Equal(t, "1", doneOrder2.ID)
	assert.Equal(t, "ask", doneOrder2.Side)
	assert.Equal(t, "50000", doneOrder2.TradePrice)
	assert.Equal(t, "3", doneOrder2.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder2.QuantityTraded)
	assert.Equal(t, "0.2", doneOrder2.Quantity)

	doneOrder3 := doneOrders[3]
	assert.Equal(t, "3", doneOrder3.ID)
	assert.Equal(t, "bid", doneOrder3.Side)
	assert.Equal(t, "50000", doneOrder3.TradePrice)
	assert.Equal(t, "1", doneOrder3.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder3.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder3.Quantity)

	doneOrder4 := doneOrders[4]
	assert.Equal(t, "1", doneOrder4.ID)
	assert.Equal(t, "ask", doneOrder4.Side)
	assert.Equal(t, "50000", doneOrder4.TradePrice)
	assert.Equal(t, "4", doneOrder4.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder4.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder4.Quantity)

	doneOrder5 := doneOrders[5]
	assert.Equal(t, "4", doneOrder5.ID)
	assert.Equal(t, "bid", doneOrder5.Side)
	assert.Equal(t, "50000", doneOrder5.TradePrice)
	assert.Equal(t, "1", doneOrder5.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder5.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder5.Quantity)

	rc.AssertExpectations(t)
}

func TestOrderBook_processOrder_Limit_Ask_MatchedWithThreeOrder_SamePrice_WithRemainingPartial(t *testing.T) {
	rc := new(mocks.RedisClient)
	matchingOrder1 := Order{
		Pair:              "BTC-USDT",
		ID:                "2",
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
	matchingOrder1Bytes, _ := json.Marshal(matchingOrder1)

	matchingOrder2 := Order{
		Pair:              "BTC-USDT",
		ID:                "3",
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

	matchingOrder2Bytes, _ := json.Marshal(matchingOrder2)

	matchingOrder3 := Order{
		Pair:              "BTC-USDT",
		ID:                "4",
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
	obp := NewRedisOrderBookProvider(rc)

	ob := newOrderBook("BTC-USDT", obp)
	o := Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "ask",
		Quantity:          "0.4",
		Price:             "50000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	doneOrders, partial, err := ob.processOrder(o)
	assert.Nil(t, err)
	assert.NotNil(t, partial)
	assert.Equal(t, "1", partial.ID)
	assert.Equal(t, "50000", partial.Price)
	assert.Equal(t, "0.1", partial.Quantity)
	assert.Equal(t, "ask", partial.Side)

	assert.Equal(t, 6, len(doneOrders))

	currentOrder := doneOrders[0]
	assert.Equal(t, "1", currentOrder.ID)
	assert.Equal(t, "ask", currentOrder.Side)
	assert.Equal(t, "50000", currentOrder.TradePrice)
	assert.Equal(t, "2", currentOrder.TradedWithOrderID)
	assert.Equal(t, "0.1", currentOrder.QuantityTraded)
	assert.Equal(t, "0.4", currentOrder.Quantity)

	doneOrder1 := doneOrders[1]
	assert.Equal(t, "2", doneOrder1.ID)
	assert.Equal(t, "bid", doneOrder1.Side)
	assert.Equal(t, "50000", doneOrder1.TradePrice)
	assert.Equal(t, "1", doneOrder1.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder1.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder1.Quantity)

	doneOrder2 := doneOrders[2]
	assert.Equal(t, "1", doneOrder2.ID)
	assert.Equal(t, "ask", doneOrder2.Side)
	assert.Equal(t, "50000", doneOrder2.TradePrice)
	assert.Equal(t, "3", doneOrder2.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder2.QuantityTraded)
	assert.Equal(t, "0.3", doneOrder2.Quantity)

	doneOrder3 := doneOrders[3]
	assert.Equal(t, "3", doneOrder3.ID)
	assert.Equal(t, "bid", doneOrder3.Side)
	assert.Equal(t, "50000", doneOrder3.TradePrice)
	assert.Equal(t, "1", doneOrder3.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder3.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder3.Quantity)

	doneOrder4 := doneOrders[4]
	assert.Equal(t, "1", doneOrder4.ID)
	assert.Equal(t, "ask", doneOrder4.Side)
	assert.Equal(t, "50000", doneOrder4.TradePrice)
	assert.Equal(t, "4", doneOrder4.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder4.QuantityTraded)
	assert.Equal(t, "0.2", doneOrder4.Quantity)

	doneOrder5 := doneOrders[5]
	assert.Equal(t, "4", doneOrder5.ID)
	assert.Equal(t, "bid", doneOrder5.Side)
	assert.Equal(t, "50000", doneOrder5.TradePrice)
	assert.Equal(t, "1", doneOrder5.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder5.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder5.Quantity)

	rc.AssertExpectations(t)
}

func TestOrderBook_processOrder_Market_Ask_MatchedWithThreeOrder_SamePrice_WithRemainingPartial(t *testing.T) {
	rc := new(mocks.RedisClient)
	matchingOrder1 := Order{
		Pair:              "BTC-USDT",
		ID:                "2",
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
	matchingOrder1Bytes, _ := json.Marshal(matchingOrder1)

	matchingOrder2 := Order{
		Pair:              "BTC-USDT",
		ID:                "3",
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

	matchingOrder2Bytes, _ := json.Marshal(matchingOrder2)

	matchingOrder3 := Order{
		Pair:              "BTC-USDT",
		ID:                "4",
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
	obp := NewRedisOrderBookProvider(rc)

	ob := newOrderBook("BTC-USDT", obp)
	o := Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "ask",
		Quantity:          "0.4",
		Price:             "",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	doneOrders, partial, err := ob.processOrder(o)
	assert.Nil(t, err)
	assert.NotNil(t, partial)
	assert.Equal(t, "1", partial.ID)
	assert.Equal(t, "", partial.Price)
	assert.Equal(t, "0.1", partial.Quantity)
	assert.Equal(t, "ask", partial.Side)

	assert.Equal(t, 6, len(doneOrders))
	currentOrder := doneOrders[0]
	assert.Equal(t, "1", currentOrder.ID)
	assert.Equal(t, "ask", currentOrder.Side)
	assert.Equal(t, "50000", currentOrder.TradePrice)
	assert.Equal(t, "2", currentOrder.TradedWithOrderID)
	assert.Equal(t, "0.1", currentOrder.QuantityTraded)
	assert.Equal(t, "0.4", currentOrder.Quantity)

	doneOrder1 := doneOrders[1]
	assert.Equal(t, "2", doneOrder1.ID)
	assert.Equal(t, "bid", doneOrder1.Side)
	assert.Equal(t, "50000", doneOrder1.TradePrice)
	assert.Equal(t, "1", doneOrder1.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder1.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder1.Quantity)

	doneOrder2 := doneOrders[2]
	assert.Equal(t, "1", doneOrder2.ID)
	assert.Equal(t, "ask", doneOrder2.Side)
	assert.Equal(t, "50000", doneOrder2.TradePrice)
	assert.Equal(t, "3", doneOrder2.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder2.QuantityTraded)
	assert.Equal(t, "0.3", doneOrder2.Quantity)

	doneOrder3 := doneOrders[3]
	assert.Equal(t, "3", doneOrder3.ID)
	assert.Equal(t, "bid", doneOrder3.Side)
	assert.Equal(t, "50000", doneOrder3.TradePrice)
	assert.Equal(t, "1", doneOrder3.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder3.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder3.Quantity)

	doneOrder4 := doneOrders[4]
	assert.Equal(t, "1", doneOrder4.ID)
	assert.Equal(t, "ask", doneOrder4.Side)
	assert.Equal(t, "50000", doneOrder4.TradePrice)
	assert.Equal(t, "4", doneOrder4.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder4.QuantityTraded)
	assert.Equal(t, "0.2", doneOrder4.Quantity)

	doneOrder5 := doneOrders[5]
	assert.Equal(t, "4", doneOrder5.ID)
	assert.Equal(t, "bid", doneOrder5.Side)
	assert.Equal(t, "50000", doneOrder5.TradePrice)
	assert.Equal(t, "1", doneOrder5.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder5.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder5.Quantity)

	rc.AssertExpectations(t)
}

func TestOrderBook_processOrder_Limit_Ask_MatchedWithThreeOrder_DifferentPrice_WithRemainingPartial(t *testing.T) {
	rc := new(mocks.RedisClient)
	matchingOrder1 := Order{
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

	matchingOrder2 := Order{
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

	matchingOrder3 := Order{
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
	obp := NewRedisOrderBookProvider(rc)

	ob := newOrderBook("BTC-USDT", obp)
	o := Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "ask",
		Quantity:          "0.4",
		Price:             "50000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	doneOrders, partial, err := ob.processOrder(o)
	assert.Nil(t, err)
	assert.NotNil(t, partial)
	assert.Equal(t, "1", partial.ID)
	assert.Equal(t, "50000", partial.Price)
	assert.Equal(t, "0.1", partial.Quantity)
	assert.Equal(t, "ask", partial.Side)

	assert.Equal(t, 6, len(doneOrders))
	currentOrder := doneOrders[0]
	assert.Equal(t, "1", currentOrder.ID)
	assert.Equal(t, "ask", currentOrder.Side)
	assert.Equal(t, "50000", currentOrder.Price)
	assert.Equal(t, "50000", currentOrder.TradePrice)
	assert.Equal(t, "4", currentOrder.TradedWithOrderID)
	assert.Equal(t, "0.1", currentOrder.QuantityTraded)
	assert.Equal(t, "0.4", currentOrder.Quantity)

	doneOrder1 := doneOrders[1]
	assert.Equal(t, "4", doneOrder1.ID)
	assert.Equal(t, "bid", doneOrder1.Side)
	assert.Equal(t, "50700", doneOrder1.Price)
	assert.Equal(t, "50000", doneOrder1.TradePrice)
	assert.Equal(t, "1", doneOrder1.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder1.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder1.Quantity)

	doneOrder2 := doneOrders[2]
	assert.Equal(t, "1", doneOrder2.ID)
	assert.Equal(t, "ask", doneOrder2.Side)
	assert.Equal(t, "50000", doneOrder2.Price)
	assert.Equal(t, "50000", doneOrder2.TradePrice)
	assert.Equal(t, "3", doneOrder2.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder2.QuantityTraded)
	assert.Equal(t, "0.3", doneOrder2.Quantity)

	doneOrder3 := doneOrders[3]
	assert.Equal(t, "3", doneOrder3.ID)
	assert.Equal(t, "bid", doneOrder3.Side)
	assert.Equal(t, "50600", doneOrder3.Price)
	assert.Equal(t, "50000", doneOrder3.TradePrice)
	assert.Equal(t, "1", doneOrder3.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder3.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder3.Quantity)

	doneOrder4 := doneOrders[4]
	assert.Equal(t, "1", doneOrder4.ID)
	assert.Equal(t, "ask", doneOrder4.Side)
	assert.Equal(t, "50000", doneOrder4.Price)
	assert.Equal(t, "50000", doneOrder4.TradePrice)
	assert.Equal(t, "2", doneOrder4.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder4.QuantityTraded)
	assert.Equal(t, "0.2", doneOrder4.Quantity)

	doneOrder5 := doneOrders[5]
	assert.Equal(t, "2", doneOrder5.ID)
	assert.Equal(t, "bid", doneOrder5.Side)
	assert.Equal(t, "50500", doneOrder5.Price)
	assert.Equal(t, "50000", doneOrder5.TradePrice)
	assert.Equal(t, "1", doneOrder5.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder5.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder5.Quantity)

	rc.AssertExpectations(t)
}

func TestOrderBook_processOrder_Market_Ask_MatchedWithThreeOrder_DifferentPrice_WithRemainingPartial(t *testing.T) {
	rc := new(mocks.RedisClient)
	matchingOrder1 := Order{
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

	matchingOrder2 := Order{
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

	matchingOrder3 := Order{
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
	obp := NewRedisOrderBookProvider(rc)

	ob := newOrderBook("BTC-USDT", obp)
	o := Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "ask",
		Quantity:          "0.4",
		Price:             "",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	doneOrders, partial, err := ob.processOrder(o)
	assert.Nil(t, err)
	assert.NotNil(t, partial)
	assert.Equal(t, "1", partial.ID)
	assert.Equal(t, "", partial.Price)
	assert.Equal(t, "0.1", partial.Quantity)
	assert.Equal(t, "ask", partial.Side)

	assert.Equal(t, 6, len(doneOrders))
	currentOrder := doneOrders[0]
	assert.Equal(t, "1", currentOrder.ID)
	assert.Equal(t, "ask", currentOrder.Side)
	assert.Equal(t, "", currentOrder.Price)
	assert.Equal(t, "50700", currentOrder.TradePrice)
	assert.Equal(t, "4", currentOrder.TradedWithOrderID)
	assert.Equal(t, "0.1", currentOrder.QuantityTraded)
	assert.Equal(t, "0.4", currentOrder.Quantity)

	doneOrder1 := doneOrders[1]
	assert.Equal(t, "4", doneOrder1.ID)
	assert.Equal(t, "bid", doneOrder1.Side)
	assert.Equal(t, "50700", doneOrder1.Price)
	assert.Equal(t, "50700", doneOrder1.TradePrice)
	assert.Equal(t, "1", doneOrder1.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder1.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder1.Quantity)

	doneOrder2 := doneOrders[2]
	assert.Equal(t, "1", doneOrder2.ID)
	assert.Equal(t, "ask", doneOrder2.Side)
	assert.Equal(t, "", doneOrder2.Price)
	assert.Equal(t, "50600", doneOrder2.TradePrice)
	assert.Equal(t, "3", doneOrder2.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder2.QuantityTraded)
	assert.Equal(t, "0.3", doneOrder2.Quantity)

	doneOrder3 := doneOrders[3]
	assert.Equal(t, "3", doneOrder3.ID)
	assert.Equal(t, "bid", doneOrder3.Side)
	assert.Equal(t, "50600", doneOrder3.Price)
	assert.Equal(t, "50600", doneOrder3.TradePrice)
	assert.Equal(t, "1", doneOrder3.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder3.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder3.Quantity)

	doneOrder4 := doneOrders[4]
	assert.Equal(t, "1", doneOrder4.ID)
	assert.Equal(t, "ask", doneOrder4.Side)
	assert.Equal(t, "", doneOrder4.Price)
	assert.Equal(t, "50500", doneOrder4.TradePrice)
	assert.Equal(t, "2", doneOrder4.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder4.QuantityTraded)
	assert.Equal(t, "0.2", doneOrder4.Quantity)

	doneOrder5 := doneOrders[5]
	assert.Equal(t, "2", doneOrder5.ID)
	assert.Equal(t, "bid", doneOrder5.Side)
	assert.Equal(t, "50500", doneOrder5.Price)
	assert.Equal(t, "50500", doneOrder5.TradePrice)
	assert.Equal(t, "1", doneOrder5.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder5.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder5.Quantity)

	rc.AssertExpectations(t)
}

//maker orders are the orders which are in orderbook and not the new coming one
func TestOrderBook_processOrder_Limit_Ask_MatchedWithThreeOrder_DifferentPrice_WithRemainingPartialForMakerOrder(t *testing.T) {
	rc := new(mocks.RedisClient)
	matchingOrder1 := Order{
		Pair:              "BTC-USDT",
		ID:                "2",
		Side:              "bid",
		Quantity:          "0.3",
		Price:             "50500",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	matchingOrder1Bytes, _ := json.Marshal(matchingOrder1)

	matchingOrder2 := Order{
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

	matchingOrder3 := Order{
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
	obp := NewRedisOrderBookProvider(rc)

	ob := newOrderBook("BTC-USDT", obp)
	o := Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "ask",
		Quantity:          "0.3",
		Price:             "50000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	doneOrders, partial, err := ob.processOrder(o)
	assert.Nil(t, err)
	assert.NotNil(t, partial)
	assert.Equal(t, "2", partial.ID)
	assert.Equal(t, "50500", partial.Price)
	assert.Equal(t, "0.2", partial.Quantity)
	assert.Equal(t, "bid", partial.Side)

	assert.Equal(t, 6, len(doneOrders))
	currentOrder := doneOrders[0]
	assert.Equal(t, "1", currentOrder.ID)
	assert.Equal(t, "ask", currentOrder.Side)
	assert.Equal(t, "50000", currentOrder.Price)
	assert.Equal(t, "50000", currentOrder.TradePrice)
	assert.Equal(t, "4", currentOrder.TradedWithOrderID)
	assert.Equal(t, "0.1", currentOrder.QuantityTraded)
	assert.Equal(t, "0.3", currentOrder.Quantity)

	doneOrder1 := doneOrders[1]
	assert.Equal(t, "4", doneOrder1.ID)
	assert.Equal(t, "bid", doneOrder1.Side)
	assert.Equal(t, "50700", doneOrder1.Price)
	assert.Equal(t, "50000", doneOrder1.TradePrice)
	assert.Equal(t, "1", doneOrder1.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder1.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder1.Quantity)

	doneOrder2 := doneOrders[2]
	assert.Equal(t, "1", doneOrder2.ID)
	assert.Equal(t, "ask", doneOrder2.Side)
	assert.Equal(t, "50000", doneOrder2.Price)
	assert.Equal(t, "50000", doneOrder2.TradePrice)
	assert.Equal(t, "3", doneOrder2.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder2.QuantityTraded)
	assert.Equal(t, "0.2", doneOrder2.Quantity)

	doneOrder3 := doneOrders[3]
	assert.Equal(t, "3", doneOrder3.ID)
	assert.Equal(t, "bid", doneOrder3.Side)
	assert.Equal(t, "50600", doneOrder3.Price)
	assert.Equal(t, "50000", doneOrder3.TradePrice)
	assert.Equal(t, "1", doneOrder3.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder3.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder3.Quantity)

	doneOrder4 := doneOrders[4]
	assert.Equal(t, "1", doneOrder4.ID)
	assert.Equal(t, "ask", doneOrder4.Side)
	assert.Equal(t, "50000", doneOrder4.Price)
	assert.Equal(t, "50000", doneOrder4.TradePrice)
	assert.Equal(t, "2", doneOrder4.TradedWithOrderID)
	assert.Equal(t, "0.1000000000000000", doneOrder4.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder4.Quantity)

	doneOrder5 := doneOrders[5]
	assert.Equal(t, "2", doneOrder5.ID)
	assert.Equal(t, "bid", doneOrder5.Side)
	assert.Equal(t, "50500", doneOrder5.Price)
	assert.Equal(t, "50000", doneOrder5.TradePrice)
	assert.Equal(t, "1", doneOrder5.TradedWithOrderID)
	assert.Equal(t, "0.1000000000000000", doneOrder5.QuantityTraded)
	assert.Equal(t, "0.3", doneOrder5.Quantity)

	rc.AssertExpectations(t)
}

func TestOrderBook_processOrder_Market_Ask_MatchedWithThreeOrder_DifferentPrice_WithRemainingPartialForMakerOrder(t *testing.T) {
	rc := new(mocks.RedisClient)
	matchingOrder1 := Order{
		Pair:              "BTC-USDT",
		ID:                "2",
		Side:              "bid",
		Quantity:          "0.3",
		Price:             "50500",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	matchingOrder1Bytes, _ := json.Marshal(matchingOrder1)

	matchingOrder2 := Order{
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

	matchingOrder3 := Order{
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
	obp := NewRedisOrderBookProvider(rc)

	ob := newOrderBook("BTC-USDT", obp)
	o := Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "ask",
		Quantity:          "0.3",
		Price:             "",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	doneOrders, partial, err := ob.processOrder(o)
	assert.Nil(t, err)
	assert.NotNil(t, partial)
	assert.Equal(t, "2", partial.ID)
	assert.Equal(t, "50500", partial.Price)
	assert.Equal(t, "0.2", partial.Quantity)
	assert.Equal(t, "bid", partial.Side)

	assert.Equal(t, 6, len(doneOrders))
	currentOrder := doneOrders[0]
	assert.Equal(t, "1", currentOrder.ID)
	assert.Equal(t, "ask", currentOrder.Side)
	assert.Equal(t, "", currentOrder.Price)
	assert.Equal(t, "50700", currentOrder.TradePrice)
	assert.Equal(t, "4", currentOrder.TradedWithOrderID)
	assert.Equal(t, "0.1", currentOrder.QuantityTraded)
	assert.Equal(t, "0.3", currentOrder.Quantity)

	doneOrder1 := doneOrders[1]
	assert.Equal(t, "4", doneOrder1.ID)
	assert.Equal(t, "bid", doneOrder1.Side)
	assert.Equal(t, "50700", doneOrder1.Price)
	assert.Equal(t, "50700", doneOrder1.TradePrice)
	assert.Equal(t, "1", doneOrder1.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder1.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder1.Quantity)

	doneOrder2 := doneOrders[2]
	assert.Equal(t, "1", doneOrder2.ID)
	assert.Equal(t, "ask", doneOrder2.Side)
	assert.Equal(t, "", doneOrder2.Price)
	assert.Equal(t, "50600", doneOrder2.TradePrice)
	assert.Equal(t, "3", doneOrder2.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder2.QuantityTraded)
	assert.Equal(t, "0.2", doneOrder2.Quantity)

	doneOrder3 := doneOrders[3]
	assert.Equal(t, "3", doneOrder3.ID)
	assert.Equal(t, "bid", doneOrder3.Side)
	assert.Equal(t, "50600", doneOrder3.Price)
	assert.Equal(t, "50600", doneOrder3.TradePrice)
	assert.Equal(t, "1", doneOrder3.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder3.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder3.Quantity)

	doneOrder4 := doneOrders[4]
	assert.Equal(t, "1", doneOrder4.ID)
	assert.Equal(t, "ask", doneOrder4.Side)
	assert.Equal(t, "", doneOrder4.Price)
	assert.Equal(t, "50500", doneOrder4.TradePrice)
	assert.Equal(t, "2", doneOrder4.TradedWithOrderID)
	assert.Equal(t, "0.1000000000000000", doneOrder4.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder4.Quantity)

	doneOrder5 := doneOrders[5]
	assert.Equal(t, "2", doneOrder5.ID)
	assert.Equal(t, "bid", doneOrder5.Side)
	assert.Equal(t, "50500", doneOrder5.Price)
	assert.Equal(t, "50500", doneOrder5.TradePrice)
	assert.Equal(t, "1", doneOrder5.TradedWithOrderID)
	assert.Equal(t, "0.1000000000000000", doneOrder5.QuantityTraded)
	assert.Equal(t, "0.3", doneOrder5.Quantity)

	rc.AssertExpectations(t)
}

func TestOrderBook_processOrder_Limit_Ask_NoMatch_OnlyPartial(t *testing.T) {
	rc := new(mocks.RedisClient)

	var data []redis.Z

	rc.On("ZRangeByScoreWithScores", mock.Anything, "order-book:bid:BTC-USDT", mock.Anything).Once().Return(data, nil)
	obp := NewRedisOrderBookProvider(rc)

	ob := newOrderBook("BTC-USDT", obp)
	o := Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "ask",
		Quantity:          "0.3",
		Price:             "50000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	doneOrders, partial, err := ob.processOrder(o)
	assert.Nil(t, err)
	assert.NotNil(t, partial)
	assert.Equal(t, "1", partial.ID)
	assert.Equal(t, "50000", partial.Price)
	assert.Equal(t, "0.3", partial.Quantity)
	assert.Equal(t, "ask", partial.Side)
	assert.Equal(t, 0, len(doneOrders))

	rc.AssertExpectations(t)
}

func TestOrderBook_processOrder_Market_Ask_NoMatch_OnlyPartial(t *testing.T) {
	rc := new(mocks.RedisClient)

	var data []redis.Z

	rc.On("ZRangeByScoreWithScores", mock.Anything, "order-book:bid:BTC-USDT", mock.Anything).Once().Return(data, nil)
	obp := NewRedisOrderBookProvider(rc)

	ob := newOrderBook("BTC-USDT", obp)
	o := Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "ask",
		Quantity:          "0.3",
		Price:             "",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	doneOrders, partial, err := ob.processOrder(o)
	assert.Nil(t, err)
	// C6: Market orders are IOC — unfilled market orders are cancelled (nil partial)
	assert.Nil(t, partial)
	assert.Equal(t, 0, len(doneOrders))

	rc.AssertExpectations(t)
}

func TestOrderBook_processOrder_Limit_Bid_MatchedWithSingleOrder_SamePrice_WithoutRemainingPartial(t *testing.T) {
	rc := new(mocks.RedisClient)
	matchingOrder := Order{
		Pair:              "BTC-USDT",
		ID:                "2",
		Side:              "ask",
		Quantity:          "0.1",
		Price:             "50000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}

	matchingOrderBytes, _ := json.Marshal(matchingOrder)
	data := []redis.Z{
		{
			Score:  1,
			Member: string(matchingOrderBytes),
		},
	}
	rc.On("ZRangeByScoreWithScores", mock.Anything, "order-book:ask:BTC-USDT", mock.Anything).Once().Return(data, nil)
	obp := NewRedisOrderBookProvider(rc)

	ob := newOrderBook("BTC-USDT", obp)
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
	doneOrders, partial, err := ob.processOrder(o)
	assert.Nil(t, err)
	assert.Nil(t, partial)
	assert.Equal(t, 2, len(doneOrders))

	currentOrder := doneOrders[0]
	assert.Equal(t, "1", currentOrder.ID)
	assert.Equal(t, "bid", currentOrder.Side)
	assert.Equal(t, "50000", currentOrder.Price)
	assert.Equal(t, "50000", currentOrder.TradePrice)
	assert.Equal(t, "2", currentOrder.TradedWithOrderID)
	assert.Equal(t, "0.1", currentOrder.QuantityTraded)

	matchingDoneOrder := doneOrders[1]
	assert.Equal(t, "2", matchingDoneOrder.ID)
	assert.Equal(t, "ask", matchingDoneOrder.Side)
	assert.Equal(t, "50000", matchingDoneOrder.Price)
	assert.Equal(t, "50000", matchingDoneOrder.TradePrice)
	assert.Equal(t, "1", matchingDoneOrder.TradedWithOrderID)
	assert.Equal(t, "0.1", matchingDoneOrder.QuantityTraded)

	rc.AssertExpectations(t)
}

func TestOrderBook_processOrder_Market_Bid_MatchedWithSingleOrder_SamePrice_WithoutRemainingPartial(t *testing.T) {
	rc := new(mocks.RedisClient)
	matchingOrder := Order{
		Pair:              "BTC-USDT",
		ID:                "2",
		Side:              "ask",
		Quantity:          "0.1",
		Price:             "50000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}

	matchingOrderBytes, _ := json.Marshal(matchingOrder)
	data := []redis.Z{
		{
			Score:  1,
			Member: string(matchingOrderBytes),
		},
	}
	rc.On("ZRangeByScoreWithScores", mock.Anything, "order-book:ask:BTC-USDT", mock.Anything).Once().Return(data, nil)
	obp := NewRedisOrderBookProvider(rc)

	ob := newOrderBook("BTC-USDT", obp)
	o := Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "bid",
		Quantity:          "0.1",
		Price:             "",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MarketPrice:       "50000.00000000",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	doneOrders, partial, err := ob.processOrder(o)
	assert.Nil(t, err)
	assert.Nil(t, partial)
	assert.Equal(t, 2, len(doneOrders))

	currentOrder := doneOrders[0]
	assert.Equal(t, "1", currentOrder.ID)
	assert.Equal(t, "bid", currentOrder.Side)
	assert.Equal(t, "", currentOrder.Price)
	assert.Equal(t, "50000", currentOrder.TradePrice)
	assert.Equal(t, "2", currentOrder.TradedWithOrderID)
	assert.Equal(t, "0.1", currentOrder.QuantityTraded)

	matchingDoneOrder := doneOrders[1]
	assert.Equal(t, "2", matchingDoneOrder.ID)
	assert.Equal(t, "ask", matchingDoneOrder.Side)
	assert.Equal(t, "50000", matchingDoneOrder.Price)
	assert.Equal(t, "50000", matchingDoneOrder.TradePrice)
	assert.Equal(t, "1", matchingDoneOrder.TradedWithOrderID)
	assert.Equal(t, "0.1", matchingDoneOrder.QuantityTraded)

	rc.AssertExpectations(t)
}

func TestOrderBook_processOrder_Limit_Bid_MatchedWithThreeOrder_SamePrice_WithoutRemainingPartial(t *testing.T) {
	rc := new(mocks.RedisClient)
	matchingOrder1 := Order{
		Pair:              "BTC-USDT",
		ID:                "2",
		Side:              "ask",
		Quantity:          "0.1",
		Price:             "50000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	matchingOrder1Bytes, _ := json.Marshal(matchingOrder1)

	matchingOrder2 := Order{
		Pair:              "BTC-USDT",
		ID:                "3",
		Side:              "ask",
		Quantity:          "0.1",
		Price:             "50000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}

	matchingOrder2Bytes, _ := json.Marshal(matchingOrder2)

	matchingOrder3 := Order{
		Pair:              "BTC-USDT",
		ID:                "4",
		Side:              "ask",
		Quantity:          "0.1",
		Price:             "50000",
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
	rc.On("ZRangeByScoreWithScores", mock.Anything, "order-book:ask:BTC-USDT", mock.Anything).Once().Return(data, nil)
	obp := NewRedisOrderBookProvider(rc)

	ob := newOrderBook("BTC-USDT", obp)
	o := Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "bid",
		Quantity:          "0.3",
		Price:             "50000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	doneOrders, partial, err := ob.processOrder(o)
	assert.Nil(t, err)
	assert.Nil(t, partial)
	assert.Equal(t, 6, len(doneOrders))

	currentOrder := doneOrders[0]
	assert.Equal(t, "1", currentOrder.ID)
	assert.Equal(t, "bid", currentOrder.Side)
	assert.Equal(t, "50000", currentOrder.Price)
	assert.Equal(t, "50000", currentOrder.TradePrice)
	assert.Equal(t, "2", currentOrder.TradedWithOrderID)
	assert.Equal(t, "0.1", currentOrder.QuantityTraded)
	assert.Equal(t, "0.3", currentOrder.Quantity)

	doneOrder1 := doneOrders[1]
	assert.Equal(t, "2", doneOrder1.ID)
	assert.Equal(t, "ask", doneOrder1.Side)
	assert.Equal(t, "50000", doneOrder1.TradePrice)
	assert.Equal(t, "1", doneOrder1.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder1.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder1.Quantity)

	doneOrder2 := doneOrders[2]
	assert.Equal(t, "1", doneOrder2.ID)
	assert.Equal(t, "bid", doneOrder2.Side)
	assert.Equal(t, "50000", doneOrder2.TradePrice)
	assert.Equal(t, "3", doneOrder2.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder2.QuantityTraded)
	assert.Equal(t, "0.2", doneOrder2.Quantity)

	doneOrder3 := doneOrders[3]
	assert.Equal(t, "3", doneOrder3.ID)
	assert.Equal(t, "ask", doneOrder3.Side)
	assert.Equal(t, "50000", doneOrder3.TradePrice)
	assert.Equal(t, "1", doneOrder3.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder3.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder3.Quantity)

	doneOrder4 := doneOrders[4]
	assert.Equal(t, "1", doneOrder4.ID)
	assert.Equal(t, "bid", doneOrder4.Side)
	assert.Equal(t, "50000", doneOrder4.TradePrice)
	assert.Equal(t, "4", doneOrder4.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder4.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder4.Quantity)

	doneOrder5 := doneOrders[5]
	assert.Equal(t, "4", doneOrder5.ID)
	assert.Equal(t, "ask", doneOrder5.Side)
	assert.Equal(t, "50000", doneOrder5.TradePrice)
	assert.Equal(t, "1", doneOrder5.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder5.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder5.Quantity)

	rc.AssertExpectations(t)
}

func TestOrderBook_processOrder_Market_Bid_MatchedWithThreeOrder_SamePrice_WithoutRemainingPartial(t *testing.T) {
	rc := new(mocks.RedisClient)
	matchingOrder1 := Order{
		Pair:              "BTC-USDT",
		ID:                "2",
		Side:              "ask",
		Quantity:          "0.1",
		Price:             "50000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	matchingOrder1Bytes, _ := json.Marshal(matchingOrder1)

	matchingOrder2 := Order{
		Pair:              "BTC-USDT",
		ID:                "3",
		Side:              "ask",
		Quantity:          "0.1",
		Price:             "50000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}

	matchingOrder2Bytes, _ := json.Marshal(matchingOrder2)

	matchingOrder3 := Order{
		Pair:              "BTC-USDT",
		ID:                "4",
		Side:              "ask",
		Quantity:          "0.1",
		Price:             "50000",
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
	rc.On("ZRangeByScoreWithScores", mock.Anything, "order-book:ask:BTC-USDT", mock.Anything).Once().Return(data, nil)
	obp := NewRedisOrderBookProvider(rc)

	ob := newOrderBook("BTC-USDT", obp)
	o := Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "bid",
		Quantity:          "0.3",
		Price:             "",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MarketPrice:       "50000.00000000",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	doneOrders, partial, err := ob.processOrder(o)
	assert.Nil(t, err)
	assert.Nil(t, partial)
	assert.Equal(t, 6, len(doneOrders))

	currentOrder := doneOrders[0]
	assert.Equal(t, "1", currentOrder.ID)
	assert.Equal(t, "bid", currentOrder.Side)
	assert.Equal(t, "", currentOrder.Price)
	assert.Equal(t, "50000", currentOrder.TradePrice)
	assert.Equal(t, "2", currentOrder.TradedWithOrderID)
	assert.Equal(t, "0.1", currentOrder.QuantityTraded)
	assert.Equal(t, "0.3", currentOrder.Quantity)

	doneOrder1 := doneOrders[1]
	assert.Equal(t, "2", doneOrder1.ID)
	assert.Equal(t, "ask", doneOrder1.Side)
	assert.Equal(t, "50000", doneOrder1.Price)
	assert.Equal(t, "50000", doneOrder1.TradePrice)
	assert.Equal(t, "1", doneOrder1.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder1.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder1.Quantity)

	doneOrder2 := doneOrders[2]
	assert.Equal(t, "1", doneOrder2.ID)
	assert.Equal(t, "bid", doneOrder2.Side)
	assert.Equal(t, "", doneOrder2.Price)
	assert.Equal(t, "50000", doneOrder2.TradePrice)
	assert.Equal(t, "3", doneOrder2.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder2.QuantityTraded)
	assert.Equal(t, "0.2000000000000000", doneOrder2.Quantity)

	doneOrder3 := doneOrders[3]
	assert.Equal(t, "3", doneOrder3.ID)
	assert.Equal(t, "ask", doneOrder3.Side)
	assert.Equal(t, "50000", doneOrder3.Price)
	assert.Equal(t, "50000", doneOrder3.TradePrice)
	assert.Equal(t, "1", doneOrder3.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder3.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder3.Quantity)

	doneOrder4 := doneOrders[4]
	assert.Equal(t, "1", doneOrder4.ID)
	assert.Equal(t, "bid", doneOrder4.Side)
	assert.Equal(t, "", doneOrder4.Price)
	assert.Equal(t, "50000", doneOrder4.TradePrice)
	assert.Equal(t, "4", doneOrder4.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder4.QuantityTraded)
	assert.Equal(t, "0.1000000000000000", doneOrder4.Quantity)

	doneOrder5 := doneOrders[5]
	assert.Equal(t, "4", doneOrder5.ID)
	assert.Equal(t, "ask", doneOrder5.Side)
	assert.Equal(t, "50000", doneOrder5.Price)
	assert.Equal(t, "50000", doneOrder5.TradePrice)
	assert.Equal(t, "1", doneOrder5.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder5.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder5.Quantity)

	rc.AssertExpectations(t)
}

func TestOrderBook_processOrder_Limit_Bid_MatchedWithThreeOrder_SamePrice_WithRemainingPartial(t *testing.T) {
	rc := new(mocks.RedisClient)
	matchingOrder1 := Order{
		Pair:              "BTC-USDT",
		ID:                "2",
		Side:              "ask",
		Quantity:          "0.1",
		Price:             "50000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	matchingOrder1Bytes, _ := json.Marshal(matchingOrder1)

	matchingOrder2 := Order{
		Pair:              "BTC-USDT",
		ID:                "3",
		Side:              "ask",
		Quantity:          "0.1",
		Price:             "50000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}

	matchingOrder2Bytes, _ := json.Marshal(matchingOrder2)

	matchingOrder3 := Order{
		Pair:              "BTC-USDT",
		ID:                "4",
		Side:              "ask",
		Quantity:          "0.1",
		Price:             "50000",
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
	rc.On("ZRangeByScoreWithScores", mock.Anything, "order-book:ask:BTC-USDT", mock.Anything).Once().Return(data, nil)
	obp := NewRedisOrderBookProvider(rc)

	ob := newOrderBook("BTC-USDT", obp)
	o := Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "bid",
		Quantity:          "0.4",
		Price:             "50000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	doneOrders, partial, err := ob.processOrder(o)
	assert.Nil(t, err)
	assert.NotNil(t, partial)
	assert.Equal(t, "1", partial.ID)
	assert.Equal(t, "50000", partial.Price)
	assert.Equal(t, "0.1", partial.Quantity)
	assert.Equal(t, "bid", partial.Side)

	assert.Equal(t, 6, len(doneOrders))

	currentOrder := doneOrders[0]
	assert.Equal(t, "1", currentOrder.ID)
	assert.Equal(t, "bid", currentOrder.Side)
	assert.Equal(t, "50000", currentOrder.Price)
	assert.Equal(t, "50000", currentOrder.TradePrice)
	assert.Equal(t, "2", currentOrder.TradedWithOrderID)
	assert.Equal(t, "0.1", currentOrder.QuantityTraded)
	assert.Equal(t, "0.4", currentOrder.Quantity)

	doneOrder1 := doneOrders[1]
	assert.Equal(t, "2", doneOrder1.ID)
	assert.Equal(t, "ask", doneOrder1.Side)
	assert.Equal(t, "50000", doneOrder1.Price)
	assert.Equal(t, "50000", doneOrder1.TradePrice)
	assert.Equal(t, "1", doneOrder1.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder1.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder1.Quantity)

	doneOrder2 := doneOrders[2]
	assert.Equal(t, "1", doneOrder2.ID)
	assert.Equal(t, "bid", doneOrder2.Side)
	assert.Equal(t, "50000", doneOrder2.Price)
	assert.Equal(t, "50000", doneOrder2.TradePrice)
	assert.Equal(t, "3", doneOrder2.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder2.QuantityTraded)
	assert.Equal(t, "0.3", doneOrder2.Quantity)

	doneOrder3 := doneOrders[3]
	assert.Equal(t, "3", doneOrder3.ID)
	assert.Equal(t, "ask", doneOrder3.Side)
	assert.Equal(t, "50000", doneOrder3.Price)
	assert.Equal(t, "50000", doneOrder3.TradePrice)
	assert.Equal(t, "1", doneOrder3.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder3.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder3.Quantity)

	doneOrder4 := doneOrders[4]
	assert.Equal(t, "1", doneOrder4.ID)
	assert.Equal(t, "bid", doneOrder4.Side)
	assert.Equal(t, "50000", doneOrder4.Price)
	assert.Equal(t, "50000", doneOrder4.TradePrice)
	assert.Equal(t, "4", doneOrder4.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder4.QuantityTraded)
	assert.Equal(t, "0.2", doneOrder4.Quantity)

	doneOrder5 := doneOrders[5]
	assert.Equal(t, "4", doneOrder5.ID)
	assert.Equal(t, "ask", doneOrder5.Side)
	assert.Equal(t, "50000", doneOrder5.TradePrice)
	assert.Equal(t, "50000", doneOrder5.Price)
	assert.Equal(t, "1", doneOrder5.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder5.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder5.Quantity)

	rc.AssertExpectations(t)
}

func TestOrderBook_processOrder_Market_Bid_MatchedWithThreeOrder_SamePrice_WithRemainingPartial(t *testing.T) {
	rc := new(mocks.RedisClient)
	matchingOrder1 := Order{
		Pair:              "BTC-USDT",
		ID:                "2",
		Side:              "ask",
		Quantity:          "0.1",
		Price:             "50000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	matchingOrder1Bytes, _ := json.Marshal(matchingOrder1)

	matchingOrder2 := Order{
		Pair:              "BTC-USDT",
		ID:                "3",
		Side:              "ask",
		Quantity:          "0.1",
		Price:             "50000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}

	matchingOrder2Bytes, _ := json.Marshal(matchingOrder2)

	matchingOrder3 := Order{
		Pair:              "BTC-USDT",
		ID:                "4",
		Side:              "ask",
		Quantity:          "0.1",
		Price:             "50000",
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
	zrangeBy := &redis.ZRangeBy{Min: "0", Max: "50500", Offset: 0, Count: 10000}
	rc.On("ZRangeByScoreWithScores", mock.Anything, "order-book:ask:BTC-USDT", zrangeBy).Once().Return(data, nil)
	obp := NewRedisOrderBookProvider(rc)

	ob := newOrderBook("BTC-USDT", obp)
	o := Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "bid",
		Quantity:          "0.4",
		Price:             "",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MarketPrice:       "50000.00000000",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	doneOrders, partial, err := ob.processOrder(o)
	assert.Nil(t, err)
	assert.NotNil(t, partial)
	assert.Equal(t, "1", partial.ID)
	assert.Equal(t, "", partial.Price)
	assert.Equal(t, "0.1000000000000000", partial.Quantity)
	assert.Equal(t, "bid", partial.Side)

	assert.Equal(t, 6, len(doneOrders))
	currentOrder := doneOrders[0]
	assert.Equal(t, "1", currentOrder.ID)
	assert.Equal(t, "bid", currentOrder.Side)
	assert.Equal(t, "", currentOrder.Price)
	assert.Equal(t, "50000", currentOrder.TradePrice)
	assert.Equal(t, "2", currentOrder.TradedWithOrderID)
	assert.Equal(t, "0.1", currentOrder.QuantityTraded)
	assert.Equal(t, "0.4", currentOrder.Quantity)

	doneOrder1 := doneOrders[1]
	assert.Equal(t, "2", doneOrder1.ID)
	assert.Equal(t, "ask", doneOrder1.Side)
	assert.Equal(t, "50000", doneOrder1.Price)
	assert.Equal(t, "50000", doneOrder1.TradePrice)
	assert.Equal(t, "1", doneOrder1.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder1.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder1.Quantity)

	doneOrder2 := doneOrders[2]
	assert.Equal(t, "1", doneOrder2.ID)
	assert.Equal(t, "bid", doneOrder2.Side)
	assert.Equal(t, "", doneOrder2.Price)
	assert.Equal(t, "50000", doneOrder2.TradePrice)
	assert.Equal(t, "3", doneOrder2.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder2.QuantityTraded)
	assert.Equal(t, "0.3000000000000000", doneOrder2.Quantity)

	doneOrder3 := doneOrders[3]
	assert.Equal(t, "3", doneOrder3.ID)
	assert.Equal(t, "ask", doneOrder3.Side)
	assert.Equal(t, "50000", doneOrder3.Price)
	assert.Equal(t, "50000", doneOrder3.TradePrice)
	assert.Equal(t, "1", doneOrder3.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder3.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder3.Quantity)

	doneOrder4 := doneOrders[4]
	assert.Equal(t, "1", doneOrder4.ID)
	assert.Equal(t, "bid", doneOrder4.Side)
	assert.Equal(t, "", doneOrder4.Price)
	assert.Equal(t, "50000", doneOrder4.TradePrice)
	assert.Equal(t, "4", doneOrder4.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder4.QuantityTraded)
	assert.Equal(t, "0.2000000000000000", doneOrder4.Quantity)

	doneOrder5 := doneOrders[5]
	assert.Equal(t, "4", doneOrder5.ID)
	assert.Equal(t, "ask", doneOrder5.Side)
	assert.Equal(t, "50000", doneOrder5.Price)
	assert.Equal(t, "50000", doneOrder5.TradePrice)
	assert.Equal(t, "1", doneOrder5.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder5.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder5.Quantity)

	rc.AssertExpectations(t)
}

func TestOrderBook_processOrder_Limit_Bid_MatchedWithThreeOrder_DifferentPrice_WithRemainingPartial(t *testing.T) {
	rc := new(mocks.RedisClient)
	matchingOrder1 := Order{
		Pair:              "BTC-USDT",
		ID:                "2",
		Side:              "ask",
		Quantity:          "0.1",
		Price:             "49500",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	matchingOrder1Bytes, _ := json.Marshal(matchingOrder1)

	matchingOrder2 := Order{
		Pair:              "BTC-USDT",
		ID:                "3",
		Side:              "ask",
		Quantity:          "0.1",
		Price:             "49600",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}

	matchingOrder2Bytes, _ := json.Marshal(matchingOrder2)

	matchingOrder3 := Order{
		Pair:              "BTC-USDT",
		ID:                "4",
		Side:              "ask",
		Quantity:          "0.1",
		Price:             "49700",
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
	rc.On("ZRangeByScoreWithScores", mock.Anything, "order-book:ask:BTC-USDT", mock.Anything).Once().Return(data, nil)
	obp := NewRedisOrderBookProvider(rc)

	ob := newOrderBook("BTC-USDT", obp)
	o := Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "bid",
		Quantity:          "0.4",
		Price:             "50000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	doneOrders, partial, err := ob.processOrder(o)
	assert.Nil(t, err)
	assert.NotNil(t, partial)
	assert.Equal(t, "1", partial.ID)
	assert.Equal(t, "50000", partial.Price)
	assert.Equal(t, "0.1", partial.Quantity)
	assert.Equal(t, "bid", partial.Side)

	assert.Equal(t, 6, len(doneOrders))
	currentOrder := doneOrders[0]
	assert.Equal(t, "1", currentOrder.ID)
	assert.Equal(t, "bid", currentOrder.Side)
	assert.Equal(t, "50000", currentOrder.Price)
	assert.Equal(t, "50000", currentOrder.TradePrice)
	assert.Equal(t, "2", currentOrder.TradedWithOrderID)
	assert.Equal(t, "0.1", currentOrder.QuantityTraded)
	assert.Equal(t, "0.4", currentOrder.Quantity)

	doneOrder1 := doneOrders[1]
	assert.Equal(t, "2", doneOrder1.ID)
	assert.Equal(t, "ask", doneOrder1.Side)
	assert.Equal(t, "49500", doneOrder1.Price)
	assert.Equal(t, "50000", doneOrder1.TradePrice)
	assert.Equal(t, "1", doneOrder1.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder1.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder1.Quantity)

	doneOrder2 := doneOrders[2]
	assert.Equal(t, "1", doneOrder2.ID)
	assert.Equal(t, "bid", doneOrder2.Side)
	assert.Equal(t, "50000", doneOrder2.Price)
	assert.Equal(t, "50000", doneOrder2.TradePrice)
	assert.Equal(t, "3", doneOrder2.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder2.QuantityTraded)
	assert.Equal(t, "0.3", doneOrder2.Quantity)

	doneOrder3 := doneOrders[3]
	assert.Equal(t, "3", doneOrder3.ID)
	assert.Equal(t, "ask", doneOrder3.Side)
	assert.Equal(t, "49600", doneOrder3.Price)
	assert.Equal(t, "50000", doneOrder3.TradePrice)
	assert.Equal(t, "1", doneOrder3.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder3.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder3.Quantity)

	doneOrder4 := doneOrders[4]
	assert.Equal(t, "1", doneOrder4.ID)
	assert.Equal(t, "bid", doneOrder4.Side)
	assert.Equal(t, "50000", doneOrder4.Price)
	assert.Equal(t, "50000", doneOrder4.TradePrice)
	assert.Equal(t, "4", doneOrder4.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder4.QuantityTraded)
	assert.Equal(t, "0.2", doneOrder4.Quantity)

	doneOrder5 := doneOrders[5]
	assert.Equal(t, "4", doneOrder5.ID)
	assert.Equal(t, "ask", doneOrder5.Side)
	assert.Equal(t, "49700", doneOrder5.Price)
	assert.Equal(t, "50000", doneOrder5.TradePrice)
	assert.Equal(t, "1", doneOrder5.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder5.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder5.Quantity)

	rc.AssertExpectations(t)
}

func TestOrderBook_processOrder_Market_Bid_MatchedWithThreeOrder_DifferentPrice_WithRemainingPartial(t *testing.T) {
	rc := new(mocks.RedisClient)
	matchingOrder1 := Order{
		Pair:              "BTC-USDT",
		ID:                "2",
		Side:              "ask",
		Quantity:          "0.1",
		Price:             "10000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	matchingOrder1Bytes, _ := json.Marshal(matchingOrder1)

	matchingOrder2 := Order{
		Pair:              "BTC-USDT",
		ID:                "3",
		Side:              "ask",
		Quantity:          "0.1",
		Price:             "20000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}

	matchingOrder2Bytes, _ := json.Marshal(matchingOrder2)

	//matchingOrder3 := Order{
	//Pair:                     "BTC-USDT",
	//ID:                       "4",
	//Side:                     "ask",
	//Quantity:                 "0.1",
	//Price:                    "49700",
	//Timestamp:                time.Now().Unix(),
	//TradedWithOrderID:        "",
	//QuantityTraded:           "",
	//TradePrice:               "",
	//MinThresholdPrice:        "49500",
	//MaxThresholdPrice:        "50500",
	//SupportsExternalExchange: true,
	//}

	//matchingOrder3Bytes, _ := json.Marshal(matchingOrder3)

	data := []redis.Z{
		{
			Score:  10000,
			Member: string(matchingOrder1Bytes),
		},
		{
			Score:  20000,
			Member: string(matchingOrder2Bytes),
		},
		//{
		//Score:  50000,
		//Member: string(matchingOrder3Bytes),
		//},
	}
	zrangeBy := &redis.ZRangeBy{Min: "0", Max: "50500", Offset: 0, Count: 10000}
	rc.On("ZRangeByScoreWithScores", mock.Anything, "order-book:ask:BTC-USDT", zrangeBy).Once().Return(data, nil)
	obp := NewRedisOrderBookProvider(rc)

	ob := newOrderBook("BTC-USDT", obp)
	o := Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "bid",
		Quantity:          "0.4",
		Price:             "",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		MarketPrice:       "10000.00000000",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	doneOrders, partial, err := ob.processOrder(o)
	assert.Nil(t, err)
	assert.NotNil(t, partial)
	assert.Equal(t, "1", partial.ID)
	assert.Equal(t, "", partial.Price)
	assert.Equal(t, "0.1000000000000000", partial.Quantity)
	assert.Equal(t, "bid", partial.Side)

	assert.Equal(t, 4, len(doneOrders))
	currentOrder := doneOrders[0]
	assert.Equal(t, "1", currentOrder.ID)
	assert.Equal(t, "bid", currentOrder.Side)
	assert.Equal(t, "", currentOrder.Price)
	assert.Equal(t, "10000", currentOrder.TradePrice)
	assert.Equal(t, "2", currentOrder.TradedWithOrderID)
	assert.Equal(t, "0.1", currentOrder.QuantityTraded)
	assert.Equal(t, "0.4", currentOrder.Quantity)

	doneOrder1 := doneOrders[1]
	assert.Equal(t, "2", doneOrder1.ID)
	assert.Equal(t, "ask", doneOrder1.Side)
	assert.Equal(t, "10000", doneOrder1.Price)
	assert.Equal(t, "10000", doneOrder1.TradePrice)
	assert.Equal(t, "1", doneOrder1.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder1.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder1.Quantity)

	doneOrder2 := doneOrders[2]
	assert.Equal(t, "1", doneOrder2.ID)
	assert.Equal(t, "bid", doneOrder2.Side)
	assert.Equal(t, "", doneOrder2.Price)
	assert.Equal(t, "20000", doneOrder2.TradePrice)
	assert.Equal(t, "3", doneOrder2.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder2.QuantityTraded)
	assert.Equal(t, "0.3000000000000000", doneOrder2.Quantity)

	doneOrder3 := doneOrders[3]
	assert.Equal(t, "3", doneOrder3.ID)
	assert.Equal(t, "ask", doneOrder3.Side)
	assert.Equal(t, "20000", doneOrder3.Price)
	assert.Equal(t, "20000", doneOrder3.TradePrice)
	assert.Equal(t, "1", doneOrder3.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder3.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder3.Quantity)
	rc.AssertExpectations(t)
}

func TestOrderBook_processOrder_Limit_Bid_MatchedWithThreeOrder_DifferentPrice_WithRemainingPartialForMakerOrder(t *testing.T) {
	rc := new(mocks.RedisClient)
	matchingOrder1 := Order{
		Pair:              "BTC-USDT",
		ID:                "2",
		Side:              "ask",
		Quantity:          "0.1",
		Price:             "49500",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	matchingOrder1Bytes, _ := json.Marshal(matchingOrder1)

	matchingOrder2 := Order{
		Pair:              "BTC-USDT",
		ID:                "3",
		Side:              "ask",
		Quantity:          "0.1",
		Price:             "49600",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}

	matchingOrder2Bytes, _ := json.Marshal(matchingOrder2)

	matchingOrder3 := Order{
		Pair:              "BTC-USDT",
		ID:                "4",
		Side:              "ask",
		Quantity:          "0.3",
		Price:             "49700",
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
	zrangeBy := &redis.ZRangeBy{Min: "0", Max: "50000", Offset: 0, Count: 10000}
	rc.On("ZRangeByScoreWithScores", mock.Anything, "order-book:ask:BTC-USDT", zrangeBy).Once().Return(data, nil)
	obp := NewRedisOrderBookProvider(rc)

	ob := newOrderBook("BTC-USDT", obp)
	o := Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "bid",
		Quantity:          "0.3",
		Price:             "50000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	doneOrders, partial, err := ob.processOrder(o)
	assert.Nil(t, err)
	assert.NotNil(t, partial)
	assert.Equal(t, "4", partial.ID)
	assert.Equal(t, "49700", partial.Price)
	assert.Equal(t, "0.2", partial.Quantity)
	assert.Equal(t, "ask", partial.Side)

	assert.Equal(t, 6, len(doneOrders))
	currentOrder := doneOrders[0]
	assert.Equal(t, "1", currentOrder.ID)
	assert.Equal(t, "bid", currentOrder.Side)
	assert.Equal(t, "50000", currentOrder.Price)
	assert.Equal(t, "50000", currentOrder.TradePrice)
	assert.Equal(t, "2", currentOrder.TradedWithOrderID)
	assert.Equal(t, "0.1", currentOrder.QuantityTraded)
	assert.Equal(t, "0.3", currentOrder.Quantity)

	doneOrder1 := doneOrders[1]
	assert.Equal(t, "2", doneOrder1.ID)
	assert.Equal(t, "ask", doneOrder1.Side)
	assert.Equal(t, "49500", doneOrder1.Price)
	assert.Equal(t, "50000", doneOrder1.TradePrice)
	assert.Equal(t, "1", doneOrder1.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder1.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder1.Quantity)

	doneOrder2 := doneOrders[2]
	assert.Equal(t, "1", doneOrder2.ID)
	assert.Equal(t, "bid", doneOrder2.Side)
	assert.Equal(t, "50000", doneOrder2.Price)
	assert.Equal(t, "50000", doneOrder2.TradePrice)
	assert.Equal(t, "3", doneOrder2.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder2.QuantityTraded)
	assert.Equal(t, "0.2", doneOrder2.Quantity)

	doneOrder3 := doneOrders[3]
	assert.Equal(t, "3", doneOrder3.ID)
	assert.Equal(t, "ask", doneOrder3.Side)
	assert.Equal(t, "49600", doneOrder3.Price)
	assert.Equal(t, "50000", doneOrder3.TradePrice)
	assert.Equal(t, "1", doneOrder3.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder3.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder3.Quantity)

	doneOrder4 := doneOrders[4]
	assert.Equal(t, "1", doneOrder4.ID)
	assert.Equal(t, "bid", doneOrder4.Side)
	assert.Equal(t, "50000", doneOrder4.Price)
	assert.Equal(t, "50000", doneOrder4.TradePrice)
	assert.Equal(t, "4", doneOrder4.TradedWithOrderID)
	assert.Equal(t, "0.1000000000000000", doneOrder4.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder4.Quantity)

	doneOrder5 := doneOrders[5]
	assert.Equal(t, "4", doneOrder5.ID)
	assert.Equal(t, "ask", doneOrder5.Side)
	assert.Equal(t, "49700", doneOrder5.Price)
	assert.Equal(t, "50000", doneOrder5.TradePrice)
	assert.Equal(t, "1", doneOrder5.TradedWithOrderID)
	assert.Equal(t, "0.1000000000000000", doneOrder5.QuantityTraded)
	assert.Equal(t, "0.3", doneOrder5.Quantity)

	rc.AssertExpectations(t)
}

func TestOrderBook_processOrder_Market_Bid_MatchedWithThreeOrder_DifferentPrice_WithRemainingPartialForMakerOrder(t *testing.T) {
	rc := new(mocks.RedisClient)
	matchingOrder1 := Order{
		Pair:              "BTC-USDT",
		ID:                "2",
		Side:              "ask",
		Quantity:          "0.1",
		Price:             "10000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	matchingOrder1Bytes, _ := json.Marshal(matchingOrder1)

	matchingOrder2 := Order{
		Pair:              "BTC-USDT",
		ID:                "3",
		Side:              "ask",
		Quantity:          "0.1",
		Price:             "20000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}

	matchingOrder2Bytes, _ := json.Marshal(matchingOrder2)

	matchingOrder3 := Order{
		Pair:              "BTC-USDT",
		ID:                "4",
		Side:              "ask",
		Quantity:          "0.3",
		Price:             "30000",
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
			Score:  10000,
			Member: string(matchingOrder1Bytes),
		},
		{
			Score:  20000,
			Member: string(matchingOrder2Bytes),
		},
		{
			Score:  30000,
			Member: string(matchingOrder3Bytes),
		},
	}
	zrangeBy := &redis.ZRangeBy{Min: "0", Max: "50500", Offset: 0, Count: 10000}
	rc.On("ZRangeByScoreWithScores", mock.Anything, "order-book:ask:BTC-USDT", zrangeBy).Once().Return(data, nil)
	obp := NewRedisOrderBookProvider(rc)

	ob := newOrderBook("BTC-USDT", obp)
	o := Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "bid",
		Quantity:          "0.4",
		Price:             "",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MarketPrice:       "10000.00000000",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	doneOrders, partial, err := ob.processOrder(o)
	assert.Nil(t, err)
	assert.NotNil(t, partial)
	assert.Equal(t, "4", partial.ID)
	assert.Equal(t, "30000", partial.Price)
	assert.Equal(t, "0.2666666666666667", partial.Quantity)
	assert.Equal(t, "ask", partial.Side)

	assert.Equal(t, 6, len(doneOrders))
	currentOrder := doneOrders[0]
	assert.Equal(t, "1", currentOrder.ID)
	assert.Equal(t, "bid", currentOrder.Side)
	assert.Equal(t, "", currentOrder.Price)
	assert.Equal(t, "10000", currentOrder.TradePrice)
	assert.Equal(t, "2", currentOrder.TradedWithOrderID)
	assert.Equal(t, "0.1", currentOrder.QuantityTraded)
	assert.Equal(t, "0.4", currentOrder.Quantity)

	doneOrder1 := doneOrders[1]
	assert.Equal(t, "2", doneOrder1.ID)
	assert.Equal(t, "ask", doneOrder1.Side)
	assert.Equal(t, "10000", doneOrder1.Price)
	assert.Equal(t, "10000", doneOrder1.TradePrice)
	assert.Equal(t, "1", doneOrder1.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder1.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder1.Quantity)

	doneOrder2 := doneOrders[2]
	assert.Equal(t, "1", doneOrder2.ID)
	assert.Equal(t, "bid", doneOrder2.Side)
	assert.Equal(t, "", doneOrder2.Price)
	assert.Equal(t, "20000", doneOrder2.TradePrice)
	assert.Equal(t, "3", doneOrder2.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder2.QuantityTraded)
	assert.Equal(t, "0.3000000000000000", doneOrder2.Quantity)

	doneOrder3 := doneOrders[3]
	assert.Equal(t, "3", doneOrder3.ID)
	assert.Equal(t, "ask", doneOrder3.Side)
	assert.Equal(t, "20000", doneOrder3.Price)
	assert.Equal(t, "20000", doneOrder3.TradePrice)
	assert.Equal(t, "1", doneOrder3.TradedWithOrderID)
	assert.Equal(t, "0.1", doneOrder3.QuantityTraded)
	assert.Equal(t, "0.1", doneOrder3.Quantity)

	doneOrder4 := doneOrders[4]
	assert.Equal(t, "1", doneOrder4.ID)
	assert.Equal(t, "bid", doneOrder4.Side)
	assert.Equal(t, "", doneOrder4.Price)
	assert.Equal(t, "30000", doneOrder4.TradePrice)
	assert.Equal(t, "4", doneOrder4.TradedWithOrderID)
	assert.Equal(t, "0.0333333333333333", doneOrder4.QuantityTraded)
	assert.Equal(t, "0.1000000000000000", doneOrder4.Quantity)

	doneOrder5 := doneOrders[5]
	assert.Equal(t, "4", doneOrder5.ID)
	assert.Equal(t, "ask", doneOrder5.Side)
	assert.Equal(t, "30000", doneOrder5.Price)
	assert.Equal(t, "30000", doneOrder5.TradePrice)
	assert.Equal(t, "1", doneOrder5.TradedWithOrderID)
	assert.Equal(t, "0.0333333333333333", doneOrder5.QuantityTraded)
	assert.Equal(t, "0.3", doneOrder5.Quantity)

	rc.AssertExpectations(t)
}

func TestOrderBook_processOrder_Limit_Bid_NoMatch_OnlyPartial(t *testing.T) {
	rc := new(mocks.RedisClient)

	var data []redis.Z

	rc.On("ZRangeByScoreWithScores", mock.Anything, "order-book:ask:BTC-USDT", mock.Anything).Once().Return(data, nil)
	obp := NewRedisOrderBookProvider(rc)

	ob := newOrderBook("BTC-USDT", obp)
	o := Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "bid",
		Quantity:          "0.3",
		Price:             "50000",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	doneOrders, partial, err := ob.processOrder(o)
	assert.Nil(t, err)
	assert.NotNil(t, partial)
	assert.Equal(t, "1", partial.ID)
	assert.Equal(t, "50000", partial.Price)
	assert.Equal(t, "0.3", partial.Quantity)
	assert.Equal(t, "bid", partial.Side)
	assert.Equal(t, 0, len(doneOrders))

	rc.AssertExpectations(t)
}

func TestOrderBook_processOrder_Market_Bid_NoMatch_OnlyPartial(t *testing.T) {
	rc := new(mocks.RedisClient)

	var data []redis.Z

	rc.On("ZRangeByScoreWithScores", mock.Anything, "order-book:ask:BTC-USDT", mock.Anything).Once().Return(data, nil)
	obp := NewRedisOrderBookProvider(rc)

	ob := newOrderBook("BTC-USDT", obp)
	o := Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "bid",
		Quantity:          "0.3",
		Price:             "",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	doneOrders, partial, err := ob.processOrder(o)
	assert.Nil(t, err)
	// C6: Market orders are IOC — unfilled market orders are cancelled (nil partial)
	assert.Nil(t, partial)
	assert.Equal(t, 0, len(doneOrders))

	rc.AssertExpectations(t)
}

func TestOrderBook_removeOrder(t *testing.T) {
	rc := new(mocks.RedisClient)
	obp := NewRedisOrderBookProvider(rc)

	ob := newOrderBook("BTC-USDT", obp)
	o := Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "bid",
		Quantity:          "0.3",
		Price:             "",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	res, err := o.MarshalForOrderbook()
	assert.Nil(t, err)
	rc.On("ZRem", mock.Anything, "order-book:bid:BTC-USDT", string(res)).Once().Return(true, nil)
	err = ob.removeOrder(o)
	assert.Nil(t, err)
	rc.AssertExpectations(t)
}

func TestOrderBook_orderExists(t *testing.T) {
	rc := new(mocks.RedisClient)
	obp := NewRedisOrderBookProvider(rc)

	ob := newOrderBook("BTC-USDT", obp)
	o := Order{
		Pair:              "BTC-USDT",
		ID:                "1",
		Side:              "bid",
		Quantity:          "0.3",
		Price:             "",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: "",
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49500",
		MaxThresholdPrice: "50500",
	}
	res, err := o.MarshalForOrderbook()
	assert.Nil(t, err)
	rc.On("ZScore", mock.Anything, "order-book:bid:BTC-USDT", string(res)).Once().Return(float64(1), nil)
	ctx := context.Background()
	exists, err := ob.orderExists(ctx, o)
	assert.Nil(t, err)
	assert.Equal(t, true, exists)

	rc.AssertExpectations(t)
}
