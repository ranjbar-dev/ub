import 'package:flutter/material.dart';

class UBRectangle extends StatelessWidget {
  final double width;
  final double height;
  final Color color;

  const UBRectangle({Key key, this.width, this.height, this.color})
      : super(key: key);
  @override
  Widget build(BuildContext context) {
    return CustomPaint(
      size: Size(
        width,
        height,
      ),
      painter: RectPainter(height: 18.0, color: color, width: width),
    );
  }
}

class RectPainter extends CustomPainter {
  final Color color;
  final double height;
  final double width;
  RectPainter({this.width, this.height, this.color});

  @override
  void paint(Canvas canvas, Size size) {
    Offset start = Offset(0, height / 2);
    Offset end = Offset(width, height / 2);

    canvas.drawLine(
        start,
        end,
        Paint()
          ..color = color
          ..strokeWidth = height);
  }

  @override
  bool shouldRepaint(covariant CustomPainter oldDelegate) {
    return true;
  }
}
