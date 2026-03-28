import 'package:flutter/material.dart';

class UBToolTip extends StatelessWidget {
  final String message;
  final Widget child;
  final bool preferBelow;

  const UBToolTip({Key key, this.message, this.child, this.preferBelow})
      : super(key: key);
  @override
  Widget build(BuildContext context) {
    return Tooltip(
      message: message,
      child: child,
      preferBelow: preferBelow,
    );
  }
}
