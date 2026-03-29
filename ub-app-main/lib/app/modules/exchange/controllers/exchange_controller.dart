import 'dart:async';

import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:get_storage/get_storage.dart';
import 'package:unitedbit/app/global/controller/authorizedMqttController.dart';
import '../../orders/pages/orderHistory/providers/order_history_provider.dart';

import '../../../../services/constants.dart';
import '../../../../services/storageKeys.dart';
import '../../../../utils/commonUtils.dart';
import '../../../../utils/extentions/basic.dart';
import '../../../../utils/logger.dart';
import '../../../../utils/marketUtils.dart';
import '../../../../utils/mixins/popups.dart';
import '../../../../utils/mixins/toast.dart';
import '../../../../utils/numUtil.dart';
import '../../../../utils/pairAndCurrencyUtils.dart';
import '../../../../utils/throttle.dart';
import '../../../global/autocompleteModel.dart';
import '../../../global/controller/globalController.dart';
import '../../../global/currency_pairs_model.dart';
import '../../home/home_page_pair_price_model.dart';
import '../../trade/controllers/trade_controller.dart';
import '../../trade/new_trade_order_model.dart';
import '../../trade/pair_balance_model.dart';
import '../models/pairLocalInfoModel.dart';
import '../providers/exchangeProvider.dart';

class ExchangeController extends GetxController with Popups, Toaster {
  //Objects To Provide Data
  final GlobalController globalController = Get.find();
  final TradeController tradeController = Get.find();
  final AuthorizedMqttController authOrderController = Get.find();

  final inputValidationToastThrottle = new Throttling(
    duration: const Duration(milliseconds: 4000),
  );
  final sparkLinePairs = <HomePagePairPriceModel>[].obs;
  final isLoadingSparkLine = false.obs;

  final totalValue = ''.obs;
  final tradeFee = ''.obs;
  final thr1500 = new Throttling(
    duration: const Duration(milliseconds: 1500),
  );
  final thr1000 = new Throttling(
    duration: const Duration(milliseconds: 1000),
  );
  final exchangeProvider = ExchangeProvider();
  final OrderHistoryProvider orderHistoryProvider = OrderHistoryProvider();
  final coinsList = Constants.currencyArray();
  final storage = GetStorage();

  final Rx<TextEditingController> inputControllerFrom =
      TextEditingController(text: '0.0000').obs;
  final Rx<TextEditingController> inputControllerTo =
      TextEditingController(text: '0.0000').obs;

  PairLocalInfoModel pairLocalInfo = PairLocalInfoModel(
      activePairID: 1.obs,
      activePairName: 'BTC-USDT'.obs,
      pairPrecision: 8.obs,
      type: 'sell'.obs,
      basisCoin: AutoCompleteItem(
        name: 'BTC-USDT',
        code: 'BTC',
        desc: 'Bitcoin',
        image: PairAndCurrencyUtils.findCoinImageByCode('BTC'),
      ).obs,
      dependantCoin: AutoCompleteItem(
        name: 'BTC-USDT',
        code: 'USDT',
        desc: 'Tether',
        image: PairAndCurrencyUtils.findCoinImageByCode('USDT'),
      ).obs,
      basisBalance: 0.0000.obs,
      possiblePairs: [
        AutoCompleteItem(
          name: 'BTC-USDT',
          code: 'USDT',
          desc: 'Tether',
          image: PairAndCurrencyUtils.findCoinImageByCode('USDT'),
        )
      ].obs,
      dependentBalance: 0.0000.obs);

  final savedCoins = <AutoCompleteItem>[].obs;
  final pairBalanceData = PairBalanceModel().obs;
  final isLoadingBalanceData = false.obs;
  final isLoadingExchangeSubmit = false.obs;
  Map<String, dynamic> currencyPairsInfo;

  Map<String, Pairs> pairsHashMap = {};

  final possiblePairs = <Pairs>[].obs;
  fillPairHashMap() {
    pairsHashMap = Map.from(PairAndCurrencyUtils.pairsMap.value);
  }

  var allPairs = PairAndCurrencyUtils.currencyPairsMap.value;

  @override
  void onInit() async {
    final List<dynamic> storedCoins =
        storage.read<List>(StorageKeys.savedDepositCoins);
    if (storedCoins != null) {
      savedCoins.assignAll(
        storedCoins.map((e) => AutoCompleteItem.fromJson(e)).toList(),
      );
    }

    filterPairsBasedOnCoin(pairLocalInfo.activePairName.value.split('-').first);
    getPairBalances();

    final Map<String, dynamic> storedPair =
        storage.read(StorageKeys.selectedPair) ?? {"id": 1, "name": 'BTC-USDT'};

    if (PairAndCurrencyUtils.pairsMap.value != null) {
      fillPairHashMap();
    } else {
      PairAndCurrencyUtils.pairsMap.listen((v) => {fillPairHashMap()});
    }

    getPriceStream(isOnInit: true);
    //getOrderStream();
    getPairPriceSparkLineChartData();

    pairLocalInfo.activePairName.value = storedPair["name"];
    pairLocalInfo.activePairID.value = storedPair["id"];

    super.onInit();
  }

  @override
  void onReady() {
    super.onReady();
  }

  @override
  void onClose() {
    inputControllerFrom.value.dispose();
    inputControllerTo.value.dispose();
    super.onClose();
  }

  FutureOr<PairBalanceModel> parseBalanceData(message) {
    return PairBalanceModel.fromJson(message);
  }

  void getOrderStream() {
    authOrderController.ordrPayload.listen((order) {
      log.e('FILLED ORDER' + order.toString());
    });
  }

  void getPriceStream({bool isOnInit}) {
    isOnInit = isOnInit ?? false;
    final List<Pairs> possibleTmp = [];

    if (isOnInit) {
      pairLocalInfo.possiblePairs.forEach((element) {
        possibleTmp.add(
          Pairs(
            pairId: element.id,
            pairName: element.name,
            price: '0.0',
            percent: '0.0',
            volume: '--',
            equivalentPrice: '--',
            formattedEquivalentPrice: '--',
            formattedPrice: '--',
            formattedVolume: '--',
          ),
        );
      });
      possiblePairs.assignAll(possibleTmp);

      tradeController.lastPrice.listen((lastPrice) {
        final index = sparkLinePairs
            .indexWhere((element) => element.pairName == lastPrice.name);
        if (index != -1 && sparkLinePairs.length > 0) {
          final sparkLineData = sparkLinePairs[index].trendData;

          sparkLineData.last.price = lastPrice.price;

          sparkLineData.last.change = lastPrice.percentage;

          if (!(sparkLineData.last.change.contains('-'))) {
            sparkLineData.last.change = '+' + sparkLineData.last.change;
          }
          sparkLinePairs[index].trendData = sparkLineData;
          if (!(GetPlatform.isWeb)) {
            sparkLinePairs.refresh();
          } else {
            thr1000.throttle(() {
              sparkLinePairs.refresh();
            });
          }
        }

        final possiblePairIndex = possiblePairs
            .indexWhere((element) => element.pairName == lastPrice.name);
        if (possiblePairIndex != -1 && pairsHashMap[lastPrice.name] != null) {
          final json = MarketUtils.priceJson(
            lastPrice,
          );

          final newPrice = Pairs.fromJson(json);

          // ignore: invalid_use_of_protected_member
          possiblePairs.value[possiblePairIndex] = newPrice;
          if (!(GetPlatform.isWeb)) {
            possiblePairs.refresh();
          } else {
            thr1500.throttle(() {
              possiblePairs.refresh();
            });
          }
          calcHowMuchYouWillGet();
        }
      });
    } else {
      pairLocalInfo.possiblePairs.forEach((element) {
        possibleTmp.add(
          Pairs(
            pairId: element.id,
            pairName: element.name,
            price: '0.0',
            percent: '0.0',
            volume: '--',
            equivalentPrice: '--',
            formattedEquivalentPrice: '--',
            formattedPrice: '--',
            formattedVolume: '--',
          ),
        );
      });

      possiblePairs.assignAll(possibleTmp);

      tradeController.lastPrice.listen((lastPrice) {
        final possiblePairIndex = possiblePairs
            .indexWhere((element) => element.pairName == lastPrice.name);
        if (possiblePairIndex != -1 && pairsHashMap[lastPrice.name] != null) {
          final json = MarketUtils.priceJson(
            lastPrice,
          );

          final newPrice = Pairs.fromJson(json);

          // ignore: invalid_use_of_protected_member
          possiblePairs.value[possiblePairIndex] = newPrice;
          if (!(GetPlatform.isWeb)) {
            possiblePairs.refresh();
          } else {
            thr1500.throttle(() {
              possiblePairs.refresh();
            });
          }
          calcHowMuchYouWillGet();
        }
      });
    }
  }

  Future getPairBalances({int pairId}) async {
    final id = pairId ?? getActivePairId();
    if (id != null) {
      try {
        isLoadingBalanceData.value = true;
        final response =
            await exchangeProvider.getCurrencyPairDetails(pairId: id);
        if (response["status"] == true) {
          pairBalanceData.value = await parseBalanceData(response['data']);
          // pairBalanceData.value =
          //     await compute(parseBalanceData, response['data']);
          for (var item in pairBalanceData.value.pairBalances) {
            if (item.currencyCode == pairLocalInfo.basisCoin.value.code) {
              pairLocalInfo.basisBalance.value = double.tryParse(item.balance);
            } else {
              pairLocalInfo.dependentBalance.value =
                  double.tryParse(item.balance);
            }
          }
        }
      } on DioError catch (e) {
        toastDioError(e);
      } finally {
        isLoadingBalanceData.value = false;
      }
      return Future.value(true);
    }
  }

  void filterPairsBasedOnCoin(String code) async {
    pairLocalInfo.possiblePairs.clear();
    allPairs.forEach((key, value) {
      if (key == code) {
        for (var pair in value.pairs) {
          pairLocalInfo
            ..possiblePairs.add(
              AutoCompleteItem(
                  id: pair.pairId,
                  name: pair.pairName,
                  code: pair.dependentCode,
                  desc: PairAndCurrencyUtils.findCoinNameByCode(
                      pair.dependentCode),
                  image: pair.image,
                  value: 'buy'),
            );
        }
      } else {
        for (var pair in value.pairs) {
          if (pair.dependentCode == code)
            pairLocalInfo
              ..possiblePairs.add(
                AutoCompleteItem(
                    id: pair.pairId,
                    name: pair.pairName,
                    code: pair.basisCode,
                    desc:
                        PairAndCurrencyUtils.findCoinNameByCode(pair.basisCode),
                    image: PairAndCurrencyUtils.findCoinImageByCode(
                        pair.basisCode),
                    value: 'sell'),
              );
        }
      }
    });

    if (pairLocalInfo.possiblePairs.isEmpty) {
      pairLocalInfo.possiblePairs.add(
        (AutoCompleteItem(name: '-', desc: 'No Pair Found', image: '')),
      );
    }
  }

  int getActivePairId() {
    return pairLocalInfo.activePairID.value ?? 1;
  }

  void calcHowMuchYouWillGet() {
    double basisAmount = double.tryParse(inputControllerFrom.value.text) ?? 0;
    if (basisAmount == 0) {
      inputControllerTo.value.text = '0.0000';
      return;
    }
    int idx;
    idx = possiblePairs.indexWhere(
        (element) => element.pairId == pairLocalInfo.activePairID.value);

    double total;
    if (pairLocalInfo.type.value == "sell") {
      total = NumUtil.multiply(
        basisAmount,
        double.tryParse(possiblePairs[idx].price),
      );
    } else {
      total = NumUtil.divide(
        basisAmount,
        double.tryParse(possiblePairs[idx].price),
      );
    }
    totalValue.value = total.parseDoubleToString();

    final fee = pairBalanceData.value.fee.makerFee;
    final newTradeFee = NumUtil.multiply(total, fee);
    final newYouGet = NumUtil.subtract(total, newTradeFee);
    inputControllerTo.value.text = newYouGet.toString();
  }

  void handleExchangeClick() async {
    final canSubmit = _canSubmit();
    if (canSubmit == true) {
      final model = NewTradeOrderModel()
        ..pairCurrencyId = pairLocalInfo.activePairID.value
        ..type = pairLocalInfo.type.value
        ..exchangeType = 'market'
        ..amount = inputControllerFrom.value.text
        ..isFastExchange = true;

      try {
        isLoadingExchangeSubmit.value = true;
        final response = await exchangeProvider.createExchange(model: model);
        if (response['status'] == true) {
          final orderID = response['data']['id'];
          await orderHistoryProvider.getOrderDetails(orderID).then(
                (value) => {
                  isLoadingExchangeSubmit.value = false,
                  openExchangeSubmitPopup(
                      spentAmount:
                          double.tryParse(inputControllerFrom.value.text),
                      coinCode: pairLocalInfo.basisCoin.value.code,
                      model: value['data'][0],
                      backTapped: () {
                        _resetAllInputs();
                        getPairBalances(
                            pairId: pairLocalInfo.activePairID.value);
                        Get.back();
                      })
                },
              );
        }
      } on DioError catch (e) {
        toastDioError(e);
      } finally {
        isLoadingExchangeSubmit.value = false;
      }
    }
  }

  bool _canSubmit() {
    String message;
    if (double.tryParse(inputControllerFrom.value.text) >
        pairLocalInfo.basisBalance.value) {
      message = 'You have exeeded the available amount';
    }
    if (message != null) {
      inputValidationToastThrottle.throttle(() => toastWarning(message));
      return false;
    }
    return true;
  }

  _resetAllInputs() {
    inputControllerFrom.value.text = '0.0000';
    inputControllerTo.value.text = '0.0000';
  }

  Future getPairPriceSparkLineChartData() async {
    isLoadingSparkLine.value = true;
    try {
      final response = await exchangeProvider.getPairsPrice();
      // final parsed = await compute(
      //     parseHomePagePairsPrice, response['data']['pairs'][0]['pairs']);
      final parsed =
          parseHomePagePairsPrice(response['data']['pairs'][0]['pairs']);
      sparkLinePairs.assignAll(parsed);
    } catch (e) {
      log.e('homePageError !, Error getting sparkline data');
      log.e(e.toString());
    } finally {
      isLoadingSparkLine.value = false;
    }
    return Future.value();
  }

  List<HomePagePairPriceModel> parseHomePagePairsPrice(response) {
    final list = List<HomePagePairPriceModel>.from(
      response.map(
        (model) => HomePagePairPriceModel.fromJson(model),
      ),
    );
    return list;
  }

  Future handleAutoExchangeDependantCoinSelected(
      {AutoCompleteItem coin}) async {
    // var currentPair = pairLocalInfo.possiblePairs
    //     .firstWhere((element) => element.code == autoExchangeCode);

    pairLocalInfo.dependantCoin.value = coin;
    pairLocalInfo.activePairID.value = coin.id;
    pairLocalInfo.activePairName.value = coin.name;
    pairLocalInfo.type.value = pairLocalInfo.dependantCoin.value.value;
  }

  Future handleAutoExchangeBaseCoinSelected(
      {AutoCompleteItem coin, String autoExchangeCode}) async {
    pairLocalInfo.basisCoin.value = coin;
    filterPairsBasedOnCoin(coin.code);

    if (!["", null].contains(autoExchangeCode)) {
      var currentExchange = pairLocalInfo.possiblePairs
          .firstWhere((element) => element.code == autoExchangeCode);
      if (!["", null].contains(currentExchange.code)) {
        pairLocalInfo.possiblePairs
            .removeWhere((element) => element.code == currentExchange.code);
        pairLocalInfo.possiblePairs.insert(0, currentExchange);
      }
    }
    pairLocalInfo.activePairID.value = pairLocalInfo.possiblePairs[0].id;
    pairLocalInfo.activePairName.value = pairLocalInfo.possiblePairs[0].name;
    pairLocalInfo.type.value = pairLocalInfo.possiblePairs[0].value;
    pairLocalInfo.dependantCoin.value = pairLocalInfo.possiblePairs[0];
  }

  Future handleCoinSelected(
      {AutoCompleteItem coin,
      bool isFrom,
      bool isSwap,
      bool isTopMarket}) async {
    isSwap = isSwap ?? false;
    isTopMarket = isTopMarket ?? false;
    _resetAllInputs();
    if (isTopMarket) {
      filterPairsBasedOnCoin(coin.code);
    } else if (isSwap) {
      pairLocalInfo.dependantCoin.value = pairLocalInfo.basisCoin.value;
      pairLocalInfo.basisCoin.value = coin;
      filterPairsBasedOnCoin(coin.code);
      if (pairLocalInfo.type.value == 'sell') {
        pairLocalInfo.type.value = 'buy';
      } else {
        pairLocalInfo.type.value = 'sell';
      }
    } else {
      if (isFrom) {
        pairLocalInfo.basisCoin.value = coin;
        filterPairsBasedOnCoin(coin.code);
        pairLocalInfo.activePairID.value = pairLocalInfo.possiblePairs[0].id;
        pairLocalInfo.activePairName.value =
            pairLocalInfo.possiblePairs[0].name;
        pairLocalInfo.type.value = pairLocalInfo.possiblePairs[0].value;
        pairLocalInfo.dependantCoin.value = pairLocalInfo.possiblePairs[0];
      } else {
        pairLocalInfo.dependantCoin.value = coin;
        pairLocalInfo.activePairID.value = coin.id;
        pairLocalInfo.activePairName.value = coin.name;
        pairLocalInfo.type.value = pairLocalInfo.dependantCoin.value.value;
      }
    }

    getPairBalances(pairId: pairLocalInfo.activePairID.value);
    getPriceStream();

    saveCoinToHistory(
      coin: coin,
      storageKey: StorageKeys.savedDepositCoins,
      stream: savedCoins,
    );
  }

  void handlePairSelect(HomePagePairPriceModel selectedPair) {
    String basisCode = selectedPair.pairName.split('-').first;
    String dependentCode = selectedPair.pairName.split('-').last;

    pairLocalInfo
      ..activePairID.value = selectedPair.pairId
      ..activePairName.value = selectedPair.pairName
      ..basisCoin.value = AutoCompleteItem(
        name: selectedPair.pairName,
        code: basisCode,
        desc: PairAndCurrencyUtils.findCoinNameByCode(basisCode),
        image: PairAndCurrencyUtils.findCoinImageByCode(basisCode),
      )
      ..dependantCoin.value = AutoCompleteItem(
        name: selectedPair.pairName,
        code: dependentCode,
        desc: PairAndCurrencyUtils.findCoinNameByCode(dependentCode),
        image: PairAndCurrencyUtils.findCoinImageByCode(dependentCode),
      );

    handleCoinSelected(coin: pairLocalInfo.basisCoin.value, isTopMarket: true);

    if (allPairs.containsKey(basisCode)) {
      for (var pair in allPairs[basisCode].pairs) {
        if (pair.dependentCode == dependentCode) {
          pairLocalInfo.type.value = 'buy';
          return;
        }
      }
    }
    if (allPairs.containsKey(dependentCode)) {
      for (var pair in allPairs[dependentCode].pairs) {
        if (pair.dependentCode == basisCode) {
          pairLocalInfo.type.value = 'sell';
        }
      }
    }
  }

  void allIn() {
    inputControllerFrom.value.text = pairLocalInfo.basisBalance.toString();
    calcHowMuchYouWillGet();
  }

  void swapCoins() {
    _resetAllInputs();
    handleCoinSelected(
        coin: pairLocalInfo.dependantCoin.value, isFrom: false, isSwap: true);
  }
}
