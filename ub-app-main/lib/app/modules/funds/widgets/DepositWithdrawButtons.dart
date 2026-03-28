import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../common/components/UBButton.dart';
import '../../../routes/app_pages.dart';
import '../../../../generated/colors.gen.dart';
import '../../../../generated/locales.g.dart';
import '../../../../utils/mixins/commonConsts.dart';
import '../../../../utils/mixins/toast.dart';
import '../../../../utils/throttle.dart';

final thr = new Throttling(duration: const Duration(milliseconds: 4000));

class DepositWithdrawButtons extends StatelessWidget with Toaster {
  final bool isUserVerified;

  const DepositWithdrawButtons({Key key, @required this.isUserVerified})
      : super(key: key);
  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 12),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          UBButton(
            height: 32.0,
            width: (Get.width / 2) - 24.0,
            onClick: () {
              if (isUserVerified == true) {
                Get.toNamed(AppPages.DEPOSITS);
              } else {
                toastToVerifyEmail();
              }
            },
            text: LocaleKeys.deposits.tr,
          ),
          hspace12,
          UBButton(
            height: 32.0,
            width: (Get.width / 2) - 24.0,
            buttonColor: ColorName.black,
            borderColor: ColorName.primaryBlue,
            variant: ButtonVariant.Outline,
            onClick: () {
              if (isUserVerified == true) {
                Get.toNamed(AppPages.WITHDRAWALS);
              } else {
                toastToVerifyEmail();
              }
            },
            text: LocaleKeys.withdrawals.tr,
          ),
        ],
      ),
    );
  }

  toastToVerifyEmail() {
    thr.throttle(() {
      toastWarning(
          "we've sent you a verification email, please verify your email first");
    });
  }
}
