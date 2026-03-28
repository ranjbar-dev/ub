import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../../routes/app_pages.dart';
import '../../../../../generated/assets.gen.dart';

class HomePageBanners extends StatelessWidget {
  const HomePageBanners({Key key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return AspectRatio(
      aspectRatio: 16 / 7,
      child: GestureDetector(
        onTap: () {
          Get.offAllNamed(AppPages.TRADE);
        },
        child: Assets.images.homePageBanner.image(),
      ),
    );
  }
}
