import 'package:flutter/material.dart';

import 'package:get/get.dart';

import '../../../common/components/UBScrollColumnExpandable.dart';
// import 'package:unitedbit/app/common/custom/captcha/captchav2Html.dart';
// import 'package:unitedbit/app/common/custom/captcha/recaptcha2.dart';
// import 'package:unitedbit/app/common/custom/ubSlideToAct/UBSlideToAct.dart';
import 'package:unitedbit/utils/mixins/commonConsts.dart';

import '../../../../generated/colors.gen.dart';
import '../../../../generated/locales.g.dart';
import '../../../common/components/UBButton.dart';
import '../../../common/components/UBSimpleInput.dart';
import '../../../common/components/UBText.dart';
import '../../../common/components/appbarTextTitle.dart';
import '../../../common/components/bottomCard.dart';
import '../../../common/components/pageContainer.dart';
import '../../../routes/app_pages.dart';
import '../controllers/login_controller.dart';

class LoginView extends GetView<LoginController> {
  // final GlobalKey<SlideActionState> _key = GlobalKey();
  @override
  Widget build(BuildContext context) {
    controller.setStartTime();

    //controller.checkForBiometricLogin();
    return PageContainer(
      protectedPage: false,
      appbarTitle: AppBarTextTitle(
        title: LocaleKeys.pageTitles_login.tr,
      ),
      child: Stack(
        children: [
          Container(
            child: UBScrollColumnExpandable(
              children: [
                fill,
                // SlideAction(
                //   key: _key,
                //   onSlide: controller.onSlide,
                //   onSubmit: () async {
                //     controller.generateKeyPair();
                //     _key.currentState.reset();
                //     controller.resetramziVars();
                //     // Future.delayed(
                //     //   Duration(seconds: 1),
                //     //   () {
                //     //     _key.currentState.reset();
                //     //   },
                //     // );
                //   },
                // ),
                Container(
                  height: 460,
                  child: BottomCard(
                    title: LocaleKeys.titles_loginToContinue.tr,
                    children: [
                      Obx(() => UBSimpleInput(
                            placeHolder: LocaleKeys.email.tr,
                            type: TextInputType.emailAddress,
                            onChange: controller.handleEmailChange,
                            error: controller.loginEmailError.value,
                          )),
                      vspace24,
                      Obx(
                        () => UBSimpleInput(
                          isPickable: true,
                          isSecure: true,
                          placeHolder: LocaleKeys.password.tr,
                          onChange: (e) => {controller.loginPassword.value = e},
                          error: controller.loginPasswordError.value,
                        ),
                      ),
                      vspace16,
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          fill,
                          UBButton(
                            text:
                                "${LocaleKeys.forgot.tr} ${LocaleKeys.password.tr}?",
                            variant: ButtonVariant.Link,
                            height: 20.0,
                            onClick: () => {Get.toNamed(AppPages.FORGOT)},
                          ),
                        ],
                      ),
                      vspace12,
                      Container(
                        padding: const EdgeInsets.symmetric(
                          vertical: 5,
                        ),
                        width: double.infinity,
                        child: Obx(
                          () => UBButton(
                            height: 36.0,
                            disabled:
                                controller.loginPassword.value.length < 8 ||
                                    controller.loginEmailError.value != '',
                            onClick: controller.login,
                            text: LocaleKeys.buttons_login.tr,
                            isLodaing: controller.isLoggingIn.value,
                          ),
                        ),
                      ),
                      vspace24,
                      vspace12,
                      Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          UBText(
                            color: ColorName.grey80,
                            text: LocaleKeys.DontHaveAccount.tr,
                          ),
                          vspace8,
                          SizedBox(
                            width: 70,
                            height: 30,
                            child: UBButton(
                              variant: ButtonVariant.Link,
                              text: LocaleKeys.buttons_register.tr,
                              textDecoration: TextDecoration.underline,
                              textColor: ColorName.white,
                              onClick: () => {Get.toNamed(AppPages.SIGNUP)},
                            ),
                          )
                        ],
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),

        ],
      ),
    );
  }
}
