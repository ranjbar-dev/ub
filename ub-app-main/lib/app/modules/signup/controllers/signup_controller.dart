import 'package:dio/dio.dart' show DioError;
import 'package:get/get.dart';
import '../providers/signupProvider.dart';
import '../signup_model.dart';
import '../../../../utils/cryptography/encoding.dart';
import '../../../../utils/emailValidator.dart';
import '../../../../utils/mixins/toast.dart';
import '../../../../utils/passwordValidator.dart';

enum signupStep { SigningUp, SignedUp }

class SignupController extends GetxController with Toaster {
  final signupProvider = SignupProvider();

  final isLoading = false.obs;
  final email = ''.obs;
  final password = ''.obs;
  final repeatPassword = ''.obs;
  final emailError = ''.obs;
  final passwordError = ''.obs;
  final repeatPasswordError = ''.obs;
  final step = signupStep.SigningUp.obs;

  DateTime startTime;

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

  handleEmailChange(String v) {
    email.value = v;
  }

  handlePasswordChange(String v) {
    password.value = v;
  }

  handleRepeatPasswordChange(String v) {
    repeatPassword.value = v;
  }

  handleSubmitClick() async {
    final canSubmit = _canSubmit();
    if (canSubmit) {
      final encoded = await genarateEnc(startTime: startTime);

      try {
        isLoading.value = true;
        final data = SignupModel(
          email: email.value,
          password: password.value,
          recaptcha: encoded,
        );
        final response = await signupProvider.signup(data: data);
        if (response['status'] == true) {
          step.value = signupStep.SignedUp;
          toastSuccess(
              'Please verify your email and login to continue identity verification',
              duration: 7000);
        }
      } on DioError catch (e) {
        toastDioError(e);
      } finally {
        isLoading.value = false;
      }
    }
  }

  _canSubmit() {
    email.value = email.value.replaceAll(' ', '');
    emailError.value = validateEmail(email.value);
    passwordError.value = validatePassword(password.value);
    repeatPasswordError.value = validatePassword(repeatPassword.value);
    if (password.value != '' &&
        repeatPassword.value != '' &&
        repeatPassword.value != password.value) {
      passwordError.value = repeatPasswordError.value =
          'password and repeat password are noe equal';
    }
    return emailError.value == '' &&
        passwordError.value == '' &&
        repeatPasswordError.value == '';
  }

  void setStartTime() {
    startTime = DateTime.now();
  }
}
