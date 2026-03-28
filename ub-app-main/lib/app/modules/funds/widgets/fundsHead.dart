import 'dart:math';

import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../common/components/UBText.dart';
import '../pages/balance/controllers/balance_controller.dart';
import '../../../../generated/assets.gen.dart';
import '../../../../generated/colors.gen.dart';
import '../../../../generated/locales.g.dart';
import '../../../../utils/commonUtils.dart';
import '../../../../utils/mixins/formatters.dart';

const headerAnimationDuration = const Duration(milliseconds: 300);

class FundsHead extends GetView<BalanceController> with Formatter {
  @override
  Widget build(BuildContext context) {
    return Obx(() {
      final showAvailableData = controller.showAvailableData.value;
      final bitcoinInOrders = controller.balancesAllData.value.btcInOrderSum;
      final inOrdersSum = controller.balancesAllData.value.inOrderSum;
      final btcTotalSum = controller.balancesAllData.value.btcTotalSum;
      final totalSum = controller.balancesAllData.value.totalSum;
      final isHeadOpen = controller.isHeadOpen.value;
      return totalSum == null
          ? const SizedBox()
          : Container(
              width: double.infinity,
              padding: const EdgeInsets.only(top: 12, left: 12, right: 12),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  UBText(
                    text: LocaleKeys.estimatedBalance.tr + ':',
                    size: 11.0,
                    color: ColorName.grey80,
                  ),
                  const SizedBox(
                    height: 4,
                  ),
                  Row(
                    children: [
                      UBText(
                        text: (showAvailableData
                            ? decimalCoin(value: btcTotalSum)
                            : '**********'),
                        size: 17,
                        color: ColorName.white,
                        weight: FontWeight.w700,
                      ),
                      const SizedBox(
                        width: 4,
                      ),
                      UBText(
                        text: 'BTC',
                        size: 17,
                        color: ColorName.white,
                      ),
                      const SizedBox(
                        width: 4,
                      ),
                      UBText(
                        text: '~',
                        size: 12,
                        color: ColorName.grey80,
                      ),
                      const SizedBox(
                        width: 4,
                      ),
                      UBText(
                        text: (showAvailableData == true
                            ? '\$' +
                                decimalCoin(value: totalSum, coinCode: "USDT")
                            : '\$**********'),
                        size: 12.0,
                        color: ColorName.grey80,
                      ),
                      const Spacer(),
                      GestureDetector(
                        onTap: () => controller.handelShowAvailableDataToggle(),
                        child: Container(
                          height: 24,
                          width: 24,
                          child: Assets.images.eye.svg(
                              color: showAvailableData
                                  ? ColorName.primaryBlue
                                  : ColorName.grey80),
                        ),
                      ),
                      const SizedBox(
                        width: 12,
                      ),
                      GestureDetector(
                        onTap: () => controller.toggleHeadOpen(),
                        child: Transform.rotate(
                          angle: isHeadOpen ? (pi) : 0,
                          child: Container(
                            height: 24,
                            width: 24,
                            child: Assets.images.doubleArrowDown.svg(),
                          ),
                        ),
                      )
                    ],
                  ),
                  if (isHeadOpen)
                    const SizedBox(
                      height: 12,
                    ),
                  if (isHeadOpen)
                    UBText(
                      text: LocaleKeys.inOrders.tr + ':',
                      size: 11.0,
                      color: ColorName.grey80,
                    ),
                  if (isHeadOpen)
                    const SizedBox(
                      height: 6,
                    ),
                  if (isHeadOpen)
                    Row(
                      children: [
                        UBText(
                          text: (showAvailableData
                              ? decimalCoin(value: bitcoinInOrders)
                              : '**********'),
                          size: 17.0,
                          color: ColorName.white,
                        ),
                        const SizedBox(
                          width: 4,
                        ),
                        UBText(
                          text: 'BTC',
                          size: 17.0,
                          color: ColorName.white,
                        ),
                        const SizedBox(
                          width: 4,
                        ),
                        UBText(
                          text: '~',
                          size: 12.0,
                          color: ColorName.grey80,
                        ),
                        const SizedBox(
                          width: 4,
                        ),
                        UBText(
                          text: (showAvailableData == true
                                  ? '\$' +
                                      decimalCoin(
                                          value: inOrdersSum, coinCode: "USDT")
                                  : '\$**********') +
                              '',
                          size: 12.0,
                          color: ColorName.grey80,
                        ),
                        const Spacer(),
                      ],
                    )
                ],
              ),
            );
    });
  }
}
