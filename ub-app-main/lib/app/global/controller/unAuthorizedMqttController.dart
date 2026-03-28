import 'package:get/get.dart';
import 'package:uuid/uuid.dart';

import '../../../mqttClient/universal_mqtt_client.dart';
import '../../../services/constants.dart';

class UnAuthorizedMqttController extends GetxController {
  UniversalMqttClient unAuthorizedClient;

  @override
  void onInit() {
    unAuthorizedClient = new UniversalMqttClient(
      broker: Uri.parse(Constants.mqttServer),
      autoReconnect: true,
      timeout: const Duration(seconds: 10),
      username: Uuid().v4(),
      password: Uuid().v4(),
    );

    super.onInit();
  }

  @override
  void onClose() {
    if (unAuthorizedClient != null) {
      unAuthorizedClient.disconnect();
    }
    super.onClose();
  }
}
