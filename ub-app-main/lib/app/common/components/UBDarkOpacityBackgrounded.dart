import 'package:flutter/material.dart';

class UBDarkOpacityBackgrounded extends StatelessWidget {
  final double opacity;
  final Widget child;
  const UBDarkOpacityBackgrounded(
      {Key key, this.opacity = 0.8, this.child = const SizedBox()})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      color: Colors.black.withOpacity(opacity),
      child: child,
    );
  }
}
