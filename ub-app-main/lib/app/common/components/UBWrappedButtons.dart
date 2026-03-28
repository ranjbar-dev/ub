import 'package:flutter/material.dart';
import '../../../generated/assets.gen.dart';
import '../../../generated/colors.gen.dart';
import '../../../utils/mixins/commonConsts.dart';

import 'UBText.dart';

class WrappedButtonModel {
  final String text;
  final String value;
  WrappedButtonModel({this.value, @required this.text});
}

class UBWrappedButtons extends StatelessWidget {
  final List<WrappedButtonModel> buttons;
  final void Function(int index) onButtonClick;
  final int selectedIndex;
  final Color unselectedButtonTextColor;
  final Color buttonBackground;
  final double minButtonWidth;
  final double buttonHeight;
  final List<int> otherNoMarginLeftIndexes;

  const UBWrappedButtons({
    Key key,
    @required this.buttons,
    @required this.onButtonClick,
    @required this.selectedIndex,
    this.buttonBackground = ColorName.black2c,
    this.minButtonWidth,
    this.unselectedButtonTextColor = ColorName.grey80,
    this.otherNoMarginLeftIndexes = const [],
    this.buttonHeight,
  });

  @override
  Widget build(BuildContext context) {
    return Wrap(
      crossAxisAlignment: WrapCrossAlignment.start,
      children: [
        for (var i = 0; i < buttons.length; i++)
          SmallButton(
            onClick: () => onButtonClick(i),
            borderColor:
                i == selectedIndex ? ColorName.primaryBlue : buttonBackground,
            color: buttonBackground,
            margin: EdgeInsets.only(
              left: ((i == 0.0) ||
                      otherNoMarginLeftIndexes.indexWhere((e) => i == e) != -1)
                  ? 0.0
                  : 8.0,
              bottom: 6.0,
              top: 6.0,
              right: (i == buttons.length - 1) ? 0.0 : 8.0,
            ),
            minWidth: minButtonWidth,
            height: buttonHeight,
            text: buttons[i].text,
            textColor: i == selectedIndex
                ? ColorName.primaryBlue
                : unselectedButtonTextColor,
          ),
      ],
    );
  }
}

class SmallButton extends StatelessWidget {
  final Function onClick;
  final Function onCloseClick;
  final double minWidth;
  final double height;
  final EdgeInsets margin;
  final Color color;
  final Color borderColor;
  final Color textColor;
  final String text;

  const SmallButton({
    Key key,
    this.onClick,
    this.minWidth,
    this.margin,
    this.color = ColorName.black1c,
    this.borderColor = ColorName.black1c,
    this.textColor = ColorName.grey80,
    this.text = '',
    this.onCloseClick,
    this.height,
  }) : super(key: key);
  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: () {
        onClick();
      },
      child: Container(
          height: height,
          width: minWidth,
          padding: EdgeInsets.only(
              left: 15,
              right: onCloseClick == null ? 15 : 0,
              top: 3,
              bottom: 3),
          margin: margin,
          decoration: BoxDecoration(
            color: color,
            border: Border.all(
              color: borderColor,
              width: 1,
            ),
            borderRadius: BorderRadius.circular(
              4,
            ),
          ),
          child: Row(
            mainAxisSize: MainAxisSize.min,
            mainAxisAlignment: onCloseClick != null
                ? MainAxisAlignment.spaceBetween
                : MainAxisAlignment.center,
            children: [
              UBText(
                align: TextAlign.center,
                text: text,
                color: textColor,
                weight: FontWeight.w600,
                size: 12,
              ),
              if (onCloseClick != null) hspace4,
              if (onCloseClick != null)
                GestureDetector(
                  onTap: onCloseClick,
                  child: Container(
                    width: 24,
                    height: 24,
                    color: color,
                    child: Assets.images.closeIcon.svg(),
                  ),
                )
            ],
          )),
    );
  }
}
