import 'package:flutter/material.dart';

import 'package:get/get.dart';

import '../../../common/components/PaddedWrapper.dart';
import '../../../common/components/UBButton.dart';
import '../../../common/components/UBColoredOverlay.dart';
import '../../../common/components/pageContainer.dart';
import '../../../routes/app_pages.dart';
import '../../../../generated/assets.gen.dart';
import '../../../../generated/locales.g.dart';
import '../../../../utils/mixins/commonConsts.dart';
import '../../../../generated/colors.gen.dart';
import '../controllers/landing_controller.dart';

class LandingView extends GetView<LandingController> {
  @override
  Widget build(BuildContext context) {
    return PageContainer(
      protectedPage: false,
      child: Stack(
        children: [
          Positioned(
            child: Container(
              clipBehavior: Clip.none,
              width: Get.width,
              child: Assets.images.watermark.svg(
                width: Get.width,
                clipBehavior: Clip.none,
              ),
            ),
          ),
          Obx(() {
            final animationStarted = controller.isLoaded.value;
            return Container(
              child: AnimatedPositioned(
                  top: animationStarted ? -(Get.height) : (-Get.height / 3.5),
                  curve: Curves.easeInOut,
                  duration: 1.seconds,
                  child: UBColoredOverlay()
                  // color: ColorName.textBlue,
                  ),
            );
          }),
          UBPaddedWrapper(
            child: Column(
              children: [
                const Spacer(),
                Hero(
                  child: Assets.images.logoSvg.svg(
                    fit: BoxFit.fitWidth,
                  ),
                  tag: 'logo',
                ),
                Obx(() {
                  return AnimatedOpacity(
                      opacity: controller.isLoaded.value == true ? 1 : 0,
                      duration: 1.seconds,
                      child: SizedBox(
                        height: 310.0,
                        child: Column(
                          crossAxisAlignment: CrossAxisAlignment.center,
                          children: [
                            vspace12,
                            RichText(
                                text: TextSpan(
                              text: 'Welcome',
                              style: const TextStyle(
                                color: ColorName.white,
                                fontSize: 18,
                              ),
                            )),
                            vspace12,
                            SizedBox(
                              width: 250.0,
                              child: RichText(
                                  textAlign: TextAlign.center,
                                  text: TextSpan(
                                    text: 'Enjoy Trading Safely',
                                    style: const TextStyle(
                                        color: ColorName.greyd8,
                                        fontSize: 14.0,
                                        height: 1.4),
                                  )),
                            ),
                            const Spacer(),
                            UBButton(
                              text: LocaleKeys.buttons_register.tr,
                              onClick: () => {Get.toNamed(AppPages.SIGNUP)},
                            ),
                            const Spacer(),
                            RichText(
                              text: TextSpan(
                                text: 'Already have an account?',
                                style: const TextStyle(
                                  color: ColorName.grey80,
                                  fontSize: 14.0,
                                  fontWeight: FontWeight.w600,
                                ),
                              ),
                            ),
                            vspace12,
                            UBButton(
                              width: 100,
                              height: 24,
                              variant: ButtonVariant.Link,
                              text: LocaleKeys.buttons_login.tr,
                              onClick: () => {Get.toNamed(AppPages.LOGIN)},
                            ),
                            vspace16
                          ],
                        ),
                      ));
                }),
              ],
            ),
          ),
        ],
      ),
    );
  }
}
