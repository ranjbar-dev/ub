import 'package:flutter_test/flutter_test.dart';
import 'package:unitedbit/utils/cryptography/rsa_encryption.dart';

void main() {
  var helper = RsaKeyHelper();
  group('Encryption --->', () {
    test('can encrypt and decrypt', () async {
      final toEncrypt = {
        "date": DateTime.now().toUtc().toString(),
        "code": "1234567890"
      }.toString();
      final publicKey = helper.parsePublicKeyFromPem(
          'MIIBCgKCAQEAget5DCIFBb97c79tc5Jhk1+0dpCoAkA+hNi9yMVuHrFM3FnRiqzNQWoJa55imC3JvitKMuZq3rOmlRC/tFxJhdMXKrDfNIwF6ODAJ6APbGhDqgPozxvX41UTdC/3R0gbQ4wL3XGZvWWZU/azqn8YOT8LRgtiRM1RH8flFi4enwfqQOBx9PoeVUTYtfUWxyJA8PaCZl9LKldaybyKUvCw38MiOH97k/0ZJm/qxAfsCTKtCcfSzOudilyyOrqSZt91rv88H1MQSmDPpZLsEfapfJuBGeO8rib31+RZenI6WnbTmn0fmNJBPYvT/N5dSmGly3yBLpWxIeZftENqRbo+xQIDAQAB');

      final privateKey = helper.parsePrivateKeyFromPem(
          'MIIFogIBAAKCAQEAget5DCIFBb97c79tc5Jhk1+0dpCoAkA+hNi9yMVuHrFM3FnRiqzNQWoJa55imC3JvitKMuZq3rOmlRC/tFxJhdMXKrDfNIwF6ODAJ6APbGhDqgPozxvX41UTdC/3R0gbQ4wL3XGZvWWZU/azqn8YOT8LRgtiRM1RH8flFi4enwfqQOBx9PoeVUTYtfUWxyJA8PaCZl9LKldaybyKUvCw38MiOH97k/0ZJm/qxAfsCTKtCcfSzOudilyyOrqSZt91rv88H1MQSmDPpZLsEfapfJuBGeO8rib31+RZenI6WnbTmn0fmNJBPYvT/N5dSmGly3yBLpWxIeZftENqRbo+xQKCAQATDP5hAxQNdbiajnV0PwDD5YLG6Ata2STRwh6CNEEjiwgkP590YEZw0yWyfDUk74Hnut1UfWkqYtmIfj4+KlI1p3B8OBdi0y2CqoJCzTu1v8w8P/qBdCnCEhWaMfZmo3IsA3sx65iJpz5Gi6Ro2d4pds8mZEDyqdC9gkhbakPfOcgdfgstb79KaD2xcvwU62/9wyI4ziMKbCIk6jWSsmct2ZKcf8GGndkHvIAAlyNPXKnjKTEeo6SdFalqrm+1NtDttqlEO9vHFsNTg4QOKzaIhNTiLzK93GD+H86Yn0OD2lfkNpdX65Re6S4BTCRo1lk5krvegB4LI52Oeway+3l5AoIBABMM/mEDFA11uJqOdXQ/AMPlgsboC1rZJNHCHoI0QSOLCCQ/n3RgRnDTJbJ8NSTvgee63VR9aSpi2Yh+Pj4qUjWncHw4F2LTLYKqgkLNO7W/zDw/+oF0KcISFZox9majciwDezHrmImnPkaLpGjZ3il2zyZkQPKp0L2CSFtqQ985yB1+Cy1vv0poPbFy/BTrb/3DIjjOIwpsIiTqNZKyZy3Zkpx/wYad2Qe8gACXI09cqeMpMR6jpJ0VqWqub7U20O22qUQ728cWw1ODhA4rNoiE1OIvMr3cYP4fzpifQ4PaV+Q2l1frlF7pLgFMJGjWWTmSu96AHgsjnY57BrL7eXkCgYEA9DyRlwzxDG1bwJYuqKVLhMY2rPcmgHK/KbGYX4Qkdg4tCkIwp6h24TUjCMNMQyOshHCYZNKKUeS/32cjkNbclblYhZwZvYWttKuiF1Ahst4NCEYIc6bgcn5+9yDrnOTKAww/SkDVH/qwEozcNJMxVjhQjFvXIaZCfzVzvd9J6f8CgYEAiC1hiavp1GtOX2sJ5Hb9VlIywSMhrGMeK6HNQhDqjFO5WesOd6+oHNXSTCgTMqjlA1yiBdMKa1PuKlpVW4E9Q3ftLGSBg7SYZvFyd7Bt9WavXD5bVywPg/NZASjVaMqajBQrTfN7caT0c77oszY8dN8X/fRBqamFfXMFc99ErzsCgYEAw7jyFzFIzmOovonbtExaW3mYgT3CPfc2mEv4xrqXmX+8ulbWtNS9B7bUb4ZKTBd/fdbZWRqbvArrdDUr/DsjJF0WwmOZARbqYDmWuMX/a16k5PdyeHPHtBkI2DQqfF2gQZcD9RZFdM4pYYQ+R2eZhvW0HvbOTOn2qgiEyyjwC7MCgYBNxCy5ZDWKmyUMlKH3mIQgMZzOcvOd8JSgMix3mBV5wa5NzVBbxTJqFSmdWB1uhskR3GqijNycYjfWc/Pe57VGvEvzWAomXpHR5/yIoXaJ9/QY53teEslhfyzK3rjQuTL839/DClLqmVsIZnOZNFXeIDEhU8XTz/1toKvyegRNVQKBgB37JLurKefhrLU1Nz1hcb+KIG/wSuV0YNkWC/AQU/ExI1AnEu51cyYsvmHnjEsdAbOSwEdt3IFZ0qzey+14pRmVwmyXXvri+CFBedp0tnciCrCfjaLby1djgKeGV3W5IJFuljNFpu5GhExQTmeYkSmzPnv/1Kx1LCysMl0f1Wa+');

      final encrypted =
          encrypt(encryptModel: EncryptModel(toEncrypt, publicKey));
      final decrypted =
          decrypt(decryptModel: DecryptModel(encrypted, privateKey));
      expect(toEncrypt, decrypted);
    });
  });
}
