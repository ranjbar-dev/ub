import 'package:flutter/material.dart';

import 'package:get/get.dart';
import '../../../common/components/UBButton.dart';
import '../../../common/components/UBScrollColumnExpandable.dart';
import '../../../common/components/UBSimpleInput.dart';
import '../../../common/components/UBText.dart';
import '../../../common/components/appbarTextTitle.dart';
import '../../../common/components/bottomCard.dart';
import '../../../common/components/pageContainer.dart';
import '../controllers/signup_controller.dart';
import 'signedUpStep.dart';
import '../../../routes/app_pages.dart';
import '../../../../generated/colors.gen.dart';
import '../../../../generated/locales.g.dart';
import '../../../../utils/mixins/commonConsts.dart';

class SignupView extends GetView<SignupController> {
  @override
  Widget build(BuildContext context) {
    controller.setStartTime();
    return PageContainer(
      protectedPage: false,
      appbarTitle: AppBarTextTitle(
        title: LocaleKeys.buttons_register.tr,
      ),
      child: UBScrollColumnExpandable(
        children: [
          fill,
          Container(
            height: controller.step.value == signupStep.SignedUp ? 325 : 440,
            child: Obx(
              () {
                if (controller.step.value == signupStep.SignedUp) {
                  return SignedUp();
                }
                return BottomCard(
                  title: LocaleKeys.create.tr + " " + LocaleKeys.account.tr,
                  children: [
                    Expanded(
                      flex: 2,
                      child: UBSimpleInput(
                        onChange: controller.handleEmailChange,
                        placeHolder: LocaleKeys.email.tr,
                        error: controller.emailError.value,
                      ),
                    ),
                    Expanded(
                      flex: 2,
                      child: UBSimpleInput(
                        onChange: controller.handlePasswordChange,
                        placeHolder: LocaleKeys.password.tr,
                        error: controller.passwordError.value,
                        isSecure: true,
                        isPickable: true,
                      ),
                    ),
                    Expanded(
                      flex: 2,
                      child: UBSimpleInput(
                        onChange: controller.handleRepeatPasswordChange,
                        placeHolder: LocaleKeys.repeatpassword.tr,
                        isSecure: true,
                        isPickable: true,
                        error: controller.repeatPasswordError.value,
                      ),
                    ),
                    Expanded(
                      flex: 2,
                      child: Container(
                          padding: const EdgeInsets.only(
                            bottom: 12,
                            top: 7,
                          ),
                          child: Obx(() {
                            final email = controller.email.value;
                            final password = controller.password.value;
                            final repeatPassword =
                                controller.repeatPassword.value;
                            return UBButton(
                              height: 30,
                              disabled: password == '' ||
                                  email == '' ||
                                  repeatPassword == '',
                              isLodaing: controller.isLoading.value,
                              onClick: controller.handleSubmitClick,
                              text: LocaleKeys.submit.tr,
                            );
                          })),
                    ),
                    Expanded(
                      flex: 1,
                      child: Container(
                        padding: const EdgeInsets.only(bottom: 8),
                        child: Row(
                          mainAxisAlignment: MainAxisAlignment.center,
                          children: [
                            UBText(
                              color: ColorName.greybf,
                              text: 'Have an account?',
                            ),
                            const SizedBox(
                              width: 8,
                            ),
                            UBButton(
                              variant: ButtonVariant.Link,
                              text: LocaleKeys.buttons_login.tr,
                              onClick: () => {Get.offAllNamed(AppPages.LOGIN)},
                            )
                          ],
                        ),
                      ),
                    ),
                  ],
                );
              },
            ),
          ),
        ],
      ),
    );
  }
}
