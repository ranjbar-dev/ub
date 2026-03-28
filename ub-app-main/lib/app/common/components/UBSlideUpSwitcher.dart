import 'package:flutter/material.dart';

class UBSlideUpSwitcher extends StatelessWidget {
  final Duration duration;
  final Widget child1;
  final Widget child2;
  final bool conditionToShowChild1;

  const UBSlideUpSwitcher({
    Key key,
    this.duration = const Duration(milliseconds: 200),
    @required this.child1,
    @required this.child2,
    @required this.conditionToShowChild1,
  }) : super(key: key);
  @override
  Widget build(BuildContext context) {
    final child = conditionToShowChild1 ? child1 : child2;
    return AnimatedSwitcher(
      child: child,
      duration: duration,
      transitionBuilder: (Widget child, Animation<double> animation) {
        final slide = Tween(begin: 200.0, end: 0.0).animate(animation);
        return AnimatedBuilder(
            animation: slide,
            child: child,
            builder: (BuildContext context, Widget child) {
              return Transform(
                transform: Matrix4.translationValues(0.0, slide.value, 0.0),
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
