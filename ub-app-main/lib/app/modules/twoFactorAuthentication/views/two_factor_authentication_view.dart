import 'package:flutter/material.dart';

import 'package:get/get.dart';
import '../../../common/components/UBText.dart';
import '../../../common/components/appbarTextTitle.dart';
import '../../../common/components/pageContainer.dart';
import '../pages/Enable_step1_install_links.dart';
import '../pages/Enable_step2_copy_code.dart';
import '../pages/Enable_step3_enter_password_and_code.dart';
import '../pages/Enable_step4_fianl_status.dart';
import '../../../../generated/colors.gen.dart';
import '../../../../generated/locales.g.dart';

import '../controllers/two_factor_authentication_controller.dart';

class TwoFactorAuthenticationView
    extends GetView<TwoFactorAuthenticationController> {
  @override
  Widget build(BuildContext context) {
    return PageContainer(
      appbarTitle: AppBarTextTitle(
          title: (controller.isEnabled.value
                  ? LocaleKeys.disable.tr
                  : LocaleKeys.enable.tr) +
              ' ' +
              LocaleKeys.googleAuthenticator.tr),
      child: Obx(
        () {
          final step = controller.step.value;
          if (step == TwoFaSteps.Enable_Step1_Install_Links) {
            return InstallLinks();
          }
          if (step == TwoFaSteps.Enable_step2_copyCode) {
            return CopyCode();
          }
          if (step == TwoFaSteps.Enable_Step3_EnterPassword_And_Code) {
            return EnterPasswordAndCode();
          }
          if (step == TwoFaSteps.Enable_Step4_Final_Status) {
            return FinalStatus();
          }
          return Container(
            child: UBText(
              text: 'implement step',
              color: ColorName.red,
            ),
          );
        },
      ),
    );
  }
}
