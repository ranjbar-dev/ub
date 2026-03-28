import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../common/components/CenterUBLoading.dart';
import '../../../common/components/UBScaleSwitcher.dart';
import '../../../common/components/UBScrollColumnExpandable.dart';
import '../../../common/components/pageContainer.dart';
import '../../orders/views/orders_view.dart';
import '../controllers/trade_controller.dart';
import '../tradePageConfig.dart';
import 'widgets/ohlcChart.dart';
import 'widgets/head.dart';
import 'widgets/newOrder.dart';
import 'widgets/orderBook.dart';
import '../../../../generated/colors.gen.dart';

class TradeView extends GetView<TradeController> {
  @override
  Widget build(BuildContext context) {
    controller.checkForEssentialData();
    return PageContainer(
      activeBottomNavIndex: 2,
      // additionalHeight: 0,
      protectedPage: false,
      child: Stack(
        children: [
          UBScrollColumnExpandable(
            children: [
              TradeHead(),
              Obx(() {
                final activeChart = controller.activeChart.value;
                final isOhlc = activeChart == TradeTopCharts.OHLC;
                return SizedBox(
                  height: getChartHeight(isOhlc: isOhlc),
                  width: double.infinity,
                  child: UBScaleSwitcher(
                    child1: OHLCChart(),
                    child2: OrderBook(),
                    conditionToShowChild1: isOhlc,
                  ),
                );
              }),
              Obx(() {
                final activeChart = controller.activeChart.value;
                if (activeChart == TradeTopCharts.OrderBook) {
                  return Container(
                    height: 4,
                    color: ColorName.black,
                  );
                }
                return SizedBox();
              }),
              // Expanded(
              //   flex: 1,
              //   child: OHLCChart(),
              // ),
              Obx(
                () => AnimatedContainer(
                  curve: Curves.easeInOut,
                  duration: const Duration(milliseconds: 300),
                  height: sizes[controller.subActiveIndex.value],
                  child: NewOrder(sizes: sizes),
                ),
              ),
              Expanded(
                child: Obx(
                  () {
                    final loggedIn = controller.globalController.loggedIn.value;
                    return loggedIn == true
                        ? OrdersView()
                        : Container(
                            color: Colors.blueGrey,
                          );
                  },
                ),
              ),
            ],
          ),
          Obx(() {
            final showLoadingOverlay = controller.showLoadingOverlay.value;
            return showLoadingOverlay == false
                ? const SizedBox()
                : Column(
                    children: [
                      Expanded(
                          child: Container(
                        color: ColorName.black.withOpacity(0.2),
                        child: CenterUbLoading(),
                      )),
                    ],
                  );
          })
        ],
      ),
    );
  }
}
