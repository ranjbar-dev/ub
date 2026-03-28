import 'package:get/get_rx/get_rx.dart';
import 'package:get_storage/get_storage.dart';
import '../services/constants.dart';
import '../app/global/currency_model.dart';
import '../app/global/currency_pairs_model.dart';
import '../services/storageKeys.dart';

class PairAndCurrencyUtils {
  static GetStream<Map<String, Pairs>> pairsMap = GetStream();
  static GetStream<Map<String, CurrencyModel>> coinsMap = GetStream();
  static GetStream<Map<String, CurrencyPairsModel>> currencyPairsMap =
      GetStream();
  static final GetStorage storage = GetStorage();
  static PairAndCurrencyUtils _singleton;
  factory PairAndCurrencyUtils() {
    if (_singleton == null) {
      _singleton = PairAndCurrencyUtils._internal();
    }
    return _singleton;
  }

  PairAndCurrencyUtils._internal() {
    final hashMapCoinsJson = storage.read(StorageKeys.coinsHashMap);
    if (hashMapCoinsJson != null) {
      final Map<String, CurrencyModel> tmp = {};
      hashMapCoinsJson.forEach((key, value) {
        tmp[key] = CurrencyModel.fromJson(value);
      });
      coinsMap.add(tmp);
    }

    final hashMapPairsJason = storage.read(StorageKeys.pairsHashMap);
    if (hashMapPairsJason != null) {
      final Map<String, Pairs> tmp = {};
      hashMapPairsJason.forEach((key, value) {
        tmp[key] = Pairs.fromJson(value);
      });
      pairsMap.add(tmp);
    }
  }

  static void setPairsHashMap(Map hashMapJson) {
    final Map<String, Pairs> tmp = {};
    hashMapJson.forEach((key, value) {
      tmp[key] = Pairs.fromJson(value);
    });
    pairsMap.add(tmp);
    storage.write(StorageKeys.pairsHashMap, hashMapJson);
  }

  static void setCoinsHashMap(Map hashMapJson) {
    final Map<String, CurrencyModel> tmp = {};
    hashMapJson.forEach((key, value) {
      tmp[key] = CurrencyModel.fromJson(value);
    });
    coinsMap.add(tmp);
    storage.write(StorageKeys.coinsHashMap, hashMapJson);
  }

  static void setCurrencyPairsHashMap(Map hashMapJson) {
    final Map<String, CurrencyPairsModel> tmp = {};
    hashMapJson.forEach((key, value) {
      tmp[key] = CurrencyPairsModel.fromJson(value);
    });
    currencyPairsMap.add(tmp);
    storage.write(StorageKeys.currencyPairsHashMap, hashMapJson);
  }

  static String findCoinNameByCode(String coineCode) {
    final coinsList = Constants.currencyArray();
    var currentCoin =
        coinsList.firstWhere((element) => element.code == coineCode);
    return currentCoin.desc;
  }

  static String findCoinImageByCode(String coineCode) {
    final coinsList = Constants.currencyArray();
    var currentCoin =
        coinsList.firstWhere((element) => element.code == coineCode);
    return currentCoin.image;
  }
}
