import 'package:flutter/material.dart' show TextEditingController;
import 'package:get/get.dart';

class TwoFaController extends GetxController {
  final userEmail = 'your email'.obs;
  final userMobile = 'your mobile phone'.obs;

  int numberOfInputs = 0;
  final canSubmit = false.obs;
  final isLoading = false.obs;
  final isLoadingResendEmail = false.obs;
  final isLoadingResendPhoneCode = false.obs;

  final twoFaController = TextEditingController();
  final emailCodeController = TextEditingController();
  final phoneCodeController = TextEditingController();
  final twoFactorValue = ''.obs;
  final emailValue = ''.obs;
  final phoneValue = ''.obs;
  List<String> filledFields = [];
  int numberOfFilledFields = 0;

  @override
  void onInit() {
    super.onInit();
  }

  @override
  void onClose() {
    twoFaController.dispose();
    emailCodeController.dispose();
    phoneCodeController.dispose();
    super.onClose();
  }

  handleEmailChange(String v) {
    emailValue.value = v;
    _checkIfCanSubmit('email', emailValue.value);
  }

  handlePhoneCodeChange(String v) {
    phoneValue.value = v;
    _checkIfCanSubmit('phone', phoneValue.value);
  }

  handleTwoFactorChange(String v) {
    twoFactorValue.value = v;
    _checkIfCanSubmit('twoFactor', twoFactorValue.value);
  }

  void _checkIfCanSubmit(String key, String value) {
    if (value.length > 0) {
      if (filledFields.indexOf(key) == -1) {
        numberOfFilledFields++;
        filledFields.add(key);
      }
      if (numberOfFilledFields == numberOfInputs) {
        canSubmit.value = true;
      }
    } else {
      numberOfFilledFields--;
      filledFields.removeWhere((element) => element == key);
      canSubmit.value = false;
    }
  }

  void resetValues() {
    canSubmit.value = false;
    twoFaController.clear();
    emailCodeController.clear();
    phoneCodeController.clear();
    twoFactorValue.value = '';
    emailValue.value = '';
    phoneValue.value = '';
    filledFields = [];
    numberOfFilledFields = 0;
    isLoading.value = false;
  }

  void onEmaiResendClick() {}

  void onPhoneResendClick() {}
}
