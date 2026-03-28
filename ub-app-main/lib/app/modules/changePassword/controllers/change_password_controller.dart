import 'package:dio/dio.dart' show DioError;
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:get/get.dart';
import 'package:get_storage/get_storage.dart';
import '../change_password_model.dart';
import '../providers/changePasswordProvider.dart';
import '../../../../services/apiService.dart';
import '../../../../services/storageKeys.dart';
import '../../../../utils/mixins/popups.dart';
import '../../../../utils/passwordValidator.dart';
import '../../../../utils/mixins/toast.dart';

enum changePasswordStep { ChangePassword, PasswordChanged }

class ChangePasswordController extends GetxController with Toaster, Popups {
  String equalPassError = 'New password is not equal to repeated one';
  final GetStorage storage = GetStorage();
  final ApiService apiService = ApiService();

  final _secureStorage = FlutterSecureStorage();
  final ChangePasswordProvider changePasswordProvider =
      ChangePasswordProvider();
  final oldPasswordValue = ''.obs;
  final oldPasswordError = ''.obs;
  final newPasswordValue = ''.obs;
  final newPasswordError = ''.obs;
  final repeatNewPasswordValue = ''.obs;
  final repeatNewPasswordError = ''.obs;
  final isSubmitting = false.obs;
  final step = changePasswordStep.ChangePassword.obs;

  @override
  void onInit() {
    super.onInit();
  }

  @override
  void onReady() {
    super.onReady();
  }

  @override
  void onClose() {}

  handleOldPasswordValueChange(String v) {
    oldPasswordError.value = validatePassword(v);
    if (oldPasswordError.value == '') {
      oldPasswordValue.value = v;
    } else {
      oldPasswordValue.value = '';
    }
  }

  handleNewPasswordValueChange(String v) {
    newPasswordError.value = validatePassword(v);
    if (newPasswordError.value == '') {
      newPasswordValue.value = v;

      if (repeatNewPasswordValue.value != '') {
        if (newPasswordValue.value == repeatNewPasswordValue.value) {
          newPasswordError.value = '';
          repeatNewPasswordError.value = '';
        } else {
          newPasswordError.value = equalPassError;
        }
      }
    } else {
      newPasswordValue.value = '';
    }
  }

  handleRepeatNewPasswordValueChange(String v) {
    repeatNewPasswordError.value = validatePassword(v);

    if (repeatNewPasswordError.value == '') {
      repeatNewPasswordValue.value = v;
      if (newPasswordValue.value != '') {
        if (newPasswordValue.value == repeatNewPasswordValue.value) {
          newPasswordError.value = '';
          repeatNewPasswordError.value = '';
        } else {
          repeatNewPasswordError.value = equalPassError;
        }
      }
    } else {
      repeatNewPasswordValue.value = '';
    }
  }

  bool canSubmit() {
    return oldPasswordValue.value != '' &&
        newPasswordValue.value != '' &&
        repeatNewPasswordValue.value != '' &&
        newPasswordError.value == '' &&
        repeatNewPasswordError.value == '';
  }

  void handleSubmitClick({ChangePasswordModel data}) async {
    try {
      isSubmitting.value = true;
      final model = data ??
          ChangePasswordModel(
            newPassword: newPasswordValue.value,
            oldPassword: oldPasswordValue.value,
            confirmed: repeatNewPasswordValue.value,
          );
      final response = await changePasswordProvider.changePassword(data: model);
      if (response['status'] == true) {
        if (response['data'] == null) {
          response["data"] = {};
        }
        if (response['data']["need2fa"] == true ||
            response['data']["needEmailCode"] == true ||
            response['data']['isNewDevice'] == true) {
          openVerificationPopup(
              need2fa: response['data']['need2fa'] ?? false,
              needEmailCode: response['data']['needEmailCode'] ?? false,
              isNewDevice: false,
              onSubmit: (v) async {
                model.s2faCode = v['s2fa'];
                //model.emailCode = v['emailCode'];
                Get.back();
                handleSubmitClick(data: model);
              });
        } else {
          if (response['data']['token'] != null) {
            ApiService.token = response['data']['token'];
            storage.write(StorageKeys.lastLoginDate, DateTime.now().toString());
            storage.write(StorageKeys.token, response['data']['token']);
          }
          step.value = changePasswordStep.PasswordChanged;
          await _secureStorage.write(
            key: SecureStorageKeys.password,
            value: model.newPassword,
          );
        }
      }
    } on DioError catch (e) {
      toastDioError(e);
    } finally {
      isSubmitting.value = false;
    }
  }
}
