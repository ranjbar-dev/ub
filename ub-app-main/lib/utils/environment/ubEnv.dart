import 'package:flutter/foundation.dart';

//const ENV = kReleaseMode ? "PRODUCTION" : "DEV";
const ENV = "PRODUCTION";
// String.fromEnvironment('ENV', defaultValue: 'DEV');
const VERSION = String.fromEnvironment('VERSION', defaultValue: '1.0.0');
