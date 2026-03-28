import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../../generated/assets.gen.dart';
import '../../../common/components/pageContainer.dart';
import '../../../routes/app_pages.dart';
import '../controllers/home_controller.dart';
import 'HomeBody.dart';
import 'widgets/sideMenu.dart';

class HomeView extends GetView<HomeController> {
  @override
  Widget build(BuildContext context) {
    controller.handlePageLoaded();

    return PageContainer(
      //appBarBottomBorderColor: Colors.transparent,
      beforePop: () {
        controller.handlePagePop();
      },
      protectedPage: false,
      activeBottomNavIndex: 0,
      sideMenu: HomePageSideMenu(),
      child: Column(
        children: [
          Container(
            height: 60.0,
            padding: const EdgeInsets.symmetric(horizontal: 12.0),
            child: Row(
              children: [
                Container(
                  width: 120.0,
                  child: Hero(
                    tag: 'logo',
                    child: Assets.images.logoSvg.svg(),
                  ),
                ),
                const Spacer(),
                Padding(
                  padding: const EdgeInsets.only(right: 40),
                  child: GestureDetector(
                    onTap: () {
                      Get.toNamed(AppPages.ACCOUNT);
                    },
                    child: Assets.images.profileIcon.svg(),
                  ),
                )
              ],
            ),
          ),
          Expanded(
            child: HomeBody(),
          ),
        ],
      ),
    );
  }
}
