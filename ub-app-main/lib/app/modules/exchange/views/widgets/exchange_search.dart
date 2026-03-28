import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../../../generated/assets.gen.dart';
import '../../../../../utils/mixins/commonConsts.dart';

import '../../../../../../generated/colors.gen.dart';
import '../../../../../../services/storageKeys.dart';
import '../../../../../../utils/mixins/popups.dart';
import '../../../../common/components/UBCounSearchHistory.dart';
import '../../../../common/components/UBDDMockButton.dart';
import '../../../../common/components/UBText.dart';
import '../../controllers/exchange_controller.dart';
import 'exchange_drop_down_results.dart';

class ExchangeSearch extends GetView<ExchangeController> with Popups {
  const ExchangeSearch({Key key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        Column(
          children: [
            SizedBox(
              height: MediaQuery.of(context).padding.top,
            ),
            Container(
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
                          padding: const EdgeInsets.symmetric(horizontal: 12),
                          child: Assets.images.appBarBack.svg(
                            color: ColorName.greybf,
                          ),
                        )),
                  ),
                  Align(
                    alignment: Alignment.center,
                    child: UBText(
                      text: 'Select Coin',
                      color: ColorName.white,
                      size: 17,
                    ),
                  ),
                ],
              ),
            ),
            Container(
              margin: EdgeInsets.only(
                top: 10,
                bottom: 16,
              ),
              child: UBDDMockButton(
                backgroundColor: ColorName.black2c,
                endIcon: const Icon(
                  Icons.search,
                  color: ColorName.greybf,
                  size: 20,
                ),
                title: 'Select Coin',
                onTap: () => openCoinSelectPopup(onCoinSelect: (coin) {
                  controller.handleCoinSelected(
                      coin: coin,
                      isFrom: Get.arguments['isFrom'],
                      isSwap: false);
                }),
              ),
            ),
            Obx(() {
              //its here because we need obx to work
              // ignore: unused_local_variable , invalid_use_of_protected_member
              final coins = controller.savedCoins.value;
              return CoinSearchHistory(
                stream: controller.savedCoins,
                storageKey: StorageKeys.savedDepositCoins,
                onCoinClick: (coin) {
                  controller.handleCoinSelected(
                      coin: coin,
                      isFrom: Get.arguments['isFrom'],
                      isSwap: false);
                },
              );
            }),
            Obx(() {
              //ignore: invalid_use_of_protected_member
              final searchedCoins = controller.savedCoins.value;
              return searchedCoins.length > 0 ? vspace12 : emptyComponent;
            }),
            ExchangeDropDownResults(
              onSelect: (coin) {
                ["", null].contains(Get.arguments['autoExchangeCode'])
                    ? controller.handleCoinSelected(
                        coin: coin,
                        isFrom: Get.arguments['isFrom'],
                        isSwap: false)
                    : controller.handleAutoExchangeDependantCoinSelected(
                        coin: coin,
                      );
                Get.back();
              },
              isFrom: Get.arguments['isFrom'],
            )
          ],
        ),
      ],
    );
  }
}
