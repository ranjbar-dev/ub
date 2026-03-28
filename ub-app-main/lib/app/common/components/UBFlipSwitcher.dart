import 'dart:math';

import 'package:flutter/material.dart';

class UBFlipSwitcher extends StatelessWidget {
  final Duration duration;
  final Widget child1;
  final Widget child2;
  final bool conditionToShowChild1;

  const UBFlipSwitcher({
    Key key,
    this.duration = const Duration(milliseconds: 300),
    @required this.child1,
    @required this.child2,
    @required this.conditionToShowChild1,
  }) : super(key: key);
  @override
  Widget build(BuildContext context) {
    final child = conditionToShowChild1 ? child1 : child2;
    return AnimatedSwitcher(
      child: child,
      //duration: Duration(seconds: 1),
      duration: duration,
      transitionBuilder: (Widget child, Animation<double> animation) {
        final rotate = Tween(begin: pi, end: 0.0).animate(animation);
        return AnimatedBuilder(
            animation: rotate,
            child: child,
            builder: (BuildContext context, Widget child) {
              final angle = min(rotate.value, pi / 2);
              return Transform(
                transform: Matrix4.rotationX(angle),
                child: child,
                alignment: Alignment.center,
              );
            });
      },
      switchInCurve: Curves.easeIn,
      switchOutCurve: Curves.easeOut,
    );
  }
}
