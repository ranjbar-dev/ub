import 'dart:async';
import 'dart:io';

import 'package:connectivity/connectivity.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_appavailability/flutter_appavailability.dart';
import 'package:get/get.dart';
import 'package:get_storage/get_storage.dart';

import '../../../services/constants.dart';
import '../../../services/storageKeys.dart';
import '../../../utils/commonUtils.dart';
import '../../../utils/environment/ubEnv.dart';
import '../../../utils/logger.dart';
import '../../../utils/mixins/popups.dart';
import '../../../utils/mixins/toast.dart';
import '../../../utils/pairAndCurrencyUtils.dart';
import '../../modules/account/controllers/account_controller.dart';
import '../../modules/exchange/controllers/exchange_controller.dart';
import '../../modules/funds/controllers/funds_controller.dart';
import '../../modules/funds/pages/balance/controllers/balance_controller.dart';
import '../../modules/funds/pages/deposits/controllers/deposits_controller.dart';
import '../../modules/funds/pages/transactionHistory/controllers/transaction_history_controller.dart';
import '../../modules/funds/pages/withdrawals/controllers/withdrawals_controller.dart';
import '../../modules/home/controllers/home_controller.dart';
import '../../modules/identityDocuments/controllers/identity_documents_controller.dart';
import '../../modules/identityInfo/controllers/identity_info_controller.dart';
import '../../modules/landing/app_version_model.dart';
import '../../modules/login/controllers/login_controller.dart';
import '../../modules/market/controllers/market_controller.dart';
import '../../modules/orders/pages/openOrders/controllers/open_orders_controller.dart';
import '../../modules/orders/pages/orderHistory/controllers/order_history_controller.dart';
import '../../modules/trade/controllers/trade_controller.dart';
import '../../routes/app_pages.dart';
import '../currency_pairs_model.dart';
import '../providers/commonDataProvider.dart';
import 'authorizedMqttController.dart';
import 'unAuthorizedMqttController.dart';

enum DeviceTypes { PHONE, TABLET }

class GlobalController extends GetxController with Toaster, Popups {
  var isAppInstalled = false;
  var isRedirectContainerDismissed = false.obs;
  var doShowRedirect = false.obs;
  Map<String, String> appAvailableMap = {};
//this is to initialize pair and con hash maps inside PairAndCurrencyUtils's factory function
  final pairAndCurrencyUtils = PairAndCurrencyUtils();

  ConnectivityResult connectivityResult;

  final currencyPairsArray = [].obs;
  DeviceTypes deviceType = DeviceTypes.PHONE;

  final allCurrencyPairs = <CurrencyPairsModel>[].obs;

  final hasConnection = true.obs;

  bool isLoggingInWithBiometrics = false;

  Map<String, Pairs> pairsHashMap = {};

  StreamSubscription<ConnectivityResult> connectionSubscription;

  final CommonDataProvider commonDataProvider = CommonDataProvider();

  AccountController accountController;

  final GetStorage storage = GetStorage();

  var loggedIn = false.obs;

  bool get isDark => storage.read(StorageKeys.darkMode) ?? false;

  ThemeData get theme => isDark ? ThemeData.dark() : ThemeData.light();

  void enableDarkTheme(bool dark) {
    if (dark == true) {
      storage.write(StorageKeys.darkMode, dark);
      Get.changeThemeMode(ThemeMode.dark);
      return;
    }
    storage.write(StorageKeys.darkMode, dark);
    Get.changeThemeMode(ThemeMode.light);
  }

  void handleLoggedOut({bool andExitApp = false}) async {
    storage.remove(StorageKeys.token);

    storage.remove(StorageKeys.lastLoginDate);

    storage.remove(StorageKeys.favPairs);

    storage.remove(StorageKeys.orderedPairs);

    storage.remove(StorageKeys.refresh);

    loggedIn.value = false;

    _purgeTheMemory();

    if (andExitApp) {
      Future.delayed(500.milliseconds).then((value) {
        exit(0);
      });
      if (GetPlatform.isAndroid) {
        await SystemChannels.platform.invokeMethod<void>(
          'SystemNavigator.pop',
        );
      } else if (GetPlatform.isIOS) {
        await SystemChannels.platform.invokeMethod<void>(
          'SystemNavigator.pop',
        );
      } else if (GetPlatform.isWeb) {
        Get.put(LoginController(), permanent: true);
        Get.offAllNamed(AppPages.LANDING);
        // window.close();
      }
    } else {
      Get.put(LoginController(), permanent: true);
      Get.offAllNamed(AppPages.LANDING);
    }

    return;
  }

  @override
  void onReady() {}

  @override
  void onInit() async {
    super.onInit();
    if (connectionSubscription != null) {
      connectionSubscription.cancel();
    }

    connectionSubscription = Connectivity()
        .onConnectivityChanged
        .listen((ConnectivityResult result) {
      if (result == ConnectivityResult.none) {
        hasConnection.value = false;

        debugPrint('Disconnected!');
      } else {
        hasConnection.value = true;

        debugPrint('Reconnected!');
      }
    });
    connectivityResult = await Connectivity().checkConnectivity();
    if (connectivityResult == ConnectivityResult.none) {
      hasConnection.value = false;
      //don't continue=>commented because api service handles the requests on connection lost
      // return;
    }

    await checkTokenValidation();
    getPairsCurrenciesCountriesAndVersion();

    checkIfRedirectIsNeeded();
    if (GetPlatform.isAndroid) {
      appAvailableMap =
          await AppAvailability.checkAvailability("com.unitedbit.app");
      if (appAvailableMap["app_name"] == 'unitedbit') isAppInstalled = true;
    } else if (GetPlatform.isIOS) {
      appAvailableMap = await AppAvailability.checkAvailability("calshow://");
      if (appAvailableMap["app_name"] == 'unitedbit') isAppInstalled = true;
    }
  }

  Future checkTokenValidation() async {
    if (storage.read(StorageKeys.token) != null) {
      final currentTime = DateTime.now();
      final lastLoginDate =
          DateTime.parse(storage.read(StorageKeys.lastLoginDate));
      final differenceInDays = currentTime.difference(lastLoginDate).inDays;
      if (differenceInDays > 28) {
        //comment top line and uncomment line beloaw to test biometrics when user was away for a long time
        //if (true) {
        Get.put(LoginController(), permanent: true);
        final LoginController loginController = Get.find();
        isLoggingInWithBiometrics = true;
        final canLogin = await loginController.checkForBiometricLogin();
        isLoggingInWithBiometrics = false;
        if (canLogin == false) {
          handleLoggedOut();
        }
        return;
      }
      final passessBiometrics = await canContinueWithBiometrics();
      if (passessBiometrics) {
        loggedIn.value = true;
        loadAuthenticatedControllers();
        final response = await commonDataProvider.getFavoritePairs();
        if (response['status'] == true) {
          storage.write(
            StorageKeys.favPairs,
            response["data"],
          );
        }
      } else {
        handleLoggedOut();
      }
      return;
    }
    return Future.value(true);
  }

  void loadAuthenticatedControllers() {
    Get.put(AuthorizedMqttController(), permanent: true);
    Get.put(UnAuthorizedMqttController(), permanent: true);
    Get.put<TradeController>(TradeController(), permanent: true);
    Get.put<AccountController>(AccountController(), permanent: true);
    accountController = Get.find();
    loggedIn.value = true;
  }

  @override
  void onClose() {
    super.onClose();
    connectionSubscription.cancel();
  }

  void _purgeTheMemory() {
    Future.delayed(100.milliseconds).then((value) {
      Get.delete<AuthorizedMqttController>(force: true);
      Get.delete<UnAuthorizedMqttController>(force: true);
      Get.delete<TradeController>(force: true);
      Get.delete<OrderHistoryController>(force: true);
      Get.delete<OpenOrdersController>(force: true);
      Get.delete<FundsController>(force: true);
      Get.delete<BalanceController>(force: true);
      Get.delete<AccountController>(force: true);
      Get.delete<TransactionHistoryController>(force: true);
      Get.delete<MarketController>(force: true);
      Get.delete<DepositsController>(force: true);
      Get.delete<IdentityDocumentsController>(force: true);
      Get.delete<IdentityInfoController>(force: true);
      Get.delete<WithdrawalsController>(force: true);
      Get.delete<HomeController>(force: true);
      Get.delete<ExchangeController>(force: true);
    });
    Get.delete<LoginController>(force: true);
  }

  Future getVersion() async {
    if (GetPlatform.isWeb != true && ENV != 'DEV') {
      try {
        final response = await commonDataProvider.getAppVerion(
            currentVersion: Constants.appVersion.split("+")[0],
            platform: GetPlatform.isAndroid ? 'android' : 'ios');
        if (response['data'] is List) {
          final List<AppVersionModel> versionList =
              parseVersionList(response['data']);
          final features = <String>[];

          final bugFixes = <String>[];

          bool forceUpdate = false;

          for (var item in versionList) {
            if (item.forceUpdate == true) {
              forceUpdate = true;
            }

            for (var feature in item.keyFeatures ?? []) {
              features.add(feature);
            }

            for (var feature in item.bugFixes ?? []) {
              bugFixes.add(feature);
            }
          }
          if (forceUpdate) {
            storage.remove(StorageKeys.lastCancelUpdate);
          }

          if (features.length > 0 || bugFixes.length > 0 || forceUpdate) {
            final lastCancelUpdate =
                storage.read(StorageKeys.lastCancelUpdate) ??
                    {'date': '', 'version': ''};

            if (lastCancelUpdate['version'] != "" &&
                lastCancelUpdate['version'] != versionList.last.version) {
              openUpdatePopup(
                forceUpdate: forceUpdate,
                features: features,
                bugFixes: bugFixes,
                version: versionList.last.version,
                url: versionList.last.url,
              );
            }
          }
        }
      } catch (e) {
        debugPrint('error getting app version');
        debugPrint(e.toString());
      } finally {}
    }
    return Future.value();
  }

  parseVersionList(data) {
    final list = <AppVersionModel>[];
    for (var item in data) {
      list.add(AppVersionModel.fromJson(item));
    }
    return list;
  }

  void getPairsCurrenciesCountriesAndVersion() async {
    try {
      final commonData = await Future.wait([
        commonDataProvider.getCurrencies(),
        commonDataProvider.getPairs(),
        commonDataProvider.getCountries(),
      ]);
      final coinsResponse = commonData[0];
      if (coinsResponse["data"]["currencies"] != null) {
        storage.write(
          StorageKeys.currencies,
          coinsResponse["data"],
        );

        final jsonMap = {};
        coinsResponse["data"]["currencies"].forEach((model) {
          return jsonMap[model['code']] = model;
        });
        PairAndCurrencyUtils.setCoinsHashMap(jsonMap);
      }

      final pairsResponse = commonData[1];
      if (pairsResponse['status'] == true) {
        final currencyPairs = List<CurrencyPairsModel>.from(
          pairsResponse['data'].map(
            (model) => CurrencyPairsModel.fromJson(model),
          ),
        );
        allCurrencyPairs.assignAll(currencyPairs);
        final tmp = [];
        final Map<String, dynamic> hashMapTmpJson = {};
        final Map<String, dynamic> currencyPairsHashMapTmpJson = {};
        currencyPairs.forEach(
          (item) => {
            currencyPairsHashMapTmpJson.putIfAbsent(
              item.code,
              () => item.toJson(),
            ),
            item.pairs.forEach((pair) {
              pairsHashMap[pair.pairName] = pair;
              hashMapTmpJson[pair.pairName] = pair.toJson();
              return tmp.add({
                "name": pair.pairName,
                "value": pair.pairName,
                "id": pair.pairId,
                "dependentName": pair.dependentName,
              });
            })
          },
        );

        storage.write(StorageKeys.pairs, tmp);
        currencyPairsArray.assignAll(tmp);
        PairAndCurrencyUtils.setCurrencyPairsHashMap(
            currencyPairsHashMapTmpJson);
        PairAndCurrencyUtils.setPairsHashMap(hashMapTmpJson);
      }

      final countriesResponse = commonData[2];
      if (countriesResponse["status"] == true) {
        storage.write(
          StorageKeys.countries,
          countriesResponse["data"],
        );
      }
    } catch (e) {
      log.e(e.toString());
    } finally {
      //Future.delayed(5.seconds).then((value) {
      //  getVersion();
      //});
    }
  }

  void checkIfRedirectIsNeeded() {
    if ((GetPlatform.isAndroid || GetPlatform.isIOS) &&
        kIsWeb &&
        !isRedirectContainerDismissed.value) {
      doShowRedirect.value = true;
    } else {
      doShowRedirect.value = false;
    }
  }

  void setRedirectConteinerDismissed(bool value) {
    isRedirectContainerDismissed.value = value;
    checkIfRedirectIsNeeded();
  }

  void doRedirect() async {
    // if (isAppInstalled)
    //   AppAvailability.launchApp(
    //       "com.unitedbit.app");
    // else {
    final response = await commonDataProvider.getAppVerion(
        currentVersion: Constants.appVersion.split("+")[0],
        platform: GetPlatform.isAndroid ? 'android' : 'ios');
    if (response['data'] is List) {
      final List<AppVersionModel> versionList =
          parseVersionList(response['data']);
      launchURL(versionList.last.url);
    }
    /*StoreRedirect.redirect();*/
  }

  setDeviceType() {
    final data = MediaQueryData.fromWindow(WidgetsBinding.instance.window);
    data.size.shortestSide < 600
        ? deviceType = DeviceTypes.PHONE
        : deviceType = DeviceTypes.TABLET;
  }
}
