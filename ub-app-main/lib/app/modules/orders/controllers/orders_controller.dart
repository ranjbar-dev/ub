import 'package:get/get.dart';
import '../pages/openOrders/controllers/open_orders_controller.dart';
import '../pages/orderHistory/controllers/order_history_controller.dart';
import '../views/orders_view.dart';
import '../../../routes/app_pages.dart';

class OrdersController extends GetxController {
  OpenOrdersController openOrdersController;

  final isFullScreen = false.obs;
  final activeTabIndex = 0.obs;
  final expanded = false.obs;
  void handleTabChange(int index) {
    activeTabIndex.value = index;
  }

  @override
  void onInit() {
    Get.put(OpenOrdersController(), permanent: true);
    Get.put(OrderHistoryController(), permanent: true);
    openOrdersController = Get.find();
    super.onInit();
  }

  @override
  void onReady() {
    super.onReady();
  }

  @override
  void onClose() {}

  void handleExpandClick() {
    expanded.toggle();
    if (!expanded.value) {
      Get.toNamed(AppPages.ORDERS);
      return;
    }
    Get.back();
  }

  void handleExpandOrdersClick() async {
    isFullScreen.toggle();
    openOrdersController.isFullScreen.value = isFullScreen.value;
    if (isFullScreen.value) {
      Get.to(() => OrdersView(
            fullScreen: true,
          ));
      return;
    }
    Get.back();
  }

  void handleWillPop() {
    if (isFullScreen.value == true) {
      isFullScreen.toggle();
    }
  }
}
