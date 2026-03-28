import 'package:flutter/widgets.dart';

import '../../../generated/colors.gen.dart';
import '../../../utils/extentions/basic.dart';

class UBCountup extends StatefulWidget {
  final double begin;
  final double end;
  final int precision;
  final Curve curve;
  final Duration duration;
  final String prefix;
  final String suffix;
  final Color color;

  UBCountup({
    Key key,
    @required this.begin,
    @required this.end,
    this.precision = 0,
    this.curve = Curves.linear,
    this.duration = const Duration(milliseconds: 250),
    this.prefix = '',
    this.suffix = '',
    this.color = ColorName.greybf,
  }) : super(key: key);

  @override
  _CountupState createState() => _CountupState();
}

class _CountupState extends State<UBCountup> with TickerProviderStateMixin {
  AnimationController _controller;
  Animation<double> _animation;
  double _latestBegin;
  double _latestEnd;

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  void initState() {
    super.initState();
    _controller = AnimationController(duration: widget.duration, vsync: this);
    _latestBegin = widget.begin;
    _latestEnd = widget.end;
  }

  @override
  Widget build(BuildContext context) {
    CurvedAnimation curvedAnimation =
        CurvedAnimation(parent: _controller, curve: widget.curve);
    _animation = Tween<double>(begin: widget.begin, end: widget.end)
        .animate(curvedAnimation);

    if (widget.begin != _latestBegin || widget.end != _latestEnd) {
      _controller.reset();
    }

    _latestBegin = widget.begin;
    _latestEnd = widget.end;
    _controller.forward();

    return _CountupAnimatedText(
      key: UniqueKey(),
      animation: _animation,
      precision: widget.precision,
      color: widget.color,
      prefix: widget.prefix,
      suffix: widget.suffix,
    );
  }
}

class _CountupAnimatedText extends AnimatedWidget {
  final RegExp reg = new RegExp(r'(\d{1,3})(?=(\d{3})+(?!\d))');
  final Animation<double> animation;
  final int precision;
  final int maxLines;
  final String prefix;
  final String suffix;
  final Color color;

  _CountupAnimatedText({
    Key key,
    @required this.animation,
    @required this.precision,
    this.maxLines,
    this.prefix,
    this.suffix,
    this.color = ColorName.greybf,
  }) : super(key: key, listenable: animation);

  @override
  Widget build(BuildContext context) => RichText(
        text: TextSpan(
          text: '$prefix' +
              (this
                  .animation
                  .value
                  .toFixedWithoutRounding(precision)
                  .currencyFormat()) +
              '$suffix',
          style: TextStyle(
            fontSize: 14.0,
            fontWeight: FontWeight.w600,
            color: color,
          ),
        ),
      );
}
