import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../common/components/UBButton.dart';
import '../../common/components/UBInputWithTitleAndPaste.dart';
import '../../common/components/UBText.dart';
import '../controllers/twofaPopupController.dart';
import '../../../generated/colors.gen.dart';
import '../../../generated/locales.g.dart';
import '../../../utils/mixins/commonConsts.dart';

const verificationHeadHeigth = 110.0;

class TwoFaPopupView extends GetView<TwoFaController> {
  final Function onSubmit;

  final bool isNewDevice;
  final bool need2fa;
  final bool needEmailCode;
  final bool withResendEmail;
  final bool withResendPhoneCode;
  TwoFaPopupView(
      {this.withResendEmail,
      this.withResendPhoneCode,
      this.isNewDevice,
      this.need2fa,
      this.needEmailCode,
      this.onSubmit});
  @override
  Widget build(BuildContext context) {
    controller.resetValues();
    int numberOfInputs = 0;
    if (need2fa) {
      numberOfInputs++;
    }
    if (needEmailCode) {
      numberOfInputs++;
    }
    if (isNewDevice) {
      numberOfInputs++;
    }
    controller.numberOfInputs = numberOfInputs;
    return Obx(() {
      final emailCodeController = controller.emailCodeController;
      final phoneCodeController = controller.phoneCodeController;
      final twoFaController = controller.twoFaController;
      final canSubmit = controller.canSubmit.value;
      final email = controller.userEmail.value;
      final phone = controller.userMobile.value;
      final isLoading = controller.isLoading.value;
      return WillPopScope(
        onWillPop: () async {
          controller.resetValues();
          return true;
        },
        child: Container(
          height: ((numberOfInputs * 85.0) + 120.0) + verificationHeadHeigth,
          decoration: const BoxDecoration(
              color: ColorName.black2c, borderRadius: roundedTop_big),
          child: Column(
            children: [
              SizedBox(
                height: verificationHeadHeigth,
                child: Center(
                  child: UBText(
                    text: 'Verification',
                    size: 18.0,
                    weight: FontWeight.w600,
                    color: ColorName.greyd8,
                  ),
                ),
              ),
              if (needEmailCode)
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 12),
                  child: Row(
                    children: [
                      UBInputWithTitleAndPaste(
                        width: withResendEmail == true
                            ? Get.width - 124.0
                            : Get.width - 24.0,
                        withPaste: true,
                        controller: emailCodeController,
                        onChange: controller.handleEmailChange,
                        placeHolder: 'Email verifocation code',
                        title: 'Code will be sent to $email.',
                      ),
                      if (withResendEmail == true) hspace8,
                      if (withResendEmail == true)
                        UBButton(
                            width: 90,
                            onClick: () {
                              controller.onEmaiResendClick();
                            },
                            text: 'Resend')
                    ],
                  ),
                ),
              if (isNewDevice)
                Padding(
                  padding: const EdgeInsets.symmetric(horizontal: 12),
                  child: Row(
                    children: [
                      UBInputWithTitleAndPaste(
                        width: withResendPhoneCode == true
                            ? Get.width - 124.0
                            : Get.width - 24.0,
                        controller: phoneCodeController,
                        withPaste: true,
                        onChange: controller.handlePhoneCodeChange,
                        placeHolder: 'SMS verifocation code',
                        title: 'Code will be sent to $phone.',
                      ),
                      if (withResendPhoneCode == true) hspace8,
                      if (withResendPhoneCode == true)
                        UBButton(
                            width: 90,
                            onClick: () {
                              controller.onPhoneResendClick();
                            },
                            text: 'Resend')
                    ],
                  ),
                ),
              if (need2fa)
                UBInputWithTitleAndPaste(
                  withPaste: true,
                  onChange: controller.handleTwoFactorChange,
                  controller: twoFaController,
                  placeHolder: 'Enter google verfication code.',
                  title: 'Google verification code',
                ),
              const Spacer(),
              Container(
                padding: const EdgeInsets.symmetric(
                  horizontal: 12,
                  vertical: 24,
                ),
                child: UBButton(
                  disabled: !canSubmit,
                  isLodaing: isLoading,
                  onClick: () {
                    controller.isLoading.value = true;
                    onSubmit({
                      's2fa': twoFaController.text.replaceAll('', ''),
                      'emailCode': emailCodeController.text.replaceAll('', ''),
                      'phoneCode': phoneCodeController.text.replaceAll('', '')
                    });
                  },
                  text: LocaleKeys.submit.tr,
                ),
              )
            ],
          ),
        ),
      );
    });
  }
}
