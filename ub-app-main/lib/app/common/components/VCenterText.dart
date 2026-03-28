import 'package:flutter/material.dart';

import 'UBText.dart';

class VCenterText extends StatelessWidget {
  final String text;
  final Color color;
  final double size;
  const VCenterText({
    Key key,
    this.text,
    this.color,
    this.size,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Align(
        alignment: Alignment.center,
        child: UBText(color: color, size: size, text: text));
  }
}
