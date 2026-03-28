import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../../common/components/UBBlurryContainer.dart';
import '../../../../common/components/UBText.dart';
import '../../../../../generated/assets.gen.dart';
import '../../../../../utils/commonUtils.dart';

class HomePageSideMenu extends StatelessWidget {
  const HomePageSideMenu({Key key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return UBBlurryContainer(
      borderRadius: BorderRadius.circular(0.0),
      padding: EdgeInsets.zero,
      blur: 5,
      child: Container(
        width: 200,
        color: Color.fromRGBO(43, 43, 50, 0.7),
        child: Stack(
          children: [
            Positioned(
              top: 140.0,
              child: Assets.images.fullWatermark.svg(),
            ),
            Drawer(
              elevation: 0,
              child: ListView(
                padding: EdgeInsets.zero,
                children: [
                  Container(
                    height: 40,
                  ),
                  ListTile(
                    title: UBText(text: 'Cryptocurrency'),
                    onTap: () {
                      Get.back();
                      launchURL('https://unitedbit.com/cryptoCurrency');
                    },
                  ),
                  ListTile(
                    title: UBText(text: 'Introducing The exchange'),
                    onTap: () {
                      Get.back();
                      launchURL('https://unitedbit.com/intro');
                    },
                  ),
                  ListTile(
                    title: UBText(text: 'Privacy policy'),
                    onTap: () {
                      Get.back();
                      launchURL('https://unitedbit.com/privacy-policies');
                    },
                  ),
                  ListTile(
                    title: UBText(text: 'Contact Us'),
                    onTap: () {
                      Get.back();
                      launchURL('https://unitedbit.com/contact');
                    },
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}
