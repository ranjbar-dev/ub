import 'package:get/get.dart';

import '../../../centrifugoClient/centrifugo_service.dart';
import '../../../services/constants.dart';

class UnAuthorizedCentrifugoController extends GetxController {
  CentrifugoService centrifugoService;

  @override
  void onInit() {
    centrifugoService = CentrifugoService();
    centrifugoService.init(
      url: Constants.centrifugoWsUrl,
    );

    super.onInit();
  }

  @override
  void onClose() {
    if (centrifugoService != null) {
      centrifugoService.disconnect();
    }
    super.onClose();
  }
}
