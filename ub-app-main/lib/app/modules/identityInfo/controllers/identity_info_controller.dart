import 'package:dio/dio.dart' show DioError;
import 'package:flutter/foundation.dart';
import 'package:get/get.dart';
import '../../../global/autocompleteModel.dart';
import '../../identityDocuments/providers/identityDocumentsProvider.dart';
import '../../identityDocuments/user_profile_model.dart';
import '../providers/identityInfoProvider.dart';
import '../update_user_info_model.dart';
import '../../../routes/app_pages.dart';
import '../../../../generated/locales.g.dart';
import '../../../../services/constants.dart';
import '../../../../utils/mixins/toast.dart';

class IdentityInfoController extends GetxController with Toaster {
  final IdentityDocumentsProvider identityDocumentsProvider =
      IdentityDocumentsProvider();
  final IdentityInfoProvider identityInfoProvider = IdentityInfoProvider();
  final isLoading = true.obs;
  final Rx<UserProfileModel> userProfile = UserProfileModel().obs;

  final firstName = ''.obs;
  final lastName = ''.obs;
  final gender = ''.obs;
  final birthday = ''.obs;
  final country = ''.obs;
  final city = ''.obs;
  final postalCode = ''.obs;
  final address = ''.obs;
  final selectedCountry = AutoCompleteItem(name: '', code: "0").obs;
  final isSubmitting = false.obs;
  bool isDataChanged = false;
  final hasAcceptedIdentityImage = false.obs;
  final hasAcceptedResidenceImage = false.obs;
  bool shouldSkip = false;

  @override
  void onInit() {
    super.onInit();
    getUserProfile();
  }

  @override
  void onReady() {
    super.onReady();
  }

  @override
  void onClose() {}

  void getUserProfile() async {
    try {
      isLoading.value = true;
      final response = await identityDocumentsProvider.getUserProfile();
      if (response['status'] == true) {
        userProfile.value = UserProfileModel.fromJson(response['data']);
        shouldSkip = _proccessImages(userProfile.value.userProfileImages);
        if (shouldSkip) {
          shouldSkip = true;
          Get.offNamed(AppPages.IDENTITYDOCUMENTS);
          return;
        }
        final profile = userProfile.value;
        firstName.value = profile.firstName ?? '';
        lastName.value = profile.lastName ?? '';
        gender.value = profile.gender ?? LocaleKeys.gender.tr;
        birthday.value = profile.dateOfBirth ?? LocaleKeys.birthday.tr;
        country.value = profile.countryName ?? LocaleKeys.country.tr;
        if (country.value != null && country.value != LocaleKeys.country.tr) {
          country.value = Constants.countriesArray()
              .firstWhere(
                  (element) => element.inPerentesis == profile.countryName)
              .name;
        }
        city.value = profile.regionAndCity ?? '';
        postalCode.value = profile.postalCode ?? '';
        address.value = profile.address ?? '';
      }
    } on DioError catch (e) {
      toastDioError(e);
    } catch (e) {
      debugPrint(e.toString());
    } finally {
      if (shouldSkip == false) {
        isLoading.value = false;
      }
    }
  }

  void handleFirstNameChange(String v) {
    isDataChanged = true;
    firstName.value = v;
  }

  void handleLastNameChange(String v) {
    isDataChanged = true;
    lastName.value = v;
  }

  void handleGenderChange(int v) {
    isDataChanged = true;
    gender.value = v == 0 ? 'male' : 'female';
  }

  void handleCountrySelected(AutoCompleteItem item) {
    isDataChanged = true;
    selectedCountry.value = item;
    country.value = item.name;
  }

  void handleBirthdayChange(DateTime date) {
    isDataChanged = true;
    birthday.value = date.toString().split(' ')[0];
  }

  void handleCityChange(String v) {
    isDataChanged = true;
    city.value = v;
  }

  void handlePostalCodeChange(String v) {
    isDataChanged = true;
    postalCode.value = v;
  }

  void handleAddressChange(String v) {
    isDataChanged = true;
    address.value = v;
  }

  void handleSubmitButtonClick() async {
    if (isDataChanged == false) {
      Get.offNamed(AppPages.IDENTITYDOCUMENTS);
      return;
    }
    final sendingData = {
      "address": address.value,
      "country_id": selectedCountry.value.id,
      "date_of_birth": birthday.value,
      "first_name": firstName.value,
      "gender": gender.value,
      "last_name": lastName.value,
      "postal_code": postalCode.value,
      "region_and_city": city.value
    };
    try {
      Get.offNamed(AppPages.IDENTITYDOCUMENTS);

      isSubmitting.value = true;

      final response = await identityInfoProvider.updateUserInfo(
          data: UpdateUserInfoModel.fromJson(sendingData));

      if (response['status'] == true) {
        Get.offNamed(AppPages.IDENTITYDOCUMENTS);
      }
    } on DioError catch (e) {
      toastDioError(e);
    } finally {
      isSubmitting.value = false;
    }
  }

  canSubmit() {
    return firstName.value != '' &&
        lastName.value != '' &&
        gender.value != LocaleKeys.gender.tr &&
        birthday.value != LocaleKeys.birthday.tr &&
        country.value != LocaleKeys.country.tr &&
        city.value != '' &&
        postalCode.value != '' &&
        address.value != '';
  }

  bool _proccessImages(List<UserProfileImages> userProfileImages) {
    //print(userProfileImages);
    bool identityAcepted = false;
    bool addressAccepted = false;
    for (var image in userProfileImages) {
      if (image.type == 'identity' &&
          (image.status == 'confirmed' ||
              image.status == 'partially_confirmed')) {
        identityAcepted = true;
      }
      if (image.type == 'address' &&
          (image.status == 'confirmed' ||
              image.status == 'partially_confirmed')) {
        addressAccepted = true;
      }
    }
    if (addressAccepted) {
      hasAcceptedResidenceImage.value = true;
    }
    if (identityAcepted) {
      hasAcceptedIdentityImage.value = true;
    }
    return identityAcepted && addressAccepted;
  }
}
