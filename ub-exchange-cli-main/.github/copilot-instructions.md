---
description: "AI coding instructions for ub-exchange-cli-main Go trading engine"
applyTo: "ub-exchange-cli-main/**"
---

# ub-exchange-cli-main — AI Coding Instructions

## Critical Rules

### Financial Precision
- **ALWAYS** use `shopspring/decimal` for monetary calculations
- **NEVER** use `float64` for prices, amounts, balances, fees, or volumes
- When converting between types, use `decimal.NewFromString()` not `decimal.NewFromFloat()`
- All repository fields storing money are `decimal.Decimal` or `string` — verify before changing

### Error Handling
- Wrap all errors with context: `fmt.Errorf("service.Method: %w", err)`
- Log errors with structured fields via `zap.String("key", val)`, `zap.Error(err)`
- Never swallow errors — if you can't handle it, propagate it
- Check all error returns from Redis, GORM, MQTT, RabbitMQ operations

### DI Container
- All services registered in `internal/di/container.go` with string constant keys
- New services: add constant in `di/`, register in container, inject via `container.Get()`
- Never instantiate services directly in handlers — always use DI

### Repository Pattern
- Data access only through `internal/repository/` — never raw SQL in services
- GORM operations: use transactions for multi-table writes
- Test with real DB (MariaDB) + real Redis in `test/` directory

### API Pattern
- Public handlers: `internal/api/handler/`
- Admin handlers: `internal/api/adminhandler/`
- All responses: `response.Success(c, data)` or `response.Error(c, statusCode, message)`
- JWT auth via middleware — user context extracted in handler

### CLI Command Pattern
```go
type MyCommand struct {
    myService SomeService
    logger    platform.Logger
}

func (cmd *MyCommand) Run(ctx context.Context, flags []string) {
    // implementation
}
```

## Package Import Rules
- Internal packages: `exchange-go/internal/<package>`
- Never import `internal/di` from business packages — DI is top-level only
- Platform package (`internal/platform/`) provides all infrastructure interfaces

## Testing
- Integration tests live in `test/` (not alongside source)
- Tests require Docker services: MariaDB 10.5 + Redis
- Use `test/main_test.go` infrastructure (getContainer, getDb, getRedis, etc.)
- HTTP tests: `makeRequest()` / `makeAuthorizedRequest()` helpers
- Mock external services only (Binance, wallet) — use real DB/Redis

## Config Access
- Read config via `platform.Configs` interface (wraps Viper)
- Keys: `exchange.environment`, `db.dsn`, `redis.dsn`, `rabbitmq.dsn`, etc.
- Env var override: `UBEXCHANGE_DB_DSN` overrides `db.dsn`
- Test config: `config/config_test.yaml` (merged over main config)

## Redis Patterns
- Order queue: `RPUSH`/`LPOP` (FIFO queue)
- Orderbook: sorted sets (`ZADD` price as score, order ID as member)
- Live data: `HSET` with key patterns like `live_data:price`, `live_data:kline:BTC_USDT:1m`
- Cache: TTL-based via go-redis/cache/v8

## MQTT Topics
- Public: `main/trade/{pair}` (depth, trades, ticker)
- Private: `main/trade/user/{userID}/open-orders`, `main/trade/user/{userID}/crypto-payments`
- Publish via `platform.MqttClient.Publish(topic, qos, retained, payload)`

## What NOT to Change Without Careful Review
- Order matching logic in `internal/engine/` — critical financial path
- Balance operations in `internal/userbalance/` — double-entry accounting
- Trade execution in `internal/order/engineresulthandler.go`
- Redis sorted set operations for orderbook — race conditions possible
- GORM model tags — data mapping to existing MySQL schema
