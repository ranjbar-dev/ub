import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../app/common/components/UBText.dart';
import '../../app/modules/orders/pages/orderHistory/widgets/orderHistoryFilterPopup.dart';
import '../../generated/assets.gen.dart';
import '../../generated/colors.gen.dart';
import '../../generated/locales.g.dart';

mixin FilterPopups {
  openOrdersFilterSelect(
      {@required Function(String) onFilterSelect, String text}) {
    return GestureDetector(
      onTap: () {
        openOpenOrdersFilter(
            onFilterSelect: onFilterSelect, selectedText: text);
      },
      child: Container(
        padding: EdgeInsets.symmetric(horizontal: 12),
        height: 26,
        child: Row(
          children: [
            UBText(
              text: text,
            ),
            const SizedBox(width: 8),
            Assets.images.roundedKeyDown.svg()
          ],
        ),
      ),
    );
  }

  orderHistoryFilterButton({String text}) {
    return GestureDetector(
      onTap: () {
        Get.bottomSheet(OrderHistoryFilterPopup());
      },
      child: Container(
        padding: EdgeInsets.symmetric(horizontal: 12),
        color: ColorName.black,
        height: 26,
        child: Row(
          children: [
            UBText(
              text: text,
            ),
            const SizedBox(width: 8),
            Assets.images.roundedKeyDown.svg()
          ],
        ),
      ),
    );
  }

  openOpenOrdersFilter(
      {@required Function(String) onFilterSelect,
      Function afterClose,
      String selectedText}) {
    Get.bottomSheet(
      Container(
        padding: const EdgeInsets.all(12),
        decoration: BoxDecoration(
          color: ColorName.black2c,
          borderRadius: const BorderRadius.only(
            topLeft: const Radius.circular(16),
            topRight: const Radius.circular(16),
          ),
        ),
        height: 170,
        child: Column(
          children: [
            openOrderFilterOption(
              onTap: () {
                onFilterSelect(LocaleKeys.allOpenOrders.tr);
                Get.back();
              },
              text: LocaleKeys.allOpenOrders.tr,
              selectedText: selectedText,
            ),
            openOrderFilterOption(
              onTap: () {
                onFilterSelect(LocaleKeys.onlySellOrders.tr);
                Get.back();
              },
              text: LocaleKeys.onlySellOrders.tr,
              selectedText: selectedText,
            ),
            openOrderFilterOption(
              onTap: () {
                onFilterSelect(LocaleKeys.onlyBuyOrders.tr);
                Get.back();
              },
              text: LocaleKeys.onlyBuyOrders.tr,
              selectedText: selectedText,
            ),
            GestureDetector(
              onTap: () {
                Get.back();
              },
              child: Container(
                height: 36,
                child: Center(
                  child: UBText(
                    text: LocaleKeys.cancel.tr,
                    weight: FontWeight.bold,
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

openOrderFilterOption({String text, Function onTap, String selectedText}) {
  final isSelected = text == selectedText;
  return GestureDetector(
    onTap: onTap,
    child: Container(
      height: 36,
      decoration: BoxDecoration(
        border: Border(
          bottom: BorderSide(
            color: ColorName.grey36,
            width: 1.0,
          ),
        ),
      ),
      child: Center(
        child: UBText(
          text: text,
          weight: FontWeight.bold,
          color: isSelected ? ColorName.textBlue : ColorName.white,
        ),
      ),
    ),
  );
}
