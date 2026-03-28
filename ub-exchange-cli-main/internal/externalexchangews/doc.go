// Package externalexchangews provides Binance WebSocket stream integration.
// It connects to Binance's real-time market data streams to receive:
//
//   - Ticker price updates
//   - Kline/candlestick data
//   - Depth (order book) updates
//   - Trade stream data
//
// Received data is forwarded to internal processors (internal/processor)
// for caching in Redis and publishing to clients via MQTT.
//
// This package powers the exchange-ws binary (cmd/exchange-ws).
package externalexchangews
