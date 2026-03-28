import 'package:get/get.dart';

import '../controllers/identity_info_controller.dart';

class IdentityInfoBinding extends Bindings {
  @override
  void dependencies() {
    Get.put<IdentityInfoController>(
      IdentityInfoController(),
    );
  }
}
