import 'package:flutter/material.dart';
import 'package:flutter_staggered_animations/flutter_staggered_animations.dart';

class UBColumnScaleAnimator extends StatelessWidget {
  final List<Widget> children;
  final MainAxisSize mainAxisSize;
  const UBColumnScaleAnimator(
      {Key key, this.children, this.mainAxisSize = MainAxisSize.min})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return AnimationLimiter(
      child: Column(
        mainAxisSize: mainAxisSize,
        children: AnimationConfiguration.toStaggeredList(
          duration: const Duration(milliseconds: 375),
          childAnimationBuilder: (widget) => ScaleAnimation(
            scale: 0.2,
            curve: Curves.easeOutCubic,
            child: FadeInAnimation(
              child: widget,
            ),
          ),
          children: children,
        ),
      ),
    );
  }
}

class UBColumnSlideAnimator extends StatelessWidget {
  final List<Widget> children;
  final MainAxisSize mainAxisSize;
  final CrossAxisAlignment crossAxisAlignment;
  final MainAxisAlignment mainAxisAlignment;
  final int animationTime;
  const UBColumnSlideAnimator({
    Key key,
    this.children,
    this.mainAxisSize = MainAxisSize.min,
    this.crossAxisAlignment = CrossAxisAlignment.center,
    this.mainAxisAlignment = MainAxisAlignment.start,
    this.animationTime = 300,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return AnimationLimiter(
      child: Column(
        mainAxisSize: mainAxisSize,
        mainAxisAlignment: mainAxisAlignment,
        crossAxisAlignment: crossAxisAlignment,
        children: AnimationConfiguration.toStaggeredList(
          duration: Duration(milliseconds: animationTime),
          childAnimationBuilder: (widget) => SlideAnimation(
            horizontalOffset: 70.0,
            curve: Curves.easeOut,
            child: FadeInAnimation(
              child: widget,
            ),
          ),
          children: children,
        ),
      ),
    );
  }
}

class UBColumnVerticalSlideAnimator extends StatelessWidget {
  final List<Widget> children;
  final MainAxisSize mainAxisSize;
  final CrossAxisAlignment crossAxisAlignment;
  final MainAxisAlignment mainAxisAlignment;
  final int animationTime;
  const UBColumnVerticalSlideAnimator({
    Key key,
    this.children,
    this.mainAxisSize = MainAxisSize.min,
    this.crossAxisAlignment = CrossAxisAlignment.center,
    this.mainAxisAlignment = MainAxisAlignment.start,
    this.animationTime = 300,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return AnimationLimiter(
      child: Column(
        mainAxisSize: mainAxisSize,
        mainAxisAlignment: mainAxisAlignment,
        crossAxisAlignment: crossAxisAlignment,
        children: AnimationConfiguration.toStaggeredList(
          duration: Duration(milliseconds: animationTime),
          childAnimationBuilder: (widget) => SlideAnimation(
            verticalOffset: 30.0,
            curve: Curves.easeOut,
            child: FadeInAnimation(
              child: widget,
            ),
          ),
          children: children,
        ),
      ),
    );
  }
}
