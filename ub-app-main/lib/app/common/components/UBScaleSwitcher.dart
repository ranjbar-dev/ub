import 'package:flutter/material.dart';

class UBScaleSwitcher extends StatelessWidget {
  final Duration duration;
  final Widget child1;
  final Widget child2;
  final bool conditionToShowChild1;

  const UBScaleSwitcher({
    Key key,
    this.duration = const Duration(milliseconds: 200),
    @required this.child1,
    @required this.child2,
    @required this.conditionToShowChild1,
  }) : super(key: key);
  @override
  Widget build(BuildContext context) {
    return AnimatedSwitcher(
        duration: duration,
        transitionBuilder: (Widget child, Animation<double> animation) {
          return ScaleTransition(child: child, scale: animation);
        },
        child: conditionToShowChild1 ? child1 : child2);
  }
}
