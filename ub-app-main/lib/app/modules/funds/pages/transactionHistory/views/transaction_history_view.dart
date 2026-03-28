import 'package:flutter/material.dart';

import 'package:get/get.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';
import '../../../../../common/components/CenterUBLoading.dart';
import '../../../../../common/components/UBDarkOpacityBackgrounded.dart';
import '../../../../../common/components/UBScrollBar.dart';
import '../../../../../common/components/UBoops.dart';
import '../widgets/transactionHistoryRow.dart';
import '../../../../../../utils/mixins/commonConsts.dart';
import '../controllers/transaction_history_controller.dart';

class TransactionHistoryView extends GetView<TransactionHistoryController> {
  final ScrollController _scrollController = ScrollController();
  final RefreshController _refreshController = RefreshController();
  @override
  Widget build(BuildContext context) {
    return Obx(
      () {
        final data = controller.transactionHistory.value.payments ?? [];
        return Expanded(
          child: controller.isLoading.value == true
              ? Container(
                  child: CenterUbLoading(),
                )
              : Stack(
                  children: [
                    UBScrollBar(
                      pullToRefreshConfig: PullToRefreshConfig(
                        controller: _refreshController,
                        onRefreshLoading: controller.handlePullToRefresh,
                        isLoading: controller.isSilentLoading.value,
                        oopsVariant: OopsVariant.NoTransactionHistory,
                        withUpdateDate: true,
                      ),
                      itemCount: data.length,
                      builder: (BuildContext context, int index) =>
                          TransactionHistoryRow(
                        data: data[index],
                        onTap: controller.handleRowClick,
                      ),
                      scrollController: _scrollController,
                    ),
                    Obx(() {
                      final isLoading = controller.showLoadingOverlay.value;
                      if (isLoading) {
                        return UBDarkOpacityBackgrounded(
                          child: CenterUbLoading(),
                        );
                      }
                      return emptyComponent;
                    })
                  ],
                ),
        );
      },
    );
  }
}
