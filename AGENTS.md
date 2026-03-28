# UnitedBit Exchange Platform — Workspace Guidelines

## Overview

This is a **cryptocurrency exchange platform** monorepo with 6 sub-projects:

| Project | Tech | Purpose |
|---------|------|---------|
| `ub-server-main` | PHP 8.2 / Symfony 5.4 LTS / Doctrine ORM 2.14 / MariaDB 10.2 | Backend REST API (13 bundles) |
| `ub-admin-main` | React 16 / TypeScript 3.9 / Redux-Saga / Material-UI 4 | Admin panel SPA |
| `ub-client-cabinet-main` | React 17 / TypeScript 4.0 / Redux-Saga / Webpack 4 | Client trading dashboard |
| `ub-app-main` | Flutter / Dart 2.11 / GetX / Dio | Mobile + Web app |
| `ub-exchange-cli-main` | Go 1.13 / Gin / GORM / Redis / gRPC | Trading engine, CLI commands, HTTP API |
| `ub-communicator-main` | Go 1.13 / RabbitMQ / MongoDB | Email/SMS notification service |

## Architecture

```
Client Apps (ub-app, ub-client-cabinet)
        ↓
   ub-server (PHP API) ←→ ub-exchange-cli (Go trading engine)
        ↓                        ↓
   ub-admin (Admin panel)    RabbitMQ → ub-communicator (Email/SMS)
        ↓                        ↓
      MySQL               Redis / MQTT (real-time)
```

### Shared Infrastructure
- **MySQL** — Primary data store (users, orders, trades, balances, currencies)
- **Redis** — Caching, order books (sorted sets), live data, session storage
- **RabbitMQ** — Async messaging between services
- **MQTT/EMQX** — Real-time WebSocket push to clients (tickers, order updates)
- **Sentry** — Error tracking across all services
- **Docker** — All services containerized

## Code Style

### PHP (ub-server-main)
- Symfony 3.4 bundle structure with service-based DI
- Doctrine entities in `src/<Bundle>/Entity/`
- Controllers return JSON via `JsonResponse`
- Service classes in `src/<Bundle>/Service/`
- EventSubscribers for decoupled cross-cutting concerns
- Custom Doctrine types (e.g., `exchange_currency`)

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

### Go (ub-exchange-cli, ub-communicator)
- Layered architecture: `internal/` for business logic, `cmd/` for entry points
- DI container pattern (sarulabs/di or custom)
- GORM for MySQL, go-redis for caching
- Structured logging via uber/zap

## Build and Test

### PHP Backend
```bash
cd ub-server-main
docker-compose up -d          # Start services
composer install              # Install deps
bin/console doctrine:migrations:migrate  # Run migrations
vendor/bin/codeception run    # Run tests
```

### Admin Panel
```bash
cd ub-admin-main
npm install
npm start                     # Dev server
npm run build                 # Production build
```

### Client Cabinet
```bash
cd ub-client-cabinet-main
npm install
npm start                     # Dev (IS_LOCAL=true)
npm run build                 # Production
npm test                      # Jest tests (98% coverage threshold)
```

### Flutter App
```bash
cd ub-app-main
flutter pub get
flutter run                   # Debug on device
./buildDevAPK.sh              # Dev APK
./buildWeb-dev.sh             # Dev web
./buildAPK.sh                 # Production APK
./buildWeb.sh                 # Production web
```

### Go Services
```bash
cd ub-exchange-cli-main
go build ./cmd/exchange-cli/
./exchange-cli <command>      # Run CLI commands

cd ub-communicator-main
go build ./cmd/rabbit-consumer/
./rabbit-consumer             # Start consumer
```

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
- MQTT topics follow `main/trade/{channel}` pattern
- Authorized vs unauthorized MQTT clients for different data access
- MQTT messages use custom cipher for auth

### Security
- Google reCAPTCHA v2 on auth endpoints
- Google Authenticator 2FA support
- RSA/AES encryption for sensitive client data
- Biometric auth on mobile (fingerprint/face)

## Upgrade Priority (Legacy Debt)

Critical upgrades needed (in priority order):
1. **Go 1.13 → Go 1.22+** (both Go services — security CVEs)
2. **PHP 7.4 → PHP 8.2+** (EOL since Nov 2022)
3. **Symfony 3.4 → Symfony 6.4 LTS** (EOL since Nov 2021)
4. **React 16/17 → React 18+** (both frontends)
5. **Dart SDK 2.x → 3.x** with null safety migration
6. **Node dependencies** — axios, Material-UI v4→v5, etc.
7. **Credentials in config files → environment variables/secrets manager**

## Gotchas & Non-Obvious Behaviors

> **For AI agents:** Read this section before making changes. These are known
> behaviors that may look like bugs but are intentional (or known issues
> that are tracked separately).

### 1. Message Type Must Be Uppercase
`CreateMessage()` normalizes `message.Type` to uppercase via `strings.ToUpper()`.
Constants are `"EMAIL"` and `"SMS"`. Producers can send any case (lowercase,
uppercase, mixed) — the normalization handles it. But if you're comparing
types manually, always compare against the uppercase constants.

### 2. autoAck=true on RabbitMQ — Messages Can Be Lost
Messages are auto-acknowledged on delivery from RabbitMQ, NOT after successful
processing. If the service crashes while a worker is processing a message,
that message is lost. This is a known trade-off for throughput. Changing to
manual ack requires also implementing retry/dead-letter-queue logic.

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
