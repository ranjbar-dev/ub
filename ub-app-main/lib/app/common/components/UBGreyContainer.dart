import 'package:flutter/material.dart';
import '../../../generated/colors.gen.dart';

class UBGreyContainer extends StatelessWidget {
  final Widget child;
  final EdgeInsets margin;
  final EdgeInsets padding;
  final double width;
  final double height;

  final Color color;

  const UBGreyContainer(
      {Key key,
      this.margin,
      this.child,
      this.width,
      this.color = ColorName.grey16,
      this.height = 36.0,
      this.padding})
      : super(key: key);
  @override
  Widget build(BuildContext context) {
    return Container(
      width: width,
      margin: margin,
      padding: padding ??
          const EdgeInsets.symmetric(
            horizontal: 8,
          ),
      height: height,
      alignment: Alignment.centerLeft,
      decoration: BoxDecoration(
        color: color,
        borderRadius: BorderRadius.circular(
          6,
        ),
        border: Border.all(
          color: ColorName.grey19,
          width: 1,
        ),
      ),
      child: child,
    );
  }
}
