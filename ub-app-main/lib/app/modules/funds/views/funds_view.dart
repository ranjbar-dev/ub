import 'package:flutter/material.dart';

import 'package:get/get.dart';
import '../../../common/components/pageContainer.dart';
import '../pages/autoExchange/views/auto_exchange_view.dart';
import '../pages/balance/views/balance_view.dart';
import '../pages/transactionHistory/views/transaction_history_view.dart';
import '../widgets/DepositWithdrawButtons.dart';
import '../widgets/fundsHead.dart';
import '../../../../generated/colors.gen.dart';
import '../../../../generated/locales.g.dart';

import '../controllers/funds_controller.dart';

class FundsView extends GetView<FundsController> {
  @override
  Widget build(BuildContext context) {
    return PageContainer(
      activeBottomNavIndex: 4,
      child: Container(
        width: double.infinity,
        color: ColorName.black,
        child: Obx(
          () => DefaultTabController(
            length: 3,
            initialIndex: controller.activeTabIndex.value,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: <Widget>[
                AnimatedContainer(
                  curve: Curves.fastOutSlowIn,
                  //margin: const EdgeInsets.only(bottom: 24),
                  duration: const Duration(milliseconds: 300),
                  decoration: BoxDecoration(
                      color: ColorName.black2c,
                      borderRadius: BorderRadius.only(
                        bottomLeft: controller.isHeadOpen.value
                            ? const Radius.circular(12)
                            : const Radius.circular(0),
                        bottomRight: controller.isHeadOpen.value
                            ? const Radius.circular(12)
                            : const Radius.circular(0),
                      )),
                  height: controller.isHeadOpen.value ? 120 : 65,
                  child: FundsHead(),
                ),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Transform.translate(
                      offset: const Offset(0, 35.0),
                      child: Container(
                        width: double.infinity,
                        height: 1,
                        color: ColorName.black2c,
                      ),
                    ),
                    Container(
                      height: 35,
                      //width: 250,
                      child: TabBar(
                        isScrollable: true,
                        onTap: controller.handleTabChange,
                        labelColor: ColorName.white,
                        unselectedLabelColor: ColorName.grey80,
                        labelStyle: const TextStyle(
                          fontSize: 14.0,
                          fontWeight: FontWeight.w600,
                        ),
                        indicatorColor: ColorName.white,
                        labelPadding: const EdgeInsets.symmetric(
                          horizontal: 11,
                        ),
                        tabs: [
                          Tab(
                            text: LocaleKeys.balance.tr,
                          ),
                          Tab(text: LocaleKeys.transactionHistory.tr),
                          Tab(text: LocaleKeys.autoExchange.tr),
                        ],
                      ),
                    ),
                  ],
                ),
                if (controller.activeTabIndex.value == 0)
                  DepositWithdrawButtons(
                    isUserVerified: controller.isUserVerified.value,
                  ),
                if (controller.activeTabIndex.value == 0)
                  BalanceView()
                else if (controller.activeTabIndex.value == 1)
                  TransactionHistoryView()
                else if (controller.activeTabIndex.value == 2)
                  AutoExchangeView()
              ],
            ),
          ),
        ),
      ),
    );
  }
}
