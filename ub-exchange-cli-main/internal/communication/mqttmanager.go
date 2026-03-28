package communication

import (
	"context"
	"exchange-go/internal/platform"
)

const (
	Topic          = "main/trade/stream"
	TradeTopic     = "main/trade/trade-book/"
	KlineTopic     = "main/trade/kline/"
	TickerTopic    = "main/trade/ticker"
	OrderbookTopic = "main/trade/order-book/"

	UserPrivateTopicPrefix = "main/trade/user/"
	UserOpenOrdersPostfix  = "/open-orders/"
	UserPaymentsPostfix    = "/crypto-payments/"

	//TradeStream  = "trade"
	//KlineStream  = "kline"
	//TickerStream = "ticker"
	//OrderStream  = "order"
)

// MqttManager publishes real-time market and user data to MQTT topics for consumption
// by connected clients via the EMQX broker.
type MqttManager interface {
	// PublishTrades publishes executed trade data to the public trade-book topic for the given pair.
	PublishTrades(ctx context.Context, pairName string, payload []byte)
	// PublishKline publishes candlestick (kline) data to the public kline topic for the given pair and time frame.
	PublishKline(ctx context.Context, pairName string, timeFrame string, payload []byte)
	// PublishTicker publishes the latest ticker snapshot to the public ticker topic.
	PublishTicker(ctx context.Context, pairName string, payload []byte)
	// PublishOrderBook publishes the current order book state to the public order-book topic for the given pair.
	PublishOrderBook(ctx context.Context, pairName string, payload []byte)
	// PublishOrderToOpenOrders publishes an order update to a user's private open-orders channel.
	PublishOrderToOpenOrders(ctx context.Context, privateChannelName string, payload []byte)
	// PublishPayment publishes a crypto payment update to a user's private payments channel.
	PublishPayment(ctx context.Context, privateChannelName string, payload []byte)
}

type manager struct {
	mqttClient platform.MqttClient
}

type finalPayload struct {
	Stream  string `json:"stream"`
	Payload string `json:"payload"`
}

func (m *manager) PublishTrades(ctx context.Context, pairName string, payload []byte) {
	//payloadString := string(payload)
	//stream := TradeStream + "@" + pairName
	//finalPayload := finalPayload{Stream: stream, Payload: payloadString}
	//data, _ := json.Marshal(finalPayload)
	topic := TradeTopic + pairName
	m.mqttClient.Publish(topic, 0, false, payload)
}

func (m *manager) PublishKline(ctx context.Context, pairName string, timeFrame string, payload []byte) {
	topic := KlineTopic + timeFrame + "/" + pairName
	////topic := KlineTopic
	//payloadString := string(payload)
	//stream := KlineStream + "@" + pairName
	//finalPayload := finalPayload{Stream: stream, Payload: payloadString}
	//data, _ := json.Marshal(finalPayload)
	m.mqttClient.Publish(topic, 0, false, payload)
}

func (m *manager) PublishTicker(ctx context.Context, pairName string, payload []byte) {
	//payloadString := string(payload)
	//stream := TickerStream + "@" + pairName
	//finalPayload := finalPayload{Stream: stream, Payload: payloadString}
	//data, _ := json.Marshal(finalPayload)
	topic := TickerTopic
	m.mqttClient.Publish(topic, 0, false, payload)
}

func (m *manager) PublishOrderBook(ctx context.Context, pairName string, payload []byte) {
	//payloadString := string(payload)
	//stream := OrderStream + "@" + pairName
	//finalPayload := finalPayload{Stream: stream, Payload: payloadString}
	//data, _ := json.Marshal(finalPayload)

	topic := OrderbookTopic + pairName
	m.mqttClient.Publish(topic, 0, false, payload)
}

func (m *manager) PublishOrderToOpenOrders(ctx context.Context, privateChannelName string, payload []byte) {
	topic := UserPrivateTopicPrefix + privateChannelName + UserOpenOrdersPostfix
	m.mqttClient.Publish(topic, 0, false, payload)
}

func (m *manager) PublishPayment(ctx context.Context, privateChannelName string, payload []byte) {
	topic := UserPrivateTopicPrefix + privateChannelName + UserPaymentsPostfix
	m.mqttClient.Publish(topic, 0, false, payload)
}

func NewMqttManager(mqttClient platform.MqttClient) MqttManager {
	return &manager{mqttClient}
}
