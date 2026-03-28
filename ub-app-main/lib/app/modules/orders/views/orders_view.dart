import 'dart:math';

import 'package:flutter/material.dart';

import 'package:get/get.dart';
import '../pages/openOrders/views/open_orders_view.dart';
import '../pages/orderHistory/views/order_history_view.dart';
import '../../../../generated/assets.gen.dart';
import '../../../../generated/colors.gen.dart';
import '../../../../generated/locales.g.dart';
import '../../../../utils/mixins/filterPopups.dart';

import '../controllers/orders_controller.dart';

class OrdersView extends GetView<OrdersController> with FilterPopups {
  final bool fullScreen;
  final int activeIndex;

  OrdersView({this.activeIndex, this.fullScreen = false});

  @override
  Widget build(BuildContext context) {
    final width = Get.width;
    return Material(
      child: WillPopScope(
        onWillPop: () {
          //toggles the isFullScreen Value
          controller.handleWillPop();
          return Future.value(true);
        },
        child: Container(
          width: double.infinity,
          height: 0,
          color: ColorName.black,
          child: Obx(
            () {
              final isFullScreen = controller.isFullScreen.value;
              final activeTabIndex = controller.activeTabIndex.value;
              return DefaultTabController(
                key: isFullScreen ? ValueKey('full') : ValueKey('min'),
                length: 2,
                initialIndex: activeTabIndex,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: <Widget>[
                    Transform.translate(
                      offset: Offset(0, 34.0),
                      child: Container(
                        width: double.infinity,
                        height: 1,
                        color: ColorName.black2c,
                      ),
                    ),
                    Padding(
                      padding: EdgeInsets.only(
                          top: isFullScreen
                              ? MediaQuery.of(context).padding.top
                              : 0),
                      child: Container(
                        width: width,
                        child: Stack(
                          children: [
                            Positioned(
                              right: 8,
                              child: GestureDetector(
                                onTap: controller.handleExpandOrdersClick,
                                child: Container(
                                  padding: const EdgeInsets.all(4),
                                  width: 32,
                                  height: 32,
                                  child: Transform.rotate(
                                    angle: isFullScreen ? 0 : pi,
                                    child: Assets.images.doubleArrowDown.svg(),
                                  ),
                                ),
                              ),
                            ),
                            Container(
                              height: 35,
                              width: 250,
                              child: TabBar(
                                onTap: (i) {
                                  if (GetPlatform.isWeb) {
                                    FocusScope.of(context).unfocus();
                                  }
                                  FocusManager.instance.primaryFocus?.unfocus();
                                  controller.handleTabChange(i);
                                },
                                labelPadding: const EdgeInsets.symmetric(
                                    horizontal: 12, vertical: 0),
                                isScrollable: true,
                                labelColor: ColorName.white,
                                unselectedLabelColor: ColorName.grey80,
                                labelStyle: const TextStyle(fontSize: 14.0),
                                indicatorColor: ColorName.grey80,
                                tabs: [
                                  Tab(
                                    text: LocaleKeys.openOrders.tr,
                                  ),
                                  Tab(
                                    text: LocaleKeys.orderHistory.tr,
                                  ),
                                ],
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),
                    if (activeTabIndex == 0)
                      OpenOrdersView(
                        fullScreen: fullScreen,
                      )
                    else if (activeTabIndex == 1)
                      Expanded(
                        child: OrderHistoryView(
                          fullScreen: fullScreen,
                          isFromExchange: false,
                        ),
                      )
                  ],
                ),
              );
            },
          ),
        ),
      ),
    );
  }
}
