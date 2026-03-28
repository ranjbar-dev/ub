import 'package:get/get.dart';

import '../controllers/qr_scan_controller.dart';

class QrScanBinding extends Bindings {
  @override
  void dependencies() {
    Get.put<QrScanController>(
      QrScanController(),
    );
  }
}
