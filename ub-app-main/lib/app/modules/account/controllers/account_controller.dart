import 'package:flutter/foundation.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:get/get.dart';
import 'package:get_storage/get_storage.dart';
import '../../../common/components/UBText.dart';
import '../../../global/controller/globalController.dart';
import '../../../global/providers/commonDataProvider.dart';
import '../providers/accountProvider.dart';
import '../user_model.dart';
import '../../../../services/localAuthService.dart';
import '../../../../services/storageKeys.dart';
import '../../../../utils/mixins/popups.dart';
import '../../../../utils/mixins/toast.dart';
import '../../../../utils/throttle.dart';

class AccountController extends GetxController with Toaster, Popups {
  final thr = new Throttling(duration: 4.seconds);
  final _secureStorage = new FlutterSecureStorage();
  final storage = GetStorage();
  final commonDataProvider = CommonDataProvider();
  final accountProvider = AccountProvider();
  final accountData = UserModel().obs;
  final requestedForEmail = 0.obs;
  final isRequestingForEmail = false.obs;
  final hasBiometrics = false.obs;
  final isBiometricsActivated = false.obs;
  @override
  void onInit() async {
    getUserData();
    final isActive = storage.read(StorageKeys.biometricsActivated);
    isBiometricsActivated.value = isActive ?? false;
    if (!(GetPlatform.isWeb)) {
      hasBiometrics.value = await BiometricsService.hasBiometrics();
    }

    super.onInit();
  }

  @override
  void onReady() {
    super.onReady();
  }

  @override
  void onClose() {}

  void toggleBiometrics() async {
    try {
      final authed = await BiometricsService().authenticateWithBiometrics();
      if (authed) {
        isBiometricsActivated.toggle();
        storage.write(
            StorageKeys.biometricsActivated, isBiometricsActivated.value);
        if (isBiometricsActivated.value) {
          toastSuccess('Biometrics enabled');
          return;
        }
        toastWarning('Biometrics disabled');
      }
    } catch (e) {
      debugPrint(e);
    }
  }

  void getUserData({bool force}) async {
    if ((force == true || accountData.value.email == null) ||
        storage.read(StorageKeys.token) != null) {
      final response = await commonDataProvider.getUserData();

      if (response['status'] == true) {
        accountData.value = UserModel.fromJson(
          response["data"],
        );
      }
    }
  }

  requestForEmailVerification() async {
    thr.throttle(() async {
      requestedForEmail.value++;
      toastInfo(
          "We've sent you a new verification email, please check your email",
          duration: 6000);
      if (requestedForEmail.value < 3) {
        isRequestingForEmail.value = true;
        try {
          await accountProvider.requestForEmail();
        } catch (e) {
          toastError('error while rquesting, please try again later');
        } finally {
          isRequestingForEmail.value = false;
        }
      }
    });
  }

  void handleLogoutClick() {
    openConfirmation(
      onConfirm: _logOutAction,
      titleWidget: UBText(text: "Logout from Unitedbit?"),
      titleDistanceFromTop: -4.0,
      cancelText: "Cancel",
      confirmText: 'Logout',
      autoBackAfterConfirm: false,
    );
  }

  _logOutAction() async {
    final GlobalController globalController = Get.find();
    if (!(GetPlatform.isWeb)) {
      await _secureStorage.deleteAll();
    }
    globalController.handleLoggedOut(andExitApp: false);
  }

  void pageLoaded() async {
    if (!(GetPlatform.isWeb)) {
      hasBiometrics.value = await BiometricsService.hasBiometrics();
    }
    getUserData();
  }
}
