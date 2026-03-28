import 'package:logger/logger.dart';
import 'environment/ubEnv.dart';

class UBLogger {
  static final UBLogger _singleton = UBLogger._internal();
  factory UBLogger() {
    return _singleton;
  }
  UBLogger._internal();

  static final log = Logger(
    filter: CustomFilter(),
    printer: PrettyPrinter(
        methodCount: 2,
        errorMethodCount: 8,
        stackTraceBeginIndex: 0,
        lineLength: 100,
        colors: true,
        printEmojis: true,
        printTime: false),
  );
}

class CustomFilter extends LogFilter {
  @override
  bool shouldLog(LogEvent event) {
    return ENV == "DEV";
  }
}

final log = UBLogger.log;
