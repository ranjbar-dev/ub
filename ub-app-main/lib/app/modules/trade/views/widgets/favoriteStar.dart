import 'package:flutter/material.dart';
import 'package:get/state_manager.dart';
import '../../../market/controllers/market_controller.dart';
import '../../../../../generated/assets.gen.dart';
import '../../../../../generated/colors.gen.dart';

class FavoriteStar extends GetView<MarketController> {
  final String pairName;
  FavoriteStar({this.pairName});
  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.only(top: 1),
      color: ColorName.black2c,
      alignment: Alignment.topCenter,
      child: Obx(() {
        final bool isFavorite =
            controller.isPairNameFavorite(pairName: pairName);
        return GestureDetector(
          onTap: () {
            controller.toggleFavoriteByPairName(pairName: pairName);
          },
          child: SizedBox(
            width: 30,
            height: 30,
            child: Assets.images.filledStar.svg(
              color: isFavorite ? ColorName.primaryBlue : ColorName.greybf,
            ),
          ),
        );
      }),
    );
  }
}
