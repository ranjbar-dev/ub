import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../../../generated/colors.gen.dart';
import '../../../../../utils/commonUtils.dart';
import "../../../../../utils/extentions/basic.dart";
import '../../../../common/components/UBShimmer.dart';
import '../../../../common/components/UBText.dart';
import '../../../../common/components/controlledInput.dart';

class TradeInput extends StatefulWidget {
  final String endText;
  final String title;
  final Function onChange;
  final String value;
  final bool gettingReady;
  final Color focusColor;
  final int precision;

  const TradeInput({
    Key key,
    @required this.endText,
    @required this.title,
    this.onChange,
    this.value,
    this.gettingReady = false,
    this.focusColor,
    this.precision,
  }) : super(key: key);

  @override
  _TradeInputState createState() => _TradeInputState();
}

class _TradeInputState extends State<TradeInput> {
  bool focused = false;

  handleFocusChanged(isFocused) {
    setState(() {
      focused = isFocused;
    });
  }

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        Container(
          height: 32,
          decoration: BoxDecoration(
            border: Border.all(
                color: focused
                    ? (widget.focusColor ?? ColorName.primaryBlue)
                    : Colors.transparent,
                width: 1),
            borderRadius: const BorderRadius.all(
              const Radius.circular(6),
            ),
            color: ColorName.black,
          ),
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              Expanded(
                flex: 10,
                child: Container(
                  height: 32,
                  child: Center(
                    child: Padding(
                      padding: EdgeInsets.only(
                        left: 8.0,
                      ),
                      child: ControlledTextField(
                        autoFocus: false,
                        noBorder: true,
                        labelText: widget.title,
                        type: TextInputType.number,
                        onFocusChanged: handleFocusChanged,
                        textStyle: const TextStyle(
                          fontWeight: FontWeight.w700,
                          fontSize: 14,
                          color: ColorName.white,
                        ),
                        text: formattedText(widget.value),
                        onChanged: handleChange,
                        maxLength: maxCalculator(widget.value),
                      ),
                    ),
                  ),
                ),
              ),
              Container(
                height: 24,
                width: 2,
                color: ColorName.black2c,
              ),
              Expanded(
                flex: 3,
                child: Center(
                  child: UBText(
                    size: 10,
                    color: ColorName.white,
                    text: widget.endText,
                  ),
                ),
              )
            ],
          ),
        ),
        if (widget.gettingReady == true)
          ClipRRect(
            borderRadius: const BorderRadius.all(Radius.circular(4)),
            child: UBShimmer(
              width: Get.width,
              height: 31.0,
            ),
          )
      ],
    );
  }

  String formattedText(String value) {
    final precision =
        widget.precision ?? coinPrecision(coinCode: widget.endText);
    //print("widget.precision ${widget.precision} && precision $precision");

    return formatCurrencyWithMaxFraction(value: value, maxFraction: precision);
  }

  formatCurrencyWithMaxFraction({
    String value,
    int maxFraction = 8,
  }) {
    String formatted = value;
    if (value.contains('.')) {
      if (value.startsWith('.')) {
        value = '0' + value;
      }
      List<String> splitted = value.split('.');
      splitted[0] = splitted[0].simpleCurrencyFormat();
      if (splitted.length == 2) {
        int currectFraction =
            splitted[1].length > maxFraction ? maxFraction : splitted[1].length;

        final tail = splitted[1].substring(0, currectFraction);

        formatted = splitted[0] + '.' + tail;
      } else {
        formatted = splitted[0] + '.';
      }

      return formatted;
    }
    return formatted.simpleCurrencyFormat();
  }

  void handleChange(String text) {
    String pure = text.removeComma();
    if (pure.startsWith('.')) {
      pure = "0" + pure;
    }
    return this.widget.onChange(pure);
  }

  int maxCalculator(String value) {
    final precision =
        widget.precision ?? coinPrecision(coinCode: widget.endText);
    final formatted =
        formatCurrencyWithMaxFraction(value: value, maxFraction: precision);
    if (formatted.contains(".")) {
      final splitted = formatted.split(".");
      return splitted[0].length + precision + 1;
    }
    return formatted.length < 15 ? formatted.length + 1 : 15;
  }
}
