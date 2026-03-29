import 'package:dio/dio.dart' show DioError;
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:get_storage/get_storage.dart';
import '../../../global/controller/globalController.dart';
import '../../../routes/app_pages.dart';
import '../../../../services/apiService.dart';
import '../../../../services/localAuthService.dart';
import '../../../../services/storageKeys.dart';
import '../../../../utils/commonUtils.dart';
import '../../../../utils/cryptography/encoding.dart';
import '../../../../utils/logger.dart';
import '../../../../utils/mixins/popups.dart';
import '../../../../utils/mixins/toast.dart';
import '../models/token_model.dart';
import '../../../../utils/emailValidator.dart';
import '../../../global/providers/commonDataProvider.dart';
import '../../account/controllers/account_controller.dart';
import '../../account/user_model.dart';
import '../models/login_model.dart';
import '../providers/authenticationProvider.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';

class LoginController extends GetxController with Toaster, Popups {
  final _secureStorage = new FlutterSecureStorage();
  AccountController accountController;
  final GlobalController globalController = Get.find();
  final storage = GetStorage();
  List<List<String>> yVals = [];
  List<List<String>> xVals = [];

  final authenticationProvider = AuthenticationProvider();
  final commonDataProvider = CommonDataProvider();

  final loginEmail = ''.obs;
  final loginPassword = ''.obs;

  final isLoggingIn = false.obs;
  bool loginLoaded = false;
  final loginEmailError = ''.obs;
  final loginPasswordError = ''.obs;
  DateTime startTime;

  Future<bool> login({LoginModel loginData}) async {
    bool isLoggedIn = false;
    loginEmail.value = loginEmail.value.replaceAll(' ', '');
    loginEmailError.value = validateEmail(loginEmail.value);
    // loginPasswordError.value = validatePassword(loginPassword.value);
    ApiService.token = null;
    if (loginEmailError.value == '') {
      isLoggingIn.value = true;
      try {
        final sendingData = loginData ??
            LoginModel(
              username: loginEmail.value.replaceAll(' ', ''),
              password: loginPassword.value,
              recaptcha: await genarateEnc(startTime: startTime),
            );
        final response = await authenticationProvider.login(
          data: sendingData,
        );
        final token = TokenModel.fromJson(response);
        if (token.token != null && token.token != '') {
          isLoggedIn = true;
          if (!(GetPlatform.isWeb)) {
            final hasBiometric = await BiometricsService.hasBiometrics();
            if (hasBiometric) {
              await _secureStorage.write(
                  key: SecureStorageKeys.email,
                  value: loginEmail.value.replaceAll(' ', ''));
              await _secureStorage.write(
                  key: SecureStorageKeys.password, value: loginPassword.value);
            }
          }
          storage.write(StorageKeys.token, token.token);
          storage.write(StorageKeys.refresh, token.refreshToken);

          storage.write(StorageKeys.lastLoginDate, DateTime.now().toString());
          Get.put(AccountController(), permanent: true);
          accountController = Get.find();

          final commonData = await Future.wait([
            commonDataProvider.getUserData(),
            commonDataProvider.getFavoritePairs(),
          ]);

          accountController.accountData.value = UserModel.fromJson(
            commonData[0]["data"],
          );
          storage.write(StorageKeys.channel,
              accountController.accountData.value.channelName);
          storage.write(
            StorageKeys.favPairs,
            commonData[1]["data"],
          );
          globalController.loadAuthenticatedControllers();
          Get.offAllNamed(AppPages.HOME);
        } else if (token.token == '') {
          if (response["need2fa"] == true ||
              response["needEmailCode"] == true ||
              response['isNewDevice'] == true) {
            await openVerificationPopup(
                need2fa: response['need2fa'],
                needEmailCode: response['needEmailCode'],
                isNewDevice: false,
                onSubmit: (v) async {
                  sendingData.s2faCode = v['s2fa'];
                  sendingData.emailCode = v['emailCode'];
                  sendingData.recaptcha =
                      await genarateEnc(startTime: startTime);
                  isLoggedIn = await login(loginData: sendingData);
                  Get.back();
                  return Future.value(isLoggedIn);
                });
          }
        }
      } on DioError catch (e) {
        final loginEmail = _secureStorage.read(key: SecureStorageKeys.email);
        if (loginEmail != null && loginLoaded == false) {
          await _secureStorage.deleteAll();
          storage.write(StorageKeys.biometricsActivated, false);
          Get.offAllNamed(AppPages.LOGIN);
        } else {
          toastDioError(e);
        }
      } catch (e) {
        log.e("Login: " + e.toString());
      } finally {
        isLoggingIn.value = false;
      }
    }
    return Future.value(isLoggedIn);
  }

  @override
  void onInit() {
    super.onInit();
  }

  @override
  void onReady() {
    super.onReady();
  }

  @override
  void onClose() {
    isLoggingIn.value = false;
  }

  void reset() {
    loginEmail.value = '';
    loginPassword.value = '';
    isLoggingIn.value = false;
    loginEmailError.value = '';
    loginPasswordError.value = '';
  }

  handleEmailChange(String v) {
    loginEmailError.value = '';
    loginEmail.value = v;
  }

  Future checkForBiometricLogin() async {
    bool canLogin = false;
    if (!(GetPlatform.isWeb)) {
      try {
        final hasBiometric = await BiometricsService.hasBiometrics();
        if (hasBiometric) {
          final isBiometricActivated =
              storage.read(StorageKeys.biometricsActivated) == true;
          final un = await _secureStorage.read(key: SecureStorageKeys.email);
          final p = await _secureStorage.read(key: SecureStorageKeys.password);
          if (un != null && p != null) {
            final canContinue = await canContinueWithBiometrics();
            if (canContinue && isBiometricActivated) {
              // generate ub_captcha
              final recaptcha = await genarateEnc(
                //create a fake start time for biometric authentication
                startTime: DateTime.now().subtract(20.seconds),
              );

              loginEmail.value = un;
              loginPassword.value = p;
              final loginModel = LoginModel(
                username: un,
                password: p,
                recaptcha: recaptcha,
              );
              canLogin = await login(loginData: loginModel);
            } else {
              await _secureStorage.deleteAll();
            }
          }
        }
      } catch (e) {
        debugPrint(e.toString());
      }
    }
    return Future.value(canLogin);
  }

  int counter = 0;
  onSlide(Map<String, String> slide) {
    final part = counter ~/ 40;
    if (xVals.length == part) {
      xVals.add([]);
    } else {
      xVals[part].add(slide['x']);
    }
    counter++;
  }

  void resetramziVars() {
    xVals = yVals = [];
    counter = 0;
  }

  void setStartTime() {
    loginLoaded = true;
    startTime = DateTime.now();
  }
}
