package communication

import (
	"context"
	"encoding/json"
	"exchange-go/internal/platform"

	"go.uber.org/zap"
)

const (
	TradeChannel     = "trade:trade-book:"
	KlineChannel     = "trade:kline:"
	TickerChannel    = "trade:ticker"
	OrderbookChannel = "trade:order-book:"

	UserChannelPrefix    = "user:"
	UserOpenOrdersSuffix = ":open-orders"
	UserPaymentsSuffix   = ":crypto-payments"
)

// CentrifugoManager publishes real-time market and user data to Centrifugo channels
// for consumption by connected clients.
type CentrifugoManager interface {
	// PublishTrades publishes executed trade data to the public trade-book channel for the given pair.
	PublishTrades(ctx context.Context, pairName string, payload []byte)
	// PublishKline publishes candlestick (kline) data to the public kline channel for the given pair and time frame.
	PublishKline(ctx context.Context, pairName string, timeFrame string, payload []byte)
	// PublishTicker publishes the latest ticker snapshot to the public ticker channel.
	PublishTicker(ctx context.Context, pairName string, payload []byte)
	// PublishOrderBook publishes the current order book state to the public order-book channel for the given pair.
	PublishOrderBook(ctx context.Context, pairName string, payload []byte)
	// PublishOrderToOpenOrders publishes an order update to a user's private open-orders channel.
	PublishOrderToOpenOrders(ctx context.Context, privateChannelName string, payload []byte)
	// PublishPayment publishes a crypto payment update to a user's private payments channel.
	PublishPayment(ctx context.Context, privateChannelName string, payload []byte)
}

type centrifugoManager struct {
	client platform.CentrifugoClient
	logger platform.Logger
}

func (m *centrifugoManager) PublishTrades(ctx context.Context, pairName string, payload []byte) {
	channel := TradeChannel + pairName
	data := m.unmarshalPayload(payload)
	if err := m.client.Publish(channel, data); err != nil {
		m.logger.Warn("failed to publish trades",
			zap.Error(err),
			zap.String("channel", channel),
		)
	}
}

func (m *centrifugoManager) PublishKline(ctx context.Context, pairName string, timeFrame string, payload []byte) {
	channel := KlineChannel + timeFrame + ":" + pairName
	data := m.unmarshalPayload(payload)
	if err := m.client.Publish(channel, data); err != nil {
		m.logger.Warn("failed to publish kline",
			zap.Error(err),
			zap.String("channel", channel),
		)
	}
}

func (m *centrifugoManager) PublishTicker(ctx context.Context, pairName string, payload []byte) {
	channel := TickerChannel
	data := m.unmarshalPayload(payload)
	if err := m.client.Publish(channel, data); err != nil {
		m.logger.Warn("failed to publish ticker",
			zap.Error(err),
			zap.String("channel", channel),
		)
	}
}

func (m *centrifugoManager) PublishOrderBook(ctx context.Context, pairName string, payload []byte) {
	channel := OrderbookChannel + pairName
	data := m.unmarshalPayload(payload)
	if err := m.client.Publish(channel, data); err != nil {
		m.logger.Warn("failed to publish order book",
			zap.Error(err),
			zap.String("channel", channel),
		)
	}
}

func (m *centrifugoManager) PublishOrderToOpenOrders(ctx context.Context, privateChannelName string, payload []byte) {
	channel := UserChannelPrefix + privateChannelName + UserOpenOrdersSuffix
	data := m.unmarshalPayload(payload)
	if err := m.client.Publish(channel, data); err != nil {
		m.logger.Warn("failed to publish open orders",
			zap.Error(err),
			zap.String("channel", channel),
		)
	}
}

func (m *centrifugoManager) PublishPayment(ctx context.Context, privateChannelName string, payload []byte) {
	channel := UserChannelPrefix + privateChannelName + UserPaymentsSuffix
	data := m.unmarshalPayload(payload)
	if err := m.client.Publish(channel, data); err != nil {
		m.logger.Warn("failed to publish payment",
			zap.Error(err),
			zap.String("channel", channel),
		)
	}
}

func (m *centrifugoManager) unmarshalPayload(payload []byte) interface{} {
	var data interface{}
	if err := json.Unmarshal(payload, &data); err != nil {
		return string(payload)
	}
	return data
}

func NewCentrifugoManager(client platform.CentrifugoClient, logger platform.Logger) CentrifugoManager {
	return &centrifugoManager{client: client, logger: logger}
}
