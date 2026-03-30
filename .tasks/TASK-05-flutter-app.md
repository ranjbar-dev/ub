# TASK-05: Replace MQTT with Centrifugo in Flutter App (ub-app-main)

## Objective
Replace the mqtt_client Dart library with centrifuge-dart for WebSocket real-time data in the Flutter mobile/web app.

## IMPORTANT: Pre-Null-Safety Constraint
⚠️ ub-app-main uses Dart SDK <3.0 (pre-null-safety). The centrifuge-dart package requires Dart 2.12+ for null safety. You MUST find a compatible version or use the package with `// @dart=2.9` annotations. Check pub.dev for the latest version that works with pre-null-safety code.

## Overview of Changes
1. Replace `mqtt_client` pub dependency with `centrifuge`
2. Replace entire `lib/mqttClient/` directory with Centrifugo client wrapper
3. Replace `authorizedMqttController.dart` with Centrifugo controller
4. Update topic constants in `lib/services/constants.dart`
5. Update all files importing/using MQTT services

## Files to Modify/Create

### 1. UPDATE: `pubspec.yaml`
- **Current**: `mqtt_client: ^9.6.1`
- **Action**: Replace with:
  ```yaml
  centrifuge: ^0.14.0  # Check latest compatible version with pre-null-safety
  ```
- Run `flutter pub get`

### 2. REPLACE: `lib/mqttClient/` directory (ENTIRE DIRECTORY)
- **Current files**:
  - `lib/mqttClient/universal_mqtt_client.dart` (barrel export)
  - `lib/mqttClient/src/universal_mqtt_client.dart` (main client class)
  - `lib/mqttClient/src/mqtt_vm.dart` (DartVM MQTT implementation)
  - `lib/mqttClient/src/mqtt_browser.dart` (Browser MQTT implementation)
  - `lib/mqttClient/src/mqtt_shared.dart` (shared enums/types)
- **Action**: Replace entire directory with `lib/centrifugoClient/`:
  ```dart
  // lib/centrifugoClient/centrifugo_service.dart
  import 'package:centrifuge/centrifuge.dart' as centrifuge;
  import 'dart:convert';

  class CentrifugoService {
    centrifuge.Client _client;
    final Map<String, centrifuge.Subscription> _subscriptions = {};
    final String _url;
    String _token;

    CentrifugoService({String url, String token})
        : _url = url,
          _token = token;

    Future<void> connect() async {
      _client = centrifuge.createClient(
        _url,
        centrifuge.ClientConfig(token: _token),
      );
      await _client.connect();
    }

    centrifuge.Subscription subscribe(String channel) {
      final sub = _client.newSubscription(channel);
      _subscriptions[channel] = sub;
      sub.subscribe();
      return sub;
    }

    void unsubscribe(String channel) {
      final sub = _subscriptions[channel];
      if (sub != null) {
        sub.unsubscribe();
        _client.removeSubscription(sub);
        _subscriptions.remove(channel);
      }
    }

    Future<void> disconnect() async {
      for (final sub in _subscriptions.values) {
        sub.unsubscribe();
      }
      _subscriptions.clear();
      await _client.disconnect();
    }
  }
  ```

### 3. REPLACE: `lib/app/global/controller/authorizedMqttController.dart`
- **Current**: GetX controller for authenticated MQTT subscriptions
  ```dart
  String get orderTopic => "main/trade/user/$channel/open-orders/";
  String get paymentsTopic => "main/trade/user/$channel/crypto-payments/";
  ```
  - Uses JWT as MQTT username
  - Subscribes to user-specific MQTT topics
  - Uses mqtt_client library directly
- **Action**: Rename to `authorizedCentrifugoController.dart`:
  ```dart
  import 'package:centrifuge/centrifuge.dart' as centrifuge;
  import 'package:get/get.dart';

  class AuthorizedCentrifugoController extends GetxController {
    centrifuge.Client _client;
    centrifuge.Subscription _openOrdersSub;
    centrifuge.Subscription _paymentsSub;
    String channel;

    String get openOrdersChannel => "user:$channel:open-orders";
    String get paymentsChannel => "user:$channel:crypto-payments";

    @override
    void onInit() async {
      super.onInit();
      final token = await _fetchCentrifugoToken();
      _client = centrifuge.createClient(
        centrifugoWsUrl,
        centrifuge.ClientConfig(token: token),
      );

      _client.connecting.listen((_) => print('Centrifugo connecting...'));
      _client.connected.listen((_) => print('Centrifugo connected'));
      _client.disconnected.listen((_) => print('Centrifugo disconnected'));

      await _client.connect();
      _subscribeToUserChannels();
    }

    void _subscribeToUserChannels() {
      _openOrdersSub = _client.newSubscription(openOrdersChannel);
      _openOrdersSub.publication.listen((event) {
        final data = json.decode(utf8.decode(event.data));
        _handleOpenOrderUpdate(data);
      });
      _openOrdersSub.subscribe();

      _paymentsSub = _client.newSubscription(paymentsChannel);
      _paymentsSub.publication.listen((event) {
        final data = json.decode(utf8.decode(event.data));
        _handlePaymentUpdate(data);
      });
      _paymentsSub.subscribe();
    }

    Future<String> _fetchCentrifugoToken() async {
      // Call backend API to get Centrifugo connection token
      final response = await dio.get('/api/v1/auth/centrifugo-token');
      return response.data['token'];
    }

    void disconnectFromTopics() {
      _openOrdersSub?.unsubscribe();
      _paymentsSub?.unsubscribe();
      _client?.disconnect();
    }

    @override
    void onClose() {
      disconnectFromTopics();
      super.onClose();
    }
  }
  ```

### 4. UPDATE: `lib/services/constants.dart`
- **Current**:
  ```dart
  static const priceTopic = 'main/trade/ticker';
  static const orderbookTopic = 'main/trade/order-book/';
  static const ohlcTopic = 'main/trade/kline/';
  ```
- **Action**: Replace with:
  ```dart
  static const centrifugoWsUrl = 'wss://app.unitedbit.com/connection/websocket';
  static const centrifugoLocalWsUrl = 'ws://localhost:8800/connection/websocket';
  static const tickerChannel = 'trade:ticker:';
  static const orderbookChannel = 'trade:order-book:';
  static const ohlcChannel = 'trade:kline:';
  static const tradeBookChannel = 'trade:trade-book:';
  ```

### 5. Search for Additional MQTT References
```bash
grep -r "mqtt" --include="*.dart" ub-app-main/lib/
grep -r "MqttClient" --include="*.dart" ub-app-main/lib/
grep -r "main/trade" --include="*.dart" ub-app-main/lib/
grep -r "mqttClient" --include="*.dart" ub-app-main/lib/
grep -r "UniversalMqttClient" --include="*.dart" ub-app-main/lib/
```
- Update ALL imports referencing old MQTT client
- Look for controllers, services, and views that subscribe to MQTT data
- Check `lib/app/modules/` for any module-specific MQTT usage (e.g., trade module, dashboard)
- Check GetX bindings that inject the MQTT controller

### 6. Update GetX Bindings
- Find where `AuthorizedMqttController` is bound (likely in a global binding or app binding)
- Update to bind `AuthorizedCentrifugoController` instead
- Update any `Get.find<AuthorizedMqttController>()` calls

### 7. Public Data Subscriptions
- Find where public MQTT topics (ticker, orderbook, kline) are subscribed
- These may be in separate controllers or services (not just the authorizedMqttController)
- Replace MQTT subscription with Centrifugo subscription pattern

## Important Notes
- MQTT uses QoS levels; Centrifugo uses best-effort delivery (no QoS)
- MQTT messages come as raw bytes; centrifuge-dart uses Protobuf by default
- Connection URL format: `ws://host:port/connection/websocket` (not MQTT-style)
- centrifuge-dart handles reconnection automatically
- The Flutter app supports both VM (mobile) and browser (web) platforms — centrifuge-dart supports both natively

## Reference Docs
- See `.docs/centrifugo-client-dart.md` for Dart SDK usage and migration guide
- See `.docs/centrifugo-overview.md` for channel naming convention
