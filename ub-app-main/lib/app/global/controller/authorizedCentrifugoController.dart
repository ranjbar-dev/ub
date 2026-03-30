import 'dart:async';
import 'dart:convert';

import 'package:dio/dio.dart';
import 'package:get/get.dart';
import 'package:get_storage/get_storage.dart';
import '../authorized_order_event_model.dart';
import '../../../centrifugoClient/centrifugo_service.dart';
import '../../../services/constants.dart';
import '../../../services/storageKeys.dart';
import '../../../utils/commonUtils.dart';
import '../../../utils/logger.dart';
import '../../../utils/mixins/toast.dart';

class AuthorizedCentrifugoController extends GetxController with Toaster {
  final GetStorage storage = GetStorage();

  var lastId;
  var lastStatus;

  CentrifugoService centrifugoService;

  StreamSubscription<String> authorizedOrderSubscription;

  StreamSubscription<String> authorizedPaymentSubscription;

  GetStream<List<RxUpdateables>> updateDataSubject = GetStream();

  final ordrPayload = AuthorizedOrderEventModel().obs;

  String get channel => storage.read(StorageKeys.channel);
  CentrifugoService get authedClient => centrifugoService;

  String get openOrdersChannel {
    return "user:$channel:open-orders";
  }

  String get paymentsChannel {
    return "user:$channel:crypto-payments";
  }

  void dsps() {
    if (authorizedOrderSubscription != null) {
      authorizedOrderSubscription.cancel();
    }
    if (authorizedPaymentSubscription != null) {
      authorizedPaymentSubscription.cancel();
    }
  }

  Future disconnectFromTopics() async {
    await Future.wait([
      purgeChannel(
        client: centrifugoService,
        topicStream: authorizedOrderSubscription,
        channel: openOrdersChannel,
      ),
      purgeChannel(
        client: centrifugoService,
        topicStream: authorizedPaymentSubscription,
        channel: paymentsChannel,
      )
    ]);
    centrifugoService.disconnect();
    return Future.value(true);
  }

  @override
  void onClose() async {
    super.onClose();
    purgeChannel(
      client: centrifugoService,
      topicStream: authorizedOrderSubscription,
      channel: openOrdersChannel,
    );
    purgeChannel(
      client: centrifugoService,
      topicStream: authorizedPaymentSubscription,
      channel: paymentsChannel,
    );
    centrifugoService.disconnect();

    storage.remove(StorageKeys.channel);
  }

  @override
  void onInit() async {
    super.onInit();

    centrifugoService = CentrifugoService();

    final token = await _fetchCentrifugoToken();
    centrifugoService.init(
      url: Constants.centrifugoWsUrl,
      token: token,
    );

    centrifugoService.status.listen(
      (status) {
        log.d('Authorized Centrifugo Status: $status');
      },
    );

    try {
      await centrifugoService.connect();
    } catch (e) {
      print(
        e.toString(),
      );
    }

    // Subscribe to user channels after connecting.
    // Centrifuge-dart handles automatic re-subscription on reconnection.
    _subscribeToUserChannels();
  }

  void _subscribeToUserChannels() {
    //connect to order stream for authed users
    authorizedOrderSubscription = centrifugoService
        .handleString(openOrdersChannel)
        .listen((message) {
      log.i('open Orders Centrifugo Initialized');
      if (message != null) {
        log.i('new order event');
        Map<String, dynamic> decoded;
        try {
          decoded = jsonDecode(message);
        } catch (e) {
          log.e('authorizedCentrifugo: failed to decode order message: $e');
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
      log.e('error listening to $openOrdersChannel');
      log.e(e);
    });
    //connect to payment stream
    authorizedPaymentSubscription = centrifugoService
        .handleString(paymentsChannel)
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
      log.e('error listening to $paymentsChannel');
      log.e(e.toString());
    });
  }

  /// Fetch a Centrifugo connection JWT from the backend.
  /// Falls back to the existing REST API token if the endpoint is not available.
  Future<String> _fetchCentrifugoToken() async {
    try {
      final dio = Dio(BaseOptions(
        baseUrl: Constants.baseUrl,
        headers: {
          'Authorization': 'Bearer ${storage.read(StorageKeys.token)}',
        },
      ));
      final response = await dio.get(
        Constants.generatemainUrl('auth/centrifugo-token'),
      );
      if (response.data != null && response.data['token'] != null) {
        return response.data['token'];
      }
    } catch (e) {
      log.d('Centrifugo token endpoint not available, using REST API token');
    }
    // Fallback: use existing REST API JWT token
    return storage.read(StorageKeys.token);
  }
}
