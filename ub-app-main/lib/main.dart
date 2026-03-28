import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:get/get.dart';
import 'package:get_storage/get_storage.dart';

import 'app/common/custom/toaster/utopic_toast.dart';
import 'app/global/binding/index.dart';
import 'app/routes/app_pages.dart';
import 'configure_nonweb.dart' if (dart.library.html) 'configure_web.dart';
import 'generated/colors.gen.dart';
import 'generated/locales.g.dart';
import 'theme/theme.dart';

void main() async {
  configureApp();
  await GetStorage.init();

  SystemChrome.setPreferredOrientations([
    DeviceOrientation.portraitDown,
    DeviceOrientation.portraitUp,
  ]);
  SystemChrome.setEnabledSystemUIMode(SystemUiMode.manual, overlays: [
    SystemUiOverlay.top,
    SystemUiOverlay.bottom,
  ]);
  SystemChrome.setSystemUIOverlayStyle(
    SystemUiOverlayStyle(
      statusBarColor: ColorName.black2c,
      statusBarBrightness: Brightness.light,
      statusBarIconBrightness: Brightness.light, //statusbar Icon Brightness
      systemNavigationBarIconBrightness:
          Brightness.light, //navigationbar Icon Brightness
    ),
  );

  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  final storage = GetStorage();

  @override
  Widget build(BuildContext context) {
    storage.writeIfNull('lightMode', false);

    final isLight = storage.read('lightMode');

    return GetMaterialApp(
      debugShowCheckedModeBanner: false,
      initialBinding: GlobalBinding(),
      builder: (context, child) {
        return ToastOverlay(child: child);
      },
      initialRoute: AppPages.INITIAL,
      theme: isLight == true ? lightThemeData : darkThemeData,
      defaultTransition: Transition.fadeIn,
      getPages: AppPages.routes,
      locale: Get.deviceLocale,
      fallbackLocale: const Locale(
        'en',
        'EN',
      ),
      translationsKeys: AppTranslation.translations,
    );
  }
}
