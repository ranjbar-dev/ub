import 'dart:async';
import 'dart:convert' show jsonDecode;

import 'package:basic_utils/basic_utils.dart';
import 'package:dio/dio.dart' show DioError;
import 'package:flutter/foundation.dart';
import 'package:get/get.dart';
import 'package:get_storage/get_storage.dart';
import 'package:unitedbit/mqttClient/universal_mqtt_client.dart'
    show UniversalMqttClient, UniversalMqttClientStatus, MqttQos;
import 'package:unitedbit/services/constants.dart';
import 'package:unitedbit/services/storageKeys.dart';
import 'package:unitedbit/utils/commonUtils.dart';
import 'package:unitedbit/utils/extentions/basic.dart';
import 'package:unitedbit/utils/logger.dart';
import 'package:unitedbit/utils/mixins/toast.dart';
import 'package:unitedbit/utils/numUtil.dart';
import 'package:unitedbit/utils/throttle.dart';
import 'package:uuid/uuid.dart';

import '../../../../generated/locales.g.dart';
import '../../../common/components/UBButton.dart';
import '../../../common/components/UBWrappedButtons.dart';
import '../../../global/autocompleteModel.dart';
import '../../../global/controller/authorizedMqttController.dart';
import '../../../global/controller/globalController.dart';
import '../../../global/controller/unAuthorizedMqttController.dart';
import '../../account/controllers/account_controller.dart';
import '../../orders/controllers/orders_controller.dart';
import '../../orders/pages/openOrders/controllers/open_orders_controller.dart';
import '../../orders/pages/orderHistory/controllers/order_history_controller.dart';
import '../labelModel.dart';
import '../models/price_model.dart';
import '../new_trade_order_model.dart';
import '../pair_balance_model.dart';
import '../provider/tradeProvider.dart';
import 'ohlcChart_controller.dart';

// calculateTimeFrame(int interval) {
//   switch (interval) {
//     case 1:
//     case 3:
//       return '1minute';
//     case 5:
//     case 15:
//     case 30:
//     case 45:
//       return '5minutes';
//     case 60:
//     case 120:
//     case 180:
//     case 240:
//       return '1hour';
//     case 60 * 24:
//     case 60 * 24 * 7:
//     case 60 * 24 * 30:
//       return '1day';

//     default:
//       return '1minute';
//   }
// }
enum TradeTopCharts { OHLC, OrderBook }
final uuid = Uuid();

class TradeController extends GetxController with Toaster {
  final balanceToastThrottle =
      new Throttling(duration: const Duration(milliseconds: 4000));
  final lastOhlcThrottle =
      new Throttling(duration: const Duration(milliseconds: 4000));
  final inputValidationToastThrottle =
      new Throttling(duration: const Duration(milliseconds: 4000));

  final GlobalController globalController = Get.find();
  final UnAuthorizedMqttController unAuthorizedMqttController = Get.find();
  final AuthorizedMqttController authorizedMqttController = Get.find();

  final log = UBLogger.log;
  final timeFrameButtons = [
    WrappedButtonModel(text: '1 Minute', value: '1minute'),
    WrappedButtonModel(text: '5 Minute', value: '5minutes'),
    WrappedButtonModel(text: '1 Hour', value: '1hour'),
    WrappedButtonModel(text: '1 Day', value: '1day'),
  ];

  final GetStorage storage = GetStorage();

  final tradeProvider = TradeProvider();

  final orderBookData = {}.obs;

  OHLCChartController chartController;

  StreamSubscription<String> priceSubscription;

  bool priceTopicInitialized = false;

  StreamSubscription<String> ohlcSubscription;

  bool ohlcTopicInitialized = false;

  StreamSubscription<String> orderbookSubscription;

  bool orderBookTopicInitialized = false;

  final activeChart = (TradeTopCharts.OHLC).obs;

  final RxString currentPairName = 'BTC-USDT'.obs;

  String equivalentAmoutFroForstCoinInPair = '0.0';

  final pairBalanceData = PairBalanceModel().obs;

  final lastOhlcValue = {}.obs;

  var pairs = [].obs;

  final isLoadingPairBalance = false.obs;

  final isCreatingOrder = false.obs;

  final mainActiveIndex = 0.obs;

  final subActiveIndex = 0.obs;

  final numberOfPercentSegments = 4;

  final selectedPercentIndex = (-1).obs;

  final totalValue = ''.obs;

  final amountValue = ''.obs;

  final priceValue = ''.obs;

  final stopValue = ''.obs;

  final tradeFee = ''.obs;

  final youGet = ''.obs;

  final selectedTimeFrame = '1hour'.obs;

  final amountInputLabel = LabelModel(
    placeHolder: "${LocaleKeys.amountTo.tr} ${LocaleKeys.buy.tr}",
    endLabel: "BTC",
  ).obs;

  final priceInputLabel = LabelModel(
    placeHolder: "${LocaleKeys.ifPriceRisesTo.tr}",
    endLabel: "BTC",
  ).obs;

  final totalInputLabel = LabelModel(
    placeHolder: "${LocaleKeys.total.tr}",
    endLabel: "USDT",
  ).obs;
  final stopPriceInputLabel = LabelModel(
    placeHolder: "${LocaleKeys.amountTo.tr} ${LocaleKeys.buy.tr}",
    endLabel: "BTC",
  ).obs;

  final currentPairPrice = PriceModel().obs;

  final lastPrice = PriceModel().obs;

  final priceArray = <PriceModel>[].obs;

  final showLoadingOverlay = false.obs;

  UniversalMqttClient unAuthorizedClient;

  String get _coinName1 => currentPairName.split('-')[1];

  String get _coinName0 => currentPairName.split('-')[0];

  LightSubscription<List<RxUpdateables>> updateSubscription;

  @override
  void onInit() async {
    unAuthorizedClient = unAuthorizedMqttController.unAuthorizedClient;
    updateSubscription =
        authorizedMqttController.updateDataSubject.listen((value) {
      if (value is List &&
          (value.indexOf(RxUpdateables.UserPairBalances) != -1)) {
        getPairBalances();
      }

      return;
    });

    Get.put(OrdersController(), permanent: true);

    Get.put(OpenOrdersController(), permanent: true);

    Get.put(OrderHistoryController(), permanent: true);

    final Map<String, dynamic> storedPair =
        storage.read(StorageKeys.selectedPair) ?? {"id": 1, "name": 'BTC-USDT'};

    currentPairName.value = storedPair["name"];

    selectedTimeFrame.value =
        storage.read(StorageKeys.selectedTimeFrame) ?? '1hour';

    Future.wait([
      _connectToUnAuthorizedMqtt(),
      //getPairBalances(pairId: storedPair["id"] ?? 1),
    ]);
    storage.write(StorageKeys.loggedInOnce, true);
    Future.delayed(100.milliseconds).then((value) {
      generateLabels(mainIndex: 0, subIndex: 0, pair: currentPairName.value);
    });

    super.onInit();
  }

  @override
  void onReady() {
    super.onReady();
  }

  @override
  void onClose() {
    if (priceSubscription != null) {
      priceSubscription.cancel();
    }
    if (ohlcSubscription != null) {
      ohlcSubscription.cancel();
    }
    if (orderbookSubscription != null) {
      orderbookSubscription.cancel();
    }
    if (updateSubscription != null) {
      updateSubscription.cancel();
    }
  }

  void handleMainTabChange(int index) {
    _resetAllInputs();
    generateLabels(
      mainIndex: index,
      subIndex: subActiveIndex.value,
      pair: currentPairName.value,
    );
    selectedPercentIndex.value = -1;
    mainActiveIndex.value = index;
  }

  void handleSubTabChange(int index) {
    _resetPercent();
    _resetAllInputs();
    generateLabels(
      mainIndex: mainActiveIndex.value,
      subIndex: index,
      pair: currentPairName.value,
    );
    youGet.value = '';
    tradeFee.value = '';
    subActiveIndex.value = index;
  }

  String get _getBalance =>
      pairBalanceData.value.pairBalances[mainActiveIndex.value].balance;

  void handlePercentClick({int index}) {
    String balance;
    if (subActiveIndex.value == 1) {
      balance =
          pairBalanceData.value.pairBalances[mainActiveIndex.value].balance;
    } else {
      balance = _getBalance;
    }
    //calculate valid amount to buy for example BTC in btc-usdt pair when we have some usdt
    if (mainActiveIndex.value == 0 &&
        subActiveIndex.value == 0 &&
        currentPairPrice.value.price != null &&
        pairBalanceData.value.pairBalances[0].balance.toDouble() > 0) {
      final currentPrice = currentPairPrice.value.price;
      balance = NumUtil.divide(balance.toDouble(), currentPrice.toDouble())
          .toStringAsFixed(8);
    }

    if (balance.toDouble() == 0) {
      _toastToDeposit();
      return;
    }

    if (pairBalanceData.value.sum != null) {
      if (index == selectedPercentIndex.value) {
        amountValue.value = ('');
        totalValue.value = '';
        selectedPercentIndex.value = -1;
        youGet.value = '';
        tradeFee.value = '';
        return;
      }
      selectedPercentIndex.value = index;
      final amountNumber =
          (balance.toDouble() * ((index + 1) / numberOfPercentSegments));
      final stringAmount = amountNumber.toString();
      handleAmountChange(stringAmount, fromInput: false);
//update youget and trade fee when subtab is set to market
      if (subActiveIndex.value == 1) {
        updateMarketYouGet();
      }
    }
  }

  void handleSubmitClick() async {
    final balance = _getBalance;
    if (balance.toDouble() == 0) {
      _toastToDeposit();
      return;
    }
    final canSubmit = _validateNewOrderInputs();
    if (canSubmit == true) {
      int pairId = getActivePairId();
      final amountV = amountValue.value.removeComma();
      var priceV = priceValue.value.removeComma();
      var stopPriceV = stopValue.value.removeComma();
      final type = mainActiveIndex.value == 0 ? 'buy' : 'sell';
      final exchangeType = subActiveIndex.value == 0
          ? 'limit'
          : subActiveIndex.value == 1
              ? 'market'
              : 'limit';
      if (subActiveIndex.value == 1) {
        priceV = null;
        stopPriceV = null;
      }
      if (subActiveIndex.value != 2) {
        stopPriceV = null;
      }
      final model = NewTradeOrderModel();
      model.type = type;
      model.exchangeType = exchangeType;
      model.amount = amountV;
      model.pairCurrencyId = pairId;
      model.price = priceV;
      model.stopPointPrice = stopPriceV;
      try {
        isCreatingOrder.value = true;
        final response = await tradeProvider.createOrder(model: model);
        if (response['status'] == true) {
          _resetAllInputs();
          _resetPercent();
          getPairBalances(pairId: pairId);
        }
      } on DioError catch (e) {
        toastDioError(e);
      } finally {
        isCreatingOrder.value = false;
      }
    }
  }

  void handlePairChange(String pairName) async {
    if (pairName != currentPairName.value) {
      _resetAllInputs();
      _resetPercent();
      generateLabels(
        mainIndex: mainActiveIndex.value,
        subIndex: subActiveIndex.value,
        pair: pairName,
      );
      purgeTopic(
        client: unAuthorizedClient,
        topicStream: orderbookSubscription,
        topic: _orderBookTopic,
      );
      purgeTopic(
        client: unAuthorizedClient,
        topicStream: ohlcSubscription,
        topic: _ohlcTopic,
      );

      orderBookData.value = {};

      currentPairPrice.value = PriceModel();
      currentPairName.value = pairName;
      int pairId = getActivePairId();
      storage.write(StorageKeys.selectedPair, {"id": pairId, "name": pairName});

      await Future.wait([
        getPairBalances(pairId: pairId),
        _connectToOrderBook(),
        _connectToOHLC()
      ]);
    }
  }

  handeTimeFrameChange(String newTimeFrame) async {
    final idx =
        timeFrameButtons.indexWhere((element) => element.value == newTimeFrame);
    if (idx == -1) {
      debugPrint('timeFrame is not valid !!!!!!!!!');
      return;
    }
    if (ohlcSubscription != null) {
      ohlcSubscription.cancel();
    }
    if (ohlcTopicInitialized) {
      unAuthorizedClient.unsubscribe(topic: _ohlcTopic);
    }
    storage.write(StorageKeys.selectedTimeFrame, newTimeFrame);
    selectedTimeFrame.value = newTimeFrame;
    await Future.wait([_connectToOHLC()]);
  }

  void generateLabels({int mainIndex, int subIndex, String pair}) {
    final pair0 = pair.split('-')[0];
    final pair1 = pair.split('-')[1];
    //total label
    final totalL = LabelModel(
      placeHolder: "${LocaleKeys.total.tr}",
      endLabel: pair1,
    );
    totalInputLabel.value = totalL;
    if (subIndex == 0 || subIndex == 2) {
      ///////////////////////////
      final amountL = LabelModel(
        placeHolder: LocaleKeys.amount.tr,
        //"${LocaleKeys.amountTo.tr} ${mainIndex == 0 ? LocaleKeys.buy.tr : LocaleKeys.sell.tr}",
        endLabel: pair0,
      );
      amountInputLabel.value = amountL;
      final priceL = LabelModel(
        placeHolder: LocaleKeys.price.tr,
        //"${mainIndex == 0 ? LocaleKeys.ifPriceRisesTo.tr : LocaleKeys.ifPriceDropsTo.tr}",
        endLabel: pair1,
      );
      if (subIndex == 2) {
        priceL.placeHolder =
            mainIndex == 0 ? LocaleKeys.buyAt.tr : LocaleKeys.sellAt.tr;
      }
      priceInputLabel.value = priceL;

      final stopL = LabelModel(
        placeHolder: LocaleKeys.ifPriceReaches.tr,
        // "${mainIndex == 0 ? LocaleKeys.buyAtMostPrice.tr : LocaleKeys.sellAtLeastPrice.tr}",
        endLabel: pair1,
      );

      stopPriceInputLabel.value = stopL;
      //////////////////////////////
    } else if (subIndex == 1) {
      final amountL = LabelModel(
        placeHolder: LocaleKeys.amount.tr,
        //"${LocaleKeys.amountTo.tr} ${mainIndex == 0 ? LocaleKeys.buy.tr : LocaleKeys.sell.tr}",
        endLabel: mainIndex == 0 ? pair1 : pair0,
      );
      amountInputLabel.value = amountL;
    }
  }

  Future getPairBalances({int pairId, bool silent}) async {
    final id = pairId ?? getActivePairId() ?? 1;
    if (id != null) {
      try {
        // if (silent != true) {
        isLoadingPairBalance.value = true;
        // }
        final response = await tradeProvider.getCurrencyPairDetails(pairId: id);
        if (response["status"] == true) {
          isLoadingPairBalance.value = false;
          pairBalanceData.value =
              await compute(parseBalanceData, response['data']);
        }
      } on DioError catch (e) {
        toastDioError(e);
      } finally {
        // if (silent != true) {
        isLoadingPairBalance.value = false;
        // }
      }
      return Future.value(true);
    }
  }

  Future _connectToUnAuthorizedMqtt() async {
    pairs = globalController.currencyPairsArray;
    unAuthorizedClient.status.listen(
      (status) {
        if (status == UniversalMqttClientStatus.disconnected) {
          purgeTopic(
            client: unAuthorizedClient,
            topicStream: orderbookSubscription,
            topic: _orderBookTopic,
          );
          purgeTopic(
            client: unAuthorizedClient,
            topicStream: ohlcSubscription,
            topic: _ohlcTopic,
          );
          purgeTopic(
            client: unAuthorizedClient,
            topicStream: priceSubscription,
            topic: _priceTopic,
          );
        }
        log.i('UnAuthorized connection Status: $status');
        if (status == UniversalMqttClientStatus.connected) {
          _connectToPriceTopic();
          _connectToOHLC();
          _connectToOrderBook();
        }
      },
    );
    try {
      await unAuthorizedClient.connect();
    } catch (e) {
      log.e(
        e.toString(),
      );
    }
    return Future.value(true);
  }

  int getActivePairId() {
    int id;
    // ignore: invalid_use_of_protected_member
    for (var item in pairs.value) {
      if (item['name'] == currentPairName.value) {
        id = item['id'];
        break;
      }
    }
    return id;
  }

//////////////////////////////////////////////////////////////// inputs
  Fee _fee() {
    return pairBalanceData.value.fee;
  }

  void handleTotalChange(String value) {
    final String v = _applyCurrections(value);
    if (v == '') {
      amountValue.value = '';
      return;
    }
    double total = v.toDouble();
    if (priceValue.value != '') {
      double price = priceValue.value.toDouble();
      double amount = NumUtil.divide(total, price);
      final newAmount = amount.parseDoubleToString();
      amountValue.value = newAmount;
      if (mainActiveIndex.value == 1 && subActiveIndex.value != 1) {
        final fee = pairBalanceData.value.fee.makerFee;
        final newTradeFee = NumUtil.multiply(total, fee);
        tradeFee.value = newTradeFee.toStringAsFixed(8);
        youGet.value = NumUtil.subtract(total, newTradeFee).toStringAsFixed(8);
      } else {
        setValueToYouGet(amount: newAmount);
      }
    } else {
      //amountValue.value = '';
    }
    updateInputValue(stream: totalValue, newValue: v);
  }

  void handleAmountChange(String value, {bool fromInput = true}) {
    if (fromInput) {
      _resetPercent();
    }
    final v = _applyCurrections(value);
    if (v == '') {
      amountValue.value = '';
      if (subActiveIndex.value == 1) {
        youGet.value = '';
        tradeFee.value = '';
      }
      return;
    }
    double amount = v.toDouble();
    if (priceValue.value != '' || subActiveIndex.value == 1) {
      double total;
      if (subActiveIndex.value != 1) {
        total = NumUtil.multiply(amount, priceValue.value.toDouble());
        totalValue.value = total.parseDoubleToString();
      }
      if (mainActiveIndex.value == 1 && subActiveIndex.value != 1) {
        final fee = pairBalanceData.value.fee.makerFee;
        final newTradeFee = NumUtil.multiply(total, fee);
        tradeFee.value = newTradeFee.toStringAsFixed(8);
        final newYouGet = NumUtil.subtract(total, fee);
        youGet.value = newYouGet.toStringAsFixed(8);
      } else {
        updateInputValue(stream: amountValue, newValue: v);
        setValueToYouGet(amount: v.removeComma());
      }
    } else {
      //totalValue.value = '';
    }
    updateInputValue(stream: amountValue, newValue: v);
  }

  void handlePriceChange(String value) {
    final String v = _applyCurrections(value);
    if (v == '') {
      priceValue.value = '';
      return;
    }

    double price = v.toDouble();

    if (amountValue.value != '') {
      double total = NumUtil.multiply(amountValue.value.toDouble(), price);
      totalValue.value = total.toStringAsFixed(8);

      if (mainActiveIndex.value == 1 && subActiveIndex.value != 1) {
        final fee = pairBalanceData.value.fee.makerFee;
        double newTradeFee = NumUtil.multiply(total, fee);
        double newYouGet = NumUtil.subtract(total, fee);
        tradeFee.value = newTradeFee.toStringAsFixed(8);
        youGet.value = newYouGet.toStringAsFixed(8);
      } else {
        setValueToYouGet(
          amount: amountValue.value.removeComma(),
        );
      }
    } else {
      totalValue.value = '';
    }
    updateInputValue(stream: priceValue, newValue: v);
  }

  updateMarketYouGet({bool force}) {
    if (subActiveIndex.value == 1) {
      final fee = _fee().takerFee;
      if (amountValue.value != '') {
        final amount = amountValue.value.toDouble();
        double eqEmount;
        if (mainActiveIndex.value == 0) {
          eqEmount =
              NumUtil.divide(amount, currentPairPrice.value.price.toDouble());
        } else {
          eqEmount =
              NumUtil.multiply(amount, currentPairPrice.value.price.toDouble());
        }
        final marketFeeValue = NumUtil.multiply(eqEmount, fee);
        tradeFee.value = marketFeeValue.toStringAsFixed(8);
        youGet.value = NumUtil.subtract(eqEmount, marketFeeValue).toStringAsFixed(8);
      } else {
        tradeFee.value = '';
        youGet.value = '';
      }
    }
  }

  void setValueToYouGet({String amount}) {
    final fee = subActiveIndex.value == 1 ? _fee().takerFee : _fee().makerFee;
    if (subActiveIndex.value != 1) {
      if (mainActiveIndex.value == 0) {
        final amountByFee = NumUtil.multiply(amount.toDouble(), fee);
        youGet.value =
            NumUtil.subtract(amount.toDouble(), amountByFee).toStringAsFixed(8);
        tradeFee.value = amountByFee.toStringAsFixed(8);
      }
    } else {
      updateMarketYouGet();
    }
    return;
  }

  _resetPercent() {
    if (selectedPercentIndex.value != -1) {
      selectedPercentIndex.value = -1;
    }
  }

  void handleStopChange(String v) {
    stopValue.value = v;
  }

  Future _connectToOrderBook() async {
    orderbookSubscription = unAuthorizedClient
        .handleString(_orderBookTopic, MqttQos.exactlyOnce)
        .listen(
      (message) async {
        if (message != null && activeChart.value == TradeTopCharts.OrderBook) {
          orderBookThrottle.throttle(
            () async {
              if (message != null) {
                orderBookData.value =
                    await compute(parseOrderBookData, message);
              }
            },
          );
        }
        orderBookTopicInitialized = true;
      },
    );
    return Future.value(true);
  }

  String get _ohlcTopic =>
      "${Constants.ohlcTopic}${selectedTimeFrame.value}/${currentPairName.value}";

  String get _orderBookTopic =>
      "${Constants.orderbookTopic}${currentPairName.value}";
  String get _priceTopic => Constants.priceTopic;

  Future _connectToOHLC() async {
    ohlcSubscription =
        unAuthorizedClient.handleString(_ohlcTopic, MqttQos.exactlyOnce).listen(
      (message) {
        ohlcTopicInitialized = true;
        // final ohlcObj = OhlcModel.fromJson(jsonDecode(message));
        if (activeChart.value == TradeTopCharts.OHLC &&
            message != null) if (!(GetPlatform.isWeb)) {
          lastOhlcValue.value = jsonDecode(message);
        } else {
          lastOhlcThrottle.throttle(() {
            lastOhlcValue.value = jsonDecode(message);
          });
        }
      },
    );
    return Future.value(true);
  }

  void _connectToPriceTopic() {
    final pairsObj = <String, dynamic>{};
    priceSubscription = unAuthorizedClient
        .handleString(Constants.priceTopic, MqttQos.exactlyOnce)
        .listen(
      (message) {
        if (message != null) {
          priceTopicInitialized = true;
          final price = PriceModel.fromJson(jsonDecode(message));
          pairsObj[price.name] = price;
          lastPrice.value = price;

          if (price.name == currentPairName.value) {
            currentPairPrice.value = price;
            if (subActiveIndex.value == 1) {
              updateMarketYouGet();
            }
          }
          if (pairsObj[price.name] == null) {
            priceArray.add(price);
          }
        }
      },
    );
  }

  void toggleCharts() {
    if (activeChart.value == TradeTopCharts.OrderBook) {
      activeChart.value = TradeTopCharts.OHLC;
      return;
    }
    activeChart.value = TradeTopCharts.OrderBook;
  }

  void updateInputValue({RxString stream, String newValue}) {
    if (double.parse(newValue) < 1) {
      stream.value = newValue;
      return;
    }
    stream.value = newValue;
  }

  _resetAllInputs() {
    amountValue.value = '';
    totalValue.value = '';
    tradeFee.value = '';
    priceValue.value = '';
    stopValue.value = '';
    youGet.value = '';
  }

  String _applyCurrections(String value) {
    String v = value.removeComma();
    if (v == '') {
      //_resetAllInputs();
      totalValue.value = '';
      youGet.value = '';
      tradeFee.value = '';
      return v;
    }
    if (StringUtils.countChars(v, '.') > 1) {
      v = '0';
    }
    if (v.startsWith('.') && v != '.') {
      v = '0' + v;
    }
    if (v.contains('..')) {
      v = '0';
    }
    return v;
  }

  handleAskClick(data) {
    applyOrderBookPrice(data: data, goToMainTabIndex: 0);
  }

  handleBidClick(data) {
    applyOrderBookPrice(data: data, goToMainTabIndex: 1);
  }

  applyOrderBookPrice({int goToMainTabIndex, data}) {
    final String price =
        double.parse(data['price']).toStringWithoutTrailingZeros();

    if (mainActiveIndex.value != goToMainTabIndex) {
      handleMainTabChange(goToMainTabIndex);
    }
    handlePriceChange(price);
  }

  bool _validateNewOrderInputs() {
    String message;
    if (subActiveIndex.value == 0 || subActiveIndex.value == 2) {
      if (priceValue.value == '' || !(priceValue.value.toDouble() > 0)) {
        message = 'price is empty';
      }
    }
    if (subActiveIndex.value == 2) {
      if (stopValue.value == '' || !(stopValue.value.toDouble() > 0)) {
        message = 'stop input is empty';
      }
    }

    if (amountValue.value == '' || !(amountValue.value.toDouble() > 0)) {
      message = 'amount is empty';
    }
    if (message != null) {
      inputValidationToastThrottle.throttle(() => toastWarning(message));
      return false;
    }
    return true;
  }

  void _toastToDeposit() {
    final String coinName =
        mainActiveIndex.value == 0 ? _coinName1 : _coinName0;

    final idx = Constants.currencyArray()
        .indexWhere((element) => element.name == coinName);

    final coin = Constants.currencyArray()[idx];

    balanceToastThrottle.throttle(() {
      toastAction(
          'your balance is too low',
          UBButton(
              height: 26,
              fontSize: 11.0,
              width: 85,
              onClick: () {
                _openDepositPage(coin: coin);
              },
              text: 'Deposit $coinName'),
          duration: 4000);
    });
  }

  void _openDepositPage({AutoCompleteItem coin}) async {
    AccountController accountController = Get.find();
    if (accountController.accountData.value.isAccountVerified != true) {
      toastWarning('Please check your email and verify your account');
      return;
    }
    showLoadingOverlay.value = true;
    try {
      await openDepositPopup(coin: coin);
    } catch (e) {
    } finally {
      showLoadingOverlay.value = false;
    }
  }

  void checkForEssentialData() {
    if (
        // pairBalanceData.value.sum == null &&
        isLoadingPairBalance.value == false) {
      getPairBalances();
    }
  }
}

FutureOr<PairBalanceModel> parseBalanceData(message) {
  return PairBalanceModel.fromJson(message);
}

final orderBookThrottle = Throttling(duration: const Duration(seconds: 1));

FutureOr<Map<dynamic, dynamic>> parseOrderBookData(data) {
  final bookObj = jsonDecode(data);
  List rev = bookObj['bids'];
  List<dynamic> bids = rev.reversed.toList();
  List<dynamic> asks = bookObj['asks'];
  if (bids.length != asks.length) {
    if (asks.length > bids.length) {
      asks = asks.sublist(0, bids.length);
    } else {
      bids = bids.sublist(0, asks.length);
    }
  }
  bookObj['bids'] = reduceOrderBook(data: bids);
  bookObj['asks'] = reduceOrderBook(data: asks);

  return bookObj;
}

reduceOrderBook({List<dynamic> data}) {
  final len = data.length;
  final tmp = [];
  if (len > 7) {
    final ev = (len ~/ 6);
    for (var i = 0; i < len; i++) {
      if (i % ev == 0) {
        tmp.add(data[i]);
      }
    }
    if (tmp.length > 7) {
      tmp.assignAll(tmp.sublist(0, 7));
    } else if (tmp.length == 6) {
      tmp.add(data.last);
    }
  } else {
    tmp.assignAll(data.toList());
  }
  return tmp;
}
