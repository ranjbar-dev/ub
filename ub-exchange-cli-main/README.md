# ub-exchange-cli — Go Trading Engine

Go backend for the UnitedBit cryptocurrency exchange platform. Provides order matching, REST API, CLI tools, and real-time market data via Binance WebSocket integration.

## Architecture

Four independent binaries share a single codebase with 29 internal packages:

| Binary | Command | Purpose | Port |
|--------|---------|---------|------|
| `exchange-httpd` | `go run cmd/exchange-httpd/main.go` | REST API (Gin) — public + admin | :8000 / :8001 |
| `exchange-engine` | `go run cmd/exchange-engine/main.go` | Order matching daemon (10 workers, Redis sorted sets) | — |
| `exchange-cli` | `go run cmd/exchange-cli/main.go <cmd>` | 16 CLI cron/maintenance commands | — |
| `exchange-ws` | `go run cmd/exchange-ws/main.go <streams>` | Binance WebSocket market data listener | — |

```
Client -> Gin HTTP API -> Auth middleware (JWT) -> Handler -> Service -> Repository (GORM/MySQL)
                                                               |                   |
                                                         Redis (cache/queue)   RabbitMQ -> ub-communicator
                                                               |
                                                   Engine (10 workers, Redis sorted sets)
                                                               |
                                                   Trade -> Centrifugo publish -> Client apps
```

## Prerequisites

- **Go 1.22+**
- **MariaDB 10.5+** (or MySQL 8.0+)
- **Redis 6.2+**
- **RabbitMQ 3.7+**
- **Centrifugo v5** (real-time messaging)
- **Docker & docker-compose** (recommended for local development)

## Quick Start

### Build All Binaries

```bash
go build -o exchange-cli    cmd/exchange-cli/main.go
go build -o exchange-httpd  cmd/exchange-httpd/main.go
go build -o exchange-engine cmd/exchange-engine/main.go
go build -o exchange-ws     cmd/exchange-ws/main.go
```

### Configuration

1. Copy and edit `config/config.yaml` — set DB, Redis, RabbitMQ, Centrifugo, JWT, and wallet credentials
2. Environment variable overrides use prefix `UBEXCHANGE_` (e.g., `UBEXCHANGE_DB_DSN`)
3. Place RSA keys in `config/jwt/private.pem` and `config/jwt/public.pem`

### Run Services

```bash
# 1. Start HTTP API (public :8000, admin :8001)
./exchange-httpd

# 2. Start order matching engine (10 worker goroutines)
./exchange-engine

# 3. Start Binance WebSocket streams (managed by supervisord in production)
./exchange-ws depth
./exchange-ws ticker trade
./exchange-ws kline_1m kline_5m kline_1h kline_1d

# 4. Run CLI commands
./exchange-cli set-user-level
./exchange-cli sync-kline
./exchange-cli check-withdrawals
```

### Run Tests

```bash
# Requires MariaDB + Redis running (typically inside Docker)
go test ./... --failfast
```

## Stack

| Category | Library | Version |
|----------|---------|---------|
| HTTP | Gin | 1.10.1 |
| ORM | GORM v2 | 1.21.15 |
| Redis | go-redis/v8 | 8.11.3 |
| Real-time | Centrifugo HTTP API | phpcent |
| RabbitMQ | amqp091-go | 1.10.0 |
| JWT | golang-jwt/v5 | 5.2.1 |
| gRPC | grpc | 1.40.0 |
| WebSocket | gorilla/websocket | 1.5.0 |
| DI | sarulabs/di | 2.0.0 |
| Decimal | shopspring/decimal | 1.2.0 |
| Config | Viper | 1.9.0 |
| Logging | uber/zap | 1.19.1 |
| Sentry | sentry-go | 0.30.0 |

## Documentation

- **[AGENTS.md](AGENTS.md)** — Complete AI-agent reference: all packages, endpoints, DI services, CLI commands, conventions
- **[ARCHITECTURE.md](ARCHITECTURE.md)** — Data flows, dependency graph, Redis structures, RabbitMQ topology, Centrifugo channels
- **[docs/](docs/)** — Domain-specific documentation (order matching, Binance WS, orderbook, user balance)
- **[tasks/](tasks/)** — 13 documented refactoring tasks with priority

## CI/CD

GitLab CI pipeline with test, dev deploy, and production deploy stages. See `.gitlab-ci.yml`.

## Docker

Services run inside Docker. Dockerfile is in `ub-server-main/.docker/go/Dockerfile`. WebSocket streams are managed by supervisord (`.docker/supervisor/go-supervisord.conf`).
