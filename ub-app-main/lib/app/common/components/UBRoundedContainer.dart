import 'package:flutter/material.dart';
import '../../../generated/colors.gen.dart';

class UBRoundedContainer extends StatelessWidget {
  final Widget child;
  final double width;
  final Color color;
  const UBRoundedContainer({
    Key key,
    this.child,
    this.width,
    this.color,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      width: width,
      decoration: BoxDecoration(
        color: color ?? ColorName.black2c,
        borderRadius: BorderRadius.circular(
          12,
        ),
      ),
      margin: const EdgeInsets.symmetric(
        horizontal: 12,
        vertical: 24,
      ),
      padding: const EdgeInsets.symmetric(
        horizontal: 12,
        vertical: 24,
      ),
      child: child,
    );
  }
}
