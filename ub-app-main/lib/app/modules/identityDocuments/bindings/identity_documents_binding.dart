import 'package:get/get.dart';

import '../controllers/identity_documents_controller.dart';

class IdentityDocumentsBinding extends Bindings {
  @override
  void dependencies() {
    Get.put<IdentityDocumentsController>(
      IdentityDocumentsController(),
    );
  }
}
