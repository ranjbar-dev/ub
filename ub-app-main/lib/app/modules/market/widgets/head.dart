import 'package:flutter/material.dart';
import 'package:get/state_manager.dart';
import '../controllers/market_controller.dart';
import '../../../../generated/assets.gen.dart';
import '../../../../generated/colors.gen.dart';
import '../../../../utils/mixins/commonConsts.dart';

class MarketWatchHead extends GetView<MarketController> {
  @override
  Widget build(BuildContext context) {
    return Obx(() {
      final coinSortDirection = controller.coinSortDirection.value;
      final lastPriceSortDirection = controller.lastPriceSortDirection.value;
      final activeTabIndex = controller.activeTabIndex.value;
      final changeSortDirection = controller.changeSortDirection.value;

      return Container(
          height: 24,
          color: ColorName.grey16,
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 12.0),
            child: Row(
              children: [
                SizedBox(
                    width: thirdWidthPlus24,
                    child: GestureDetector(
                      onTap: () {
                        if (activeTabIndex != 0)
                          controller.handleCoinSortClick();
                      },
                      child: _header(
                          title: 'Coin',
                          sortDirection: coinSortDirection,
                          activeTabIndex: activeTabIndex),
                    )),
                GestureDetector(
                  onTap: () {
                    if (activeTabIndex != 0)
                      controller.handleLastPriceSortClick();
                  },
                  child: _header(
                      title: 'Last Price',
                      sortDirection: lastPriceSortDirection,
                      activeTabIndex: activeTabIndex),
                ),
                const Spacer(),
                SizedBox(
                  child: Row(
                    children: [
                      GestureDetector(
                        onTap: () {
                          if (activeTabIndex != 0)
                            controller.hanldeChangeSortClick();
                        },
                        child: _header(
                            title: 'Chg%',
                            sortDirection: changeSortDirection,
                            activeTabIndex: activeTabIndex),
                      ),
                      if (activeTabIndex == 0) hspace4,
                      if (activeTabIndex == 0)
                        GestureDetector(
                          onTap: controller.toggleEditFavs,
                          child: SizedBox(
                            width: 30,
                            height: 24,
                            child: Assets.images.editWithUnderline.svg(
                              color: ColorName.primaryBlue,
                            ),
                          ),
                        )
                    ],
                  ),
                )
              ],
            ),
          ));
    });
  }

  _header({String title, SortDirection sortDirection, int activeTabIndex}) {
    return Container(
      color: Colors.transparent,
      child: Row(
        children: [
          Text(
            title,
            style: grey80Bold13,
          ),
          hspace4,
          if (activeTabIndex != 0)
            Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Assets.images.keyboardUp.svg(
                  color: sortDirection == SortDirection.ASC
                      ? ColorName.white
                      : ColorName.grey80,
                ),
                Assets.images.keyBoardDown.svg(
                  color: sortDirection == SortDirection.DESC
                      ? ColorName.white
                      : ColorName.grey80,
                )
              ],
            )
        ],
      ),
    );
  }
}
