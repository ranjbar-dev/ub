import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../../../common/components/CenterUBLoading.dart';
import '../../../../../common/components/CoinList.dart';
import '../../../../../common/components/UBCounSearchHistory.dart';
import '../../../../../common/components/UBDDMockButton.dart';
import '../../../../../common/components/appbarTextTitle.dart';
import '../../../../../common/components/pageContainer.dart';
import '../../../../../../generated/colors.gen.dart';
import '../../../../../../generated/locales.g.dart';
import '../../../../../../services/storageKeys.dart';
import '../../../../../../utils/mixins/commonConsts.dart';
import '../../../../../../utils/mixins/formatters.dart';
import '../../../../../../utils/mixins/popups.dart';
import '../../../../../../utils/mixins/toast.dart';

import '../controllers/withdrawals_controller.dart';

class WithdrawalsView extends GetView<WithdrawalsController>
    with Formatter, Toaster, Popups {
  @override
  Widget build(BuildContext context) {
    return PageContainer(
      appbarTitle: AppBarTextTitle(
        title: LocaleKeys.withdrawals.tr,
      ),
      child: Stack(
        children: [
          Column(
            children: [
              //_coinSelector(),
              Container(
                  margin: EdgeInsets.only(
                    top: 24,
                    bottom: 16,
                  ),
                  child: UBDDMockButton(
                    backgroundColor: ColorName.black2c,
                    endIcon: const Icon(
                      Icons.search,
                      color: ColorName.greybf,
                      size: 20,
                    ),
                    title: LocaleKeys.pleaseselectanycoin.tr,
                    onTap: () => openCoinSelectPopup(
                      onCoinSelect: controller.handleCoinSelected,
                      afterClose: () => _openDetailsPopup(),
                    ),
                  )),
              Obx(() {
                //we need this for obx to work
                // ignore: invalid_use_of_protected_member , unused_local_variable
                final coins = controller.savedCoins.value;
                return CoinSearchHistory(
                  key: UniqueKey(),
                  stream: controller.savedCoins,
                  storageKey: StorageKeys.savedWithdrawalCoins,
                  onCoinClick: (coin) {
                    controller.handleCoinSelected(coin);
                    _openDetailsPopup();
                  },
                );
              }),
              Obx(() {
                // ignore: invalid_use_of_protected_member
                final searchedCoins = controller.savedCoins.value;
                return searchedCoins.length > 0 ? vspace12 : emptyComponent;
              }),
              CoinsList(
                onSelect: (coin) {
                  controller.handleCoinSelected(coin);
                  _openDetailsPopup();
                },
              )
            ],
          ),
          Obx(() {
            final isLoading = controller.isLoadingWithdrawAndDepositData.value;
            return isLoading
                ? Container(
                    height: Get.height,
                    width: Get.width,
                    color: ColorName.black.withOpacity(0.2),
                    child: CenterUbLoading(),
                  )
                : emptyComponent;
          }),
        ],
      ),
    );
  }

  //Container _coinSelector() {
  //  return Container(
  //    margin: EdgeInsets.only(
  //      top: 24,
  //      bottom: 24,
  //    ),
  //    child: Obx(
  //      () {
  //        final isLoading = controller.isLoadingWithdrawAndDepositData.value;
  //        final selectedCoin = controller.selectedCoin.value;
  //        if (isLoading == true) {
  //          return UBGreyContainer(
  //            color: ColorName.black2c,
  //            margin: const EdgeInsets.symmetric(horizontal: 12),
  //            child: Row(
  //              mainAxisAlignment: MainAxisAlignment.spaceBetween,
  //              children: [
  //                SizedBox(
  //                  width: 24,
  //                  height: 24,
  //                  child: CircularProgressIndicator(),
  //                ),
  //                const SizedBox()
  //              ],
  //            ),
  //          );
  //        }
  //        return UBDDMockButton(
  //          backgroundColor: ColorName.black2c,
  //          title: selectedCoin.id == null
  //              ? LocaleKeys.pleaseselectanycoin.tr
  //              : selectedCoin.name,
  //          titleAppendix: selectedCoin.desc,
  //          iconAddress: selectedCoin.image,
  //          onTap: () => openCoinSelectPopup(
  //              onCoinSelect: controller.handleCoinSelected),
  //        );
  //      },
  //    ),
  //  );
  //}

  void _openDetailsPopup() {}
}
