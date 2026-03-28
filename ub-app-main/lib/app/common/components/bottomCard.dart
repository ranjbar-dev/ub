import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../utils/mixins/commonConsts.dart';

import '../../../generated/assets.gen.dart';
import '../../../generated/colors.gen.dart';

class BottomCard extends StatelessWidget {
  final List<Widget> children;
  final String title;
  final bool withCloseButton;
  const BottomCard({
    Key key,
    @required this.children,
    @required this.title,
    this.withCloseButton = false,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        Container(
          width: double.infinity,
          padding: const EdgeInsets.symmetric(
            horizontal: 12.0,
            vertical: 24,
          ),
          decoration: const BoxDecoration(
            color: ColorName.black2c,
            borderRadius: const BorderRadius.only(
              topLeft: const Radius.circular(24),
              topRight: const Radius.circular(24),
            ),
          ),
          child: Column(
            children: [
              vspace24,
              Container(
                height: 28,
                width: double.infinity,
                child: Hero(
                  child: Assets.images.logoSvg.svg(
                    fit: BoxFit.fitHeight,
                  ),
                  tag: 'logo',
                ),
              ),
              Container(
                padding: EdgeInsets.only(
                  top: 8.0,
                  bottom: 16,
                ),
                child: Text(
                  title,
                  style: const TextStyle(
                    color: ColorName.greyd8,
                    fontSize: 13,
                    fontWeight: FontWeight.w600,
                  ),
                ),
              ),
              vspace24,
              Expanded(
                child: Column(
                  children: children,
                ),
              )
            ],
          ),
        ),
        if (withCloseButton)
          Positioned(
              right: 12.0,
              top: 12.0,
              child: GestureDetector(
                onTap: () {
                  Get.back();
                },
                child: Container(
                  padding: const EdgeInsets.all(6),
                  width: 36,
                  height: 36,
                  color: Colors.transparent,
                  child: Assets.images.closeIcon.svg(width: 24, height: 24),
                ),
              ))
      ],
    );
  }
}
