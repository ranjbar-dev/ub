# Exchange CLI — ub-exchange-cli-main

## Stack

| Category | Library | Version (go.mod) | Purpose |
|----------|---------|-------------------|---------|
| Language | Go | 1.22 | Runtime |
| HTTP | gin-gonic/gin | 1.10.1 | REST API framework |
| ORM | gorm.io/gorm | 1.21.15 | MySQL ORM |
| DB Driver | gorm.io/driver/mysql | 1.1.2 | MySQL/MariaDB |
| Redis | go-redis/redis/v8 | 8.11.3 | Cache, order books, queues |
| Redis Cache | go-redis/cache/v8 | 8.4.3 | In-process LRU + Redis |
| MQTT | eclipse/paho.mqtt.golang | 1.3.5 | Real-time publish to clients |
| RabbitMQ | rabbitmq/amqp091-go | 1.10.0 | Email/SMS queue to ub-communicator |
| JWT | golang-jwt/jwt/v5 | 5.2.1 | Authentication tokens |
| gRPC | google.golang.org/grpc | 1.40.0 | Candle service client |
| Protobuf | google.golang.org/protobuf | 1.36.1 | gRPC serialization |
| WebSocket | gorilla/websocket | 1.5.0 | Binance WS streams |
| DI | sarulabs/di | 2.0.0 | Dependency injection container |
| Decimal | shopspring/decimal | 1.2.0 | Precision financial math |
| Config | spf13/viper | 1.9.0 | YAML config + env override |
| Logging | go.uber.org/zap | 1.19.1 | Structured logging |
| Sentry | getsentry/sentry-go | 0.30.0 | Error tracking |
| Validation | go-playground/validator/v10 | 10.23.0 | Request validation |
| CORS | gin-contrib/cors | 1.7.3 | Cross-origin middleware |
| 2FA/TOTP | pquerna/otp | 1.3.0 | Google Authenticator |
| Phone | ttacon/libphonenumber | 1.2.1 | Phone number validation |
| User-Agent | avct/uasurfer | — | Browser/device detection |
| UUID | google/uuid | 1.6.0 | Unique identifiers |
| Crypto | golang.org/x/crypto | 0.31.0 | bcrypt password hashing |
| Testing | stretchr/testify | 1.9.0 | Assertions + mocks |
| Test DB | DATA-DOG/go-sqlmock | 1.5.0 | SQL mock (unit tests) |
| Test Redis | alicebob/miniredis | 2.5.0 | In-memory Redis mock |

---

## Architecture

### 4 Binary Entry Points

| Binary | Source | Purpose | Port | Key Goroutines |
|--------|--------|---------|------|----------------|
| `exchange-cli` | `cmd/exchange-cli/main.go` | 16 CLI cron commands (runs one, then exits) | N/A | None |
| `exchange-engine` | `cmd/exchange-engine/main.go` | Order matching daemon with worker pool | N/A | 1 dispatcher (BLPop) + 10 worker goroutines |
| `exchange-httpd` | `cmd/exchange-httpd/main.go` | REST API (Gin) + crash recovery | :8000 (public), :8001 (admin) | Gin servers + `UnmatchedOrderHandler.Match()` loop |
| `exchange-ws` | `cmd/exchange-ws/main.go` | Binance WebSocket market data listener | N/A | WS connection + reconnect loop |

**Startup sequence (`exchange-httpd`):**
```
di.NewContainer()           ← builds ~109 services lazily
  ├─ Config, Logger, DB, Redis, MQTT, RabbitMQ
  ├─ 26 Repositories (GORM)
  ├─ Domain services (business logic)
  └─ HTTPServer (Gin router with all routes)

httpServer.ListenAndServeAdmin(":8001")   ← admin routes (separate Gin instance)
unmatchedOrdersHandler.Match()            ← crash-recovery goroutine
httpServer.ListenAndServe(":8000")        ← public + MQTT auth routes
```

**Startup sequence (`exchange-engine`):**
```
di.NewContainer()
engine.Run(workerCount=10, shouldStartDispatcher=true)
  ├─ Dispatcher: BLPop("engine:queue:orders", 1s)
  └─ 10 workers waiting on workChan
```

**Startup sequence (`exchange-ws`):**
```
di.NewContainer()
ExternalExchangeWsService.GetActiveExternalExchangeWs()
  └─ Connects to Binance WebSocket with configured streams
      Supervised by .docker/supervisor/go-supervisord.conf:
        depth-stream         → ./ws depth
        ticker-trade-stream  → ./ws ticker trade
        kline-stream         → ./ws kline_1m kline_5m kline_1h kline_1d
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

---

## Internal Package Map (29 packages)

```
internal/
├── api/                    # REST handlers (Gin), middleware, admin handlers
│   ├── handler/            # 14 files: Public endpoint handlers /api/v1/*
│   ├── adminhandler/       # 4 files: Admin endpoint handlers /api/v1/admin/*
│   └── middleware/         # 4 files: Auth, AdminAuth, NonRequiredAuth, recovery
├── auth/                   # 7 files: Login, register, 2FA (TOTP), MQTT ACL
├── command/                # 20 files: 16 CLI ConsoleCommand implementations
├── communication/          # 6 files: MQTT publishing + RabbitMQ queue publishing
├── configuration/          # 4 files: App settings, app version service
├── country/                # 4 files: Country management
├── currency/               # 9+ files: Coins, pairs, prices, kline, gRPC candle client
├── di/                     # 8 files: DI container (121 service registrations)
├── engine/                 # 16 files: Order matching: worker pool + Redis sorted sets
├── externalexchange/       # 7 files: Binance REST API integration
├── externalexchangews/     # 3+ files: Binance WebSocket streams (depth, trade, kline, ticker)
├── jwt/                    # 2 files: JWT token creation/validation
├── livedata/               # 3 files: Redis live data (price, kline, depth, trade book)
├── mocks/                  # 82 files: Generated test mock implementations
├── order/                  # 20+ files: Order CRUD, matching, trades, bot aggregation
├── orderbook/              # 4 files: Orderbook aggregation from Redis + Binance REST
├── payment/                # 10 files: Withdrawal, deposit, auto-exchange, internal transfer
├── platform/               # 13 files: Infrastructure abstractions:
│                           #   jwt.go, logger.go, redis.go, sql.go, cache.go,
│                           #   mqtt.go, rabbitmq.go, http.go, bcrypt.go,
│                           #   config.go, error.go, ws.go, sentry.go
├── processor/              # 3 files: Real-time Binance WS data → Redis + MQTT
├── repository/             # 26 GORM data access repositories (see §Repository Layer)
├── response/               # 2 files: Standard API response envelope
├── transaction/            # 3 files: Bank transfer / ledger logic
├── user/                   # 22 files: User CRUD, profiles, KYC, 2FA, permissions, captcha
├── userbalance/            # 4 files: Balance management (lock/unlock/update)
├── userdevice/             # 2 files: Device fingerprinting from User-Agent
├── userwithdrawaddress/    # 4 files: Withdrawal address CRUD + wallet validation
├── wallet/                 # 6 files: External wallet service HTTP client
└── ws/                     # 2 files: WebSocket handler (Binance connection wrapper)
```

### Package Details

#### `internal/api/` — HTTP Server & Routing
- **Purpose**: Gin HTTP server setup, route registration, request handling
- **Key types**: `HTTPServer` (httpserver.go), handler structs, middleware funcs
- **Servers**: Public API on `:8000`, Admin API on `:8001` (separate Gin instances)
- **CORS**: Production restricts to `*.unitedbit.com`; non-prod allows `*`
- **Dependencies**: auth, order, payment, user, currency, configuration, userbalance, orderbook, userwithdrawaddress, communication

#### `internal/auth/` — Authentication & Authorization
- **Purpose**: Login, registration, 2FA (TOTP), password reset, MQTT broker auth
- **Key types**: `Service` (interface), `AuthEventsHandler`, `MqttAuthService`
- **Key functions**: `Login()`, `Register()`, `GetUser(token)`, `GetAdminUser(token)`, `ForgotPassword()`, `VerifyEmail()`
- **External deps**: Redis (forgot-password tokens), RabbitMQ (email notifications)
- **MQTT auth**: `/api/v1/emqtt/login`, `/acl`, `/superuser` — called by EMQX as webhook backend

#### `internal/command/` — CLI Commands
- **Purpose**: 16 scheduled/manual commands run by `exchange-cli` binary
- **Key interface**: `ConsoleCommand { Run(ctx context.Context, flags []string) }`
- **Dispatch**: `CommandService` maps command name → `ConsoleCommand` via DI
- **Dependencies**: varies per command (see §CLI Commands)

#### `internal/communication/` — Message Publishing
- **Purpose**: MQTT topic publishing + RabbitMQ queue publishing
- **Key types**: `MqttManager` (MQTT publish), `QueueManager` (RabbitMQ publish), `Service` (email/SMS orchestration)
- **MQTT topics**: `main/trade/*` (public), `main/trade/user/<channel>/*` (private)
- **RabbitMQ routing keys**: `messages.command.send` (email/SMS), `livedata.event.kline-created` (candle data)
- **External deps**: MQTT broker (EMQX), RabbitMQ

#### `internal/configuration/` — App Settings
- **Purpose**: Runtime configuration, app version management
- **Key types**: `Service`, `AppVersionRepository`
- **Dependencies**: DB (configurations table, app_version table), communication service

#### `internal/country/` — Country Management
- **Purpose**: Country list for user registration, KYC
- **Key types**: `Service`, `Repository`
- **Dependencies**: DB (countries table), config (supported countries)

#### `internal/currency/` — Coins, Pairs, Prices, Kline
- **Purpose**: Currency/coin management, trading pair management, price data, kline/OHLC sync
- **Key types**: `Service`, `PriceGenerator`, `KlineService`, `Repository`, `PairRepository`, `FavoritePairRepository`, `KlineSyncRepository`
- **Subpackage**: `candle/` — gRPC client for external candle service (`candle-app:50051`)
- **Dependencies**: DB, Redis (live data), gRPC (candle service)

#### `internal/di/` — Dependency Injection Container
- **Purpose**: Wire all ~109 services into a single `sarulabs/di` App-scoped container
- **Files**: `container.go` (constants + factory), `di_commands.go`, `di_http.go`, `di_repositories.go`, `di_services.go`, `di_trading.go`, `di_external.go`, `doc.go`
- **Registration**: All `builder.Add()` calls use string constants; registration order reflects dependency order
- **See**: §DI Container for complete service list

#### `internal/engine/` — Order Matching Engine
- **Purpose**: Worker pool that matches buy/sell orders using Redis sorted sets
- **Key types**: `Engine` (interface), `Order`, `OrderbookProvider` (interface), `ResultHandler` (interface), `Worker`, `Pool`, `Queue`
- **Architecture**: Dispatcher (BLPop) → Pool → Workers → OrderBook match → ResultHandler callback
- **Redis**: Sorted sets `order-book:bid:<pair>` / `order-book:ask:<pair>`, queue `engine:queue:orders`
- **See**: §Order Matching Engine for detailed flow

#### `internal/externalexchange/` — Binance REST API
- **Purpose**: Binance account operations, order submission, withdrawal checks, order sync
- **Key types**: `Service`, `OrderService`, `OrderFromExternalService`, `Repository`, `OrderRepository`
- **Operations**: Submit orders to Binance, check withdrawal status, retrieve external orders/trades
- **Dependencies**: HTTP client, Redis (exchange config cache), DB

#### `internal/externalexchangews/` — Binance WebSocket
- **Purpose**: Real-time market data streams from Binance
- **Key types**: `Service` (manages WS lifecycle), `binance.WsClient`
- **Streams**: depth, trade, kline (1m/5m/1h/1d), ticker — managed by supervisord
- **Dependencies**: WebSocket client, DataProcessor, currency service

#### `internal/jwt/` — JWT Token Helpers
- **Purpose**: Create and validate JWT tokens for user authentication
- **Key types**: `Service`
- **Dependencies**: Config (RSA keys path, passphrase, TTL), platform.JwtHandler

#### `internal/livedata/` — Redis Live Market Data
- **Purpose**: Store/retrieve real-time price, kline, depth, trade book data in Redis hashes
- **Key types**: `Service`
- **Redis pattern**: Hash `live_data:pair_currency:<pair>` with fields: `price`, `volume`, `change_price_percentage`, `trade_book`, `kline_<tf>`, `pre_kline_<tf>`, `depth_snapshot`, `order_book`
- **Dependencies**: Redis only

#### `internal/mocks/` — Test Mocks
- **Purpose**: Mock implementations for all service interfaces (82 files)
- **Pattern**: Each mock implements the corresponding interface with `testify/mock`

#### `internal/order/` — Order Management (Core Domain)
- **Purpose**: Order CRUD, matching orchestration, trade execution, bot aggregation, stop orders
- **Key types**: `Service`, `CreateManager`, `EventsHandler`, `PostOrderMatchingService`, `EngineCommunicator`, `EngineResultHandler`, `DecisionManager`, `ForceTrader`, `BotAggregationService`, `StopOrderSubmissionManager`, `InQueueOrderManager`, `UnmatchedOrdersHandler`, `AdminOrderManager`, `TradeEventsHandler`, `RedisManager`
- **Key flow**: Create → Decide (internal/external) → Engine queue → Match → PostOrderMatching → Balance update → MQTT publish
- **Dependencies**: DB, Redis, MQTT, engine, currency, userbalance, externalexchange

#### `internal/orderbook/` — Orderbook Aggregation
- **Purpose**: Build aggregated order book from Redis data + Binance REST snapshots
- **Key types**: `Service`
- **Dependencies**: Redis (livedata), HTTP client (Binance REST), currency service

#### `internal/payment/` — Payment Processing
- **Purpose**: Withdrawals, deposits, auto-exchange, internal transfers, email confirmation
- **Key types**: `Service`, `AutoExchangeManager`, `WithdrawEmailConfirmationManager`, `InternalTransferService`, `Repository`
- **Flow**: Pre-withdraw → 2FA verify → Email confirm → Lock balance → Wallet service → DB persist → MQTT
- **Dependencies**: DB, Redis (confirmation codes), wallet service, MQTT, communication

#### `internal/platform/` — Infrastructure Abstractions
- **Purpose**: Thin wrappers around external libraries; all other packages depend on these interfaces
- **Key types/files**:
  - `config.go` — `Configs` interface (wraps Viper)
  - `logger.go` — `Logger` interface (wraps zap)
  - `redis.go` — `RedisClient` interface (wraps go-redis)
  - `sql.go` — `SqlClient` (wraps `*gorm.DB`)
  - `cache.go` — `Cache` interface (in-process LRU + Redis)
  - `mqtt.go` — `MqttClient` interface (wraps paho)
  - `rabbitmq.go` — `RabbitMqClient` interface (wraps amqp091)
  - `http.go` — `HTTPClient` interface (wraps `net/http`)
  - `ws.go` — `WsClient` interface (wraps gorilla/websocket)
  - `jwt.go` — `JwtHandler` interface (wraps golang-jwt)
  - `bcrypt.go` — `PasswordEncoder` interface (wraps bcrypt)
  - `error.go` — Custom error types
  - `sentry.go` — Sentry integration

#### `internal/processor/` — Real-time Market Data Processor
- **Purpose**: Transform Binance WS events into Redis updates + MQTT publishes
- **Key types**: `Processor` (interface), `WsDataProcessor` (implementation)
- **Methods**: `ProcessTrade()`, `ProcessDepth()`, `ProcessKline()`, `ProcessTicker()`
- **Ticker processing**: Also triggers `StopOrderSubmissionManager.Check()` and publishes to Redis pub/sub channel `channel:ticker`
- **Dependencies**: Redis, livedata, MQTT, kline, orderbook, currency, stop order manager

#### `internal/response/` — API Response Envelope
- **Purpose**: Standard `{status, message, data}` response struct
- **Key types**: `APIResponse`
- **Usage**: `response.Success(c, data)`, `response.Error(c, code, msg)`

#### `internal/transaction/` — Bank Transfers
- **Purpose**: Fiat/bank transfer logic and ledger
- **Key types**: `Service`
- **Dependencies**: DB

#### `internal/user/` — User Management
- **Purpose**: User CRUD, profiles, KYC levels, 2FA, permissions, captcha, phone verification
- **Key types**: `Service`, `LevelService`, `ConfigService`, `LoginHistoryService`, `Repository`, `ProfileRepository`, `ProfileImageRepository`, `LevelRepository`, `ConfigRepository`, `PermissionManager`, `TwoFaManager`, `RecaptchaManager`, `UbCaptchaManager`, `ForgotPasswordManager`, `PhoneConfirmationManager`, `PermissionRepository`, `UsersPermissionsRepository`
- **Dependencies**: DB, Redis (OTP codes, reset tokens), communication, jwt, platform

#### `internal/userbalance/` — Balance Management
- **Purpose**: User balance operations (lock, unlock, credit, debit) for trading and withdrawals
- **Key types**: `Service`, `Repository`, `UserWalletBalanceRepository`
- **Key operations**: `LockBalance()`, `UpdateBalances()` (inside DB transactions), `GetAllBalances()`
- **Dependencies**: DB, currency, wallet, user, permission

#### `internal/userdevice/` — Device Fingerprinting
- **Purpose**: Parse User-Agent headers to extract browser, OS, device info
- **Key types**: `Service`
- **Dependencies**: avct/uasurfer

#### `internal/userwithdrawaddress/` — Withdrawal Addresses
- **Purpose**: CRUD for user's saved withdrawal addresses
- **Key types**: `Service`, `Repository`
- **Dependencies**: DB, currency, wallet service (address validation)

#### `internal/wallet/` — External Wallet Service Client
- **Purpose**: HTTP client for external blockchain wallet service
- **Key types**: `Service`, `AuthorizationService`
- **Operations**: Create withdrawal, generate address, get balance, validate address
- **Auth**: JWT-based auth cached in Redis (`wallet:auth:<token>`)
- **Dependencies**: HTTP client, Redis, config

#### `internal/ws/` — WebSocket Handler
- **Purpose**: Low-level WebSocket connection management for Binance streams
- **Key types**: `Handler`
- **Dependencies**: gorilla/websocket, platform.WsClient

---

## Order Matching Engine — Detailed

### Architecture
```
exchange-engine binary starts
  └─ engine.Run(workerCount=10, shouldStartDispatcher=true)

Dispatcher goroutine:
  └─ BLPop("engine:queue:orders", 1s timeout)  ← blocks until order arrives
       └─ JSON-decode → Order struct → work{order}
            └─ pool.addWork(&work) → sends to worker.workChan

Worker goroutine (×10):
  └─ processOrder(order)
       ├─ ob := newOrderBook(pair, redisOrderBookProvider)
       ├─ ob.loadOrders(oppositeSide)
       │    └─ Redis ZRangeByScoreWithScores("order-book:bid:<pair>" or "order-book:ask:<pair>")
       ├─ processLimitOrder() or processMarketOrder()
       │    └─ price-time matching loop → tradeOrders()
       │         └─ builds doneOrders[], partialOrder
       └─ callBackManager.callBack(doneOrders, partialOrder)
            └─ EngineResultHandler.CallBack()
                 └─ PostOrderMatchingService.Handle()
                      ├─ DB: persist trades, update order status
                      ├─ UserBalanceService.UpdateBalances()
                      ├─ TradeEventsHandler.OnTrade() → BotAggregationService
                      ├─ LiveDataService.UpdateTradeBook()
                      ├─ MqttManager.PublishTrades()
                      └─ MqttManager.PublishOrderToOpenOrders()
                 └─ ob.rewriteOrderBook()
                      └─ Redis TxPipeline: ZRem done + ZAdd partial
```

### Redis Sorted Set Usage
| Key Pattern | Type | Score | Member | Matching Direction |
|-------------|------|-------|--------|-------------------|
| `order-book:bid:<pair>` | Sorted Set | price (float64) | JSON-serialized `Order` | ZPopMax (highest bid first) |
| `order-book:ask:<pair>` | Sorted Set | price (float64) | JSON-serialized `Order` | ZPopMin (lowest ask first) |
| `engine:queue:orders` | List | N/A | JSON-serialized `Order` | RPUSH submit / BLPop consume / LPush priority |

### Order Lifecycle
1. **Submit**: Client POST → `OrderService.Create()` → `OrderCreateManager.Create()` (DB persist + lock balance)
2. **Route**: `OrderEventsHandler.OnOrderCreated()` → `DecisionManager.Decide()` → INTERNAL or EXTERNAL
3. **Queue** (internal): `EngineCommunicator.SubmitToEngine()` → `Engine.SubmitOrder()` → Redis RPUSH
4. **Dispatch**: Engine dispatcher BLPop → deserialize → send to worker pool
5. **Match**: Worker loads opposite orderbook from Redis sorted set → price-time matching loop
6. **Fill**: Fully filled orders → `doneOrders[]`; partially filled → `partialOrder`
7. **Settle**: `PostOrderMatchingService.Handle()` → DB trade records + balance updates + MQTT publish
8. **Rewrite**: Atomically update Redis sorted sets (ZRem done, ZAdd partial with updated amount)

### Partial Fill Handling
- When a limit order partially fills, the remaining quantity is stored back in the Redis sorted set via `ZAdd`
- The `partialOrder` retains the original order ID with updated `Amount` reflecting the unfilled remainder
- Both the filled portion (as a trade) and the partial update are published via MQTT

### Crash Recovery
- `UnmatchedOrderHandler.Match()` runs as a background goroutine in `exchange-httpd`
- Periodically scans open orders from DB that are missing from Redis queue/orderbook
- Re-submits via `engine.RetrieveOrder()` which uses `LPush` (head of queue) for priority

---

## HTTP API — Complete Endpoint Reference

### Middleware Chain

**Public API (`:8000`):**
```
globalRecover (panic recovery)
  → CORS (origin validation)
    → Route group middleware:
      → [No auth] for public routes
      → NonRequiredAuthMiddleware + AuthMiddleware for optional-auth routes
      → AuthMiddleware for protected routes
```

**Admin API (`:8001`):**
```
globalRecover (panic recovery)
  → CORS (origin validation)
    → AdminAuthMiddleware (validates JWT + checks admin role)
```

### Authentication Endpoints — No Auth Required

| Method | Path | Handler | Purpose |
|--------|------|---------|---------|
| POST | `/api/v1/auth/login` | `handler.Login` | Email/password login → JWT token |
| POST | `/api/v1/auth/register` | `handler.Register` | New user registration → JWT token |
| POST | `/api/v1/auth/forgot-password` | `handler.ForgotPassword` | Send password reset email |
| POST | `/api/v1/auth/forgot-password/update` | `handler.ForgotPasswordUpdate` | Complete password reset with token |
| POST | `/api/v1/auth/verify` | `handler.VerifyEmail` | Verify email address with token |

### MQTT Broker Auth Endpoints — No Auth Required (called by EMQX)

| Method | Path | Handler | Purpose |
|--------|------|---------|---------|
| POST | `/api/v1/emqtt/login` | `handler.MqttLogin` | MQTT broker login webhook |
| POST | `/api/v1/emqtt/acl` | `handler.MqttACL` | MQTT ACL check (form-encoded) |
| POST | `/api/v1/emqtt/superuser` | `handler.MqttSuperUser` | MQTT superuser check (form-encoded) |

### Public Data Endpoints — No Auth Required

| Method | Path | Handler | Purpose |
|--------|------|---------|---------|
| GET | `/api/v1/main-data/check` | `handler.Check` | Health check → `{"status":true,"message":"OK"}` |
| GET | `/api/v1/main-data/country-list` | `handler.Countries` | List all countries |
| GET | `/api/v1/main-data/common` | `handler.GetRecaptchaKey` | Get reCAPTCHA public key |
| GET | `/api/v1/main-data/version` | `handler.GetAppVersion` | App version info by platform |
| POST | `/api/v1/main-data/contact-us` | `handler.ContactUs` | Submit contact form |

### Currency Endpoints — Mixed Auth

| Method | Path | Auth | Handler | Purpose |
|--------|------|------|---------|---------|
| GET | `/api/v1/currencies` | ❌ | `handler.GetCurrencies` | List all currencies/coins |
| GET | `/api/v1/currencies/pairs` | ❌ | `handler.GetPairs` | List all trading pairs |
| GET | `/api/v1/currencies/pairs-list` | ❌ | `handler.GetPairsList` | Pairs list (alternate format) |
| GET | `/api/v1/currencies/pairs-statistic` | ❌ | `handler.GetPairsStatistic` | Pair statistics (price, volume) |
| GET | `/api/v1/currencies/pairs-ratio` | ❌ | `handler.GetPairRatio` | Pair price ratio |
| GET | `/api/v1/currencies/fees` | ❌ | `handler.GetFees` | Trading fee structure |
| POST | `/api/v1/currencies/favorite` | ✅ | `handler.AddOrRemoveFavoritePair` | Toggle favorite pair |
| GET | `/api/v1/currencies/favorite-pairs` | ✅ | `handler.GetFavoritePairs` | User's favorite pairs |

### Order Book Endpoints — No Auth Required

| Method | Path | Handler | Purpose |
|--------|------|---------|---------|
| GET | `/api/v1/order-book` | `handler.OrderBook` | Aggregated order book by pair |
| GET | `/api/v1/trade-book` | `handler.TradeBook` | Recent trades by pair |

### Order Endpoints — Auth Required

| Method | Path | Handler | Purpose |
|--------|------|---------|---------|
| POST | `/api/v1/order/create` | `handler.CreateOrder` | Create buy/sell order |
| POST | `/api/v1/order/cancel` | `handler.CancelOrder` | Cancel open order |
| GET | `/api/v1/order/open-orders` | `handler.OpenOrders` | User's open orders (paginated) |
| GET | `/api/v1/order/history` | `handler.OrdersHistory` | Order history (paginated) |
| GET | `/api/v1/order/full-history` | `handler.FullOrdersHistory` | Complete order history |
| GET | `/api/v1/order/detail` | `handler.GetOrderDetail` | Single order detail |

### Trade Endpoints — Auth Required

| Method | Path | Handler | Purpose |
|--------|------|---------|---------|
| GET | `/api/v1/trade/history` | `handler.TradesHistory` | Trade history (paginated) |
| GET | `/api/v1/trade/full-history` | `handler.FullTradesHistory` | Complete trade history |

### User Balance Endpoints — Auth Required

| Method | Path | Handler | Purpose |
|--------|------|---------|---------|
| GET | `/api/v1/user-balance/pair-balance` | `handler.PairBalances` | Balance for specific pair |
| GET | `/api/v1/user-balance/balance` | `handler.AllBalances` | All wallet balances |
| GET | `/api/v1/user-balance/withdraw-deposit` | `handler.WithdrawAndDeposit` | Withdraw/deposit history |
| POST | `/api/v1/user-balance/auto-exchange` | `handler.SetAutoExchange` | Set auto-exchange coin |

### Crypto Payment Endpoints — Auth Required

| Method | Path | Handler | Purpose |
|--------|------|---------|---------|
| GET | `/api/v1/crypto-payment` | `handler.GetPayments` | List payments/withdrawals |
| GET | `/api/v1/crypto-payment/detail` | `handler.GetPaymentDetail` | Payment detail |
| POST | `/api/v1/crypto-payment/pre-withdraw` | `handler.PreWithdraw` | Validate withdrawal params |
| POST | `/api/v1/crypto-payment/withdraw` | `handler.Withdraw` | Execute withdrawal |
| POST | `/api/v1/crypto-payment/cancel` | `handler.Cancel` | Cancel pending withdrawal |

### Withdrawal Address Endpoints — Auth Required

| Method | Path | Handler | Purpose |
|--------|------|---------|---------|
| GET | `/api/v1/withdraw-address` | `handler.GetWithdrawAddresses` | List saved addresses |
| GET | `/api/v1/withdraw-address/former-addresses` | `handler.GetFormerAddresses` | Historical addresses |
| POST | `/api/v1/withdraw-address/new` | `handler.NewWithdrawAddress` | Create withdrawal address |
| POST | `/api/v1/withdraw-address/favorite` | `handler.AddToFavorites` | Mark address as favorite |
| POST | `/api/v1/withdraw-address/delete` | `handler.Delete` | Delete withdrawal address |

### User Endpoints — Auth Required

| Method | Path | Handler | Purpose |
|--------|------|---------|---------|
| POST | `/api/v1/user/set-user-profile` | `handler.SetUserProfile` | Update profile |
| GET | `/api/v1/user/get-user-profile` | `handler.GetUserProfile` | Get profile |
| GET | `/api/v1/user/user-data` | `handler.GetUserData` | Get full user data |
| GET | `/api/v1/user/google-2fa-barcode` | `handler.Get2FaBarcode` | Get 2FA QR code/secret |
| POST | `/api/v1/user/google-2fa-enable` | `handler.Enable2Fa` | Enable 2FA |
| POST | `/api/v1/user/google-2fa-disable` | `handler.Disable2Fa` | Disable 2FA |
| POST | `/api/v1/user/change-password` | `handler.ChangePassword` | Change password |
| POST | `/api/v1/user/sms-send` | `handler.SendSms` | Send SMS OTP |
| POST | `/api/v1/user/sms-enable` | `handler.EnableSms` | Enable SMS 2FA |
| POST | `/api/v1/user/sms-disable` | `handler.DisableSms` | Disable SMS 2FA |
| POST | `/api/v1/user/send-verification-email` | `handler.SendVerificationEmail` | Resend verification |

### User Profile Image Endpoints — Auth Required

| Method | Path | Handler | Purpose |
|--------|------|---------|---------|
| POST | `/api/v1/user-profile-image/multiple-upload` | `handler.MultipleUpload` | Upload KYC images (multipart) |
| POST | `/api/v1/user-profile-image/delete` | `handler.DeleteProfileImage` | Delete profile image |

### Admin Endpoints (Port 8001) — Admin Auth Required

| Method | Path | Handler | Purpose |
|--------|------|---------|---------|
| POST | `/api/v1/order/fulfill` | `adminhandler.FulFillOrder` | Admin force-fill order |
| POST | `/api/v1/payment/callback` | `adminhandler.Callback` | Wallet callback handler |
| POST | `/api/v1/payment/update-withdraw` | `adminhandler.UpdateWithdraw` | Update withdrawal status |
| POST | `/api/v1/payment/update-deposit` | `adminhandler.UpdateDeposit` | Update deposit status |
| POST | `/api/v1/user-balance/update` | `adminhandler.UpdateUserBalance` | Admin balance adjustment |

### API Response Format
```json
// Success
{ "status": true, "message": "OK", "data": { ... } }

// Error
{ "status": false, "message": "error description", "data": { "field": "validation error" } }
```

**HTTP status codes**: 200 (success), 400/422 (validation), 401 (unauthorized), 500 (server error)

**Content types**: JSON (`application/json`), form (`application/x-www-form-urlencoded` for MQTT auth), multipart (image uploads)

**Auth header**: `Authorization: Bearer <jwt_token>`

---

## CLI Commands (16 total)

| # | Command | Cron Schedule | Purpose | DB Write | Redis Write | External API |
|---|---------|--------------|---------|----------|-------------|--------------|
| 1 | `set-user-level` | `0 1 * * *` (daily 1AM) | Calculate trading level from 30-day volume | ✓ users | — | — |
| 2 | `initialize-balance` | Manual | Create zero-balance records for all users × all active coins | ✓ user_balances | — | — |
| 3 | `generate-address` | `0 5 * * *` (daily 5AM) | Generate wallet addresses for users missing them | ✓ user_balances | — | Wallet service |
| 4 | `retrieve-open-orders` | `*/15 * * * *` (15min) | Sync external open orders to Redis; trigger stop orders | — | ✓ stop orders | Engine queue |
| 5 | `submit-bot-orders` | `* * * * *` (1min) | Read aggregated bot trades from Redis, submit to Binance | ✓ external orders | ✓ read+del | Binance REST |
| 6 | `sync-kline` | `* * * * *` (1min) | Fetch OHLC candle data from gRPC/Binance, publish via MQTT | ✓ ohlc_sync | — | gRPC/Binance |
| 7 | `generate-kline-sync` | `0 2 * * *` (daily 2AM) | Create kline sync jobs for all pairs × timeframes | ✓ ohlc_sync | — | — |
| 8 | `check-withdrawals` | `*/10 * * * *` (10min) | Check withdrawal status on Binance, update payment records | ✓ payments | — | Binance REST |
| 9 | `update-orders-from-external` | `*/30 * * * *` (30min) | Sync filled/cancelled external orders back to local DB | ✓ ext orders | — | Binance REST |
| 10 | `retrieve-external-orders` | Disabled (was `*/20 * * * *`) | Cache external exchange orders into Redis | — | ✓ orders | — |
| 11 | `delete-cache` | Manual | Flush all Redis keys | — | ✓ FLUSHDB | — |
| 12 | `ub-captcha-generate-keys` | Manual | Generate RSA key pair for UB captcha encryption | — | — | — |
| 13 | `ub-captcha-encryption` | Manual | Encrypt captcha data with RSA public key (debug tool) | — | — | — |
| 14 | `ub-captcha-decryption` | Manual | Decrypt captcha data with RSA private key (debug tool) | — | — | — |
| 15 | `ub-update-user-wallet-balances` | `0 0 * * *` (daily midnight) | Sync wallet service balances to local DB | ✓ wallet balances | — | Wallet service |
| 16 | `ws-health-check` | Supervisor (30–60s) | Check WebSocket stream health, restart via supervisord if unhealthy | — | — | Supervisor API |

---

## DI Container — Complete Service Registry (~109 services)

### Infrastructure Layer (no app dependencies)

| Constant | Type | Dependencies |
|----------|------|-------------|
| `ConfigService` | `platform.Configs` | none (Viper YAML) |
| `LoggerService` | `platform.Logger` | ConfigService |
| `cacheService` | `platform.Cache` | ConfigService, LoggerService |
| `dbClient` | `*gorm.DB` | ConfigService |
| `RedisClient` | `platform.RedisClient` | ConfigService |
| `wsClient` | `platform.WsClient` | none |
| `mqttClient` | `platform.MqttClient` | ConfigService, LoggerService |
| `httpClient` | `platform.HTTPClient` | none |
| `rabbitmqClient` | `platform.RabbitMqClient` | ConfigService, LoggerService |
| `jwtHandler` | `platform.JwtHandler` | none |
| `passwordEncoder` | `platform.PasswordEncoder` | none |

### Communication Layer

| Constant | Type | Dependencies |
|----------|------|-------------|
| `mqttManager` | `communication.MqttManager` | mqttClient |
| `queueManager` | `communication.QueueManager` | rabbitmqClient, LoggerService |
| `communicationService` | `communication.Service` | queueManager, LoggerService |

### Repository Layer (26 repositories — all depend on dbClient)

| Constant | Type | Extra Dependencies |
|----------|------|-------------------|
| `orderRepository` | `order.Repository` | — |
| `userRepository` | `user.Repository` | cacheService |
| `currencyRepository` | `currency.Repository` | — |
| `pairRepository` | `currency.PairRepository` | cacheService |
| `favoritePairRepository` | `currency.FavoritePairRepository` | — |
| `userBalanceRepository` | `userbalance.Repository` | — |
| `paymentRepository` | `payment.Repository` | — |
| `tradeRepository` | `order.TradeRepository` | — |
| `userProfileRepository` | `user.ProfileRepository` | — |
| `profileImageRepository` | `user.ProfileImageRepository` | — |
| `countryRepository` | `country.Repository` | cacheService |
| `userLevelRepository` | `user.LevelRepository` | cacheService |
| `loginHistoryRepository` | `user.LoginHistoryRepository` | — |
| `userWalletBalanceRepository` | `userbalance.UserWalletBalanceRepository` | — |
| `klineSyncRepository` | `currency.KlineSyncRepository` | — |
| `appVersionRepository` | `configuration.AppVersionRepository` | cacheService |
| `internalTransferRepository` | `payment.InternalTransferRepository` | — |
| `userConfigRepository` | `user.ConfigRepository` | — |
| `userWithdrawAddressRepository` | `userwithdrawaddress.Repository` | — |
| `tradeFromExternalRepository` | `externalexchange.TradeFromExternalRepository` | — |
| `orderFromExternalRepository` | `externalexchange.OrderFromExternalRepository` | — |
| `configurationRepository` | `configuration.Repository` | — |
| `externalExchangeOrderRepository` | `externalexchange.OrderRepository` | — |
| `externalExchangeRepository` | `externalexchange.Repository` | — |
| `permissionRepository` | `user.PermissionRepository` | — |
| `usersPermissionsRepository` | `user.UsersPermissionsRepository` | — |

### Domain Services

| Constant | Type | Key Dependencies |
|----------|------|-----------------|
| `liveDataService` | `livedata.Service` | RedisClient |
| `priceGenerator` | `currency.PriceGenerator` | liveDataService, klineService, pairRepository |
| `klineService` | `currency.KlineService` | klineSyncRepository, liveDataService, candleGRPCClient |
| `candleGRPCClient` | `candle.CandleGRPCClient` | ConfigService, LoggerService |
| `currencyService` | `currency.Service` | currencyRepository, pairRepository, liveDataService, priceGenerator, klineService, favoritePairRepository, ConfigService, LoggerService |
| `orderbookService` | `orderbook.Service` | liveDataService, httpClient, currencyService, LoggerService |
| `countryService` | `country.Service` | countryRepository, ConfigService |
| `jwtService` | `jwt.Service` | ConfigService, jwtHandler |
| `userLevelService` | `user.LevelService` | userLevelRepository |
| `loginHistoryService` | `user.LoginHistoryService` | loginHistoryRepository |
| `userConfigService` | `user.ConfigService` | userConfigRepository |
| `twoFaManager` | `user.TwoFaManager` | none |
| `permissionManager` | `user.PermissionManager` | usersPermissionsRepository, permissionRepository |
| `ubCaptchaManager` | `user.UbCaptchaManager` | LoggerService |
| `recaptchaManager` | `user.RecaptchaManager` | httpClient, ConfigService, LoggerService, ubCaptchaManager |
| `forgotPasswordManager` | `user.ForgotPasswordManager` | RedisClient, communicationService, ConfigService |
| `phoneConfirmationManager` | `user.PhoneConfirmationManager` | RedisClient, communicationService |
| `userService` | `user.Service` | dbClient, userRepository, userProfileRepository, profileImageRepository, countryService, twoFaManager, passwordEncoder, communicationService, phoneConfirmationManager, jwtService, ConfigService, LoggerService |
| `authEventsHandler` | `auth.EventsHandler` | loginHistoryService, communicationService, userService, ConfigService, LoggerService |
| `authService` | `auth.Service` | dbClient, userRepository, userLevelService, permissionManager, userBalanceService, jwtService, passwordEncoder, communicationService, authEventsHandler, forgotPasswordManager, recaptchaManager, twoFaManager, ConfigService, LoggerService |
| `mqttAuthService` | `auth.MqttAuthService` | authService, ConfigService, LoggerService |
| `walletAuthorizationService` | `wallet.AuthorizationService` | RedisClient, httpClient, ConfigService, LoggerService |
| `walletService` | `wallet.Service` | walletAuthorizationService, httpClient, ConfigService, LoggerService |
| `userBalanceService` | `userbalance.Service` | dbClient, userBalanceRepository, currencyService, priceGenerator, permissionManager, walletService, userService, userWalletBalanceRepository, ConfigService, LoggerService |
| `userWithdrawAddressService` | `userwithdrawaddress.Service` | dbClient, userWithdrawAddressRepository, currencyService, walletService, LoggerService |
| `configurationService` | `configuration.Service` | configurationRepository, appVersionRepository, communicationService, ConfigService, LoggerService |
| `tradeService` | `order.TradeService` | tradeRepository |
| `internalTransferService` | `payment.InternalTransferService` | internalTransferRepository |
| `withdrawEmailConfirmationManager` | `payment.WithdrawEmailConfirmationManager` | RedisClient, communicationService |
| `autoExchangeManager` | `payment.AutoExchangeManager` | dbClient, paymentRepository, orderCreateManager, orderEventsHandler, userService, currencyService, priceGenerator, LoggerService |
| `paymentService` | `payment.Service` | dbClient, paymentRepository, currencyService, walletService, userConfigService, twoFaManager, withdrawEmailConfirmationManager, permissionManager, userService, userBalanceService, userWithdrawAddressService, communicationService, priceGenerator, internalTransferService, externalExchangeService, autoExchangeManager, mqttManager, ConfigService, LoggerService |

### Order & Trading Pipeline

| Constant | Type | Key Dependencies |
|----------|------|-----------------|
| `orderRedisManager` | `order.RedisManager` | RedisClient |
| `decisionManager` | `order.DecisionManager` | ConfigService |
| `forceTrader` | `order.ForceTrader` | priceGenerator, currencyService |
| `botAggregationService` | `order.BotAggregationService` | RedisClient |
| `tradeEventsHandler` | `order.TradeEventsHandler` | botAggregationService, ConfigService, LoggerService |
| `engineService` | `engine.Engine` | RedisClient, orderbookProvider, engineResultHandler, LoggerService, ConfigService |
| `engineCommunicator` | `order.EngineCommunicator` | forceTrader, engineService |
| `EngineResultHandler` | `order.EngineResultHandler` | postOrderMatchingService |
| `orderEventsHandler` | `order.EventsHandler` | orderRedisManager, decisionManager, mqttManager, externalExchangeOrderService, engineCommunicator, postOrderMatchingService, LoggerService |
| `postOrderMatchingService` | `order.PostOrderMatchingService` | dbClient, orderRepository, userBalanceService, forceTrader, priceGenerator, tradeEventsHandler, mqttManager, RedisClient, currencyService, userService, userLevelService, ConfigService, LoggerService |
| `orderCreateManager` | `order.CreateManager` | dbClient, userBalanceService, userLevelService, priceGenerator |
| `orderService` | `order.Service` | dbClient, orderRepository, orderCreateManager, orderEventsHandler, currencyService, priceGenerator, userBalanceService, orderRedisManager, userConfigService, permissionManager, adminOrderManager, engineCommunicator, ConfigService, LoggerService |
| `adminOrderManager` | `order.AdminOrderManager` | currencyService, klineService, priceGenerator, postOrderMatchingService, stopOrderSubmissionManager, orderEventsHandler, LoggerService |
| `StopOrderSubmissionManager` | `order.StopOrderSubmissionManager` | dbClient, orderRepository, liveDataService, orderRedisManager, orderEventsHandler, LoggerService |
| `inQueueOrderManager` | `order.InQueueOrderManager` | engineService, LoggerService |
| `UnmatchedOrderHandler` | `order.UnmatchedOrdersHandler` | RedisClient, orderRepository, engineCommunicator, ConfigService, LoggerService |

### External Exchange

| Constant | Type | Key Dependencies |
|----------|------|-----------------|
| `externalExchangeService` | `externalexchange.Service` | externalExchangeRepository, RedisClient, httpClient, priceGenerator, ConfigService, LoggerService |
| `externalExchangeOrderService` | `externalexchange.OrderService` | externalExchangeOrderRepository, externalExchangeService, LoggerService |
| `orderFromExternalService` | `externalexchange.OrderFromExternalService` | orderFromExternalRepository, tradeFromExternalRepository |
| `ExternalExchangeWsService` | `externalexchangews.Service` | wsClient, wsDataProcessor, ConfigService, LoggerService, currencyService |
| `wsDataProcessor` | `processor.Processor` | RedisClient, liveDataService, priceGenerator, klineService, orderbookService, mqttManager, stopOrderSubmissionManager, inQueueOrderManager, queueManager, LoggerService, currencyService |

### HTTP Server

| Constant | Type | Key Dependencies |
|----------|------|-----------------|
| `HTTPServer` | `api.HTTPServer` | All handler dependencies injected via service layer |

### CLI Commands (13 registered in DI)

Registered in `di_commands.go`: `SetUserLevelCommand`, `InitializeBalanceCommand`, `GenerateAddressCommand`, `RetrieveOpenOrdersCommand`, `SubmitBotOrdersCommand`, `SyncKlineCommand`, `GenerateKlineSyncCommand`, `CheckWithdrawalsCommand`, `UpdateExternalOrdersCommand`, `RetrieveExternalOrdersCommand`, `DeleteCacheCommand`, `UpdateUserWalletBalancesCommand`, `WsHealthCheckCommand`

(Captcha commands are inline in main.go, not registered in DI)

---

## Repository Layer — All 26 Repositories

| # | Repository | Entity Struct | DB Table | Key Methods |
|---|-----------|--------------|----------|-------------|
| 1 | `AppVersionRepository` | `AppVersion` | `app_version` | `FindNewAppVersion`, `FindNewAppVersions` |
| 2 | `ConfigurationRepository` | `Configuration` | `configurations` | GORM standard CRUD |
| 3 | `CountryRepository` | `Country` | `countries` | `All`, `GetCountryByID` |
| 4 | `CurrencyRepository` | `Coin` | `currencies` | `GetCoinsAlphabetically`, `GetActiveCoins`, `GetCoinByCode` |
| 5 | `PairRepository` | `Pair` | `pair_currencies` | `GetPairByID`, `GetPairByName`, `GetActivePairCurrenciesList`, `GetAllPairs` |
| 6 | `KlineSyncRepository` | `KlineSync` | `ohlc_sync` | `Create`, `Update`, `GetKlineSyncsByStatusAndLimit` |
| 7 | `FavoritePairRepository` | `FavoritePair` | `user_favorite_pair_currency` | `Create`, `Delete`, `GetFavoritePair`, `GetUserFavoritePairs` |
| 8 | `OrderRepository` | `Order` | `orders` | `GetOrdersByIds`, `GetOrderByID`, `GetUserOpenOrders`, `GetLeafOrders`, `GetUserTradedOrders` |
| 9 | `TradeRepository` | `Trade` | `trades` | `Create`, `GetTradesOfUserBetweenTimes`, `GetBotTradesByIDAndCreatedAtGreaterThan` |
| 10 | `PaymentRepository` | `Payment` | `crypto_payments` | `GetUserPayments`, `GetPaymentByID`, `GetInProgressWithdrawalsInExternalExchange` |
| 11 | `UserBalanceRepository` | `UserBalance` | `user_balances` | `GetBalanceOfUserByCoinID`, `GetBalancesOfUsersForCoins`, `GetUserAllBalances` |
| 12 | `UserWalletBalanceRepository` | `UserWalletBalance` | `user_wallet_balance` | `FindUserWalletBalance` |
| 13 | `UserWithdrawAddressRepository` | `UserWithdrawAddress` | `user_withdraw_address` | `Create`, `GetUserWithdrawAddressesByCoinID`, `GetUserWithdrawAddresses` |
| 14 | `ExternalExchangeRepository` | `ExternalExchange` | `external_exchanges` | `GetEnabledPrivateExternalExchange` |
| 15 | `ExternalExchangeOrderRepository` | `ExternalExchangeOrder` | `external_exchange_orders` | `Create`, `Update`, `GetExternalExchangeOrdersLastTradeIds` |
| 16 | `OrderFromExternalRepository` | `OrderFromExternal` | `order_from_external` | `GetLastOrderFromExternalByPairID`, `Create`, `GetOrderByExternalOrderID` |
| 17 | `TradeFromExternalRepository` | `TradeFromExternal` | `trade_from_external` | `GetLastTradeFromExternalByPairID`, `Create` |
| 18 | `InternalTransferRepository` | `InternalTransfer` | `crypto_internal_transfer` | `GetFromExternalInProgressTransfers`, `GetInternalTransferByID`, `Update` |
| 19 | `UserRepository` | `User` | `users` | `GetUserByUsername`, `GetUserByID`, `GetUsersByPagination`, `GetUsersDataForOrderMatching` |
| 20 | `UserProfileRepository` | `Profile` | `user_profiles` | `GetProfileByUserID`, `GetProfileByUserIDUsingTx` |
| 21 | `ProfileImageRepository` | `ProfileImage` | `user_profile_image` | `GetImagesByIds`, `GetImageByID`, `GetLatestImagesDataByProfileID` |
| 22 | `UserLoginHistoryRepository` | `LoginHistory` | `user_login_history` | `Create`, `GetLastLoginHistoryByUserID` |
| 23 | `UserLevelRepository` | `Level` | `user_levels` | `GetAllLevels`, `GetLevelByID`, `GetLevelByCode` |
| 24 | `UserConfigRepository` | `Config` | `user_configs` | `GetUserConfigByUserID` |
| 25 | `UsersPermissionsRepository` | `UsersPermissions` | `users_permissions` | `GetUserPermissions` |
| 26 | `PermissionRepository` | `Permission` | `user_permissions` | `GetAllPermissions` |

---

## Real-time Data — MQTT Topics, Redis Pub/Sub, WebSocket

### MQTT Topics (QoS 0, non-retained, via EMQX v4)

**Public market data (broadcast to all):**

| Topic | Payload | Publisher | Trigger |
|-------|---------|----------|---------|
| `main/trade/trade-book/<pair>` | JSON trade list | `PostOrderMatchingService`, `Processor.ProcessTrade()` | Internal match or Binance trade |
| `main/trade/kline/<timeframe>/<pair>` | JSON `RedisKline` | `Processor.ProcessKline()` | New candle from Binance WS |
| `main/trade/ticker` | JSON ticker snapshot | `Processor.ProcessTicker()` | Binance ticker update |
| `main/trade/order-book/<pair>` | JSON order book | `Processor.ProcessDepth()` | Binance depth update |

**Private user topics (per-user, EMQX ACL-protected):**

| Topic Pattern | Payload | Publisher | Trigger |
|---------------|---------|----------|---------|
| `main/trade/user/<privateChannel>/open-orders/` | JSON order update | `OrderEventsHandler`, `PostOrderMatchingService` | Order created/filled/partial |
| `main/trade/user/<privateChannel>/crypto-payments/` | JSON payment | `PaymentService` | Withdrawal/deposit status change |

`<privateChannel>` is a per-user identifier from JWT claims. MQTT auth is handled by EMQX webhook calling `/api/v1/emqtt/login`, `/acl`, `/superuser`.

### Redis Pub/Sub

| Channel | Publisher | Subscriber | Data |
|---------|----------|-----------|------|
| `channel:ticker` | `Processor.ProcessTicker()` | Any Redis subscriber | Ticker snapshot (real-time price updates) |

### Redis Data Structures

| Key Pattern | Type | Purpose | TTL |
|-------------|------|---------|-----|
| `engine:queue:orders` | List | Order queue (RPUSH submit, BLPop consume, LPush priority) | None |
| `order-book:bid:<pair>` | Sorted Set | Bid orderbook (score=price, member=JSON Order) | None |
| `order-book:ask:<pair>` | Sorted Set | Ask orderbook (score=price, member=JSON Order) | None |
| `queue:stop:order:<type>:<pair>` | Sorted Set | Stop orders by price trigger (score=stop-price) | None |
| `live_data:pair_currency:<pair>` | Hash | All live market data (price, volume, kline, depth, trades) | None |
| `wallet:auth` | Hash | Wallet service JWT (fields: token, expiredAt) | 5 hours |
| `withdraw-confirmation:<userId>` | Hash+Expire | Withdrawal email OTP (fields: code, expiredAt, coin, amount, address) | 3 hours |
| `forgot-password:<userId>` | Hash+Expire | Password reset (fields: userId, code, expiredAt) | 3 hours |
| `phone-confirmation:<userId>` | Hash+Expire | Phone verification OTP (fields: userId, code, expiredAt, phone) | 3 hours |

---

## External Exchange — Binance Integration

### REST API (`internal/externalexchange/`)
- **Order submission**: `ExternalExchangeOrderService.Submit()` → Binance POST /order
- **Withdrawal check**: `ExternalExchangeService.CheckWithdrawals()` → Binance withdrawal status
- **Order sync**: `UpdateOrdersFromExternalExchange()` → fetch filled/cancelled orders
- **Order retrieval**: Cache external orders to Redis for local use

### WebSocket (`internal/externalexchangews/` + `internal/processor/`)
- **Connection**: `ExternalExchangeWsService` manages Binance WS lifecycle
- **Streams** (managed by supervisord):
  - `depth` — Order book depth snapshots
  - `trade` — Recent trades
  - `ticker` — 24h price/volume ticker
  - `kline_1m`, `kline_5m`, `kline_1h`, `kline_1d` — Candlestick data
- **Processing**: `WsDataProcessor` transforms raw WS events → Redis live data + MQTT publish
- **Recovery**: On depth stream out-of-sync, re-fetches snapshot from Binance REST API

---

## Configuration

### config/config.yaml Structure
```yaml
exchange:
  environment: "dev"                          # dev | test | production
  active_external_exchange: "binance"         # External exchange provider
  domain: "exchange.local:8000"               # API domain
  supportemail: "support@unitedbit.com"       # Support contact

db:
  name: "mysql"
  dsn: "user:pass@tcp(host:3306)/dbname?parseTime=true"

redis:
  dsn: "redis:6379"
  password: ""
  db: 0

rabbitmq:
  dsn: "amqp://user:pass@host:5672/"
  queue_name: "email_queue_1"

sentry:
  dsn: "https://...@sentry.io/..."
  debug: false

mqtt:
  dsn: "emqtt:1883"
  clientid: "mqtt_client"
  username: "..."
  password: "..."

recaptcha:
  secretkey: "..."
  sitekey: "..."
  androidsecretkey: "..."
  androidsitekey: "..."

jwt:
  private_key: "config/jwt/private.pem"       # RSA private key path
  public_key: "config/jwt/public.pem"         # RSA public key path
  passphrase: "..."                           # Key passphrase
  ttl: 2592000                                # Token TTL in minutes (30 days)

wallet:
  host: "..."                                 # External wallet service URL
  username: "..."
  password: "..."

candle:
  grpcaddr: "candle-app:50051"                # gRPC candle service address
```

### Environment Variable Overrides
- **Prefix**: `UBEXCHANGE_`
- **Pattern**: Dots → underscores. Example: `UBEXCHANGE_DB_DSN` overrides `db.dsn`
- **Loading**: `config/config.go` → `SetConfigs()` returns configured Viper instance
- **Detection**: `flag.Lookup("test.v")` detects test environment → loads `config_test.yaml` overlay
- **Production path**: `/app/config/config.yaml` (inside Docker)
- **Test path**: `./../config/` (relative from test directory)

### RSA Keys
- `config/jwt/private.pem` / `public.pem` — JWT signing/verification
- `config/ub-captcha/private.pem` / `public.pem` — Captcha encryption

---

## Testing

### Infrastructure (`test/main_test.go`)
- **Database**: MariaDB 10.5, separate `exchange_go_test` database
- **Schema**: `test/data/db.sql` (60+ tables, loaded via `setupDb()`)
- **Redis**: Shared Redis instance
- **DI container**: Singleton across all tests via `getContainer()`

### Seeders (`test/data/seed/`)
| Seeder | Seeds | Key Data |
|--------|-------|----------|
| `currencySeeder` | `currencies` | USDT (ID=1), BTC (2), ETH (3), GRS (4) with limits/fees |
| `externalExchangeSeed` | `external_exchanges` | Binance (ID=1) |
| `userLevelSeed` | `user_levels` | VIP1 (100 limit, 1% fee), VIP2 |
| `userPermissionSeed` | `user_permissions` | Deposit, Withdraw, Exchange, FiatDeposit, FiatWithdraw |
| `roleSeed` | `roles` | ROLE_SUPER_ADMIN (1), ROLE_ADMIN (2) |

### Test Helpers
- `getContainer()` — Singleton DI container
- `getDb()` — Singleton GORM connection
- `getRedis()` — Singleton Redis client
- `getUserActor()` — Pre-authenticated test user (email: `test@test.com`, password: `123456789`, 2FA code: `HWOAQZBGXCKJZQVH`)
- `getAdminUserActor()` — Admin user (RoleID=1)
- `getNewUserActor()` — Creates new user with random email/UUID
- `makeRequest(method, path, body)` — Raw HTTP request helper
- `makeAuthorizedRequest(method, path, body, token)` — Authenticated HTTP request

### Test Files (28 files)
```
test/
├── main_test.go                     # Infrastructure setup, seeders, helpers
├── auto_exchange_test.go            # Auto-exchange feature
├── check_withdrawals_in_external_exchange_command_test.go
├── configuration_test.go            # App settings
├── country_test.go                  # Country endpoints
├── currency_test.go                 # Currency/pair endpoints
├── generate_address_command_test.go
├── generate_kline_sync_command_test.go
├── initialize_balance_command_test.go
├── mqtt_auth_test.go                # MQTT broker auth
├── order_cancel_test.go             # Order cancellation
├── order_create_test.go             # Order creation
├── order_list_test.go               # Order listing
├── order_matching_test.go           # Engine matching
├── payment_test.go                  # Payment/withdrawal
├── retrieve_external_orders_command_test.go
├── retrieve_open_orders_command_test.go
├── set_user_level_command_test.go
├── stop_order_submission_test.go    # Stop order triggers
├── submit_bot_aggregated_order_command_test.go
├── sync_kline_command_test.go
├── unmatched_orders_test.go         # Crash recovery
├── update_orders_in_external_exchange_command_test.go
├── update_user_wallet_balances_command_test.go
├── userbalance_test.go              # Balance operations
├── userwithdrawaddess_test.go       # Withdraw addresses
├── user_auth_test.go                # Login/register
└── user_test.go                     # User profile
```

### Running Tests
```bash
# Inside Docker with MariaDB + Redis services
go test ./... --failfast

# CI pipeline uses golang:1.22 image with MariaDB 10.5 + Redis Alpine
```

---

## Security

### JWT Authentication
- RSA-based JWT (RS256) via `golang-jwt/jwt/v5`
- Private/public key pair in `config/jwt/`
- Token TTL: 30 days (production), 10 minutes (test)
- Middleware extracts user from `Authorization: Bearer <token>` header

### MQTT ACL
- EMQX v4 broker with HTTP webhook auth backend
- Endpoints: `/api/v1/emqtt/login`, `/acl`, `/superuser`
- Per-user private channels via `<privateChannel>` identifier
- No MQTT message encryption (plaintext QoS 0)

### Password Security
- bcrypt hashing via `golang.org/x/crypto`
- `platform.PasswordEncoder` interface

### 2FA
- Google Authenticator TOTP via `pquerna/otp`
- SMS OTP via Redis-stored codes + RabbitMQ → ub-communicator

### Known Security Concerns
- ⚠️ `config/config.yaml` contains hardcoded credentials (DB, RabbitMQ, MQTT, JWT passphrase `123456789`, Sentry DSN)
- ⚠️ RSA keys committed to `config/jwt/` and `config/ub-captcha/` — should use secrets manager
- ⚠️ No rate limiting on any API endpoints (login brute-force, order spam, SMS flood)
- ⚠️ CORS allows all origins in non-production (`*`)
- ⚠️ JWT TTL is 30 days — unusually long for a financial application
- ⚠️ No CSRF protection (relies on CORS + Bearer token)
- ⚠️ Admin API on separate port (:8001) but no IP whitelisting
- 🔴 MQTT login endpoint always returns success — no actual authentication (handler/mqtt.go)
- 🔴 Test email backdoor bypasses recaptcha in production (auth/service.go:516)
- 🔴 No JWT algorithm enforcement — RS256/HS256 confusion attack possible (platform/jwt.go:77)
- 🔴 Admin auth missing 2FA change check (auth/service.go:747)
- 🔴 Balance race condition — TOCTOU in postmatch_balance.go, no negative balance check
- 🔴 Binance API keys stored plaintext in database
- 🔴 Payment webhook callback has no signature/HMAC validation

---

## Conventions

### Code Style
- **Financial math**: ALWAYS use `shopspring/decimal` — NEVER `float64` for money
- **Error handling**: Wrap errors with context: `fmt.Errorf("service.Method: %w", err)`
- **Logging**: `uber/zap` structured logger + Sentry error tracking
- **Naming**: Go standard (camelCase locals, PascalCase exports)
- **File naming**: Lowercase, no underscores in package names; underscores in filenames for multi-word

### Architecture Patterns
- **DI**: `sarulabs/di` container with string constants in `internal/di/`
- **Repository pattern**: Each entity has its own repository; data access only through repositories
- **Service layer**: Business logic in service structs; handlers are thin (bind → call service → respond)
- **CLI commands**: All implement `ConsoleCommand` interface: `Run(ctx context.Context, flags []string)`
- **API responses**: `{ "status": bool, "message": string, "data": interface{} }`

### Config Access
- Via `platform.Configs` interface (wraps Viper)
- Env override prefix: `UBEXCHANGE_` (e.g., `UBEXCHANGE_DB_DSN`)

### Redis Key Patterns
- Order books: `order-book:bid:<pair>` / `order-book:ask:<pair>` (sorted sets)
- Stop orders: `queue:stop:order:<type>:<pair>` (sorted sets, score=stop-price)
- Engine queue: `engine:queue:orders` (list)
- Live data: `live_data:pair_currency:<pair>` (hash with fields: price, volume, kline_*, depth_snapshot, etc.)
- Wallet auth: `wallet:auth` (hash, TTL 5h)
- Confirmations: `withdraw-confirmation:<userId>`, `forgot-password:<userId>`, `phone-confirmation:<userId>` (hash, TTL 3h)
- Cache: `<entity>:<id>` pattern (via go-redis/cache)

### MQTT Topic Patterns
- Public: `main/trade/<data-type>/<pair>`
- Private: `main/trade/user/<privateChannel>/<data-type>/`

### Testing Pattern
- Integration tests in `test/` directory (not alongside source)
- Require Docker services (MariaDB 10.5 + Redis)
- Use helpers: `getContainer()`, `getUserActor()`, `makeAuthorizedRequest()`
- Mock interfaces defined in `internal/mocks/` (82 files)

---

## CI/CD

### GitLab CI (`.gitlab-ci.yml`)

| Stage | Job | Trigger | Image | Action |
|-------|-----|---------|-------|--------|
| test | `test` | MR only | `golang:1.22` + MariaDB 10.5 + Redis | `go test ./... --failfast` |
| dev_deploy | `dev_deploy` | Branch: `develop` | — | SSH → `git checkout develop` + `docker-compose restart` |
| deploy | `deploy` | Branch: `master` | — | SSH → `git checkout master` + `docker-compose restart` |

**Deploy targets:**
- Dev: `exchange@135.181.30.255:7268` — restarts 3 containers
- Prod: `exchange@116.203.76.196:7268` — restarts 4 containers (+ nginx)

### Linting (`.golangci.yml`)
- Single linter: **revive** with `var-naming` rule only
- All other revive rules explicitly disabled

---

## Build Commands

```bash
go build -o cli    cmd/exchange-cli/main.go
go build -o httpd  cmd/exchange-httpd/main.go
go build -o engine cmd/exchange-engine/main.go
go build -o ws     cmd/exchange-ws/main.go
```

## Docker

- Go services run inside Docker (Dockerfile at `ub-server-main/.docker/go/Dockerfile`)
- docker-compose services: `exchange-go`, `exchange-httpd-go`, `exchange-engine-go`
- Supervisor for WS streams: `.docker/supervisor/go-supervisord.conf`
  - `depth-stream`: `./ws depth` (3s sleep start)
  - `ticker-trade-stream`: `./ws ticker trade`
  - `kline-stream`: `./ws kline_1m kline_5m kline_1h kline_1d`
- Dependencies: MariaDB 10.5+, Redis 6.2+, RabbitMQ 3.7+, EMQX v4

---

## Known Issues & Upgrade Roadmap

### Active Issues
1. **DI error suppression**: All 109 `builder.Add()` calls use `mustAdd()` wrapper — but 318 `ctn.Get()` type assertions are unchecked (panic on wrong type)
2. **No graceful shutdown**: `exchange-httpd` uses `panic()` on fatal, no SIGTERM handling
3. **Variable shadowing**: Local vars shadow DI constants throughout di files
4. **Magic strings**: HTTP headers repeated as string literals across handlers
5. **Response inconsistency**: Multiple response struct patterns (being unified)
6. **Filename typos**: `ordercreatemanger.go`, `stopordersubmissionmanger.go`, `NewUnmachtechOrdersHandler`

### Upgrade Priority
1. 🟡 Gin 1.10+ already in go.mod ✓ (was 1.7)
2. 🟡 golang-jwt/jwt/v5 already in go.mod ✓ (migrated from dgrijalva)
3. 🟡 rabbitmq/amqp091-go already in go.mod ✓ (migrated from streadway)
4. 🟡 sentry-go 0.30.0 already in go.mod ✓ (was 0.11)
5. 🟠 go-redis v8 → v9 (breaking changes, better perf)
6. 🟠 golang.org/x/* packages — keep current for security patches
7. 🟠 GORM 1.21 → latest v2 (performance + features)
8. 🟢 CI pipeline: already using Go 1.22 ✓
9. 🟢 Docker base image update to match Go 1.22

### Refactoring Tasks (13 documented in `tasks/`)
| Task | Priority | Description |
|------|----------|-------------|
| 001 | HIGH | Unify response envelope (3 different response structs → 1) |
| 002 | MEDIUM | Handler boilerplate wrapper (reduce ~150 lines) |
| 003 | MEDIUM | Fix typos in filenames and identifiers |
| 004 | MEDIUM | Handle DI builder errors (add `mustAdd()` helper) |
| 005 | LOW | Extract magic strings to constants |
| 006 | MEDIUM | Expand validation error handling (4 tags → 15+) |
| 007 | MEDIUM | Split route registration by domain |
| 008 | HIGH | Add graceful shutdown (SIGINT/SIGTERM) |
| 009 | LOW | Fix DI variable shadowing |
| 010 | LOW | Split 407-line sync_kline command |
| 011 | LOW | Add test file documentation |
| 012 | MEDIUM | Create ARCHITECTURE.md ✓ (completed) |
| 013 | MEDIUM | Document DI registration order |

---

## Deep Audit Findings (Line-by-Line Code Review)

> Comprehensive audit performed by reading every source file. Full details in
> `deep-exchange-report.md`.

### DI Container Correction

**Actual service count: 109** (not ~121). Breakdown by file:
| File | `mustAdd()` Calls |
|------|------------------|
| `di_infrastructure.go` | 13 |
| `di_repositories.go` | 28 (27 repos + 1 Redis manager) |
| `di_services.go` | 32 |
| `di_order_services.go` | 21 |
| `di_commands.go` | 15 |
| `di_http.go` | 1 |
| **Total** | **109** |

**318 unchecked type assertions** across all DI files — `ctn.Get(x).(Type)` without comma-ok pattern.

### Engine Bugs Found

| ID | File:Area | Severity | Description |
|----|-----------|----------|-------------|
| E-1 | `redisorderbookprovider.go:46` | **CRITICAL** | Bid limit order price range inverted — sets `min=price` instead of `max=price`, matches only MORE expensive asks |
| E-2 | `queue.go:91` | **HIGH** | `exists()` returns `pos > 0` but `LPos` is 0-based — order at head falsely reported as "not exists" |
| E-3 | `orderbook.go:121` | **MEDIUM** | Market order remainder divides by original marketPrice, not actual trade price |
| E-4 | `orderbook.go:243-286` | **MEDIUM** | Same-price sort uses array index `i < j` as primary — unstable price-time priority |
| E-5 | `redisorderbookprovider.go:183-193` | **LOW** | `Exists()` checks `score > 0` — fails for price=0 edge case (should check `err == nil`) |
| E-6 | `callbackmanager.go:10-14` | **HIGH** | Single global mutex serializes ALL 10 worker callbacks — bottleneck |
| E-7 | `worker.go:24-38` | **HIGH** | No `recover()` in worker goroutine — single panic permanently kills a worker |
| E-8 | `redisorderbookprovider.go:73-104` | **HIGH** | `RewriteOrderBook` early return on marshal error leaves TxPipeline unexecuted → orphaned orders |
| E-9 | `engine.go:159-188` | **MEDIUM** | `HandleInQueueOrders` fire-and-forget goroutine — errors swallowed |

### Security Vulnerabilities Found

| ID | Location | Severity | Description |
|----|----------|----------|-------------|
| S-1 | `handler/mqtt.go:10-16` | **CRITICAL** | `MqttLogin()` always returns success — no actual authentication |
| S-2 | `auth/service.go:516` | **CRITICAL** | Test email `behkamegit@gmail.com` bypasses recaptcha in production |
| S-3 | `postmatch_balance.go:15-75` | **CRITICAL** | TOCTOU race in balance updates — no pessimistic lock, no negative check |
| S-4 | `platform/jwt.go:77` | **HIGH** | No explicit RS256 algorithm enforcement — HS256 confusion attack possible |
| S-5 | `auth/service.go:747-763` | **HIGH** | Admin auth skips `TwoFaChangedAt` check — old tokens survive 2FA toggle |
| S-6 | `config/config.yaml` | **HIGH** | 7 hardcoded secrets (JWT passphrase `123456789`, DB, MQTT, RabbitMQ, Wallet credentials) |
| S-7 | `adminhandler/payment.go` | **HIGH** | Payment webhook callback has no HMAC/signature validation |
| S-8 | All endpoints | **HIGH** | No rate limiting middleware (login brute-force, SMS flood, withdrawal spam) |
| S-9 | `handler/userprofileimage.go` | **MEDIUM** | KYC image upload has no file size validation in handler |
| S-10 | `auth/service.go:737 vs 759` | **MEDIUM** | Inconsistent time comparison: user auth uses `.Sub() > 0`, admin uses `.After()` |

### Concurrency Issues Found

| ID | Location | Severity | Description |
|----|----------|----------|-------------|
| C-1 | `auth/service.go:145-156` | **HIGH** | 4+ goroutines on login — no WaitGroup, no context, no error handling |
| C-2 | `payment/service.go` | **HIGH** | 8+ notification goroutines — fire-and-forget, no tracking |
| C-3 | `processor/dataprocessor.go` | **MEDIUM** | Package-level `sync.Mutex` — potential deadlock under load |
| C-4 | `engine/engine.go:78-80` | **MEDIUM** | `Stop()` sends on quit channel — double `Stop()` call deadlocks |
| C-5 | `pool.go:32` | **MEDIUM** | 1000-item buffered channel — blocks dispatcher if workers stall |
| C-6 | All goroutines | **MEDIUM** | No `context.WithTimeout()` on external API calls (wallet, Binance) |

### Binance Integration Issues

| ID | Area | Severity | Description |
|----|------|----------|-------------|
| B-1 | All REST API calls | **CRITICAL** | Zero retry logic — network timeout = immediate failure, no backoff |
| B-2 | `externalexchangews/binance/ws.go` | **CRITICAL** | `os.Exit(1)` on ANY WebSocket error — no graceful shutdown or reconnection |
| B-3 | `ratelimithandler.go` | **HIGH** | Rate limit map updated via goroutine without mutex — race condition |
| B-4 | `handler/binance/service.go` | **HIGH** | API keys stored plaintext in database, no encryption at rest |
| B-5 | `ws_health_check.go` | **MEDIUM** | 20-second gap before restart detected; no message buffering during reconnect |
| B-6 | `submit_bot_aggregated_order.go` | **HIGH** | Deletes Redis aggregation data regardless of Binance submission success |

### CLI Command Issues

| Command | Issue | Severity |
|---------|-------|----------|
| `submit-bot-orders` | Deletes Redis data even on Binance failure — lost bot orders | **HIGH** |
| `ws-health-check` | No auth on `supervisorctl` call; 20s detection gap | **MEDIUM** |
| `check-withdrawals` | No retry on Binance API failure | **MEDIUM** |
| `sync-kline` | 407-line monolith, hard to maintain | **LOW** |

### Binance REST API Endpoints (Complete Reference)

| Method | URI | Purpose | Weight |
|--------|-----|---------|--------|
| POST | `api.binance.com/api/v3/order` | Place order | 1 |
| POST | `api.binance.com/api/v3/order/test` | Test order (test env) | 1 |
| GET | `api.binance.com/api/v3/klines` | OHLC candles | 1 |
| GET | `sapi/v1/capital/withdraw/history` | Withdrawal status | 1 |
| GET | `api.binance.com/api/v3/allOrders` | All orders | 10 |
| GET | `api.binance.com/api/v3/myTrades` | Trade history | 10 |
| GET | `api.binance.com/api/v3/exchangeInfo` | Symbol metadata | 10 |
| POST | `sapi/v1/capital/withdraw/apply` | Initiate withdrawal | N/A |

### MQTT Trade Throttling
`ProcessTrade()` in `processor/dataprocessor.go` publishes only **1 in every 10 Binance trades** to reduce MQTT load. Counter is mutex-protected per-pair.

---

## Wave 4: Matching Engine Deep Audit (30 bugs found)

### CRITICAL — Price-Time Priority Broken

| # | Bug | File:Line | Impact |
|---|-----|-----------|--------|
| C1 | Float64 precision loss in Redis sorted set scores | redisorderbookprovider.go:90-94 | Orders at same price may execute out of FIFO sequence |
| C2 | Bid price sorting compares INDICES not PRICES (`return i < j`) | orderbook.go:281 | Highest-price bids NOT matched first — PTP violation |
| C3 | Ask orders never sorted after fetch | orderbook.go:262-263 | Lowest-price asks NOT matched first — PTP violation |
| C4 | Global variables (orderbookProvider, cbm, etc.) without synchronization | engine.go:12-15 | Data race per Go memory model |
| C5 | Zero/negative order amounts accepted — no validation | order.go:36-42 | Division by zero, negative-quantity trades |
| C6 | Market orders persist as partial fills instead of cancel | orderbook.go:173-212 | Market orders violate IOC semantics |
| C7 | Unbuffered channel deadlock in worker.stop() | worker.go:71-75 | Goroutine leak on shutdown |
| C8 | Unbuffered channel deadlock in engine.Stop() | engine.go:77-82 | Engine shutdown hangs |

### HIGH — Concurrency, Validation, Error Handling

| # | Bug | File:Line | Impact |
|---|-----|-----------|--------|
| H1 | No self-trade detection (wash trading possible) | orderbook.go:83-171 | Illegal in most jurisdictions |
| H2 | Queue position check off-by-one (`pos > 0` should be `>= 0`) | queue.go:92-96 | First order in queue invisible |
| H3 | ZScore comparison broken for zero prices | redisorderbookprovider.go:195-199 | Incorrect existence checks |
| H4 | Market order infinite loop on zero quantity | orderbook.go:203 | Worker thread permanently blocked |
| H5 | HandleInQueueOrders spawns goroutine that races with workers | engine.go:159-188 | Concurrent order book corruption |
| H6 | Worker continues processing after error | worker.go:41-68 | Corrupted state from failed matching |
| H7 | dispatchOrder infinite retry without backoff | engine.go:102-130 | 100% CPU during Redis outage |
| H8 | StringFixed(16) truncates partial fill precision | orderbook.go:124-127 | Micro amounts lost |
| H9 | MarshalForOrderbook errors silently swallowed | redisorderbookprovider.go:90-92 | Partial orders fail to persist |
| H10 | Pool buffer overflow — hardcoded 1000, no backpressure | pool.go:31-37 | System freeze under load |
| H11 | ParseInt errors silently ignored in sort | orderbook.go:276-277 | Sort corruption on invalid IDs |
| H12 | No timeout on any context (all use context.Background()) | Multiple files | Goroutines hang indefinitely |

### MEDIUM — Edge Cases

| # | Bug | File:Line | Impact |
|---|-----|-----------|--------|
| M1 | No minPrice <= maxPrice validation | orderbook.go:178 | Unfillable orders accepted |
| M2 | bestOrder mutates asks/bids arrays directly | orderbook.go:216-230 | Thread-unsafe |
| M3 | String price equality comparison | orderbook.go:270 | `"9" != "09"` but same price |
| M4 | Market bid value calc uses wrong price | orderbook.go:94-104 | Wrong quantity on market buys |
| M5 | BLPOP timeout hardcoded 1s | engine.go:116 | No perf tuning |
| M6 | loadOrders error continues silently | orderbook.go:252-260 | Silent failures |
| M7 | Error log says wrong method name | orderbook.go:254 | Misleading logs |
| M8 | RemoveOrder continues after queue error | engine.go:137-157 | Orphaned orders |
| M9 | rewriteOrderBook errors not propagated | orderbook.go:233-241 | Silent persistence failure |
| M10 | Integer overflow on large order IDs | orderbook.go:276 | Sort corruption |

### Architectural Recommendation

The matching engine has **fundamental Price-Time Priority violations** (C2+C3) — the most basic rule of fair exchange matching. Combined with float64 precision loss (C1), no input validation (C5), and unsafe concurrency (C4+H5), the engine requires a **significant refactor** before production use. Every trade executed could be at a suboptimal price.
