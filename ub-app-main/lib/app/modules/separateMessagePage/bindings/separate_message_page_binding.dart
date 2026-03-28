import 'package:get/get.dart';

import '../controllers/separate_message_page_controller.dart';

class SeparateMessagePageBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut<SeparateMessagePageController>(
      () => SeparateMessagePageController(),
    );
  }
}
