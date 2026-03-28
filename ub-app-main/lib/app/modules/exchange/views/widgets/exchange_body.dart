import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:supercharged/supercharged.dart';

import '../../../../../generated/assets.gen.dart';
import '../../../../../generated/colors.gen.dart';
import '../../../../../utils/extentions/basic.dart';
import '../../../../../utils/mixins/commonConsts.dart';
import '../../../../common/components/UBButton.dart';
import '../../../../common/components/UBShimmer.dart';
import '../../../../common/components/UBText.dart';
import '../../controllers/exchange_controller.dart';
import 'exchange_drop_down.dart';
import 'exchange_input.dart';
import 'top_markets.dart';

class ExchangeBody extends GetView<ExchangeController> {
  @override
  Widget build(BuildContext context) {
    //controller.getPairBalances();
    return SingleChildScrollView(
      child: Padding(
        padding: const EdgeInsets.symmetric(horizontal: 12.0),
        child: Column(
          children: [
            AspectRatio(
              aspectRatio: 16 / 7,
              child: Assets.images.exchangeBanner.image(),
            ),
            vspace8,
            Container(
              height: ((MediaQuery.of(context).size.height / 3) - 25),
              decoration: BoxDecoration(
                color: ColorName.black2c,
                borderRadius: BorderRadius.circular(15),
              ),
              child: Column(
                children: [
                  Expanded(
                    flex: 30,
                    child: Padding(
                      padding: const EdgeInsets.all(12),
                      child: Row(
                        //mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Obx(
                            () => Expanded(
                              flex: 46,
                              child: ExchangeDropDown(
                                name: controller
                                    .pairLocalInfo.basisCoin.value.code,
                                desc: controller
                                    .pairLocalInfo.basisCoin.value.desc,
                                imageUrl: controller
                                    .pairLocalInfo.basisCoin.value.image,
                                isFrom: true,
                              ),
                            ),
                          ),
                          Expanded(
                            flex: 8,
                            child: Center(
                              child: UBText(
                                text: 'To',
                                size: 13,
                                color: ColorName.grey97,
                              ),
                            ),
                          ),
                          Obx(
                            () => Expanded(
                              flex: 46,
                              child: ExchangeDropDown(
                                name: controller
                                    .pairLocalInfo.dependantCoin.value.code,
                                desc: controller
                                    .pairLocalInfo.dependantCoin.value.desc,
                                imageUrl: controller
                                    .pairLocalInfo.dependantCoin.value.image,
                                isFrom: false,
                              ),
                            ),
                          )
                        ],
                      ),
                    ),
                  ),
                  Expanded(
                    flex: 14,
                    child: Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 12),
                      child: Row(
                        //mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Expanded(
                            flex: 46,
                            child: UBText(
                              text: 'You Will Give',
                              size: 12,
                              color: ColorName.greybf,
                            ),
                          ),
                          Expanded(
                            flex: 8,
                            child: SizedBox(),
                          ),
                          Expanded(
                            flex: 46,
                            child: UBText(
                              text: 'You Will Get ~',
                              size: 12,
                              color: ColorName.greybf,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                  Expanded(
                    flex: 22,
                    child: Align(
                      alignment: Alignment.center,
                      child: Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 12.0),
                        child: Stack(
                          children: [
                            Positioned(
                              left: 0,
                              child: Container(
                                width: (MediaQuery.of(context).size.width / 2) -
                                    28,
                                child: ExchangeInput(
                                  hasBadge: true,
                                ),
                              ),
                            ),
                            Positioned(
                              right: 0,
                              child: Container(
                                width: (MediaQuery.of(context).size.width / 2) -
                                    28,
                                child: ExchangeInput(
                                  hasBadge: false,
                                ),
                              ),
                            ),
                            Positioned(
                              left:
                                  (MediaQuery.of(context).size.width / 2) - 49,
                              right:
                                  (MediaQuery.of(context).size.width / 2) - 49,
                              child: Container(
                                width: 43.0,
                                height: 43.0,
                                decoration: BoxDecoration(
                                  color: '34343D'.toColor(),
                                  shape: BoxShape.circle,
                                  border: Border.all(
                                    color: '525261'.toColor(),
                                  ),
                                ),
                                child: Padding(
                                  padding: const EdgeInsets.all(6.0),
                                  child: InkWell(
                                    child: Assets.images.swap.svg(),
                                    onTap: () => controller.swapCoins(),
                                  ),
                                ),
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),
                  ),
                  Expanded(
                    flex: 14,
                    child: Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 12),
                      child: Align(
                        alignment: Alignment.topCenter,
                        child: Row(
                            //mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              Expanded(
                                flex: 46,
                                child: Row(
                                  mainAxisAlignment: MainAxisAlignment.start,
                                  children: [
                                    UBText(
                                      text: 'Aval: ',
                                      size: 12,
                                      color: ColorName.grey80,
                                    ),
                                    Obx(
                                      () => controller
                                              .isLoadingBalanceData.value
                                          ? ClipRRect(
                                              borderRadius: BorderRadius.all(
                                                Radius.circular(8),
                                              ),
                                              child: UBShimmer(
                                                width: (MediaQuery.of(context)
                                                            .size
                                                            .width /
                                                        2) -
                                                    70,
                                                height: 15,
                                                background: ColorName.black2c,
                                              ),
                                            )
                                          : UBText(
                                              text: controller.pairLocalInfo
                                                  .basisBalance.value
                                                  .toStringAsPrecision(
                                                      controller.pairLocalInfo
                                                          .pairPrecision.value)
                                                  .currencyFormat(
                                                      centFormat: true),
                                              size: 12,
                                              color: ColorName.greyd8,
                                            ),
                                    )
                                  ],
                                ),
                              ),
                              Expanded(
                                flex: 8,
                                child: SizedBox(),
                              ),
                              Expanded(
                                flex: 46,
                                child: Row(
                                  mainAxisAlignment: MainAxisAlignment.start,
                                  children: [
                                    UBText(
                                      text: 'Aval: ',
                                      size: 12,
                                      color: ColorName.grey80,
                                    ),
                                    Obx(
                                      () => controller
                                              .isLoadingBalanceData.value
                                          ? ClipRRect(
                                              borderRadius: BorderRadius.all(
                                                Radius.circular(8),
                                              ),
                                              child: UBShimmer(
                                                width: (MediaQuery.of(context)
                                                            .size
                                                            .width /
                                                        2) -
                                                    70,
                                                height: 15,
                                                background: ColorName.black2c,
                                              ),
                                            )
                                          : UBText(
                                              text: controller.pairLocalInfo
                                                  .dependentBalance.value
                                                  .toStringAsPrecision(
                                                      controller.pairLocalInfo
                                                          .pairPrecision.value)
                                                  .currencyFormat(
                                                      centFormat: true),
                                              size: 12,
                                              color: ColorName.greyd8,
                                            ),
                                    ),
                                  ],
                                ),
                              ),
                            ]),
                      ),
                    ),
                  ),
                  Expanded(
                    flex: 30,
                    child: Padding(
                      padding: const EdgeInsets.all(14.0),
                      child: Obx(
                        () => UBButton(
                          onClick: () {
                            controller.handleExchangeClick();
                          },
                          height: 40,
                          isLodaing: controller.isLoadingExchangeSubmit.value,
                          text: 'Exchange',
                          fontSize: 18,
                          borderRadius: 6.0,
                        ),
                      ),
                    ),
                  )
                ],
              ),
            ),
            vspace8,
            SafeArea(child: TopMarkets())
          ],
        ),
      ),
    );
  }
}
