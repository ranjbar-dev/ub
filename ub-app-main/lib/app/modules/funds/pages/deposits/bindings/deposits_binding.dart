import 'package:get/get.dart';

import '../controllers/deposits_controller.dart';

class DepositsBinding extends Bindings {
  @override
  void dependencies() {
    Get.put<DepositsController>(DepositsController(), permanent: true);
  }
}
