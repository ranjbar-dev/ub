import 'package:get/get.dart';
import '../pages/balance/controllers/balance_controller.dart';
import '../pages/deposits/controllers/deposits_controller.dart';
import '../pages/transactionHistory/controllers/transaction_history_controller.dart';
import '../pages/withdrawals/controllers/withdrawals_controller.dart';

class FundsController extends GetxController {
  TransactionHistoryController transactionHistoryController;
  final activeTabIndex = 0.obs;
  final isHeadOpen = false.obs;
  final isUserVerified = false.obs;
  @override
  void onInit() {
    Get.put(BalanceController(), permanent: true);
    Get.put(DepositsController(), permanent: true);
    Get.put<WithdrawalsController>(WithdrawalsController(), permanent: true);
    final BalanceController balanceController = Get.find();
    balanceController.isHeadOpen.listen((v) {
      isHeadOpen.value = v;
    });
    super.onInit();
  }

  @override
  void onReady() {
    super.onReady();
  }

  void handleTabChange(int index) {
    if (transactionHistoryController == null) {
      Get.put<TransactionHistoryController>(TransactionHistoryController(),
          permanent: true);
      transactionHistoryController = Get.find();
    }
    if (index == 1) {
      transactionHistoryController.getTransactionHistory();
    }

    activeTabIndex.value = index;
    return;
  }

  @override
  void onClose() {}
}
