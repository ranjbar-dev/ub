import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:url_launcher/url_launcher.dart';

import '../../../../generated/assets.gen.dart';
import '../../../../generated/colors.gen.dart';
import '../../../../generated/locales.g.dart';
import '../../../common/components/UBButton.dart';
import '../../../common/components/UBText.dart';
import '../../../routes/app_pages.dart';
import '../controllers/two_factor_authentication_controller.dart';

class FinalStatus extends GetView<TwoFactorAuthenticationController> {
  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        _space24(),
        if (controller.isEnabled.value)
          Assets.images.twofaEnabled.svg()
        else
          Assets.images.twofaicon.svg(),
        _space24(),
        UBText(
          text: LocaleKeys.done.tr + "!",
          size: 14,
          weight: FontWeight.bold,
          color: ColorName.primaryBlue,
        ),
        _space12(),
        Row(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            UBText(
              text: LocaleKeys.googleAuthenticator.tr,
              size: 13,
              color: ColorName.grey80,
            ),
            const SizedBox(width: 2),
            UBText(
              text: controller.isEnabled.value
                  ? LocaleKeys.enabled.tr
                  : LocaleKeys.disabled.tr,
              size: 13,
              color: ColorName.primaryBlue,
            ),
          ],
        ),
        _space24(),
        _space24(),
        Container(
          width: 130,
          child: UBButton(
            variant: ButtonVariant.Rounded,
            buttonColor: ColorName.grey23,
            textColor: ColorName.primaryBlue,
            onClick: () {
              Get.offNamed(AppPages.ACCOUNT);
              Future.delayed(100.milliseconds).then((value) {
                Get.delete<TwoFactorAuthenticationController>();
              });
            },
            text: LocaleKeys.gotoDashboard.tr,
          ),
        ),
      ],
    );
  }

  SizedBox _space24() {
    return const SizedBox(
      height: 24,
    );
  }

  SizedBox _space12() {
    return const SizedBox(
      height: 12,
    );
  }

  launchURL(String url) async {
    if (await canLaunchUrl(Uri.parse(url))) {
      await launchUrl(Uri.parse(url));
    } else {
      throw 'Could not launch $url';
    }
  }
}
