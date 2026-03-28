import 'dart:async';

import 'package:flutter/foundation.dart';
import 'package:get/get.dart';

import '../../../../utils/logger.dart';
import '../../../../utils/marketUtils.dart';
import '../../../../utils/pairAndCurrencyUtils.dart';
import '../../../../utils/throttle.dart';
import '../../../global/controller/globalController.dart';
import '../../../global/currency_pairs_model.dart';
import '../../../global/providers/commonDataProvider.dart';
import '../../../routes/app_pages.dart';
import '../../account/user_model.dart';
import '../../trade/controllers/trade_controller.dart';
import '../home_page_pair_price_model.dart';
import '../news_model.dart';
import '../providers/homePageProvider.dart';

class HomeController extends GetxController {
  final thr1500 = new Throttling(duration: const Duration(milliseconds: 1500));
  final thr1000 = new Throttling(duration: const Duration(milliseconds: 1000));

  final HomePageProvider homePageProvider = HomePageProvider();
  final CommonDataProvider commonDataProvider = CommonDataProvider();
  final TradeController tradeController = Get.find();
  final GlobalController globalController = Get.find();
  final isLoadingSparkLine = true.obs;
  final isSilentLoadingSparkLine = false.obs;
  final isLoadingNews = true.obs;
  final isRefreshing = false.obs;
  final isSilentLoadingNews = false.obs;
  bool isPageActive = true;
  final latestNews = <NewsModel>[].obs;
  final sparkLinePairs = <HomePagePairPriceModel>[].obs;
  final isUserVerified = false.obs;

  Map<String, Pairs> pairsHashMap = {};
  List<String> initialPopularPairs = [
    'BTC-USDT',
    'ETH-USDT',
    'BCH-USDT',
    'DASH-USDT',
    'DOGE-USDT',
    'MKR-USDT',
    'LTC-USDT',
    'ETH-BTC',
    'TRX-USDT',
  ];
  final popularPairs = <Pairs>[].obs;

  Timer timer;
  fillPairHashMap() {
    pairsHashMap = Map.from(PairAndCurrencyUtils.pairsMap.value);
  }

  @override
  void onInit() async {
    globalController.getVersion();
    //globalController.checkIfRedirectIsNeeded();
    if (PairAndCurrencyUtils.pairsMap.value != null) {
      fillPairHashMap();
    } else {
      PairAndCurrencyUtils.pairsMap.listen((v) => {fillPairHashMap()});
    }
    getPairPriceSparkLineChartData(silent: false);
    getHomaPageRestData();
    final List<Pairs> tmp = [];
    initialPopularPairs.forEach((element) {
      tmp.add(Pairs(
        pairName: element,
        price: '0.0',
        percent: '0.0',
        volume: '--',
        equivalentPrice: '--',
        formattedEquivalentPrice: '--',
        formattedPrice: '--',
        formattedVolume: '--',
      ));
    });
    popularPairs.assignAll(tmp);

    tradeController.lastPrice.listen((lastPrice) {
      final idx = popularPairs
          .indexWhere((element) => element.pairName == lastPrice.name);
      if (idx != -1 &&
          isPageActive == true &&
          pairsHashMap[lastPrice.name] != null) {
        final json = MarketUtils.priceJson(
          lastPrice,
        );

        final newPrice = Pairs.fromJson(json);

        // ignore: invalid_use_of_protected_member
        popularPairs.value[idx] = newPrice;
        if (!(GetPlatform.isWeb)) {
          popularPairs.refresh();
        } else {
          thr1500.throttle(() {
            popularPairs.refresh();
          });
        }
      }
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
    });
    checkIfUserEmailIsVerified();
    super.onInit();
  }

  void getHomaPageRestData() {
    getLatestNews(silent: false);
  }

  @override
  void onReady() {
    getAppEssentialData();
    super.onReady();
  }

  @override
  void onClose() {
    if (timer != null) {
      timer.cancel();
    }
  }

  Future getPairPriceSparkLineChartData({bool silent}) async {
    if (silent == true) {
      isSilentLoadingSparkLine.value = true;
    } else {
      isLoadingSparkLine.value = true;
    }
    try {
      final response = await homePageProvider.getPairsPrice();
      final parsed = await compute(
          parseHomePagePairsPrice, response['data']['pairs'][0]['pairs']);
      sparkLinePairs.assignAll(parsed);
    } catch (e) {
      log.e('homePageError !, Error getting sparkline data');
      log.e(e.toString());
    } finally {
      isSilentLoadingSparkLine.value = false;
      isLoadingSparkLine.value = false;
    }
    return Future.value();
  }

  Future getLatestNews({bool silent}) async {
    if (silent == true) {
      isSilentLoadingNews.value = true;
    } else {
      isLoadingNews.value = true;
    }
    try {
      final List response = await homePageProvider.getLastNews();
      final list = await compute(parseNews, response);
      latestNews.assignAll(list);
    } catch (e) {
      debugPrint('homePageError !,  Error getting latest news');
    } finally {
      isSilentLoadingNews.value = false;
      isLoadingNews.value = false;
    }
    return Future.value();
  }

  void handlePairClick(String pairName) {
    tradeController.handlePairChange(pairName);
    Get.offAllNamed(AppPages.TRADE);
  }

  // void startTimer() {
  //   if (timer != null) {
  //     timer.cancel();
  //   }
  //   timer = Timer.periodic(new Duration(seconds: 30), (timer) {
  //     getPairPriceSparkLineChartData(silent: false);
  //   });
  // }

  void handlePageLoaded() {
    // startTimer();
    isPageActive = true;
  }

  void handlePagePop() {
    // if (timer != null) {
    //   timer.cancel();
    // timer = null;
    // }
    isPageActive = false;
  }

  refreshHomePage() async {
    await Future.wait([
      getLatestNews(silent: true),
      getPairPriceSparkLineChartData(silent: true)
    ]);
    checkIfUserEmailIsVerified();
  }

  checkIfUserEmailIsVerified() async {
    final data = await commonDataProvider.getUserData();
    final userData = data;
    if (userData['status'] == true) {
      final userInfo = UserModel.fromJson(
        userData["data"],
      );

      isUserVerified.value = userInfo.isAccountVerified;
    }
  }

  void getAppEssentialData() {
    if (globalController.currencyPairsArray.isEmpty) {
      globalController.getPairsCurrenciesCountriesAndVersion();
    }
  }
}

List<HomePagePairPriceModel> parseHomePagePairsPrice(response) {
  final list = List<HomePagePairPriceModel>.from(
    response.map(
      (model) => HomePagePairPriceModel.fromJson(model),
    ),
  );
  return list;
}

List<NewsModel> parseNews(response) {
  final list = List<NewsModel>.from(
    response.map(
      (e) => NewsModel.fromJson(e),
    ),
  );
  return list;
}
