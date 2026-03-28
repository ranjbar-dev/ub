import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../../generated/assets.gen.dart';
import '../../../../generated/colors.gen.dart';
import '../controllers/market_controller.dart';

class MarketTopTabs extends GetView<MarketController> {
  @override
  Widget build(BuildContext context) {
    return Obx(
      () {
        final tabs = controller.tabCurrencies;
        final activeTab = controller.activeTabIndex.value;
        return tabs.length == 0
            ? const SizedBox(
                height: 46,
              )
            : Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Transform.translate(
                    offset: const Offset(0, 47.5),
                    child: Container(
                      width: double.infinity,
                      height: 1,
                      color: ColorName.grey36,
                    ),
                  ),
                  DefaultTabController(
                    length: controller.tabCurrencies.length + 2,
                    initialIndex: controller.activeTabIndex.value,
                    child: TabBar(
                      isScrollable: true,
                      onTap: controller.handleTabChange,
                      labelColor: ColorName.primaryBlue,
                      unselectedLabelColor: ColorName.greybf,
                      indicatorColor: ColorName.primaryBlue,
                      labelPadding: const EdgeInsets.symmetric(
                        horizontal: 11,
                      ),
                      tabs: [
                        Container(
                          width: 20,
                          height: 20,
                          child: activeTab == 0
                              ? Assets.images.filledStar.svg()
                              : Assets.images.filledStar
                                  .svg(color: ColorName.greybf),
                        ),
                        Tab(
                          child: RichText(
                            text: TextSpan(
                              text: 'All',
                              style: TextStyle(
                                fontSize: 14,
                                fontWeight: FontWeight.w600,
                                color: activeTab == 1
                                    ? ColorName.primaryBlue
                                    : ColorName.grey80,
                              ),
                            ),
                          ),
                        ),
                        for (var i = 0;
                            i < controller.tabCurrencies.length;
                            i++)
                          Tab(
                            // text: item,
                            child: RichText(
                              text: TextSpan(
                                text: controller.tabCurrencies[i],
                                style: TextStyle(
                                  fontSize: 14,
                                  fontWeight: FontWeight.w600,
                                  color: activeTab == i + 2
                                      ? ColorName.primaryBlue
                                      : ColorName.grey80,
                                ),
                              ),
                            ),
                          ),
                      ],
                    ),
                  ),
                ],
              );
      },
    );
  }
}
