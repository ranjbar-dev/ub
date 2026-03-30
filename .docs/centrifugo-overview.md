# Centrifugo Overview & Migration Reference

## What is Centrifugo?

Centrifugo is a scalable, language-agnostic, real-time messaging server. It delivers messages over WebSocket, HTTP-streaming, SSE, WebTransport, and gRPC. It's a user-facing pub/sub server where business logic stays in your backend and Centrifugo manages real-time connections.

- **GitHub**: https://github.com/centrifugal/centrifugo
- **Docs**: https://centrifugal.dev/
- **Docker Hub**: https://hub.docker.com/r/centrifugo/centrifugo/

## Key Differences from EMQX/MQTT

| Aspect | EMQX (MQTT) | Centrifugo (WebSocket) |
|--------|-------------|----------------------|
| Protocol | MQTT over TCP/WebSocket | Native WebSocket, SSE, HTTP-streaming |
| Topics | Hierarchical with wildcards (`+`, `#`) | Flat channel names (string-based) |
| QoS | 0, 1, 2 (at-most-once to exactly-once) | Best-effort delivery |
| Retained Messages | Yes | Channel history (configurable) |
| Auth | Username/password, TLS certs, JWT, plugins | JWT tokens or proxy auth to backend |
| Session | Stateful, persistent sessions | Stateless per WebSocket connection |
| Publishing | Any MQTT client can publish | Backend publishes via HTTP/gRPC API |
| Scaling | Clustering built-in | Redis/NATS/Tarantool for horizontal scaling |

## Architecture

```
┌──────────────────────────┐
│  Clients (Browser/Mobile) │
│  centrifuge-js / dart SDK │
└────────────┬─────────────┘
             │ WebSocket (ws:// or wss://)
             ▼
┌──────────────────────────┐
│      Centrifugo Server    │
│  - JWT auth validation    │
│  - Channel pub/sub        │
│  - Presence & history     │
│  - Admin web UI           │
│  Port 8000                │
└─────┬──────────┬─────────┘
      │          │
      │ HTTP API │ Proxy callbacks
      │ (publish)│ (connect/subscribe)
      ▼          ▼
┌──────────────────────────┐
│      Your Backend         │
│  - Generate JWT tokens    │
│  - Publish via HTTP API   │
│  - Handle proxy auth      │
└──────────────────────────┘
```

## Channel Naming (Topic Mapping)

EMQX MQTT topics → Centrifugo channels:

| MQTT Topic | Centrifugo Channel |
|---|---|
| `main/trade/ticker/{pair}` | `trade:ticker:{pair}` |
| `main/trade/order-book/{pair}` | `trade:order-book:{pair}` |
| `main/trade/trade-book/{pair}` | `trade:trade-book:{pair}` |
| `main/trade/chart/{timeFrame}/{pair}` | `trade:chart:{timeFrame}:{pair}` |
| `main/trade/kline/{timeFrame}/{pair}` | `trade:kline:{timeFrame}:{pair}` |
| `main/trade/change-price/{pair}` | `trade:change-price:{pair}` |
| `main/trade/market-price/{pair}` | `trade:market-price:{pair}` |
| `main/trade/user/{privateChannel}/open-orders/` | `user:{privateChannel}:open-orders` |
| `main/trade/user/{privateChannel}/crypto-payments/` | `user:{privateChannel}:crypto-payments` |

**Namespace convention**: Use `:` as separator. Configure namespace `trade` and `user` in Centrifugo config.

## Authentication Flow

### Current EMQX Flow:
1. Client connects to EMQX via MQTT over WebSocket (port 8443)
2. EMQX calls `/api/v1/emqtt/login` on backend to authenticate
3. EMQX calls `/api/v1/emqtt/acl` for per-topic authorization
4. Both PHP and Go backends serve these auth endpoints

### New Centrifugo Flow:
1. Client requests a connection JWT from backend (`POST /api/v1/auth/centrifugo-token`)
2. Backend generates JWT with `sub` (user ID) and `exp` claims, signed with HMAC secret
3. Client connects to Centrifugo WebSocket with JWT token
4. Centrifugo validates JWT automatically
5. For private channels: use subscription JWT or proxy subscribe endpoint
6. EMQX auth endpoints (`/api/v1/emqtt/*`) can be **removed**

## Publishing Messages

### Current (MQTT):
- PHP: `EmqttManager` publishes via MQTT protocol to EMQX broker
- Go: `mqttmanager` publishes via MQTT protocol to EMQX broker

### New (Centrifugo HTTP API):
- PHP: Use `centrifugal/phpcent` library to publish via HTTP API
- Go: Use HTTP POST to `http://centrifugo:8000/api/publish` with API key

```php
// PHP - using phpcent
$client = new \phpcent\Client("http://centrifugo:8000/api");
$client->setApiKey("your-api-key");
$client->publish("trade:ticker:BTC-USDT", ["price" => "50000", "volume" => "1.5"]);
```

```go
// Go - HTTP API publish
func publishToCentrifugo(channel string, data interface{}) error {
    payload := map[string]interface{}{
        "channel": channel,
        "data":    data,
    }
    body, _ := json.Marshal(payload)
    req, _ := http.NewRequest("POST", "http://centrifugo:8000/api/publish", bytes.NewReader(body))
    req.Header.Set("X-API-Key", apiKey)
    req.Header.Set("Content-Type", "application/json")
    resp, err := http.DefaultClient.Do(req)
    // handle response...
}
```

## Client SDKs

| Platform | Package | Install |
|----------|---------|---------|
| JavaScript (React) | `centrifuge` | `npm install centrifuge` |
| Dart (Flutter) | `centrifuge` | `flutter pub add centrifuge` |
| Go | `centrifuge-go` | `go get github.com/centrifugal/centrifuge-go` |
| PHP (server-side) | `phpcent` | `composer require centrifugal/phpcent:~6.0` |
