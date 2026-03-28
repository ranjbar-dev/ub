import 'dart:async';

import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:fluttertoast/fluttertoast.dart';
import 'package:get/get.dart';
import 'package:get_storage/get_storage.dart';
import 'package:mqtt_client/mqtt_client.dart' show MqttSubscriptionStatus;
import 'package:permission_handler/permission_handler.dart';
import 'package:unitedbit/mqttClient/universal_mqtt_client.dart'
    show UniversalMqttClient;
import 'package:unitedbit/services/localAuthService.dart';
import 'package:unitedbit/services/storageKeys.dart';
import 'package:unitedbit/utils/mixins/commonConsts.dart';
import 'package:unitedbit/utils/pairAndCurrencyUtils.dart';
import 'package:url_launcher/url_launcher.dart';

import "./extentions/basic.dart";
import '../app/global/autocompleteModel.dart';
import '../app/global/currency_model.dart';
import '../app/global/currency_pairs_model.dart';
import '../app/modules/funds/pages/deposits/controllers/deposits_controller.dart';
import '../app/modules/funds/pages/deposits/views/depositDetails.dart';
import '../generated/colors.gen.dart';
// bool isNumeric(String s) {
//   if (s == null) {
//     return false;
//   }
//   return double.tryParse(s) != null;
// }

bool canExitApp = false;
showDoubleTapToast() {
  Fluttertoast.showToast(
      msg: "Tab back again to exit",
      toastLength: Toast.LENGTH_SHORT,
      gravity: ToastGravity.BOTTOM,
      timeInSecForIosWeb: 1,
      backgroundColor: ColorName.grey22,
      textColor: ColorName.white,
      fontSize: 14.0);
}

void saveCoinToHistory({
  @required AutoCompleteItem coin,
  @required RxList<AutoCompleteItem> stream,
  @required String storageKey,
}) {
  final storage = GetStorage();
  // ignore: invalid_use_of_protected_member
  final List<AutoCompleteItem> tmp = List.from(stream.value);
  if (tmp.indexWhere((element) => element.id == coin.id) == -1) {
    if (tmp.length < 8) {
      tmp.insert(0, coin);
      stream.assignAll(tmp);
    } else {
      tmp.insert(0, coin);
      tmp.removeLast();
      stream.assignAll(tmp);
    }
    final tmp2 = [];
    stream.forEach((element) {
      tmp2.add(element.toJson());
    });
    storage.write(storageKey, tmp2);
  }
}

openDepositPopup({@required AutoCompleteItem coin}) async {
  Get.put<DepositsController>(DepositsController());
  final DepositsController depositController = Get.find();
  try {
    await depositController.handleCoinSelected(coin);
    // Get.toNamed(AppPages.DEPOSITDETAILS);
    Get.dialog(
      Container(
        clipBehavior: Clip.antiAlias,
        margin: const EdgeInsets.all(12.0),
        decoration: BoxDecoration(
          borderRadius: rounded_big,
          border: Border.all(color: ColorName.black2c, width: 1),
        ),
        height: Get.height,
        width: Get.width,
        child: DepostDetailsView(coin: coin),
      ),
    );
  } catch (e) {}
}

void launchURL(String url) async => await canLaunchUrl(Uri.parse(url))
    ? await launchUrl(Uri.parse(url))
    : throw 'Could not launch $url';

Future checkCameraPermission(
    {@required Function onGranted,
    @required Function onDenied,
    bool showPrompForSetting = true}) async {
  PermissionStatus status = await Permission.camera.status;
  if (status.isDenied) {
    PermissionStatus requestStatus = await Permission.camera.request();
    if (requestStatus.isGranted) {
      onGranted();
    } else {
      if (showPrompForSetting) {
        promptForPermissionInSetting(
          onDenied: onDenied,
          title: 'Camera Permission',
          desc: 'Unitedbit needs camera access to scan qr codes',
        );
      } else {
        onDenied();
      }
    }
  } else {
    onGranted();
  }
  return Future.value();
}

Future<bool> canContinueWithBiometrics() async {
  final storage = GetStorage();
  final biometricsActivated = storage.read(StorageKeys.biometricsActivated);
  if (!(GetPlatform.isWeb)) {
    if (biometricsActivated == true) {
      final authed = await BiometricsService().authenticateWithBiometrics();
      return authed;
    }
  }
  return true;
}

Future checkGalleryPermission(
    {@required Function onGranted,
    @required Function onDenied,
    bool showPrompForSetting = true}) async {
  PermissionStatus status = await Permission.photos.status;
  if (status.isDenied) {
    PermissionStatus requestStatus = await Permission.photos.request();
    if (requestStatus.isGranted) {
      onGranted();
    } else {
      if (showPrompForSetting) {
        promptForPermissionInSetting(
          onDenied: () {},
          title: 'Photos Permission',
          desc: 'Unitedbit needs access to images',
        );
      } else {
        onDenied();
      }
    }
  } else {
    onGranted();
  }
  return Future.value();
}

Future purgeTopic({
  @required UniversalMqttClient client,
  @required StreamSubscription topicStream,
  @required String topic,
}) async {
  if (client != null) {
    final topicStatus = client.getSubscriptionsStatus(topic: topic);
    if (topicStatus == MqttSubscriptionStatus.active) {
      print('debugPrint:' + topicStatus.toString());
      client.unsubscribe(topic: topic);
      if (topicStream != null) {
        return await topicStream.cancel();
      }
    }
  }
}

void promptForPermissionInSetting({
  @required String title,
  @required String desc,
  Function onDenied,
}) {
  if (GetPlatform.isIOS) {
    Get.dialog(CupertinoAlertDialog(
      title: Text(title),
      content: Text(desc),
      actions: <Widget>[
        CupertinoDialogAction(
          child: Text('Deny'),
          onPressed: () {
            Get.back();
            onDenied();
          },
        ),
        CupertinoDialogAction(
            child: Text('Settings'),
            onPressed: () {
              Get.back();
              openAppSettings();
            })
      ],
    ));
  } else if (GetPlatform.isAndroid) {
    Get.dialog(
      AlertDialog(
        title: Text(title),
        content: Text(desc),
        actions: <Widget>[
          TextButton(
            child: Text('Deny'),
            onPressed: () {
              Get.back();
              onDenied();
            },
          ),
          TextButton(
              child: Text('Settings'),
              onPressed: () {
                Get.back();
                openAppSettings();
              })
        ],
      ),
    );
  }
}

Map coinMap = {};

int coinPrecision(
    {String coinCode = "BTC", Map<String, CurrencyModel> coinHashMap}) {
  String code = coinCode;

  if (coinMap["BTC"] == null) {
    coinMap = coinHashMap ?? Map.from(PairAndCurrencyUtils.coinsMap.value);
  }
  try {
    return coinMap[code].showDigits;
  } catch (e) {
    debugPrint("error getting coinPrecission");
    return 8;
  }
}

String _precisionFormatter({String value, int precision}) {
  return double.parse(value)
      .toFixedWithoutRounding(precision)
      .simpleCurrencyFormat();
}

String decimalCoin(
    {@required String value,
    String coinCode = "BTC",
    Map<String, CurrencyModel> coinHashMap}) {
  final splitted = value.split(" ");
  String code = coinCode;
  if (value.contains(" ")) {
    if (splitted.length == 2) {
      code = splitted[1];
    }
  }
  if (coinMap["BTC"] == null) {
    coinMap = coinHashMap ?? Map.from(PairAndCurrencyUtils.coinsMap.value);
  }
  try {
    final formatted =
        "${_precisionFormatter(value: splitted[0], precision: coinMap[code].showDigits)}${splitted.length == 2 ? (' ' + splitted[1]) : ''}";
    return formatted;
  } catch (e) {
    debugPrint('coin not found when calculating decimalCoin');
    return value;
  }
}

Map pairMap = {};

String decimalPair(
    {@required String value,
    String pairName = "BTC-USDT",
    Map<String, Pairs> pairHashMap}) {
  if (coinMap["BTC"] == null) {
    pairMap = pairHashMap ??
        Map.from(PairAndCurrencyUtils.pairsMap.value ??
            {"BTC-USDT": Pairs(showDigits: 8)});
  }
  try {
    final formatted =
        "${_precisionFormatter(value: value, precision: pairMap[pairName].showDigits)}";
    return formatted;
  } catch (e) {
    debugPrint('pair not found when calculating decimalPair');
    return value;
  }
}

int pairPrecision(
    {String pairName = "BTC-USDT", Map<String, Pairs> pairsHashMap}) {
  String name = pairName;

  if (pairMap["BTC-USDT"] == null) {
    pairMap = pairsHashMap ?? Map.from(PairAndCurrencyUtils.pairsMap.value);
  }
  try {
    return pairMap[name].showDigits;
  } catch (e) {
    debugPrint("error getting coinPrecission");
    return 8;
  }
}

censorAddress({String address}) {
  final toCensor = (address.length) ~/ 3;
  return address.substring(0, toCensor) +
      "*******" +
      address.substring(address.length - toCensor, address.length);
}

const outKPem = """-----BEGIN RSA PUBLIC KEY-----
MIIBCgKCAQEAu4/3ri2L0f60SSmrL2V0KJtKTgno9opoY/N3jiQIVh5xLeqA3NJw
k+rUjrG0/UnANLLeV7LVj1a9bwW5VZEb5TW7PzAy7z9Q1VR7gGSHfvt6DWToIYzD
KkaEBX3bUxwwsg6jC/Jx44jQ8mBQdkf0AW0X5aEzWjZAIo2Vq6wvV2G1fBHekH8t
XIhWAoYPibSEa59oYDtrmkE2J8hVTOITXjJTEA+U2SZ3rwb5C9O7RvmP48HG9UrX
x16r40mDc8R5tW5x3RQk1tBCAEdGhI9ISlPwd6eHbqm4KNEyxe6vCLKrHSipqFBY
etyTnd2rzKe+RLEcUZ2AvlmXplUBeqdrzQIDAQAB
-----END RSA PUBLIC KEY-----
""";
const outK =
    """MIIBCgKCAQEAu4/3ri2L0f60SSmrL2V0KJtKTgno9opoY/N3jiQIVh5xLeqA3NJwk+rUjrG0/UnANLLeV7LVj1a9bwW5VZEb5TW7PzAy7z9Q1VR7gGSHfvt6DWToIYzDKkaEBX3bUxwwsg6jC/Jx44jQ8mBQdkf0AW0X5aEzWjZAIo2Vq6wvV2G1fBHekH8tXIhWAoYPibSEa59oYDtrmkE2J8hVTOITXjJTEA+U2SZ3rwb5C9O7RvmP48HG9UrXx16r40mDc8R5tW5x3RQk1tBCAEdGhI9ISlPwd6eHbqm4KNEyxe6vCLKrHSipqFBYetyTnd2rzKe+RLEcUZ2AvlmXplUBeqdrzQIDAQAB""";
