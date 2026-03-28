import 'package:get/get.dart';

import '../controllers/open_orders_controller.dart';

class OpenOrdersBinding extends Bindings {
  @override
  void dependencies() {
    Get.put<OpenOrdersController>(OpenOrdersController(), permanent: true);
  }
}
