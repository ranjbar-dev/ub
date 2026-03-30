// Package processor_test tests the data processor that handles incoming market events. Covers:
//   - Processing trade events (updating trade book and publishing via Centrifugo)
//   - Processing depth events (updating external order book and publishing order book via Centrifugo)
//   - Processing kline events (publishing candlestick data via Centrifugo and queue manager)
//   - Processing ticker events (updating live price data, publishing tickers,
//     triggering stop-order submissions and in-queue order handling)
//
// Test data: mock live data, price generator, kline service, order book service,
// Centrifugo manager, queue manager, stop-order submission manager, in-queue order manager,
// currency service, and Redis client with BTC-USDT pair and market data fixtures.
package processor_test

import (
	"context"
	"exchange-go/internal/currency"
	"exchange-go/internal/mocks"
	"exchange-go/internal/orderbook"
	"exchange-go/internal/processor"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

func TestProcessor_ProcessTrade(t *testing.T) {
	redisClient := new(mocks.RedisClient)
	liveData := new(mocks.LiveData)
	liveData.On("UpdateTradeBook", mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)
	priceGenerator := new(mocks.PriceGenerator)
	klineService := new(mocks.KlineService)
	orderbookService := new(mocks.OrderbookService)
	mqttManager := new(mocks.CentrifugoManager)
	mqttManager.On("PublishTrades", mock.Anything, mock.Anything, mock.Anything).Once().Return()
	inQueueOrderManager := new(mocks.InQueueOrderManager)
	queueManager := new(mocks.QueueManager)
	stopOrderSubmissionManager := new(mocks.StopOrderSubmissionManager)

	logger := new(mocks.Logger)
	currencyService := new(mocks.CurrencyService)

	dataProcessor := processor.NewProcessor(redisClient, liveData, priceGenerator, klineService, orderbookService,
		mqttManager, stopOrderSubmissionManager, inQueueOrderManager, queueManager, logger, currencyService)
	ctx := context.Background()
	trade := processor.Trade{
		Pair:      "BTC-USDT",
		Price:     "32000",
		Amount:    "0.1",
		CreatedAt: "2021-01-12 20:20:20",
		IsMaker:   true,
		Ignore:    true,
	}
	dataProcessor.ProcessTrade(ctx, trade)
	time.Sleep(50 * time.Millisecond)
	liveData.AssertExpectations(t)
	mqttManager.AssertExpectations(t)

}

func TestProcessor_ProcessDepth(t *testing.T) {
	redisClient := new(mocks.RedisClient)
	liveData := new(mocks.LiveData)
	priceGenerator := new(mocks.PriceGenerator)
	klineService := new(mocks.KlineService)
	orderbookService := new(mocks.OrderbookService)
	ob := orderbook.OrderBook{
		Asks: []orderbook.BookItem{
			{
				Price:      "32100",
				Amount:     "0.01",
				Value:      "320",
				Percentage: "1.5",
				Sum:        "321",
				Type:       "ask",
			},
		},
		Bids: []orderbook.BookItem{
			{
				Price:      "32100",
				Amount:     "0.01",
				Value:      "320",
				Percentage: "1.5",
				Sum:        "321",
				Type:       "bid",
			},
		},
	}

	orderbookService.On("UpdateExternalOrderBook", mock.Anything, mock.Anything, mock.Anything,
		mock.Anything, mock.Anything, mock.Anything).Once().Return(ob, nil)

	mqttManager := new(mocks.CentrifugoManager)
	mqttManager.On("PublishOrderBook", mock.Anything, mock.Anything, mock.Anything).Once().Return()

	inQueueOrderManager := new(mocks.InQueueOrderManager)
	queueManager := new(mocks.QueueManager)
	stopOrderSubmissionManager := new(mocks.StopOrderSubmissionManager)

	logger := new(mocks.Logger)
	currencyService := new(mocks.CurrencyService)
	BTCUSDT := currency.Pair{
		ID:              1,
		Name:            "BTC-USDT",
		IsActive:        true,
		Spread:          1.5,
		ShowDigits:      6,
		BasisCoinID:     1,
		DependentCoinID: 2,
	}
	currencyService.On("GetPairByName", "BTC-USDT").Once().Return(BTCUSDT, nil)
	dataProcessor := processor.NewProcessor(redisClient, liveData, priceGenerator, klineService, orderbookService,
		mqttManager, stopOrderSubmissionManager, inQueueOrderManager, queueManager, logger, currencyService)
	ctx := context.Background()
	externalExchangeName := "binance"
	pairName := "BTC-USDT"
	externalExchangePairName := "BTC-USDT"
	data := []byte("")
	dataProcessor.ProcessDepth(ctx, externalExchangeName, pairName, externalExchangePairName, data)
	time.Sleep(50 * time.Millisecond)
	orderbookService.AssertExpectations(t)
	mqttManager.AssertExpectations(t)

}

func TestProcessor_ProcessKline(t *testing.T) {
	redisClient := new(mocks.RedisClient)

	liveData := new(mocks.LiveData)
	//liveData.On("SetKline", mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)

	priceGenerator := new(mocks.PriceGenerator)
	klineService := new(mocks.KlineService)
	//klineService.Wg.Add(1)
	//klineService.On("SaveKline", mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)
	orderbookService := new(mocks.OrderbookService)

	mqttManager := new(mocks.CentrifugoManager)
	mqttManager.On("PublishKline", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return()
	inQueueOrderManager := new(mocks.InQueueOrderManager)
	queueManager := new(mocks.QueueManager)
	queueManager.On("PublishKline", mock.Anything).Once().Return()
	stopOrderSubmissionManager := new(mocks.StopOrderSubmissionManager)
	logger := new(mocks.Logger)
	currencyService := new(mocks.CurrencyService)
	dataProcessor := processor.NewProcessor(redisClient, liveData, priceGenerator, klineService, orderbookService,
		mqttManager, stopOrderSubmissionManager, inQueueOrderManager, queueManager, logger, currencyService)
	ctx := context.Background()
	kline := processor.Kline{
		Pair:                "BTC-USDT",
		TimeFrame:           "1hour",
		KlineStartTime:      "2021-01-12 12:00:00",
		KlineCloseTime:      "2021-01-12 12:59:59",
		OpenPrice:           "32000",
		ClosePrice:          "32100",
		HighPrice:           "32400",
		LowPrice:            "31900",
		BaseVolume:          "6231.0000",
		QuoteVolume:         "6231.0000",
		TakerBuyBaseVolume:  "6231.0000",
		TakerBuyQuoteVolume: "6231.0000",
	}
	dataProcessor.ProcessKline(ctx, kline)
	time.Sleep(50 * time.Millisecond)
	queueManager.AssertExpectations(t)
	mqttManager.AssertExpectations(t)
}

func TestProcessor_ProcessTicker(t *testing.T) {
	redisClient := new(mocks.RedisClient)

	liveData := new(mocks.LiveData)
	liveData.On("GetPrice", mock.Anything, mock.Anything).Once().Return("32100", nil)
	liveData.On("SetPriceData", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)

	priceGenerator := new(mocks.PriceGenerator)
	priceGenerator.On("GetPairPriceBasedOnUSDT", mock.Anything, mock.Anything).Once().Return("32000", nil)
	klineService := new(mocks.KlineService)
	orderbookService := new(mocks.OrderbookService)

	mqttManager := new(mocks.CentrifugoManager)
	mqttManager.On("PublishTicker", mock.Anything, mock.Anything, mock.Anything).Once().Return()

	stopOrderSubmissionManager := new(mocks.StopOrderSubmissionManager)
	stopOrderSubmissionManager.On("Submit", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return()

	inQueueOrderManager := new(mocks.InQueueOrderManager)
	queueManager := new(mocks.QueueManager)
	inQueueOrderManager.On("HandleInQueueOrders", mock.Anything, mock.Anything, mock.Anything).Once().Return()

	logger := new(mocks.Logger)
	BTCUSDT := currency.Pair{
		ID:              1,
		Name:            "BTC-USDT",
		IsActive:        true,
		Spread:          1.5,
		ShowDigits:      6,
		BasisCoinID:     1,
		DependentCoinID: 2,
	}
	currencyService := new(mocks.CurrencyService)
	currencyService.On("GetPairByName", "BTC-USDT").Once().Return(BTCUSDT, nil)

	dataProcessor := processor.NewProcessor(redisClient, liveData, priceGenerator, klineService, orderbookService,
		mqttManager, stopOrderSubmissionManager, inQueueOrderManager, queueManager, logger, currencyService)
	ctx := context.Background()
	ticker := processor.Ticker{
		Pair:            "BTC-USDT",
		Price:           "32000",
		Percentage:      "5.1",
		ID:              1,
		EquivalentPrice: "32000",
		Volume:          "3214.2150000",
		High:            "3214.2150000",
		Low:             "3214.2150000",
	}
	dataProcessor.ProcessTicker(ctx, ticker)
	time.Sleep(50 * time.Millisecond)
	redisClient.AssertExpectations(t)
	liveData.AssertExpectations(t)
	mqttManager.AssertExpectations(t)
}
