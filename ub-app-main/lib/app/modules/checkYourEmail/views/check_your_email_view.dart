import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../../generated/assets.gen.dart';
import '../../../../generated/colors.gen.dart';
import '../../../../utils/mixins/commonConsts.dart';
import '../../../common/components/UBText.dart';
import '../../../common/components/bottomCard.dart';
import '../../../common/components/pageContainer.dart';
import '../controllers/check_your_email_controller.dart';

class CheckYourEmailBottom {
  final String text;
  final Function onClick;

  CheckYourEmailBottom({this.text, this.onClick});
}

class CheckYourEmailView extends GetView<CheckYourEmailController> {
  final String title;
  final String sub;
  final String warningText;
  final CheckYourEmailBottom bottom;

  CheckYourEmailView({this.title, this.sub, this.warningText, this.bottom});
  @override
  Widget build(BuildContext context) {
    return PageContainer(
      child: Column(
        children: [
          fill,
          Container(
            height: 460,
            child: BottomCard(
              withCloseButton: true,
              title: title,
              children: [
                Assets.images.emailPageMainImage.svg(),
                if (sub != null) vspace24,
                if (sub != null)
                  UBText(
                    text: sub,
                    size: 14.0,
                    color: ColorName.greybf,
                  ),
                if (warningText != null) vspace24,
                if (warningText != null)
                  Container(
                    width: Get.width - 24.0,
                    child: UBText(
                      align: TextAlign.center,
                      text: warningText,
                      size: 14.0,
                      color: ColorName.orange,
                    ),
                  ),
              ],
            ),
          )
        ],
      ),
    );
  }
}
