import 'package:dio/dio.dart' show DioError;
import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:get_storage/get_storage.dart';
import '../../../global/controller/globalController.dart';
import '../../../global/currency_pairs_model.dart';
import '../../trade/controllers/trade_controller.dart';
import '../../trade/provider/favoritePairsProvider.dart';
import '../../../routes/app_pages.dart';
import '../../../../services/storageKeys.dart';
import '../../../../utils/marketUtils.dart';
import '../../../../utils/mixins/formatters.dart';
import '../../../../utils/mixins/toast.dart';
import '../../../../utils/pairAndCurrencyUtils.dart';
import '../../../../utils/throttle.dart';

enum SortDirection { ASC, DESC, NONE }

class SearchParameters {
  final bool isOpen;
  final bool isInputVisible;
  final String searchValue;

  SearchParameters(
      {@required this.isOpen,
      @required this.isInputVisible,
      @required this.searchValue});
}

class MarketController extends GetxController with Toaster, Formatter {
  final thr1500 = Throttling(duration: const Duration(milliseconds: 1500));

  final GlobalController globalController = Get.find();
  final TradeController tradeController = Get.find();
  final pairsScrollController = ScrollController();
  final favsScrollController = ScrollController();
  final isPageActive = false.obs;

  final searchComponentParameters = SearchParameters(
    isOpen: false,
    isInputVisible: false,
    searchValue: '',
  ).obs;

  final GetStorage storage = GetStorage();
  final FavoritePairsProvider favoritePairsProvider = FavoritePairsProvider();
  final pairs = <Pairs>[].obs;
  final favorites = <Pairs>[].obs;
  final orderedPairs = <Pairs>[].obs;
  final sorted = <Pairs>[].obs;
  List<Pairs> tmpPairs = [];
  final tabCurrencies = [].obs;
  final activeTabIndex = 1.obs;
  final activeTabString = 'All'.obs;

  final coinSortDirection = (SortDirection.NONE).obs;
  final lastPriceSortDirection = (SortDirection.NONE).obs;
  final changeSortDirection = (SortDirection.NONE).obs;

  var pairsMapper = {};
  Map<String, Pairs> pairsHashMap = {};

  @override
  void onInit() {
    final storedActiveTabIndex =
        storage.read(StorageKeys.activeMarketTabIndex) ?? 1;
    handleTabChange(storedActiveTabIndex);
    final storedOrderedPairs = storage.read(StorageKeys.orderedPairs);
    if (storedOrderedPairs != null) {
      final tmp = <Pairs>[];
      for (var item in storedOrderedPairs) {
        tmp.add(Pairs.fromJson(item));
      }
      orderedPairs.assignAll(tmp);
      addIsFavoriteToPairs(stream: orderedPairs);
      assignFavoritesFromOrderedPairs();
    }

    if (globalController.allCurrencyPairs.length > 0) {
      _fillTabs();
    } else {
      globalController.allCurrencyPairs.listen((v) => {_fillTabs()});
    }

    super.onInit();
  }

  @override
  void onReady() {
    super.onReady();
  }

  @override
  void onClose() {}

  void _fillTabs() {
    pairsHashMap = Map.from(PairAndCurrencyUtils.pairsMap.value);
    final tmp =
        globalController.allCurrencyPairs.map((element) => element.code);
    tabCurrencies.assignAll(tmp);
    List<Pairs> tmpCurrencyPairs = [];
    //extract pairs array
    globalController.allCurrencyPairs.forEach(
      (item) => {
        item.pairs.forEach((pair) {
          if (pair.price == null) {
            pair.price = '0.0';
            pair.percent = '0.0';
            pair.volume = '--';
            pair.equivalentPrice = '--';
            pair.formattedEquivalentPrice = '--';
            pair.formattedVolume = '--';
            pair.formattedPrice = '--';
          }
          tmpCurrencyPairs.add(pair);
        })
      },
    );
    for (var i = 0; i < tmpCurrencyPairs.length; i++) {
      final el = tmpCurrencyPairs[i];
      pairsMapper[el.pairName] = i;
    }
    pairs.assignAll(tmpCurrencyPairs);
    addIsFavoriteToPairs(stream: pairs);
    if (orderedPairs.isEmpty) {
      assignFavoritesFromUnOrderedPairsAndFillOrderedPairsList();
    }

    // constantPairs.assignAll(tmpCurrencyPairs);
    tmpPairs = tmpCurrencyPairs;
    tradeController.lastPrice.listen(
      (lastPrice) {
        if (pairsMapper[lastPrice.name] != null && isPageActive.value == true) {
          //print(lastPrice.volume
          //    .currencyFormat(removeInsignificantZeros: true, compact: true));

          final newPrice = Pairs.fromJson(
            MarketUtils.priceJson(
              lastPrice,
            ),
          );
          if (orderedPairs.isNotEmpty) {
            if (orderedPairs.indexWhere(
                    (element) => element.pairId == newPrice.pairId) ==
                -1) {
              orderedPairs.add(newPrice);
              saveOrderedPairs();
            }
          }
          if (favorites.isNotEmpty) {
            final idx = favorites
                .indexWhere((element) => element.pairId == newPrice.pairId);
            if (idx != -1) {
              // ignore: invalid_use_of_protected_member
              favorites.value[idx] = newPrice;
              favorites.refresh();
            }
          }
          if (sorted.isNotEmpty && activeTabIndex.value != 0) {
            final idx = sorted
                .indexWhere((element) => element.pairName == newPrice.pairName);
            // ignore: invalid_use_of_protected_member
            sorted.value[idx] = newPrice;
            sorted.refresh();
          } else if (activeTabIndex.value != 0) {
            tmpPairs[pairsMapper[lastPrice.name]] = newPrice;
            if (!(GetPlatform.isWeb)) {
              pairs.assignAll(tmpPairs);
            } else {
              thr1500.throttle(() {
                pairs.assignAll(tmpPairs);
              });
            }
          }
        }
      },
    );
  }

  void handleTabChange(int index) {
    if (index == 0) {
      activeTabString.value = 'Favs';
    } else if (index == 1) {
      activeTabString.value = "All";
    } else {
      activeTabString.value = tabCurrencies[index - 2];
    }
    activeTabIndex.value = index;

    //tabCurrencies might change, but these are constant
    if (index == 0 || index == 1) {
      storage.write(StorageKeys.activeMarketTabIndex, index);
    }
  }

  void handleCoinSortClick() {
    lastPriceSortDirection.value = SortDirection.NONE;
    changeSortDirection.value = SortDirection.NONE;
    _filterSelect(
      stream: coinSortDirection,
      sortDesc: (Pairs a, Pairs b) => b.pairName.compareTo(a.pairName),
      sortAsc: (Pairs a, Pairs b) => a.pairName.compareTo(b.pairName),
    );
  }

  void handleLastPriceSortClick() {
    coinSortDirection.value = SortDirection.NONE;
    changeSortDirection.value = SortDirection.NONE;
    _filterSelect(
      stream: lastPriceSortDirection,
      sortDesc: (Pairs a, Pairs b) =>
          double.parse(b.price).compareTo(double.parse(a.price)),
      sortAsc: (Pairs a, Pairs b) =>
          double.parse(a.price).compareTo(double.parse(b.price)),
    );
  }

  void hanldeChangeSortClick() {
    coinSortDirection.value = SortDirection.NONE;
    lastPriceSortDirection.value = SortDirection.NONE;
    _filterSelect(
      stream: changeSortDirection,
      sortDesc: (Pairs a, Pairs b) =>
          double.parse(b.percent).compareTo(double.parse(a.percent)),
      sortAsc: (Pairs a, Pairs b) =>
          double.parse(a.percent).compareTo(double.parse(b.percent)),
    );
  }

  void _filterSelect({
    Rx<SortDirection> stream,
    Function(Pairs a, Pairs b) sortDesc,
    Function(Pairs a, Pairs b) sortAsc,
  }) {
    // ignore: invalid_use_of_protected_member
    final List<Pairs> pairArray = List.from((pairs.value));
    switch (stream.value) {
      case SortDirection.NONE:
        stream.value = SortDirection.DESC;
        pairArray.sort((a, b) => sortDesc(a, b));
        sorted.assignAll(pairArray);
        break;
      case SortDirection.DESC:
        stream.value = SortDirection.ASC;
        pairArray.sort((a, b) => sortAsc(a, b));
        sorted.assignAll(pairArray);
        break;
      case SortDirection.ASC:
        stream.value = SortDirection.NONE;
        sorted.assignAll([]);
        break;
    }
    return;
  }

  void toggleEditFavs() {
    Get.toNamed(AppPages.EDIT_FAVORITES);
  }

  void handleFavPairsReorder(int oldIndex, int newIndex) {
    if (newIndex > oldIndex) {
      newIndex -= 1;
    }
    final tmp = <Pairs>[];
    if (orderedPairs.isEmpty) {
      // ignore: invalid_use_of_protected_member
      tmp.assignAll(pairs.value);
    } else {
      // ignore: invalid_use_of_protected_member
      tmp.assignAll(orderedPairs.value);
    }
    final Pairs item = tmp.removeAt(oldIndex);
    tmp.insert(newIndex, item);
    orderedPairs.assignAll(tmp);
    _saveOrderedPairsAndAssignFavorites();
  }

  void addIsFavoriteToPairs({RxList<Pairs> stream}) {
    final favsJsonList = storage.read(StorageKeys.favPairs);
    if (favsJsonList != null && favsJsonList.length > 0) {
      for (var item in favsJsonList) {
        final idx =
            stream.indexWhere((element) => element.pairId == item['id']);
        if (idx != -1) {
          // ignore: invalid_use_of_protected_member
          stream.value[idx].isFavorite = true;
        }
      }
      stream.refresh();
    }
  }

  void assignFavoritesFromOrderedPairs() {
    final favs =
        orderedPairs.where((element) => element.isFavorite == true).toList();
    favorites.assignAll(favs);
  }

  void assignFavoritesFromUnOrderedPairsAndFillOrderedPairsList() {
    final favs = pairs.where((element) => element.isFavorite == true).toList();
    // ignore: invalid_use_of_protected_member
    orderedPairs.assignAll(pairs.value);
    favorites.assignAll(favs);
  }

  void movePairToTopOfOrderedList({Pairs pair}) {
    // ignore: invalid_use_of_protected_member
    final List<Pairs> tmp = List.from(orderedPairs.value);
    final idx = tmp.indexWhere((element) => element.pairId == pair.pairId);
    final removedPair = tmp.removeAt(idx);
    tmp.assignAll([removedPair, ...tmp]);
    orderedPairs.assignAll(tmp);
    _saveOrderedPairsAndAssignFavorites();
  }

  void saveOrderedPairs() {
    final storageJsonList = [];
    // ignore: invalid_use_of_protected_member
    for (var item in orderedPairs.value) {
      storageJsonList.add(item.toJson());
    }
    storage.write(StorageKeys.orderedPairs, storageJsonList);
  }

  void _saveOrderedPairsAndAssignFavorites() {
    saveOrderedPairs();
    assignFavoritesFromOrderedPairs();
  }

  void handleFavCheckboxClick({@required Pairs pair}) async {
    final idx =
        orderedPairs.indexWhere((element) => element.pairId == pair.pairId);
    // ignore: invalid_use_of_protected_member
    orderedPairs.value[idx].isFavorite = pair.isFavorite != true ? true : null;
    orderedPairs.refresh();
    _saveOrderedPairsAndAssignFavorites();
    try {
      await favoritePairsProvider.toggleFav(
          pairId: pair.pairId,
          //because we changed it before calling api
          add: pair.isFavorite == true);
    } on DioError catch (e) {
      toastDioError(e);
    }
  }

  toggleFavoriteByPairName({String pairName}) {
    final idx =
        orderedPairs.indexWhere((element) => element.pairName == pairName);
    // ignore: invalid_use_of_protected_member
    final pair = orderedPairs.value[idx];
    handleFavCheckboxClick(pair: pair);
  }

  isPairNameFavorite({String pairName}) {
    final idx = favorites.indexWhere((element) => element.pairName == pairName);
    if (idx != -1) {
      return true;
    }
    return false;
  }

  handlePriceRowClick(String pairName) {
    tradeController.handlePairChange(pairName);
    Get.offAllNamed(AppPages.TRADE);
  }

  void toggleIsSearchOpen() {
    final isOpen = searchComponentParameters.value.isOpen;
    if (isOpen) {
      final newValue = SearchParameters(
          isOpen: false, isInputVisible: false, searchValue: '');
      searchComponentParameters.value = newValue;
      return;
    }
    SearchParameters newValue =
        SearchParameters(isOpen: true, isInputVisible: false, searchValue: '');
    searchComponentParameters.value = newValue;
    Future.delayed(300.milliseconds).then((value) {
      newValue =
          SearchParameters(isOpen: true, isInputVisible: true, searchValue: '');
      searchComponentParameters.value = newValue;
      return;
    });
    return;
  }

  void handleSearchChange(String v) {
    if (v == '') {}
    final newValue =
        SearchParameters(isOpen: true, isInputVisible: true, searchValue: v);
    searchComponentParameters.value = newValue;
  }
}
