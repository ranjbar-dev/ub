import 'dart:convert';

import 'package:centrifuge/centrifuge.dart' as centrifuge;
import 'package:rxdart/rxdart.dart';

/// Connection status states for the Centrifugo client.
enum CentrifugoConnectionStatus {
  connecting,
  connected,
  disconnected,
}

/// Wrapper around the centrifuge-dart client providing a simplified API
/// compatible with both VM (mobile) and browser (web) platforms.
///
/// Replaces the old UniversalMqttClient with Centrifugo WebSocket transport.
/// Centrifuge-dart handles automatic reconnection and re-subscription.
class CentrifugoService {
  centrifuge.Client _client;
  final Map<String, centrifuge.Subscription> _subscriptions = {};

  final _status = BehaviorSubject<CentrifugoConnectionStatus>.seeded(
    CentrifugoConnectionStatus.disconnected,
  );

  /// Stream of connection status changes.
  ValueStream<CentrifugoConnectionStatus> get status => _status.stream;

  /// Initialize the Centrifugo client with connection [url] and optional JWT [token].
  /// For anonymous (public-only) access, pass an empty or null token.
  void init({String url, String token}) {
    _client = centrifuge.createClient(
      url,
      centrifuge.ClientConfig(token: token ?? ''),
    );

    _client.connecting.listen((_) {
      _status.add(CentrifugoConnectionStatus.connecting);
    });
    _client.connected.listen((_) {
      _status.add(CentrifugoConnectionStatus.connected);
    });
    _client.disconnected.listen((_) {
      _status.add(CentrifugoConnectionStatus.disconnected);
    });
  }

  /// Connect to the Centrifugo server.
  Future<void> connect() async {
    _status.add(CentrifugoConnectionStatus.connecting);
    await _client.connect();
  }

  /// Subscribe to a [channel] and return a stream of decoded string messages.
  ///
  /// Centrifuge-dart sends data as protobuf-encoded bytes; this method
  /// decodes them to UTF-8 strings automatically.
  /// If already subscribed to the channel, returns a new mapped stream
  /// from the existing subscription.
  Stream<String> handleString(String channel) {
    if (!_subscriptions.containsKey(channel)) {
      final sub = _client.newSubscription(channel);
      _subscriptions[channel] = sub;
      sub.subscribe();
    }
    return _subscriptions[channel].publication.map(
      (event) => utf8.decode(event.data),
    );
  }

  /// Unsubscribe from a [channel] and remove the subscription.
  void unsubscribe({String channel}) {
    final sub = _subscriptions[channel];
    if (sub != null) {
      sub.unsubscribe();
      _client.removeSubscription(sub);
      _subscriptions.remove(channel);
    }
  }

  /// Check if currently subscribed to a [channel].
  bool isSubscribed(String channel) {
    return _subscriptions.containsKey(channel);
  }

  /// Disconnect from the server and clean up all subscriptions.
  void disconnect() {
    for (final sub in _subscriptions.values) {
      sub.unsubscribe();
    }
    _subscriptions.clear();
    if (_client != null) {
      _client.disconnect();
    }
  }

  /// Dispose of the service, disconnecting and closing status stream.
  void dispose() {
    disconnect();
    _status.close();
  }
}
