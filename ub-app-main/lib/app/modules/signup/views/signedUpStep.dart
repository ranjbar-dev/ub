import 'package:flutter/material.dart';
import '../../../common/components/UBButton.dart';
import '../../../common/components/UBText.dart';
import '../../../common/components/bottomCard.dart';
import '../../../routes/app_pages.dart';
import '../../../../generated/assets.gen.dart';
import '../../../../generated/colors.gen.dart';
import '../../../../generated/locales.g.dart';
import 'package:get/get.dart';
import 'package:android_intent_plus/android_intent.dart';
import '../../../../utils/logger.dart';
import 'package:url_launcher/url_launcher.dart';

class SignedUp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return BottomCard(
      title: LocaleKeys.create.tr + " " + LocaleKeys.account.tr,
      children: [
        Expanded(
          flex: 1,
          child: Column(
            children: [
              Assets.images.envelopeWithNumberOne.svg(),
              _space24(),
              UBText(text: LocaleKeys.youraccountiscreated.tr),
              _space8(),
              UBText(
                text: LocaleKeys.pleasecheckyouremail.tr,
                color: ColorName.orange,
              ),
              _space24(),
              SizedBox(
                width: GetPlatform.isWeb ? 80 : 130,
                child: UBButton(
                  variant: ButtonVariant.Rounded,
                  buttonColor: ColorName.grey16,
                  padding:
                      const EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                  textColor: ColorName.primaryBlue,
                  onClick: () {
                    if (GetPlatform.isWeb) {
                      Get.offNamed(AppPages.LOGIN);
                    } else {
                      if (GetPlatform.isAndroid) {
                        AndroidIntent intent = AndroidIntent(
                          action: 'android.intent.action.MAIN',
                          category: 'android.intent.category.APP_EMAIL',
                        );
                        intent.launch().catchError((e) {
                          log.e(e.toString());
                        });
                      } else if (GetPlatform.isIOS) {
                        launchUrl(Uri.parse("message://")).catchError((e) {
                          log.e(e.toString());
                        });
                      }
                    }
                  },
                  text: GetPlatform.isWeb
                      ? LocaleKeys.buttons_login.tr
                      : LocaleKeys.openEmailApp.tr,
                ),
              ),
            ],
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

  SizedBox _space8() {
    return const SizedBox(
      height: 8,
    );
  }
}
