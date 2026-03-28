import 'dart:async';

import 'package:connectivity/connectivity.dart';
import 'package:dio/dio.dart' show DioError, Response, Dio;
import 'package:meta/meta.dart';

class DioConnectivityRequestRetrier {
  final Dio dio;
  final Connectivity connectivity;

  DioConnectivityRequestRetrier({
    @required this.dio,
    @required this.connectivity,
  });

  Future<Response> scheduleRequestRetry(DioError err) async {
    StreamSubscription streamSubscription;
    final responseCompleter = Completer<Response>();
    streamSubscription = connectivity.onConnectivityChanged.listen(
      (connectivityResult) async {
        if (connectivityResult != ConnectivityResult.none) {
          streamSubscription.cancel();
          // Complete the completer instead of returning
          responseCompleter.complete(dio.fetch(err.requestOptions));
        }
      },
    );

    return responseCompleter.future;
  }
}
