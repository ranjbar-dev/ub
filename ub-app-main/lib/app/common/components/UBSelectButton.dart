import 'package:flutter/material.dart';
import 'UBGreyContainer.dart';
import '../../../generated/colors.gen.dart';

import 'UBText.dart';

class UBSelectButton extends StatelessWidget {
  final IconData icon;
  final String valueText;
  final Function onClick;
  final double iconSize;
  final EdgeInsets padding;
  final Color backgroundColor;
  const UBSelectButton({
    Key key,
    this.icon = Icons.keyboard_arrow_down,
    @required this.valueText,
    @required this.onClick,
    this.iconSize,
    this.padding,
    this.backgroundColor = ColorName.black2c,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: onClick,
      child: UBGreyContainer(
        padding:
            padding ?? const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
        color: backgroundColor,
        width: double.infinity,
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            UBText(
              text: valueText,
            ),
            Icon(
              icon,
              color: ColorName.greybf,
              size: iconSize,
            )
          ],
        ),
      ),
    );
  }
}
