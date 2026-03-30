# TASK-02: Replace EMQX with Centrifugo in PHP Backend (ub-server-main)

## Objective
Replace all MQTT/EMQX publishing and authentication code in the PHP Symfony backend with Centrifugo HTTP API integration using the `phpcent` library.

## Overview of Changes
1. Install `centrifugal/phpcent` composer package
2. Replace `EmqttClient.php` with a `CentrifugoClient.php` service
3. Replace `EmqttManager.php` with `CentrifugoManager.php` (update all publish methods)
4. Replace `EmqttConstants.php` with `CentrifugoConstants.php` (channel names instead of MQTT topics)
5. Remove or repurpose `EmqttController.php` (EMQX auth webhooks no longer needed)
6. Update config files (parameters, services, security)

## Files to Modify/Create

### 1. Install phpcent library
```bash
cd ub-server-main
composer require centrifugal/phpcent:~6.0
```

### 2. REMOVE: `src/Exchange/CoreBundle/Services/EmqttClient.php`
- **Current**: Low-level MQTT client wrapper using `Bluerhinos/phpMQTT`
- **Action**: Delete this file. Replace with new Centrifugo client.

### 3. CREATE: `src/Exchange/CoreBundle/Services/CentrifugoClient.php`
- **Purpose**: Centrifugo HTTP API client wrapper
- **Implementation**:
  ```php
  <?php
  namespace Exchange\CoreBundle\Services;

  use phpcent\Client;
  use Symfony\Component\DependencyInjection\ParameterBag\ParameterBagInterface;

  class CentrifugoClient
  {
      private Client $client;

      public function __construct(ParameterBagInterface $parameterBag)
      {
          $apiUrl = $parameterBag->get('centrifugo.api_url');
          $apiKey = $parameterBag->get('centrifugo.api_key');
          $this->client = new Client($apiUrl);
          $this->client->setApiKey($apiKey);
      }

      public function publish(string $channel, array $data): void
      {
          $this->client->publish($channel, $data);
      }

      public function broadcast(array $channels, array $data): void
      {
          $this->client->broadcast($channels, $data);
      }
  }
  ```

### 4. REPLACE: `src/Exchange/CommunicationBundle/Constants/EmqttConstants.php`
- **Current contents**:
  ```php
  const EMQTT_PUBLIC_TOPIC_CHART_PREFIX = 'main/trade/chart/';
  const EMQTT_PUBLIC_TOPIC_ORDER_BOOK_PREFIX = 'main/trade/order-book/';
  const EMQTT_PUBLIC_TOPIC_TRADE_BOOK_PREFIX = 'main/trade/trade-book/';
  const EMQTT_PUBLIC_TOPIC_TICKER = 'main/trade/ticker/';
  const EMQTT_PUBLIC_TOPIC_CHANGE_PRICE_PREFIX = 'main/trade/change-price/';
  const EMQTT_PUBLIC_TOPIC_KLINE_PREFIX = 'main/trade/kline/';
  const EMQTT_PRIVATE_TOPIC_USER_PREFIX = 'main/trade/user/';
  const EMQTT_USER_OPEN_ORDERS_POSTFIX = '/open-orders/';
  const EMQTT_USER_CRYPTO_PAYMENTS_POSTFIX = '/crypto-payments/';
  ```
- **Action**: Rename file to `CentrifugoConstants.php`, update namespace and constants:
  ```php
  const CENTRIFUGO_CHANNEL_CHART_PREFIX = 'trade:chart:';
  const CENTRIFUGO_CHANNEL_ORDER_BOOK_PREFIX = 'trade:order-book:';
  const CENTRIFUGO_CHANNEL_TRADE_BOOK_PREFIX = 'trade:trade-book:';
  const CENTRIFUGO_CHANNEL_TICKER = 'trade:ticker:';
  const CENTRIFUGO_CHANNEL_CHANGE_PRICE_PREFIX = 'trade:change-price:';
  const CENTRIFUGO_CHANNEL_KLINE_PREFIX = 'trade:kline:';
  const CENTRIFUGO_CHANNEL_USER_PREFIX = 'user:';
  const CENTRIFUGO_USER_OPEN_ORDERS_POSTFIX = ':open-orders';
  const CENTRIFUGO_USER_CRYPTO_PAYMENTS_POSTFIX = ':crypto-payments';
  ```

### 5. REPLACE: `src/Exchange/CommunicationBundle/Services/EmqttManager.php`
- **Current**: Uses EmqttClient to publish MQTT messages to EMQX broker
- **Key methods to rewrite** (change MQTT publish to Centrifugo HTTP API publish):
  - `publishOrdersToOrderBook($pairCurrencyName, $ordersData)` → publish to `trade:order-book:{pair}`
  - `publishOrderToOpenOrders(Order $order)` → publish to `user:{channel}:open-orders`
  - `publishToOpenOrders($trades)` → publish to `user:{channel}:open-orders`
  - `publishTradesFromExternalExchangeToTradeBook($pairCurrencyName, $trades)` → publish to `trade:trade-book:{pair}`
  - `publishOhlcFromExternalExchange($ohlc, $pairCurrencyName, $timeFrame)` → publish to `trade:kline:{timeFrame}:{pair}`
  - `publishCurrentPrice(...)` → publish to `trade:ticker:{pair}`
  - `publishChangePricePercentage(...)` → publish to `trade:change-price:{pair}`
  - `pushCryptoPaymentStatusToUser(CryptoPayment $cryptoPayment)` → publish to `user:{channel}:crypto-payments`
- **Action**: Rename to `CentrifugoManager.php`. Replace EmqttClient dependency with CentrifugoClient. Replace all topic strings with Centrifugo channel names. Replace `$this->emqttClient->publish()` with `$this->centrifugoClient->publish()`.

### 6. REMOVE/REPURPOSE: `src/Exchange/ApiBundle/Controller/V1/EmqttController.php`
- **Current**: Handles EMQX webhook callbacks for auth/ACL
  - `loginAction()` - `/api/v1/emqtt/login`
  - `aclAction()` - `/api/v1/emqtt/acl`
  - `superuserAction()` - `/api/v1/emqtt/superuser`
- **Action**: Remove this controller entirely. Centrifugo uses JWT-based auth, no webhook callbacks needed.
- **Also**: Remove the route definitions for `/api/v1/emqtt/*`

### 7. UPDATE: `app/config/parameters.docker.yml`
- **Current**:
  ```yaml
  emqtt_publisher_username: mqtt_abbas
  emqtt_publisher_password: mqtt_abbas
  emqtt_publisher_clientid: mqtt_abbas
  ```
- **Action**: Replace with:
  ```yaml
  centrifugo.api_url: "http://centrifugo:8000/api"
  centrifugo.api_key: "your-api-key"
  centrifugo.token_hmac_secret: "your-secret-key"
  ```

### 8. UPDATE: `app/config/parameters.yml.dist`
- **Current**: Contains `emqtt_publisher_*` parameter templates
- **Action**: Replace with `centrifugo.*` parameter templates

### 9. UPDATE: `app/config/security.yml`
- **Current**: Has `/api/v1/emqtt/*` → `PUBLIC_ACCESS` rule
- **Action**: Remove the emqtt security rule. Add a new endpoint for Centrifugo token generation if needed:
  ```yaml
  api_centrifugo_token:
      pattern: ^/api/v1/auth/centrifugo-token
      # Requires authenticated user
  ```

### 10. UPDATE: `app/config/config.yml`
- **Current**: Contains EMQX/MQTT comments and configuration
- **Action**: Replace EMQX references with Centrifugo. Update service definitions.

### 11. UPDATE: Service definitions (DI config)
- Find where `EmqttClient` and `EmqttManager` are registered as services
- Update to reference `CentrifugoClient` and `CentrifugoManager`
- Remove `Bluerhinos/phpMQTT` from composer dependencies

### 12. Create endpoint: Centrifugo connection token
- **New endpoint**: `POST /api/v1/auth/centrifugo-token`
- **Purpose**: Generate JWT for Centrifugo connection (clients call this before connecting)
- **Implementation**: Use the user's ID as `sub` claim, sign with HMAC secret
- **Auth**: Requires valid JWT bearer token (existing auth)

### 13. Create endpoint: Centrifugo subscription token (for private channels)
- **New endpoint**: `GET /api/v1/auth/centrifugo-subscribe-token?channel=...`
- **Purpose**: Generate channel-specific JWT for private channel subscription
- **Implementation**: Verify user owns the channel, generate JWT with `sub` + `channel` claims
- **Auth**: Requires valid JWT bearer token

## Search for Additional References
- Search entire ub-server-main for `emqtt`, `mqtt`, `EmqttClient`, `EmqttManager`, `EmqttConstants`, `Bluerhinos`, `phpMQTT`
- Update any EventSubscribers that reference MQTT/EMQX services
- Update any command classes that use MQTT

## Composer Changes
- ADD: `centrifugal/phpcent:~6.0`
- ADD: `firebase/php-jwt` (if not already present, for generating Centrifugo tokens)
- REMOVE: `bluerhinos/phpmqtt` (or `bluerhinos/phpmqtt:dev-master`)

## Reference Docs
- See `.docs/centrifugo-server-api.md` for PHP client usage and JWT token generation
- See `.docs/centrifugo-overview.md` for channel naming convention
