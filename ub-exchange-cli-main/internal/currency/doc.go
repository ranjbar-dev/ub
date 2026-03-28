// Package currency manages cryptocurrency and fiat currency data for the
// exchange platform:
//
//   - Coin definitions (name, symbol, decimals, status)
//   - Trading pair configuration (base/quote currencies, fees, limits)
//   - Price data retrieval and caching
//   - gRPC candle/kline data service
//   - Favorite pair management per user
//
// Currency data is cached in Redis for fast access by the trading engine
// and API handlers.
package currency
