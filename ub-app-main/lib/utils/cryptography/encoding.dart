// import 'package:encrypt/encrypt.dart';
// import 'package:encrypt/encrypt_io.dart';
import 'package:encrypt/encrypt.dart';
// import 'package:encrypt/encrypt_io.dart';
// import 'package:fast_rsa/rsa.dart' as Fast;
// import 'package:pointycastle/pointycastle.dart';
// import 'package:pointycastle/asymmetric/api.dart';
import 'package:unitedbit/utils/commonUtils.dart';
import 'package:meta/meta.dart' show required;

const encPrefix = "ub-captcha_";

encrypt(String string) async {
  final parser = RSAKeyParser();
  final publicKey = parser.parse(outKPem);
  final encrypter = Encrypter(RSA(
    publicKey: publicKey,
  ));
  final encrypted = encrypter.encrypt(string).base64;

  // final encrypted = await Fast.RSA.encryptPKCS1v15(string, outKPem);
  return encPrefix + encrypted;
}

// encryptor(String string)async {
//     final publicKey = await parseKeyFromFile<RSAPublicKey>('assets/files/public.pem');
//   final encrypter = Encrypter(RSA(publicKey: publicKey, ));

// return    encrypter.encrypt(string);

// }
genarateEnc({@required DateTime startTime}) async {
  final endTime = DateTime.now();

  final now = DateTime.now().toUtc().millisecondsSinceEpoch ~/ 1000;
  final spent = endTime.difference(startTime).inSeconds;

  final map = {};
  map['"timestamp"'] = now;
  map['"spent_time"'] = spent;

  final string = map.toString();
  final encoded = await encrypt(string);

  return encoded;
}
