import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../generated/colors.gen.dart';
import '../../../utils/mixins/commonConsts.dart';
import 'UBText.dart';

class ConnectionLost extends StatelessWidget {
  final String text;

  const ConnectionLost({
    Key key,
    this.text = 'Connection to internet is lost!',
  }) : super(key: key);
  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Expanded(
          child: Container(
            color: ColorName.black.withOpacity(0.2),
            child: Center(
              child: Container(
                decoration: const BoxDecoration(
                  borderRadius: rounded7,
                  color: ColorName.black2c,
                ),
                width: Get.width / 2,
                height: 30,
                child: Align(
                  child: UBText(
                    color: ColorName.red,
                    text: text,
                    weight: FontWeight.w600,
                  ),
                ),
              ),
            ),
          ),
        ),
      ],
    );
  }
}
