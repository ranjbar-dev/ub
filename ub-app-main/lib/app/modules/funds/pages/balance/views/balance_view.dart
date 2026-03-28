import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart'
    hide RefreshIndicator, RefreshIndicatorState;
import 'package:flutter_switch/flutter_switch.dart';
import 'package:get/get.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
import 'package:unitedbit/app/common/components/CenterUBLoading.dart';
import 'package:unitedbit/app/common/components/UBScrollBar.dart';
import 'package:unitedbit/app/common/components/UBText.dart';
import 'package:unitedbit/app/common/components/UBoops.dart';
import 'package:unitedbit/app/modules/funds/pages/balance/widgets/balanceRow.dart';
import 'package:unitedbit/generated/colors.gen.dart';
import 'package:unitedbit/generated/locales.g.dart';

import '../controllers/balance_controller.dart';

class BalanceView extends GetView<BalanceController> {
  final scrollController = ScrollController();
  final RefreshController _refreshController =
      RefreshController(initialRefresh: false);

  @override
  Widget build(BuildContext context) {
    return Expanded(
      child: Column(
        children: [
          Container(
            height: 40,
            color: ColorName.grey16,
            padding: const EdgeInsets.symmetric(horizontal: 12),
            child: Row(
              children: [
                UBText(
                  text: LocaleKeys.smallBalances.tr,
                  size: 12,
                  color: ColorName.greybf,
                ),
                const SizedBox(
                  width: 12,
                ),
                Container(
                  width: 35,
                  padding: const EdgeInsets.symmetric(vertical: 10),
                  child: Obx(() => FlutterSwitch(
                        padding: 3,
                        toggleSize: 15,
                        activeColor: ColorName.primaryBlue,
                        value: controller.showSmallBalances.value,
                        onToggle: (val) {
                          controller.handleShowSmallBalancesChange(val);
                        },
                      )),
                ),
                const SizedBox(
                  width: 12,
                ),
                Obx(() => UBText(
                      text: controller.showSmallBalances.value == true
                          ? 'Show'
                          : 'Hide',
                      size: 12,
                      color: ColorName.greybf,
                    )),
              ],
            ),
          ),
          Expanded(
            child: Obx(
              () {
                final balances = controller.balancesAllData.value.balances;
                final showSmallBalances = controller.showSmallBalances.value;

                if (balances == null && controller.isLoading.value == false) {
                  return UBoops(
                    variant: OopsVariant.ErrorOops,
                  );
                }
                return controller.isLoading.value == true
                    ? CenterUbLoading()
                    : UBScrollBar(
                        pullToRefreshConfig: PullToRefreshConfig(
                          oopsVariant: OopsVariant.NoBalancesOops,
                          controller: _refreshController,
                          isLoading: controller.isSilentLoading.value,
                          onRefreshLoading: controller.handleRefreshBalances,
                          withUpdateDate: true,
                        ),
                        scrollController: scrollController,
                        itemCount: balances.length,
                        builder: (BuildContext context, int index) {
                          final balance = balances[index];
                          final isGreaterThanSmallValue =
                              double.parse(balance.totalAmount) >
                                  double.parse(controller.balancesAllData.value
                                      .minimumOfSmallBalances);
                          if (showSmallBalances == true ||
                              (showSmallBalances == false &&
                                  isGreaterThanSmallValue)) {
                            if (index == balances.length - 1)
                              return Column(
                                children: [
                                  BalanceRow(
                                    balance: balances[index],
                                  ),
                                  Container(
                                    height: 49,
                                  )
                                ],
                              );
                            else
                              return BalanceRow(
                                balance: balances[index],
                              );
                          }
                          return const SizedBox();
                        },
                      );
              },
            ),
          )
        ],
      ),
    );
  }
}
