import 'dart:math' as math;

import 'package:flutter/material.dart';

import '../../../../../generated/colors.gen.dart';

/// Strategy used when drawing individual data points over the sparkline.
enum PointsMode {
  /// Do not draw individual points.
  none,

  /// Draw all the points in the data set.
  all,

  /// Draw only the last point in the data set.
  last,
}

/// A widget that draws a sparkline chart.
///
/// Represents the given [data] in a sparkline chart that spans the available
/// space.
///
/// By default only the sparkline is drawn, with its looks defined by
/// the [lineWidth], [lineColor], and [lineGradient] properties.
///
///
/// The area above or below the sparkline can be filled with the provided
/// [fillColor] or [fillGradient] by setting the desired [fillMode].
///
/// [pointsMode] controls how individual points are drawn over the sparkline
/// at the provided data point. Their appearance is determined by the
/// [pointSize] and [pointColor] properties.
///
/// By default, the sparkline is sized to fit its container. If the
/// sparkline is in an unbounded space, it will size itself according to the
/// given [fallbackWidth] and [fallbackHeight].
class Sparkline extends StatelessWidget {
  /// Creates a widget that represents provided [data] in a Sparkline chart.
  Sparkline({
    Key key,
    @required this.data,
    this.lineWidth = 2.0,
    this.lineColor = ColorName.primaryBlue,
    this.lineGradient,
    this.pointSize = 4.0,
    this.pointColor = const Color(0xFF0277BD), //Colors.lightBlue[800]
    this.fillColor = const Color(0xFF81D4FA), //Colors.lightBlue[200]
    this.fillGradient,
    this.fallbackHeight = 100.0,
    this.fallbackWidth = 300.0,
  })  : assert(data != null),
        assert(lineWidth != null),
        assert(lineColor != null),
        assert(pointSize != null),
        assert(pointColor != null),
        assert(fillColor != null),
        assert(fallbackHeight != null),
        assert(fallbackWidth != null),
        super(key: key);

  /// List of values to be represented by the sparkline.
  ///
  /// Each data entry represents an individual point on the chart, with a path
  /// drawn connecting the consecutive points to form the sparkline.
  ///
  /// The values are normalized to fit within the bounds of the chart.
  final List<double> data;

  /// The width of the sparkline.
  ///
  /// Defaults to 2.0.
  final double lineWidth;

  /// The color of the sparkline.
  ///
  /// Defaults to Colors.lightBlue.
  ///
  /// This is ignored if [lineGradient] is non-null.
  final Color lineColor;

  /// A gradient to use when coloring the sparkline.
  ///
  /// If this is specified, [lineColor] has no effect.
  final Gradient lineGradient;

  /// Determines how individual data points should be drawn over the sparkline.
  ///

  /// The size to use when drawing individual data points over the sparkline.
  ///
  /// Defaults to 4.0.
  final double pointSize;

  /// The color used when drawing individual data points over the sparkline.
  ///
  /// Defaults to Colors.lightBlue[800].
  final Color pointColor;

  /// Determines if the sparkline path should have sharp corners where two
  /// segments intersect.
  ///
  /// Defaults to false.

  /// Determines the area that should be filled with [fillColor].
  ///

  /// The fill color used in the chart, as determined by [fillMode].
  ///
  /// Defaults to Colors.lightBlue[200].
  ///
  /// This is ignored if [fillGradient] is non-null.
  final Color fillColor;

  /// A gradient to use when filling the chart, as determined by [fillMode].
  ///
  /// If this is specified, [fillColor] has no effect.
  final Gradient fillGradient;

  /// The width to use when the sparkline is in a situation with an unbounded
  /// width.
  ///
  /// See also:
  ///
  ///  * [fallbackHeight], the same but vertically.
  final double fallbackWidth;

  /// The height to use when the sparkline is in a situation with an unbounded
  /// height.
  ///
  /// See also:
  ///
  ///  * [fallbackWidth], the same but horizontally.
  final double fallbackHeight;

  @override
  Widget build(BuildContext context) {
    return new LimitedBox(
      maxWidth: fallbackWidth,
      maxHeight: fallbackHeight,
      child: new CustomPaint(
        size: Size.infinite,
        painter: new _SparklinePainter(
          data,
          lineWidth: lineWidth,
          lineColor: lineColor,
          lineGradient: lineGradient,
          fillColor: fillColor,
          fillGradient: fillGradient,
          pointSize: pointSize,
          pointColor: pointColor,
        ),
      ),
    );
  }
}

class _SparklinePainter extends CustomPainter {
  _SparklinePainter(
    this.dataPoints, {
    @required this.lineWidth,
    @required this.lineColor,
    this.lineGradient,
    @required this.fillColor,
    this.fillGradient,
    @required this.pointSize,
    @required this.pointColor,
  })  : _max = dataPoints.reduce(math.max),
        _min = dataPoints.reduce(math.min);

  final List<double> dataPoints;

  final double lineWidth;
  final Color lineColor;
  final Gradient lineGradient;

  final Color fillColor;
  final Gradient fillGradient;

  final double pointSize;
  final Color pointColor;

  final double _max;
  final double _min;

  @override
  void paint(Canvas canvas, Size size) {
    final Path path = new Path();

    Offset startPoint;

    final double width = size.width - lineWidth;
    final double height = size.height - lineWidth;

    final yMin = _min;
    final yMax = _max;
    final yHeight = yMax - yMin;
    final xAxisStep = size.width / dataPoints.length;
    var xValue = 0.0;
    for (var i = 0; i < dataPoints.length; i++) {
      final value = dataPoints[i];
      final yValue = yHeight == 0
          ? (0.5 * size.height)
          : ((yMax - value) / yHeight) * (size.height);
      if (xValue == 0) {
        startPoint = new Offset(xValue, yValue);
        path.moveTo(xValue, yValue);
      } else {
        final previousValue = dataPoints[i - 1];
        final xPrevious = xValue - xAxisStep;
        final yPrevious = yHeight == 0
            ? (0.5 * size.height)
            : ((yMax - previousValue) / yHeight) * size.height;
        final controlPointX = xPrevious + (xValue - xPrevious) / 2;
        // HERE is the main line of code making your line smooth
        path.cubicTo(
            controlPointX, yPrevious, controlPointX, yValue, xValue, yValue);
      }
      xValue += xAxisStep;
    }
    Paint paint = new Paint()
      ..strokeWidth = lineWidth
      ..color = lineColor
      ..strokeCap = StrokeCap.round
      ..strokeJoin = StrokeJoin.round
      ..style = PaintingStyle.stroke;

    if (lineGradient != null) {
      final Rect lineRect = new Rect.fromLTWH(0.0, 0.0, width, height);
      paint.shader = lineGradient.createShader(lineRect);
    }

    fillUnder(
      path: path,
      size: size,
      startPoint: startPoint,
      width: width,
      height: height,
      canvas: canvas,
    );

    canvas.drawPath(path, paint);
  }

  fillUnder({
    @required Path path,
    @required Size size,
    @required Offset startPoint,
    @required double width,
    @required double height,
    @required Canvas canvas,
  }) {
    Path fillPath = new Path()..addPath(path, Offset.zero);

    fillPath.relativeLineTo(lineWidth / 2, 0.0);
    fillPath.lineTo(size.width, size.height);
    fillPath.lineTo(0.0, size.height);
    fillPath.lineTo(startPoint.dx - lineWidth / 2, startPoint.dy);

    fillPath.close();

    Paint fillPaint = new Paint()
      ..strokeWidth = 0.0
      ..color = fillColor
      ..style = PaintingStyle.fill;

    if (fillGradient != null) {
      final Rect fillRect = new Rect.fromLTWH(0.0, 0.0, width, height);
      fillPaint.shader = fillGradient.createShader(fillRect);
    }
    canvas.drawPath(fillPath, fillPaint);
  }

  @override
  bool shouldRepaint(_SparklinePainter old) {
    return dataPoints != old.dataPoints ||
        lineWidth != old.lineWidth ||
        lineColor != old.lineColor ||
        lineGradient != old.lineGradient ||
        fillColor != old.fillColor ||
        fillGradient != old.fillGradient ||
        pointSize != old.pointSize ||
        pointColor != old.pointColor;
  }

  void drawBezierLine(Canvas canvas, Offset start, Offset end) {
    Offset middle = end / 2 + start / 2;
    Offset baselineVector = end - start;
    Offset controlPoint =
        middle + Offset(-baselineVector.dy, baselineVector.dx);
    Paint paint = Paint()
      ..color = Colors.red
      ..style = PaintingStyle.stroke
      ..strokeWidth = 8.0;

    var path = Path();
    path.moveTo(start.dx, start.dy);
    path.quadraticBezierTo(controlPoint.dx, controlPoint.dy, end.dx, end.dy);
    canvas.drawPath(path, paint);
  }
}
