# AGENTS.md — ub-exchange-cli-main

> Universal AI instruction file. Read by Copilot, Claude, Cursor, and all AI coding tools.

## Project

Exchange-CLI: Go-based cryptocurrency trading engine for the UnitedBit exchange platform.
Provides order matching, REST API (43 endpoints), CLI tools (16 commands), and real-time market data via Binance WebSocket.

**Stack:** Go 1.22, Gin 1.10, GORM 1.21/MySQL, Redis 6.2, Centrifugo v4, RabbitMQ, JWT (RS256), gRPC, shopspring/decimal, Viper, Zap, Sentry.

## Setup & Commands

```bash
# Build
go build ./cmd/exchange-cli/ && go build ./cmd/exchange-httpd/ && go build ./cmd/exchange-engine/ && go build ./cmd/exchange-ws/

# Test
go test ./... --failfast
go test -race ./...        # before committing concurrency changes

# Lint
golangci-lint run

# Run
go run cmd/exchange-httpd/main.go         # :8000 public, :8001 admin
go run cmd/exchange-engine/main.go        # 10 worker goroutines
go run cmd/exchange-cli/main.go <command> # one-shot CLI command
```

## Architecture

**4 binaries:**

| Binary | Purpose | Notes |
|--------|---------|-------|
| `exchange-httpd` | REST API | :8000 public, :8001 admin, + crash-recovery goroutine |
| `exchange-engine` | Order matching daemon | 1 dispatcher (BLPop) + 10 worker goroutines |
| `exchange-cli` | 16 cron/manual commands | Runs one command then exits |
| `exchange-ws` | Binance WebSocket listener | Streams: depth, trade, kline, ticker |

**Data flow:**
```
Client → Gin HTTP → Auth middleware (JWT) → Handler → Service → Repository (GORM/MySQL)
    ↕ Redis (cache/queue) → Engine (10 workers, Redis sorted sets) → Trade → Centrifugo → Client apps
    ↕ RabbitMQ → ub-communicator (email/SMS)
```

**29 internal packages (grouped):**

- **API layer:** `api/` (handlers + admin + middleware), `response/` (standard envelope), `auth/` (login/register/2FA/Centrifugo tokens)
- **Trading core:** `engine/` (order matching worker pool), `order/` (CRUD + matching + trades), `orderbook/` (aggregation), `userbalance/` (lock/unlock/credit/debit), `payment/` (withdraw/deposit/auto-exchange)
- **Data:** `repository/` (26 GORM repos), `currency/` (coins/pairs/kline/gRPC candle), `livedata/` (Redis market data), `configuration/`, `country/`, `transaction/`
- **Communication:** `communication/` (Centrifugo + RabbitMQ publish), `jwt/` (token create/validate)
- **External:** `externalexchange/` (Binance REST), `externalexchangews/` (Binance WS), `processor/` (WS data → Redis + Centrifugo), `wallet/` (external wallet HTTP client), `ws/` (WebSocket handler)
- **Infrastructure:** `platform/` (13 interface abstractions), `di/` (109 service registrations), `command/` (16 CLI commands), `user/` (profiles/KYC/2FA/permissions), `mocks/` (82 test mocks)
- **Utility:** `userdevice/`, `userwithdrawaddress/`

## Code Style & Conventions

- **Financial precision:** ALWAYS `shopspring/decimal` for money — NEVER `float64`
- **Error handling:** `fmt.Errorf("service.Method: %w", err)` — never swallow errors
- **Logging:** zap structured fields — `zap.String("key", val)`, `zap.Error(err)`
- **Naming:** Go standard — PascalCase exports, camelCase locals
- **DI:** sarulabs/di container with string constants in `internal/di/` — never instantiate services directly
- **Repository pattern:** data access only via `internal/repository/` — never raw SQL in services
- **API responses:** `response.Success(c, data)` or `response.Error(c, statusCode, message)` → `{"status": bool, "message": string, "data": ...}`
- **Config:** `platform.Configs` interface (Viper), env override prefix `UBEXCHANGE_`
- **CLI commands:** implement `ConsoleCommand` interface — `Run(ctx context.Context, flags []string)`

## Do

- Use `decimal.NewFromString()` not `decimal.NewFromFloat()`
- Wrap all errors with context: `fmt.Errorf("pkg.Func: %w", err)`
- Use GORM transactions for multi-table writes
- Check all error returns from Redis, GORM, Centrifugo, RabbitMQ
- Run `go test -race` before committing concurrency changes
- Use `platform` package interfaces (`Logger`, `RedisClient`, `SqlClient`, `Cache`, etc.)
- Follow handler → service → repository layering
- Register new services in `internal/di/` with string constant keys

## Don't

- NEVER use `float64` for prices, amounts, balances, fees, or volumes
- NEVER import `internal/di` from business packages — DI is top-level only
- NEVER instantiate services directly in handlers — use DI container
- NEVER commit secrets or hardcode credentials
- NEVER modify GORM model tags without verifying MySQL schema compatibility
- NEVER create goroutines without context/WaitGroup
- NEVER use global variables for shared state
- NEVER put raw SQL in service layer — use repository

## Testing

- **Location:** integration tests in `test/` directory (NOT alongside source)
- **Requirements:** Docker services — MariaDB 10.5 + Redis
- **Infrastructure:** `test/main_test.go` provides `getContainer()`, `getDb()`, `getRedis()`, `getUserActor()`, `getAdminUserActor()`
- **HTTP tests:** `makeRequest(method, path, body)`, `makeAuthorizedRequest(method, path, body, token)`
- **Mocking:** mock external services only (Binance, wallet) — use real DB/Redis
- **Coverage:** 28 test files covering API, CLI commands, order matching, payments, balances

## Safety & Critical Paths

> Review these areas carefully — bugs here cause financial loss.

- **Order matching:** `internal/engine/` — critical financial path
- **Balance operations:** `internal/userbalance/` — double-entry accounting
- **Trade execution:** `internal/order/engineresulthandler.go`
- **Redis orderbook:** sorted set operations — race conditions possible
- **Payment flow:** pre-withdraw → 2FA verify → email confirm → lock balance → wallet service → DB persist
- **Engine queue:** Redis `RPUSH`/`BLPop` on `engine:queue:orders`, orderbook in sorted sets `order-book:bid/ask:<pair>`

## Key Redis Patterns

| Pattern | Type | Purpose |
|---------|------|---------|
| `order-book:bid:<pair>` / `order-book:ask:<pair>` | Sorted set (score=price) | Order books |
| `engine:queue:orders` | List (RPUSH/BLPop) | Engine work queue |
| `queue:stop:order:<type>:<pair>` | Sorted set | Stop order triggers |
| `live_data:pair_currency:<pair>` | Hash | Price, volume, klines, depth |
| `withdraw-confirmation:<userId>` | Hash (TTL 3h) | Withdrawal email confirm |
| `forgot-password:<userId>` | Hash (TTL 3h) | Password reset tokens |
| `phone-confirmation:<userId>` | Hash (TTL 3h) | Phone verification codes |

## References

- `ARCHITECTURE.md` — Detailed data flows, dependency graph, Redis structures, DI service registry
- `docs/` — Domain documentation (order matching, Binance WS, orderbook, user balance)
- `tasks/` — 13 documented refactoring tasks with priority