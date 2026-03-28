import 'package:get/get.dart';

import '../controllers/two_factor_authentication_controller.dart';

class TwoFactorAuthenticationBinding extends Bindings {
  @override
  void dependencies() {
    Get.put<TwoFactorAuthenticationController>(
      TwoFactorAuthenticationController(),
    );
  }
}
