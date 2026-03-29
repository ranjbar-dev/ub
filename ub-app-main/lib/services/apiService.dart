import 'dart:async';
import 'dart:io';

import 'package:connectivity/connectivity.dart';
import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:get/get.dart' hide Response, FormData;
import 'package:get_storage/get_storage.dart';

import '../utils/computes.dart';
import '../utils/environment/ubEnv.dart';
import '../utils/logger.dart';
import 'constants.dart';
import 'interceptors/request_retrier.dart';
import 'storageKeys.dart';

class ApiService {
  static GetStorage storage;
  static DioConnectivityRequestRetrier connectionRequestRetrier;
  static Dio dio;
  static Dio rawDio;
  static BaseOptions options;
  static ApiService _apiService;
  static String token;
  static Connectivity connectivity;
  static bool _isRefreshing = false;
  static Completer<String> _refreshCompleter;
  static int _refreshRetryCount = 0;
  static const int _maxRefreshRetries = 3;
  static final isDebug = ENV == "DEV";
  static final FlutterSecureStorage _secureStorage = FlutterSecureStorage();
  static final platform = (GetPlatform.isAndroid ? 'ubandroid' : 'ubandroid') +
      '-v' +
      Constants.appVersion.split("+")[0];

  factory ApiService() {
    if (_apiService == null) {
      _apiService = ApiService._internal();
    }
    return _apiService;
  }

  ApiService._internal() {
    if (storage == null) {
      storage = GetStorage();
    }
    if (connectivity == null) {
      connectivity = Connectivity();
    }
    if (options == null) {
      options = new BaseOptions(
          baseUrl: Constants.baseUrl,
          contentType: 'application/json',
          connectTimeout: 10000,
          receiveTimeout: 10000,
          headers: {if (!(GetPlatform.isWeb)) 'platform': platform});
    }
    if (dio == null) {
      rawDio = new Dio();
      dio = new Dio(options);
      (dio.transformer as DefaultTransformer).jsonDecodeCallback = parseJson;
      if (connectionRequestRetrier == null) {
        connectionRequestRetrier = DioConnectivityRequestRetrier(
          dio: dio,
          connectivity: connectivity,
        );
      }
    }

    dio.interceptors.add(
      InterceptorsWrapper(
        onError: (DioError err, ErrorInterceptorHandler handler) async {
          if (_shouldRetryForInternetLoss(err)) {
            try {
              log.w('retry scheduled!');
              return connectionRequestRetrier.scheduleRequestRetry(err);
            } catch (e) {
              handler.next(e);
            }
          } else {
            return handler.next(err);
          }
        },
      ),
    );

    //handle refresh token
    dio.interceptors.add(
      InterceptorsWrapper(
        onError: (DioError err, ErrorInterceptorHandler handler) async {
          if (_shouldRefreshToken(err)) {
            try {
              String newToken;
              if (_isRefreshing) {
                // Another refresh is in progress, wait for it
                newToken = await _refreshCompleter.future;
              } else {
                _refreshRetryCount++;
                if (_refreshRetryCount > _maxRefreshRetries) {
                  _refreshRetryCount = 0;
                  return handler.next(err);
                }
                _refreshCompleter = Completer<String>();
                _isRefreshing = true;
                try {
                  dio.interceptors.requestLock.lock();
                  dio.interceptors.responseLock.lock();
                  final refreshObj = await post(
                    url: "auth/refresh",
                    data: {"refresh": await _secureStorage.read(key: SecureStorageKeys.refresh)},
                  );
                  newToken = refreshObj["token"];
                  token = newToken;
                  await _secureStorage.write(
                      key: SecureStorageKeys.refresh,
                      value: refreshObj["refreshToken"]);
                  dio.interceptors.requestLock.unlock();
                  dio.interceptors.responseLock.unlock();
                  _refreshRetryCount = 0;
                  _refreshCompleter.complete(newToken);
                } catch (e) {
                  dio.interceptors.requestLock.unlock();
                  dio.interceptors.responseLock.unlock();
                  _refreshCompleter.completeError(e);
                  rethrow;
                } finally {
                  _isRefreshing = false;
                }
              }
              RequestOptions options = err.requestOptions;
              options.headers["Authorization"] = "Bearer " + newToken;
              return dio.fetch(options);
            } catch (e) {
              return handler.next(err);
            }
          } else {
            return handler.next(err);
          }
        },
      ),
    );

    // dio.interceptors.add(
    //   TimeoutRetryInterceptor(
    //     dio: dio,
    //     options: RetryOptions(
    //       retries: 1, // Number of retries before a failure
    //       retryInterval: 5.seconds, // Interval between each retry
    //       retryEvaluator: (error) =>
    //           error.type ==
    //           DioErrorType
    //               .connectTimeout, // Evaluating if a retry is necessary regarding the error. It is a good candidate for updating authentication token in case of a unauthorized error (be careful with concurrency though)
    //     ),
    //   ),
    // );

    dio.interceptors.add(
      InterceptorsWrapper(
        onRequest:
            (RequestOptions options, RequestInterceptorHandler handler) async {
          if (token == null) {
            token = await _secureStorage.read(key: SecureStorageKeys.token);
          }
          if (token != null) {
            options.headers['Authorization'] = 'Bearer' + ' $token';
          } else {
            if (isDebug) {
              log.w(
                'unAuthorized Request to: ${options.baseUrl}${options.path}',
              );
            }
          }

          return handler.next(options); //continue
        },
        onResponse:
            (Response response, ResponseInterceptorHandler handler) async {
          return handler.next(response); // continue
        },
        onError: (DioError err, ErrorInterceptorHandler handler) async {
          if (isDebug) {
            log.e(err.message.toString());
          }
          return handler.next(err); //continue
        },
      ),
    );
    // if (isDebug) {
    //   dio.interceptors.add(PrettyDioLogger(
    //     requestHeader: true,
    //     requestBody: true,
    //     responseBody: true,
    //     error: true,
    //     request: true,
    //     responseHeader: false,
    //     compact: true,
    //   ));
    // }
  }
  bool _shouldRetryForInternetLoss(DioError err) {
    return err.type == DioErrorType.other &&
        err.error != null &&
        err.error is SocketException;
  }

  bool _shouldRefreshToken(DioError err) {
    return err.response != null && err.response.statusCode == 401;
  }

  Future get({
    String url,
    Function urlGenerator,
    dynamic data,
    String rawUrl,
  }) async {
    if (rawUrl != null) {
      final response = await dio.get(
        rawUrl,
        queryParameters: data == null ? {} : data,
      );
      return response.data;
    }
    if (urlGenerator == null) {
      final response = await dio.get(
        Constants.generatemainUrl(url),
        queryParameters: data == null ? {} : data,
      );
      return response.data;
    }
    final response = await dio.get(
      urlGenerator(url),
      queryParameters: data == null ? {} : data,
    );
    return response.data;
  }

  Future rawGet({
    @required String rawUrl,
  }) async {
    final response = await rawDio.get(
      rawUrl,
    );
    return response.data;
  }

  Future post({
    @required String url,
    @required dynamic data,
  }) async {
    final response = await dio.post(
      Constants.generatemainUrl(url),
      data: data,
    );
    return response.data;
  }

  Future upload(
      {FormData form,
      RxInt stream,
      String url,
      CancelToken cancelToken}) async {
    final response = await dio.post(
      Constants.generatemainUrl(url),
      data: form,
      cancelToken: cancelToken,
      onSendProgress: (int sent, int total) {
        final progress = (sent / total) * 100;
        stream.value = progress.toInt();
      },
    );
    return response.data;
  }

  // Future<void> refreshToken() async {
  //   final refreshToken = storage.read(StorageKeys.refresh);
  //   final response =
  //       await post(url: "auth/refresh", data: {'refresh': refreshToken});
  //   if (response.statusCode == 200) {
  //     token = response.data['token'];
  //   }
  // }

  // Future<Response<dynamic>> _retry(RequestOptions requestOptions) async {
  //   final options = new Options(
  //     method: requestOptions.method,
  //     headers: requestOptions.headers,
  //   );
  //   return dio.request<dynamic>(requestOptions.path,
  //       data: requestOptions.data,
  //       queryParameters: requestOptions.queryParameters,
  //       options: options);
  // }
}

final apiService = ApiService();
