import 'package:flutter/material.dart';
import 'UBText.dart';
import '../../../generated/assets.gen.dart';
import '../../../generated/colors.gen.dart';
import '../../../utils/mixins/commonConsts.dart';

class UBWarningRow extends StatelessWidget {
  final String text;
  final Color background;
  final Color separatorColor;
  final double height;
  final EdgeInsets margin;
  final EdgeInsets padding;
  const UBWarningRow({
    Key key,
    @required this.text,
    this.background = Colors.transparent,
    this.height = 70.0,
    this.margin = const EdgeInsets.symmetric(
      horizontal: 12.0,
    ),
    this.separatorColor = ColorName.black2c,
    this.padding = const EdgeInsets.symmetric(
      horizontal: 0.0,
    ),
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(color: background, borderRadius: rounded_big),
      height: height,
      margin: margin,
      padding: padding,
      child: Row(
        children: [
          Assets.images.warningInCircle.svg(),
          Container(
            height: 50.0,
            width: 2.0,
            margin: const EdgeInsets.symmetric(
              horizontal: 12.0,
            ),
            color: separatorColor,
          ),
          UBText(
            wrapped: true,
            text: text,
            size: 13,
            color: ColorName.greybf,
          )
        ],
      ),
    );
  }
}
