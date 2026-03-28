import 'package:dio/dio.dart' show DioError;
import 'package:get/get.dart';
import 'package:get_storage/get_storage.dart';
import '../../account/controllers/account_controller.dart';
import '../../account/user_model.dart';
import '../providers/twoFactorAuthenticationProvider.dart';
import '../toggle2fa_model.dart';
import '../views/twoFactorIntroPopup.dart';
import '../../../../services/apiService.dart';
import '../../../../services/storageKeys.dart';
import '../../../../utils/mixins/toast.dart';

enum TwoFaSteps {
  Enable_Step1_Install_Links,
  Enable_step2_copyCode,
  Enable_Step3_EnterPassword_And_Code,
  Enable_Step4_Final_Status,
}

class TwoFactorAuthenticationController extends GetxController with Toaster {
  final AccountController accountController = Get.find();
  final storage = GetStorage();
  final twofaProvider = TwoFactorAuthenticationProviderProvider();
  final step = (TwoFaSteps.Enable_Step1_Install_Links).obs;
  final codeCoppied = false.obs;
  final isLoadingCharCode = false.obs;

  final isFinalSubmitLoading = false.obs;
  final characterCode = ''.obs;
  final qrImageAddress = ''.obs;
  final code = ''.obs;
  final password = ''.obs;
  final isEnabled = false.obs;

  @override
  void onInit() {
    isEnabled.value = accountController.accountData.value.google2faEnabled;
    if (isEnabled.value) {
      step.value = TwoFaSteps.Enable_Step3_EnterPassword_And_Code;
    }
    super.onInit();
  }

  @override
  void onReady() {
    super.onReady();
  }

  @override
  void onClose() {}

  void handleGoToCharacterCodeClick() async {
    try {
      isLoadingCharCode.value = true;
      final response = await twofaProvider.getCharacterCode();
      if (response['status'] == true) {
        characterCode.value = response['data']['code'];
        qrImageAddress.value = response['data']['url'];
        step.value = TwoFaSteps.Enable_step2_copyCode;
      }
    } on DioError catch (e) {
      toastDioError(e);
    } finally {
      isLoadingCharCode.value = false;
    }
  }

  void handleFinalSubmitClick({bool enable}) async {
    try {
      isFinalSubmitLoading.value = true;
      final model = Toggle2faModel(
        code: code.value,
        password: password.value,
      );
      final response =
          await twofaProvider.toggle2Fa(data: model, setEnabled: enable);
      if (response['status'] == true) {
        final jsoned = accountController.accountData.value.toJson();
        jsoned["google2faEnabled"] = enable;
        accountController.accountData.value = UserModel.fromJson(jsoned);
        isEnabled.value = enable;
        step.value = TwoFaSteps.Enable_Step4_Final_Status;
        if (response['data'] != null && response['data']['token'] != null) {
          ApiService.token = response['data']['token'];
          storage.write(StorageKeys.lastLoginDate, DateTime.now().toString());
          storage.write(StorageKeys.token, response['data']['token']);
        }
      }
    } on DioError catch (e) {
      toastDioError(e);
    } finally {
      isFinalSubmitLoading.value = false;
    }
  }

  void handlePasswordChange(String v) {
    password.value = v;
  }

  void handleCodeChange(String v) {
    code.value = v;
  }

  openIntro() {
    Get.dialog(
      TwoFactorIntroPopup(),
    );
  }
}
