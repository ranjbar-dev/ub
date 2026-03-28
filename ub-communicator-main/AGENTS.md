# Communicator — ub-communicator-main

## Overview
Go-based notification microservice that consumes messages from RabbitMQ and delivers them via email (SendGrid, Mailjet, Mailgun, SMTP) or SMS (Twilio). Stores delivery audit logs in MongoDB.

## Stack
| Component | Version | Status |
|-----------|---------|--------|
| Go | 1.24.0 (go.mod) / 1.24 (Docker) | ✅ Current |
| RabbitMQ client | rabbitmq/amqp091-go v1.10.0 | ✅ Current |
| MongoDB driver | go.mongodb.org/mongo-driver v1.17.2 | ✅ Current |
| Logging | go.uber.org/zap v1.27.0 | ✅ Current |
| Config | spf13/viper v1.19.0 | ✅ Current |
| Sentry | getsentry/sentry-go v0.31.1 | ✅ Current |
| SendGrid | sendgrid-go v3.16.0 | ✅ Current |
| Mailgun | mailgun-go v2.0.0+incompatible | ⚠️ Archived — needs v4 migration |
| Mailjet | mailjet-apiv3-go/v3 v3.2.0 | ✅ Current (v4 available) |
| SMTP | gomail.v2 | ⚠️ Abandoned 2016 — functional but unmaintained |

## Architecture
```
cmd/rabbit-consumer/main.go    → Entry point (long-running consumer)
pkg/di/container.go            → Custom DI container (singleton lazy-init)
pkg/consumer/service.go        → RabbitMQ consumer (exchange/queue/bind/consume)
pkg/consumer/pool.go           → Worker pool dispatcher (5 workers)
pkg/consumer/worker.go         → Worker goroutine (processes Work items)
pkg/messaging/service.go       → Message orchestration (routes EMAIL/SMS)
pkg/messaging/mailService.go   → Email delivery delegation
pkg/messaging/smsService.go    → Twilio SMS delivery via HTTP
pkg/messaging/repository.go    → Message struct + Repository interface
pkg/platform/config.go         → Viper config wrapper with env var support
pkg/platform/http.go           → HTTP client (GET/POST/form/basic-auth)
pkg/platform/logger.go         → Zap + Sentry logging
pkg/platform/mail.go           → Mail provider factory (4 providers)
pkg/platform/mailgunmail.go    → Mailgun implementation
pkg/platform/mailjetmail.go    → Mailjet implementation
pkg/platform/sendgridmail.go   → SendGrid implementation
pkg/platform/smtpmail.go       → SMTP/gomail implementation
pkg/platform/mongo.go          → MongoDB client initialization
pkg/platform/rabbitmq.go       → RabbitMQ connection pool (mutex, lazy connect)
pkg/repository/messageRepository.go → MongoDB message persistence
config/config.go               → Viper setup (reads config.yaml, env prefix: ubcommunicator)
config/config.yaml             → All configuration values (⚠️ contains secrets)
```

## Message Flow
1. RabbitMQ exchange `messages` (topic) → queue `messages.command.send.consumer`
2. Consumer dispatches to 5-worker pool via channel-based dispatcher
3. Workers call `messaging.Service.Send()` which routes by type (EMAIL/SMS)
4. Email goes through configured provider (smtp/sendgrid/mailjet/mailgun)
5. SMS goes through Twilio HTTP API
6. Result (success/fail) stored in MongoDB `messages` collection

## Build & Run
```bash
# Build
go build -mod=vendor cmd/rabbit-consumer/main.go

# Run (requires RabbitMQ + MongoDB)
./main

# Docker dev
docker-compose up -d

# Docker prod
docker-compose -f docker-compose-prod.yml up -d --build
```

## Configuration
- File: `config/config.yaml`
- Env prefix: `UBCOMMUNICATOR_` (e.g., `UBCOMMUNICATOR_SENDGRID_API_KEY`)
- Key: `mailer_broker` selects email provider: `smtp`, `sendgrid`, `mailjet`, `mailgun`

## Key Interfaces
- `consumer.Service` — `Consume()` (blocks forever consuming RabbitMQ)
- `consumer.Pool` — `StartDispatcher(workerCount int) Collector`
- `messaging.Service` — `Send(message)`, `CreateMessage(data []byte)`
- `messaging.MailService` — `Send(subject, receiver, content) (bool, error)`
- `messaging.SmsService` — `Send(subject, receiver, content) (bool, error)`
- `messaging.Repository` — `NewMessage(message *Message) error`
- `platform.MailerClient` — `Send(subject, receiver, content) (bool, error)`
- `platform.RabbitMqClient` — `GetChannel() (*amqp.Channel, error)`
- `platform.HttpClient` — `HttpGet`, `HttpPost`, `HttpPostForm`, `BasicAuth`
- `platform.Configs` — `GetString`, `GetInt`, `GetBool`, etc.
- `platform.Logger` — `Info`, `Warn`, `Error`, `Panic`, `Fatal`

## Known Issues
- **Hardcoded secrets** in config.yaml (API keys, passwords, tokens) — use env vars
- **context.TODO()** in MongoDB operations (no timeouts) — use context.WithTimeout
- **No graceful shutdown** — `<-forever` blocks without signal handling
- **No health check endpoint** — no way to monitor service health
- **Sentry re-initialized** per error in captureError() — should init once in NewLogger
- **InsecureSkipVerify: true** in SMTP TLS config — MitM risk
- **autoAck=true** on RabbitMQ consume — messages lost on crash
- **No reconnection logic** for RabbitMQ channel failures
- **mailgun-go v2 archived** — needs migration to v4 (breaking API changes)
- **gomail.v2 abandoned** (2016) — functional but no security patches
- **Telegram bot token** hardcoded in .gitlab-ci.yml

## Testing Guidelines
- Unit tests: mock interfaces (MailerClient, HttpClient, Repository, etc.)
- Integration tests: use testcontainers for MongoDB and RabbitMQ
- Test message routing (EMAIL → mail service, SMS → twilio)
- Test each mail provider Send() with mock HTTP
- Test worker pool dispatch and shutdown

## Conventions
- All services use interface-based design for testability
- DI container in `pkg/di/` with lazy singleton initialization
- Platform integrations isolated in `pkg/platform/`
- Business logic in `pkg/messaging/` and `pkg/consumer/`
- Factory pattern for mail providers in `platform.NewMailerClient()`
