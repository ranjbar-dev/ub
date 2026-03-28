// Package order_test tests the AdminOrderManager. Covers:
//   - TryToFulfillOrder for a stop order: validates kline range, generates price, submits in DB
//   - TryToFulfillOrder for a limit order: validates kline range, delegates to post-order matching and events handler
//   - TryToFulfillOrder for a market order: validates kline range, delegates to post-order matching and events handler
//
// Test data: mocked currency service, kline service with high/low price data,
// price generator, post-order matching service, stop order submission manager,
// events handler, and BTC-USDT order fixtures.
package order_test

import (
	"database/sql"
	"exchange-go/internal/currency"
	"exchange-go/internal/currency/candle"
	"exchange-go/internal/mocks"
	"exchange-go/internal/order"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAdminOrderManager_TryToFulfillOrder_StopOrder(t *testing.T) {
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetPairByID", int64(1)).Once().Return(currency.Pair{Name: "BTC-USDT"}, nil)
	klineService := new(mocks.KlineService)
	highAndLowData := candle.HighAndLowPrice{
		High: "55000.00",
		Low:  "40000.00",
	}
	klineService.On("GetHighAndLowPriceFromDateForPairByPairName", "BTC-USDT", mock.Anything).Once().Return(highAndLowData, nil)

	pg := new(mocks.PriceGenerator)
	pg.On("GetPrice", mock.Anything, "BTC-USDT").Once().Return("50000.00000000", nil)

	poms := new(mocks.PostOrderMatchingService)
	matchingResult := order.MatchingResult{
		Err:                   nil,
		RemainingPartialOrder: nil,
		RemovingDoneOrderIds:  []int64{},
	}
	poms.On("HandlePostOrderMatching", mock.Anything, mock.Anything, true).Once().Return(matchingResult)

	stopOrderSubmissionManager := new(mocks.StopOrderSubmissionManager)
	stopOrderSubmissionManager.On("SubmitOrderInDb", mock.Anything, mock.Anything, "50000.00000000").Once().Return(nil)
	eh := new(mocks.EventsHandler)
	logger := new(mocks.Logger)
	adminOrderManager := order.NewAdminOrderManager(currencyService, klineService, pg, poms, stopOrderSubmissionManager, eh, logger)
	pair := currency.Pair{ID: 1, Name: "BTC-USDT"}
	o := order.Order{
		ID:                 1,
		UserID:             1,
		Type:               "BUY",
		ExchangeType:       "LIMIT",
		Price:              sql.NullString{String: "50000.00000", Valid: true},
		Status:             "OPEN",
		DemandedAmount:     sql.NullString{String: "0.1", Valid: true},
		PayedByAmount:      sql.NullString{String: "5000.00", Valid: true},
		PairID:             1,
		Pair:               pair,
		Level:              sql.NullInt64{Int64: 1, Valid: true},
		Path:               sql.NullString{String: "1,", Valid: true},
		StopPointPrice:     sql.NullString{String: "45000.00000", Valid: true},
		CurrentMarketPrice: sql.NullString{String: "50000.00000000", Valid: true},
	}

	err := adminOrderManager.TryToFulfillOrder(o)
	assert.Nil(t, err)

	klineService.AssertExpectations(t)
	poms.AssertExpectations(t)
	stopOrderSubmissionManager.AssertExpectations(t)
	pg.AssertExpectations(t)

}

func TestAdminOrderManager_TryToFulfillOrder_LimitOrder(t *testing.T) {
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetPairByID", int64(1)).Once().Return(currency.Pair{Name: "BTC-USDT"}, nil)
	klineService := new(mocks.KlineService)
	highAndLowData := candle.HighAndLowPrice{
		High: "55000.00",
		Low:  "40000.00",
	}
	klineService.On("GetHighAndLowPriceFromDateForPairByPairName", "BTC-USDT", mock.Anything).Once().Return(highAndLowData, nil)
	pg := new(mocks.PriceGenerator)
	poms := new(mocks.PostOrderMatchingService)
	matchingResult := order.MatchingResult{
		Err:                   nil,
		RemainingPartialOrder: nil,
		RemovingDoneOrderIds:  []int64{},
	}
	poms.On("HandlePostOrderMatching", mock.Anything, mock.Anything, true).Once().Return(matchingResult)

	stopOrderSubmissionManager := new(mocks.StopOrderSubmissionManager)
	eh := new(mocks.EventsHandler)
	eh.On("HandleOrderFulfillByAdmin", mock.Anything).Once().Return()

	logger := new(mocks.Logger)
	adminOrderManager := order.NewAdminOrderManager(currencyService, klineService, pg, poms, stopOrderSubmissionManager, eh, logger)
	o := order.Order{
		ID:             1,
		UserID:         1,
		Type:           "BUY",
		ExchangeType:   "LIMIT",
		Price:          sql.NullString{String: "50000.00000", Valid: true},
		Status:         "OPEN",
		DemandedAmount: sql.NullString{String: "0.1", Valid: true},
		PayedByAmount:  sql.NullString{String: "5000.00", Valid: true},
		PairID:         1,
		Level:          sql.NullInt64{Int64: 1, Valid: true},
		Path:           sql.NullString{String: "1,", Valid: true},
	}

	err := adminOrderManager.TryToFulfillOrder(o)
	assert.Nil(t, err)

	klineService.AssertExpectations(t)
	poms.AssertExpectations(t)
	eh.AssertExpectations(t)

}

func TestAdminOrderManager_TryToFulfillOrder_marketOrder(t *testing.T) {
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetPairByID", int64(1)).Once().Return(currency.Pair{Name: "BTC-USDT"}, nil)
	klineService := new(mocks.KlineService)
	highAndLowData := candle.HighAndLowPrice{
		High: "55000.00",
		Low:  "40000.00",
	}
	klineService.On("GetHighAndLowPriceFromDateForPairByPairName", "BTC-USDT", mock.Anything).Once().Return(highAndLowData, nil)

	pg := new(mocks.PriceGenerator)
	poms := new(mocks.PostOrderMatchingService)
	matchingResult := order.MatchingResult{
		Err:                   nil,
		RemainingPartialOrder: nil,
		RemovingDoneOrderIds:  []int64{},
	}
	poms.On("HandlePostOrderMatching", mock.Anything, mock.Anything, true).Once().Return(matchingResult)

	stopOrderSubmissionManager := new(mocks.StopOrderSubmissionManager)
	eh := new(mocks.EventsHandler)
	eh.On("HandleOrderFulfillByAdmin", mock.Anything).Once().Return()

	logger := new(mocks.Logger)
	adminOrderManager := order.NewAdminOrderManager(currencyService, klineService, pg, poms, stopOrderSubmissionManager, eh, logger)
	o := order.Order{
		ID:             1,
		UserID:         1,
		Type:           "BUY",
		ExchangeType:   "MARKET",
		Status:         "OPEN",
		DemandedAmount: sql.NullString{String: "0.1", Valid: true},
		PayedByAmount:  sql.NullString{String: "5000.00", Valid: true},
		PairID:         1,
		Level:          sql.NullInt64{Int64: 1, Valid: true},
		Path:           sql.NullString{String: "1,", Valid: true},
	}

	err := adminOrderManager.TryToFulfillOrder(o)
	assert.Nil(t, err)

	klineService.AssertExpectations(t)
	poms.AssertExpectations(t)
	eh.AssertExpectations(t)

}
