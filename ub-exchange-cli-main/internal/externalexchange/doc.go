// Package externalexchange provides the Binance REST API integration layer.
// It handles:
//
//   - Order submission to Binance (limit, market orders)
//   - Order status queries and synchronization
//   - Account balance retrieval
//   - Trade history fetching
//
// This package is consumed by bot aggregation services and CLI commands
// that sync external exchange state with the local database.
package externalexchange
