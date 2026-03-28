import 'package:flutter/material.dart';
import 'package:local_auth/auth_strings.dart';
import 'package:local_auth/local_auth.dart';

import '../generated/colors.gen.dart';
import '../utils/mixins/popups.dart';

const iosStrings = const IOSAuthMessages(
    cancelButton: 'cancel',
    goToSettingsButton: 'settings',
    goToSettingsDescription: 'Please set up your Touch ID.',
    lockOut: 'Please reenable your Touch ID');

const androidAuthStrings = AndroidAuthMessages(
  biometricHint: "",
  cancelButton: 'Cancel',
  goToSettingsButton: 'Setting',
  goToSettingsDescription: 'Please set up your Touch ID.',
  biometricSuccess: "Authentication Successfull",
);

class BiometricsService with Popups {
  Future<bool> authenticateWithBiometrics(
      {bool promptToEnableInSetting = true}) async {
    final LocalAuthentication localAuthentication = LocalAuthentication();
    bool isBiometricSupported = await localAuthentication.isDeviceSupported();
    bool canCheckBiometrics = await localAuthentication.canCheckBiometrics;

    bool isAuthenticated = false;

    if (isBiometricSupported && canCheckBiometrics) {
      try {
        isAuthenticated = await localAuthentication.authenticate(
          localizedReason: 'Please complete the biometrics to proceed.',
          biometricOnly: true,
          stickyAuth: true,
          useErrorDialogs: true,
          iOSAuthStrings: iosStrings,
          androidAuthStrings: androidAuthStrings,
        );
      } catch (e) {
        if (e.toString().contains('biometrics enrolled on this device') &&
            promptToEnableInSetting) {
          openConfirmation(
            onConfirm: () {},
            cancelText: 'Ok',
            titleDistanceFromTop: -15.0,
            titleWidget: Container(
              width: 200.0,
              child: RichText(
                textAlign: TextAlign.center,
                text: TextSpan(
                  children: [
                    TextSpan(
                      text: 'Please set the biometrics from device settings',
                      style: TextStyle(
                        color: ColorName.white,
                        height: 1.4,
                      ),
                    ),
                  ],
                ),
              ),
            ),
          );
        }
      }
    }

    return isAuthenticated;
  }

  static Future<bool> hasBiometrics() async {
    final LocalAuthentication localAuthentication = LocalAuthentication();
    return await localAuthentication.isDeviceSupported();
  }
}
