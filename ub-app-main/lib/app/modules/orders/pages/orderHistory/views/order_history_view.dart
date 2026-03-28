import 'package:flutter/material.dart';

import 'package:get/get.dart';
import 'package:unitedbit/app/common/components/UBText.dart';
import 'package:unitedbit/generated/assets.gen.dart';
import '../../../../../common/components/CenterUBLoading.dart';
import '../../../../../common/components/UBButton.dart';
import '../../../../../common/components/UBDarkOpacityBackgrounded.dart';
import '../../../../../common/components/UBScrollBar.dart';
import '../../../../../common/components/UBoops.dart';
import '../widgets/orderHistoryFilterPopup.dart';
import '../../../widgets/orderRow/orderHistoryRow.dart';
import '../../../../../../generated/colors.gen.dart';
import '../../../../../../utils/mixins/commonConsts.dart';
import '../../../../../../utils/mixins/filterPopups.dart';
import '../../../../../../utils/mixins/formatters.dart';
import '../../../../../../utils/throttle.dart';
import '../controllers/order_history_controller.dart';

/*
Debouncer	Wait for changes to stop before notifying.
Throttle	Notifies once per Duration for a value that keeps changing.

 */
final thr = Throttling(duration: const Duration(milliseconds: 200));
//final deb = Debouncing(duration: const Duration(seconds: 2));

class OrderHistoryView extends GetView<OrderHistoryController>
    with FilterPopups, Formatter {
  final bool fullScreen;
  final bool isFromExchange;

  OrderHistoryView({this.isFromExchange = false, this.fullScreen});

  @override
  Widget build(BuildContext context) {
    controller.showFilterButton.value = true;
    var _isFromExchange;
    var _isFullScreen;
    if (Get.arguments != null) {
      _isFromExchange = Get.arguments['isFromExchange'];
      _isFullScreen = Get.arguments['fullScreen'];
    } else {
      _isFromExchange = false;
      _isFullScreen = false;
    }

    double lastScrollPosition = 0.0;
    final ScrollController scrollController = ScrollController();
    listener() {
      double newPosition = scrollController.offset;

      if (newPosition >= scrollController.position.maxScrollExtent &&
          !scrollController.position.outOfRange) {
        controller.onListEndRiched();
      }

      if (newPosition < 10.0) {
        controller.handleListScroll(direction: ScrollDirection.Up);
        return;
      }
      thr.throttle(() {
        if (newPosition > lastScrollPosition) {
          controller.handleListScroll(direction: ScrollDirection.Down);
        } else {
          controller.handleListScroll(direction: ScrollDirection.Up);
        }
        lastScrollPosition = newPosition;
      });
    }

    scrollController.removeListener(listener);
    scrollController.addListener(listener);
    return Obx(
      () {
        final loadingId = controller.loadingId.value;
        final hideCanceledOrders = !(controller.showCanceledOrders.value);
        final orderHistory = hideCanceledOrders
            ? controller.orderHistory
                .where((element) => element.status != 'canceled')
                .toList()
            : controller.orderHistory;
        final exchangeHistory = hideCanceledOrders
            ? controller.orderHistory
                .where((element) =>
                    element.status != 'canceled' &&
                    element.type == 'fast exchange')
                .toList()
            : controller.orderHistory
                .where((element) => element.type == 'fast exchange')
                .toList();
        final loadingData = controller.loadingData.value;
        final isSilentLoading = controller.silentLoadingData.value;
        final filtered = controller.filtered.value;
        return loadingData == true
            ? Container(
                child: CenterUbLoading(),
              )
            : orderHistory.length == 0 ||
                    (_isFromExchange && exchangeHistory.length == 0)
                ? Stack(
                    children: [
                      Column(
                        children: [
                          UBoops(
                            variant: filtered
                                ? OopsVariant.NoFilterResultsOops
                                : OopsVariant.OrderHistoryOops,
                          ),
                          if (filtered == true) vspace24,
                          if (filtered == true)
                            SizedBox(
                              width: 99,
                              child: UBButton(
                                fontSize: 13.0,
                                buttonColor: ColorName.black1c,
                                variant: ButtonVariant.Rounded,
                                textColor: ColorName.primaryBlue,
                                height: 32.0,
                                onClick: () {
                                  controller.handleResetFiltersClick(
                                      andPop: false);
                                },
                                text: 'Reset Filters',
                              ),
                            )
                        ],
                      ),
                      if (isSilentLoading)
                        UBDarkOpacityBackgrounded(
                          child: CenterUbLoading(),
                        ),
                      if (fullScreen == true || _isFullScreen == true)
                        ButtomFilterButton(),
                    ],
                  )
                : Stack(
                    children: [
                      Column(
                        children: [
                          Visibility(
                            visible: _isFromExchange || isFromExchange,
                            child: Padding(
                              padding: EdgeInsets.only(
                                  top: MediaQuery.of(context).padding.top),
                              child: Container(
                                height: 50,
                                decoration: BoxDecoration(
                                  border: Border(
                                    bottom: BorderSide(color: ColorName.grey36),
                                  ),
                                ),
                                //width: MediaQuery.of(context).size.width,
                                child: new Stack(
                                  //  Stack places the objects in the upper left corner
                                  children: <Widget>[
                                    Align(
                                      alignment: Alignment.centerLeft,
                                      child: GestureDetector(
                                          onTap: () => Get.back(),
                                          child: Padding(
                                            padding: const EdgeInsets.symmetric(
                                                horizontal: 12),
                                            child: Assets.images.appBarBack.svg(
                                              color: ColorName.greybf,
                                            ),
                                          )),
                                    ),
                                    Align(
                                      alignment: Alignment.center,
                                      child: UBText(
                                        text: 'Exchange History',
                                        color: ColorName.white,
                                        size: 17,
                                      ),
                                    ),
                                  ],
                                ),
                              ),
                            ),
                          ),
                          Expanded(
                            child: MediaQuery.removePadding(
                              context: context,
                              removeTop: true,
                              child: UBScrollBar(
                                scrollDirection: Axis.vertical,
                                builder: (BuildContext context, int index) {
                                  final order =
                                      _isFromExchange || isFromExchange
                                          ? exchangeHistory[index]
                                          : orderHistory[index];
                                  return OrderHistoryRow(
                                    order: order,
                                    onDetailsClick:
                                        controller.hadleOrderDetailsClick,
                                    isLoading: order.id == loadingId,
                                    formatter: currencyFormatter,
                                  );
                                },
                                itemCount: _isFromExchange || isFromExchange
                                    ? exchangeHistory.length
                                    : orderHistory.length,
                                scrollController: scrollController,
                              ),
                            ),
                          )
                        ],
                      ),
                      if (isSilentLoading)
                        UBDarkOpacityBackgrounded(
                          child: CenterUbLoading(),
                        ),
                      if (_isFullScreen == true || fullScreen == true)
                        ButtomFilterButton(),
                    ],
                  );
      },
    );
  }
}

class ButtomFilterButton extends GetView<OrderHistoryController> {
  @override
  Widget build(BuildContext context) {
    return Positioned(
      bottom: 36.0,
      left: (Get.width / 2) - 50.0,
      child: Obx(
        () {
          final showFilterButton = controller.showFilterButton.value;
          return AnimatedSwitcher(
            duration: 100.milliseconds,
            transitionBuilder: (Widget child, Animation<double> animation) {
              return ScaleTransition(child: child, scale: animation);
            },
            child: showFilterButton
                ? Container(
                    key: ValueKey('q1'),
                    width: 85,
                    height: 28,
                    child: UBButton(
                        height: 28,
                        variant: ButtonVariant.Rounded,
                        buttonColor: ColorName.grey42,
                        endWidget: const Icon(
                          Icons.filter_alt_sharp,
                          color: ColorName.greybf,
                          size: 14.0,
                        ),
                        fontSize: 14.0,
                        onClick: () {
                          Get.bottomSheet(OrderHistoryFilterPopup());
                        },
                        text: 'Filter'),
                  )
                : AbsorbPointer(
                    child: Container(
                      key: ValueKey('q2'),
                      width: 85,
                      height: 28,
                      color: Colors.transparent,
                    ),
                  ),
          );
        },
      ),
    );
  }
}
