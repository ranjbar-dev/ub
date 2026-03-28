import 'package:get_storage/get_storage.dart';

import '../app/global/autocompleteModel.dart';
import '../app/global/currency_model.dart';
import '../app/modules/phoneVerification/country_model.dart';
import '../utils/environment/ubEnv.dart';
import 'storageKeys.dart';

enum RxUpdateables {
  Balances,
  TransactionHistory,
  UserPairBalances,
  OpenOrders,
  OrderHistory
}

class Constants {
  static final GetStorage storage = GetStorage();
  static final Constants _singleton = Constants._internal();
  factory Constants() {
    return _singleton;
  }

  Constants._internal();
  static const _urlPrefix = (ENV == "DEV") ? 'dev-' : '';

  static const appVersion = VERSION;
  static const priceTopic = 'main/trade/ticker';
  static const orderbookTopic = 'main/trade/order-book/';
  static const ohlcTopic = 'main/trade/kline/';
  static const landingPageAddress = 'https://www.unitedbit.com';
  static const cmsAddress = 'https://content.unitedbit.com';
  static const webLandingAddress = 'https://www.unitedbit.com';
  static const newsAddress = '$webLandingAddress/news/';

  static const initialPair = 'BTC-USDT';
  static const pre = 'https';
  static const mainUrl = '${_urlPrefix}app.unitedbit.com';
  static const _mqttMainUrl = '${_urlPrefix}app.unitedbit.com';

  static const _mqttProtocol = 'wss';
  static const mqttServer = "$_mqttProtocol://$_mqttMainUrl:8443";
  static const appUrl = "$pre://$mainUrl";
  static const baseUrl = appUrl;
  static const tradingView = "$appUrl/tv/api/v1/main-route";
  static const jsAPI = "$appUrl/tv/api/v1/js";
  static const urlPrefix = '/api/v1/';
  static const tvUrlPrefix = '/tv/api/v1/js/';

  static String generatemainUrl(String url) {
    return urlPrefix + url;
  }

  static String generatetvUrl(String url) {
    return tvUrlPrefix + url;
  }

  static List<AutoCompleteItem> currencys;

  static List<AutoCompleteItem> currencyArray() {
    if (currencys == null) {
      final currenciesJson = storage.read(StorageKeys.currencies)["currencies"];
      final currencies = List<CurrencyModel>.from(
        currenciesJson.map(
          (model) => CurrencyModel.fromJson(model),
        ),
      );
      List<AutoCompleteItem> autoCompleteArray = [];
      for (var item in currencies) {
        autoCompleteArray.add(
          AutoCompleteItem(
            name: item.code,
            image: item.image,
            desc: item.name,
            code: item.code,
            id: item.id,
            searchPhrase: item.code + '_' + item.name,
            otherNetworks: item.otherBlockChainNetworks,
            mainNetwork: item.mainNetwork,
          ),
        );
      }
      currencys = autoCompleteArray;
    }

    return currencys;
  }

  static List<AutoCompleteItem> countriesArray() {
    final countriesJson = storage.read(StorageKeys.countries);
    final countries = List<CountryModel>.from(
      countriesJson.map(
        (model) => CountryModel.fromJson(model),
      ),
    );
    List<AutoCompleteItem> autoCompleteArray = [];
    for (var item in countries) {
      autoCompleteArray.add(
        AutoCompleteItem(
          name: item.fullName,
          //image: item.image,
          //desc: item.name,
          code: item.code,
          inPerentesis: item.name,
          id: item.id,
        ),
      );
    }
    return autoCompleteArray;
  }

  static List<AutoCompleteItem> pairsAutoCompleteArray() {
    var pairsJson = storage.read(StorageKeys.pairs);
    if (pairsJson == null) {
      print('pairs are not!');
      pairsJson = [];
    }

    List<AutoCompleteItem> autoCompleteArray = [];
    for (var item in pairsJson) {
      autoCompleteArray.add(
        AutoCompleteItem(
          name: item['name'],
          //image: item.image,
          //desc: item.name,
          value: item['value'],
          id: item['id'],
        ),
      );
    }
    return autoCompleteArray;
  }
}
