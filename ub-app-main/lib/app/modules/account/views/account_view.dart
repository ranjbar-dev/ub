import 'package:flutter/material.dart';
import 'package:flutter_switch/flutter_switch.dart';
import 'package:get/get.dart';

import '../../../../generated/assets.gen.dart';
import '../../../../generated/colors.gen.dart';
import '../../../../generated/locales.g.dart';
import '../../../../utils/environment/ubEnv.dart';
import '../../../../utils/mixins/commonConsts.dart';
import '../../../../utils/mixins/toast.dart';
import '../../../../utils/throttle.dart';
import '../../../common/components/CenterUBLoading.dart';
import '../../../common/components/PaddedWrapper.dart';
import '../../../common/components/UBColumnAnimator.dart';
import '../../../common/components/UBText.dart';
import '../../../routes/app_pages.dart';
import '../accountRowModel.dart';
import '../controllers/account_controller.dart';
import '../widgets/accountCard.dart';

final thr = Throttling(duration: const Duration(seconds: 2));

class AccountView extends GetView<AccountController> with Toaster {
  throttleToast(String message) {
    thr.throttle(() => toastInfo(message));
  }

  @override
  Widget build(BuildContext context) {
    controller.pageLoaded();
    return SafeArea(
      child: UBPaddedWrapper(
        padding: 12,
        child: Obx(
          () {
            final userData = controller.accountData.value;

            String email = userData.email;
            if (email.length > 20) {
              email = email.substring(0, 10) +
                  "****" +
                  email.substring(email.length - 9, email.length);
            }

            return userData.email == null
                ? CenterUbLoading()
                : SizedBox(
                    height: Get.height,
                    child: SingleChildScrollView(
                      child: UBColumnVerticalSlideAnimator(
                        animationTime: 275,
                        children: [
                          vspace24,
                          Row(
                            children: [
                              Assets.images.maleAccountImage.svg(),
                              hspace8,
                              Column(
                                crossAxisAlignment: CrossAxisAlignment.start,
                                children: [
                                  UBText(
                                    text: LocaleKeys.userProfile.tr,
                                    color: ColorName.white,
                                    weight: FontWeight.w400,
                                    size: 32.0,
                                  ),
                                  UBText(
                                    text: 'User ID: ' + userData.ubId,
                                    color: ColorName.white,
                                    weight: FontWeight.w400,
                                    size: 13,
                                  )
                                ],
                              ),
                            ],
                          ),
                          vspace24,
                          AccountCard(
                            rows: [
                              AccountCardRowModel(
                                title: LocaleKeys.email.tr + ":",
                                value: email,
                                icon: Assets.images.emailIcon,
                                endButtonTitle: LocaleKeys.verify.tr,
                                inactiveEndButtonTitle: LocaleKeys.verifid.tr,
                                endButtonActive:
                                    userData.isAccountVerified != true,
                                onEndButtonClick:
                                    controller.requestForEmailVerification,
                              ),
                              //AccountCardRowModel(
                              //  icon: Assets.images.user,
                              //  title: LocaleKeys.userId.tr + ":",
                              //  value: userData.ubId,
                              //),
                              AccountCardRowModel(
                                icon: Assets.images.lock,
                                title: LocaleKeys.password.tr + ":",
                                value: '**********',
                                endButtonTitle: LocaleKeys.change.tr,
                                onEndButtonClick: () {
                                  if (userData.isAccountVerified == true) {
                                    Get.toNamed(AppPages.CHANGEPASSWORD);
                                    return;
                                  } else {
                                    throttleToast("Please verify your email");
                                  }
                                },
                                endButtonActive: true,
                              ),
                              AccountCardRowModel(
                                  isLast: true,
                                  icon: Assets.images.phone,
                                  title: LocaleKeys.phone.tr + ":",
                                  value: userData.phone != ''
                                      ? userData.phone.substring(0, 3) +
                                          '******' +
                                          userData.phone.substring(
                                              11, userData.phone.length)
                                      : "Please verify your phone",
                                  endButtonTitle: userData.phone != ''
                                      ? LocaleKeys.change.tr
                                      : LocaleKeys.verify.tr,
                                  endButtonActive: true,
                                  onEndButtonClick: () {
                                    if (userData.isAccountVerified == true) {
                                      Get.toNamed(AppPages.PHONE_VERIFICATION,
                                          arguments: {
                                            "twofaEnabled":
                                                userData.google2faEnabled
                                          });
                                      return;
                                    } else {
                                      throttleToast("Please verify your email");
                                    }
                                  }),
                            ],
                          ),
                          AccountCard(
                            rows: [
                              AccountCardRowModel(
                                icon: Assets.images.identity,
                                value: LocaleKeys.identityVerification.tr,
                                endButtonTitle: LocaleKeys.verify.tr,
                                inactiveEndButtonTitle: LocaleKeys.verifid.tr,
                                onEndButtonClick: () {
                                  if (userData.isAccountVerified == true) {
                                    Get.toNamed(AppPages.IDENTITYVERIFICATION);
                                    return;
                                  } else {
                                    throttleToast("Please verify your email");
                                  }
                                },
                                endButtonActive:
                                    userData.profileStatus != 'confirmed'
                                        ? true
                                        : false,
                              ),
                              AccountCardRowModel(
                                icon: Assets.images.g2fa,
                                value: LocaleKeys.googleAuthentication.tr,
                                endButtonTitle:
                                    userData.google2faEnabled == true
                                        ? LocaleKeys.disable.tr
                                        : LocaleKeys.enable.tr,
                                onEndButtonClick: () {
                                  if (userData.isAccountVerified == true) {
                                    Get.toNamed(
                                        AppPages.TWOFACTORAUTHENTICATION);
                                  } else {
                                    throttleToast("Please verify your email");
                                  }
                                },
                                isLast: !(controller.hasBiometrics.value),
                                endButtonActive: true,
                              ),
                              if (controller.hasBiometrics.value)
                                AccountCardRowModel(
                                  isLast: true,
                                  icon: Assets.images.biometrics,
                                  value: 'Biometric Authentication',
                                  //endButtonTitle:
                                  //    controller.isBiometricsActivated.value ==
                                  //            true
                                  //        ? LocaleKeys.disable.tr
                                  //        : LocaleKeys.enable.tr,
                                  //onEndButtonClick: () {
                                  //  controller.toggleBiometrics();
                                  //},
                                  //endButtonActive: true,
                                  endWidget: Obx(
                                    () => Container(
                                      width: 35,
                                      padding: const EdgeInsets.symmetric(
                                          vertical: 12.0),
                                      child: FlutterSwitch(
                                        padding: 3,
                                        toggleSize: 15,
                                        activeColor: ColorName.primaryBlue,
                                        value: controller
                                            .isBiometricsActivated.value,
                                        onToggle: (val) {
                                          controller.toggleBiometrics();
                                        },
                                      ),
                                    ),
                                  ),
                                ),
                            ],
                          ),
                          AccountCard(
                            rows: [
                              AccountCardRowModel(
                                icon: Assets.images.addressManagement,
                                value: LocaleKeys.withdrawAddressManagement.tr,
                                endButtonTitle: LocaleKeys.manage.tr,
                                onEndButtonClick: () {
                                  if (userData.isAccountVerified == true) {
                                    Get.toNamed(
                                        AppPages.WITHDRAWADDRESSMANAGEMENT);
                                    return;
                                  } else {
                                    throttleToast("Please verify your email");
                                  }
                                },
                                endButtonActive: true,
                              ),
                              AccountCardRowModel(
                                icon: Assets.images.listCheck,
                                value: LocaleKeys.loginWhiteList.tr,
                                endButtonTitle: LocaleKeys.disabled.tr,
                                endButtonActive: false,
                              ),
                              AccountCardRowModel(
                                icon: Assets.images.history,
                                value: LocaleKeys.recentLoginHistory.tr,
                                endButtonTitle: LocaleKeys.disabled.tr,
                                endButtonActive: false,
                              ),
                              AccountCardRowModel(
                                onEndButtonClick: () {
                                  controller.handleLogoutClick();
                                },
                                endButtonTitle: "Logout",
                                isLast: true,
                                icon: Assets.images.user,
                                value: 'Logout',
                                endButtonActive: true,
                              ),
                            ],
                          ),
                          UBText(text: 'Version: ${VERSION.split('+')[0]}')
                        ],
                      ),
                    ),
                  );
          },
        ),
      ),
    );
  }
}
