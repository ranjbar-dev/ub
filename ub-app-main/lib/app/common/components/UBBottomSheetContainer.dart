import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../generated/colors.gen.dart';
import 'UBText.dart';

class UBButtomSheetContainer extends StatelessWidget {
  const UBButtomSheetContainer({
    Key key,
    this.height = 110.0,
    this.title,
    @required this.child,
    this.backgroundColor = ColorName.grey16,
  }) : super(key: key);
  final double height;
  final Color backgroundColor;
  final String title;
  final Widget child;
  @override
  Widget build(BuildContext context) {
    return Container(
      height: height,
      margin: const EdgeInsets.symmetric(horizontal: 12),
      decoration: BoxDecoration(
        color: backgroundColor,
        borderRadius: const BorderRadius.only(
          topLeft: const Radius.circular(8),
          topRight: const Radius.circular(8),
        ),
      ),
      child: Column(
        children: [
          Padding(
            padding: const EdgeInsets.all(12),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                if (title != null)
                  UBText(text: title, size: 13)
                else
                  SizedBox(),
                GestureDetector(
                  onTap: () {
                    Get.back();
                  },
                  child: Icon(
                    Icons.close,
                    color: ColorName.greybf,
                  ),
                ),
              ],
            ),
          ),
          child
        ],
      ),
    );
  }
}
