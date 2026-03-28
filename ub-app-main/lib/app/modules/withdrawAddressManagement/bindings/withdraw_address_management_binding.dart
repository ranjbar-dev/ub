import 'package:get/get.dart';

import '../controllers/withdraw_address_management_controller.dart';

class WithdrawAddressManagementBinding extends Bindings {
  @override
  void dependencies() {
    Get.put<WithdrawAddressManagementController>(
        WithdrawAddressManagementController(),
        permanent: true);
  }
}
