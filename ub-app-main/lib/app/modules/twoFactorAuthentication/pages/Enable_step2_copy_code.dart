import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:get/get.dart';
import 'package:transparent_image/transparent_image.dart';
import '../../../common/components/UBButton.dart';
import '../../../common/components/UBGreyContainer.dart';
import '../../../common/components/UBText.dart';
import '../controllers/two_factor_authentication_controller.dart';
import '../../../common/components/UBWarningRow.dart';
import '../../../../utils/mixins/commonConsts.dart';
import '../../../../utils/mixins/toast.dart';

import '../../../../generated/colors.gen.dart';
import '../../../../generated/locales.g.dart';
import '../../../../utils/throttle.dart';

final thr = Throttling(duration: const Duration(milliseconds: 4000));

class CopyCode extends GetView<TwoFactorAuthenticationController> with Toaster {
  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        fill,
        Container(
          height: 510.0,
          decoration: const BoxDecoration(
              color: ColorName.black1c, borderRadius: roundedTop_big),
          child: Column(
            children: [
              vspace24,
              SizedBox(
                width: 300,
                child: UBText(
                  align: TextAlign.center,
                  text: 'SAVE THIS CODE ON PAPER OR TAKE A SCREENSHOT',
                  size: 13,
                  color: ColorName.orange,
                ),
              ),
              vspace24,
              Obx(() {
                return Container(
                    alignment: Alignment.center,
                    width: 160.0,
                    decoration: const BoxDecoration(
                      color: ColorName.white,
                    ),
                    child: FadeInImage.memoryNetwork(
                      height: 160,
                      placeholder: kTransparentImage,
                      image: controller.qrImageAddress.value,
                      fadeInDuration: const Duration(
                        milliseconds: 300,
                      ),
                    )

                    // QrImage(
                    //   data: code,
                    //   version: QrVersions.auto,
                    //   size: 160.0,
                    // ),
                    );
              }),
              vspace12,
              UBWarningRow(
                text: LocaleKeys.twofaCopyCodeTopMessage.tr,
              ),
              vspace12,
              UBGreyContainer(
                color: ColorName.black,
                height: 40,
                width: Get.width - 24,
                child: Obx(
                  () {
                    final code = controller.characterCode.value;
                    return Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        UBText(text: code),
                        UBButton(
                          height: 36,
                          width: 40,
                          onClick: () {
                            thr.throttle(() {
                              Clipboard.setData(
                                ClipboardData(
                                  text: code,
                                ),
                              );
                              toastInfo('Code copied to clipboard');
                              controller.codeCoppied.value = true;
                            });
                          },
                          text: LocaleKeys.copy.tr,
                          textColor: ColorName.primaryBlue,
                          variant: ButtonVariant.Link,
                        )
                      ],
                    );
                  },
                ),
              ),
              const Spacer(),
              Padding(
                padding: const EdgeInsets.symmetric(
                  horizontal: 12,
                ),
                child: Obx(
                  () {
                    final codeCoppied = controller.codeCoppied.value;
                    return UBButton(
                        onClick: () {
                          if (codeCoppied) {
                            controller.step.value =
                                TwoFaSteps.Enable_Step3_EnterPassword_And_Code;
                          } else {
                            toastInfo('Please copy the code');
                          }
                        },
                        text: LocaleKeys.next.tr);
                  },
                ),
              ),
              Container(
                width: 120,
                padding:
                    const EdgeInsets.symmetric(horizontal: 12, vertical: 24),
                child: UBButton(
                  height: 24.0,
                  onClick: () {
                    Get.back();
                  },
                  variant: ButtonVariant.TransparentBackground,
                  textColor: ColorName.grey80,
                  text: LocaleKeys.cancel.tr,
                ),
              )
            ],
          ),
        ),
      ],
    );
  }
}
