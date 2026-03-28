import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../generated/colors.gen.dart';

enum ScrollDirection { Up, Down }

//vertical space
const emptyComponent = const SizedBox();
const fill = const Spacer();
const vspace2 = const SizedBox(
  height: 2.0,
);
const vspace4 = const SizedBox(
  height: 4.0,
);
const vspace8 = const SizedBox(
  height: 8.0,
);
const vspace10 = const SizedBox(
  height: 10.0,
);
const vspace12 = const SizedBox(
  height: 12.0,
);
const vspace16 = const SizedBox(
  height: 16.0,
);
const vspace20 = const SizedBox(
  height: 20.0,
);
const vspace24 = const SizedBox(
  height: 24.0,
);
const vspace48 = const SizedBox(
  height: 48.0,
);
//horizontal space
const hspace4 = const SizedBox(
  width: 4.0,
);
const hspace6 = const SizedBox(
  width: 6.0,
);
const hspace2 = const SizedBox(
  width: 2.0,
);
const hspace8 = const SizedBox(
  width: 8.0,
);
const hspace12 = const SizedBox(
  width: 12.0,
);
const hspace16 = const SizedBox(
  width: 16.0,
);
const hspace20 = const SizedBox(
  width: 20.0,
);
const hspace24 = const SizedBox(
  width: 24.0,
);
// text styles
const grey80Bold10 = const TextStyle(
  color: ColorName.grey80,
  fontSize: 10.0,
  fontWeight: FontWeight.w600,
);
const grey80Bold11 = const TextStyle(
  color: ColorName.grey80,
  fontSize: 11.0,
  fontWeight: FontWeight.w600,
);
const grey80Bold12 = const TextStyle(
  color: ColorName.grey80,
  fontSize: 12.0,
  fontWeight: FontWeight.w600,
);
const grey80Bold13 = const TextStyle(
  color: ColorName.grey80,
  fontSize: 13.0,
  fontWeight: FontWeight.w600,
);
const redBold13 = const TextStyle(
  color: ColorName.red,
  fontSize: 13.0,
  fontWeight: FontWeight.w600,
);

const whiteBold8 = const TextStyle(
  color: ColorName.white,
  fontSize: 8.0,
  fontWeight: FontWeight.w600,
);
const whiteBold12 = const TextStyle(
  color: ColorName.white,
  fontSize: 12.0,
  fontWeight: FontWeight.w600,
);
const whiteBold13 = const TextStyle(
  color: ColorName.white,
  fontSize: 13.0,
  fontWeight: FontWeight.w600,
);
const whiteBold14 = const TextStyle(
  color: ColorName.white,
  fontSize: 14.0,
  fontWeight: FontWeight.w600,
);
const whiteBold24 = const TextStyle(
  color: ColorName.white,
  fontSize: 24.0,
  fontWeight: FontWeight.w600,
);
const grey97Bold24 = const TextStyle(
  color: ColorName.grey97,
  fontSize: 24.0,
  fontWeight: FontWeight.w600,
);
const whiteBold10 = const TextStyle(
  fontSize: 10,
  color: ColorName.white,
  fontWeight: FontWeight.bold,
);
const whiteBold11 = const TextStyle(
  fontSize: 10,
  color: ColorName.white,
  fontWeight: FontWeight.bold,
);
const greyBold10 = const TextStyle(
  fontSize: 10,
  color: ColorName.greybf,
  fontWeight: FontWeight.bold,
);
const grey97Bold10 = const TextStyle(
  fontSize: 10,
  color: ColorName.grey97,
  fontWeight: FontWeight.bold,
);

//borders
const rowDecoration = const BoxDecoration(
  border: const Border(
    bottom: const BorderSide(
      width: 1.0,
      color: ColorName.grey16,
    ),
  ),
);

const borderBottomGrey42 = const Border(
  bottom: const BorderSide(
    width: 1.0,
    color: ColorName.grey42,
  ),
);
const borderBottomBlack2c = const Border(
  bottom: const BorderSide(
    width: 1.0,
    color: ColorName.black2c,
  ),
);
//padding
const px12 = const EdgeInsets.symmetric(horizontal: 12.0);
//borderRadius
const roundedTop_big = const BorderRadius.only(
  topLeft: const Radius.circular(16.0),
  topRight: const Radius.circular(16.0),
);
const rounded2 = const BorderRadius.all(
  const Radius.circular(2.0),
);
const rounded6 = const BorderRadius.all(
  const Radius.circular(6.0),
);
const rounded7 = const BorderRadius.all(
  const Radius.circular(7.0),
);
const rounded_big = const BorderRadius.all(
  const Radius.circular(12.0),
);
//icons
const closeIcon32Icon16 = const SizedBox(
  width: 32.0,
  height: 32.0,
  child: const Icon(
    Icons.close,
    color: ColorName.grey80,
    size: 16.0,
  ),
);

//Divider
const dividerGrey42 = const Divider(
  color: ColorName.grey42,
  height: 1,
);
final thirdWidthPlus24 = (Get.width / 3) + 24;

headerWithCloseButton(
    {@required String title, bool centerTitle, bool noContentPadding}) {
  return Container(
    height: 36,
    padding: noContentPadding == true ? const EdgeInsets.all(0) : px12,
    decoration: const BoxDecoration(
      border: borderBottomGrey42,
    ),
    child: Stack(
      children: [
        SizedBox(
          child: Align(
            alignment:
                centerTitle == true ? Alignment.center : Alignment.centerLeft,
            child: Text(
              title,
              style: whiteBold13,
            ),
          ),
        ),
        Positioned(
          right: 0,
          child: GestureDetector(
            onTap: () {
              Get.back();
            },
            child: Container(
              child: closeIcon32Icon16,
              color: Colors.transparent,
            ),
          ),
        )
      ],
    ),
  );
}
