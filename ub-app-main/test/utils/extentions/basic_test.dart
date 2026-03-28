import 'package:flutter_test/flutter_test.dart';
import 'package:unitedbit/utils/extentions/basic.dart';

void main() {
  group('Extentions-basic-strings--->', () {
    test('can add comma in between numbers for currency strings', () {
      final initialString = '1234567890.1234567890';
      final formatted = '1,234,567,890.1234567890';
      expect(initialString.currencyFormat(), formatted);
    });
    test('can remove insignificant zeros', () {
      final initialString = '1234567890.1234000000';
      final formatted = '1,234,567,890.1234';
      expect(
          initialString.currencyFormat(
            removeInsignificantZeros: true,
          ),
          formatted);
    });
    test('can format to cent from no fraction', () {
      final initialString = '1234567890';
      final formatted = '1,234,567,890.00';
      expect(initialString.currencyFormat(centFormat: true), formatted);
    });
    test('can format to cent from .0', () {
      final initialString = '1234567890.0';
      final formatted = '1,234,567,890.00';
      expect(initialString.currencyFormat(centFormat: true), formatted);
    });
    test('can convert formatted to double', () {
      final initialString = '1,234,567,890.00';
      final doubled = 1234567890.0;
      expect(initialString.toDouble(), doubled);
    });
    test('can remove comma', () {
      final initialString = '1,234,567,890.00';
      final removed = '1234567890.00';
      expect(initialString.removeComma(), removed);
    });
  });
}

class Counter {
  int value = 0;

  void increment() => value++;

  void decrement() => value--;
}
