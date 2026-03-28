// Package livedata_test tests the live data service for Redis-backed market data. Covers:
//   - Getting and setting pair prices with percentage and volume
//   - Updating and retrieving trade books (recent trade lists in Redis)
//   - Setting and getting kline (OHLCV candlestick) data per pair and timeframe
//   - Getting previous kline data for a given timeframe
//   - Setting and getting depth snapshots (order book snapshots with bids/asks)
//   - Setting and getting order books (price-to-amount maps)
//   - Batch retrieval of price data for multiple pairs
//   - Getting last insert timestamps for depth snapshots
//   - Getting kline data by pair name and timeframe key
//
// Test data: mock Redis client, miniredis for integration-style tests,
// and JSON-serialized trade/kline/depth/order book fixtures.
package livedata_test

import (
	"context"
	"exchange-go/internal/livedata"
	"exchange-go/internal/mocks"
	"exchange-go/internal/platform"
	"testing"

	"github.com/alicebob/miniredis"
	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_GetPrice(t *testing.T) {
	redisClient := new(mocks.RedisClient)
	redisClient.On("HGet", mock.Anything, mock.Anything, mock.Anything).Once().Return("32000", nil)
	liveDataService := livedata.NewLiveDataService(redisClient)
	ctx := context.Background()
	pairName := "BTC-USDT"
	price, err := liveDataService.GetPrice(ctx, pairName)
	assert.Nil(t, err)
	assert.Equal(t, "32000", price)
	redisClient.AssertExpectations(t)
}

func TestService_SetPriceData(t *testing.T) {
	redisClient := new(mocks.RedisClient)
	redisClient.On("HSet", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)
	liveDataService := livedata.NewLiveDataService(redisClient)
	ctx := context.Background()
	pairName := "BTC-USDT"
	price := "32000"
	percentage := "1.2"
	volume := "1000.10"
	err := liveDataService.SetPriceData(ctx, pairName, price, percentage, volume)
	assert.Nil(t, err)
	redisClient.AssertExpectations(t)
}

func TestService_UpdateTradeBook(t *testing.T) {
	redisClient := new(mocks.RedisClient)
	tradesString := "[{" +
		"\"price\":\"31794.41000000\"," +
		"\"amount\":\"0.00300000\"," +
		"\"createdAt\":\"2021-01-23 12:46:26\"," +
		"\"isMaker\":true," +
		"\"ignore\":true" +
		"}]"
	redisClient.On("Exists", mock.Anything, mock.Anything).Once().Return(true)
	redisClient.On("HGet", mock.Anything, mock.Anything, mock.Anything).Once().Return(tradesString, nil)
	redisClient.On("HSet", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)
	liveDataService := livedata.NewLiveDataService(redisClient)
	ctx := context.Background()
	pairName := "BTC-USDT"
	redisTrade := livedata.RedisTrade{
		Price:     "32152",
		Amount:    "1.5",
		CreatedAt: "2021-01-12 20:20:20",
		IsMaker:   true,
		Ignore:    true,
	}

	err := liveDataService.UpdateTradeBook(ctx, pairName, redisTrade)
	assert.Nil(t, err)
	redisClient.AssertExpectations(t)
}

func TestService_GetTradeBook(t *testing.T) {
	redisClient := new(mocks.RedisClient)
	tradesString := "[{" +
		"\"price\":\"31794.41000000\"," +
		"\"amount\":\"0.00300000\"," +
		"\"createdAt\":\"2021-01-23 12:46:26\"," +
		"\"isMaker\":true," +
		"\"ignore\":true" +
		"}]"
	redisClient.On("Exists", mock.Anything, mock.Anything).Once().Return(true)
	redisClient.On("HGet", mock.Anything, mock.Anything, mock.Anything).Once().Return(tradesString, nil)
	liveDataService := livedata.NewLiveDataService(redisClient)
	ctx := context.Background()
	pairName := "BTC-USDT"

	trades, err := liveDataService.GetTradeBook(ctx, pairName)
	assert.Nil(t, err)
	assert.Equal(t, "31794.41000000", trades[0].Price)
	assert.Equal(t, "0.00300000", trades[0].Amount)
	assert.Equal(t, true, trades[0].IsMaker)
	assert.Equal(t, true, trades[0].Ignore)
	redisClient.AssertExpectations(t)

}

func TestService_SetKline(t *testing.T) {
	redisClient := new(mocks.RedisClient)
	currentKlineData := "{" +
		"\"timeFrame\":\"1hour\"," +
		"\"ohlcStartTime\":\"2021-01-23 11:00:00\"," +
		"\"ohlcCloseTime\":\"2021-01-23 11:59:59\"," +
		"\"openPrice\":\"31643.90000000\"," +
		"\"closePrice\":\"31643.90000000\"," +
		"\"highPrice\":\"31643.90000000\"," +
		"\"lowPrice\":\"31643.90000000\"," +
		"\"baseVolume\":\"2293.69229500\"," +
		"\"quoteVolume\":\"2293.69229500\"," +
		"\"takerBuyBaseVolume\":\"2293.69229500\"," +
		"\"takerBuyQuoteVolume\":\"2293.69229500\"" +
		"}"

	redisClient.On("HGet", mock.Anything, mock.Anything, mock.Anything).Once().Return(currentKlineData, nil)
	redisClient.On("HSet", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Twice().Return(nil)
	liveDataService := livedata.NewLiveDataService(redisClient)
	ctx := context.Background()
	pairName := "BTC-USDT"
	kline := livedata.RedisKline{
		TimeFrame:           "1hour",
		KlineStartTime:      "2021-01-12 12:00:00",
		KlineCloseTime:      "2021-01-12 12:59:59",
		OpenPrice:           "32000",
		ClosePrice:          "32100",
		HighPrice:           "32400",
		LowPrice:            "31800",
		BaseVolume:          "2364.12",
		QuoteVolume:         "2354.15",
		TakerBuyBaseVolume:  "2364.12",
		TakerBuyQuoteVolume: "2354.15",
	}
	err := liveDataService.SetKline(ctx, pairName, kline)
	assert.Nil(t, err)
	redisClient.AssertExpectations(t)
}

func TestService_GetPreKline(t *testing.T) {
	redisClient := new(mocks.RedisClient)
	preKlineData := "{" +
		"\"timeFrame\":\"1hour\"," +
		"\"ohlcStartTime\":\"2021-01-23 11:00:00\"," +
		"\"ohlcCloseTime\":\"2021-01-23 11:59:59\"," +
		"\"openPrice\":\"31643.90000000\"," +
		"\"closePrice\":\"31642.90000000\"," +
		"\"highPrice\":\"31645.90000000\"," +
		"\"lowPrice\":\"31641.90000000\"," +
		"\"baseVolume\":\"2293.69229500\"," +
		"\"quoteVolume\":\"2293.69229500\"," +
		"\"takerBuyBaseVolume\":\"2293.69229500\"," +
		"\"takerBuyQuoteVolume\":\"2293.69229500\"" +
		"}"
	redisClient.On("HGet", mock.Anything, mock.Anything, mock.Anything).Once().Return(preKlineData, nil)
	liveDataService := livedata.NewLiveDataService(redisClient)
	ctx := context.Background()
	pairName := "BTC-USDT"
	TimeFrame := "1hour"
	preKline, err := liveDataService.GetPreKline(ctx, pairName, TimeFrame)
	assert.Nil(t, err)
	assert.Equal(t, "1hour", preKline.TimeFrame)
	assert.Equal(t, "2021-01-23 11:00:00", preKline.KlineStartTime)
	assert.Equal(t, "2021-01-23 11:59:59", preKline.KlineCloseTime)
	assert.Equal(t, "31643.90000000", preKline.OpenPrice)
	assert.Equal(t, "31642.90000000", preKline.ClosePrice)
	assert.Equal(t, "31645.90000000", preKline.HighPrice)
	assert.Equal(t, "31641.90000000", preKline.LowPrice)
	assert.Equal(t, "2293.69229500", preKline.BaseVolume)
	assert.Equal(t, "2293.69229500", preKline.QuoteVolume)
	assert.Equal(t, "2293.69229500", preKline.TakerBuyBaseVolume)
	assert.Equal(t, "2293.69229500", preKline.TakerBuyQuoteVolume)
	redisClient.AssertExpectations(t)
}

func TestService_SetDepthSnapshot(t *testing.T) {
	redisClient := new(mocks.RedisClient)
	redisClient.On("HSet", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)
	liveDataService := livedata.NewLiveDataService(redisClient)
	ctx := context.Background()
	pairName := "BTC-USDT"

	snapShot := livedata.RedisDepthSnapshot{
		LastUpdatedID: 21354545,
		Bids:          [][3]string{{"32000", "0.1", "0"}, {"32100", "0.2", "0"}},
		Asks:          [][3]string{{"321500", "0.1", "0"}, {"32200", "0.2", "0"}},
	}

	err := liveDataService.SetDepthSnapshot(ctx, pairName, snapShot)
	assert.Nil(t, err)
	redisClient.AssertExpectations(t)
}

func TestService_GetDepthSnapshot(t *testing.T) {
	redisClient := new(mocks.RedisClient)
	depthSnapShotData := "{\"lastUpdatedId\":8128975988,\"bids\":[[\"31451.63000000\",\"0.00064400\"]],\"asks\":[[\"31451.63000000\",\"0.00064400\"]]}"
	redisClient.On("HGet", mock.Anything, mock.Anything, mock.Anything).Once().Return(depthSnapShotData, nil)
	liveDataService := livedata.NewLiveDataService(redisClient)
	ctx := context.Background()
	pairName := "BTC-USDT"

	snapshot, err := liveDataService.GetDepthSnapshot(ctx, pairName)
	assert.Nil(t, err)
	assert.Equal(t, int64(8128975988), snapshot.LastUpdatedID)
	assert.Equal(t, "31451.63000000", snapshot.Bids[0][0])
	assert.Equal(t, "0.00064400", snapshot.Bids[0][1])
	assert.Equal(t, "31451.63000000", snapshot.Asks[0][0])
	assert.Equal(t, "0.00064400", snapshot.Asks[0][1])
	redisClient.AssertExpectations(t)

}

func TestService_SetOrderBook(t *testing.T) {
	redisClient := new(mocks.RedisClient)
	redisClient.On("HSet", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(nil)
	liveDataService := livedata.NewLiveDataService(redisClient)
	ctx := context.Background()
	pairName := "BTC-USDT"
	orderBook := livedata.RedisOrderBook{
		Bids: map[string]string{"32100": "0.01"},
		Asks: map[string]string{"33200": "0.3"},
	}

	err := liveDataService.SetOrderBook(ctx, pairName, orderBook)
	assert.Nil(t, err)
	redisClient.AssertExpectations(t)

}

func TestService_GetOrderBook(t *testing.T) {
	redisClient := new(mocks.RedisClient)
	orderBookData := "{\"bids\":{\"31451.63000000\":\"0.00064400\"},\"asks\":{\"31452.13000000\":\"0.00164400\"}}"
	redisClient.On("HGet", mock.Anything, mock.Anything, mock.Anything).Once().Return(orderBookData, nil)
	liveDataService := livedata.NewLiveDataService(redisClient)
	ctx := context.Background()
	pairName := "BTC-USDT"

	orderBook, err := liveDataService.GetOrderBook(ctx, pairName)
	assert.Nil(t, err)
	assert.Equal(t, "0.00064400", orderBook.Bids["31451.63000000"])
	assert.Equal(t, "0.00164400", orderBook.Asks["31452.13000000"])
	redisClient.AssertExpectations(t)

}

func TestService_GetPairsPricesData(t *testing.T) {
	s := miniredis.NewMiniRedis()
	defer s.Close()
	_ = s.Start()
	rc := redis.NewClient(&redis.Options{Addr: s.Addr()})
	ctx := context.Background()
	rc.HMSet(ctx, "live_data:pair_currency:BTC-USDT", "price", "50000")
	rc.HMSet(ctx, "live_data:pair_currency:BTC-USDT", "change_price_percentage", "1.2")
	rc.HMSet(ctx, "live_data:pair_currency:ETH-USDT", "price", "2000")
	rc.HMSet(ctx, "live_data:pair_currency:ETH-USDT", "change_price_percentage", "2.2")
	redisClient := platform.NewRedisTestClient(rc)
	liveDataService := livedata.NewLiveDataService(redisClient)
	pairNames := []string{"BTC-USDT", "ETH-USDT"}

	pairsPriceData, err := liveDataService.GetPairsPriceData(ctx, pairNames)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(pairsPriceData))
	for _, pairPriceData := range pairsPriceData {
		if pairPriceData.PairName == "BTC-USDT" {
			assert.Equal(t, "50000", pairPriceData.Price)
			assert.Equal(t, "1.2", pairPriceData.Percentage)
		}

		if pairPriceData.PairName == "ETH-USDT" {
			assert.Equal(t, "2000", pairPriceData.Price)
			assert.Equal(t, "2.2", pairPriceData.Percentage)
		}
	}

}

func TestService_GetLastInsertTime(t *testing.T) {
	s := miniredis.NewMiniRedis()
	defer s.Close()
	_ = s.Start()
	rc := redis.NewClient(&redis.Options{Addr: s.Addr()})
	ctx := context.Background()
	rc.HMSet(ctx, "live_data:pair_currency:BTC-USDT", "depth_snapshot_last_insert_time", 127123123456)
	redisClient := platform.NewRedisTestClient(rc)
	liveDataService := livedata.NewLiveDataService(redisClient)
	pairName := "BTC-USDT"
	lastInsertTime, err := liveDataService.GetLastInsertTime(ctx, pairName, livedata.DepthSnapshotLastInsertTime)
	assert.Nil(t, err)
	assert.Equal(t, "127123123456", lastInsertTime)

}

func TestService_GetKline(t *testing.T) {
	redisClient := new(mocks.RedisClient)
	klineData := "{" +
		"\"timeFrame\":\"1hour\"," +
		"\"ohlcStartTime\":\"2021-01-23 11:00:00\"," +
		"\"ohlcCloseTime\":\"2021-01-23 11:59:59\"," +
		"\"openPrice\":\"31643.90000000\"," +
		"\"closePrice\":\"31642.90000000\"," +
		"\"highPrice\":\"31645.90000000\"," +
		"\"lowPrice\":\"31641.90000000\"," +
		"\"baseVolume\":\"2293.69229500\"," +
		"\"quoteVolume\":\"2293.69229500\"," +
		"\"takerBuyBaseVolume\":\"2293.69229500\"," +
		"\"takerBuyQuoteVolume\":\"2293.69229500\"" +
		"}"
	redisClient.On("HGet", mock.Anything, "live_data:pair_currency:BTC-USDT", "kline_1hour").Once().Return(klineData, nil)
	liveDataService := livedata.NewLiveDataService(redisClient)
	kline, err := liveDataService.GetKline(context.Background(), "BTC-USDT", "1hour")
	assert.Nil(t, err)
	assert.Equal(t, "1hour", kline.TimeFrame)
	assert.Equal(t, "2021-01-23 11:00:00", kline.KlineStartTime)
	assert.Equal(t, "2021-01-23 11:59:59", kline.KlineCloseTime)
	assert.Equal(t, "31643.90000000", kline.OpenPrice)
	assert.Equal(t, "31642.90000000", kline.ClosePrice)
	assert.Equal(t, "31645.90000000", kline.HighPrice)
	assert.Equal(t, "31641.90000000", kline.LowPrice)
	assert.Equal(t, "2293.69229500", kline.BaseVolume)
	assert.Equal(t, "2293.69229500", kline.QuoteVolume)
	assert.Equal(t, "2293.69229500", kline.TakerBuyBaseVolume)
	assert.Equal(t, "2293.69229500", kline.TakerBuyQuoteVolume)
	redisClient.AssertExpectations(t)
}
