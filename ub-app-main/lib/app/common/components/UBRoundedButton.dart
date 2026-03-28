import 'package:flutter/material.dart';

import '../../../generated/colors.gen.dart';

class UBRoundButton extends StatelessWidget {
  final Widget child;
  final double size;
  final Color color;
  final EdgeInsets padding;
  final Function onClick;

  const UBRoundButton(
      {Key key,
      this.size = 32,
      this.color,
      this.padding = const EdgeInsets.all(4),
      @required this.onClick,
      @required this.child})
      : super(key: key);
  @override
  Widget build(BuildContext context) {
    return ClipOval(
      child: GestureDetector(
        onTap: onClick,
        child: Container(
          width: size,
          height: size,
          padding: padding,
          color: color ??
              ColorName.black.withOpacity(
                0.4,
              ),
          child: child,
        ),
      ),
    );
  }
}
