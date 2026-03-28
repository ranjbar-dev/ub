// Package order_test tests the EngineResultHandler. Covers:
//   - CallBack conversion from engine done orders and partial order into a MatchingResult
//   - Correct mapping of order fields (pair, side, quantity, price, threshold prices)
//   - Delegation to PostOrderMatchingService for post-match processing
//
// Test data: mocked PostOrderMatchingService, engine.Order structs, and
// CallBackOrderData with BTC-USDT pair, BUY type, and threshold prices.
package order_test

import (
	"exchange-go/internal/engine"
	"exchange-go/internal/mocks"
	"exchange-go/internal/order"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestEngineResultHandler_CallBack(t *testing.T) {
	poms := new(mocks.PostOrderMatchingService)
	callBackOrderData := order.CallBackOrderData{
		ID:                1,
		PairName:          "BTC-USDT",
		OrderType:         "BUY",
		Quantity:          "0.1",
		Price:             "",
		Timestamp:         time.Now().Unix(),
		TradedWithOrderID: 0,
		QuantityTraded:    "",
		TradePrice:        "",
		MinThresholdPrice: "49950",
		MaxThresholdPrice: "50050",
	}
	result := order.MatchingResult{
		Err:                   nil,
		RemainingPartialOrder: &callBackOrderData,
		RemovingDoneOrderIds:  []int64{},
	}
	poms.On("HandlePostOrderMatching", mock.Anything, mock.Anything, mock.Anything).Once().Return(result)
	erh := order.NewEngineResultHandler(poms)
	doneOrders := []engine.Order{
		{},
	}

	partialOrder := &engine.Order{}
	matchingResult := erh.CallBack(doneOrders, partialOrder)
	assert.Nil(t, matchingResult.Err)
	assert.Equal(t, "BTC-USDT", matchingResult.RemainingPartialOrder.Pair)
	assert.Equal(t, "00000000001", matchingResult.RemainingPartialOrder.ID)
	assert.Equal(t, "bid", matchingResult.RemainingPartialOrder.Side)
	assert.Equal(t, "0.1", matchingResult.RemainingPartialOrder.Quantity)
	assert.Equal(t, "", matchingResult.RemainingPartialOrder.Price)
	assert.Equal(t, "", matchingResult.RemainingPartialOrder.TradedWithOrderID)
	assert.Equal(t, "", matchingResult.RemainingPartialOrder.QuantityTraded)
	assert.Equal(t, "", matchingResult.RemainingPartialOrder.TradePrice)
	assert.Equal(t, "49950", matchingResult.RemainingPartialOrder.MinThresholdPrice)
	assert.Equal(t, "50050", matchingResult.RemainingPartialOrder.MaxThresholdPrice)
	poms.AssertExpectations(t)
}
