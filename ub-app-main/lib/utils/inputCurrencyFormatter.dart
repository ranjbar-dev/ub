import 'package:flutter/services.dart';

import 'mixins/formatters.dart';

class InputCurrencyFormatter extends TextInputFormatter with Formatter {
  InputCurrencyFormatter({this.maxDigits});
  final int maxDigits;

  TextEditingValue formatEditUpdate(
      TextEditingValue oldValue, TextEditingValue newValue) {
    if (newValue.selection.baseOffset == 0) {
      return newValue;
    }

    if (maxDigits != null && newValue.selection.baseOffset > maxDigits) {
      return oldValue;
    }

    var newText = newValue.text;
    return newValue.copyWith(
        text: newValue.text,
        selection: new TextSelection.collapsed(offset: newText.length));
  }
}
