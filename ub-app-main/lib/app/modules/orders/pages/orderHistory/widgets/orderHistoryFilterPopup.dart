import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../../../common/components/UBButton.dart';
import '../../../../../common/components/UBDatePicker.dart';
import '../../../../../common/components/UBSection.dart';
import '../../../../../common/components/UBText.dart';
import '../../../../../common/components/UBWrappedButtons.dart';
import '../controllers/order_history_controller.dart';
import '../../../../../../generated/assets.gen.dart';
import '../../../../../../generated/colors.gen.dart';
import '../../../../../../utils/mixins/commonConsts.dart';
import '../../../../../../utils/mixins/popups.dart';

final _smallButtonHeight = 26.0;

class OrderHistoryFilterPopup extends GetView<OrderHistoryController>
    with Popups {
  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(
        vertical: 4,
        horizontal: 12,
      ),
      decoration: BoxDecoration(
        color: ColorName.black2c,
        borderRadius: const BorderRadius.only(
          topLeft: const Radius.circular(16),
          topRight: const Radius.circular(16),
        ),
      ),
      height: 400,
      child: Column(
        children: [
          headerWithCloseButton(
            title: 'Filter',
            centerTitle: true,
            noContentPadding: true,
          ),
          vspace12,
          UBSection(
            title: 'Date',
            hTitlePadding: 0.0,
            child: Container(
              child: Column(
                children: [
                  Obx(
                    () => UBWrappedButtons(
                      buttonHeight: _smallButtonHeight,
                      minButtonWidth: (Get.width / 4) - 18,
                      buttonBackground: ColorName.black1c,
                      buttons: controller.dateButtons,
                      onButtonClick: (idx) {
                        controller.handleDateButtonClick(index: idx);
                      },
                      selectedIndex: controller.selectedDateButtonIndex.value,
                    ),
                  ),
                  vspace4,
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Obx(() {
                        final value = controller.filterStartDate.value;
                        return UBDatePicker(
                          onClearClick: () {
                            controller.filterStartDate.value = '';
                          },
                          onDateSelect: controller.handleStartDateSelect,
                          width: ((Get.width / 2) - 20),
                          placeHolder: 'Start Date',
                          value: value,
                        );
                      }),
                      Obx(() {
                        final value = controller.filterEndDate.value;
                        return UBDatePicker(
                          onClearClick: () {
                            controller.filterEndDate.value = '';
                          },
                          onDateSelect: controller.handleEndDateSelect,
                          width: ((Get.width / 2) - 20),
                          placeHolder: 'End Date',
                          value: value,
                        );
                      }),
                    ],
                  )
                ],
              ),
            ),
          ),
          vspace12,
          UBSection(
            title: 'Pair',
            hTitlePadding: 0.0,
            child: Column(
              children: [
                vspace4,
                Container(
                  height: _smallButtonHeight,
                  child: Obx(
                    () {
                      final pair = controller.filterPair.value;
                      final hasValue = pair != 'all-all';
                      return AnimatedSwitcher(
                        duration: 200.milliseconds,
                        transitionBuilder:
                            (Widget child, Animation<double> animation) {
                          return ScaleTransition(
                              child: child, scale: animation);
                        },
                        child: !hasValue
                            ? SmallButton(
                                key: ValueKey('AllPairs'),
                                height: _smallButtonHeight,
                                minWidth: 100,
                                text: 'All Pairs',
                                onClick: () {
                                  openPairSelectPopup(onSelect: (v) {
                                    controller.handlePairSelect(v);
                                  });
                                },
                              )
                            : SmallButton(
                                height: _smallButtonHeight,
                                key: ValueKey('pair'),
                                text: pair,
                                onCloseClick: () {
                                  controller.filterPair.value = 'all-all';
                                },
                                onClick: () {
                                  openPairSelectPopup(onSelect: (v) {
                                    controller.handlePairSelect(v);
                                  });
                                },
                              ),
                      );
                    },
                  ),
                ),
              ],
            ),
          ),
          vspace12,
          UBSection(
            title: 'Type',
            hTitlePadding: 0.0,
            child: Container(
              height: 30,
              child: Obx(
                () => UBWrappedButtons(
                  buttonHeight: _smallButtonHeight,
                  minButtonWidth: 75.0,
                  buttonBackground: ColorName.black1c,
                  buttons: controller.filterTypeButtons,
                  onButtonClick: (idx) {
                    controller.handleTypeButtonClick(index: idx);
                  },
                  selectedIndex: controller.selecteTypeButtonIndex.value,
                ),
              ),
            ),
          ),
          vspace12,
          GestureDetector(
            onTap: () {
              controller.handleShowCanceledOrdersToggle();
            },
            child: Container(
              height: 24,
              color: ColorName.black2c,
              child: Obx(
                () {
                  final showCancelOrders = controller.showCanceledOrders.value;
                  return Row(
                    children: [
                      showCancelOrders
                          ? Assets.images.emptyCheckbox.svg()
                          : Assets.images.filledCheckbox.svg(),
                      hspace8,
                      UBText(
                        text: 'Hide all cancelled orders',
                        size: 14.0,
                        weight: FontWeight.w600,
                        color: showCancelOrders
                            ? ColorName.greybf
                            : ColorName.white,
                      )
                    ],
                  );
                },
              ),
            ),
          ),
          const Spacer(),
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              UBButton(
                onClick: () {
                  controller.handleResetFiltersClick();
                },
                width: (Get.width / 2) - 16,
                height: 36.0,
                buttonColor: ColorName.black1c,
                text: 'Reset',
              ),
              UBButton(
                onClick: () {
                  controller.handleFilterSubmitClick();
                },
                width: (Get.width / 2) - 16,
                height: 36.0,
                text: 'Apply',
              ),
            ],
          ),
          vspace8
        ],
      ),
    );
  }
}
