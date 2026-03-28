import 'package:get/get.dart';

import '../controllers/balance_controller.dart';

class BalanceBinding extends Bindings {
  @override
  void dependencies() {
    Get.put<BalanceController>(BalanceController(), permanent: true);
  }
}
