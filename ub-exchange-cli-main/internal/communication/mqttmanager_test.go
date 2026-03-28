// Package communication_test tests the MQTT manager for real-time message publishing. Covers:
//   - Publishing trade data to pair-specific MQTT topics
//   - Publishing kline (candlestick) data with timeframe routing
//   - Publishing ticker updates (price, percentage, volume) to pair topics
//   - Publishing order book snapshots (bids/asks) to pair topics
//   - Publishing order updates to user-specific private channels
//
// Test data: mock MQTT client and token with JSON-serialized trade, kline,
// ticker, and order book payloads for BTC-USDT pair.
package communication_test

import (
	"context"
	"exchange-go/internal/communication"
	"exchange-go/internal/mocks"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestManager_PublishTrades(t *testing.T) {
	mqttToken := new(mocks.MqttToken)
	mqttClient := new(mocks.MqttClient)
	mqttClient.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(mqttToken)
	mqttManager := communication.NewMqttManager(mqttClient)
	ctx := context.Background()
	pairName := "BTC-USDT"
	tradeData := "{" +
		"\"price\":\"31794.41000000\"," +
		"\"amount\":\"0.00300000\"," +
		"\"createdAt\":\"2021-01-23 12:46:26\"," +
		"\"isMaker\":true," +
		"\"ignore\":true" +
		"}"
	payload := []byte(tradeData)
	mqttManager.PublishTrades(ctx, pairName, payload)
	mqttClient.AssertExpectations(t)
}

func TestManager_PublishKline(t *testing.T) {
	mqttToken := new(mocks.MqttToken)
	mqttClient := new(mocks.MqttClient)
	mqttClient.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(mqttToken)
	mqttManager := communication.NewMqttManager(mqttClient)
	ctx := context.Background()
	pairName := "BTC-USDT"
	timeFrame := "1hour"
	klineData := "{" +
		"\"timeFrame\":\"1hour\"," +
		"\"ohlcStartTime\":\"2021-01-23 12:00:00\"," +
		"\"ohlcCloseTime\":\"2021-01-23 12:59:59\"," +
		"\"openPrice\":\"31643.90000000\"," +
		"\"closePrice\":\"31643.90000000\"," +
		"\"highPrice\":\"31643.90000000\"," +
		"\"lowPrice\":\"31643.90000000\"," +
		"\"baseVolume\":\"2293.69229500\"," +
		"\"quoteVolume\":\"2293.69229500\"," +
		"\"takerBuyBaseVolume\":\"2293.69229500\"," +
		"\"takerBuyQuoteVolume\":\"2293.69229500\"" +
		"}"

	payload := []byte(klineData)
	mqttManager.PublishKline(ctx, pairName, timeFrame, payload)
	mqttClient.AssertExpectations(t)
}

func TestManager_PublishTicker(t *testing.T) {
	mqttToken := new(mocks.MqttToken)
	mqttClient := new(mocks.MqttClient)
	mqttClient.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(mqttToken)
	mqttManager := communication.NewMqttManager(mqttClient)
	ctx := context.Background()

	pairName := "BTC-USDT"
	tickerData := "{" +
		"\"name\":\"BTC-USDT\"," +
		"\"price\":\"32400\"," +
		"\"formerPrice\":\"32300\"," +
		"\"percentage\":\"1.2\"," +
		"\"id\":\"1\"," +
		"\"equivalentPrice\":\"32400\"," +
		"\"volume\":\"2293.69229500\"," +
		"\"high\":\"2293.69229500\"," +
		"\"low\":\"2293.69229500\"" +
		"}"
	payload := []byte(tickerData)
	mqttManager.PublishTicker(ctx, pairName, payload)
	mqttClient.AssertExpectations(t)
}

func TestManager_PublishOrderBook(t *testing.T) {
	mqttToken := new(mocks.MqttToken)
	mqttClient := new(mocks.MqttClient)
	mqttClient.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(mqttToken)
	mqttManager := communication.NewMqttManager(mqttClient)
	ctx := context.Background()
	pairName := "BTC-USDT"
	tickerData := "{" +
		"\"bids\":[]," +
		"\"asks\":[]" +
		"}"

	payload := []byte(tickerData)
	mqttManager.PublishOrderBook(ctx, pairName, payload)
	mqttClient.AssertExpectations(t)
}

func TestManager_PublishOrderToOpenOrders(t *testing.T) {
	mqttToken := new(mocks.MqttToken)
	mqttClient := new(mocks.MqttClient)
	mqttClient.On("Publish", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Once().Return(mqttToken)
	mqttManager := communication.NewMqttManager(mqttClient)
	ctx := context.Background()


	orderPushPayLoad := `{
			"id":"1",
			"tradeAmount":"1.0",
			"tradePrice":"5000",
			"status":"FILLED",
			"type":"BUY"
		}`

	payload := []byte(orderPushPayLoad)
	privateChannelName := "someUniqueId"
	mqttManager.PublishOrderToOpenOrders(ctx, privateChannelName, payload)
	mqttClient.AssertExpectations(t)

}
