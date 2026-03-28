import 'package:carousel_slider/carousel_slider.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../generated/colors.gen.dart';

class UBCarousel extends StatefulWidget {
  final CarouselController controller;
  final List<Widget> slides;
  final double height;
  final double aspectRatio;
  final double fraction;
  final bool infiniteScroll;
  final Function(int index) onChange;
  final int selectedIndex;
  final Axis scrollDirection;
  final bool showNavigationArrows;
  final bool disabledrag;
  final bool autoPlay;
  final int autoplayInterval;
  final bool showIndicators;

  const UBCarousel({
    Key key,
    @required this.controller,
    @required this.slides,
    @required this.height,
    this.aspectRatio = 1.0,
    this.infiniteScroll = false,
    @required this.onChange,
    this.selectedIndex,
    this.scrollDirection = Axis.horizontal,
    this.showNavigationArrows = true,
    this.autoPlay = false,
    this.autoplayInterval = 3000,
    this.disabledrag = false,
    this.showIndicators = false,
    this.fraction = 0.9,
  }) : super(key: key);

  @override
  _UBCarouselState createState() => _UBCarouselState();
}

class _UBCarouselState extends State<UBCarousel> {
  int _current = 0;
  @override
  Widget build(BuildContext context) {
    return Stack(
      clipBehavior: Clip.none,
      children: [
        CarouselSlider(
            carouselController: widget.controller,
            items: widget.slides,
            options: CarouselOptions(
              scrollPhysics: widget.disabledrag
                  ? const NeverScrollableScrollPhysics()
                  : null,
              height: widget.height,
              aspectRatio: widget.aspectRatio,
              viewportFraction: widget.fraction,
              initialPage: 0,
              enableInfiniteScroll: widget.infiniteScroll,
              autoPlay: widget.autoPlay,
              autoPlayInterval: Duration(milliseconds: widget.autoplayInterval),
              autoPlayAnimationDuration: Duration(milliseconds: 800),
              enlargeCenterPage: true,
              onPageChanged: (idx, d) {
                setState(() {
                  _current = idx;
                });
                widget.onChange(idx);
                return;
              },
              scrollDirection: widget.scrollDirection,
            )),
        if (widget.slides.length > 1 && widget.showNavigationArrows)
          Positioned(
            top: ((widget.height) / 2) - 25,
            left: 12.0,
            child: GestureDetector(
              onTap: () {
                widget.controller.previousPage();
              },
              child: Container(
                height: 50,
                width: 50,
                color: Colors.transparent,
                child: Center(
                  child: const Icon(
                    Icons.keyboard_arrow_left_rounded,
                    color: ColorName.white,
                  ),
                ),
              ),
            ),
          ),
        if (widget.slides.length > 1 && widget.showNavigationArrows)
          Positioned(
            top: ((widget.height) / 2) - 25,
            right: 12.0,
            child: GestureDetector(
              onTap: () {
                widget.controller.nextPage();
              },
              child: Container(
                height: 50,
                width: 50,
                color: Colors.transparent,
                child: Center(
                  child: const Icon(
                    Icons.keyboard_arrow_right_rounded,
                    color: ColorName.white,
                  ),
                ),
              ),
            ),
          ),
        if (widget.showIndicators)
          Positioned(
            bottom: -24.0,
            left: (Get.width / 2) - 66,
            child: Container(
              width: 100,
              height: 24,
              child: Center(
                child: Row(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: widget.slides.asMap().entries.map((entry) {
                    return Container(
                      width: 12.0,
                      height: 12.0,
                      margin:
                          EdgeInsets.symmetric(vertical: 8.0, horizontal: 4.0),
                      decoration: BoxDecoration(
                          shape: BoxShape.circle,
                          color: _current == entry.key
                              ? ColorName.primaryBlue
                              : Color.fromRGBO(82, 82, 97, 1)),
                    );
                  }).toList(),
                ),
              ),
            ),
          ),
      ],
    );
  }
}
