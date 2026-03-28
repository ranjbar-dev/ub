import 'package:flutter_test/flutter_test.dart';
import 'package:mockito/annotations.dart';
import 'package:mockito/mockito.dart';
import 'package:unitedbit/app/modules/account/providers/accountProvider.dart';

import 'accountProvider_test.mocks.dart';

@GenerateMocks([AccountProvider])
void main() {
  final accountProvider = MockAccountProvider();
  group('Account provider--->', () {
    test('can resend verification email', () async {
      when(accountProvider.requestForEmail()).thenAnswer(
          (_) async => Future.value({"status": true, "message": ""}));

      expect(await accountProvider.requestForEmail(),
          equals({"status": true, "message": ""}));
    });
  });
}
