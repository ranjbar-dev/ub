import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:get/get.dart';

import '../../../generated/colors.gen.dart';
import '../../../utils/mixins/commonConsts.dart';
import 'UBBorderlessInput.dart';
import 'UBButton.dart';
import 'UBGreyContainer.dart';

class UBInputWithTitleAndPaste extends StatelessWidget {
  final bool withPaste;
  final String title;
  final String placeHolder;
  final TextEditingController controller;
  final Function(String) onChange;
  final double width;
  const UBInputWithTitleAndPaste({
    Key key,
    @required this.title,
    @required this.placeHolder,
    @required this.controller,
    this.withPaste = false,
    this.onChange,
    this.width,
  }) : super(key: key);
  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.only(bottom: 24.0),
      width: width ?? Get.width - 24,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          RichText(
            text: TextSpan(
              text: title,
              style: const TextStyle(
                color: ColorName.grey80,
                fontSize: 13.0,
                fontWeight: FontWeight.w600,
              ),
            ),
          ),
          vspace12,
          UBGreyContainer(
            color: ColorName.black,
            child: Stack(
              children: [
                UBBorderlessInput(
                  type: TextInputType.number,
                  placeholder: placeHolder,
                  onChange: onChange,
                  controller: controller,
                ),
                if (withPaste)
                  Positioned(
                    right: 0.0,
                    top: 3.0,
                    child: UBButton(
                      onClick: () async {
                        ClipboardData data =
                            await Clipboard.getData('text/plain');
                        handleTextFieldControllerChanged(
                          controller: controller,
                          newValue: data.text,
                        );
                        onChange(data.text);
                      },
                      text: 'Paste',
                      variant: ButtonVariant.Link,
                      textColor: ColorName.primaryBlue,
                      width: 60,
                      height: 30,
                      fontSize: 13.0,
                    ),
                  )
              ],
            ),
          )
        ],
      ),
    );
  }
}

handleTextFieldControllerChanged({
  @required TextEditingController controller,
  @required String newValue,
}) {
  controller.text = newValue;
  controller.selection = TextSelection.collapsed(offset: newValue.length);
}
