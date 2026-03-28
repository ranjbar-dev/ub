import 'package:flutter/material.dart';

class UBHorizontalDivider extends StatelessWidget {
  final double thickness;
  final Color color;
  final EdgeInsets verticalMargin;
  const UBHorizontalDivider({
    Key key,
    this.thickness = 1.0,
    this.color = const Color(0xFF525261),
    this.verticalMargin = const EdgeInsets.symmetric(vertical: 24.0),
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      height: thickness,
      margin: verticalMargin,
      color: color,
    );
  }
}
