import 'package:get/get.dart';
import '../../../global/controller/globalController.dart';

class LandingController extends GetxController {
  GlobalController _globalController = GlobalController();
  final isLoaded = false.obs;
  final showHero = false.obs;
  @override
  void onInit() {
    _globalController.getVersion();
    super.onInit();
  }

  @override
  void onReady() {
    isLoaded.value = true;
    Future.delayed(200.milliseconds).then((v) {
      showHero.value = true;
    });
    super.onReady();
  }

  @override
  void onClose() {}
}
