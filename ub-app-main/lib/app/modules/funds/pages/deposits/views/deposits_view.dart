import 'package:flutter/material.dart';

import 'package:get/get.dart';
import '../../../../../common/components/CoinList.dart';
import '../../../../../common/components/UBCounSearchHistory.dart';
import '../../../../../common/components/UBDDMockButton.dart';
import '../../../../../common/components/appbarTextTitle.dart';
import '../../../../../common/components/pageContainer.dart';
import '../../../../../global/autocompleteModel.dart';
import 'depositDetails.dart';
import '../../../../../../generated/colors.gen.dart';
import '../../../../../../generated/locales.g.dart';
import '../../../../../../services/storageKeys.dart';
import '../../../../../../utils/mixins/commonConsts.dart';
import '../../../../../../utils/mixins/popups.dart';

import '../controllers/deposits_controller.dart';

class DepositsView extends GetView<DepositsController> with Popups {
  @override
  Widget build(BuildContext context) {
    return PageContainer(
      appbarTitle: AppBarTextTitle(
        title: LocaleKeys.deposits.tr,
      ),
      child: Stack(
        children: [
          Column(
            children: [
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
                  title: LocaleKeys.selectCoin.tr,
                  onTap: () => openCoinSelectPopup(onCoinSelect: (coin) {
                    controller.handleCoinSelected(coin);
                  }, afterClose: () {
                    _openDetailsPopup(controller.selectedCoin.value);
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
                    controller.handleCoinSelected(coin);
                    _openDetailsPopup(coin);
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
                  _openDetailsPopup(coin);
                },
              )
            ],
          ),
          //FloatingSearchBar(
          //  hint: 'Search currencies',
          //  height: 32,
          //  backgroundColor: ColorName.black2c,
          //  iconColor: ColorName.grey97,

          //  hintStyle: TextStyle(
          //    color: ColorName.greybf,
          //    fontSize: 11.0,
          //    fontWeight: FontWeight.w600,
          //  ),
          //  automaticallyImplyBackButton: false,
          //  scrollPadding: const EdgeInsets.only(top: 16, bottom: 56),
          //  transitionDuration: const Duration(milliseconds: 200),
          //  transitionCurve: Curves.easeInOut,
          //  physics: const BouncingScrollPhysics(),
          //  axisAlignment: 0.0,
          //  openAxisAlignment: 0.0,
          //  debounceDelay: const Duration(milliseconds: 500),
          //  onQueryChanged: (query) {
          //    // Call your model, bloc, controller here.
          //  },
          //  // Specify a custom transition to be used for
          //  // animating between opened and closed stated.
          //  transition: CircularFloatingSearchBarTransition(),
          //  actions: [
          //    //FloatingSearchBarAction(
          //    //  showIfOpened: false,
          //    //  child: CircularButton(
          //    //    icon: const Icon(Icons.search),
          //    //    onPressed: () {},
          //    //  ),
          //    //),
          //    FloatingSearchBarAction.searchToClear(
          //      showIfClosed: true,
          //    ),
          //  ],
          //  builder: (context, transition) {
          //    return ClipRRect(
          //      borderRadius: BorderRadius.circular(8),
          //      child: Material(
          //        color: ColorName.black,
          //        elevation: 0.0,
          //        child: Column(
          //          mainAxisSize: MainAxisSize.min,
          //          children: Colors.accents.map((color) {
          //            return Container(height: 112, color: color);
          //          }).toList(),
          //        ),
          //      ),
          //    );
          //  },
          //),
        ],
      ),
    );
  }

  _openDetailsPopup(AutoCompleteItem coin) {
    Get.dialog(
      Container(
        margin: const EdgeInsets.all(12.0),
        decoration: BoxDecoration(
          borderRadius: rounded_big,
          border: Border.all(color: ColorName.black2c, width: 1),
        ),
        width: Get.width,
        child: DepostDetailsView(coin: coin),
      ),
    );
  }
}
