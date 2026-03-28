import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:unitedbit/app/global/autocompleteModel.dart';

import '../../../../../../generated/assets.gen.dart';
import '../../../../../../generated/colors.gen.dart';
import '../../../../../../generated/locales.g.dart';
import '../../../../../../utils/mixins/commonConsts.dart';
import '../../../../../../utils/mixins/popups.dart';
import '../../../../../common/components/CenterUBLoading.dart';
import '../../../../../common/components/UBDDMockButton.dart';
import '../../../../../common/components/UBText.dart';
import '../../../../../common/components/UBoops.dart';
import '../../../../exchange/controllers/exchange_controller.dart';
import '../../widgets/auto_exchange_row.dart';
import '../controllers/auto_exchange_controller.dart';

class AutoExchangeView extends GetView<AutoExchangeController> with Popups {
  @override
  Widget build(BuildContext context) {
    Get.put<ExchangeController>(ExchangeController(), permanent: true);
    Get.put<AutoExchangeController>(AutoExchangeController(), permanent: true);
    return Expanded(
      child: Container(
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 12.0),
          child: Column(
            children: [
              AspectRatio(
                aspectRatio: 3,
                child: Assets.images.autoExchangeBanner.image(),
              ),
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 16),
                child: Row(
                  children: [
                    Assets.images.bigInfoIcon.svg(),
                    hspace8,
                    UBText(
                        color: ColorName.greyd8,
                        size: 12,
                        weight: FontWeight.w400,
                        wrapped: true,
                        text:
                            'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Amet.')
                  ],
                ),
              ),
              vspace16,
              Container(
                  width: Get.width,
                  child: Obx(
                    () => UBDDMockButton(
                      horizontalPadding: 0,
                      backgroundColor: ColorName.black2c,
                      endIcon: controller.searchedCoin.value.name == ""
                          ? const Icon(
                              Icons.search,
                              color: ColorName.greybf,
                              size: 20,
                            )
                          : InkWell(
                              onTap: () async {
                                controller.getBalances();
                                controller.searchedCoin.value =
                                    AutoCompleteItem(name: "");
                              },
                              child: const Icon(
                                Icons.close,
                                color: ColorName.greybf,
                                size: 20,
                              ),
                            ),
                      title: controller.searchedCoin.value.name == ""
                          ? LocaleKeys.selectCoin.tr
                          : controller.searchedCoin.value.desc,
                      onTap: () => openCoinSelectPopup(onCoinSelect: (coin) {
                        Get.find<AutoExchangeController>()
                            .handleCoinSelected(coin: coin, isFromSearch: true);
                      }),
                    ),
                  )),
              vspace16,
              Align(
                alignment: Alignment.centerLeft,
                child: UBText(
                  text: 'Coin List',
                  size: 13,
                  color: ColorName.grey80,
                ),
              ),
              vspace10,
              Expanded(
                child: Obx(
                  () {
                    if (controller.balances == null &&
                        controller.isCoinsListLoading.value == false) {
                      return UBoops(
                        variant: OopsVariant.ErrorOops,
                      );
                    }
                    return controller.isCoinsListLoading.value == true
                        ? CenterUbLoading()
                        : ListView.builder(
                            //physics: const NeverScrollableScrollPhysics(),
                            shrinkWrap: true,
                            scrollDirection: Axis.vertical,
                            itemCount: controller.balances.length,
                            itemBuilder: (ctx, index) {
                              final balance = controller.balances[index];
                              return AutoExchangeRow(
                                balance: balance,
                              );
                            },
                          );
                  },
                ),
              )
            ],
          ),
        ),
      ),
    );
  }
}
