import 'package:get/get.dart';

import '../controllers/twofaPopupController.dart';

class GlobalBinding implements Bindings {
  @override
  void dependencies() {
    Get.put<TwoFaController>(TwoFaController());
  }
}
