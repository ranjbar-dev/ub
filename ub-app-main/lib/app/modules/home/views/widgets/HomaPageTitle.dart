import 'package:flutter/material.dart';
import '../../../../common/components/UBText.dart';
import '../../../../../generated/colors.gen.dart';

class HomePageTitle extends StatelessWidget {
  final String text;
  const HomePageTitle({Key key, this.text}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Padding(
        padding: const EdgeInsets.symmetric(
          vertical: 12,
        ),
        child: UBText(
          text: text,
          size: 15.0,
          weight: FontWeight.w700,
          color: ColorName.white,
        ));
  }
}
