# Centrifugo JavaScript Client SDK (centrifuge-js)

## Installation

```bash
npm install centrifuge
```

## Connection with JWT

```javascript
import { Centrifuge } from 'centrifuge';

// Option 1: Direct token
const centrifuge = new Centrifuge('ws://localhost:8800/connection/websocket', {
  token: 'USER-JWT-TOKEN'
});

// Option 2: Token refresh function (recommended)
const centrifuge = new Centrifuge('ws://localhost:8800/connection/websocket', {
  getToken: async () => {
    const response = await fetch('/api/v1/auth/centrifugo-token', {
      headers: { 'Authorization': `Bearer ${jwtToken}` }
    });
    const data = await response.json();
    return data.token;
  }
});

// Connection lifecycle events
centrifuge.on('connecting', (ctx) => console.log('Connecting:', ctx));
centrifuge.on('connected', (ctx) => console.log('Connected:', ctx));
centrifuge.on('disconnected', (ctx) => console.log('Disconnected:', ctx));

centrifuge.connect();
```

## Subscribing to Channels

```javascript
// Public channel (e.g., ticker for all users)
const tickerSub = centrifuge.newSubscription('trade:ticker:BTC-USDT');
tickerSub.on('publication', (ctx) => {
  console.log('Ticker update:', ctx.data);
});
tickerSub.on('subscribing', (ctx) => console.log('Subscribing:', ctx));
tickerSub.on('subscribed', (ctx) => console.log('Subscribed:', ctx));
tickerSub.on('unsubscribed', (ctx) => console.log('Unsubscribed:', ctx));
tickerSub.subscribe();

// Order book
const orderBookSub = centrifuge.newSubscription('trade:order-book:BTC-USDT');
orderBookSub.on('publication', (ctx) => {
  console.log('Order book update:', ctx.data);
});
orderBookSub.subscribe();

// Private channel (user-specific, requires subscription token)
const openOrdersSub = centrifuge.newSubscription('user:abc123:open-orders', {
  getToken: async () => {
    const resp = await fetch('/api/v1/auth/centrifugo-subscribe-token?channel=user:abc123:open-orders');
    const data = await resp.json();
    return data.token;
  }
});
openOrdersSub.on('publication', (ctx) => {
  console.log('Open order update:', ctx.data);
});
openOrdersSub.subscribe();
```

## Unsubscribing & Disconnecting

```javascript
// Unsubscribe from a channel
tickerSub.unsubscribe();
centrifuge.removeSubscription(tickerSub);

// Disconnect entirely
centrifuge.disconnect();
```

## Migration from MQTT.js

### Before (MQTT.js with EMQX):
```javascript
import mqtt from 'mqtt';
const client = mqtt.connect('wss://domain:8443/mqtt', {
  username: 'mqtt_abbas',
  password: 'mqtt_abbas'
});
client.subscribe('main/trade/ticker/BTC-USDT');
client.on('message', (topic, message) => {
  const data = JSON.parse(message.toString());
});
```

### After (centrifuge-js with Centrifugo):
```javascript
import { Centrifuge } from 'centrifuge';
const centrifuge = new Centrifuge('wss://domain/connection/websocket', {
  getToken: async () => { /* fetch JWT from backend */ }
});
centrifuge.connect();
const sub = centrifuge.newSubscription('trade:ticker:BTC-USDT');
sub.on('publication', (ctx) => {
  const data = ctx.data; // already parsed JSON
});
sub.subscribe();
```

## Key Differences from MQTT.js

| MQTT.js | centrifuge-js |
|---------|---------------|
| `client.subscribe(topic)` | `centrifuge.newSubscription(channel).subscribe()` |
| `client.on('message', cb)` | `sub.on('publication', cb)` |
| `client.publish(topic, msg)` | Not available (publish from backend only) |
| Topic: `main/trade/ticker/BTC-USDT` | Channel: `trade:ticker:BTC-USDT` |
| Auth: username/password | Auth: JWT token |
| Message: raw buffer → JSON.parse | Message: already parsed `ctx.data` |
