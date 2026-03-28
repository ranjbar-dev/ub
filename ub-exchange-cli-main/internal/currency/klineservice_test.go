// Package currency_test tests the kline service component. Covers:
//   - Retrieving the last close price for a trading pair via gRPC candle client
//   - Fetching 24-hour high and low prices for a pair from a given date
//   - Creating new kline sync records with time range and sync type parameters
//   - Updating existing kline sync records
//
// Test data: mock candle gRPC client, kline sync repository, and live data service
// with BTC-USDT pair fixtures.
package currency_test

import (
	"context"
	"exchange-go/internal/currency"
	"exchange-go/internal/currency/candle"
	"exchange-go/internal/mocks"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestKlineService_GetClosePriceForPair(t *testing.T) {
	liveDataService := new(mocks.LiveData)
	candleGRPCClient := new(mocks.CandleGRPCClient)
	candleGRPCClient.On("GetLastPriceForPair", "BTC-USDT", mock.Anything).Once().Return("50000", nil)
	klineSyncRepo := new(mocks.KlineSyncRepository)

	klineService := currency.NewKlineService(klineSyncRepo, liveDataService, candleGRPCClient)
	price, err := klineService.GetLastPriceForPair(context.Background(), "BTC-USDT", time.Now())
	assert.Nil(t, err)
	assert.Equal(t, "50000", price)
	candleGRPCClient.AssertExpectations(t)
	liveDataService.AssertExpectations(t)

}

func TestKlineService_GetHighAndLowPriceFromDateForPairByPairName(t *testing.T) {
	data := candle.HighAndLowPrice{
		High: "55000.0",
		Low:  "50000.0",
	}
	liveDataService := new(mocks.LiveData)
	candleGRPCClient := new(mocks.CandleGRPCClient)
	candleGRPCClient.On("GetHighAndLowPriceForPairFromDate", "BTC-USDT", mock.Anything).Once().Return(data, nil)
	klineSyncRepo := new(mocks.KlineSyncRepository)

	klineService := currency.NewKlineService(klineSyncRepo, liveDataService, candleGRPCClient)
	result, err := klineService.GetHighAndLowPriceFromDateForPairByPairName("BTC-USDT", time.Now())
	assert.Nil(t, err)
	assert.Equal(t, "55000.0", result.High)
	assert.Equal(t, "50000.0", result.Low)
	candleGRPCClient.AssertExpectations(t)

}

func TestKlineService_CreateKlineSync(t *testing.T) {
	liveDataService := new(mocks.LiveData)
	candleGRPCClient := new(mocks.CandleGRPCClient)
	klineSyncRepo := new(mocks.KlineSyncRepository)
	klineSyncRepo.On("Create", mock.Anything).Once().Return(nil)

	klineService := currency.NewKlineService(klineSyncRepo, liveDataService, candleGRPCClient)
	params := currency.CreateKlineSyncParams{
		StartTime:  "2021-05-21 00:00:00",
		EndTime:    "2021-06-21 00:00:00",
		PairID:     1,
		TimeFrame:  currency.Timeframe1minute,
		WithUpdate: false,
		Type:       currency.SyncTypeAuto,
	}
	err := klineService.CreateKlineSync(params)
	assert.Nil(t, err)
	klineSyncRepo.AssertExpectations(t)

}

func TestUpdateKlineSync(t *testing.T) {
	liveDataService := new(mocks.LiveData)
	candleGRPCClient := new(mocks.CandleGRPCClient)
	klineSyncRepo := new(mocks.KlineSyncRepository)
	klineSyncRepo.On("Update", mock.Anything).Once().Return(nil)

	klineService := currency.NewKlineService(klineSyncRepo, liveDataService, candleGRPCClient)
	err := klineService.UpdateKlineSync(&currency.KlineSync{})
	assert.Nil(t, err)
	klineSyncRepo.AssertExpectations(t)
}
