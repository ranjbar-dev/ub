# ub-communicator

Go-based notification microservice for the UnitedBit cryptocurrency exchange platform.
Consumes messages from RabbitMQ and delivers them via email (SendGrid, Mailjet, Mailgun, SMTP)
or SMS (Twilio). Stores delivery audit logs in MongoDB.

## Architecture

```
RabbitMQ (topic exchange: "messages")
  └── Queue: messages.command.send.consumer
        └── Consumer Service (pkg/consumer/)
              └── Worker Pool (configurable, default: 5 workers)
                    │
                    ├── type=EMAIL → Mail Provider Factory
                    │                 ├── SendGrid (API)
                    │                 ├── Mailjet  (API)
                    │                 ├── Mailgun  (API)
                    │                 └── SMTP     (gomail)
                    │
                    └── type=SMS  → Twilio HTTP API
                    │
                    └── MongoDB (audit log: every attempt persisted)
```

## Prerequisites

- **Go 1.24+**
- **RabbitMQ** — AMQP message broker
- **MongoDB** — audit log storage
- **Docker & Docker Compose** — for containerized deployment

## Quick Start

```bash
# 1. Copy env file and fill in values
cp .env.example .env

# 2. Build
go build -mod=vendor cmd/rabbit-consumer/main.go

# 3. Run (requires RabbitMQ + MongoDB)
./main

# Docker — development
docker-compose up -d
docker exec communicator-app go build -mod=vendor cmd/rabbit-consumer/main.go

# Docker — production
docker-compose -f docker-compose-prod.yml up -d --build
```

## Configuration

Configuration is loaded from `config/config.yaml` with environment variable overrides.

**Env var format:** `UBCOMMUNICATOR_<SECTION>_<KEY>` (dots become underscores)

| Config Key | Env Variable | Default | Description |
|-----------|-------------|---------|-------------|
| `communicator.environment` | `UBCOMMUNICATOR_COMMUNICATOR_ENVIRONMENT` | `dev` | `dev`, `test`, or `prod` |
| `consumer.worker_count` | `UBCOMMUNICATOR_CONSUMER_WORKER_COUNT` | `5` | Worker pool size |
| `rabbitmq.dsn` | `UBCOMMUNICATOR_RABBITMQ_DSN` | — | RabbitMQ AMQP URI |
| `rabbitmq.exchange` | `UBCOMMUNICATOR_RABBITMQ_EXCHANGE` | `messages` | Exchange name |
| `rabbitmq.queue_name` | `UBCOMMUNICATOR_RABBITMQ_QUEUE_NAME` | `messages.command.send.consumer` | Queue name |
| `rabbitmq.binding` | `UBCOMMUNICATOR_RABBITMQ_BINDING` | `messages.command.send` | Routing key |
| `mongodb.dsn` | `UBCOMMUNICATOR_MONGODB_DSN` | — | MongoDB connection string |
| `mongodb.name` | `UBCOMMUNICATOR_MONGODB_NAME` | `ubMessages` | Database name |
| `mailer_broker` | `UBCOMMUNICATOR_MAILER_BROKER` | `smtp` | `smtp` / `sendgrid` / `mailjet` / `mailgun` |
| `mail.name` | `UBCOMMUNICATOR_MAIL_NAME` | `UNITEDBIT` | Sender display name |
| `mail.from_address` | `UBCOMMUNICATOR_MAIL_FROM_ADDRESS` | `no-reply@unitedbit.com` | Sender email |
| `sms.account_sid` | `UBCOMMUNICATOR_SMS_ACCOUNT_SID` | — | Twilio Account SID |
| `sms.auth_token` | `UBCOMMUNICATOR_SMS_AUTH_TOKEN` | — | Twilio Auth Token |
| `sms.from` | `UBCOMMUNICATOR_SMS_FROM` | — | Twilio sender phone |
| `sentry.dsn` | `UBCOMMUNICATOR_SENTRY_DSN` | — | Sentry DSN (empty = disabled) |
| `logging.file_path` | `UBCOMMUNICATOR_LOGGING_FILE_PATH` | `stdout` | Log file path |

Provider-specific keys: see `config/config.yaml` and `.env.example` for the full list.

## Message Format

Messages are consumed from RabbitMQ as JSON:

```json
{
  "type": "email",
  "receiver": "user@example.com",
  "subject": "Welcome",
  "content": "<h1>Hello</h1>",
  "priority": 1,
  "scheduledAt": ""
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `type` | string | Yes | `"email"` or `"sms"` (case-insensitive, normalized to uppercase) |
| `receiver` | string | Yes | Email address or E.164 phone number (`+1234567890`) |
| `subject` | string | No | Email subject line (not used for SMS) |
| `content` | string | Yes | HTML body (email) or text body (SMS) |
| `priority` | int | No | Stored but not currently used |
| `scheduledAt` | string | No | Stored but not currently used |

## Project Structure

```
cmd/rabbit-consumer/main.go       Entry point
config/                            Configuration (Viper + YAML)
pkg/consumer/                      RabbitMQ consumer + worker pool
pkg/messaging/                     Message routing + delivery orchestration
pkg/platform/                      Infrastructure adapters (HTTP, logging, mail, DB, MQ)
pkg/repository/                    MongoDB persistence
pkg/di/                            Dependency injection container
.docker/                           Dockerfiles (dev + prod)
```

See `AGENTS.md` for detailed file-by-file documentation, all exported types,
DI wiring, known issues, and spec-driven development guide.

## Monitoring

- **Sentry** — Error reporting in production (configure `sentry.dsn`)
- **Logs** — Structured JSON via zap (stdout + optional file)
- **MongoDB** — Audit trail of all delivery attempts in `messages` collection

## Development

```bash
# Run tests
go test ./...

# Static analysis
go vet ./...

# Build binary
go build -mod=vendor cmd/rabbit-consumer/main.go
```

## Docker Compose Files

| File | Purpose |
|------|---------|
| `docker-compose.yml` | Local development (code volume-mounted) |
| `docker-compose-dev.yml` | Dev server deployment |
| `docker-compose-prod.yml` | Production deployment (pre-built binary) |

All compose files include MongoDB and the Go app, connected to the external
`ub-server_rabbit` network for RabbitMQ access.
