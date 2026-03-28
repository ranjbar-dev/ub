// Package engine implements the order matching engine that processes buy/sell
// orders against the order book. It uses a worker pool architecture:
//
//   - Orders arrive via a Redis queue (RPUSH/LPOP)
//   - A dispatcher pulls orders and assigns them to workers
//   - Workers match orders against the order book stored in Redis sorted sets
//   - Matched trades flow to the ResultHandler for settlement
//   - Partially filled orders are re-queued for further matching
//
// Key types:
//   - Engine: Main engine orchestrator that manages the worker pool
//   - OrderbookProvider: Reads bid/ask levels from Redis sorted sets
//   - QueueHandler: Manages the Redis-based order queue
//   - ResultHandler: Callback interface for processing matched trades
//
// The engine is started as the exchange-engine binary (cmd/exchange-engine)
// with a configurable number of workers (default: 10).
package engine
