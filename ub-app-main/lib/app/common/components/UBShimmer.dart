import 'package:flutter/material.dart';
import 'package:shimmer/shimmer.dart';
import '../../../generated/colors.gen.dart';

class UBShimmer extends StatelessWidget {
  final double width;
  final double height;
  final double opacity;
  final Color background;

  const UBShimmer(
      {Key key,
      this.width = 50.0,
      this.height = 14.0,
      this.opacity = 1.0,
      this.background = Colors.black})
      : super(key: key);
  @override
  Widget build(BuildContext context) {
    return Shimmer.fromColors(
      period: const Duration(milliseconds: 500),
      child: Container(
        color: background.withOpacity(opacity),
        height: height,
        width: width,
      ),
      baseColor: background,
      highlightColor: ColorName.grey19,
    );
  }
}
