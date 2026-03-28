import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../generated/colors.gen.dart';

enum ButtonVariant {
  Filled,
  Outline,
  TransparentBackground,
  Rounded,
  Link,
}

class UBButton extends StatelessWidget {
  final double height;
  final double width;
  final Function onClick;
  final String text;
  final Color textColor;
  final Color buttonColor;
  final Color borderColor;
  final ButtonVariant variant;
  final double fontSize;
  final bool isLodaing;
  final bool smallLoading;

  final bool disabled;
  final TextDecoration textDecoration;
  final Widget endWidget;
  final EdgeInsets padding;
  final double borderRadius;
  const UBButton({
    Key key,
    this.height = 38.0,
    @required this.onClick,
    @required this.text,
    this.textColor = Colors.white,
    this.variant,
    this.buttonColor,
    this.fontSize = 14.0,
    this.isLodaing,
    this.disabled,
    this.padding,
    this.borderColor,
    this.width,
    this.endWidget,
    this.borderRadius,
    this.textDecoration = TextDecoration.none,
    this.smallLoading,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final loadingSize = smallLoading == true ? 12.0 : 17.0;
    final loadingWidth = smallLoading == true ? 1.0 : 2.0;
    final mainColor = disabled == true
        ? ColorName.grey19
        : buttonColor ?? ColorName.primaryBlue;
    final br = BorderRadius.circular(
      borderRadius ?? (variant == ButtonVariant.Rounded ? 100.0 : 8.0),
    );
    final background = (variant == ButtonVariant.Link ||
            variant == ButtonVariant.TransparentBackground)
        ? Colors.transparent
        : mainColor;
    return Container(
      width: width,
      decoration: BoxDecoration(
        borderRadius: br,
        border: Border.all(
          color: borderColor ?? background,
          width: variant == ButtonVariant.Outline ? 1 : 0,
        ),
      ),
      height: height,
      child: Material(
        borderRadius: br,
        color: buttonBackgroundColor(
            variant: variant ?? ButtonVariant.Filled, color: background),
        child: InkWell(
          borderRadius: br,
          child: Center(
            child: isLodaing == true
                ? Container(
                    width: loadingSize,
                    height: loadingSize,
                    child: CircularProgressIndicator(
                      valueColor:
                          AlwaysStoppedAnimation<Color>(ColorName.white),
                      strokeWidth: loadingWidth,
                      backgroundColor: Colors.transparent,
                    ),
                  )
                : Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    crossAxisAlignment: CrossAxisAlignment.center,
                    children: [
                      Container(
                        padding: padding,
                        child: RichText(
                          text: TextSpan(
                            text: text,
                            style: TextStyle(
                              color: disabled == true
                                  ? ColorName.grey36
                                  : variant == ButtonVariant.Outline
                                      ? (textColor ?? background)
                                      : variant == ButtonVariant.Link
                                          ? (textColor ?? ColorName.primaryBlue)
                                          : textColor,
                              fontSize: fontSize,
                              fontWeight: FontWeight.w600,
                              decoration: textDecoration,
                            ),
                          ),
                        ),
                      ),
                      if (endWidget != null)
                        const SizedBox(
                          width: 4,
                        ),
                      if (endWidget != null) endWidget
                    ],
                  ),
          ),
          onTap: () {
            if (isLodaing != true && disabled != true) {
              if (GetPlatform.isWeb) {
                FocusScope.of(context).unfocus();
              }
              FocusManager.instance.primaryFocus?.unfocus();
              onClick();
              return;
            } else {
              return null;
            }
          },
        ),
      ),
    );
  }

  buttonBackgroundColor({@required ButtonVariant variant, @required color}) {
    switch (variant) {
      case ButtonVariant.Filled:
      case ButtonVariant.Rounded:
        return color;
      case ButtonVariant.Link:
        return Colors.transparent;
      case ButtonVariant.Outline:
      case ButtonVariant.TransparentBackground:
        return Colors.transparent;
      default:
        return ColorName.primaryBlue;
    }
  }
}
