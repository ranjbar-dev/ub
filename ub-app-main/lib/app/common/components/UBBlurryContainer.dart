import 'dart:ui';

import 'package:flutter/material.dart';

const double kBlur = 1.0;
const EdgeInsetsGeometry kDefaultPadding = EdgeInsets.all(16);
const Color kDefaultColor = Colors.transparent;
const BorderRadius kBorderRadius = BorderRadius.all(Radius.circular(20));
const double kColorOpacity = 0.0;

class UBBlurryContainer extends StatelessWidget {
  final Widget child;
  final double blur;
  final EdgeInsetsGeometry padding;
  final Color bgColor;

  final BorderRadius borderRadius;

  //final double colorOpacity;

  UBBlurryContainer({
    this.child,
    this.blur = 1,
    this.padding = kDefaultPadding,
    this.bgColor = kDefaultColor,
    this.borderRadius = kBorderRadius,
    //this.colorOpacity = kColorOpacity,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      transform: Matrix4.translationValues(0.0, -1.0, 0.0),
      child: ClipRRect(
        borderRadius: borderRadius,
        child: BackdropFilter(
          filter: ImageFilter.blur(sigmaX: blur, sigmaY: blur),
          child: Container(
            padding: padding,
            color: bgColor == Colors.transparent
                ? bgColor
                : bgColor.withOpacity(0.5),
            child: child,
          ),
        ),
      ),
    );
  }
}
