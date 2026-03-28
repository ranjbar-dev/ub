// Package order_test tests the TradeEventsHandler. Covers:
//   - HandleTradesCreation: processes multiple trade events by adding BUY bot aggregation
//     data to Redis when the currency pair aggregation status is active
//
// Test data: mocked RedisClient with LRem/LPush expectations, BotAggregationService,
// mocked Configs, and TradeData fixtures for pair ID 1 with BUY bot order type.
package order_test

import (
	"database/sql"
	"exchange-go/internal/currency"
	"exchange-go/internal/mocks"
	"exchange-go/internal/order"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestTradeEventsHandler_HandleTradesCreation(t *testing.T) {
	rc := new(mocks.RedisClient)

	//data1 := `{"tradeId":1,"pairId":1,"robotType":"BUY","amount":"0.10000000","price":"40000.00000000","lastOrderId":1,"userId":1}`
	data1 := `{"tradeId":1,"pairId":1,"robotType":"BUY","amount":"0.10000000","price":"40000.00000000","lastOrderId":1}`

	//data2 := `{"tradeId":2,"pairId":1,"robotType":"BUY","amount":"0.20000000","price":"30000.00000000","lastOrderId":2,"userId":1}`
	data2 := `{"tradeId":2,"pairId":1,"robotType":"BUY","amount":"0.20000000","price":"30000.00000000","lastOrderId":2}`

	rc.On("LRem", mock.Anything, "not-calculated:trades:1", int64(0), data1).Once().Return(int64(1), nil)
	rc.On("LRem", mock.Anything, "not-calculated:trades:1", int64(0), data2).Once().Return(int64(1), nil)
	rc.On("LPush", mock.Anything, "not-calculated:trades:1", data1).Once().Return(int64(1), nil)
	rc.On("LPush", mock.Anything, "not-calculated:trades:1", data2).Once().Return(int64(1), nil)

	bas := order.NewBotAggregationService(rc)
	configs := new(mocks.Configs)
	logger := new(mocks.Logger)
	teh := order.NewTradeEventsHandler(bas, configs, logger)

	tradesData := []order.TradeData{
		{
			Trade: order.Trade{
				ID:           1,
				PairID:       1,
				BotOrderType: sql.NullString{String: "BUY", Valid: true},
				Amount:       sql.NullString{String: "0.10000000", Valid: true},
				Price:        sql.NullString{String: "40000.00000000", Valid: true},
				SellOrderID:  sql.NullInt64{Int64: 1, Valid: true},
			},
			UserEmail: "test#gmail.com",
			UserID:    1,
		},
		{
			Trade: order.Trade{
				ID:           2,
				PairID:       1,
				BotOrderType: sql.NullString{String: "BUY", Valid: true},
				Amount:       sql.NullString{String: "0.20000000", Valid: true},
				Price:        sql.NullString{String: "30000.00000000", Valid: true},
				SellOrderID:  sql.NullInt64{Int64: 2, Valid: true},
			},
			UserEmail: "test#gmail.com",
			UserID:    1,
		},
	}

	pair := currency.Pair{
		AggregationStatus: currency.AggregationStatusRun,
	}

	teh.HandleTradesCreation(tradesData, pair)
}
