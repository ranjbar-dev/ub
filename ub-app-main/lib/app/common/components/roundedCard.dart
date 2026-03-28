import 'package:flutter/material.dart';
import '../../../generated/colors.gen.dart';

class RoundedCard extends StatelessWidget {
  final Widget child;

  const RoundedCard({Key key, @required this.child}) : super(key: key);
  @override
  Widget build(BuildContext context) {
    return Container(
      margin: const EdgeInsets.only(bottom: 8),
      width: double.infinity,
      decoration: BoxDecoration(
        color: ColorName.black2c,
        borderRadius: BorderRadius.circular(16),
      ),
      padding: const EdgeInsets.symmetric(
        horizontal: 12,
        vertical: 10,
      ),
      child: child,
    );
  }
}
