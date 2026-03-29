# UnitedBit Exchange Platform — Workspace Guidelines

> **For AI agents:** This is the definitive top-level guide. Read this before
> making changes to any sub-project. Each sub-project also has its own `AGENTS.md`
> with project-specific details.

## Overview

This is a **cryptocurrency exchange platform** monorepo with 6 sub-projects:

| Project | Tech | Purpose |
|---------|------|---------|
| `ub-server-main` | PHP 8.1+ / Symfony 5.4 LTS / Doctrine ORM 2.14 / MariaDB 10.2 | Backend REST API (13 bundles) |
| `ub-admin-main` | React 17.0.2 / TypeScript 5.4.5 / Redux-Saga / Material-UI 4.12 | Admin panel SPA |
| `ub-client-cabinet-main` | React 18.3.1 / TypeScript 5.4.5 / Redux-Saga / Webpack 4.44 | Client trading dashboard |
| `ub-app-main` | Flutter 2.x / Dart ≥2.11 <3.0 (pre-null-safety) / GetX 4.3 / Dio 4.0 | Mobile + Web app |
| `ub-exchange-cli-main` | Go 1.22 / Gin 1.10 / GORM 1.21 / Redis / gRPC | Trading engine, CLI commands, HTTP API |
| `ub-communicator-main` | Go 1.24 / RabbitMQ (amqp091-go) / MongoDB 1.17 | Email/SMS notification service |

## Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│  Client Apps                                                        │
│  ub-app-main (Flutter)  ·  ub-client-cabinet-main (React SPA)      │
└────────────────────┬────────────────────────────────────────────────┘
                     │  HTTPS  /api/v1/*
                     ▼
              ┌──── nginx ────┐    host-based routing
              │   (reverse    │    public domain → exchange-httpd-go
              │    proxy)     │    admin domain  → ub-admin-main
              └───┬──────┬───┘
                  │      │
      ┌───────────▼┐  ┌──▼──────────────────────┐
      │ub-server   │  │ub-exchange-cli-main (Go) │
      │(PHP-FPM)   │  │ exchange-httpd :8000/8001│
      │Symfony 5.4 │  │ exchange-ws   (WS)       │
      │            │  │ exchange-engine           │
      └──┬──┬──┬───┘  └──┬──────┬────────────────┘
         │  │  │          │      │
         │  │  │  Shared  │      │  gRPC (candle-grpc network)
         │  │  │  MySQL   │      │  between Go services
         │  │  │  ◄───────┘      │
         │  │  │                 │
         │  │  └──── MQTT ──► EMQX broker ──► Client WebSockets
         │  │        publish    :1883/:8083   (tickers, order books,
         │  │        market     (auth via     trade books, user events)
         │  │        data       /api/v1/emqtt)
         │  │
         │  └──── RabbitMQ ──► ub-communicator-main (Go)
         │        publish       consumer workers (5 pool)
         │        email/sms     → SendGrid / Mailjet / Mailgun / SMTP
         │        messages      → Twilio SMS
         │                      → MongoDB audit log
         │
         └──── Redis 6.2 ── caching, order books (sorted sets),
                             session storage, live data, pub/sub
```

### Shared Infrastructure

| Service | Version | Purpose |
|---------|---------|---------|
| **MariaDB** | 10.2 | Primary data store (users, orders, trades, balances, currencies) |
| **Redis** | 6.2.2-alpine | Caching, order books (sorted sets), live data, session storage, pub/sub |
| **RabbitMQ** | 3.7 | Async messaging: ub-server/ub-exchange-cli → ub-communicator |
| **EMQX** | 3.0 (prod) / 4.0 (dev) | Real-time WebSocket push (tickers, order updates) |
| **MongoDB** | latest | Audit log for sent messages (ub-communicator only) |
| **Sentry** | — | Error tracking across all services |
| **nginx** | — | Reverse proxy, SSL termination, host-based routing |

## Cross-Service Integration

### Communication Patterns

| Pattern | From → To | Protocol | Details |
|---------|-----------|----------|---------|
| REST API | clients → ub-exchange-cli | HTTP `/api/v1/*` | All client-facing API (orders, trades, balances) via Gin |
| REST API | clients → ub-server | HTTP `/api/v1/*` | Auth, user management, crypto payments via Symfony |
| Shared DB | ub-exchange-cli ↔ ub-server | MySQL (GORM / Doctrine) | Both services read/write the same MySQL database |
| RabbitMQ | ub-exchange-cli → ub-communicator | AMQP topic exchange `messages` | Email/SMS notifications (async); Go CLI is the working publisher path |
| RabbitMQ | ub-server → ub-communicator | AMQP direct exchange `email_exchange` | ⚠️ **BROKEN** — exchange/type/routing mismatch with consumer (see Gotcha #11) |
| MQTT Pub | ub-server + ub-exchange-cli → EMQX → clients | MQTT topics `main/trade/*` | Real-time market data (ticker, order book, trades, klines) |
| gRPC | Go services (ws, httpd, engine) | Internal `candle-grpc` network | Inter-service Go communication |
| JWT | ub-server issues → ub-exchange-cli validates | HTTP `Authorization: Bearer` header | Shared auth via Lexik JWT Bundle |

### RabbitMQ Message Flow

```
ub-server (PHP)                          ub-communicator (Go)
 CommunicationManager                    rabbit-consumer
   └─ RabbitMQProducerService              └─ Consumer Service
       exchange: "email_exchange"              exchange: "messages" (topic)
       queue:    "email_queue_1"               queue:    "messages.command.send.consumer"
       type:     direct, durable               binding:  "messages.command.send"
       routing:  [] (empty)                    workers:  5 (channel-of-channels pool)
                                               autoAck:  true

ub-exchange-cli (Go)                      ↑ COMPATIBLE with consumer
 communication/queuemanager.go
   └─ PublishEmailOrSms()
       exchange: "messages" (topic)
       routing:  "messages.command.send"   ← matches consumer binding
       Also publishes kline data:
       exchange: "livedata" (topic)
       routing:  "livedata.event.kline-created"

Message payload: { receiver, subject, content, priority, type ("email"|"sms"), scheduledAt }
Note: PHP only publishes in "prod" environment and uses INCOMPATIBLE exchange config.
      Go CLI is the working path for email/SMS notifications to reach ub-communicator.
```

### MQTT Topic Structure

```
Public topics (all clients):
  main/trade/ticker/{pair}              — Price tickers (or just main/trade/ticker for all-pairs array)
  main/trade/order-book/{pair}          — Live order books (optional /{precision} suffix)
  main/trade/trade-book/{pair}          — Executed trades
  main/trade/chart/{timeFrame}/{pair}   — OHLC chart data (by time frame)
  main/trade/kline/{timeFrame}/{pair}   — K-line/candlestick data (by time frame)
  main/trade/change-price/{pair}        — Price changes (⚠️ event subscriber commented out)
  main/trade/market-price/{pair}        — Current market prices

Private topics (authenticated users):
  main/trade/user/{privateChannel}/open-orders/      — User's open orders
  main/trade/user/{privateChannel}/crypto-payments/  — User's payment status

Publishers: Both ub-server (PHP/EmqttManager) and ub-exchange-cli (Go/mqttmanager) publish
            to the same topics. Both use credentials mqtt_abbas:mqtt_abbas on emqtt:1883.

Auth: EMQX calls /api/v1/emqtt/login, /acl, /superuser — served by BOTH PHP and Go.
      Subscribers connect via WSS on port 8443. Public topics allow anonymous access.
      Private topics validated against user's privateChannelName via JWT.
```

## Database Schema Overview

All entities live in `ub-server-main/src/Exchange/` (Doctrine ORM). Both ub-server (Doctrine)
and ub-exchange-cli (GORM) share the same MySQL database.

| Entity | Bundle | Key Fields |
|--------|--------|------------|
| **User** | UserBundle | id, email, phone, roles, 2FA, KYC status |
| **UserBalance** | UserBundle | user, currency, available, frozen (Money embeddable) |
| **Order** | OrderBundle | user, pair, type (buy/sell), price, amount, status |
| **Trade** | TradeBundle | buyOrder, sellOrder, price, amount, fee |
| **Currency** | CryptoBundle | symbol, name, network, decimals, enabled |
| **PairCurrency** | CryptoBundle | baseCurrency, quoteCurrency, fees, limits |
| **CryptoPayment** | TransactionBundle | user, currency, amount, type (deposit/withdraw), status |
| **UserWithdrawAddress** | TransactionBundle | user, currency, address, label |

Key conventions:
- All monetary values use `Money` embeddable (precision handling)
- Materialized path tree for hierarchical orders
- MySQL JSON functions for flexible extra-info fields
- 250+ Doctrine migrations (append-only — **never edit existing migrations**)

## Code Style

### PHP (ub-server-main)
- Symfony 5.4 bundle structure with service-based DI
- Doctrine entities in `src/Exchange/<Bundle>/Entity/`
- Controllers return JSON via `JsonResponse`
- Service classes in `src/Exchange/<Bundle>/Services/`
- EventSubscribers for decoupled cross-cutting concerns
- Custom Doctrine types (e.g., `exchange_currency`)
- God-service refactoring: large services split into focused services + facade (see `refactoring-summary.md`)

### TypeScript/React (ub-admin, ub-client-cabinet)
- Container pattern: each page = `index.tsx` + `saga.ts` + `slice.ts` + `selectors.ts` + `types.ts`
- Redux-Saga for side effects, dynamic injection via `useInjectReducer()`/`useInjectSaga()`
- Singleton `ApiService` pattern for HTTP calls
- Styled-components + Material-UI for styling
- i18next (admin) / react-intl (client) for translations

### Flutter/Dart (ub-app-main)
- GetX-based MVC modules: `bindings/` + `controllers/` + `views/` + `providers/`
- Dio HTTP client with interceptors (retry, token refresh, logging)
- GetStorage for local persistence, FlutterSecureStorage for credentials
- Named routes via `GetPages` (40+ routes)
- ⚠️ **Pre-null-safety** — Dart SDK <3.0, no sound null safety

### Go (ub-exchange-cli, ub-communicator)
- Layered architecture: `internal/` for business logic, `cmd/` for entry points
- DI container pattern (sarulabs/di or custom)
- GORM for MySQL, go-redis for caching
- Structured logging via uber/zap

## Build and Test

### PHP Backend
```bash
cd ub-server-main
docker-compose up -d                          # Start all services (nginx, PHP, Go, DB, Redis, RabbitMQ, EMQX)
composer install                              # Install deps
bin/console doctrine:migrations:migrate       # Run migrations
vendor/bin/codeception run                    # Run tests (227 tests, 1777 assertions)
```

### Admin Panel
```bash
cd ub-admin-main
npm install                                   # or yarn
npm start                                     # Dev server (port 3000)
npm run build                                 # Production build
npm test                                      # Jest tests (⚠️ some suites may fail)
npm run lint                                  # ESLint
```

### Client Cabinet
```bash
cd ub-client-cabinet-main
yarn install                                  # Uses yarn
npm start                                     # Dev (IS_LOCAL=true)
npm run build                                 # Production build
npm test                                      # Jest tests (98% coverage threshold)
npm run lint                                  # ESLint
```

### Flutter App
```bash
cd ub-app-main
flutter pub get
flutter run                                   # Debug on device
./buildDevAPK.sh                              # Dev APK
./buildWeb-dev.sh                             # Dev web build
./buildAPK.sh                                 # Production APK
./buildWeb.sh                                 # Production web build
```

### Go Trading Engine
```bash
cd ub-exchange-cli-main
go build ./cmd/exchange-cli/                  # CLI tool
go build ./cmd/exchange-httpd/                # HTTP API server (:8000 public, :8001 admin)
go build ./cmd/exchange-ws/                   # WebSocket server
go build ./cmd/exchange-engine/               # Matching engine
```

### Go Communicator
```bash
cd ub-communicator-main
go build -mod=vendor ./cmd/rabbit-consumer/   # Uses vendored deps
./rabbit-consumer                             # Start RabbitMQ consumer
```

## Environment Setup

### Docker Compose (Full Stack)

The primary Docker Compose lives in `ub-server-main/`:

```bash
cd ub-server-main
docker-compose up -d    # Starts: nginx, PHP-FPM, 3 Go services, MariaDB, Redis, RabbitMQ, EMQX
```

| Compose file | Purpose |
|---|---|
| `docker-compose.yml` | Local development (all services, exposed ports) |
| `docker-compose-dev.yml` | Dev server (external env, EMQX v4, no nginx/DB) |
| `docker-compose-prod.yml` | Production (SSL, localhost-only DB/Redis, deploy scripts) |

The communicator has its own Compose in `ub-communicator-main/`:
```bash
cd ub-communicator-main
docker-compose up -d    # Starts: MongoDB, Go consumer (connects to ub-server's RabbitMQ via external network)
```

### Key Ports (Local Dev)

| Port | Service |
|------|---------|
| 8081 | nginx → PHP API |
| 8082 | nginx (secondary) |
| 8000 | exchange-httpd-go (public API) |
| 8001 | exchange-httpd-go (admin API) |
| 3308 | MariaDB |
| 6379 | Redis |
| 5672 | RabbitMQ |
| 1883 | EMQX MQTT |
| 8083 | EMQX WebSocket |
| 27017 | MongoDB (communicator) |

### Environment Variables

- PHP: `app/config/parameters.yml` (Symfony), `.env` files
- Go exchange-cli: environment variables or config files
- Go communicator: `config/config.yaml` + env vars with `UBCOMMUNICATOR_` prefix (Viper)
- Frontend apps: `.env` files (see `.env.example` in each sub-project)

## API Contract Summary

### Endpoint Structure

All APIs versioned under `/api/v1/`:

| Route prefix | Service | Purpose |
|---|---|---|
| `/api/v1/auth/*` | ub-server + ub-exchange-cli | Login, register, forgot-password, 2FA |
| `/api/v1/order/*` | ub-exchange-cli | Create, cancel, list orders |
| `/api/v1/trade/*` | ub-exchange-cli | Trade history |
| `/api/v1/currencies/*` | ub-exchange-cli | Pairs, fees, statistics |
| `/api/v1/user-balance/*` | ub-exchange-cli | Balances, auto-exchange |
| `/api/v1/user/*` | ub-exchange-cli | Profile, 2FA, password, SMS |
| `/api/v1/crypto-payment/*` | ub-exchange-cli | Deposit, withdraw, cancel |
| `/api/v1/emqtt/*` | ub-server + ub-exchange-cli | MQTT auth (login, ACL, superuser) — both backends serve these |
| `/tv/api/v1/*` | ub-server | TradingView charting integration |

Admin API uses **host-based routing** (admin subdomain → port 8001).

### Authentication Flow

1. Client sends `POST /api/v1/auth/login` with credentials
2. ub-server issues JWT token (Lexik JWT Bundle)
3. Client includes `Authorization: Bearer <JWT>` on subsequent requests
4. ub-exchange-cli validates JWT via `authService.GetUser(token)`
5. Token invalidated if user changes password or enables 2FA after token issuance

### Response Format
```json
{ "status": "success|error", "message": "...", "data": {...}, "token": "..." }
```
Error codes: `401` (auth), `422` (validation), `500` (server)

## Deployment Overview

### CI/CD (GitLab CI)

Both `ub-server-main` and `ub-communicator-main` use `.gitlab-ci.yml`:

| Branch | Pipeline |
|--------|----------|
| `develop` | Build Docker images → deploy to dev server → Telegram notification |
| `merge_requests` | Run Codeception tests (ub-server only, with MariaDB + Redis services) |
| `master` | SSH deploy to production → run `deploy.sh` → Telegram notification |

Production deployment (`ub-server-main/deploy.sh`):
```bash
# Clears caches, metadata, runs migrations
docker-compose -f docker-compose-prod.yml exec exchange-app php bin/console c:c --env=prod
docker-compose -f docker-compose-prod.yml exec exchange-app php bin/console doctrine:migrations:migrate --no-interaction --env=prod
```

### Docker Images

| Sub-project | Base image | Dockerfile(s) |
|---|---|---|
| ub-server-main | php:8.2-fpm | `.docker/php/Dockerfile`, `.docker/nginx/Dockerfile`, `.docker/go/Dockerfile` |
| ub-admin-main | node:18 | `DockerfileProd` |
| ub-client-cabinet-main | node:18 | `Dockerfile`, `DockerfileProd` |
| ub-app-main | Flutter 2.10 | `Dockerfiledev`, `Dockerfileprod`, `Dockerfileapkprod` |
| ub-communicator-main | golang:1.24 | `.docker/go/Dockerfile.dev`, `.docker/go/Dockerfile.prod` |

## Conventions

### API Patterns
- All APIs are versioned under `/api/v1/`
- Admin API uses host-based routing (admin subdomain)
- JWT authentication via `Authorization: Bearer <token>`
- Standard response: `{ status, message, data, token }`
- Error codes: 401 (auth), 422 (validation), 500 (server)

### Database
- All monetary values use `Money` embeddable (precision handling)
- Materialized path tree for hierarchical orders
- MySQL JSON functions for flexible extra-info fields
- 250+ Doctrine migrations (append-only, never edit existing)

### Real-time
- MQTT topics follow `main/trade/{channel}/{pair}` pattern
- Authorized vs unauthorized MQTT clients for different data access
- EMQX authenticates clients via HTTP callback to `/api/v1/emqtt/*`

### Security
- Google reCAPTCHA v2 on auth endpoints
- Google Authenticator 2FA support
- RSA/AES encryption for sensitive client data
- Biometric auth on mobile (fingerprint/face)

## Upgrade Priority (Legacy Debt)

Critical upgrades needed (in priority order):

1. **Dart SDK 2.x → 3.x** with null-safety migration (ub-app-main is pre-null-safety — blocks all modernization)
2. **React 17 → 18** (ub-admin-main) — ub-client-cabinet already at React 18
3. **Material-UI v4 → v5** (both frontends)
4. **Symfony 5.4 → 6.4 LTS** (ub-server-main)
5. **Webpack 4 → 5** (ub-client-cabinet-main)
6. **mailgun-go v2 → v4** (ub-communicator — v2 is archived) and **gomail.v2 → go-mail** (abandoned since 2016)
7. **axios 0.21 → 1.x** (ub-client-cabinet — 0.21 has known CVEs)
8. **Credentials in config files → environment variables / secrets manager**

> **Already completed:** Go services were upgraded (exchange-cli to Go 1.22, communicator to Go 1.24).
> PHP is already at 8.1+/Symfony 5.4. TypeScript is already at 5.4.5 in both frontends.

## Gotchas & Non-Obvious Behaviors

> **For AI agents:** Read this section before making changes. These are known
> behaviors that may look like bugs but are intentional (or known issues
> that are tracked separately).

### 1. Message Type Must Be Uppercase
`CreateMessage()` normalizes `message.Type` to uppercase via `strings.ToUpper()`.
Constants are `"EMAIL"` and `"SMS"`. Producers can send any case — the normalization
handles it. But if you're comparing types manually, always compare against
the uppercase constants.

### 2. autoAck=true on RabbitMQ — Messages Can Be Lost
Messages are auto-acknowledged on delivery from RabbitMQ, NOT after successful
processing. If the service crashes while processing, that message is lost.
This is a known trade-off for throughput. Changing to manual ack requires also
implementing retry/dead-letter-queue logic.

### 3. Worker Pool Uses Channel-of-Channels Pattern
Workers register availability by sending their work channel onto a shared
`workerChannel`. The dispatcher picks the first available worker. If all workers
are busy, the dispatcher blocks. Don't change this to a simple goroutine-per-message
model without understanding the back-pressure implications.

### 4. MailService Is an Intentional Thin Wrapper
`messaging.MailService` wraps `platform.MailerClient` with identical signature.
It exists as an extension point for future middleware (logging, metrics, retry).
Don't remove it — add cross-cutting concerns here instead of in the platform layer.

### 5. Sentry Only Active in Production
`captureError()` checks `l.env != "prod"` and returns early. Errors in dev/test
are only logged to zap, not sent to Sentry. This is intentional — don't remove
the check.

### 6. Config Precedence: Env Vars Override YAML
Viper reads `config.yaml` first, then environment variables with prefix
`UBCOMMUNICATOR_` override any yaml values. Dots become underscores:
`rabbitmq.dsn` → `UBCOMMUNICATOR_RABBITMQ_DSN`. If a value appears empty
in config.yaml, check environment variables.

### 7. MongoDB Collection Name Is Hardcoded
All audit logs go to the `"messages"` collection (hardcoded in
`pkg/repository/messageRepository.go`). The database name is configurable
via `mongodb.name` config key, but the collection name is not.

### 8. Subject Prefix Is Inconsistent Across Providers
Mailjet and SendGrid prepend `[UNITEDBIT]` to subjects that don't already
contain `[`. Mailgun and SMTP do NOT add any prefix. This means emails
look different depending on the configured provider.

### 9. SMS Uses Twilio HTTP API Directly (Not SDK)
The SMS service calls Twilio's REST API via raw HTTP POST with Basic Auth,
not the official Twilio Go SDK. The endpoint URL is constructed from the
Account SID: `https://api.twilio.com/2010-04-01/Accounts/{SID}/Messages.json`.

### 10. HTTP Client Response Format
All `HttpClient` methods return `(body []byte, headers http.Header, statusCode int, error)`.
On connection errors, statusCode defaults to 0 (after fix) — don't trust
statusCode unless error is nil.

### 11. RabbitMQ Exchange/Queue Name Mismatch Between Producer and Consumer
ub-server publishes to exchange `"email_exchange"` / queue `"email_queue_1"` (direct),
while ub-communicator consumes from exchange `"messages"` / queue
`"messages.command.send.consumer"` (topic). These must be reconciled via
RabbitMQ config or both sides updated if changing messaging patterns.

### 12. ub-server RabbitMQ Only Publishes in Production
`RabbitMQProducerService` only publishes in `prod` environment. In dev/test,
no messages reach ub-communicator. If testing the full email flow locally,
you must either change this check or publish manually.

### 13. Mailgun API Parameters Are Swapped (ub-communicator BUG)
In `pkg/platform/mail.go:59`, the code calls `mailgun.NewMailgun(apiKey, domain)`
but the SDK signature is `NewMailgun(domain, apiKey string)`. The domain is used
as the API key and vice versa. **All Mailgun emails fail with auth errors.**
Fix: swap to `mailgun.NewMailgun(domain, apiKey)`.

### 14. Sentry Is Silently Broken in Communicator Production
`platform.EnvConfigKey` reads `"wallet.environment"` but `config.yaml` uses
`"communicator.environment"`. `GetEnv()` always returns `""`, so `captureError()`
never sends to Sentry (the `env != "prod"` check always passes). Workaround:
set `UBCOMMUNICATOR_WALLET_ENVIRONMENT=prod` as an env var. Fix: change
the constant to `"communicator.environment"`.

### 15. No Graceful Shutdown in Communicator
`main.go` passes `context.Background()` to `Consume()` — SIGTERM/SIGINT are
ignored. Docker `stop` waits its timeout then SIGKILL's the process. Worker
pool shutdown logic exists but is unreachable via OS signals.

## Spec-Driven Development Guide

### Which Sub-Project to Modify

| Feature area | Primary project | May also need |
|---|---|---|
| User auth, registration, 2FA | ub-server-main | ub-exchange-cli-main (JWT validation) |
| Trading (orders, matching) | ub-exchange-cli-main | ub-server-main (Doctrine entities) |
| Admin panel features | ub-admin-main | ub-server-main (admin API endpoints) |
| Client dashboard UI | ub-client-cabinet-main | ub-exchange-cli-main (API) |
| Mobile app features | ub-app-main | ub-exchange-cli-main (API) |
| Email/SMS notifications | ub-communicator-main | ub-server-main (publisher) |
| Real-time data (WebSocket) | ub-server-main (MQTT publish) | ub-exchange-cli-main (WS server) |
| New currency/pair | ub-server-main (entity + migration) | ub-exchange-cli-main (GORM model) |
| New database table | ub-server-main (Doctrine migration) | ub-exchange-cli-main (GORM model if shared) |

### Cross-Service Change Checklist

When modifying shared entities (User, Order, Trade, Balance, Currency):
1. Update Doctrine entity in ub-server-main
2. Create Doctrine migration (`bin/console doctrine:migrations:diff`)
3. Update GORM model in ub-exchange-cli-main (if the Go service uses this table)
4. Verify both services can read/write the updated schema
5. Update API response DTOs in both services if the field is exposed
