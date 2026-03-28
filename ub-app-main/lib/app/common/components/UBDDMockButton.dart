import 'package:flutter/material.dart';
import 'UBCircularImage.dart';
import 'UBGreyContainer.dart';
import 'UBText.dart';
import '../../../generated/colors.gen.dart';

class UBDDMockButton extends StatelessWidget {
  final String title;
  final String titleAppendix;
  final String iconAddress;
  final double horizontalPadding;
  final Color backgroundColor;
  final Widget endIcon;
  final void Function() onTap;
  const UBDDMockButton({
    Key key,
    this.title,
    this.onTap,
    this.titleAppendix,
    this.iconAddress,
    this.horizontalPadding = 12,
    this.backgroundColor = ColorName.black,
    this.endIcon = const Icon(
      Icons.keyboard_arrow_down,
      color: ColorName.greybf,
      size: 24,
    ),
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: () => onTap(),
      child: UBGreyContainer(
        color: backgroundColor,
        margin: EdgeInsets.symmetric(
          horizontal: horizontalPadding,
        ),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Container(
              child: Row(
                children: [
                  if (iconAddress != null)
                    UBCircularImage(
                      imageAddress: iconAddress,
                      //size: 20,
                    ),
                  if (iconAddress != null)
                    const SizedBox(
                      width: 4,
                    ),
                  UBText(
                    text: title,
                    color: ColorName.greybf,
                  ),
                  const SizedBox(
                    width: 4,
                  ),
                  if (titleAppendix != null)
                    UBText(
                      text: "($title)",
                      size: 8,
                      color: ColorName.grey80,
                    ),
                ],
              ),
            ),
            endIcon
          ],
        ),
      ),
    );
  }
}
