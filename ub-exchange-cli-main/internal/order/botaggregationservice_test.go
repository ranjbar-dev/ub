// Package order_test tests the BotAggregationService. Covers:
//   - AddToList: adding bot trade data to the Redis aggregation list (LRem + LPush)
//   - DeleteList: removing the aggregation list for a pair from Redis
//   - GetAggregationResultForPair: table-driven aggregation of BUY-only, SELL-only,
//     and mixed BUY/SELL trades with correct net amount, price, and order type calculation
//
// Test data: mocked RedisClient with LRem/LPush/LRange/Del expectations and
// JSON-encoded BotAggregationData for pair ID 1.
package order_test

import (
	"exchange-go/internal/mocks"
	"exchange-go/internal/order"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBotAggregationService_AddToList(t *testing.T) {
	rc := new(mocks.RedisClient)

	data1 := `{"tradeId":1,"pairId":1,"robotType":"BUY","amount":"0.10000000","price":"40000.00000000","lastOrderId":1}`
	rc.On("LRem", mock.Anything, "not-calculated:trades:1", int64(0), data1).Once().Return(int64(1), nil)
	rc.On("LPush", mock.Anything, "not-calculated:trades:1", data1).Once().Return(int64(1), nil)

	botAggregationService := order.NewBotAggregationService(rc)
	data := order.BotAggregationData{
		TradeID:     1,
		PairID:      1,
		RobotType:   "BUY",
		Amount:      "0.10000000",
		Price:       "40000.00000000",
		LastOrderID: 1,
	}
	err := botAggregationService.AddToList(data)

	assert.Nil(t, err)
	rc.AssertExpectations(t)
}

func TestBotAggregationService_DeleteList(t *testing.T) {
	rc := new(mocks.RedisClient)
	rc.On("Del", mock.Anything, "not-calculated:trades:1").Once().Return(int64(1), nil)

	botAggregationService := order.NewBotAggregationService(rc)
	err := botAggregationService.DeleteList(int64(1))

	assert.Nil(t, err)
	rc.AssertExpectations(t)

}

var getAggregationResultForPairTable = []struct {
	pairID      int64
	dataInRedis []string
	result      order.AggregationResult
}{
	{
		pairID: 1,
		dataInRedis: []string{
			`{
  				"tradeId":1,
				"pairId":1,
				"robotType":"BUY",
				"amount":"0.1",
				"price":"40000",
				"lastOrderId":1,
				"userId":1
  			 }`,
		},
		result: order.AggregationResult{
			LastTradeID: 1,
			OrderIds:    []string{"1"},
			AggregatedOrderData: order.FinalAggregatedOrderData{
				Type:         "SELL",
				Amount:       "0.10000000",
				ExchangeType: "MARKET",
				Price:        "40000.00000000",
				BuyAmount:    "0.00000000",
				BuyPrice:     "0.00000000",
				SellAmount:   "0.10000000",
				SellPrice:    "40000.00000000",
			},
		},
	},
	{
		pairID: 1,
		dataInRedis: []string{
			`{
  				"tradeId":2,
				"pairId":1,
				"robotType":"SELL",
				"amount":"0.1",
				"price":"40000",
				"lastOrderId":2,
				"userId":2
  			 }`,
		},
		result: order.AggregationResult{
			LastTradeID: 2,
			OrderIds:    []string{"2"},
			AggregatedOrderData: order.FinalAggregatedOrderData{
				Type:         "BUY",
				Amount:       "0.10000000",
				ExchangeType: "MARKET",
				Price:        "40000.00000000",
				BuyAmount:    "0.10000000",
				BuyPrice:     "40000.00000000",
				SellAmount:   "0.00000000",
				SellPrice:    "0.00000000",
			},
		},
	},
	{
		pairID: 1,
		dataInRedis: []string{
			`{
  				"tradeId":1,
				"pairId":1,
				"robotType":"BUY",
				"amount":"0.1",
				"price":"40000",
				"lastOrderId":1,
				"userId":1
  			 }`,
			`{
  				"tradeId":2,
				"pairId":1,
				"robotType":"BUY",
				"amount":"0.1",
				"price":"40000",
				"lastOrderId":2,
				"userId":1
  			 }`,
			`{
  				"tradeId":3,
				"pairId":1,
				"robotType":"SELL",
				"amount":"0.05",
				"price":"40000",
				"lastOrderId":3,
				"userId":1
  			 }`,
		},
		result: order.AggregationResult{
			LastTradeID: 3,
			OrderIds:    []string{"1", "2", "3"},
			AggregatedOrderData: order.FinalAggregatedOrderData{
				Type:         "SELL",
				Amount:       "0.15000000",
				ExchangeType: "MARKET",
				Price:        "40000.00000000",
				BuyAmount:    "0.05000000",
				BuyPrice:     "40000.00000000",
				SellAmount:   "0.20000000",
				SellPrice:    "40000.00000000",
			},
		},
	},
	{
		pairID: 1,
		dataInRedis: []string{
			`{
  				"tradeId":1,
				"pairId":1,
				"robotType":"SELL",
				"amount":"0.1",
				"price":"40000",
				"lastOrderId":1,
				"userId":2
  			 }`,
			`{
  				"tradeId":2,
				"pairId":1,
				"robotType":"SELL",
				"amount":"0.1",
				"price":"40000",
				"lastOrderId":2,
				"userId":2
  			 }`,
			`{
  				"tradeId":3,
				"pairId":1,
				"robotType":"BUY",
				"amount":"0.05",
				"price":"40000",
				"lastOrderId":3,
				"userId":2
  			 }`,
		},
		result: order.AggregationResult{
			LastTradeID: 3,
			OrderIds:    []string{"1", "2", "3"},
			AggregatedOrderData: order.FinalAggregatedOrderData{
				Type:         "BUY",
				Amount:       "0.15000000",
				ExchangeType: "MARKET",
				Price:        "40000.00000000",
				BuyAmount:    "0.20000000",
				BuyPrice:     "40000.00000000",
				SellAmount:   "0.05000000",
				SellPrice:    "40000.00000000",
			},
		},
	},
}

func TestBotAggregationService_GetAggregationResultForPair(t *testing.T) {
	rc := new(mocks.RedisClient)
	botAggregationService := order.NewBotAggregationService(rc)
	for _, item := range getAggregationResultForPairTable {
		rc.On("LRange", mock.Anything, "not-calculated:trades:1", int64(0), int64(-1)).Once().Return(item.dataInRedis, nil)
		result, err := botAggregationService.GetAggregationResultForPair(item.pairID)

		assert.Nil(t, err)
		assert.Equal(t, item.result, result)

		rc.AssertExpectations(t)
	}

}
