import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../../generated/assets.gen.dart';
import '../../../../generated/colors.gen.dart';
import '../../../../generated/locales.g.dart';
import '../../../../utils/mixins/commonConsts.dart';
import '../../../../utils/mixins/toast.dart';
import '../../../common/components/UBButton.dart';
import '../../../common/components/UBGreyContainer.dart';
import '../../../common/components/UBScrollColumnExpandable.dart';
import '../../../common/components/UBSimpleInput.dart';
import '../../../common/components/UBText.dart';
import '../../../common/custom/rflutter_alert/rflutter_alert.dart';
import '../controllers/two_factor_authentication_controller.dart';

class EnterPasswordAndCode extends GetView<TwoFactorAuthenticationController>
    with Toaster {
  @override
  Widget build(BuildContext context) {
    final isDisabling = controller.isEnabled.value == true;
    return UBScrollColumnExpandable(
      children: [
        Container(
          child: Column(
            children: [
              vspace12,
              UBGreyContainer(
                height: 70,
                margin: EdgeInsets.symmetric(horizontal: 12),
                child: Row(
                  children: [
                    const Icon(
                      Icons.warning_amber_rounded,
                      color: ColorName.orange,
                    ),
                    Container(
                      height: 50,
                      width: 2,
                      margin: const EdgeInsets.symmetric(horizontal: 12),
                      color: ColorName.black2c,
                    ),
                    UBText(
                      wrapped: true,
                      text: isDisabling == true
                          ? LocaleKeys.disable2FaWarning.tr
                          : "For security reasons, after enabling/disabling 2fa,you will not be able to withdraw for 24 hours",
                      size: 13,
                      color: ColorName.orange,
                    ),
                  ],
                ),
              ),
              vspace24,
              if (isDisabling) vspace24,
              if (isDisabling) vspace12,
              if (isDisabling) Assets.images.twofaicon.svg(),
              if (isDisabling) vspace24,
            ],
          ),
        ),
        fill,
        Container(
          height: isDisabling ? 300.0 : 410.0,
          decoration: const BoxDecoration(
              color: ColorName.black1c, borderRadius: roundedTop_big),
          padding: const EdgeInsets.symmetric(horizontal: 12.0),
          child: UBScrollColumnExpandable(
            children: [
              if (!isDisabling) vspace24,
              if (!isDisabling) Assets.images.shieldWithKey.svg(),
              vspace24,
              UBSimpleInput(
                onChange: controller.handlePasswordChange,
                placeHolder: LocaleKeys.password.tr,
                isSecure: true,
                isPickable: true,
                //error: controller.repeatPasswordError.value,
              ),
              vspace24,
              UBSimpleInput(
                placeHolder: 'Google authentication code',
                onChange: controller.handleCodeChange,
                //error: controller.repeatPasswordError.value,
              ),
              fill,
              Padding(
                padding: const EdgeInsets.symmetric(
                  horizontal: 12,
                ),
                child: Obx(
                  () {
                    final password = controller.password.value;
                    final code = controller.code.value;
                    final isLoading = controller.isFinalSubmitLoading.value;
                    return UBButton(
                        isLodaing: isLoading,
                        disabled: password.length < 8 || code.length < 6,
                        onClick: () {
                          if (isDisabling) {
                            _openConfirmation();
                          } else {
                            controller.handleFinalSubmitClick(
                              enable: !controller.isEnabled.value,
                            );
                          }
                        },
                        text: LocaleKeys.submit.tr);
                  },
                ),
              ),
              Container(
                width: 120,
                padding:
                    const EdgeInsets.symmetric(horizontal: 12, vertical: 24),
                child: UBButton(
                  onClick: () {
                    //_openConfirmation();
                    Get.back();
                  },
                  variant: ButtonVariant.TransparentBackground,
                  textColor: ColorName.grey80,
                  text: LocaleKeys.cancel.tr,
                ),
              )
            ],
          ),
        ),
      ],
    );
  }

  void _openConfirmation() {
    Alert(
      style: AlertStyle(
        animationType: AnimationType.grow,
      ),
      context: Get.context,
      content: Container(
        child: Column(
          children: [
            UBText(
              text: LocaleKeys.disable2Fa.tr + ' ?',
              color: ColorName.red,
            ),
            vspace8,
            RichText(
                textAlign: TextAlign.center,
                text: TextSpan(children: [
                  TextSpan(
                    text: "After this action, withdrawal will be disabled for ",
                    style: grey80Bold13,
                  ),
                  TextSpan(
                    text: "24 hours",
                    style: redBold13,
                  ),
                  TextSpan(
                    text: ", are you sure you want to proceed?",
                    style: grey80Bold13,
                  ),
                ])),
            vspace24,
            Container(
              child: Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  GestureDetector(
                      child: UBText(text: LocaleKeys.cancel.tr),
                      onTap: () {
                        Get.back();
                        Get.back();
                      }),
                  Container(
                    height: 24,
                    width: 1,
                    margin: const EdgeInsets.symmetric(horizontal: 24),
                    color: ColorName.white,
                  ),
                  GestureDetector(
                    child: UBText(
                      align: TextAlign.center,
                      text: "I'm Sure",
                      color: ColorName.red,
                    ),
                    onTap: () {
                      Get.back();
                      controller.handleFinalSubmitClick(
                        enable: !controller.isEnabled.value,
                      );
                    },
                  ),
                ],
              ),
            )
          ],
        ),
      ),
    ).show();
  }
}
