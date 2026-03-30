# Centrifugo Server API Reference

## HTTP API

Base URL: `http://centrifugo:8000/api`

### Authentication

All API requests require the API key header:
```
X-API-Key: your-api-key
```

### Publish

Send a message to a channel:

```http
POST /api/publish HTTP/1.1
X-API-Key: your-api-key
Content-Type: application/json

{
  "channel": "trade:ticker:BTC-USDT",
  "data": {
    "price": "50000.00",
    "volume": "1.5",
    "change": "+2.5%"
  }
}
```

Response: `{"result": {}}`

### Broadcast

Send to multiple channels at once:

```http
POST /api/broadcast HTTP/1.1
X-API-Key: your-api-key
Content-Type: application/json

{
  "channels": ["trade:ticker:BTC-USDT", "trade:ticker:ETH-USDT"],
  "data": {"type": "market_update"}
}
```

### Presence

Get list of connected clients on a channel:

```http
POST /api/presence HTTP/1.1
X-API-Key: your-api-key
Content-Type: application/json

{"channel": "trade:ticker:BTC-USDT"}
```

### History

Get recent messages in a channel:

```http
POST /api/history HTTP/1.1
X-API-Key: your-api-key
Content-Type: application/json

{"channel": "trade:ticker:BTC-USDT", "limit": 10}
```

### Disconnect

Force disconnect a user:

```http
POST /api/disconnect HTTP/1.1
X-API-Key: your-api-key
Content-Type: application/json

{"user": "user123"}
```

### Unsubscribe

Force unsubscribe a user from a channel:

```http
POST /api/unsubscribe HTTP/1.1
X-API-Key: your-api-key
Content-Type: application/json

{"channel": "user:abc:open-orders", "user": "user123"}
```

## PHP Client (phpcent)

```bash
composer require centrifugal/phpcent:~6.0
```

```php
use phpcent\Client;

$client = new Client("http://centrifugo:8000/api");
$client->setApiKey("your-api-key");

// Publish
$client->publish("trade:ticker:BTC-USDT", [
    "price" => "50000.00",
    "volume" => "1.5"
]);

// Broadcast
$client->broadcast(
    ["trade:ticker:BTC-USDT", "trade:ticker:ETH-USDT"],
    ["type" => "market_update"]
);

// Disconnect user
$client->disconnect("user123");

// Generate connection JWT token
// Use firebase/php-jwt or similar library
use Firebase\JWT\JWT;

$payload = [
    'sub' => (string)$userId,
    'exp' => time() + 3600,
];
$token = JWT::encode($payload, $hmacSecret, 'HS256');
```

## Go HTTP API Client

```go
package centrifugo

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
)

type CentrifugoClient struct {
    baseURL string
    apiKey  string
    client  *http.Client
}

func NewClient(baseURL, apiKey string) *CentrifugoClient {
    return &CentrifugoClient{
        baseURL: baseURL,
        apiKey:  apiKey,
        client:  &http.Client{},
    }
}

func (c *CentrifugoClient) Publish(channel string, data interface{}) error {
    payload := map[string]interface{}{
        "channel": channel,
        "data":    data,
    }
    return c.doRequest("/api/publish", payload)
}

func (c *CentrifugoClient) Broadcast(channels []string, data interface{}) error {
    payload := map[string]interface{}{
        "channels": channels,
        "data":     data,
    }
    return c.doRequest("/api/broadcast", payload)
}

func (c *CentrifugoClient) Disconnect(user string) error {
    payload := map[string]interface{}{
        "user": user,
    }
    return c.doRequest("/api/disconnect", payload)
}

func (c *CentrifugoClient) doRequest(path string, payload interface{}) error {
    body, err := json.Marshal(payload)
    if err != nil {
        return err
    }
    req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewReader(body))
    if err != nil {
        return err
    }
    req.Header.Set("X-API-Key", c.apiKey)
    req.Header.Set("Content-Type", "application/json")
    resp, err := c.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("centrifugo API error: status %d", resp.StatusCode)
    }
    return nil
}
```

## JWT Token Generation

### PHP
```php
use Firebase\JWT\JWT;

function generateCentrifugoToken(int $userId, string $secret): string {
    $payload = [
        'sub' => (string)$userId,
        'exp' => time() + 3600,
    ];
    return JWT::encode($payload, $secret, 'HS256');
}
```

### Go
```go
import "github.com/golang-jwt/jwt/v5"

func GenerateCentrifugoToken(userID string, secret string) (string, error) {
    claims := jwt.MapClaims{
        "sub": userID,
        "exp": time.Now().Add(time.Hour).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}
```

### Subscription Token (for private channels)
```go
func GenerateSubscriptionToken(userID, channel, secret string) (string, error) {
    claims := jwt.MapClaims{
        "sub":     userID,
        "channel": channel,
        "exp":     time.Now().Add(time.Hour).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}
```
