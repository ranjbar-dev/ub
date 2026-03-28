import 'package:get/get.dart';

import '../controllers/after_splash_controller.dart';

class AfterSplashBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut<AfterSplashController>(
      () => AfterSplashController(),
    );
  }
}
