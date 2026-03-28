import 'package:flutter/material.dart';
import 'package:get/get.dart';

class UBColoredOverlay extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return Container(
      width: Get.width,
      height: Get.height,
      decoration: new BoxDecoration(
        gradient: new LinearGradient(
            colors: [
              Colors.transparent,
              const Color.fromRGBO(54, 255, 231, 0.03),
              const Color.fromRGBO(54, 255, 231, 0.05),
              const Color.fromRGBO(54, 255, 231, 0.08),
              const Color.fromRGBO(54, 255, 231, 0.1),
              const Color.fromRGBO(84, 125, 240, 0.14),
              const Color.fromRGBO(112, 0, 255, 0.14),
              const Color.fromRGBO(112, 0, 255, 0.1),
              Colors.transparent,
            ],
            begin: const FractionalOffset(0.0, 0.0),
            end: const FractionalOffset(0.0, 1.0),
            stops: [
              0.0,
              0.05,
              0.1,
              0.125,
              .25,
              .5,
              .75,
              .85,
              1,
            ],
            tileMode: TileMode.clamp),
      ),
    );
  }
}
