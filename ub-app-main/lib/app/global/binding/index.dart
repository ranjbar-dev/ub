import 'package:get/get.dart';

import '../controller/globalController.dart';

class GlobalBinding implements Bindings {
  @override
  void dependencies() {
    Get.put<GlobalController>(GlobalController(), permanent: true);
  }
}
