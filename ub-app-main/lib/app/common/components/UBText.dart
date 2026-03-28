import 'package:flutter/material.dart';

import '../../../generated/colors.gen.dart';

class UBText extends StatelessWidget {
  const UBText({
    Key key,
    @required this.text,
    this.color = ColorName.greybf,
    this.size = 13.0,
    this.weight = FontWeight.w600,
    this.align = TextAlign.start,
    this.wrapped,
    this.lineHeight,
    this.fontStyle = FontStyle.normal,
  }) : super(key: key);

  final String text;
  final Color color;
  final bool wrapped;
  final double size;
  final FontWeight weight;
  final TextAlign align;
  final FontStyle fontStyle;

  /// The lineHeight of this text span, as a multiple of the font size.
  final double lineHeight;

  @override
  Widget build(BuildContext context) {
    final RichText txt = RichText(
      textAlign: align,
      text: TextSpan(
        text: text,
        style: TextStyle(
          height: lineHeight,
          color: color,
          fontSize: size,
          fontStyle: fontStyle,
          fontWeight: weight,
        ),
      ),
    );
    return wrapped == true ? Flexible(child: txt) : txt;
  }
}
