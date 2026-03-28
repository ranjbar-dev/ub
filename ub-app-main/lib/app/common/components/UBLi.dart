import 'package:flutter/material.dart';

import '../../../generated/colors.gen.dart';
import 'UBText.dart';

class UBLi extends StatelessWidget {
  final String text;
  final Color stringColor;
  final Color dotColor;
  const UBLi(
      {Key key,
      @required this.text,
      this.stringColor = ColorName.white,
      this.dotColor = ColorName.white})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.only(left: 12.0),
      margin: const EdgeInsets.only(top: 8.0),
      child: Row(
        children: [
          Container(
            margin: const EdgeInsets.only(right: 4.0),
            width: 5,
            height: 5,
            decoration: BoxDecoration(
              color: dotColor,
              borderRadius: BorderRadius.circular(12),
            ),
          ),
          UBText(
            wrapped: true,
            text: text,
            color: stringColor,
            size: 11,
          ),
        ],
      ),
    );
  }
}
