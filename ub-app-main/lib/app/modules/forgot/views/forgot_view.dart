import 'package:flutter/material.dart';

import 'package:get/get.dart';
import '../../../common/components/UBButton.dart';
import '../../../common/components/UBScrollColumnExpandable.dart';
import '../../../common/components/UBSimpleInput.dart';
import '../../../common/components/appbarTextTitle.dart';
import '../../../common/components/bottomCard.dart';
import '../../../common/components/pageContainer.dart';
import '../controllers/forgot_controller.dart';
import '../../../../generated/assets.gen.dart';
import '../../../../generated/colors.gen.dart';
import '../../../../generated/locales.g.dart';
import '../../../../utils/mixins/commonConsts.dart';

class ForgotView extends GetView<ForgotController> {
  @override
  Widget build(BuildContext context) {
    controller.setStartTime();
    return PageContainer(
      protectedPage: false,
      appbarTitle: AppBarTextTitle(
        title: LocaleKeys.forgot.tr + " " + LocaleKeys.password.tr,
      ),
      child: UBScrollColumnExpandable(
        children: [
          fill,
          Container(
            height: 385,
            child: BottomCard(
              title: LocaleKeys.forgot.tr + " " + LocaleKeys.password.tr,
              children: [
                Expanded(
                  flex: 3,
                  child: Assets.images.forgotPasswordLogo.svg(),
                ),
                Expanded(
                  flex: 1,
                  child: const SizedBox(),
                ),
                Expanded(
                  flex: 2,
                  child: Obx(
                    () {
                      return UBSimpleInput(
                        error: controller.emailError.value,
                        onChange: (e) => {controller.handleEmailChange(e)},
                        placeHolder: LocaleKeys.email.tr,
                      );
                    },
                  ),
                ),
                Expanded(
                  flex: 2,
                  child: Container(
                    padding: const EdgeInsets.only(
                      bottom: 12,
                      top: 7,
                    ),
                    width: 90,
                    child: Obx(() {
                      return UBButton(
                        disabled: controller.email.value.length < 5,
                        isLodaing: controller.isLoading.value,
                        onClick: () {
                          controller.handleSubmitClick();
                          return;
                        },
                        text: LocaleKeys.submit.tr,
                      );
                    }),
                  ),
                ),
                Expanded(
                  flex: 1,
                  child: Container(
                    width: 90,
                    child: UBButton(
                      onClick: () => {
                        Get.back(),
                      },
                      text: LocaleKeys.cancel.tr,
                      textColor: ColorName.grey80,
                      variant: ButtonVariant.TransparentBackground,
                    ),
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
