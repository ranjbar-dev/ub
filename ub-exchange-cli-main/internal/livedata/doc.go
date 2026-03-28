// Package livedata manages real-time market data stored in Redis. It handles:
//
//   - Live price data per trading pair
//   - Kline (candlestick) OHLCV data at multiple intervals
//   - Order book depth snapshots
//   - Recent trade book entries
//
// Data is written by processors (internal/processor) that consume Binance
// WebSocket streams, and read by API handlers and MQTT publishers to push
// updates to connected clients.
package livedata
