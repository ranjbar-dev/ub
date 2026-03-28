import 'package:flutter/material.dart';

import 'package:get/get.dart';
import '../../../../../common/components/CenterUBLoading.dart';
import '../../../../../common/components/UBDarkOpacityBackgrounded.dart';
import '../../../../../common/components/UBScrollBar.dart';
import '../../../../../common/components/UBoops.dart';
import '../../../widgets/orderRow/openOrderRow.dart';
import '../../../../../../utils/mixins/commonConsts.dart';
import '../../../../../../utils/mixins/filterPopups.dart';

import '../controllers/open_orders_controller.dart';

class OpenOrdersView extends GetView<OpenOrdersController> with FilterPopups {
  final bool fullScreen;

  OpenOrdersView({Key key, this.fullScreen});

  @override
  Widget build(BuildContext context) {
    final ScrollController scrollController = ScrollController();

    //final key = controller.keySelector(fullScreen ?? false);
    return Obx(() {
      // ignore: invalid_use_of_protected_member
      final loadingIds = controller.loadingIds.value;
      final openOrders = controller.openOrders;
      final handleCancelClick = controller.handleCancelClick;
      final filterText = controller.selectedOpenOrderFilterText.value;
      return Expanded(
        child: controller.loadingData.value == true
            ? Container(
                child: CenterUbLoading(),
              )
            : controller.openOrders.length == 0
                ? Stack(
                    children: [
                      Column(
                        children: [
                          if (fullScreen == true &&
                              filterText != 'All Open Orders')
                            openOrdersFilterSelect(
                              onFilterSelect: controller.filterOpenOrders,
                              text:
                                  controller.selectedOpenOrderFilterText.value,
                            ),
                          Expanded(
                            child: UBoops(
                              variant: OopsVariant.OpenOrderOops,
                            ),
                          ),
                        ],
                      ),
                      Obx(() {
                        final isSilentLoading =
                            controller.isSilentLoading.value;
                        if (isSilentLoading == true) {
                          return UBDarkOpacityBackgrounded(
                              child: CenterUbLoading());
                        }
                        return emptyComponent;
                      })
                    ],
                  )
                : Stack(
                    children: [
                      Column(
                        children: [
                          if (fullScreen == true)
                            openOrdersFilterSelect(
                              onFilterSelect: controller.filterOpenOrders,
                              text:
                                  controller.selectedOpenOrderFilterText.value,
                            ),
                          Expanded(
                              child: UBScrollBar(
                            scrollController: scrollController,
                            itemCount: openOrders.length,
                            builder: (
                              BuildContext context,
                              int index,
                            ) {
                              if (openOrders[index] != null) {
                                return OpenOrderRow(
                                  onCancelClick: (id) =>
                                      handleCancelClick(id: id, index: index),
                                  isCancelLoading:
                                      loadingIds.contains(openOrders[index].id),
                                  order: openOrders[index],
                                );
                              }
                              return SizedBox();
                            },
                          )),

                          //AnimatedList(
                          //  shrinkWrap: true,
                          //  key: key,
                          //  initialItemCount: openOrders.length,
                          //  itemBuilder: (BuildContext context, int index,
                          //      Animation animation) {
                          //    if (openOrders[index] != null) {
                          //      return SizeTransition(
                          //        sizeFactor: animation,
                          //        child: FadeTransition(
                          //          opacity: animation,
                          //          child: OpenOrderRow(
                          //            onCancelClick: (id) => handleCancelClick(
                          //                id: id, index: index),
                          //            isCancelLoading: loadingIds
                          //                .contains(openOrders[index].id),
                          //            order: openOrders[index],
                          //          ),
                          //        ),
                          //      );
                          //    }
                          //    return SizedBox();
                          //  },
                          //),
                        ],
                      ),
                      Obx(() {
                        final isSilentLoading =
                            controller.isSilentLoading.value;
                        if (isSilentLoading == true) {
                          return UBDarkOpacityBackgrounded(
                            child: CenterUbLoading(),
                          );
                        }
                        return emptyComponent;
                      })
                    ],
                  ),
      );
    });
  }
}
