import 'package:flutter/material.dart';

import 'package:get/get.dart';
import '../../../common/components/UBText.dart';
import '../../../common/components/appbarTextTitle.dart';
import '../../../common/components/pageContainer.dart';
import '../pages/EnterPhoneNumber.dart';
import '../pages/enterPassword.dart';
import '../pages/enter_s_m_s_verification_code.dart';
import '../../../../generated/colors.gen.dart';
import '../../../../generated/locales.g.dart';

import '../controllers/phone_verification_controller.dart';

class PhoneVerificationView extends GetView<PhoneVerificationController> {
  @override
  Widget build(BuildContext context) {
    final recivingData = Get.arguments;
    controller.reset(recivingData: recivingData);
    return PageContainer(
      appbarTitle: AppBarTextTitle(title: LocaleKeys.phoneVerifivation.tr),
      child: Obx(
        () {
          final step = controller.step.value;

          if (step == PhoneVerificationSteps.EnterPhoneNumber) {
            return EnterPhoneNumber();
          }
          if (step == PhoneVerificationSteps.EnterSMSVerificationCode) {
            return EnterSMSVerificationCode();
          }
          if (step == PhoneVerificationSteps.EnterPassword) {
            return EnterPassword();
          }
          return Container(
            child: UBText(
              text: 'Implement Page',
              color: ColorName.red,
            ),
          );
        },
      ),
    );
  }
}
