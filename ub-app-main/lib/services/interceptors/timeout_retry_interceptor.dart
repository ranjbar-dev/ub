import 'package:dio/dio.dart';
// import 'package:logging/logging.dart';
import 'package:meta/meta.dart';

import 'options.dart';

/// An interceptor that will try to send failed request again
class TimeoutRetryInterceptor extends Interceptor {
  final Dio dio;

  final Function(ErrorResult err) errorCallback;

  final RetryOptions options;

  TimeoutRetryInterceptor(
      {@required this.dio, RetryOptions options, this.errorCallback})
      : this.options = options ?? const RetryOptions();

  @override
  onError(DioError err, ErrorInterceptorHandler handler) async {
    var extra = this.options;

//     var shouldRetry = extra.retries > 0 && await extra.retryEvaluator(err); (bugged, as per https://github.com/aloisdeniel/dio_retry/pull/5)
    var shouldRetry = extra.retries > 0 && await options.retryEvaluator(err);
    if (shouldRetry) {
      if (extra.retryInterval.inMilliseconds > 0) {
        await Future<void>.delayed(extra.retryInterval);
      }

      // Update options to decrease retry count before new try
      extra = extra.copyWith(retries: extra.retries - 1);

      try {
        errorCallback(ErrorResult(err.requestOptions.uri.toString(),
            err.message, extra.retries, err.error.toString()));

        // logger?.warning(
        //     "[${err.request.uri}] An error occured during request, trying a again (remaining tries: ${extra.retries}, error: ${err.error})");
        // We retry with the updated options
        await dio.fetch(err.requestOptions);
      } catch (e) {
        return e;
      }
    }
  }
}

class ErrorResult {
  final String uri;
  final String response;
  final int retryCount;
  final String error;

  // final allowedRetries;

  ErrorResult(this.uri, this.response, this.retryCount, this.error);
}
