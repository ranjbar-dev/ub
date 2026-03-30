// Package processor implements real-time market data processors that consume
// data from Binance WebSocket streams and write to Redis/Centrifugo:
//
//   - TradeProcessor: Processes individual trade events
//   - DepthProcessor: Processes order book depth updates
//   - KlineProcessor: Processes candlestick/OHLCV data
//   - TickerProcessor: Processes 24hr ticker statistics
//
// Each processor transforms Binance-format data into the exchange's internal
// format, stores it in Redis (via livedata package), and publishes updates
// to connected clients via Centrifugo.
package processor
