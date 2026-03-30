package communication_test

import (
	"context"
	"exchange-go/internal/communication"
	"exchange-go/internal/mocks"
	"testing"

	"github.com/stretchr/testify/mock"
)

func TestCentrifugoManager_PublishTrades(t *testing.T) {
	centrifugoClient := new(mocks.CentrifugoClient)
	logger := new(mocks.Logger)
	centrifugoClient.On("Publish", "trade:trade-book:BTC-USDT", mock.Anything).Once().Return(nil)
	mgr := communication.NewCentrifugoManager(centrifugoClient, logger)
	ctx := context.Background()
	pairName := "BTC-USDT"
	tradeData := `{"price":"31794.41000000","amount":"0.00300000","createdAt":"2021-01-23 12:46:26","isMaker":true,"ignore":true}`
	payload := []byte(tradeData)
	mgr.PublishTrades(ctx, pairName, payload)
	centrifugoClient.AssertExpectations(t)
}

func TestCentrifugoManager_PublishKline(t *testing.T) {
	centrifugoClient := new(mocks.CentrifugoClient)
	logger := new(mocks.Logger)
	centrifugoClient.On("Publish", "trade:kline:1hour:BTC-USDT", mock.Anything).Once().Return(nil)
	mgr := communication.NewCentrifugoManager(centrifugoClient, logger)
	ctx := context.Background()
	pairName := "BTC-USDT"
	timeFrame := "1hour"
	klineData := `{"timeFrame":"1hour","openPrice":"31643.90000000","closePrice":"31643.90000000"}`
	payload := []byte(klineData)
	mgr.PublishKline(ctx, pairName, timeFrame, payload)
	centrifugoClient.AssertExpectations(t)
}

func TestCentrifugoManager_PublishTicker(t *testing.T) {
	centrifugoClient := new(mocks.CentrifugoClient)
	logger := new(mocks.Logger)
	centrifugoClient.On("Publish", "trade:ticker", mock.Anything).Once().Return(nil)
	mgr := communication.NewCentrifugoManager(centrifugoClient, logger)
	ctx := context.Background()
	pairName := "BTC-USDT"
	tickerData := `{"name":"BTC-USDT","price":"32400","volume":"2293.69229500"}`
	payload := []byte(tickerData)
	mgr.PublishTicker(ctx, pairName, payload)
	centrifugoClient.AssertExpectations(t)
}

func TestCentrifugoManager_PublishOrderBook(t *testing.T) {
	centrifugoClient := new(mocks.CentrifugoClient)
	logger := new(mocks.Logger)
	centrifugoClient.On("Publish", "trade:order-book:BTC-USDT", mock.Anything).Once().Return(nil)
	mgr := communication.NewCentrifugoManager(centrifugoClient, logger)
	ctx := context.Background()
	pairName := "BTC-USDT"
	orderBookData := `{"bids":[],"asks":[]}`
	payload := []byte(orderBookData)
	mgr.PublishOrderBook(ctx, pairName, payload)
	centrifugoClient.AssertExpectations(t)
}

func TestCentrifugoManager_PublishOrderToOpenOrders(t *testing.T) {
	centrifugoClient := new(mocks.CentrifugoClient)
	logger := new(mocks.Logger)
	centrifugoClient.On("Publish", "user:someUniqueId:open-orders", mock.Anything).Once().Return(nil)
	mgr := communication.NewCentrifugoManager(centrifugoClient, logger)
	ctx := context.Background()
	orderPushPayLoad := `{"id":"1","tradeAmount":"1.0","tradePrice":"5000","status":"FILLED","type":"BUY"}`
	payload := []byte(orderPushPayLoad)
	privateChannelName := "someUniqueId"
	mgr.PublishOrderToOpenOrders(ctx, privateChannelName, payload)
	centrifugoClient.AssertExpectations(t)
}

func TestCentrifugoManager_PublishPayment(t *testing.T) {
	centrifugoClient := new(mocks.CentrifugoClient)
	logger := new(mocks.Logger)
	centrifugoClient.On("Publish", "user:someUniqueId:crypto-payments", mock.Anything).Once().Return(nil)
	mgr := communication.NewCentrifugoManager(centrifugoClient, logger)
	ctx := context.Background()
	paymentData := `{"id":"1","amount":"0.5","status":"completed"}`
	payload := []byte(paymentData)
	privateChannelName := "someUniqueId"
	mgr.PublishPayment(ctx, privateChannelName, payload)
	centrifugoClient.AssertExpectations(t)
}
