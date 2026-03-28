import 'package:supercharged/supercharged.dart';
import 'package:dio/dio.dart' show DioError;
import 'package:flutter/foundation.dart';
import 'package:get/get.dart';
import '../../../common/components/UBWrappedButtons.dart';
import '../../../common/custom/k_chart/flutter_k_chart.dart';
import 'trade_controller.dart';
import '../provider/ohlcProvider.dart';
import '../../../../utils/mixins/toast.dart';

class OHLCChartController extends GetxController with Toaster {
  final TradeController tradeController = Get.find();
  final OHLCProvider ohlcProvider = OHLCProvider();
  final isOhlcDetailsOpen = false.obs;
  final isLoadingOhlc = false.obs;
  final chartData = <KLineEntity>[].obs;
  final mainState = (MainState.NONE).obs;
  final secondaryState = (SecondaryState.NONE).obs;
  final isLine = false.obs;
  final bids = <DepthEntity>[].obs;
  final asks = <DepthEntity>[].obs;
  var currentPairName = '';
  var timeFrame = '';
  var currentResolution;
  final isTimeFramePoupOpen = false.obs;
  List<WrappedButtonModel> timeFrameButtons;
  final selectedTimeFrameButtonIndex = (0).obs;
  bool canLoadNewData = false;

  @override
  void onInit() async {
    timeFrameButtons = tradeController.timeFrameButtons;
    currentPairName = tradeController.currentPairName.value;
    timeFrame = tradeController.selectedTimeFrame.value;
    selectedTimeFrameButtonIndex.value =
        timeFrameButtons.indexWhere((element) => element.value == timeFrame);
    getPairOhlc();
    tradeController.currentPairName.listen((v) {
      currentPairName = v;
      getPairOhlc();
    });

    tradeController.selectedTimeFrame.listen((v) {
      timeFrame = v;
      getPairOhlc();
    });

    tradeController.lastOhlcValue.listen((v) {
      if (v["pair"] == currentPairName) {
        updateLastRecord(v);
      }
    });
    super.onInit();
  }

  @override
  void onReady() {
    super.onReady();
  }

  @override
  void onClose() {}

  void getPairOhlc() async {
    canLoadNewData = false;
    var resolution;
    currentResolution = resolution;
    try {
      final now = DateTime.now();
      var from = _parseSecond(1.days.ago());
      var to = _parseSecond(now);

      if (timeFrame == '1minute') {
        final past = 1.days.ago();
        resolution = '1';
        from = _parseSecond(past);
        to = _parseSecond(now);
      }
      if (timeFrame == '5minutes') {
        resolution = '5';
        final past = 3.days.ago();
        from = _parseSecond(past);
        to = _parseSecond(now);
      }
      if (timeFrame == '1hour') {
        resolution = '60';
        final past = 30.days.ago();
        from = _parseSecond(past);
        to = _parseSecond(now);
      }
      if (timeFrame == '1day') {
        resolution = '1D';
        final past = 365.days.ago();
        from = _parseSecond(past);
        to = _parseSecond(now);
      }

      isLoadingOhlc.value = true;
      final response = await ohlcProvider.getOHLC(
        symbol: currentPairName.replaceFirst('-', '/'),
        resolution: resolution,
        from: from,
        to: to,
      );
      final List<dynamic> dataArray = response['bars'];
      final data = await compute(_parseChartData, dataArray);
      chartData.assignAll(data);
      canLoadNewData = true;
    } on DioError catch (e) {
      toastDioError(e);
    } finally {
      isLoadingOhlc.value = false;
    }
  }

  String _parseSecond(DateTime past) {
    final tmp = past.millisecondsSinceEpoch.toString();
    return tmp.substring(0, tmp.length - 3);
  }

  void updateLastRecord(Map<String, dynamic> newData) {
    if (canLoadNewData && chartData.isNotEmpty) {
      // ignore: invalid_use_of_protected_member
      final tmp = List<KLineEntity>.from(chartData.value);
      final lastItem = tmp.last;
      final now = DateTime.now();
      //final offset = now.timeZoneOffset.inMilliseconds;
//final parsedNewDataDate=DateTime.parse(newData["closeTime"])+now.timeZoneOffset;
      final DateTime newDataTime =
          DateTime.parse(newData["startTime"]) + now.timeZoneOffset;

      final DateTime lastItemTime =
          new DateTime.fromMillisecondsSinceEpoch(lastItem.time);

      final shouldUpdate = lastItemTime == newDataTime;

      // _shouldUpdate(
      //     differenceDuration: lastItemTime.difference(newDataTime),
      //     timePeriod: newData["timeFrame"]);
// MarketModel{open: 56143.43, high: 56321.18, low: 55460.16, close: 56020.57, vol: 4864.080125, time: 1620129600000, amount: null, ratio: null, change: null}
//{pair: BTC-USDT, timeFrame: 1minute, startTime: 2021-05-04 13:24:00, closeTime: 2021-05-04 13:24:59, openPrice: 56012.74000000, closePrice: 55996.13000000, highPrice: 56012.74000000, lowPrice: 55994.61000000, baseVolume: 11.24162200, quoteVolume: 629553.77281168, takerBuyBaseVolume: 5.09207200, takerBuyQuoteVolume: 285162.54234474}

      final entity = KLineEntity.fromJson(
        {
          "time": newDataTime.millisecondsSinceEpoch,
          "open": double.parse(newData["openPrice"]),
          "close": double.parse(newData["closePrice"]),
          "high": double.parse(newData["highPrice"]),
          "low": double.parse(newData["lowPrice"]),
          "vol": double.parse(newData["baseVolume"]) -
              double.parse(newData["takerBuyBaseVolume"])
        },
      );
      final lastParsedData = DataUtil.calculate([entity])[0];
      if (shouldUpdate) {
        tmp.last = lastParsedData;
        chartData.assignAll(tmp);
        return;
      }
      tmp.add(lastParsedData);
      chartData.assignAll(tmp);
    }
  }

  void handleIsOhlcDetailsOpen({bool isOpen}) {
    isOhlcDetailsOpen.value = isOpen;
  }

  void handleTimeFrameChange(int index) {
    isTimeFramePoupOpen.value = false;
    if (index != selectedTimeFrameButtonIndex.value) {
      selectedTimeFrameButtonIndex.value = index;
      tradeController.handeTimeFrameChange(timeFrameButtons[index].value);
    }
  }
}

List<KLineEntity> _parseChartData(List dataArray) {
  final kList = List.filled(dataArray.length, KLineEntity.fromJson({}));
  for (var i = 0; i < kList.length; i++) {
    final rawData = dataArray[i];
    kList[i] = KLineEntity.fromJson(
      {
        "time": rawData["time"],
        "open": rawData["open"],
        "close": rawData["close"],
        "high": rawData["high"],
        "low": rawData["low"],
        "vol": rawData["volume"]
      },
    );
  }
  final data = kList.toList().cast<KLineEntity>();
  // return data;
  return DataUtil.calculate(data);
}
