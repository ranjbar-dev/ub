import 'package:get/get.dart';

import '../controllers/web_view_page_controller.dart';

class WebViewPageBinding extends Bindings {
  @override
  void dependencies() {
    Get.lazyPut<WebViewPageController>(
      () => WebViewPageController(),
    );
  }
}
