package types

import "context"

// ExternalWs represents a WebSocket connection to an external exchange for
// receiving real-time market data streams.
type ExternalWs interface {
	// Run starts the WebSocket connection and subscribes to the specified stream names
	// (e.g. ticker, kline, depth). It blocks until the context is cancelled.
	Run(ctx context.Context, streams []string)
}
