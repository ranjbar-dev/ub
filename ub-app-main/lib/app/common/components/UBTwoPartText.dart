import 'package:flutter/material.dart';

import '../../../generated/colors.gen.dart';
import 'UBText.dart';

class UBTwoPartText extends StatelessWidget {
  final String title;
  final String value;
  final Widget valueWidget;
  final double firstTextWidth;
  final EdgeInsets padding;
  final double size;
  const UBTwoPartText({
    Key key,
    @required this.title,
    this.value,
    this.firstTextWidth,
    this.padding,
    this.size,
    this.valueWidget,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding:
          padding ?? const EdgeInsets.only(left: 12, right: 12, bottom: 12),
      child: Row(
        children: [
          Container(
            width: firstTextWidth ?? 123.0,
            child: UBText(
              text: "$title : ",
              color: ColorName.grey80,
              size: size,
            ),
          ),
          valueWidget ??
              UBText(
                text: value,
                color: ColorName.white,
                size: size,
              ),
        ],
      ),
    );
  }
}
