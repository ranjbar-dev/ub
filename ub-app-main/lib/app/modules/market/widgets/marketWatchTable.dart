import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../common/components/UBScrollBar.dart';
import '../../../common/components/UBoops.dart';
import '../../../global/currency_pairs_model.dart';
import '../controllers/market_controller.dart';
import '../views/favoritesList.dart';
import 'marketPriceRow.dart';

class MarketWatchTable extends GetView<MarketController> {
  @override
  Widget build(BuildContext context) {
    final scrollController = controller.pairsScrollController;
    return Expanded(
      child: Obx(
        () {
          final activeTabIndex = controller.activeTabIndex.value;
          final activeTabString = controller.activeTabString.value;
          final searchValue =
              controller.searchComponentParameters.value.searchValue;
          if (activeTabIndex != 0) {
            return SingleChildScrollView(
              child: Column(
                children: [
                  CoinList(
                    searchValue: searchValue,
                    onRowClick: controller.handlePriceRowClick,
                    pairsArray: controller.sorted.isNotEmpty
                        ? controller.sorted
                        // ignore: invalid_use_of_protected_member
                        : controller.pairs.value,
                    scrollController: scrollController,
                    activeTabString: activeTabString,
                    activeTabIndex: activeTabIndex,
                  ),
                ],
              ),
            );
          }
          return FavoritesList();
        },
      ),
    );
  }
}

class CoinList extends StatelessWidget {
  final List<Pairs> pairsArray;
  final String activeTabString;
  final String searchValue;
  final int activeTabIndex;
  final ScrollController scrollController;
  final Function(String) onRowClick;

  const CoinList({
    Key key,
    @required this.pairsArray,
    this.activeTabString,
    this.activeTabIndex,
    @required this.scrollController,
    @required this.onRowClick,
    this.searchValue = '',
  }) : super(key: key);
  @override
  Widget build(BuildContext context) {
    final filteredArray = [];
    pairsArray.forEach((element) {
      if (activeTabIndex != 0 &&
          element.pairName
              .replaceAll('-', '')
              .contains(searchValue.toUpperCase()) &&
          (activeTabIndex == 1 ||
              element.pairName.split('-')[1] == activeTabString)) {
        filteredArray.add(element);
      }
    });
    return filteredArray.length == 0
        ? UBoops(
            variant: OopsVariant.SearchOops,
          )
        : UBScrollBar(
            scrollController: scrollController,
            itemCount: filteredArray.length,
            builder: (BuildContext context, int index) {
              final data = filteredArray[index];
              return MarketPriceRow(
                onClick: onRowClick,
                key: ValueKey(data.pairName),
                price: data.formattedPrice,
                volume: data.formattedVolume,
                equivalentPrice: data.formattedEquivalentPrice,
                data: data,
              );
            },
          );
  }
}
