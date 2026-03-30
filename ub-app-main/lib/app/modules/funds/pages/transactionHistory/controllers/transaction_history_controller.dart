import 'dart:async';

import 'package:get/get.dart';
import '../../../../../common/components/UBText.dart';
import '../../../../../global/controller/authorizedCentrifugoController.dart';
import '../providers/transactionHistoryProvider.dart';
import '../transaction_history_model.dart';
import '../../../../../../services/constants.dart';
import '../../../../../../utils/mixins/popups.dart';

class TransactionHistoryController extends GetxController with Popups {
  final transactionHistoryProvider = TransactionHistoryProvider();
  final AuthorizedCentrifugoController authorizedCentrifugoController = Get.find();

  final transactionHistory = TransactionHistoryModel().obs;
  final isLoading = false.obs;
  final showLoadingOverlay = false.obs;
  final isSilentLoading = false.obs;
  LightSubscription<List<RxUpdateables>> updateSubscription;

  @override
  void onInit() {
    updateSubscription =
        authorizedCentrifugoController.updateDataSubject.listen((value) {
      if (value is List &&
          (value.indexOf(RxUpdateables.TransactionHistory) != -1)) {
        getTransactionHistory(silent: true);
      }
      return;
    });
    super.onInit();
  }

  @override
  void onReady() {
    super.onReady();
  }

  @override
  void onClose() {
    if (updateSubscription != null) {
      updateSubscription.cancel();
    }
  }

  Future<void> getTransactionHistory({bool silent}) async {
    if (transactionHistory.value.payments == null && silent != true) {
      isLoading.value = true;
    }
    if (silent == true) {
      isSilentLoading.value = true;
    }
    try {
      final response = await transactionHistoryProvider.getTransactionHistory();
      if (response["status"] == true) {
        transactionHistory.value =
            TransactionHistoryModel.fromJson(response["data"]);
      }
    } catch (e) {
    } finally {
      isSilentLoading.value = false;
      if (isLoading.value == true) {
        isLoading.value = false;
      }
    }
    return Future.value();
  }

  handlePullToRefresh() {
    getTransactionHistory(silent: true);
    Future.delayed(2000.milliseconds).then((v) {
      isSilentLoading.value = false;
    });
  }

  handleRowClick(Payments data) {
    openTransactionDetailsPopup(data: data, onCancelClick: handleCancelClick);
  }

  handleCancelClick(int id) {
    openConfirmation(
      onConfirm: () {
        requestCancelWithdraw(id: id);
      },
      titleWidget: UBText(text: 'Cancel this withdraw?'),
      cancelText: 'No',
      confirmText: "Yes",
    );
  }

  void requestCancelWithdraw({int id}) async {
    //pop popups
    Get.back();
    Get.back();
    //
    try {
      showLoadingOverlay.value = true;
      final response = await transactionHistoryProvider.cancelWithdraw(id: id);
      if (response["status"] == true) {
        await getTransactionHistory(silent: true);
      }
    } catch (e) {
    } finally {
      isSilentLoading.value = false;
      showLoadingOverlay.value = false;
    }
  }
}
