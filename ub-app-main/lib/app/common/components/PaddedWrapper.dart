import 'package:flutter/material.dart';

class UBPaddedWrapper extends StatelessWidget {
  final Widget child;
  final double padding;
  const UBPaddedWrapper({Key key, this.child, this.padding}) : super(key: key);
  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: EdgeInsets.symmetric(horizontal: padding ?? 32),
      child: child,
    );
  }
}
