# TASK-03: Replace EMQX with Centrifugo in Go Backend (ub-exchange-cli-main)

## Objective
Replace all MQTT/EMQX publishing and authentication code in the Go trading engine with Centrifugo HTTP API integration.

## Overview of Changes
1. Replace `internal/platform/mqtt.go` (paho MQTT client) with a Centrifugo HTTP API client
2. Replace `internal/communication/mqttmanager.go` with Centrifugo publisher manager
3. Remove `internal/auth/mqttauthservice.go` (EMQX auth webhooks not needed)
4. Remove `internal/api/handler/mqtt.go` (EMQX webhook HTTP handlers)
5. Update config files
6. Update Go module dependencies
7. Update tests

## Files to Modify/Create

### 1. REPLACE: `internal/platform/mqtt.go`
- **Current**: MQTT client using `github.com/eclipse/paho.mqtt.golang` v1.3.5
  ```go
  type MqttClient interface {
      Publish(topic string, qos byte, retained bool, payload interface{}) (token mqtt.Token)
  }
  func NewMqttClient(configs Configs, logger Logger) MqttClient {
      dsn := configs.GetString("mqtt.dsn")
      clientID := configs.GetString("mqtt.clientid") + uuid
      username := configs.GetString("mqtt.username")
      password := configs.GetString("mqtt.password")
      // ... paho mqtt.Client initialization
  }
  ```
- **Action**: Replace with `internal/platform/centrifugo.go`:
  ```go
  package platform

  import (
      "bytes"
      "encoding/json"
      "fmt"
      "net/http"
  )

  type CentrifugoClient interface {
      Publish(channel string, data interface{}) error
      Broadcast(channels []string, data interface{}) error
  }

  type centrifugoClient struct {
      baseURL string
      apiKey  string
      client  *http.Client
      logger  Logger
  }

  func NewCentrifugoClient(configs Configs, logger Logger) CentrifugoClient {
      return &centrifugoClient{
          baseURL: configs.GetString("centrifugo.api_url"),
          apiKey:  configs.GetString("centrifugo.api_key"),
          client:  &http.Client{},
          logger:  logger,
      }
  }

  func (c *centrifugoClient) Publish(channel string, data interface{}) error {
      payload := map[string]interface{}{"channel": channel, "data": data}
      return c.doRequest("/api/publish", payload)
  }

  func (c *centrifugoClient) Broadcast(channels []string, data interface{}) error {
      payload := map[string]interface{}{"channels": channels, "data": data}
      return c.doRequest("/api/broadcast", payload)
  }

  func (c *centrifugoClient) doRequest(path string, payload interface{}) error {
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
          c.logger.Error("centrifugo publish error", err)
          return err
      }
      defer resp.Body.Close()
      if resp.StatusCode != http.StatusOK {
          return fmt.Errorf("centrifugo API error: status %d", resp.StatusCode)
      }
      return nil
  }
  ```

### 2. REPLACE: `internal/communication/mqttmanager.go`
- **Current MQTT topic constants**:
  ```go
  const (
      TradeTopic     = "main/trade/trade-book/"
      KlineTopic     = "main/trade/kline/"
      TickerTopic    = "main/trade/ticker"
      OrderbookTopic = "main/trade/order-book/"
      UserPrivateTopicPrefix = "main/trade/user/"
      UserOpenOrdersPostfix  = "/open-orders/"
      UserPaymentsPostfix    = "/crypto-payments/"
  )
  ```
- **Current publish methods**:
  ```go
  func (m *manager) PublishTrades(ctx context.Context, pairName string, payload []byte) {
      topic := TradeTopic + pairName
      m.mqttClient.Publish(topic, 0, false, payload)
  }
  ```
- **Action**: Rename to `centrifugomanager.go`. Update constants to Centrifugo channels:
  ```go
  const (
      TradeChannel     = "trade:trade-book:"
      KlineChannel     = "trade:kline:"
      TickerChannel    = "trade:ticker:"
      OrderbookChannel = "trade:order-book:"
      UserChannelPrefix = "user:"
      UserOpenOrdersSuffix  = ":open-orders"
      UserPaymentsSuffix    = ":crypto-payments"
  )
  ```
- Update all publish methods to use CentrifugoClient.Publish() instead of mqttClient.Publish():
  ```go
  func (m *manager) PublishTrades(ctx context.Context, pairName string, payload []byte) error {
      channel := TradeChannel + pairName
      var data interface{}
      json.Unmarshal(payload, &data)
      return m.centrifugoClient.Publish(channel, data)
  }
  ```
- **Note**: MQTT publish takes raw bytes; Centrifugo publish takes a JSON-serializable interface{}. The payload must be unmarshaled before publishing.

### 3. REMOVE: `internal/auth/mqttauthservice.go`
- **Current**: EMQX webhook authentication service with Login, ACL, SuperUser methods
- **Action**: Delete entirely. Centrifugo handles auth via JWT tokens, no webhook callbacks needed.

### 4. REMOVE: `internal/api/handler/mqtt.go`
- **Current**: HTTP handlers for EMQX webhook endpoints
  - `MqttLogin()` → `/api/v1/emqtt/login`
  - `MqttACL()` → `/api/v1/emqtt/acl`
  - `MqttSuperUser()` → `/api/v1/emqtt/superuser`
- **Action**: Delete entirely.

### 5. UPDATE: Route registration
- Find where `/api/v1/emqtt/*` routes are registered (likely in a router setup file)
- Remove the three EMQX webhook routes
- Add a new route for Centrifugo token generation:
  - `POST /api/v1/auth/centrifugo-token` → Generate connection JWT
  - `GET /api/v1/auth/centrifugo-subscribe-token` → Generate subscription JWT for private channels

### 6. CREATE: Centrifugo token endpoint handler
- **Purpose**: Generate JWT tokens for Centrifugo client authentication
- **Connection token**: Contains `sub` (user ID) and `exp` claims
- **Subscription token**: Contains `sub`, `channel`, and `exp` claims
- **Signing**: HMAC-SHA256 with shared secret from config
- Use `github.com/golang-jwt/jwt/v5` (likely already a dependency)

### 7. UPDATE: `config/config.yaml`
- **Current**:
  ```yaml
  mqtt:
    dsn: "emqtt:1883"
    clientid: "mqtt_abbas"
    username: "mqtt_abbas"
    password: "mqtt_abbas"
  ```
- **Action**: Replace with:
  ```yaml
  centrifugo:
    api_url: "http://centrifugo:8000"
    api_key: "your-api-key"
    token_hmac_secret: "your-secret-key"
  ```

### 8. UPDATE: `config/config.docker.yaml`
- Same changes as config.yaml but with Docker service name

### 9. UPDATE: `config/config_test.yaml`
- Same changes for test environment

### 10. UPDATE: `test/mqtt_auth_test.go`
- **Current**: Tests MQTT auth service with topics like `main/trade/order-book/BTC-USDT`
- **Action**: Delete this test file (mqttauthservice.go is being removed) or replace with Centrifugo token generation tests

### 11. UPDATE: DI Container
- Find where `MqttClient` is registered in the DI container (sarulabs/di)
- Replace with `CentrifugoClient` registration
- Find where mqtt manager is registered and update

### 12. UPDATE: `go.mod`
- REMOVE: `github.com/eclipse/paho.mqtt.golang v1.3.5`
- Ensure `github.com/golang-jwt/jwt/v5` is present (for token generation)
- Run `go mod tidy`

## Search for Additional References
- `grep -r "mqtt" --include="*.go" ub-exchange-cli-main/` for any remaining MQTT references
- `grep -r "emqtt" --include="*.go" ub-exchange-cli-main/` for EMQX-specific references
- Check cmd/ entry points for MQTT initialization code
- Check any middleware or interceptors that reference MQTT

## Reference Docs
- See `.docs/centrifugo-server-api.md` for Go client implementation and JWT generation
- See `.docs/centrifugo-overview.md` for channel naming convention
