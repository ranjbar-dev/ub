import 'package:flutter/material.dart';

import 'package:get/get.dart';
import '../../../common/components/UBButton.dart';
import '../../../common/components/UBMessagePage.dart';
import '../../../common/components/UBScrollColumnExpandable.dart';
import '../../../common/components/UBSimpleInput.dart';
import '../../../common/components/UBText.dart';
import '../../../common/components/appbarTextTitle.dart';
import '../../../common/components/pageContainer.dart';
import '../../../routes/app_pages.dart';
import '../../../../generated/assets.gen.dart';
import '../../../../generated/colors.gen.dart';
import '../../../../generated/locales.g.dart';
import '../../../../utils/mixins/commonConsts.dart';

import '../controllers/change_password_controller.dart';

class ChangePasswordView extends GetView<ChangePasswordController> {
  @override
  Widget build(BuildContext context) {
    return Obx(
      () {
        return PageContainer(
          appbarTitle: AppBarTextTitle(
            title: LocaleKeys.changePassword.tr,
          ),
          child: controller.step.value == changePasswordStep.PasswordChanged
              ? UBMessagePage(
                  children: [
                    Assets.images.passwordChangedIcon.svg(),
                    UBText(
                      text: LocaleKeys.yourpasswordhasbeenchanged.tr,
                    ),
                    SizedBox(
                      width: 140,
                      child: UBButton(
                        variant: ButtonVariant.Rounded,
                        buttonColor: ColorName.black2c,
                        textColor: ColorName.primaryBlue,
                        onClick: () {
                          Get.offAllNamed(AppPages.ACCOUNT);
                        },
                        text: LocaleKeys.gotoDashboard.tr,
                      ),
                    )
                  ],
                )
              : UBScrollColumnExpandable(
                  children: [
                    fill,
                    Container(
                      decoration: const BoxDecoration(
                        borderRadius: roundedTop_big,
                        color: ColorName.black2c,
                      ),
                      padding: const EdgeInsets.symmetric(horizontal: 12),
                      child: Column(
                        children: [
                          vspace48,
                          Assets.images.changePassword.svg(),
                          vspace48,
                          Padding(
                            padding: const EdgeInsets.only(bottom: 24),
                            child: Obx(
                              () => UBSimpleInput(
                                isPickable: true,
                                isSecure: true,
                                placeHolder: LocaleKeys.oldPawwsord.tr,
                                error: controller.oldPasswordError.value,
                                onChange:
                                    controller.handleOldPasswordValueChange,
                              ),
                            ),
                          ),
                          Padding(
                            padding: const EdgeInsets.only(bottom: 24),
                            child: Obx(
                              () => UBSimpleInput(
                                isPickable: true,
                                isSecure: true,
                                placeHolder: LocaleKeys.newPassword.tr,
                                error: controller.newPasswordError.value,
                                onChange:
                                    controller.handleNewPasswordValueChange,
                              ),
                            ),
                          ),
                          Padding(
                            padding: const EdgeInsets.only(bottom: 24),
                            child: Obx(
                              () => UBSimpleInput(
                                isPickable: true,
                                isSecure: true,
                                placeHolder: LocaleKeys.repeatNewPassword.tr,
                                error: controller.repeatNewPasswordError.value,
                                onChange: controller
                                    .handleRepeatNewPasswordValueChange,
                              ),
                            ),
                          ),
                          Padding(
                            padding: const EdgeInsets.only(bottom: 24),
                            child: UBButton(
                              isLodaing: controller.isSubmitting.value,
                              disabled: !controller.canSubmit(),
                              onClick: () {
                                controller.handleSubmitClick();
                              },
                              text: LocaleKeys.changePassword.tr,
                            ),
                          ),
                          Padding(
                            padding: const EdgeInsets.only(bottom: 24),
                            child: Container(
                              width: 140,
                              child: UBButton(
                                onClick: () {
                                  Navigator.pop(context);
                                },
                                text: LocaleKeys.cancel.tr,
                                variant: ButtonVariant.TransparentBackground,
                                buttonColor: ColorName.grey80,
                              ),
                            ),
                          ),
                        ],
                      ),
                    )
                  ],
                ),
        );
      },
    );
  }
}
