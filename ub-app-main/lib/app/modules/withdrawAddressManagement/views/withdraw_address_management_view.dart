import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';

import '../../../../generated/locales.g.dart';
import '../../../common/components/CenterUBLoading.dart';
import '../../../common/components/UBButton.dart';
import '../../../common/components/UBScrollBar.dart';
import '../../../common/components/UBoops.dart';
import '../../../common/components/appbarTextTitle.dart';
import '../../../common/components/pageContainer.dart';
import '../../../routes/app_pages.dart';
import '../controllers/withdraw_address_management_controller.dart';
import '../widgets/withdrawAddressRow.dart';

class WithdrawAddressManagementView
    extends GetView<WithdrawAddressManagementController> {
  final bool selectAddress;
  final String code;

  WithdrawAddressManagementView({this.code, this.selectAddress});

  @override
  Widget build(BuildContext context) {
    final scrollController = ScrollController();
    RefreshController _refreshController =
        RefreshController(initialRefresh: false);
    return PageContainer(
      appbarTitle: AppBarTextTitle(
        title: LocaleKeys.withdrawAddressManagement.tr,
      ),
      child: Padding(
        padding: const EdgeInsets.only(
          bottom: 12,
        ),
        child: Stack(
          children: [
            Column(
              children: [
                Obx(
                  () {
                    final data = code == null
                        // ignore: invalid_use_of_protected_member
                        ? controller.withdrawAddresses.value
                        // ignore: invalid_use_of_protected_member
                        : controller.withdrawAddresses.value
                            .where((item) => item.code == code)
                            .toList();
                    return Expanded(
                        child: controller.loadingData.value == true
                            ? Container(
                                child: CenterUbLoading(),
                              )
                            : UBScrollBar(
                                pullToRefreshConfig: PullToRefreshConfig(
                                  controller: _refreshController,
                                  isLoading: controller.isSilentLoading.value,
                                  onRefreshLoading:
                                      controller.handlePullToRefresh,
                                  oopsVariant:
                                      OopsVariant.AddressManagementOops,
                                ),
                                itemCount: data.length,
                                scrollController: scrollController,
                                builder: (BuildContext context, int index) {
                                  final item = data[index];
                                  if (selectAddress == true) {
                                    return GestureDetector(
                                      onTap: () {
                                        Get.back(result: item.address);
                                      },
                                      child: WithdrawAddressRow(
                                        onDeleteClick: () =>
                                            controller.handleDeleteClick(index),
                                        item: item,
                                        onlySelectable: true,
                                      ),
                                    );
                                  }
                                  return WithdrawAddressRow(
                                    item: item,
                                    onDeleteClick: () =>
                                        controller.handleDeleteClick(index),
                                  );
                                },
                              ));
                  },
                ),
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 12.0),
                  child: Obx(
                    () {
                      if (controller.loadingData.value == false &&
                          selectAddress != true) {
                        return UBButton(
                          onClick: () {
                            Get.toNamed(AppPages.ADD_NEW_ADDRESS);
                          },
                          text: LocaleKeys.addNewAddress.tr,
                        );
                      }
                      return const SizedBox();
                    },
                  ),
                )
              ],
            ),
            Obx(() {
              final isRefreshing = controller.isRefreshing.value;
              return isRefreshing ? CenterUbLoading() : const SizedBox();
            })
          ],
        ),
      ),
    );
  }
}
