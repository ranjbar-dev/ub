# UnitedBit Exchange Platform

A full-stack cryptocurrency exchange platform built as a monorepo with 6 specialized sub-projects spanning PHP, Go, TypeScript/React, and Flutter/Dart.

## Sub-Projects

| Project | Description | Tech Stack |
|---------|-------------|------------|
| [`ub-server-main`](ub-server-main/) | Backend REST API — 13 Symfony bundles, Doctrine ORM, JWT auth | PHP 8.1+ / Symfony 5.4 / MariaDB |
| [`ub-admin-main`](ub-admin-main/) | Admin panel SPA — user management, order oversight, system config | React 17 / TypeScript 5.4 / Material-UI 4 |
| [`ub-client-cabinet-main`](ub-client-cabinet-main/) | Client trading dashboard — portfolio, orders, deposits/withdrawals | React 18 / TypeScript 5.4 / Webpack 4 |
| [`ub-app-main`](ub-app-main/) | Mobile + Web app — iOS, Android, Web via Flutter | Flutter 2.x / Dart 2.11 / GetX |
| [`ub-exchange-cli-main`](ub-exchange-cli-main/) | Trading engine — order matching, HTTP API, WebSocket, CLI tools | Go 1.22 / Gin / GORM / gRPC |
| [`ub-communicator-main`](ub-communicator-main/) | Notification service — email (SendGrid/Mailgun/Mailjet) & SMS (Twilio) | Go 1.24 / RabbitMQ / MongoDB |

## Architecture

```
  Flutter App / React SPA (clients)
           │
           ▼
        nginx (reverse proxy, SSL, host-based routing)
        ┌──────────┬──────────────────┐
        ▼          ▼                  ▼
   ub-server    ub-exchange-cli   ub-admin
   (PHP API)    (Go HTTP/WS/      (React SPA)
                 Engine)
        │          │
        ├── Shared MySQL (MariaDB 10.2) ──┤
        ├── Redis 6.2 (cache, order books)│
        ├── Centrifugo (WebSocket real-time push) ─┘
        │
        └── RabbitMQ ──► ub-communicator (Go)
            (async)       Email / SMS delivery
                          MongoDB audit log
```

See [`AGENTS.md`](AGENTS.md) for detailed architecture, API contracts, and cross-service integration patterns.

## Quick Start

### Prerequisites
- Docker & Docker Compose
- PHP 8.1+ & Composer (for ub-server)
- Node.js 18+ & npm/yarn (for frontends)
- Go 1.22+ (for Go services)
- Flutter SDK 2.x (for mobile app)

### Start the Full Stack (Docker)
```bash
# 1. Start core infrastructure (DB, Redis, RabbitMQ, Centrifugo, PHP, Go services, nginx)
cd ub-server-main
docker-compose up -d

# 2. Start communicator (connects to ub-server's RabbitMQ)
cd ../ub-communicator-main
docker-compose up -d

# 3. Start frontend dev servers
cd ../ub-admin-main && npm install && npm start         # :3000
cd ../ub-client-cabinet-main && yarn install && npm start  # :3000

# 4. Start mobile app
cd ../ub-app-main && flutter pub get && flutter run
```

### Key Ports (Local Development)
| Port | Service |
|------|---------|
| 8081 | nginx → PHP API |
| 8000 | Go HTTP API (public) |
| 8001 | Go HTTP API (admin) |
| 3308 | MariaDB |
| 6379 | Redis |
| 5672 | RabbitMQ |
| 8000 | Centrifugo HTTP API |

## Project Documentation

- **[`AGENTS.md`](AGENTS.md)** — Comprehensive workspace guidelines (tech stack, architecture, conventions, gotchas, spec-driven development guide)
- **[`refactoring-summary.md`](docs/refactoring-summary.md)** — History of god-service refactoring (UserBalanceService, ExternalExchangeOrderService)
- Each sub-project has its own `AGENTS.md` and/or `README.md` with project-specific details

## Development Workflow

1. **Read** [`AGENTS.md`](AGENTS.md) before making changes
2. **Identify** which sub-project(s) to modify (see "Which Sub-Project to Modify" in AGENTS.md)
3. **Check** cross-service dependencies (shared MySQL schema, RabbitMQ contracts, JWT auth)
4. **Build & test** in the affected sub-project(s)
5. **Verify** integration if modifying shared entities or API contracts

## CI/CD

GitLab CI pipelines per sub-project:
- **`develop`** → Build + deploy to dev server + Telegram notification
- **`merge_requests`** → Automated tests (Codeception for PHP)
- **`master`** → Deploy to production via SSH + Telegram notification

## Contributing

1. Create a feature branch from `develop`
2. Make changes following the conventions in [`AGENTS.md`](AGENTS.md)
3. Ensure tests pass in affected sub-projects
4. Open a merge request targeting `develop`
5. After review and CI passes, merge to `develop` for dev deployment
