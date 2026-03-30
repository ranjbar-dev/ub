# Centrifugo Dart/Flutter Client SDK (centrifuge-dart)

## Installation

```yaml
# pubspec.yaml
dependencies:
  centrifuge: ^0.14.0
```

```bash
flutter pub get
```

## Connection with JWT

```dart
import 'package:centrifuge/centrifuge.dart' as centrifuge;

// Create client with JWT token in URL or config
final client = centrifuge.createClient(
  'ws://localhost:8800/connection/websocket',
  centrifuge.ClientConfig(
    token: 'USER-JWT-TOKEN',
    // OR use getToken for refresh:
    // getToken: (centrifuge.ConnectionTokenEvent event) async {
    //   final response = await dio.get('/api/v1/auth/centrifugo-token');
    //   return response.data['token'];
    // },
  ),
);

// Connection state events
client.connecting.listen((event) => print('Connecting...'));
client.connected.listen((event) => print('Connected!'));
client.disconnected.listen((event) => print('Disconnected: ${event.reason}'));

await client.connect();
```

## Subscribing to Channels

```dart
// Public channel - ticker
final tickerSub = client.newSubscription('trade:ticker:BTC-USDT');
tickerSub.publication.listen((event) {
  final data = json.decode(utf8.decode(event.data));
  print('Ticker: $data');
});
tickerSub.subscribing.listen((_) => print('Subscribing to ticker...'));
tickerSub.subscribed.listen((_) => print('Subscribed to ticker'));
await tickerSub.subscribe();

// Order book
final orderBookSub = client.newSubscription('trade:order-book:BTC-USDT');
orderBookSub.publication.listen((event) {
  final data = json.decode(utf8.decode(event.data));
  print('Order book: $data');
});
await orderBookSub.subscribe();

// Private channel (user-specific)
final openOrdersSub = client.newSubscription(
  'user:$privateChannelName:open-orders',
  centrifuge.SubscriptionConfig(
    getToken: (centrifuge.SubscriptionTokenEvent event) async {
      final response = await dio.get(
        '/api/v1/auth/centrifugo-subscribe-token',
        queryParameters: {'channel': event.channel},
      );
      return response.data['token'];
    },
  ),
);
openOrdersSub.publication.listen((event) {
  final data = json.decode(utf8.decode(event.data));
  print('Open orders: $data');
});
await openOrdersSub.subscribe();
```

## Unsubscribing & Disconnecting

```dart
// Unsubscribe from a channel
await tickerSub.unsubscribe();
client.removeSubscription(tickerSub);

// Disconnect entirely
await client.disconnect();
```

## Migration from mqtt_client

### Before (MQTT with EMQX):
```dart
import 'package:mqtt_client/mqtt_client.dart';
import 'package:mqtt_client/mqtt_server_client.dart';

final client = MqttServerClient.withPort('wss://domain', 'flutter_client', 8443);
client.websocketProtocols = MqttClientConstants.protocolsSingleDefault;
await client.connect('mqtt_abbas', 'mqtt_abbas');
client.subscribe('main/trade/ticker/BTC-USDT', MqttQos.atMostOnce);
client.updates!.listen((messages) { ... });
```

### After (centrifuge-dart with Centrifugo):
```dart
import 'package:centrifuge/centrifuge.dart' as centrifuge;

final client = centrifuge.createClient(
  'wss://domain/connection/websocket',
  centrifuge.ClientConfig(token: jwtToken),
);
await client.connect();
final sub = client.newSubscription('trade:ticker:BTC-USDT');
sub.publication.listen((event) {
  final data = json.decode(utf8.decode(event.data));
  // handle data
});
await sub.subscribe();
```

## Key Differences from mqtt_client

| mqtt_client (MQTT) | centrifuge-dart (Centrifugo) |
|---|---|
| `client.subscribe(topic, qos)` | `client.newSubscription(channel).subscribe()` |
| `client.updates.listen(cb)` | `sub.publication.listen(cb)` |
| Topics: `main/trade/ticker/BTC-USDT` | Channels: `trade:ticker:BTC-USDT` |
| Auth: username/password on connect | Auth: JWT token in ClientConfig |
| Protocol: MQTT over WebSocket | Protocol: Native WebSocket |
| Port: 8443 (WSS) | Port: 8800 (mapped from 8000) |
| Protobuf: default for centrifuge-dart | Protobuf: default (more efficient) |

## Notes for Pre-Null-Safety (ub-app-main)

The ub-app-main Flutter project uses Dart SDK <3.0 (pre-null-safety). You may need to use an older version of centrifuge-dart that supports this:
- Check `centrifuge: ^0.10.0` or similar version that doesn't require null safety
- Alternatively, pin an older version and add `// @dart=2.9` at top of files
- **Recommended**: This migration is a good opportunity to also migrate to null safety
