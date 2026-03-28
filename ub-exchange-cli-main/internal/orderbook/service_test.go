// Package orderbook_test tests the order book service. Covers:
//   - Updating external order book from WebSocket depth data with precision formatting
//   - Retrieving the full order book with bid/ask amounts, values, sums, and percentages
//   - Fetching the trade book with recent trade details (price, amount, maker flag)
//
// Test data: mock HTTP client, live data service, currency service, and pre-built
// depth snapshot and trade fixtures for BTC-USDT pair.
package orderbook_test

import (
	"context"
	"exchange-go/internal/currency"
	"exchange-go/internal/livedata"
	"exchange-go/internal/mocks"
	"exchange-go/internal/orderbook"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_UpdateExternalOrderBook(t *testing.T) {
	httpClient := new(mocks.HttpClient)
	//bodyString := "{\"lastUpdateId\":155,\"bids\":[[\"32725.19000000\",\"0.55989300\"],[\"32723.57000000\",\"1.85000000\"]],\"asks\":[[\"32725.19000000\",\"0.55989300\"],[\"32723.57000000\",\"1.85000000\"]]}"
	//httpClient.On("HttpGet", mock.Anything, mock.Anything).Once().Return([]byte(bodyString), http.Header{}, http.StatusOK, nil)
	liveDataService := new(mocks.LiveData)
	liveDataService.On("SetDepthSnapshot", mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)
	redisDepthSnapshot := livedata.RedisDepthSnapshot{}
	liveDataService.On("GetDepthSnapshot", mock.Anything, mock.Anything).Once().Return(redisDepthSnapshot, nil)
	currencyService := new(mocks.CurrencyService)
	logger := new(mocks.Logger)
	orderBookService := orderbook.NewOrderBookService(httpClient, liveDataService, currencyService, logger)
	ctx := context.Background()
	externalExchangeName := "binance"
	pairName := "BTC-USDT"
	externalExchangePairName := "BTCUSDT"
	precision := 6

	data := "{" +
		"\"Ub\": 156," +
		"\"u\": 158," +
		"\"B\": [[\"32725.19000000\",\"11.3120000\"]]," +
		"\"A\": [[\"32723.57000000\",\"1.86000000\"]]" +
		"}"

	res, err := orderBookService.UpdateExternalOrderBook(ctx, externalExchangeName, pairName, externalExchangePairName, precision, []byte(data))
	assert.Nil(t, err)
	assert.Equal(t, "32725.190000", res.Bids[0].Price)
	assert.Equal(t, "11.312000", res.Bids[0].Amount)
	assert.Equal(t, "32723.570000", res.Asks[0].Price)
	assert.Equal(t, "1.860000", res.Asks[0].Amount)
	httpClient.AssertExpectations(t)
	liveDataService.AssertExpectations(t)
}

func TestService_GetOrderBook(t *testing.T) {
	httpClient := new(mocks.HttpClient)
	liveDataService := new(mocks.LiveData)
	redisOrderBook := livedata.RedisDepthSnapshot{
		Bids: [][3]string{
			{"50000", "0.5", "25000"},
			{"51000", "0.5", "25500"},
		},
		Asks: [][3]string{
			{"51000", "0.5", "25500"},
			{"52000", "0.5", "26000"},
		},
	}

	liveDataService.On("GetDepthSnapshot", mock.Anything, "BTC-USDT").Once().Return(redisOrderBook, nil)
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetPairByName", "BTC-USDT").Once().Return(currency.Pair{Name: "BTC-USDT"}, nil)
	logger := new(mocks.Logger)
	orderBookService := orderbook.NewOrderBookService(httpClient, liveDataService, currencyService, logger)
	params := orderbook.GetOrderBookParams{
		Pair: "BTC-USDT",
	}
	res, statusCode := orderBookService.GetOrderBook(params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)
	orderbook, ok := res.Data.(orderbook.OrderBook)
	if !ok {
		t.Error("can not cast response to struct")
	}

	assert.Equal(t, 2, len(orderbook.Asks))
	assert.Equal(t, 2, len(orderbook.Bids))

	assert.Equal(t, "0.50000000", orderbook.Asks[0].Amount)
	assert.Equal(t, "51000.00000000", orderbook.Asks[0].Price)
	assert.Equal(t, "0.00", orderbook.Asks[0].Percentage)
	assert.Equal(t, "0.50000000", orderbook.Asks[0].Sum)
	assert.Equal(t, "ask", orderbook.Asks[0].Type)
	assert.Equal(t, "25500.00000000", orderbook.Asks[0].Value)

	assert.Equal(t, "0.50000000", orderbook.Bids[0].Amount)
	assert.Equal(t, "50000.00000000", orderbook.Bids[0].Price)
	assert.Equal(t, "0.00", orderbook.Bids[0].Percentage)
	assert.Equal(t, "0.50000000", orderbook.Bids[0].Sum)
	assert.Equal(t, "bid", orderbook.Bids[0].Type)
	assert.Equal(t, "25000.00000000", orderbook.Bids[0].Value)

	currencyService.AssertExpectations(t)
	liveDataService.AssertExpectations(t)
}

func TestService_GetTradeBook(t *testing.T) {
	httpClient := new(mocks.HttpClient)
	liveDataService := new(mocks.LiveData)
	trades := []livedata.RedisTrade{
		{
			Price:   "50000.00000000",
			Amount:  "0.5",
			IsMaker: true,
			Ignore:  false,
		},
	}
	liveDataService.On("GetTradeBook", mock.Anything, "BTC-USDT").Once().Return(trades, nil)
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetPairByName", "BTC-USDT").Once().Return(currency.Pair{Name: "BTC-USDT"}, nil)
	logger := new(mocks.Logger)
	orderBookService := orderbook.NewOrderBookService(httpClient, liveDataService, currencyService, logger)
	params := orderbook.GetTradeBookParams{
		Pair: "BTC-USDT",
	}
	res, statusCode := orderBookService.GetTradeBook(params)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Equal(t, true, res.Status)
	assert.Equal(t, "", res.Message)
	tradebook, ok := res.Data.([]livedata.RedisTrade)
	if !ok {
		t.Error("can not cast response to struct")
	}

	assert.Equal(t, 1, len(tradebook))
	assert.Equal(t, "0.5", tradebook[0].Amount)
	assert.Equal(t, "50000.00000000", tradebook[0].Price)
	assert.Equal(t, true, tradebook[0].IsMaker)
	assert.Equal(t, false, tradebook[0].Ignore)

	currencyService.AssertExpectations(t)
	liveDataService.AssertExpectations(t)
}
