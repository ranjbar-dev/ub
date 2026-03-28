import 'dart:async';

import 'package:dio/dio.dart' show DioError;
import 'package:get/get.dart';
import '../../../global/autocompleteModel.dart';
import '../../account/controllers/account_controller.dart';
import '../../account/user_model.dart';
import '../phone_verification_post_model.dart';
import '../providers/phoneVerificationProvider.dart';
import '../../../routes/app_pages.dart';
import '../../../../utils/mixins/popups.dart';
import '../../../../utils/mixins/toast.dart';
import '../smsRequestModel.dart';

enum PhoneVerificationSteps {
  EnterPhoneNumber,
  EnterSMSVerificationCode,
  EnterPassword
}

class PhoneVerificationController extends GetxController with Toaster, Popups {
  final AccountController accountController = Get.find();
  final is2faActivated = false.obs;
  final s2faCode = ''.obs;
  final initialCountdownValue = '01:00';
  final countDownValue = "1:00".obs;
  final canResend = false.obs;
  bool twofaEnabled = false;
  Timer timerController;
  final interval = const Duration(milliseconds: 1050);
  final int timerMaxSeconds = 60; // 60 * 5;
  int currentSeconds = 0;
  String get timerText =>
      '${((timerMaxSeconds - currentSeconds) ~/ 60).toString().padLeft(2, '0')}: ${((timerMaxSeconds - currentSeconds) % 60).toString().padLeft(2, '0')}';
  startTimeout([int milliseconds]) {
    var duration = interval;
    Timer.periodic(duration, (timer) {
      timerController = timer;
      currentSeconds = timerController.tick;
      countDownValue.value = timerText;
      if (timerController.tick >= timerMaxSeconds) {
        canResend.value = true;
        timerController.cancel();
      }
    });
  }

  final phoneVerificationProvider = PhoneVerificationProvider();
  final step = (PhoneVerificationSteps.EnterPhoneNumber).obs;

  final selectedCountry = AutoCompleteItem(name: '', code: "0").obs;
  final isRequestingForSms = false.obs;
  final isRequestingtoSubmitCode = false.obs;
  final phoneNumber = ''.obs;
  final password = ''.obs;
  final verificationCode = ''.obs;

  @override
  void onInit() {
    is2faActivated.value = accountController.accountData.value.google2faEnabled;
    super.onInit();
  }

  @override
  void onReady() {
    super.onReady();
  }

  @override
  void onClose() {}

  void handleCountrySelected(AutoCompleteItem item) {
    selectedCountry.value = item;
  }

  void handlePhoneNumberChange(String v) {
    if (canResend.value == true) {
      canResend.value = false;
    }
    phoneNumber.value = v;
  }

  void phoneNumberSubmitted() async {
    final phone = "+" + selectedCountry.value.code + phoneNumber.value;
    try {
      final model = SMSRequestModel.fromJson({"phone": phone});
      isRequestingForSms.value = true;
      final response = await phoneVerificationProvider.requestSMS(data: model);
      if (response["status"] == true) {
        step.value = PhoneVerificationSteps.EnterSMSVerificationCode;
        _resetTimer();
      }
    } on DioError catch (e) {
      toastDioError(e);
    } finally {
      isRequestingForSms.value = false;
    }
  }

  void handleEditPhoneClicked() {
    phoneNumber.value = '';
    step.value = PhoneVerificationSteps.EnterPhoneNumber;
  }

  void _resetTimer() {
    if (timerController != null) {
      timerController.cancel();
    }
    countDownValue.value = initialCountdownValue;
    startTimeout();
    canResend.value = false;
  }

  void handleVerificationCodeChange(String v) {
    verificationCode.value = v;
    if (v.length > 5) {
      if (twofaEnabled == true) {
        finalSubmitClicked();
        return;
      }
      step.value = PhoneVerificationSteps.EnterPassword;
    }
  }

  void handleGoTo2faPageClick() {}

  void handlePasswordChange(String v) {
    password.value = v;
  }

  void finalSubmitClicked({PhoneVerificationPostModel model}) async {
    final phone = "+" + selectedCountry.value.code + phoneNumber.value;
    final requestModel = model ??
        PhoneVerificationPostModel.fromJson({
          "phone": phone,
          "code": verificationCode.value,
          if (password.value != '') "password": password.value,
        });
    try {
      isRequestingtoSubmitCode.value = true;
      final response = await phoneVerificationProvider
          .submitSMSVerificationCode(data: requestModel);
      if (response['status'] == true) {
        if (response['data'] == null) {
          response['data'] = {};
        }
        if (response['data']["need2fa"] == true ||
            response['data']["needEmailCode"] == true ||
            response['data']['isNewDevice'] == true) {
          openVerificationPopup(
              need2fa: response['data']['need2fa'] ?? false,
              needEmailCode: response['data']['needEmailCode'] ?? false,
              isNewDevice: false,
              onSubmit: (v) async {
                requestModel.s2faCode = v['s2fa'];
                //requestModel.emailCode = v['emailCode'];
                Get.back();
                finalSubmitClicked(model: requestModel);
              });
        } else {
          toastSuccess('Phone number verified');
          final userData = accountController.accountData.value;
          final jsoned = userData.toJson();
          jsoned["phone"] = phone;
          accountController.accountData.value = UserModel.fromJson(jsoned);
          Get.offAllNamed(AppPages.ACCOUNT);
        }
      } else {
        print(response);
      }
    } on DioError catch (e) {
      print(e.toString());
      toastDioError(e);
      reset();
    } finally {
      isRequestingtoSubmitCode.value = false;
    }
  }

  void reset({Map<String, dynamic> recivingData}) {
    if (recivingData != null && recivingData["twofaEnabled"] != null) {
      twofaEnabled = recivingData["twofaEnabled"];
    }
    Future.delayed(100.milliseconds).then((value) {
      step.value = PhoneVerificationSteps.EnterPhoneNumber;
      selectedCountry.value = AutoCompleteItem(name: '', code: "0");
      isRequestingForSms.value = false;
      isRequestingtoSubmitCode.value = false;
      phoneNumber.value = '';
      password.value = '';
      verificationCode.value = '';
    });
  }
}
