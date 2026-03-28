import 'package:get/get.dart';

import '../controllers/withdrawals_controller.dart';

class WithdrawalsBinding extends Bindings {
  @override
  void dependencies() {
    Get.put<WithdrawalsController>(WithdrawalsController(), permanent: true);
  }
}
