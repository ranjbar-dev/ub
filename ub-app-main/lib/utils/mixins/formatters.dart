RegExp commaFormat = RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))');

RegExp removeTrailingZerosRegexp = RegExp(r"([.]*0)(?!.*\d)");
// ignore: prefer_function_declarations_over_variables
String Function(Match match) mathFunc = (Match match) => '${match[1]},';

mixin Formatter {
  bool isNumeric(String s) {
    if (s == null) {
      return false;
    }
    return double.tryParse(s) != null;
  }

  String currencyFormatter(String v,
      {bool twoFractionFix = true, bool removeTrailingZeros = true}) {
    if (v.contains('..')) {
      return '0';
    }
    v = v.replaceAll(' ', '');
    if (v == '0.') {
      return v;
    }
    if (isNumeric(v) && double.parse(v) == 0) {
      return v;
    }
    if (v == '.') {
      return '0.';
    }

    if (v.contains('.')) {
      var fixedZero = v;
      if (v.split('.')[1] != '0' && removeTrailingZeros) {
        fixedZero = fixedZero.replaceAll(removeTrailingZerosRegexp, "");
      }
      var main =
          fixedZero.split('.')[0].replaceAllMapped(commaFormat, mathFunc);
      var trailing = fixedZero.split('.')[1];
      if (trailing.length == 1 && twoFractionFix) {
        trailing = trailing + '0';
      }
      return main + '.' + trailing;
    }
    var main = v.replaceAllMapped(commaFormat, mathFunc);
    return main;
  }

  String removeDecimalZeroFormat(double n) {
    return n.toStringAsFixed(n.truncateToDouble() == n ? 0 : 1);
  }
}
