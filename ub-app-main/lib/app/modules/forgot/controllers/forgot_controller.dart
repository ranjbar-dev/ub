import 'package:dio/dio.dart' show DioError;
import 'package:get/get.dart';
import '../forgot_password_model.dart';
import '../providers/forgotProvider.dart';
import '../../../../utils/cryptography/encoding.dart';
import '../../../../utils/emailValidator.dart';
import '../../../../utils/mixins/toast.dart';

class ForgotController extends GetxController with Toaster {
  final forgotProvider = ForgotProvider();
  final email = ''.obs;
  final emailError = ''.obs;
  final isLoading = false.obs;
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

  handleEmailChange(String e) {
    email.value = e;
  }

  void handleSubmitClick() async {
    email.value = email.value.replaceAll(' ', '');
    emailError.value = validateEmail(email.value);
    if (emailError.value == '') {
      try {
        final encoded = await genarateEnc(startTime: startTime);

        isLoading.value = true;
        final response = await forgotProvider.forgot(
          data: ForgotPasswordModel(
            email: email.value,
            recaptcha: encoded,
          ),
        );
        if (response['status'] == true) {
          toastSuccess('Please check your email inbox', duration: 10);
          // Get.offNamed(AppPages.LOGIN);
        }
      } on DioError catch (e) {
        toastDioError(e);
      } finally {
        isLoading.value = false;
      }
    }
  }

  void setStartTime() {
    startTime = DateTime.now();
  }
}
