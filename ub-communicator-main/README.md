# ub-communicator

Go-based notification microservice that consumes messages from RabbitMQ and delivers them via email (SendGrid, Mailjet, Mailgun, SMTP) or SMS (Twilio). Stores delivery audit logs in MongoDB.

## Architecture

```
RabbitMQ (messages exchange, topic)
  └── messages.command.send.consumer queue
        └── Consumer Service
              └── Worker Pool (configurable, default: 5 workers)
                    ├── Email → Mail Provider (SendGrid/Mailjet/Mailgun/SMTP)
                    └── SMS → Twilio HTTP API
                          └── MongoDB (audit log)
```

## Prerequisites

- Go 1.24+
- RabbitMQ
- MongoDB
- Docker & Docker Compose (for containerized deployment)

## Quick Start

```bash
# Clone and build
go build -mod=vendor cmd/rabbit-consumer/main.go

# Run (requires RabbitMQ + MongoDB)
./main

# Docker development
docker-compose up -d

# Docker production
docker-compose -f docker-compose-prod.yml up -d --build
```

## Configuration

Configuration is loaded from `config/config.yaml` and can be overridden with environment variables.

**Env var format:** `UBCOMMUNICATOR_<SECTION>_<KEY>` (dots become underscores)

| Config Key | Env Variable | Description |
|-----------|-------------|-------------|
| `communicator.environment` | `UBCOMMUNICATOR_COMMUNICATOR_ENVIRONMENT` | `dev`, `test`, or `prod` |
| `rabbitmq.dsn` | `UBCOMMUNICATOR_RABBITMQ_DSN` | RabbitMQ connection string |
| `rabbitmq.exchange` | `UBCOMMUNICATOR_RABBITMQ_EXCHANGE` | Exchange name (default: `messages`) |
| `rabbitmq.queue_name` | `UBCOMMUNICATOR_RABBITMQ_QUEUE_NAME` | Queue name |
| `rabbitmq.binding` | `UBCOMMUNICATOR_RABBITMQ_BINDING` | Routing key |
| `mongodb.dsn` | `UBCOMMUNICATOR_MONGODB_DSN` | MongoDB connection string |
| `mongodb.name` | `UBCOMMUNICATOR_MONGODB_NAME` | Database name |
| `mailer_broker` | `UBCOMMUNICATOR_MAILER_BROKER` | `smtp`, `sendgrid`, `mailjet`, or `mailgun` |
| `consumer.worker_count` | `UBCOMMUNICATOR_CONSUMER_WORKER_COUNT` | Worker pool size (default: 5) |
| `sms.account_sid` | `UBCOMMUNICATOR_SMS_ACCOUNT_SID` | Twilio Account SID |
| `sms.auth_token` | `UBCOMMUNICATOR_SMS_AUTH_TOKEN` | Twilio Auth Token |
| `sentry.dsn` | `UBCOMMUNICATOR_SENTRY_DSN` | Sentry DSN (optional) |

See `.env.example` for the complete list.

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
| `type` | string | Yes | `"email"` or `"sms"` (case-insensitive) |
| `receiver` | string | Yes | Email address or E.164 phone number |
| `subject` | string | Yes | Email subject / SMS label |
| `content` | string | Yes | HTML email body or SMS text |
| `priority` | int | No | Message priority (not currently used) |
| `scheduledAt` | string | No | Future scheduling (not currently used) |

## Project Structure

```
cmd/rabbit-consumer/main.go    Entry point
config/                        Configuration (Viper + YAML)
pkg/consumer/                  RabbitMQ consumer + worker pool
pkg/messaging/                 Message routing + delivery orchestration
pkg/platform/                  Infrastructure adapters (HTTP, logging, mail, DB)
pkg/repository/                MongoDB persistence
```

See `AGENTS.md` for detailed file-by-file documentation.

## Monitoring

- **Sentry:** Error reporting in production (configure `sentry.dsn`)
- **Logs:** Structured JSON via zap (stdout + optional file)
- **MongoDB:** Audit trail of all delivery attempts in `messages` collection

## Development

```bash
# Run with hot reload (requires air or similar)
air

# Run tests
go test ./...

# Lint
go vet ./...
```
