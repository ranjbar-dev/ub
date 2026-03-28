import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../../generated/assets.gen.dart';
import '../../../../generated/colors.gen.dart';
import '../../../../generated/locales.g.dart';
import '../../../../utils/mixins/commonConsts.dart';
import '../../../common/components/UBBorderlessInput.dart';
import '../../../common/components/UBButton.dart';
import '../../../common/components/UBGreyContainer.dart';
import '../../../common/components/UBText.dart';
import '../../../routes/app_pages.dart';
import '../controllers/phone_verification_controller.dart';

class EnterSMSVerificationCode extends GetView<PhoneVerificationController> {
  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        fill,
        Container(
          height: 385.0,
          decoration: const BoxDecoration(
              color: ColorName.black2c, borderRadius: roundedTop_big),
          child: Column(
            children: [
              vspace24,
              Assets.images.smsLogo.svg(),
              vspace24,
              Container(
                width: 240,
                child: UBText(
                  lineHeight: 1.6,
                  size: 13,
                  color: ColorName.grey80,
                  align: TextAlign.center,
                  text: LocaleKeys.smsVerificationPageTopTitle.tr,
                ),
              ),
              Container(
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    UBText(
                      text: LocaleKeys.wesenttheverificationCodeTo.tr,
                      size: 13,
                      color: ColorName.grey80,
                    ),
                    _space(),
                    UBText(
                      text:
                          '+${controller.selectedCountry.value.code}${controller.phoneNumber.value}',
                      color: ColorName.primaryBlue,
                      size: 13,
                    ),
                  ],
                ),
              ),
              const SizedBox(
                height: 6,
              ),
              Container(
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    UBText(
                      fontStyle: FontStyle.italic,
                      text: LocaleKeys.didntreceivecode.tr,
                      color: ColorName.grey80,
                      size: 13,
                    ),
                    _space(),
                    UBButton(
                      height: 16,
                      onClick: () {
                        controller.handleEditPhoneClicked();
                      },
                      text: LocaleKeys.editphonenumber.tr,
                      variant: ButtonVariant.Link,
                      textColor: ColorName.primaryBlue,
                      textDecoration: TextDecoration.underline,
                    ),
                    Obx(() =>
                        controller.canResend.value ? _space() : SizedBox()),
                    Obx(() {
                      final canResend = controller.canResend.value;
                      return canResend
                          ? UBText(
                              fontStyle: FontStyle.italic,
                              text: LocaleKeys.or.tr,
                              color: ColorName.grey80,
                              size: 13,
                            )
                          : const SizedBox();
                    }),
                    Obx(() =>
                        controller.canResend.value ? _space() : SizedBox()),
                    Obx(() {
                      final canResend = controller.canResend.value;
                      return canResend
                          ? UBButton(
                              height: 16,
                              onClick: () {
                                controller.phoneNumberSubmitted();
                              },
                              text: LocaleKeys.resend.tr,
                              textColor: ColorName.primaryBlue,
                              variant: ButtonVariant.Link,
                              textDecoration: TextDecoration.underline,
                            )
                          : const SizedBox();
                    }),
                    _space(),
                  ],
                ),
              ),
              Stack(
                children: [
                  UBGreyContainer(
                    color: ColorName.black,
                    margin: const EdgeInsets.only(left: 12, right: 12, top: 24),
                    child: UBBorderlessInput(
                      type: TextInputType.number,
                      placeholder: LocaleKeys.enterCodeHere.tr,
                      onChange: controller.handleVerificationCodeChange,
                    ),
                  ),
                  Positioned(
                    child: Obx(() {
                      final timer = controller.countDownValue.value;
                      return UBText(
                        text: timer,
                        size: 13,
                        color: ColorName.grey80,
                      );
                    }),
                    right: 24,
                    bottom: 10,
                  )
                ],
              ),
              fill,
              Container(
                margin: EdgeInsets.symmetric(vertical: 24),
                width: 140,
                child: UBButton(
                  onClick: () {
                    Get.toNamed(AppPages.ACCOUNT);
                  },
                  variant: ButtonVariant.Rounded,
                  buttonColor: ColorName.grey23,
                  textColor: ColorName.primaryBlue,
                  text: 'Go To Dashboard',
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }

  SizedBox _space() {
    return const SizedBox(
      width: 4,
    );
  }
}
