package platform

import (
	"context"
	"github.com/gorilla/websocket"
	"net/http"
)

// WsConnection represents an active WebSocket connection for reading and writing
// messages, with support for ping/pong keep-alive handlers.
type WsConnection interface {
	// WriteMessage sends a message of the given type (e.g., TextMessage, BinaryMessage) over the connection.
	WriteMessage(messageType int, data []byte) error
	// ReadMessage reads the next message from the connection, returning the message type and payload.
	ReadMessage() (messageType int, p []byte, err error)
	// SetPongHandler sets the handler invoked when a pong control message is received from the peer.
	SetPongHandler(h func(appData string) error)
	// SetPingHandler sets the handler invoked when a ping control message is received from the peer.
	SetPingHandler(h func(appData string) error)
	// Close performs the WebSocket closing handshake and releases the underlying connection.
	Close() error
}

// WsClient provides WebSocket dialing for establishing connections to external
// services such as Binance WebSocket streams.
type WsClient interface {
	// Dial opens a new WebSocket connection to the given URL with optional HTTP headers.
	Dial(ctx context.Context, urlStr string, requestHeader http.Header) (WsConnection, error)
}

type wsClient struct {
	dialer *websocket.Dialer
}

func (wsClient *wsClient) Dial(ctx context.Context, urlStr string, requestHeader http.Header) (WsConnection, error) {
	c, _, err := wsClient.dialer.DialContext(ctx, urlStr, requestHeader)
	return c, err
}

func NewWsClient() WsClient {
	dialer := websocket.DefaultDialer
	return &wsClient{dialer}
}
