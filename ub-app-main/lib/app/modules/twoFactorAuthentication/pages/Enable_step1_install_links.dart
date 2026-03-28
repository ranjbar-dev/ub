import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../common/components/UBButton.dart';
import '../../../common/components/UBText.dart';
import '../controllers/two_factor_authentication_controller.dart';
import '../../../../utils/commonUtils.dart';
import '../../../../utils/mixins/commonConsts.dart';
import '../../../../generated/assets.gen.dart';
import '../../../../generated/colors.gen.dart';
import '../../../../generated/locales.g.dart';

class InstallLinks extends GetView<TwoFactorAuthenticationController> {
  @override
  Widget build(BuildContext context) {
    // Future.delayed(200.milliseconds).then((value) {
    //   if (controller.isEnabled.value == false) {
    //     controller.openIntro();
    //   }
    // });
    return Column(
      children: [
        vspace24,
        Assets.images.twofaicon.svg(),
        vspace24,
        UBText(
          text: LocaleKeys.twofaMessageLine1.tr,
          size: 13,
          color: ColorName.grey80,
        ),
        UBText(
          text: LocaleKeys.twofaMessageLine2.tr,
          size: 13,
          color: ColorName.grey80,
        ),
        vspace24,
        SizedBox(
          height: 40,
          child: Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              GestureDetector(
                onTap: () {
                  const url =
                      'https://apps.apple.com/us/app/google-authenticator/id388497605';
                  launchURL(url);
                },
                child: Assets.images.appStoreIcon.svg(),
              ),
              const SizedBox(width: 8),
              GestureDetector(
                onTap: () {
                  const url =
                      'https://play.google.com/store/apps/details?id=com.google.android.apps.authenticator2&';
                  launchURL(url);
                },
                child: Assets.images.googlePlayIcon.svg(),
              ),
            ],
          ),
        ),
        fill,
        Padding(
          padding: const EdgeInsets.symmetric(
            horizontal: 12,
          ),
          child: Obx(() {
            final isLoding = controller.isLoadingCharCode.value;
            return UBButton(
              isLodaing: isLoding,
              onClick: () {
                controller.handleGoToCharacterCodeClick();
              },
              text: LocaleKeys.next.tr,
            );
          }),
        ),
        Container(
          width: 120,
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 24),
          child: UBButton(
            onClick: () {
              Get.back();
            },
            variant: ButtonVariant.TransparentBackground,
            textColor: ColorName.grey80,
            text: LocaleKeys.cancel.tr,
          ),
        )
      ],
    );
  }
}
