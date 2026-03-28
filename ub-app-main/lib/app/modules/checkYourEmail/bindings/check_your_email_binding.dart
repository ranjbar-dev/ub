import 'package:get/get.dart';

import '../controllers/check_your_email_controller.dart';

class CheckYourEmailBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut<CheckYourEmailController>(
      () => CheckYourEmailController(),
    );
  }
}
