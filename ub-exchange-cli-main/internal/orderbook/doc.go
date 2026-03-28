// Package orderbook provides order book aggregation and retrieval from Redis.
// It reads bid/ask price levels from Redis sorted sets maintained by the
// matching engine and returns aggregated order book snapshots for API
// responses and market data feeds.
package orderbook
