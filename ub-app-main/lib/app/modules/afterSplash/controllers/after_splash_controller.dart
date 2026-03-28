import 'dart:async';

import 'package:get/get.dart';
import 'package:get_storage/get_storage.dart';
import '../../../global/controller/globalController.dart';
import '../../../global/providers/commonDataProvider.dart';
import '../providers/afterSplashProvider.dart';
import '../../../routes/app_pages.dart';
import '../../../../services/storageKeys.dart';
import '../../../../utils/logger.dart';
import '../../../../utils/mixins/popups.dart';

final afterSplashAnimationDuration = 2.seconds;

class AfterSplashController extends GetxController with Popups {
  final bool test;
  AfterSplashController({this.test});

  CommonDataProvider commonDataProvider = CommonDataProvider();

  final GlobalController globalController = Get.find();
  final AfterSplashProvider afterSplashProvider = AfterSplashProvider();
  final storage = GetStorage();
  StreamSubscription<bool> connectionSubscription;
  StreamSubscription<dynamic> currencyArraySubscription;

  final animationStarted = false.obs;
  final isConnected = true.obs;
  final showRetryButton = false.obs;
  final isLoadingRetry = false.obs;
  bool navigated = false;

  @override
  void onInit() async {
    if (this.test != true) {
      connectionSubscription =
          globalController.hasConnection.listen((connected) {
        if (!connected) {
          isConnected.value = false;
        } else {
          isConnected.value = true;
          Future.delayed(15.seconds).then((value) {
            if (isConnected.value) {
              showRetryButton.value = true;
            }
          });
        }
      });
    }
    super.onInit();
  }

  @override
  void onReady() {
    if (this.test != true) {
      animationStarted.value = true;
      _checkInitialDataAndNavigate();
    }
    super.onReady();
  }

  @override
  void onClose() {
    if (this.test != true) {
      connectionSubscription.cancel();
      if (currencyArraySubscription != null) {
        currencyArraySubscription.cancel();
      }
    }
  }

  void _checkInitialDataAndNavigate() {
    Future.delayed(afterSplashAnimationDuration).then((v) {
      // ignore: invalid_use_of_protected_member
      if (globalController.currencyPairsArray.value.isNotEmpty &&
          globalController.isLoggingInWithBiometrics == false) {
        _navigate();
      } else {
        currencyArraySubscription =
            globalController.currencyPairsArray.listen((v) {
          if (v.isNotEmpty &&
              globalController.isLoggingInWithBiometrics == false) {
            _navigate();
          }
        });
      }
    });
  }

  _navigate() {
    navigated = true;
    if (globalController.loggedIn.value == true) {
      Get.offNamed(AppPages.HOME);
    } else {
      if (storage.read(StorageKeys.loggedInOnce) == true) {
        // commented because mehdi
        //Get.offNamed(AppPages.LOGIN);

        Get.offNamed(AppPages.LANDING);
        return;
      }
      Get.offNamed(AppPages.LANDING);
    }
  }

  void handleRetryButton() {
    isLoadingRetry.value = true;
    try {
      globalController.onInit();
    } catch (e) {
      log.e(e.toString());
      isLoadingRetry.value = false;
    }
  }
}
