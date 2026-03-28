import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../common/components/UBText.dart';
import '../../../common/components/bottomCard.dart';
import '../../../common/components/pageContainer.dart';
import '../../../../generated/colors.gen.dart';
import '../../../../utils/mixins/commonConsts.dart';

import '../controllers/separate_message_page_controller.dart';

class SeparateMessagePageText {
  final String text;
  final Color color;
  SeparateMessagePageText({this.text, this.color = ColorName.greybf});
}

class SeparateMessagePageView extends GetView<SeparateMessagePageController> {
  final dynamic image;
  final List<SeparateMessagePageText> texts;
  final double contentHeight;
  SeparateMessagePageView({
    @required this.image,
    this.texts,
    this.contentHeight = 460.0,
  });

  @override
  Widget build(BuildContext context) {
    return PageContainer(
      child: Column(
        children: [
          fill,
          Container(
            height: contentHeight,
            child: BottomCard(
              withCloseButton: true,
              title: '',
              children: [
                image,
                vspace24,
                for (var item in texts)
                  Container(
                    margin: const EdgeInsets.only(bottom: 12.0),
                    width: Get.width - 24.0,
                    child: UBText(
                      align: TextAlign.center,
                      text: item.text,
                      size: 14.0,
                      color: item.color,
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
