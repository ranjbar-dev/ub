import 'dart:math';

import 'package:intl/intl.dart';

RegExp commaFormat = RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))');
String Function(Match match) mathFunc = (Match match) => '${match[1]},';

bool isNumeric(String s) {
  if (s == null) {
    return false;
  }
  return double.tryParse(s) != null;
}

extension TruncateDoubles on double {
  String toFixedWithoutRounding(int fractionalDigits, {bool addZeroes = true}) {
    var fixed = ((this * pow(10, fractionalDigits)).truncate() /
        pow(10, fractionalDigits));
    return addZeroes ? fixed.toStringAsFixed(fractionalDigits) : fixed;
  }
}

extension UBDoubleExtensions on double {
  String toStringWithoutTrailingZeros({bool addZeros}) {
    if (this == null) return null;
    if (addZeros == true) {
      final val = truncateToDouble() == this ? toInt().toString() : toString();
      if (val.split('.').length == 1) {
        return val + '.00';
      } else if (val.split('.')[1].length == 1) {
        return val + '0';
      }
    }
    return truncateToDouble() == this ? toInt().toString() : toString();
  }

  String parseDoubleToString() {
    final s = this.toString();
    final splitted = s.split('.');
    if (splitted[1] == '0') {
      return splitted[0];
    }
    if (splitted[1].length > 12) {
      return splitted[0] + '.' + splitted[1].substring(0, 12);
    }
    return s;
  }
}

extension UBStringExtentions on String {
  String simpleCurrencyFormat() {
    final splitted = this.split('.');
    return splitted[0].replaceAllMapped(commaFormat, mathFunc) +
        "${splitted.length == 2 ? ('.' + splitted[1] ?? '') : ''}";
  }

  String removeComma() {
    if (this == null) return null;
    if (this == '') return '';

    if (this.contains(',')) {
      return this.replaceAll(',', '');
    }
    return this;
  }

  /*also removes commas*/
  double toDouble() {
    if (this == null) return null;
    if (this == '') return null;
    String main = this;
    if (main.contains(',')) {
      main = main.removeComma();
    }
    if (main.startsWith('.')) {
      main = '0' + main;
    }
    if (main.contains('..')) {
      return 0.0;
    }
    // if (!isNumeric(this)) {
    //   return 0.0;
    // }
    return double.parse(main);
  }

  String operator &(String other) {
    return '$this $other';
  }

  String currencyFormat({
    bool removeInsignificantZeros = false,
    bool centFormat = false,
    int toFixed,
    bool compact = false,
    bool formatSmall = false,
  }) {
    if (this == null) return null;
    if (this == '') return '';
    if (compact == true) {
      final v = NumberFormat.compactCurrency(decimalDigits: 2, symbol: '\$')
          .format(double.parse(this.replaceAll(',', '')));
      return v;
    }
    String main = this;
    if (main.contains(',')) {
      return main = main.replaceAll(',', '');
    }
    if (!isNumeric(main) && !main.endsWith('.')) {
      return '';
    }
    if (main.contains('..')) {
      return '0.0';
    }
    if (removeInsignificantZeros) {
      main = double.parse(main).toStringWithoutTrailingZeros();
    }
    final splitted = main.split('.');
    int trailingLength = 0;
    //for small numbers
    if (splitted[0].length < 4 && formatSmall) {
      if (removeInsignificantZeros && !centFormat) {
        return main;
      }
      if (centFormat) {
        if (splitted.length == 2) {
          trailingLength = splitted[1].length;
        }
        if (trailingLength == 0) {
          main = main + '.00';
          return main;
        }
        if (trailingLength == 1) {
          main = main + '0';
          return main;
        }
      }

      return main;
    }

    main = splitted[0].replaceAllMapped(commaFormat, mathFunc);
    if (splitted.length == 2) {
      trailingLength = splitted[1].length;

      main = main + (splitted[1] != '' ? ('.' + splitted[1]) : '');
    }

    if (centFormat == true) {
      if (trailingLength == 0) {
        main = main + '.00';
        return main;
      }
      if (trailingLength == 1) {
        main = main + '0';
        return main;
      }
    }
    if (toFixed != null && splitted.length == 2) {
      if (splitted[1].length > toFixed) {
        return splitted[0] + '.' + splitted[1].substring(0, toFixed);
      }
    }
    return main;
  }
}
