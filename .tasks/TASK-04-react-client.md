# TASK-04: Replace MQTT with Centrifugo in React Client (ub-client-cabinet-main)

## Objective
Replace the MQTT.js client library with centrifuge-js for WebSocket real-time data in the React client cabinet.

## Overview of Changes
1. Replace `mqtt` npm package with `centrifuge`
2. Replace `MqttService2.ts` (public data) with Centrifugo service
3. Replace `RegisteredMqttService.ts` (authenticated user data) with Centrifugo service
4. Update topic constants to Centrifugo channel names
5. Update React hooks for connection management
6. Update message handling services

## Files to Modify/Create

### 1. Package Changes
```bash
cd ub-client-cabinet-main
npm uninstall mqtt
npm install centrifuge
```
- **Current**: `"mqtt": "^4.3.8"` in `package.json`
- **Action**: Replace with `"centrifuge": "^5.0.0"` (or latest)

### 2. REPLACE: `app/services/constants.ts`
- **Current MQTT config**:
  ```typescript
  export const MqttProtocol = 'wss';
  export const mqttServer = `${MqttProtocol}://${mainUrl}:8443`;
  export const MqttAdditionalConfig: IClientOptions = {
      protocol: MqttProtocol,
      connectTimeout: 30 * 1000,
      reconnectPeriod: 2 * 1000,
      keepalive: 60,
  };
  ```
- **Action**: Replace with Centrifugo config:
  ```typescript
  export const centrifugoUrl = `wss://${mainUrl}/connection/websocket`;
  // For local dev: export const centrifugoUrl = `ws://localhost:8800/connection/websocket`;
  ```
- Also remove `IClientOptions` import from mqtt and any MQTT-specific type imports

### 3. REPLACE: `app/services/MqttService2.ts` (Public data service)
- **Current**: Singleton MQTT service subscribing to public topics:
  - `main/trade/ticker` → `MarketWatchMessageService`
  - `main/trade/order-book/` → `OrderBookMessageService`
  - `main/trade/trade-book/` → `MarketTradeMessageService`
  - `main/trade/kline/` → `TradeChartMessageService`
- **Action**: Create `app/services/CentrifugoPublicService.ts`:
  ```typescript
  import { Centrifuge } from 'centrifuge';

  class CentrifugoPublicService {
      private static instance: CentrifugoPublicService;
      private centrifuge: Centrifuge;
      private subscriptions: Map<string, any> = new Map();

      private constructor() {
          // Public data doesn't need auth token (anonymous connection)
          this.centrifuge = new Centrifuge(centrifugoUrl);
          this.centrifuge.connect();
      }

      static getInstance(): CentrifugoPublicService { ... }

      subscribeToTicker(pair: string, callback: (data: any) => void) {
          const channel = `trade:ticker:${pair}`;
          const sub = this.centrifuge.newSubscription(channel);
          sub.on('publication', (ctx) => callback(ctx.data));
          sub.subscribe();
          this.subscriptions.set(channel, sub);
      }

      subscribeToOrderBook(pair: string, callback: (data: any) => void) {
          const channel = `trade:order-book:${pair}`;
          // ... same pattern
      }

      subscribeToTradeBook(pair: string, callback: (data: any) => void) {
          const channel = `trade:trade-book:${pair}`;
          // ... same pattern
      }

      subscribeToKline(pair: string, timeFrame: string, callback: (data: any) => void) {
          const channel = `trade:kline:${timeFrame}:${pair}`;
          // ... same pattern
      }

      unsubscribe(channel: string) {
          const sub = this.subscriptions.get(channel);
          if (sub) {
              sub.unsubscribe();
              this.centrifuge.removeSubscription(sub);
              this.subscriptions.delete(channel);
          }
      }

      disconnect() {
          this.centrifuge.disconnect();
      }
  }
  ```

### 4. REPLACE: `app/services/RegisteredMqttService.ts` (Authenticated user service)
- **Current**: Singleton MQTT service using JWT as username for auth:
  ```typescript
  RegisteredMqttService.mqttCl = mqttConnect(mqttServer, {
      password: RegisteredMqttService.clId,
      username: updatedToken ?? cookies.get(CookieKeys.Token),
      clientId: RegisteredMqttService.clId,
      ...MqttAdditionalConfig,
  });
  ```
  Subscribes to:
  - `main/trade/user/{channel}/open-orders/`
  - `main/trade/user/{channel}/crypto-payments/`
- **Action**: Create `app/services/CentrifugoAuthService.ts`:
  ```typescript
  import { Centrifuge } from 'centrifuge';

  class CentrifugoAuthService {
      private static instance: CentrifugoAuthService;
      private centrifuge: Centrifuge;
      private subscriptions: Map<string, any> = new Map();

      private constructor(token: string) {
          this.centrifuge = new Centrifuge(centrifugoUrl, {
              getToken: async () => {
                  // Fetch Centrifugo connection token from backend
                  const response = await ApiService.get('/api/v1/auth/centrifugo-token');
                  return response.data.token;
              }
          });
          this.centrifuge.connect();
      }

      static getInstance(token?: string): CentrifugoAuthService { ... }

      subscribeToOpenOrders(privateChannel: string, callback: (data: any) => void) {
          const channel = `user:${privateChannel}:open-orders`;
          const sub = this.centrifuge.newSubscription(channel, {
              getToken: async () => {
                  const resp = await ApiService.get(`/api/v1/auth/centrifugo-subscribe-token?channel=${channel}`);
                  return resp.data.token;
              }
          });
          sub.on('publication', (ctx) => callback(ctx.data));
          sub.subscribe();
          this.subscriptions.set(channel, sub);
      }

      subscribeToCryptoPayments(privateChannel: string, callback: (data: any) => void) {
          const channel = `user:${privateChannel}:crypto-payments`;
          // ... same pattern with subscription token
      }

      disconnect() {
          this.subscriptions.forEach((sub) => sub.unsubscribe());
          this.centrifuge.disconnect();
      }
  }
  ```

### 5. UPDATE: `app/containers/App/constants.ts`
- **Current**:
  ```typescript
  enum MqttTopicsPrefixes {
      MarketTradeAddress = 'main/trade/trade-book/',
      OrderBookAddress = 'main/trade/order-book/',
      MarketWatchAddress = 'main/trade/ticker',
      TradeChartAddress = 'main/trade/kline/',
  }
  ```
- **Action**: Replace with:
  ```typescript
  enum CentrifugoChannels {
      MarketTradePrefix = 'trade:trade-book:',
      OrderBookPrefix = 'trade:order-book:',
      TickerPrefix = 'trade:ticker:',
      TradeChartPrefix = 'trade:kline:',
  }
  ```

### 6. UPDATE: `app/containers/App/hooks/connectToMqtt2.tsx`
- **Current hook**: `useConnectToAuthorizedMqtt2` — manages MQTT connection lifecycle
  ```typescript
  mqtt2.current.ConnectToSubject({
      subject: `main/trade/user/${channel}/open-orders/`,
  });
  ```
- **Action**: Rewrite hook to use CentrifugoAuthService instead:
  ```typescript
  export function useConnectToCentrifugo() {
      const centrifugoRef = useRef<CentrifugoAuthService | null>(null);
      // Connect on mount, disconnect on unmount
      // Subscribe to user channels using CentrifugoAuthService
  }
  ```

### 7. UPDATE: Message Service Files
- **Current**: Services like `MarketWatchMessageService`, `OrderBookMessageService`, etc. receive raw MQTT message buffers and parse them
- **Action**: Update to receive already-parsed JSON from Centrifugo `ctx.data` (no need for `TextDecoder` or `JSON.parse` on raw buffer)
- Files to check:
  - `app/services/RegisteredUserMessageService.ts` (or similar)
  - Any message parsing in saga files that handle MQTT data

### 8. Search for Additional MQTT References
```bash
grep -r "mqtt" --include="*.ts" --include="*.tsx" ub-client-cabinet-main/app/
grep -r "MqttService" --include="*.ts" --include="*.tsx" ub-client-cabinet-main/app/
grep -r "RegisteredMqttService" --include="*.ts" --include="*.tsx" ub-client-cabinet-main/app/
grep -r "main/trade" --include="*.ts" --include="*.tsx" ub-client-cabinet-main/app/
```
- Update ALL imports referencing old MQTT services
- Update ALL components/sagas that use MQTT services

## Important Notes
- MQTT messages come as raw buffers → need `JSON.parse(buffer.toString())` 
- Centrifugo messages come as parsed objects in `ctx.data` → much simpler
- The `mqtt` npm package handles reconnection differently from centrifuge-js
- centrifuge-js has built-in reconnection with exponential backoff
- Topic pattern: MQTT uses `/` separator, Centrifugo uses `:` separator

## Reference Docs
- See `.docs/centrifugo-client-js.md` for JavaScript SDK usage and migration guide
- See `.docs/centrifugo-overview.md` for channel naming convention
