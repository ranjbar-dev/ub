# Communicator — ub-communicator-main

## Overview

Go-based notification microservice for the UnitedBit cryptocurrency exchange platform.
Consumes messages from a RabbitMQ topic exchange and delivers them via email
(SendGrid, Mailjet, Mailgun, or SMTP) or SMS (Twilio). Every delivery attempt —
success or failure — is persisted to MongoDB as an audit log.

The service is a **long-running consumer** with no HTTP endpoints. It blocks on
the RabbitMQ delivery channel until the context is cancelled or the channel closes.

---

## Stack

All versions verified from `go.mod` and Docker files.

| Component | Package | Version | Status |
|-----------|---------|---------|--------|
| Go | — | 1.24.0 (go.mod) / toolchain 1.24.2 / Docker `golang:1.24` | ✅ Current |
| RabbitMQ client | `github.com/rabbitmq/amqp091-go` | v1.10.0 | ✅ Current |
| MongoDB driver | `go.mongodb.org/mongo-driver` | v1.17.2 | ✅ Current |
| Structured logging | `go.uber.org/zap` | v1.27.0 | ✅ Current |
| Configuration | `github.com/spf13/viper` | v1.19.0 | ✅ Current |
| Error reporting | `github.com/getsentry/sentry-go` | v0.31.1 | ✅ Current |
| SendGrid | `github.com/sendgrid/sendgrid-go` | v3.16.0+incompatible | ✅ Current |
| Mailgun | `github.com/mailgun/mailgun-go` | v2.0.0+incompatible | ⚠️ Archived — needs v4 migration |
| Mailjet | `github.com/mailjet/mailjet-apiv3-go/v3` | v3.2.0 | ✅ Current (v4 available) |
| SMTP | `gopkg.in/gomail.v2` | v2.0.0-20160411212932 | ⚠️ Abandoned 2016 — functional but unmaintained |

**Indirect dependencies of note:** `go.uber.org/multierr` v1.10.0, `github.com/sendgrid/rest` v2.6.9.

---

## Architecture

### File Map

```
cmd/rabbit-consumer/main.go       → Entry point (creates DI container, calls Consume)
config/config.go                   → Viper setup (reads config.yaml, env prefix UBCOMMUNICATOR_)
config/config.yaml                 → Default configuration values
pkg/di/container.go                → Custom DI container (lazy singleton initialization)
pkg/consumer/service.go            → RabbitMQ consumer (exchange/queue/bind/consume loop)
pkg/consumer/pool.go               → Worker pool dispatcher (channel-of-channels pattern)
pkg/consumer/worker.go             → Worker goroutine (processes Work items via Send())
pkg/messaging/service.go           → Message orchestration (validation, routing EMAIL/SMS, audit)
pkg/messaging/mailService.go       → Email delivery delegation (wraps MailerClient)
pkg/messaging/smsService.go        → Twilio SMS delivery via HTTP API
pkg/messaging/repository.go        → Message struct + Repository interface definition
pkg/platform/config.go             → Viper config wrapper with typed accessors
pkg/platform/http.go               → HTTP client (GET/POST/form/basic-auth, 10s timeout)
pkg/platform/logger.go             → Zap structured logging + Sentry error reporting
pkg/platform/mail.go               → Mail provider factory (4 providers, switch on mailer_broker)
pkg/platform/mailgunmail.go        → Mailgun provider implementation
pkg/platform/mailjetmail.go        → Mailjet provider implementation
pkg/platform/sendgridmail.go       → SendGrid provider implementation
pkg/platform/smtpmail.go           → SMTP/gomail provider implementation
pkg/platform/mongo.go              → MongoDB client init (connect + ping, 10s timeout)
pkg/platform/rabbitmq.go           → RabbitMQ connection (mutex-protected, lazy connect/reconnect)
pkg/repository/messageRepository.go → MongoDB message persistence (InsertOne, 5s timeout)
```

### Message Flow

```
┌──────────────────────────────────────────────────────────────────┐
│  RabbitMQ                                                        │
│  Exchange: "messages" (topic, durable)                           │
│  Queue:    "messages.command.send.consumer" (durable)            │
│  Binding:  "messages.command.send"                               │
└───────────────────────┬──────────────────────────────────────────┘
                        │ amqp.Delivery
                        ▼
┌──────────────────────────────────────────────────────────────────┐
│  consumer.service.Consume()                                      │
│  - Declares exchange (topic)                                     │
│  - Declares queue (durable)                                      │
│  - Binds queue to exchange                                       │
│  - Starts consuming (autoAck=true)                               │
│  - Calls messaging.Service.CreateMessage(body) to parse JSON     │
│  - Pushes Work{Message} onto Collector.Work channel              │
└───────────────────────┬──────────────────────────────────────────┘
                        │ chan Work (buffered, cap 100)
                        ▼
┌──────────────────────────────────────────────────────────────────┐
│  consumer.pool — Dispatcher goroutine                            │
│  - Reads from Collector.Work channel                             │
│  - Picks first available worker via workerChannel (chan chan Work)│
│  - Sends work to that worker's private channel                   │
└───────────────────────┬──────────────────────────────────────────┘
                        │ per-worker chan Work
                        ▼
┌──────────────────────────────────────────────────────────────────┐
│  consumer.Worker.Start() goroutine                               │
│  - Registers on workerChannel when idle                          │
│  - Receives Work, calls messaging.Service.Send(message)          │
└───────────────────────┬──────────────────────────────────────────┘
                        │
                        ▼
┌──────────────────────────────────────────────────────────────────┐
│  messaging.service.Send(message)                                 │
│  1. validateMessage() — checks receiver, content, type           │
│  2. Routes by message.Type:                                      │
│     EMAIL → messaging.MailService.Send(subject, receiver, content)│
│     SMS   → messaging.SmsService.Send(subject, receiver, content)│
│  3. Sets message.Status = "successful" or "failed"               │
│  4. Persists to MongoDB via messaging.Repository.NewMessage()    │
└───────────────────┬───────────────┬──────────────────────────────┘
                    │               │
        ┌───────────┘               └──────────────┐
        ▼                                          ▼
┌─────────────────────────┐  ┌─────────────────────────────────────┐
│  MailService → factory  │  │  SmsService → Twilio HTTP API       │
│  ┌ sendgrid            │  │  POST /2010-04-01/Accounts/{sid}/   │
│  ├ mailjet             │  │       Messages.json                  │
│  ├ mailgun             │  │  Basic Auth: accountSid:authToken    │
│  └ smtp                │  │  Form body: To, From, Body           │
└────────────┬────────────┘  └─────────────────┬───────────────────┘
             │                                 │
             └────────────┬────────────────────┘
                          ▼
            ┌───────────────────────────┐
            │  MongoDB — ubMessages DB  │
            │  Collection: "messages"   │
            │  Document: Message struct │
            └───────────────────────────┘
```

---

## Package Documentation

### `cmd/rabbit-consumer` — Entry Point

**File:** `cmd/rabbit-consumer/main.go`

Creates the DI container, retrieves the consumer service, and calls `Consume(ctx)`.
On any initialization or runtime error, the process exits via `log.Fatalf`.

```go
func main() {
    container, err := di.NewContainer()  // Eagerly validates config + logger
    consumer := container.GetConsumer()   // Lazy-inits all remaining services
    consumer.Consume(context.Background())
}
```

### `config` — Configuration Loading

**File:** `config/config.go`

| Export | Signature | Description |
|--------|-----------|-------------|
| `FileName` | `const "config"` | Viper config file name (without extension) |
| `SetConfigs` | `func() (*viper.Viper, error)` | Initializes Viper: reads `config.yaml`, enables env override |

**Env var mapping:** prefix `UBCOMMUNICATOR_`, dots replaced with underscores.
Example: `rabbitmq.dsn` → `UBCOMMUNICATOR_RABBITMQ_DSN`.

**Config path:** defaults to `./config`, overridable via `CONFIG_PATH` env var.

**Precedence:** Environment variables > `config.yaml` defaults.

### `pkg/di` — Dependency Injection Container

**File:** `pkg/di/container.go`

Lazy singleton pattern — each service is created on first access and cached.

| Export | Type | Description |
|--------|------|-------------|
| `Container` | interface | Public interface with `GetConsumer() consumer.Service` |
| `NewContainer` | `func() (Container, error)` | Constructor; eagerly validates config + logger |

**Initialization order** (lazy, triggered by `GetConsumer()`):

```
NewContainer()
  ├── getConfigs() → config.SetConfigs() → platform.NewConfigs(viper)
  ├── getLogger() → platform.NewLogger(configs)          [validates sentry DSN]
  └── GetConsumer()
        ├── getRabbitMq() → platform.NewRabbitMqClient(configs, logger)
        ├── getMessagingService()
        │     ├── getMessageRepository()
        │     │     └── getDb() → platform.NewDbClient(configs) [connects + pings MongoDB]
        │     ├── getMailService()
        │     │     └── getMailerClient() → platform.NewMailerClient(configs, logger) [factory]
        │     ├── getSmsService()
        │     │     └── getHttpClient() → platform.NewHttpClient()
        │     └── messaging.NewMessagingService(repo, mail, sms, logger)
        ├── getPool() → consumer.NewPool(messagingService)
        └── consumer.NewConsumerService(rabbitMq, messaging, pool, logger, configs)
```

**Registered services and their dependencies:**

| Service | Constructor | Dependencies |
|---------|-------------|--------------|
| `configs` | `platform.NewConfigs(viper)` | `config.SetConfigs()` output |
| `logger` | `platform.NewLogger(configs)` | `configs` |
| `httpClient` | `platform.NewHttpClient()` | none |
| `db` | `platform.NewDbClient(configs)` | `configs` |
| `rabbitMq` | `platform.NewRabbitMqClient(configs, logger)` | `configs`, `logger` |
| `mailerClient` | `platform.NewMailerClient(configs, logger)` | `configs`, `logger` |
| `messageRepository` | `repository.NewMessageRepository(db, configs)` | `db`, `configs` |
| `mailService` | `messaging.NewMailService(mailerClient)` | `mailerClient` |
| `smsService` | `messaging.NewSmsService(httpClient, configs)` | `httpClient`, `configs` |
| `messagingService` | `messaging.NewMessagingService(repo, mail, sms, logger)` | `messageRepository`, `mailService`, `smsService`, `logger` |
| `pool` | `consumer.NewPool(messagingService)` | `messagingService` |
| `consumer` | `consumer.NewConsumerService(...)` | `rabbitMq`, `messagingService`, `pool`, `logger`, `configs` |

### `pkg/consumer` — RabbitMQ Consumer & Worker Pool

#### `service.go` — Consumer Service

| Export | Type | Description |
|--------|------|-------------|
| `Service` | interface | `Consume(ctx context.Context) error` — blocks until ctx cancelled |
| `NewConsumerService` | func | Constructor: `(rc, ms, pool, logger, configs) → Service` |

**Consume() behavior:**
1. Gets AMQP channel from `RabbitMqClient.GetChannel()`
2. Declares topic exchange (name from `rabbitmq.exchange`, default `"messages"`)
3. Declares durable queue (name from `rabbitmq.queue_name`, default `"messages.command.send.consumer"`)
4. Binds queue with routing key (from `rabbitmq.binding`, default `"messages.command.send"`)
5. Starts consuming with `autoAck=true`
6. Reads `consumer.worker_count` from config (default 5), starts dispatcher
7. Loops: reads deliveries → `CreateMessage(body)` → pushes `Work{Message}` to collector
8. On `ctx.Done()`: signals `collector.End`, returns `ctx.Err()`
9. On channel close: signals `collector.End`, returns error

#### `pool.go` — Worker Pool Dispatcher

| Export | Type | Description |
|--------|------|-------------|
| `Collector` | struct | `Work chan Work` (buffered 100) + `End chan bool` |
| `Pool` | interface | `StartDispatcher(workerCount int) Collector` |
| `NewPool` | func | Constructor: `(ms messaging.Service) → Pool` |

**Channel-of-channels pattern:**
- `pool.workerChannel` is `chan chan Work` — a meta-channel
- Each worker sends its private `chan Work` onto `workerChannel` when idle
- Dispatcher reads from `workerChannel` to find an available worker, then sends work to it
- This ensures work is only dispatched to workers ready to process

**Dispatcher goroutine** (started by `StartDispatcher`):
```
for {
    select {
    case <-end:           → stop all workers, return
    case work := <-input: → workerChan := <-workerChannel; workerChan <- work
    }
}
```

**Back-pressure:** The `Collector.Work` channel is buffered (capacity 100). If all workers
are busy AND the buffer is full, the consumer's select loop blocks on `collector.Work <- Work{...}`,
which back-pressures the RabbitMQ delivery channel. The dispatcher also blocks on
`<-p.workerChannel` when no workers are available, creating a natural flow-control cascade.

#### `worker.go` — Worker Goroutine

| Export | Type | Description |
|--------|------|-------------|
| `Work` | struct | `ID int64`, `Message messaging.Message` |
| `Worker` | struct | `ID int`, `WorkerChannel chan chan Work`, `Channel chan Work`, `End chan bool`, `Ms messaging.Service` |
| `Worker.Start()` | method | Launches goroutine: registers as available, processes work |
| `Worker.Stop()` | method | Sends `true` to `End` channel, blocking until worker receives |

**Worker lifecycle:**
1. `Start()` launches a goroutine
2. Worker sends its `Channel` onto `WorkerChannel` (registers as available)
3. Waits on `select { case work := <-Channel: ... case <-End: return }`
4. On work: calls `w.Ms.Send(work.Message)`; errors are logged via `log.Printf`
5. On End: goroutine exits
6. Loop repeats (re-registers as available after each work item)

### `pkg/messaging` — Message Routing & Delivery

#### `service.go` — Message Orchestrator

| Export | Type | Description |
|--------|------|-------------|
| `Service` | interface | `Send(message Message) error`, `CreateMessage(data []byte) (Message, error)` |
| `NewMessagingService` | func | Constructor: `(repo, mail, sms, logger) → Service` |

**`CreateMessage(data []byte)`:**
1. `json.Unmarshal` into `Message` struct
2. Normalizes `Type` to uppercase (`strings.ToUpper`)
3. Sets `Status = "pending"`, `CreatedAt = time.Now()`

**`validateMessage(msg Message)`:**
- Checks `Receiver` is non-empty
- Checks `Content` is non-empty
- For EMAIL: validates with `net/mail.ParseAddress()`
- For SMS: validates E.164 format via regex `^\+[1-9]\d{6,14}$`
- Unknown type returns error

**`Send(message Message)`:**
1. Validates message; on failure: sets status "failed", persists, returns error
2. Routes by `message.Type`:
   - `EMAIL` → `MailService.Send(subject, receiver, content)`
   - `SMS` → `SmsService.Send(subject, receiver, content)`
3. Sets `Status` to `"successful"` or `"failed"` based on `(isSent, err)` return
4. Persists message to MongoDB via `Repository.NewMessage()`
5. Returns `nil` (errors are logged but not propagated to caller)

#### `mailService.go` — Email Delivery Wrapper

| Export | Type | Description |
|--------|------|-------------|
| `MailService` | interface | `Send(subject, receiver, content string) (bool, error)` |
| `NewMailService` | func | Constructor: `(mc platform.MailerClient) → MailService` |

Direct passthrough to `platform.MailerClient.Send()`. Exists as an extension point for
future cross-cutting concerns (retry, circuit breaker, metrics) at the messaging layer.

#### `smsService.go` — Twilio SMS

| Export | Type | Description |
|--------|------|-------------|
| `SmsService` | interface | `Send(subject, receiver, content string) (bool, error)` |
| `NewSmsService` | func | Constructor: `(httpClient, configs) → SmsService` |

**Twilio HTTP API details:**
- **Endpoint:** `https://api.twilio.com/2010-04-01/Accounts/{accountSid}/Messages.json`
- **Method:** POST
- **Auth:** HTTP Basic Auth (`accountSid:authToken`, base64 encoded)
- **Content-Type:** `application/x-www-form-urlencoded`
- **Accept:** `application/json`
- **Form body:** `To={receiver}&From={from}&Body={content}`
- **Success:** HTTP 2xx → parses JSON response → returns `(true, nil)`
- **Failure:** non-2xx → returns `(false, error)` with status code

Note: The `subject` parameter is accepted but not used by the Twilio API call.

#### `repository.go` — Message Model & Repository Interface

| Export | Type | Description |
|--------|------|-------------|
| `MessageStatusPending` | const `"pending"` | Initial status set by `CreateMessage` |
| `MessageStatusFailed` | const `"failed"` | Set on delivery failure or validation error |
| `MessageStatusSuccessful` | const `"successful"` | Set on successful delivery |
| `MessageTypeEmail` | const `"EMAIL"` | Uppercase email type |
| `MessageTypeSms` | const `"SMS"` | Uppercase SMS type |
| `Message` | struct | Notification data model (see schema below) |
| `Repository` | interface | `NewMessage(message *Message) error` |

### `pkg/platform` — Infrastructure Adapters

#### `config.go` — Configuration Wrapper

| Export | Type | Description |
|--------|------|-------------|
| `AllowedIpsConfigKey` | const `"wallet.allowed_ips"` | Config key for IP allowlist |
| `SentryDsnKey` | const `"sentry.dsn"` | Config key for Sentry DSN |
| `EnvConfigKey` | const `"wallet.environment"` | Config key for environment |
| `EnvProd`, `EnvTest`, `EnvDev` | const | Environment name constants |
| `SmsUrlPrefix` | const `"https://api.twilio.com/2010-04-01/Accounts/"` | Twilio base URL |
| `SmsUrlPostfix` | const `"/Messages.json"` | Twilio endpoint suffix |
| `Configs` | interface | Typed config accessors (see below) |
| `NewConfigs` | func | Constructor: `(viper *viper.Viper) → Configs` |

**`Configs` interface methods:**
| Method | Signature | Description |
|--------|-----------|-------------|
| `GetString` | `(name string) string` | String config value |
| `GetInt` | `(name string) int` | Integer config value |
| `GetBool` | `(name string) bool` | Boolean config value |
| `GetStringSlice` | `(name string) []string` | String slice config value |
| `UnmarshalKey` | `(key string, i interface{}) error` | Unmarshal a config section |
| `GetAllowedIps` | `() []string` | Shortcut for `wallet.allowed_ips` |
| `GetEnv` | `() string` | Shortcut for `wallet.environment` |
| `GetSentryDsn` | `() string` | Shortcut for `sentry.dsn` |
| `GetSmsUrl` | `(sId string) string` | Builds Twilio URL from account SID |

Note: `NewConfigs` auto-detects `go test` via `flag.Lookup("test.v")` and sets env to `"test"`.

#### `http.go` — HTTP Client

| Export | Type | Description |
|--------|------|-------------|
| `HttpClient` | interface | HTTP methods with Basic Auth |
| `NewHttpClient` | func | Constructor: returns client with 10s timeout |

**`HttpClient` interface methods:**
| Method | Signature | Description |
|--------|-----------|-------------|
| `HttpGet` | `(url string) ([]byte, http.Header, int, error)` | GET request |
| `HttpPost` | `(url string, body interface{}, headers map[string]string) ([]byte, http.Header, int, error)` | POST with JSON body |
| `HttpPostForm` | `(url string, body *strings.Reader, headers map[string]string) ([]byte, http.Header, int, error)` | POST with form body |
| `BasicAuth` | `(username, password string) string` | Base64-encodes `user:pass` |

All methods return `(responseBody, responseHeaders, statusCode, error)`.

#### `logger.go` — Zap + Sentry Logging

| Export | Type | Description |
|--------|------|-------------|
| `Logger` | interface | `Info`, `Warn`, `Error`, `Fatal` (each with `...zap.Field`) |
| `NewLogger` | func | Constructor: `(configs Configs) → Logger` |

**Logger interface methods:**
| Method | Description |
|--------|-------------|
| `Info(msg string, fields ...zap.Field)` | Informational log |
| `Warn(msg string, fields ...zap.Field)` | Warning log |
| `Error(msg string, fields ...zap.Field)` | Error log + sends to Sentry (prod only) |
| `Fatal(msg string, fields ...zap.Field)` | Fatal log (calls `os.Exit(1)`) |

**Sentry integration:**
- Initialized once in `NewLogger()` if `sentry.dsn` is configured
- `Error()` extracts `error` fields and calls `sentry.CaptureException()`
- Only reports in production (`env == "prod"`); dev/test are excluded
- Flushes with a 2-second timeout after each capture

**Log output:**
- Default: stdout (zap production JSON format)
- Optional: dual output to stdout + file path (from `logging.file_path` config)
- Creates log directory if it doesn't exist; falls back to stdout on failure

#### `mail.go` — Mail Provider Factory

| Export | Type | Description |
|--------|------|-------------|
| `MailerSendGrid` | const `"sendgrid"` | Provider name constant |
| `MailerMailJet` | const `"mailjet"` | Provider name constant |
| `MailerMailGun` | const `"mailgun"` | Provider name constant |
| `MailerSMTP` | const `"smtp"` | Provider name constant |
| `MailerClient` | interface | `Send(subject, receiver, content string) (bool, error)` |
| `NewMailerClient` | func | Factory: `(configs, logger) → MailerClient` (returns `nil` for unknown provider) |

**Factory behavior** (switches on `configs.GetString("mailer_broker")`):

| Provider | Config Keys | SDK / Library |
|----------|-------------|---------------|
| `sendgrid` | `sendgrid.api_key` | `sendgrid.NewSendClient(apiKey)` |
| `mailjet` | `mailjet.api_public_key`, `mailjet.api_private_key` | `mailjet.NewMailjetClient(pub, priv)` |
| `mailgun` | `mailgun.api_key`, `mailgun.domain`, `mailgun.api_base` | `mailgun.NewMailgun(apiKey, domain)` + `SetAPIBase` |
| `smtp` | `smtp.host`, `smtp.port`, `smtp.username`, `smtp.password` | `gomail.NewDialer(host, port, user, pass)` |

All providers also read: `mail.name` (sender display name), `mail.from_address` (sender email).

#### `sendgridmail.go` — SendGrid Provider

**Struct:** `sendGridMailerClient` (fields: `*sendgrid.Client`, `name`, `fromAddress`, `logger`)

**Send() behavior:**
1. Creates `mail.NewEmail` for From and To
2. Prefixes subject with `[UNITEDBIT]` unless it already contains `[`
3. Calls `mail.NewSingleEmail(from, subject, to, plainText, html)`
4. Calls `client.Send(message)`
5. Checks response status code is 200-300
6. Returns `(true, nil)` on success, `(false, error)` on failure

**API:** Uses SendGrid Go SDK which calls `https://api.sendgrid.com/v3/mail/send` internally.
**Auth:** API key passed to `sendgrid.NewSendClient()`.

#### `mailjetmail.go` — Mailjet Provider

**Struct:** `mailJetMailerClient` (fields: `*mailjet.Client`, `name`, `fromAddress`, `logger`)

**Send() behavior:**
1. Prefixes subject with `[UNITEDBIT]` unless it already contains `[`
2. Builds `mailjet.MessagesV31` with From, To, Subject, TextPart, HTMLPart
3. Calls `client.SendMailV31(&messages)`
4. Returns `(true, nil)` on success, `(false, error)` on failure

**API:** Uses Mailjet Go SDK which calls `https://api.mailjet.com/v3.1/send` internally.
**Auth:** API public key + private key passed to `mailjet.NewMailjetClient()`.

#### `mailgunmail.go` — Mailgun Provider

**Struct:** `mailGunClient` (fields: `*mailgun.MailgunImpl`, `name`, `fromAddress`, `logger`)

**Send() behavior:**
1. Creates message with `client.NewMessage(from, subject, plainText, receiver)`
2. Sets HTML content with `message.SetHtml(content)`
3. Calls `client.Send(message)`
4. Returns `(true, nil)` on success, `(false, error)` on failure

**API:** Uses Mailgun Go SDK; API base configurable (default: `https://api.eu.mailgun.net/v3`).
**Auth:** API key + domain passed to `mailgun.NewMailgun()`.

Note: Does NOT prefix subject with `[UNITEDBIT]` — inconsistent with SendGrid and Mailjet.

#### `smtpmail.go` — SMTP Provider

**Struct:** `smtpClient` (fields: `*gomail.Dialer`, `name`, `fromAddress`, `logger`)

**Send() behavior:**
1. Creates `gomail.NewMessage()` with From, To, Subject headers
2. Sets body as `text/html`
3. Calls `client.DialAndSend(message)` — connects, sends, disconnects per email
4. Returns `(true, nil)` on success, `(false, error)` on failure

**TLS config:** `InsecureSkipVerify: false`, `MinVersion: tls.VersionTLS12`.

Note: Does NOT prefix subject with `[UNITEDBIT]` — inconsistent with SendGrid and Mailjet.
Note: `DialAndSend` creates a new TCP connection per email (no connection pooling).

#### `mongo.go` — MongoDB Client

| Export | Type | Description |
|--------|------|-------------|
| `NewDbClient` | func | `(configs Configs) (*mongo.Client, error)` — connect + ping with 10s timeout |

Reads `mongodb.dsn` from config. Returns error if connection or ping fails.

#### `rabbitmq.go` — RabbitMQ Connection

| Export | Type | Description |
|--------|------|-------------|
| `RabbitMqClient` | interface | `GetChannel() (*amqp.Channel, error)` |
| `NewRabbitMqClient` | func | `(configs, logger) → RabbitMqClient` — deferred connection |

**Connection management:**
- Lazy: actual AMQP dial happens on first `GetChannel()` call
- Mutex-protected: `sync.Mutex` prevents concurrent connection creation
- Reconnect-aware: checks `connection.IsClosed()` and re-dials if needed
- On channel open failure: marks `isConnected = false` for next-call reconnect
- Reads `rabbitmq.dsn` from config

### `pkg/repository` — MongoDB Persistence

**File:** `pkg/repository/messageRepository.go`

| Export | Type | Description |
|--------|------|-------------|
| `CollectionName` | const `"messages"` | MongoDB collection name |
| `NewMessageRepository` | func | `(db *mongo.Client, configs Configs) → messaging.Repository` |

**`NewMessage(message *Message)` behavior:**
- Uses `context.WithTimeout(5 * time.Second)`
- Calls `collection.InsertOne(ctx, message)`
- Database name from `mongodb.name` config (default: `"ubMessages"`)

---

## Message Format

### RabbitMQ Input (JSON)

```json
{
  "type": "email",
  "receiver": "user@example.com",
  "subject": "Welcome to UnitedBit",
  "content": "<h1>Hello</h1><p>Your account is ready.</p>",
  "priority": 1,
  "scheduledAt": ""
}
```

### Field Schema

| Field | JSON Key | BSON Key | Go Type | Required | Validation |
|-------|----------|----------|---------|----------|------------|
| Type | `type` | `type` | `string` | Yes | Normalized to uppercase; must be `"EMAIL"` or `"SMS"` |
| Receiver | `receiver` | `receiver` | `string` | Yes | EMAIL: valid email (`net/mail.ParseAddress`); SMS: E.164 (`^\+[1-9]\d{6,14}$`) |
| Subject | `subject` | `subject` | `string` | No | Used as email subject; ignored by Twilio SMS |
| Content | `content` | `content` | `string` | Yes | HTML for email; plain text for SMS |
| Priority | `priority` | `priority` | `int` | No | Stored but not currently used for ordering |
| ScheduledAt | `scheduledAt` | `scheduledAt` | `string` | No | Stored but not currently used for scheduling |
| Status | `status` | `status` | `string` | Auto | Set by system: `"pending"` → `"successful"` or `"failed"` |
| CreatedAt | `createdAt` | `createdAt` | `time.Time` | Auto | Set to `time.Now()` on message creation |

### Type Normalization

The `CreateMessage()` function calls `strings.ToUpper(message.Type)`, so producers can
send `"email"`, `"Email"`, `"EMAIL"`, etc. — all are normalized to `"EMAIL"`.

---

## MongoDB Schema

**Database:** `ubMessages` (from `mongodb.name` config)
**Collection:** `messages`

**Document structure** (maps from `messaging.Message` struct via BSON tags):

```json
{
  "receiver": "user@example.com",
  "subject": "Welcome",
  "content": "<h1>Hello</h1>",
  "priority": 1,
  "scheduledAt": "",
  "type": "EMAIL",
  "status": "successful",
  "createdAt": ISODate("2025-01-15T10:30:00Z")
}
```

**Status values:** `"pending"`, `"successful"`, `"failed"`
**Type values:** `"EMAIL"`, `"SMS"`

**Note:** No indexes are defined by the application. For production queries, consider adding:
```js
db.messages.createIndex({ "createdAt": -1 })
db.messages.createIndex({ "status": 1 })
db.messages.createIndex({ "receiver": 1, "createdAt": -1 })
```

**MongoDB init script** (`.docker/mongo/mongo-init.js`):
Creates user `ub_mongo_user` with `readWrite` role on `ubMessages` database.

---

## Configuration Reference

### `config/config.yaml` — Complete Key Reference

| Section | Key | Type | Default | Env Override | Description |
|---------|-----|------|---------|-------------|-------------|
| `communicator` | `environment` | string | `"dev"` | `UBCOMMUNICATOR_COMMUNICATOR_ENVIRONMENT` | Runtime environment (`dev`/`test`/`prod`) |
| `communicator` | `allowed_ips` | []string | `["127.0.0.1"]` | `UBCOMMUNICATOR_COMMUNICATOR_ALLOWED_IPS` | IP allowlist (not currently used by consumer) |
| `consumer` | `worker_count` | int | `5` | `UBCOMMUNICATOR_CONSUMER_WORKER_COUNT` | Worker pool size |
| `mail` | `name` | string | `"UNITEDBIT"` | `UBCOMMUNICATOR_MAIL_NAME` | Sender display name |
| `mail` | `from_address` | string | `"no-reply@unitedbit.com"` | `UBCOMMUNICATOR_MAIL_FROM_ADDRESS` | Sender email address |
| `sms` | `account_sid` | string | `""` | `UBCOMMUNICATOR_SMS_ACCOUNT_SID` | Twilio Account SID |
| `sms` | `auth_token` | string | `""` | `UBCOMMUNICATOR_SMS_AUTH_TOKEN` | Twilio Auth Token |
| `sms` | `from` | string | `""` | `UBCOMMUNICATOR_SMS_FROM` | Twilio sender phone number |
| `mongodb` | `dsn` | string | `""` | `UBCOMMUNICATOR_MONGODB_DSN` | MongoDB connection URI |
| `mongodb` | `name` | string | `"ubMessages"` | `UBCOMMUNICATOR_MONGODB_NAME` | MongoDB database name |
| `rabbitmq` | `dsn` | string | `""` | `UBCOMMUNICATOR_RABBITMQ_DSN` | RabbitMQ AMQP URI |
| `rabbitmq` | `queue_name` | string | `"messages.command.send.consumer"` | `UBCOMMUNICATOR_RABBITMQ_QUEUE_NAME` | Queue name |
| `rabbitmq` | `exchange` | string | — | `UBCOMMUNICATOR_RABBITMQ_EXCHANGE` | Exchange name (code default: `"messages"`) |
| `rabbitmq` | `binding` | string | — | `UBCOMMUNICATOR_RABBITMQ_BINDING` | Routing key (code default: `"messages.command.send"`) |
| — | `mailer_broker` | string | `"smtp"` | `UBCOMMUNICATOR_MAILER_BROKER` | Email provider: `smtp`/`sendgrid`/`mailjet`/`mailgun` |
| `sendgrid` | `api_key` | string | `""` | `UBCOMMUNICATOR_SENDGRID_API_KEY` | SendGrid API key |
| `mailjet` | `api_public_key` | string | `""` | `UBCOMMUNICATOR_MAILJET_API_PUBLIC_KEY` | Mailjet public API key |
| `mailjet` | `api_private_key` | string | `""` | `UBCOMMUNICATOR_MAILJET_API_PRIVATE_KEY` | Mailjet private API key |
| `mailgun` | `api_key` | string | `""` | `UBCOMMUNICATOR_MAILGUN_API_KEY` | Mailgun API key |
| `mailgun` | `domain` | string | `""` | `UBCOMMUNICATOR_MAILGUN_DOMAIN` | Mailgun sending domain |
| `mailgun` | `api_base` | string | `"https://api.eu.mailgun.net/v3"` | `UBCOMMUNICATOR_MAILGUN_API_BASE` | Mailgun API base URL |
| `smtp` | `host` | string | `""` | `UBCOMMUNICATOR_SMTP_HOST` | SMTP server hostname |
| `smtp` | `port` | int | `587` | `UBCOMMUNICATOR_SMTP_PORT` | SMTP server port |
| `smtp` | `username` | string | `""` | `UBCOMMUNICATOR_SMTP_USERNAME` | SMTP username |
| `smtp` | `password` | string | `""` | `UBCOMMUNICATOR_SMTP_PASSWORD` | SMTP password |
| `sentry` | `dsn` | string | `""` | `UBCOMMUNICATOR_SENTRY_DSN` | Sentry DSN (empty = disabled) |
| `sentry` | `debug` | bool | `false` | `UBCOMMUNICATOR_SENTRY_DEBUG` | Sentry debug mode |
| `logging` | `file_path` | string | `"stdout"` | `UBCOMMUNICATOR_LOGGING_FILE_PATH` | Log file path (`"stdout"` for console only) |

---

## Docker & Deployment

### Docker Compose Files

| File | Purpose | Runner Tag |
|------|---------|------------|
| `docker-compose.yml` | Local development | — |
| `docker-compose-dev.yml` | Dev server deployment | `communicator-app-dev` |
| `docker-compose-prod.yml` | Production deployment | `communicator-app-prod` |

### Services

**All compose files define two services:**

| Service | Container | Image | Networks |
|---------|-----------|-------|----------|
| `mongodb` | `communicator-db` | `mongo` (official) | `default` |
| `go` | `communicator-app` | Build from Dockerfile | `default`, `ub-server_rabbit` (external) |

**External network:** `ub-server_rabbit` connects the communicator to the RabbitMQ instance
managed by the `ub-server` service's docker-compose network.

### Dockerfiles

| File | Use | Description |
|------|-----|-------------|
| `.docker/go/Dockerfile.dev` | Development | `golang:1.24`, working dir `/app`, no build step (code is mounted) |
| `.docker/go/Dockerfile.prod` | Production | `golang:1.24`, copies source, builds `go build -mod=vendor cmd/rabbit-consumer/main.go` |

### Dev vs Prod Differences

| Aspect | Dev (`docker-compose.yml`) | Prod (`docker-compose-prod.yml`) |
|--------|---------------------------|----------------------------------|
| Code | Volume-mounted `.:/app` | Copied into image at build |
| Build | Manual: `docker exec communicator-app go build ...` | Built in Dockerfile |
| Command | `tail -F docker-compose.yml` (keeps alive) | `bash -c "./main"` |
| MongoDB port | `27017:27017` (all interfaces) | `127.0.0.1:27017:27017` (localhost only) |
| MongoDB data | `./../communicator-db` | `/home/exchange/communicator-app/var` |
| Env vars | Docker-compose environment | `.env` file at `/home/gitlab-runner/communicator-app/.env` |

### CI/CD — `.gitlab-ci.yml`

**Stages:** `build-dev` → `deploy-dev` → `dev-notification` → `build-prod` → `deploy-prod` → `prod-notification`

| Branch | Action |
|--------|--------|
| `develop` | Build + deploy via `docker-compose-dev.yml`, then `docker exec` to compile |
| `master` | Build + deploy via `docker-compose-prod.yml` (binary pre-built in image) |

Both branches send Telegram notifications on success (🟢) or failure (🔴).

---

## Worker Pool — Technical Detail

### Pattern: Channel-of-Channels Dispatcher

```
                            workerChannel (chan chan Work)
                           ┌───────────────────────────────┐
                           │                               │
   ┌─────────┐   Work     │   ┌────────────┐             │   ┌────────────┐
   │Consumer │──────────▶│   │ Dispatcher │─────────────▶│   │  Worker 1 │
   │  Loop   │  buffered  │   │ goroutine  │  pick idle   │   │ goroutine  │
   └─────────┘  chan(100)  │   └────────────┘  worker      │   └────────────┘
                           │                               │   ┌────────────┐
                           │   Workers register when idle: │   │  Worker 2 │
                           │   w.WorkerChannel <- w.Channel│   │ goroutine  │
                           │                               │   └────────────┘
                           │                               │   ┌────────────┐
                           │                               │   │  Worker N │
                           └───────────────────────────────┘   └────────────┘
```

**Flow:**
1. Consumer pushes `Work` onto `Collector.Work` (buffered channel, capacity 100)
2. Dispatcher reads from `Collector.Work`
3. Dispatcher reads from `pool.workerChannel` — blocks until a worker registers
4. Dispatcher sends work to the worker's private channel
5. Worker processes, then loops back to step 3 (re-registers)

**Shutdown sequence:**
1. Consumer sends `true` to `Collector.End`
2. Dispatcher receives on `End`, calls `Stop()` on each worker
3. `Worker.Stop()` sends `true` to worker's `End` channel
4. Worker goroutine returns

**Error handling:**
- Worker catches errors from `messaging.Service.Send()` and logs via `log.Printf`
- Errors do NOT propagate back to the dispatcher or consumer
- A failed delivery does not crash the worker — it re-registers and processes the next item

---

## Logging & Monitoring

### Zap Structured Logging

- **Format:** JSON (zap production config)
- **Output:** stdout (always) + optional file (from `logging.file_path`)
- **Levels used:** Info, Warn, Error, Fatal
- **Fields:** Uses `zap.Field` for structured context (e.g., `zap.Error(err)`)

### Sentry Error Reporting

- **Trigger:** Every `Logger.Error()` call with a `zap.Error(err)` field
- **Environment gate:** Only active when `wallet.environment == "prod"`
- **Init:** Once in `NewLogger()` using `sentry.dsn` config value
- **Flush:** 2-second timeout after each `CaptureException`

### MongoDB Audit Trail

- Every delivery attempt is persisted to MongoDB regardless of outcome
- Failed validation, failed delivery, and successful delivery are all logged
- Query audit trail: `db.messages.find({ receiver: "user@example.com" }).sort({ createdAt: -1 })`

---

## Known Issues

### 🔴 Critical

| Issue | Location | Impact |
|-------|----------|--------|
| **autoAck=true** on RabbitMQ consume | `pkg/consumer/service.go:76` | Messages are ACK'd immediately on delivery from RabbitMQ, not after processing. Messages are lost if the service crashes mid-processing. Fix: use `autoAck=false` and call `d.Ack(false)` after `Send()` succeeds. |
| **No reconnection loop** for RabbitMQ | `pkg/consumer/service.go:89-107` | If the RabbitMQ connection drops, the delivery channel closes and the consumer exits with an error. The process must be restarted externally. No automatic reconnect/retry logic. |
| **Telegram bot token hardcoded** | `.gitlab-ci.yml:37-38` | Bot token `1416070700:AAGSBy7q...` is committed in plaintext. Should use CI/CD variables. |

### 🟡 Medium

| Issue | Location | Impact |
|-------|----------|--------|
| **No health check endpoint** | — (no HTTP server) | No way to probe service liveness. Docker/K8s cannot health-check. |
| **mailgun-go v2 archived** | `go.mod:9` | Library is archived; should migrate to `github.com/mailgun/mailgun-go/v4`. Breaking API changes required. |
| **gomail.v2 abandoned** (2016) | `go.mod:16` | Functional but receives no security patches. |
| **Subject prefix inconsistency** | `sendgridmail.go:24`, `mailjetmail.go:21` vs `mailgunmail.go`, `smtpmail.go` | SendGrid and Mailjet prefix subjects with `[UNITEDBIT]`; Mailgun and SMTP do not. |
| **SMTP DialAndSend per email** | `pkg/platform/smtpmail.go:23` | Creates new TCP+TLS connection per email. No connection pooling. |
| **No message deduplication** | `pkg/messaging/service.go` | No idempotency key; duplicate messages from RabbitMQ will send duplicate notifications. |

### 🟢 Low / Informational

| Issue | Location | Impact |
|-------|----------|--------|
| **`priority` field unused** | `pkg/messaging/repository.go:24` | Stored but never used for ordering or priority routing. |
| **`scheduledAt` field unused** | `pkg/messaging/repository.go:25` | Stored but no scheduling logic exists. |
| **Worker uses `log.Printf`** | `pkg/consumer/worker.go:32` | Workers log via stdlib `log` instead of the structured `platform.Logger`. |
| **Pool uses `fmt.Printf`** | `pkg/consumer/pool.go:33` | "starting worker: N" printed via `fmt.Printf` instead of structured logger. |
| **`wallet.environment` vs `communicator.environment`** | `config.go:19` vs `config.yaml:2` | Config key mismatch: `platform.EnvConfigKey` reads `"wallet.environment"` but config.yaml has `"communicator.environment"`. The env var override still works. |
| **`wallet.allowed_ips`** | `config.go:18` | Config key references `wallet.allowed_ips` but config.yaml has `communicator.allowed_ips`. |
| **No MongoDB indexes** | `pkg/repository/messageRepository.go` | Collection has no application-managed indexes. |
| **Compiled binary in repo** | `rabbit-consumer`, `rabbit-consumer.exe` | Binary artifacts committed to git. Should be in `.gitignore`. |

---

## Testing Guidelines

### Unit Testing Strategy

All services use interface-based design, making mocking straightforward.

**Mocking dependencies:**

| Package | Interface to Mock | Test Focus |
|---------|-------------------|------------|
| `pkg/messaging` | `messaging.Repository`, `messaging.MailService`, `messaging.SmsService` | Message routing, validation, status updates |
| `pkg/consumer` | `messaging.Service`, `consumer.Pool` | Consumer loop, message parsing |
| `pkg/platform` | `platform.MailerClient`, `platform.HttpClient`, `platform.Configs`, `platform.Logger` | Provider logic, HTTP calls |

**Key test scenarios:**

1. **Message validation** (`pkg/messaging/service.go`):
   - Valid email address → passes
   - Invalid email → fails with error
   - Valid E.164 phone → passes
   - Invalid phone format → fails with error
   - Empty receiver → fails
   - Empty content → fails
   - Unknown type → fails

2. **Message routing** (`pkg/messaging/service.go`):
   - Type="EMAIL" → calls `MailService.Send()`, not `SmsService`
   - Type="SMS" → calls `SmsService.Send()`, not `MailService`
   - Successful send → status="successful", persisted
   - Failed send → status="failed", persisted
   - Validation failure → status="failed", persisted, returns error

3. **Worker pool** (`pkg/consumer/pool.go`, `worker.go`):
   - N workers started → N goroutines running
   - Work dispatched → `messaging.Service.Send()` called
   - End signal → all workers stop
   - Send error → worker continues processing (doesn't crash)

4. **Mail providers** (`pkg/platform/*mail.go`):
   - Each provider: mock the underlying SDK client
   - Test subject prefix behavior (SendGrid/Mailjet add prefix; Mailgun/SMTP don't)
   - Test error propagation from SDK

5. **SMS service** (`pkg/messaging/smsService.go`):
   - Mock `HttpClient.HttpPostForm()`
   - Test success path (2xx response)
   - Test failure path (non-2xx response)
   - Verify request headers (Basic Auth, Content-Type)
   - Verify form body (To, From, Body)

### Integration Testing

- Use `testcontainers-go` for MongoDB and RabbitMQ
- Test end-to-end: publish message to RabbitMQ → verify document in MongoDB
- Test MongoDB connectivity: `NewDbClient` with real container

### Running Tests

```bash
go test ./...                  # Run all tests
go test ./pkg/messaging/...    # Test messaging package
go test -v -run TestValidation # Run specific test
go vet ./...                   # Static analysis
```

---

## Conventions

### Interface-Based Design
Every service is defined as an interface with a private struct implementation.
Constructors follow the `New<Type>(deps...) <Interface>` pattern.

```go
type Service interface { ... }
type service struct { ... }
func NewService(deps...) Service { return &service{...} }
```

### Dependency Injection
- Manual DI via `pkg/di/container.go` (no framework)
- Lazy singleton initialization (create on first access, cache for reuse)
- Constructor returns `(Container, error)` — no panics

### Platform Abstraction
All external service integrations live in `pkg/platform/`:
- Infrastructure concerns are isolated from business logic
- `pkg/messaging/` and `pkg/consumer/` depend only on interfaces

### Factory Pattern
`platform.NewMailerClient()` is a factory that selects the implementation
based on the `mailer_broker` config value. Returns `nil` for unknown values.

### Error Handling
- Constructors return `error` where initialization can fail
- `fmt.Errorf("context: %w", err)` for error wrapping
- `messaging.Service.Send()` logs errors but returns `nil` to prevent worker crashes
- Repository errors are logged but don't stop delivery tracking

### Configuration
- Viper with YAML config file + environment variable overrides
- Env prefix: `UBCOMMUNICATOR_`; dots replaced with underscores
- Secrets should always be set via env vars, never committed to config.yaml

---

## Spec-Driven Development Guide

### Adding a New Email Provider

1. Create `pkg/platform/<name>mail.go`
2. Define a private struct implementing `MailerClient` interface:
   ```go
   type myProviderClient struct {
       client  *myProviderSDK.Client
       name    string
       fromAddress string
       logger  Logger
   }
   func (m *myProviderClient) Send(subject, receiver, content string) (bool, error) { ... }
   ```
3. Add a const in `pkg/platform/mail.go`: `MailerMyProvider = "myprovider"`
4. Add a case in `NewMailerClient()` switch statement in `pkg/platform/mail.go`
5. Add config keys to `config/config.yaml` and `.env.example`
6. No other files need to change — the factory pattern handles wiring

### Adding a New Notification Channel (e.g., Push, Telegram)

1. Create `pkg/messaging/<channel>Service.go` with interface + implementation
2. Add a new `MessageType<Channel>` constant in `pkg/messaging/repository.go`
3. Add a routing case in `messaging.service.Send()` in `pkg/messaging/service.go`
4. Add validation logic in `validateMessage()` for the new type
5. Wire the new service in `pkg/di/container.go`:
   - Add field to `container` struct
   - Add getter method
   - Inject into `getMessagingService()`
6. Update `messaging.NewMessagingService()` signature to accept the new service

### Adding Health Checks

No HTTP server exists. To add one:
1. Create `pkg/health/handler.go` with an HTTP handler
2. Add health check logic: ping MongoDB, check RabbitMQ connection
3. Register in DI container
4. Start HTTP server in `main.go` (separate goroutine from consumer)
5. Add Docker `HEALTHCHECK` instruction in Dockerfiles

### Adding Graceful Shutdown

The consumer already supports context cancellation (`ctx.Done()`). To wire it up:
1. In `main.go`, create a context with `signal.NotifyContext(ctx, syscall.SIGTERM, syscall.SIGINT)`
2. Pass this context to `consumer.Consume(ctx)`
3. Add `defer sentry.Flush(5 * time.Second)` in main
4. Add MongoDB client disconnect: `defer db.Disconnect(ctx)`

### Adding Reconnection Logic

Current behavior: if RabbitMQ connection drops, consumer exits.
To add reconnection:
1. Wrap the `Consume()` body in a retry loop with exponential backoff
2. On delivery channel close, log warning and reconnect
3. Consider using `connection.NotifyClose()` for proactive detection
4. The `RabbitMqClient.connect()` already handles reconnection at the connection level;
   the gap is at the consumer loop level

### Writing Tests for Any Component

All interfaces can be mocked. Recommended mock generation:
```bash
# Install mockgen
go install go.uber.org/mock/mockgen@latest

# Generate mocks
mockgen -source=pkg/messaging/repository.go -destination=pkg/messaging/mocks/repository_mock.go
mockgen -source=pkg/platform/mail.go -destination=pkg/platform/mocks/mail_mock.go
mockgen -source=pkg/platform/http.go -destination=pkg/platform/mocks/http_mock.go
```

Or use hand-written mocks since interfaces are small (1-4 methods each).
