import 'dart:async';

import 'package:meta/meta.dart';

/// Throttling
/// Have method [throttle]
class Throttling {
  final Duration duration;
  bool canDo = true;

  throttle(Function runningFunc) {
    if (canDo == false) return;

    runningFunc();
    canDo = false;

    Future.delayed(duration).then((value) {
      canDo = true;
    });
  }

  Throttling({@required this.duration});
}
