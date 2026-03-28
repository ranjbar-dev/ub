import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../common/components/UBScaleSwitcher.dart';
import '../../../common/components/controlledInput.dart';
import '../../../common/components/pageContainer.dart';
import '../controllers/market_controller.dart';
import '../widgets/head.dart';
import '../widgets/marketWatchTable.dart';
import '../widgets/topTabs.dart';
import '../../../../generated/assets.gen.dart';
import '../../../../generated/colors.gen.dart';
import '../../../../utils/mixins/commonConsts.dart';

class MarketView extends GetView<MarketController> {
  @override
  Widget build(BuildContext context) {
    controller.isPageActive.value = true;
    return PageContainer(
      activeBottomNavIndex: 1,
      beforePop: () {
        controller.isPageActive.value = false;
      },
      child: Obx(() {
        final activeTopTab = controller.activeTabIndex.value;
        return Stack(
          children: [
            Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Container(
                  width: activeTopTab != 0 ? Get.width - 50 : Get.width,
                  child: MarketTopTabs(),
                ),
                MarketWatchHead(),
                MarketWatchTable(),
              ],
            ),
            if (activeTopTab != 0)
              Positioned(
                child: SearchBar(key: ValueKey('MarketSearchBar')),
                top: 0.0,
                right: 0.0,
              )
          ],
        );
      }),
    );
  }
}

class SearchBar extends GetView<MarketController> {
  const SearchBar({Key key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Obx(() {
      final parameters = controller.searchComponentParameters.value;
      return GestureDetector(
        onTap: controller.toggleIsSearchOpen,
        child: AnimatedContainer(
          padding: const EdgeInsets.only(top: 6.0),
          height: 49.0,
          width: parameters.isOpen ? Get.width : 50,
          decoration: const BoxDecoration(
            color: Colors.black,
            border: const Border(
              bottom: const BorderSide(
                width: 2.0,
                color: ColorName.grey23,
              ),
            ),
          ),
          duration: 300.milliseconds,
          child: Row(
            children: [
              if (parameters.isInputVisible) hspace12,
              if (parameters.isInputVisible)
                SizedBox(
                  width: 140.0,
                  child: ControlledTextField(
                    key: ValueKey('marketSearch'),
                    noBorder: true,
                    autoFocus: true,
                    text: parameters.searchValue,
                    labelText: 'Search pair...',
                    onChanged: (v) {
                      controller.handleSearchChange(v);
                    },
                  ),
                ),
              fill,
              UBScaleSwitcher(
                child1: Assets.images.closeIcon
                    .svg(width: 24, key: ValueKey('topClose')),
                child2: Assets.images.searchIcon
                    .svg(width: 24, key: ValueKey('topSearchIcon')),
                conditionToShowChild1: parameters.isInputVisible,
              ),
              hspace12,
            ],
          ),
        ),
      );
    });
  }
}
