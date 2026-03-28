import 'package:ai_barcode/ai_barcode.dart';
import 'package:flutter/material.dart';

import 'package:get/get.dart';
import '../../../common/components/UBButton.dart';
import '../../../common/components/UBText.dart';
import '../../../../generated/assets.gen.dart';
import '../../../../generated/colors.gen.dart';
import '../../../../utils/commonUtils.dart';
import '../../../../utils/mixins/commonConsts.dart';

import '../controllers/qr_scan_controller.dart';

class QrScanView extends GetView<QrScanController> {
  @override
  Widget build(BuildContext context) {
    controller.startCameraUsage();
    return Material(
      child: Stack(
        children: [
          SizedBox(
            height: Get.height,
            width: Get.width,
            child: Column(
              children: [
                Container(
                  height: Get.height,
                  width: Get.width,
                  color: ColorName.black2c,
                  child: Column(
                    children: [
                      Expanded(
                        child: PlatformAiBarcodeScannerWidget(
                          platformScannerController: controller.scanner,
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),
          Positioned(
            bottom: 24,
            child: SizedBox(
                height: 80,
                width: Get.width,
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    _columnButton(
                        text: 'Back',
                        onTap: () {
                          Get.back();
                        },
                        icon: Assets.images.arrowLeft.svg()),
                    if (!(GetPlatform.isWeb))
                      const SizedBox(
                        width: 24,
                      ),
                    if (!(GetPlatform.isWeb))
                      _columnButton(
                        text: 'Gallery',
                        onTap: () {
                          controller.handleBrowsClick();
                        },
                        icon: Assets.images.gallery.svg(),
                      ),
                  ],
                )),
          ),
          Positioned(
            top: 100,
            left: (Get.width / 2) - 50,
            child: UBText(
              text: 'Scan Qr Code',
              size: 16,
              color: ColorName.white,
            ),
          ),
          Center(child: Obx(() {
            final cameraPermissionDenied =
                controller.deniedCameraPermission.value;
            return Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                const SizedBox(
                  height: 400.0,
                ),
                if (!cameraPermissionDenied)
                  UBText(
                    text: 'Please Scan the Qr code or',
                    size: 16,
                    color: ColorName.white,
                  ),
                if (cameraPermissionDenied) vspace12,
                if (cameraPermissionDenied)
                  Container(
                    width: 280,
                    child: UBButton(
                        height: 24,
                        onClick: () {
                          Get.back();
                          Get.back();
                          promptForPermissionInSetting(
                            onDenied: () {},
                            title: 'Camera Permission',
                            desc:
                                'This app needs camera access to scan qr codes',
                          );
                        },
                        text: 'Grant camera permission to scan'),
                  ),
                if (cameraPermissionDenied) vspace12,
                UBText(
                  text: 'or',
                  size: 16,
                  color: ColorName.white,
                ),
                if (cameraPermissionDenied) vspace12,
                UBText(
                  text: 'Select Qr code image from gallery',
                  size: 16,
                  color: ColorName.white,
                )
              ],
            );
          })),
        ],
      ),
    );
  }

  _space8() {
    return const SizedBox(
      height: 8,
    );
  }

  _columnButton({String text, Null Function() onTap, Widget icon}) {
    return SizedBox(
      width: 50,
      height: 70,
      child: Column(
        children: [
          GestureDetector(
            onTap: () {
              onTap();
            },
            child: Container(
              width: 40,
              height: 40,
              decoration: BoxDecoration(
                color: ColorName.grey16,
                borderRadius: BorderRadius.circular(100),
              ),
              child: Align(
                child: icon,
              ),
            ),
          ),
          _space8(),
          UBText(
            text: text,
            color: ColorName.white,
          )
        ],
      ),
    );
  }
}
