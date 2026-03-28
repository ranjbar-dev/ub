import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:rive/rive.dart';
import 'package:supercharged/supercharged.dart';

import '../../../generated/assets.gen.dart';
import '../../../generated/colors.gen.dart';
import '../../../generated/locales.g.dart';
import '../../../utils/commonUtils.dart';
import '../../global/controller/globalController.dart';
import '../../modules/home/views/widgets/redirectContainer.dart';
import '../../routes/app_pages.dart';
import 'UBConnectionLost.dart';

class PageContainer extends GetView<GlobalController> {
  final Widget appbarTitle;
  final Widget child;
  final bool singleChildScroll;
  final bool protectedPage;
  final int activeBottomNavIndex;
  final Color backgroundColor;
  final Color activeColor = ColorName.textBlue;
  final Color deactiveColor = ColorName.greybf;
  final Color appBarBottomBorderColor;
  final void Function() beforePop;
  final int additionalHeight;
  final Widget sideMenu;
  final Widget drawerOpener;

  const PageContainer(
      {Key key,
      this.drawerOpener,
      this.appBarBottomBorderColor = ColorName.grey36,
      @required this.child,
      this.backgroundColor = Colors.transparent,
      this.appbarTitle,
      this.singleChildScroll,
      this.activeBottomNavIndex,
      this.additionalHeight,
      this.protectedPage = false,
      this.sideMenu,
      this.beforePop})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    final GlobalKey<ScaffoldState> _key = GlobalKey(); // Create a key
    if (protectedPage) {
      if (controller.loggedIn.value != true) {
        debugPrint(
            '!!!!!!!!!!!!!!!!!!!!!!!!!protected page!!!!!!!!!!!!!!!!!!!!!!!!!!!');
        Get.back();
        return const SizedBox();
      }
    }
    return WillPopScope(
      onWillPop: () async {
        if (canExitApp == true || activeBottomNavIndex == null) {
          return Future<bool>.value(true);
        }
        canExitApp = true;
        Future.delayed(2.seconds).then((value) => canExitApp = false);
        if (!(GetPlatform.isWeb)) {
          showDoubleTapToast();
          return Future<bool>.value(false);
        }
        return true;
      },
      child: GestureDetector(
        //to unfocus when tapped outside keyboard
        onTap: () {
          FocusManager.instance.primaryFocus?.unfocus();
        },
        child: SafeArea(
          child: Theme(
            data: Theme.of(context).copyWith(canvasColor: Colors.transparent),
            child: Scaffold(
              key: _key,
              extendBody: true,
              endDrawer: sideMenu == null ? null : sideMenu,
              appBar: appbarTitle != null
                  ? AppBar(
                      title: appbarTitle,
                      backgroundColor: ColorName.black,
                      centerTitle: true,
                      bottom: PreferredSize(
                        child: Container(
                          color: appBarBottomBorderColor,
                          height: 1.0,
                        ),
                        preferredSize: Size.fromHeight(
                          0.0,
                        ), // here the desired height
                      ),
                    )
                  : null,
              body: Stack(
                children: [
                  singleChildScroll == true
                      ? SingleChildScrollView(
                          child: Container(
                            color: Colors.transparent,
                            //height: Get.height - (65 + (additionalHeight ?? 0)),
                            child: child,
                          ),
                        )
                      : Container(
                          width: double.infinity,
                          //height: Get.height - (65 + (additionalHeight ?? 0)),
                          color: Colors.transparent,
                          child: child,
                        ),
                  Obx(() {
                    final connected = controller.hasConnection.value;
                    return connected ? const SizedBox() : ConnectionLost();
                  }),
                  if (sideMenu != null)
                    Positioned(
                      top: 16.0,
                      right: 12,
                      child: Container(
                        width: 30,
                        height: 30,
                        color: Colors.transparent,
                        child: GestureDetector(
                          onTap: () {
                            // Scaffold.of(context).openDrawer();
                            _key.currentState.openEndDrawer();
                          },
                          child: Assets.images.bigInfoIcon.svg(),
                        ),
                      ),
                    ),
                  Positioned(
                    bottom: 2,
                    left: 0,
                    right: 0,
                    child: Obx(
                      () => Visibility(
                        visible: controller.doShowRedirect.value &&
                            activeBottomNavIndex != null,
                        child: RedirectContainer(),
                      ),
                    ),
                  )
                ],
              ),
              bottomNavigationBar: activeBottomNavIndex != null
                  ? Stack(
                      children: [
                        controller.deviceType == DeviceTypes.PHONE
                            ? Assets.images.bottomNavigationBar.svg(
                                color: ColorName.black2c.withOpacity(1),
                                width: Get.width,
                                fit: BoxFit.fitHeight)
                            : Assets.images.bottomNavigationBar.svg(
                                color: ColorName.black2c.withOpacity(1),
                                width: Get.width,
                                fit: BoxFit.fitHeight),
                        controller.deviceType == DeviceTypes.PHONE
                            ? Assets.images.bottomNavigationBar.svg(
                                color: ColorName.black2c.withOpacity(0.8),
                                width: Get.width,
                                fit: BoxFit.fitHeight)
                            : Assets.images.bottomNavigationBar.svg(
                                color: ColorName.black2c.withOpacity(0.8),
                                width: Get.width,
                                fit: BoxFit.fitHeight),
                        BottomNavigationBar(
                          selectedItemColor: activeColor,
                          unselectedItemColor: deactiveColor,
                          backgroundColor: Colors.transparent,
                          unselectedLabelStyle: const TextStyle(
                            fontWeight: FontWeight.w600,
                            fontSize: 10.0,
                          ),
                          selectedLabelStyle: const TextStyle(
                            fontWeight: FontWeight.w600,
                            fontSize: 10.0,
                          ),
                          showUnselectedLabels: true,
                          type: BottomNavigationBarType.fixed,

                          onTap: (int index) {
                            FocusManager.instance.primaryFocus?.unfocus();
                            if (GetPlatform.isWeb) {
                              FocusScope.of(context).unfocus();
                            }
                            if (index != activeBottomNavIndex) {
                              if (beforePop != null) {
                                beforePop();
                              }

                              switch (index) {
                                case 0:
                                  Get.offAllNamed(AppPages.HOME);
                                  break;
                                case 1:
                                  Get.offAllNamed(AppPages.MARKET);
                                  break;
                                case 2:
                                  Get.offAllNamed(AppPages.TRADE);
                                  break;
                                case 3:
                                  Get.offAllNamed(AppPages.EXCHANGE);
                                  break;
                                case 4:
                                  Get.offAllNamed(AppPages.FUNDS);
                                  break;
                              }
                            }
                          },
                          currentIndex:
                              activeBottomNavIndex, // this will be set when a new tab is tapped
                          items: [
                            BottomNavigationBarItem(
                              icon: Padding(
                                padding: const EdgeInsets.only(top: 8),
                                child: BottomBarItem(
                                  key: ValueKey('homeAnimatedIcon'),
                                  isActive: activeBottomNavIndex == 0,
                                  rivAnimationFileName: 'homeIcon',
                                  svgIcon: Assets.images.home,
                                  isMain: false,
                                ),
                              ),
                              label: LocaleKeys.home.tr,
                            ),
                            BottomNavigationBarItem(
                              icon: Padding(
                                padding: const EdgeInsets.only(top: 8),
                                child: BottomBarItem(
                                  key: ValueKey('marketAnimatedIcon'),
                                  isActive: activeBottomNavIndex == 1,
                                  rivAnimationFileName: 'marketIcon',
                                  svgIcon: Assets.images.market,
                                  isMain: false,
                                ),
                              ),
                              label: 'Markets',
                            ),
                            BottomNavigationBarItem(
                              icon: Column(
                                children: [
                                  SizedBox(height: 3),
                                  Container(
                                    width: 35,
                                    height: 35,
                                    child: BottomBarItem(
                                      key: ValueKey('tradeAnimatedIcon'),
                                      isActive: activeBottomNavIndex == 2,
                                      rivAnimationFileName: 'tradeIcon',
                                      svgIcon: Assets.images.trade,
                                      isMain: true,
                                    ),
                                  ),
                                  SizedBox(height: 3)
                                ],
                              ),
                              label: LocaleKeys.trade.tr,
                            ),
                            BottomNavigationBarItem(
                              icon: Padding(
                                padding: const EdgeInsets.only(top: 8),
                                child: BottomBarItem(
                                  key: ValueKey('exchangeAnimatedIcon'),
                                  isActive: activeBottomNavIndex == 3,
                                  rivAnimationFileName: 'exchangeIcon',
                                  svgIcon: Assets.images.exchange,
                                  isMain: false,
                                ),
                              ),
                              label: LocaleKeys.exchange.tr,
                            ),
                            BottomNavigationBarItem(
                              icon: Padding(
                                padding: const EdgeInsets.only(top: 8),
                                child: BottomBarItem(
                                  key: ValueKey('fundAnimatedIcon'),
                                  isActive: activeBottomNavIndex == 4,
                                  rivAnimationFileName: 'fundIcon',
                                  svgIcon: Assets.images.fund,
                                  isMain: false,
                                ),
                              ),
                              label: 'Funds',
                            ),
                          ],
                        ),
                      ],
                    )
                  : null,
            ),
          ),
        ),
      ),
    );
  }

  iconColor(int i) {
    return activeBottomNavIndex == i ? activeColor : deactiveColor;
  }
}

class BottomBarItem extends StatelessWidget {
  final bool isActive;
  final bool isMain;
  final SvgGenImage svgIcon;
  final String rivAnimationFileName;

  const BottomBarItem(
      {Key key,
      this.isActive,
      this.svgIcon,
      this.rivAnimationFileName,
      this.isMain = false})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return isActive
        ? isMain
            ? Container(
                width: 35.0,
                height: 35.0,
                decoration: BoxDecoration(
                  borderRadius: BorderRadius.circular(10),
                  color: '16161A'.toColor(),
                ),
                child: Padding(
                  padding: const EdgeInsets.all(5.0),
                  child: RiveAnimation.asset(
                    'assets/rive/$rivAnimationFileName.riv',
                  ),
                ),
              )
            : SizedBox(
                width: 32.0,
                height: 32.0,
                child: Padding(
                  padding: const EdgeInsets.only(top: 8.0),
                  child: RiveAnimation.asset(
                      'assets/rive/$rivAnimationFileName.riv'),
                ),
              )
        : isMain
            ? Container(
                width: 35.0,
                height: 35.0,
                padding: const EdgeInsets.all(0),
                decoration: BoxDecoration(
                    borderRadius: BorderRadius.circular(10),
                    color: ColorName.black3c),
                child: Padding(
                  padding: const EdgeInsets.all(5.0),
                  child: svgIcon.svg(color: ColorName.greybf),
                ),
              )
            : SizedBox(
                width: 32.0,
                height: 32.0,
                child: Padding(
                  padding: const EdgeInsets.only(top: 8.0),
                  child: svgIcon.svg(color: ColorName.greybf),
                ),
              );
  }
}
