import 'package:flutter/material.dart';

import '../generated/colors.gen.dart';
import '../generated/fonts.gen.dart';

final ThemeData darkThemeData = ThemeData(
  fontFamily: FontFamily.openSans,
  highlightColor: ColorName.primaryBlue.shade300,
  hintColor: ColorName.inputPlaceholderText,
  primaryColor: ColorName.black,
  splashColor: ColorName.primaryBlue.shade100,
  canvasColor: Colors.transparent,
  textTheme: const TextTheme(
    headline1: const TextStyle(fontSize: 72.0, fontWeight: FontWeight.bold),
  ),
  colorScheme:
      ColorScheme.fromSwatch().copyWith(secondary: ColorName.yellowOcher),
);
final ThemeData lightThemeData = ThemeData(
  primaryColor: ColorName.primaryBlue,
  splashColor: ColorName.primaryBlue.shade500,
  highlightColor: ColorName.crimsonRed,
  canvasColor: Colors.transparent,
  fontFamily: FontFamily.openSans,
  textTheme: const TextTheme(
    headline1: const TextStyle(fontSize: 72.0, fontWeight: FontWeight.bold),
  ),
  colorScheme:
      ColorScheme.fromSwatch().copyWith(secondary: ColorName.yellowOcher),
);
