import 'package:dotted_border/dotted_border.dart';
import 'package:flutter/material.dart';
import '../../../generated/colors.gen.dart';

class UBDottedBorder extends StatelessWidget {
  final Widget child;
  final Color color;
  final double borderRadius;
  const UBDottedBorder({
    Key key,
    @required this.child,
    this.color = ColorName.primaryBlue,
    this.borderRadius = 12.0,
  }) : super(key: key);
  @override
  Widget build(BuildContext context) {
    return DottedBorder(
      color: color,
      dashPattern: [2, 2],
      strokeWidth: 1,
      strokeCap: StrokeCap.round,
      borderType: BorderType.RRect,
      radius: Radius.circular(borderRadius),
      child: child,
    );
  }
}
