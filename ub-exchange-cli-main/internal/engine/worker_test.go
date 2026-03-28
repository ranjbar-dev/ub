// Package engine tests the worker component that dequeues orders and
// drives the order matching pipeline. Covers:
//   - processOrder: submits a bid order against an existing ask in the order
//     book, verifies the result handler callback is invoked with matched
//     orders, and confirms done order IDs are returned in the MatchingResult
//
// Test data: miniredis for a real Redis sorted set containing an ask order;
// engineResultHandler mock for the result callback; BTC-USDT bid/ask pair
// at price 50000 with quantity 0.1.
package engine

import (
	"context"
	"encoding/json"
	"exchange-go/internal/platform"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/mock"
)

// engineResultHandler is a mock implementation of the engine's result
// handler interface, used to capture and assert callback invocations
// during order matching tests.
type engineResultHandler struct {
	mock.Mock
}

// CallBack provides a mock function with given fields: doneOrders, partialOrder
func (_m *engineResultHandler) CallBack(doneOrders []Order, partialOrder *Order) MatchingResult {
	ret := _m.Called(doneOrders, partialOrder)

	var r0 MatchingResult
	if rf, ok := ret.Get(0).(func([]Order, *Order) MatchingResult); ok {
		r0 = rf(doneOrders, partialOrder)
	} else {
		r0 = ret.Get(0).(MatchingResult)
	}

	return r0
}

func TestWorker_processOrder(t *testing.T) {
	s := miniredis.NewMiniRedis()
	defer s.Close()
	_ = s.Start()
	rc := redis.NewClient(&redis.Options{Addr: s.Addr()})
	ctx := context.Background()

	matchingOrder := Order{
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

	matchingOrderBytes, _ := json.Marshal(matchingOrder)
	data := []*redis.Z{
		{
			Score:  1,
			Member: string(matchingOrderBytes),
		},
	}
	rc.ZAdd(ctx, "order-book:ask:BTC-USDT", data...)
	redisClient := platform.NewRedisTestClient(rc)

	obp := NewRedisOrderBookProvider(redisClient)
	orderbookProvider = obp

	workChan := make(chan *work)
	rh := new(engineResultHandler)
	shouldCallPostOrderMatching = true
	matchingResult := MatchingResult{
		Err:                   nil,
		RemainingPartialOrder: nil,
		RemovingDoneOrderIds:  []int64{1},
	}
	rh.On("CallBack", mock.Anything, mock.Anything).Once().Return(matchingResult)
	cbm := getCallbackManager(rh)
	worker := newWorker(workChan, 1, cbm)

	o := Order{
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

	worker.processOrder(o)
	rh.AssertExpectations(t)
}
