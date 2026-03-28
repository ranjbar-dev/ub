import 'package:flutter/material.dart';
import '../../../generated/colors.gen.dart';

class UBSection extends StatelessWidget {
  final String title;
  final Widget child;
  final double hTitlePadding;
  final Widget titleEndWidget;

  const UBSection({
    Key key,
    @required this.title,
    @required this.child,
    this.hTitlePadding = 8.0,
    this.titleEndWidget = const SizedBox(),
  }) : super(key: key);
  @override
  Widget build(BuildContext context) {
    return Column(crossAxisAlignment: CrossAxisAlignment.start, children: [
      Container(
        alignment: Alignment.centerLeft,
        padding: EdgeInsets.symmetric(horizontal: hTitlePadding),
        height: 20,
        child: Row(
          children: [
            RichText(
              text: TextSpan(
                text: title,
                style: const TextStyle(
                  fontSize: 13,
                  fontWeight: FontWeight.w600,
                  color: ColorName.greybf,
                ),
              ),
            ),
            const Spacer(),
            titleEndWidget
          ],
        ),
      ),
      child
    ]);
  }
}
