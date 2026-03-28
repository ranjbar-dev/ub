import 'dart:ui' as ui;

import 'package:flutter/material.dart';
import '../../../generated/colors.gen.dart';

// class UBCustomPainter extends CustomPainter {
//   @override
//   void paint(Canvas canvas, Size size) {
//     Paint paint = Paint();
//     Path path = Path();

//     // Path number 1

//     paint.color = /*Color(0xffffffff).withOpacity(0)*/ Colors.cyanAccent;
//     path = Path();
//     path.lineTo(size.width, size.height * 1.11);
//     path.cubicTo(size.width, size.height * 1.11, size.width * 0.76,
//         size.height * 1.11, size.width * 0.76, size.height * 1.11);
//     path.cubicTo(size.width * 0.74, size.height * 1.11, size.width * 0.72,
//         size.height * 1.02, size.width * 0.7, size.height * 0.86);
//     path.cubicTo(size.width * 0.7, size.height * 0.86, size.width * 0.66,
//         size.height * 0.52, size.width * 0.66, size.height * 0.52);
//     path.cubicTo(size.width * 0.64, size.height * 0.26, size.width * 0.6,
//         size.height * 0.11, size.width * 0.56, size.height * 0.11);
//     path.cubicTo(size.width * 0.56, size.height * 0.11, size.width * 0.51,
//         size.height * 0.11, size.width * 0.51, size.height * 0.11);
//     path.cubicTo(size.width * 0.51, size.height * 0.11, size.width * 0.46,
//         size.height * 0.11, size.width * 0.46, size.height * 0.11);
//     path.cubicTo(size.width * 0.42, size.height * 0.11, size.width * 0.39,
//         size.height / 4, size.width * 0.36, size.height * 0.48);
//     path.cubicTo(size.width * 0.36, size.height * 0.48, size.width * 0.32,
//         size.height * 0.88, size.width * 0.32, size.height * 0.88);
//     path.cubicTo(size.width * 0.3, size.height * 1.03, size.width * 0.28,
//         size.height * 1.11, size.width * 0.26, size.height * 1.11);
//     path.cubicTo(size.width * 0.26, size.height * 1.11, 0, size.height * 1.11,
//         0, size.height * 1.11);
//     canvas.drawPath(path, paint);
//   }

//   @override
//   bool shouldRepaint(CustomPainter oldDelegate) {
//     return true;
//   }
// }

//Don't panic dude. This is generated from https://fluttershapemaker.com/
class UBCustomPainter extends CustomPainter {
  @override
  void paint(Canvas canvas, Size size) {
    Path path_0 = Path();
    path_0.moveTo(size.width, size.height * 0.9375000);
    path_0.lineTo(size.width * 0.7612472, size.height * 0.9375000);
    path_0.cubicTo(
        size.width * 0.7377087,
        size.height * 0.9375000,
        size.width * 0.7151843,
        size.height * 0.8614625,
        size.width * 0.6988575,
        size.height * 0.7268813);
    path_0.lineTo(size.width * 0.6638598, size.height * 0.4384031);
    path_0.cubicTo(
        size.width * 0.6371425,
        size.height * 0.2181794,
        size.width * 0.6002843,
        size.height * 0.09375000,
        size.width * 0.5617669,
        size.height * 0.09375000);
    path_0.lineTo(size.width * 0.5078740, size.height * 0.09375000);
    path_0.lineTo(size.width * 0.4587787, size.height * 0.09375000);
    path_0.cubicTo(
        size.width * 0.4220811,
        size.height * 0.09375000,
        size.width * 0.3868134,
        size.height * 0.2067319,
        size.width * 0.3603969,
        size.height * 0.4089250);
    path_0.lineTo(size.width * 0.3165031, size.height * 0.7448938);
    path_0.cubicTo(
        size.width * 0.3003591,
        size.height * 0.8684563,
        size.width * 0.2788071,
        size.height * 0.9375000,
        size.width * 0.2563811,
        size.height * 0.9375000);
    path_0.lineTo(0, size.height * 0.9375000);

    Paint paint0Stroke = Paint()
      ..style = PaintingStyle.stroke
      ..strokeWidth = size.width * 0.01881102;
    paint0Stroke.shader = ui.Gradient.radial(Offset(0, 0), size.width * .0, [
      Color(0xf49806).withOpacity(0),
      Color(0xf49806).withOpacity(1),
      Color(0xf49806).withOpacity(0)
    ], [
      0,
      0.536458,
      1
    ]);
    canvas.drawPath(path_0, paint0Stroke);

    Paint paint0Fill = Paint()..style = PaintingStyle.fill;
    paint0Fill.color = ColorName.black2c.withOpacity(1.0);
    canvas.drawPath(path_0, paint0Fill);
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) {
    return true;
  }
}
