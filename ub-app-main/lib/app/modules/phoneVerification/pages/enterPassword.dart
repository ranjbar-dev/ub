import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../../generated/assets.gen.dart';
import '../../../../generated/colors.gen.dart';
import '../../../../generated/locales.g.dart';
import '../../../../utils/mixins/commonConsts.dart';
import '../../../common/components/UBButton.dart';
import '../../../common/components/UBSimpleInput.dart';
import '../../../common/components/UBText.dart';
import '../controllers/phone_verification_controller.dart';

class EnterPassword extends GetView<PhoneVerificationController> {
  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        fill,
        Container(
          height: 410.0,
          decoration: const BoxDecoration(
              color: ColorName.black2c, borderRadius: roundedTop_big),
          child: Column(
            children: [
              vspace24,
              Assets.images.shieldWithKey.svg(),
              vspace24,
              Container(
                width: 240,
                child: UBText(
                  lineHeight: 1.6,
                  size: 13,
                  color: ColorName.grey80,
                  align: TextAlign.center,
                  text: LocaleKeys.pleaseenteryourpassword.tr,
                ),
              ),
              Obx(
                () {
                  final is2faActivated = controller.is2faActivated.value;
                  return is2faActivated
                      ? const SizedBox()
                      : Container(
                          child: UBText(
                            lineHeight: 1.6,
                            size: 13,
                            color: ColorName.grey80,
                            align: TextAlign.center,
                            text: LocaleKeys
                                .forsecurityreasonswerecommendtoenable.tr,
                          ),
                        );
                },
              ),
              Obx(
                () {
                  final is2faActivated = controller.is2faActivated.value;
                  return is2faActivated
                      ? const SizedBox()
                      : Container(
                          width: 150,
                          alignment: Alignment.center,
                          child: Row(
                            children: [
                              UBText(
                                size: 13,
                                color: ColorName.grey80,
                                align: TextAlign.center,
                                text: LocaleKeys.your2Fa.tr,
                              ),
                              _space(),
                              UBButton(
                                height: 16,
                                onClick: () {
                                  controller.handleGoTo2faPageClick();
                                },
                                text: LocaleKeys.goto2Fapage.tr,
                                textColor: ColorName.primaryBlue,
                                variant: ButtonVariant.Link,
                                textDecoration: TextDecoration.underline,
                              ),
                            ],
                          ),
                        );
                },
              ),
              fill,
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 12.0),
                child: UBSimpleInput(
                  isPickable: true,
                  isSecure: true,
                  placeHolder: LocaleKeys.password.tr,
                  onChange: controller.handlePasswordChange,
                ),
              ),
              vspace24,
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 12),
                child: Obx(
                  () {
                    final password = controller.password.value;
                    return UBButton(
                      height: 36,
                      isLodaing: controller.isRequestingtoSubmitCode.value,
                      disabled: (!(password.length > 7)),
                      onClick: controller.finalSubmitClicked,
                      text: LocaleKeys.submit.tr,
                    );
                  },
                ),
              ),
              Container(
                margin: const EdgeInsets.symmetric(vertical: 24),
                width: 120,
                child: UBButton(
                  variant: ButtonVariant.TransparentBackground,
                  onClick: () {
                    Get.back();
                  },
                  text: LocaleKeys.cancel.tr,
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
