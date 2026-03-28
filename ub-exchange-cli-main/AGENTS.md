# Exchange CLI — ub-exchange-cli-main

## Stack
- **Go 1.13** (target: Go 1.22+) | Gin 1.7 (HTTP) | GORM v2 1.21 (MySQL ORM)
- go-redis v8 | gRPC 1.40 + Protobuf | MQTT (eclipse/paho 1.3)
- RabbitMQ (streadway/amqp) | JWT (dgrijalva/jwt-go → golang-jwt/v5)
- shopspring/decimal (precision math) | sarulabs/di (DI) | uber/zap (logging)
- Sentry 0.11 | Viper 1.9 (config) | Binance REST+WS integration

## Architecture

### 4 Binary Entry Points
| Binary | File | Purpose | Port |
|--------|------|---------|------|
| exchange-cli | `cmd/exchange-cli/main.go` | 16 CLI cron commands | N/A |
| exchange-engine | `cmd/exchange-engine/main.go` | Order matching daemon (10 workers) | N/A |
| exchange-httpd | `cmd/exchange-httpd/main.go` | REST API (Gin) | 8000 (public), 8001 (admin) |
| exchange-ws | `cmd/exchange-ws/main.go` | Binance WebSocket listener | N/A |

### Internal Package Map (29 packages)
```
internal/
├── api/              # REST handlers (Gin), middleware, admin handlers
│   ├── handler/      # Public endpoints: /api/v1/*
│   ├── adminhandler/ # Admin endpoints: /api/v1/admin/*
│   └── middleware/    # Auth, recovery, CORS
├── auth/             # Login, register, 2FA (TOTP), MQTT ACL
├── command/          # 16 CLI ConsoleCommand implementations
├── communication/    # RabbitMQ message publishing
├── configuration/    # App settings service
├── country/          # Country management
├── currency/         # Coins, pairs, prices, gRPC candle service
├── di/               # DI container (~110+ services registered)
├── engine/           # Order matching: worker pool + Redis sorted sets
├── externalexchange/ # Binance REST API integration
├── externalexchangews/ # Binance WebSocket streams
├── jwt/              # JWT token helpers
├── livedata/         # Redis live data: price, kline, depth, trade book
├── mocks/            # Test mock implementations
├── order/            # Order CRUD, matching, trades, bot aggregation (20+ files)
├── orderbook/        # Orderbook aggregation from Redis
├── payment/          # Payment processing
├── platform/         # Infrastructure abstraction (12 files):
│   │                 # jwt.go, logger.go, redis.go, sql.go, cache.go,
│   │                 # mqtt.go, rabbitmq.go, http.go, bcrypt.go,
│   │                 # config.go, error.go, ws.go
├── processor/        # Real-time market data: trade, depth, kline, ticker
├── repository/       # 26 GORM data access repositories
├── response/         # APIResponse struct
├── transaction/      # Bank transfer logic
├── user/             # User CRUD, profiles, KYC levels
├── userbalance/      # Balance management
├── userdevice/       # Device fingerprinting
├── userwithdrawaddress/ # Withdrawal address management
├── wallet/           # External wallet service client
└── ws/               # WebSocket handler
```

### Data Flow
```
Client → Gin HTTP API (8000/8001)
                ↓
         Auth middleware (JWT)
                ↓
    Handler → Service → Repository (GORM/MySQL)
                ↓                     ↓
          Redis (cache/queue)    RabbitMQ → ub-communicator
                ↓
    Engine (10 workers, Redis sorted sets for orderbook)
                ↓
    Trade → MQTT publish → Client apps
```

### Order Matching Engine Flow
```
Order submitted → Redis queue (RPUSH)
  → Dispatcher pulls (LPOP) → Work struct
  → Worker matches against orderbook (Redis sorted sets: bid/ask)
  → ResultHandler: DB persist + MQTT publish + balance updates
  → Done or partial fill returned to queue
```

## CLI Commands (16 total)
| Command | Cron | Purpose |
|---------|------|---------|
| set-user-level | daily | Calculate trading level from volume |
| initialize-balance | manual | Create balances for new coins |
| generate-address | manual | Generate wallet addresses |
| retrieve-open-orders | 15min | Cache open orders to Redis |
| submit-bot-orders | 1min | Execute bot orders to Binance |
| sync-kline | 1min | OHLC candle data sync |
| check-withdrawals | 10min | Verify withdrawal status |
| update-orders-from-external | 30min | Sync external order status |
| retrieve-external-orders | periodic | Cache external orders to Redis |
| generate-kline-sync | manual | Initialize kline sync |
| ub-update-user-wallet-balances | daily | Sync wallet balances |
| ub-captcha-generate-keys | manual | Generate captcha RSA keys |
| ub-captcha-encryption | manual | Encrypt captcha data |
| ub-captcha-decryption | manual | Decrypt captcha data |
| delete-cache | manual | Clear Redis cache |

## Conventions
- **Config**: `config/config.yaml` loaded by Viper, env override prefix `UBEXCHANGE_`
- **DI**: sarulabs/di container with string constants in `internal/di/`
- **Logging**: uber/zap structured logger + Sentry error tracking
- **CLI**: All commands implement `ConsoleCommand` interface: `Run(ctx context.Context, flags []string)`
- **Financial math**: Always use `shopspring/decimal` — NEVER float64 for money
- **Repository pattern**: Each entity has its own repository in `internal/repository/`
- **API responses**: `{ status: bool, message: string, data: interface{} }`
- **Error handling**: Wrap errors with context: `fmt.Errorf("failed to X: %w", err)`
- **Redis keys**: Order books use sorted sets (ZADD/ZRANGEBYSCORE), cache keys `<entity>:<id>`

## Testing
- **30+ integration test files** in `test/` directory
- Tests require MySQL (MariaDB 10.5) + Redis running
- Test infrastructure: `test/main_test.go` (seeders, DB setup, HTTP helpers)
- Schema: `test/data/db.sql` | Seeders: `test/data/seed/`
- Run: `go test ./... --failfast` (inside Docker with DB+Redis services)

## CI/CD
- **GitLab CI**: `.gitlab-ci.yml`
- Test stage: `golang:1.16.7` image with MariaDB 10.5 + Redis
- Deploy: SSH → docker-compose restart on dev/prod servers
- Linting: `.golangci.yml` (revive linter)

## Build Commands
```bash
go build -o cli cmd/exchange-cli/main.go
go build -o httpd cmd/exchange-httpd/main.go
go build -o engine cmd/exchange-engine/main.go
go build -o ws cmd/exchange-ws/main.go
```

## Docker
- Go services run inside Docker (base image: `golang:1.16`)
- Dockerfile at `ub-server-main/.docker/go/Dockerfile`
- docker-compose services: `exchange-go`, `exchange-httpd-go`, `exchange-engine-go`
- Dependencies: MariaDB 10.2+, Redis 6.2, RabbitMQ 3.7, EMQX v4

## Security Concerns
- ⚠️ `config/config.yaml` contains hardcoded credentials (DB, RabbitMQ, MQTT, JWT passphrase)
- ⚠️ `dgrijalva/jwt-go` is ARCHIVED — must migrate to `golang-jwt/jwt/v5`
- ⚠️ RSA keys in `config/jwt/` — should use secrets manager
- ⚠️ No rate limiting on API endpoints

## Upgrade Priority
1. 🔴 Go 1.13 → 1.22+ (security CVEs, module improvements)
2. 🔴 dgrijalva/jwt-go → golang-jwt/jwt/v5 (archived, security)
3. 🟡 streadway/amqp → rabbitmq/amqp091-go (streadway archived)
4. 🟡 Gin 1.7 → 1.10+ (security fixes, performance)
5. 🟡 sentry-go 0.11 → 0.30+ (major improvements)
6. 🟠 go-redis v8 → v9 (breaking changes, better perf)
7. 🟠 golang.org/x/* packages (security patches)
8. 🟢 CI pipeline Go version update
9. 🟢 Docker base image update
