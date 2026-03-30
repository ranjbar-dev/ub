// Package order implements the core order management domain. This is the
// largest package in the codebase, handling the full order lifecycle:
//
// Order Creation:
//   - CreateManager: Validates and creates new orders (limit, market, stop)
//   - DecisionManager: Determines order routing (internal match vs external)
//
// Order Matching:
//   - EngineCommunicator: Submits orders to the matching engine via Redis queue
//   - EngineResultHandler: Processes matched trade results from the engine
//   - PostOrderMatchingService: Post-trade settlement (balance updates, DB persist, Centrifugo publish)
//
// Order Management:
//   - Service: CRUD operations, listing, cancellation, status queries
//   - AdminOrderManager: Admin-level order operations (fulfillment, manual trades)
//   - RedisManager: Redis cache operations for order data
//   - InQueueOrderManager: Manages orders waiting in the engine queue
//
// External Exchange:
//   - BotAggregationService: Aggregates and submits bot orders to Binance
//   - StopOrderSubmissionManager: Manages stop-order trigger and submission
//   - ForceTrader: Forced trade execution for admin operations
//
// Events:
//   - EventsHandler: Publishes order lifecycle events (created, matched, cancelled)
//   - TradeEventsHandler: Publishes trade execution events via Centrifugo
//
// All financial calculations use shopspring/decimal.
package order
