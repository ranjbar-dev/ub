import 'dart:async';
import 'dart:convert';

import 'package:get/get.dart';
import 'package:get_storage/get_storage.dart';
import '../authorized_order_event_model.dart';
import '../../../mqttClient/universal_mqtt_client.dart';
import '../../../services/constants.dart';
import '../../../services/storageKeys.dart';
import '../../../utils/commonUtils.dart';
import '../../../utils/logger.dart';
import '../../../utils/mixins/toast.dart';
import 'package:uuid/uuid.dart';

class AuthorizedMqttController extends GetxController with Toaster {
  final GetStorage storage = GetStorage();

  var lastId;
  var lastStatus;

  UniversalMqttClient authorizedClient;

  StreamSubscription<String> authorizedOrderSubscription;

  StreamSubscription<String> authorizedPaymentSubscription;

  GetStream<List<RxUpdateables>> updateDataSubject = GetStream();

  final ordrPayload = AuthorizedOrderEventModel().obs;

  Timer _timer;

  String get channel => storage.read(StorageKeys.channel);
  UniversalMqttClient get authedClient => authorizedClient;

  String get orderTopic {
    return "main/trade/user/$channel/open-orders/";
  }

  String get paymentsTopic {
    return "main/trade/user/$channel/crypto-payments/";
  }

  void dsps() {
    authorizedOrderSubscription.cancel();
    authorizedPaymentSubscription.cancel();
  }

  Future disconnectFromTopics() async {
    await Future.wait([
      purgeTopic(
        client: authorizedClient,
        topicStream: authorizedOrderSubscription,
        topic: orderTopic,
      ),
      purgeTopic(
        client: authorizedClient,
        topicStream: authorizedPaymentSubscription,
        topic: paymentsTopic,
      )
    ]);
    authorizedClient.disconnect();
    return Future.value(true);
  }

  @override
  void onClose() async {
    super.onClose();
    purgeTopic(
      client: authorizedClient,
      topicStream: authorizedOrderSubscription,
      topic: orderTopic,
    );
    purgeTopic(
      client: authorizedClient,
      topicStream: authorizedPaymentSubscription,
      topic: paymentsTopic,
    );
    _timer.cancel();
    _timer = null;
    authorizedClient.disconnect();
    // await disconnectFromTopics();

    storage.remove(StorageKeys.channel);
  }

  @override
  void onInit() async {
    super.onInit();

    // SECURITY RISK: The JWT access token is used as the MQTT username.
    // If MQTT traffic is not encrypted end-to-end (TLS), or broker logs
    // are exposed, the token can be stolen and replayed against the REST API.
    // TODO: Replace with a short-lived MQTT-only credential obtained from
    //       a dedicated server endpoint (e.g. POST /auth/mqtt-token) that
    //       issues credentials scoped only to the MQTT broker.
    authorizedClient = UniversalMqttClient(
      broker: Uri.parse(Constants.mqttServer),
      autoReconnect: true,
      timeout: const Duration(seconds: 10),
      username: storage.read(StorageKeys.token),
      password: Uuid().v4(),
    );

    authorizedClient.status.listen(
      (status) {
        log.i('Authorized Connection Status: $status');
        if (status == UniversalMqttClientStatus.disconnected) {
          purgeTopic(
            client: authorizedClient,
            topicStream: authorizedOrderSubscription,
            topic: orderTopic,
          );
          purgeTopic(
            client: authorizedClient,
            topicStream: authorizedPaymentSubscription,
            topic: paymentsTopic,
          );
        }
        if (status == UniversalMqttClientStatus.connected) {
          // if (_timer != null) {
          //   _timer.cancel();
          //   _timer = null;
          // }
          // // keep the connection alive
          // _timer = new Timer.periodic(const Duration(seconds: 5), (t) {
          //   if (authedClient.status.value ==
          //       UniversalMqttClientStatus.connected) {
          //     try {
          //       authedClient.publishString('test', '0', MqttQos.atLeastOnce);
          //     } catch (e) {
          //       log.e(e);
          //     }
          //   }
          // });

          //connect to order stream for authed users
          authorizedOrderSubscription = authorizedClient
              .handleString(orderTopic, MqttQos.exactlyOnce)
              .listen((message) {
            log.i('open Orders Mqtt Initialized');
            if (message != null) {
              log.i('new order event');
              Map<String, dynamic> decoded;
              try {
                decoded = jsonDecode(message);
              } catch (e) {
                log.e('authorizedMqtt: failed to decode order message: $e');
                return;
              }
              final payload =
                  AuthorizedOrderEventModel.fromJson(decoded);
              ordrPayload.value = payload;
              toastAuthorizedEvent(payload);
              final messageStatus = payload.status.toLowerCase();
              log.i("order event, status : ", messageStatus);
              if (messageStatus == 'open' || messageStatus == 'placed') {
                updateDataSubject.add([
                  RxUpdateables.UserPairBalances,
                  RxUpdateables.OpenOrders,
                ]);
              } else if (messageStatus == 'filled' ||
                  messageStatus == 'canceled') {
                updateDataSubject.add([
                  RxUpdateables.UserPairBalances,
                  RxUpdateables.OpenOrders,
                  RxUpdateables.OrderHistory
                ]);
              }

              updateDataSubject.add([RxUpdateables.Balances]);
            }
          }, onError: (e) {
            log.e('error listening to $orderTopic');
            log.e(e);
          });
          //connect to payment stream
          authorizedPaymentSubscription = authorizedClient
              .handleString(paymentsTopic, MqttQos.exactlyOnce)
              .listen((message) {
            if (message != null) {
              log.e('new payment added');
              toastInfo('Balance updated');
              updateDataSubject.add([
                RxUpdateables.Balances,
                RxUpdateables.TransactionHistory,
              ]);
            }
          }, onError: (e) {
            log.e('error listening to $paymentsTopic');
            log.e(e.toString());
          });
        }
      },
    );
    try {
      await authorizedClient.connect();
    } catch (e) {
      print(
        e.toString(),
      );
    }
  }
}
