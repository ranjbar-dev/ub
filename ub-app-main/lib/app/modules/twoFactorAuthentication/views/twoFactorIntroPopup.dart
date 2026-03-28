import 'package:carousel_slider/carousel_controller.dart';
import 'package:flutter/material.dart';

import 'package:get/get.dart';
import '../../../common/components/UBCarousel.dart';
import '../../../common/components/UBText.dart';
import '../../../../generated/assets.gen.dart';
import '../../../../generated/colors.gen.dart';
import '../../../../utils/mixins/commonConsts.dart';

import '../controllers/two_factor_authentication_controller.dart';

class TwoFactorIntroPopup extends GetView<TwoFactorAuthenticationController> {
  @override
  Widget build(BuildContext context) {
    CarouselController controller = CarouselController();
    return Container(
      height: 500,
      child: Align(
        child: Container(
          width: Get.width - 32,
          height: 420.0,
          decoration: const BoxDecoration(
            borderRadius: BorderRadius.all(
              Radius.circular(8.0),
            ),
            color: ColorName.black2c,
          ),
          child: Column(
            children: [
              vspace16,
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 16.0),
                child: Row(
                  children: [
                    const Spacer(),
                    GestureDetector(
                        onTap: () => Get.back(),
                        child: Container(
                          color: Colors.transparent,
                          child: SizedBox(
                            width: 24.0,
                            height: 24.0,
                            child: Assets.images.closeIcon.svg(
                              color: Colors.white,
                            ),
                          ),
                        )),
                  ],
                ),
              ),
              SizedBox(
                width: 240.0,
                child: UBText(
                    align: TextAlign.center,
                    text:
                        'You’ll need to add a UnitetBit account to your Google Authenticator app and manually enter the 16-digit key.'),
              ),
              vspace24,
              UBCarousel(
                controller: controller,
                showNavigationArrows: false,
                showIndicators: true,
                slides: [
                  Assets.images.twofaIntroSlider1Png.image(),
                  Assets.images.twofaIntroSlider2Png.image(),
                  Assets.images.twofaIntroSlider3.svg(),
                  Assets.images.twofaIntroSlider4.svg(),
                ],
                height: 240,
                onChange: (i) {},
              ),
            ],
          ),
        ),
      ),
    );
  }
}
