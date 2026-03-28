import 'package:get/get.dart';

import '../controllers/phone_verification_controller.dart';

class PhoneVerificationBinding extends Bindings {
  @override
  void dependencies() {
    Get.put<PhoneVerificationController>(
      PhoneVerificationController(),
    );
  }
}
