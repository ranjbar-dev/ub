import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../common/components/UBButton.dart';
import '../../../common/components/UBColoredOverlay.dart';
import '../../../common/components/UBConnectionLost.dart';
import '../../../common/components/UBText.dart';
import '../../../../generated/assets.gen.dart';
import '../../../../generated/colors.gen.dart';
import 'package:supercharged/supercharged.dart';
import '../../../../utils/environment/ubEnv.dart';

import '../controllers/after_splash_controller.dart';

class AfterSplashView extends GetView<AfterSplashController> {
  @override
  Widget build(BuildContext context) {
    return Material(
      child: SafeArea(
        child: Stack(
          children: [
            Container(
              width: Get.width,
              color: ColorName.black,
            ),
            Positioned(
              child: Obx(
                () {
                  final animationStarted = controller.animationStarted.value;
                  return AnimatedOpacity(
                    curve: Curves.easeInOut,
                    opacity: animationStarted ? 1 : 0,
                    duration: afterSplashAnimationDuration + 100.milliseconds,
                    child: Assets.images.watermark.svg(
                      width: Get.width,
                    ),
                  );
                },
              ),
            ),
            Obx(() {
              final animationStarted = controller.animationStarted.value;
              return Container(
                child: AnimatedPositioned(
                    top: animationStarted ? (-Get.height / 3.5) : Get.height,
                    curve: Curves.easeInOut,
                    duration: afterSplashAnimationDuration,
                    child: UBColoredOverlay()
                    // color: ColorName.textBlue,
                    ),
              );
            }),
            Column(
              children: [
                Expanded(
                  child: Center(
                    child: Obx(
                      () {
                        final animationStarted =
                            controller.animationStarted.value;
                        return AnimatedOpacity(
                          curve: Curves.easeInOut,
                          opacity: animationStarted ? 1 : 0,
                          duration: afterSplashAnimationDuration,
                          child: Hero(
                            tag: 'logo',
                            child: Assets.images.logoSvg.svg(),
                          ),
                        );
                      },
                    ),
                    // Assets.images.fullWidthIconwithAppName.svg(),
                  ),
                ),
              ],
            ),
            Obx(() {
              final connected = controller.isConnected.value;
              return connected
                  ? const SizedBox()
                  : ConnectionLost(text: 'Check your internet Connection');
            }),
            if (ENV == "DEV")
              UBText(
                text: "Version: $VERSION",
                size: 18,
                color: ColorName.red,
              ),
            Obx(() {
              final showRetryButton = controller.showRetryButton.value;
              final isLoading = controller.isLoadingRetry.value;
              return !showRetryButton
                  ? const SizedBox()
                  : Positioned(
                      bottom: 24.0,
                      right: (Get.width / 2) - 35,
                      child: Container(
                        width: 70,
                        child: UBButton(
                          isLodaing: isLoading,
                          variant: ButtonVariant.Rounded,
                          onClick: () {
                            controller.handleRetryButton();
                          },
                          text: 'Retry',
                        ),
                      ),
                    );
            }),
          ],
        ),
      ),
    );
  }
}
