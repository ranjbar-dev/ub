import 'package:basic_utils/basic_utils.dart';
import 'package:flutter/material.dart';
import 'package:flutter_datetime_picker/flutter_datetime_picker.dart';
import 'package:get/get.dart';

import '../../../../generated/colors.gen.dart';
import '../../../../generated/locales.g.dart';
import '../../../../utils/mixins/commonConsts.dart';
import '../../../../utils/mixins/popups.dart';
import '../../../common/components/CenterUBLoading.dart';
import '../../../common/components/UBBlackContainer.dart';
import '../../../common/components/UBBottomSheetContainer.dart';
import '../../../common/components/UBButton.dart';
import '../../../common/components/UBDDMockButton.dart';
import '../../../common/components/UBScrollColumnExpandable.dart';
import '../../../common/components/UBSelectButton.dart';
import '../../../common/components/UBText.dart';
import '../../../common/components/UBToastOnTap.dart';
import '../../../common/components/UBWrappedButtons.dart';
import '../../../common/components/appbarTextTitle.dart';
import '../../../common/components/controlledInput.dart';
import '../../../common/components/pageContainer.dart';
import '../../../routes/app_pages.dart';
import '../controllers/identity_info_controller.dart';

class IdentityInfoView extends GetView<IdentityInfoController> with Popups {
  @override
  Widget build(BuildContext context) {
    return PageContainer(
      appbarTitle: AppBarTextTitle(
        title: LocaleKeys.identityVerification.tr,
      ),
      child: Container(
        width: Get.width,
        color: ColorName.black,
        child: Obx(
          () {
            final isLoading = controller.isLoading.value;
            return isLoading == true
                ? Container(
                    child: CenterUbLoading(),
                    height: Get.height - 60,
                  )
                : UBScrollColumnExpandable(
                    children: [
                      vspace24,
                      _topText(),
                      vspace24,
                      Expanded(
                          child: Container(
                        padding:
                            const EdgeInsets.only(left: 12, right: 12, top: 12),
                        decoration: BoxDecoration(
                          borderRadius: roundedTop_big,
                          color: ColorName.black2c,
                        ),
                        child: Column(
                          children: [
                            _basicInfoHeaderText(),
                            UBToastOnTap(
                              active: controller.hasAcceptedIdentityImage.value,
                              toastText:
                                  "You can't change this field, you have an accepted identity document",
                              child: Column(
                                children: [
                                  _firstName(),
                                  vspace20,
                                  _lastName(),
                                  vspace20,
                                  _genderSelect(),
                                  vspace20,
                                  _birthdaySelect(),
                                ],
                              ),
                            ),
                            const Spacer(),
                            _residentialInfoHeaderText(),
                            UBToastOnTap(
                              active:
                                  controller.hasAcceptedResidenceImage.value,
                              toastText:
                                  "You can't change this field, you have an accepted Residence document",
                              child: Column(
                                children: [
                                  _countrySelect(),
                                  vspace20,
                                  _cityInput(),
                                  vspace20,
                                  _postalCodeInput(),
                                  vspace20,
                                  _addressInput(),
                                ],
                              ),
                            ),
                            vspace20,
                            _submitButton(),
                            vspace20,
                            _cancelButton(),
                            vspace20,
                          ],
                        ),
                      ))
                    ],
                  );
          },
        ),
      ),
    );
  }

  RichText _topText() {
    return RichText(
      textAlign: TextAlign.center,
      text: TextSpan(
          text: LocaleKeys.identityVerificationTopText.tr,
          style: const TextStyle(
            color: ColorName.greybf,
            fontSize: 13,
            fontWeight: FontWeight.w600,
          )),
    );
  }

  void _openGenderSelect({@required int selectedGender}) {
    Get.bottomSheet(
      UBButtomSheetContainer(
        title: LocaleKeys.gender.tr,
        child: Container(
          child: UBWrappedButtons(
            buttons: [
              WrappedButtonModel(
                text: 'Male',
              ),
              WrappedButtonModel(
                text: 'Female',
              ),
            ],
            onButtonClick: (v) {
              Navigator.of(Get.context).pop(); // Navigator
              controller.handleGenderChange(
                v,
              );
            },
            selectedIndex: selectedGender,
          ),
        ),
      ),
    );
  }

  void _openDatePicker() {
    DatePicker.showDatePicker(Get.context,
        showTitleActions: true,
        minTime: DateTime(1900, 1, 1),
        maxTime: DateTime.now(),
        theme: DatePickerTheme(
          headerColor: ColorName.black2c,
          backgroundColor: ColorName.grey16,
          itemStyle: TextStyle(
            color: Colors.white,
            fontWeight: FontWeight.bold,
            fontSize: 18,
          ),
          doneStyle: TextStyle(
            color: Colors.white,
            fontSize: 16,
          ),
          cancelStyle: TextStyle(
            color: Colors.white,
            fontSize: 16,
          ),
        ), onConfirm: (date) {
      controller.handleBirthdayChange(date);
    }, currentTime: DateTime.now(), locale: LocaleType.en);
  }

  _basicInfoHeaderText() {
    return Container(
      alignment: Alignment.centerLeft,
      padding: const EdgeInsets.symmetric(vertical: 12.0, horizontal: 2.0),
      child: UBText(
        text: LocaleKeys.basicInfo.tr,
        color: ColorName.greybf,
      ),
    );
  }

  _firstName() {
    return UBBlackContainer(
      child: Obx(
        () => ControlledTextField(
          labelText: LocaleKeys.firstName.tr,
          text: controller.firstName.value,
          noBorder: true,
          onChanged: (v) {
            controller.handleFirstNameChange(v);
          },
        ),
      ),
    );
  }

  _lastName() {
    return UBBlackContainer(
      child: Obx(
        () => ControlledTextField(
          labelText: LocaleKeys.lastName.tr,
          text: controller.lastName.value,
          noBorder: true,
          onChanged: (v) {
            controller.handleLastNameChange(v);
          },
        ),
      ),
    );
  }

  _genderSelect() {
    return Obx(() {
      final gender = controller.gender.value;
      return UBSelectButton(
        backgroundColor: ColorName.black,
        onClick: () {
          _openGenderSelect(selectedGender: gender == 'female' ? 1 : 0);
        },
        valueText: gender == ''
            ? LocaleKeys.gender.tr
            : StringUtils.capitalize(gender),
      );
    });
  }

  _birthdaySelect() {
    return Obx(() {
      final birthday = controller.birthday.value;
      return UBSelectButton(
        backgroundColor: ColorName.black,
        onClick: () {
          _openDatePicker();
        },
        valueText: birthday == '' ? LocaleKeys.birthday.tr : birthday,
        icon: Icons.calendar_today_outlined,
        iconSize: 16,
        padding: const EdgeInsets.only(left: 12.0, right: 16.0),
      );
    });
  }

  _residentialInfoHeaderText() {
    return Container(
      alignment: Alignment.centerLeft,
      padding: const EdgeInsets.symmetric(vertical: 12.0, horizontal: 2.0),
      child: UBText(
        text: LocaleKeys.residentialAddress.tr,
        color: ColorName.greybf,
      ),
    );
  }

  _countrySelect() {
    return Obx(() {
      final selectedCountry = controller.selectedCountry.value;
      String country;

      if (controller.country.value != '') {
        country = controller.country.value;
      } else if (selectedCountry.id == null) {
        country = LocaleKeys.selectCountry.tr;
      } else {
        country = selectedCountry.name;
      }

      return UBDDMockButton(
        horizontalPadding: 0.0,
        title: country,
        onTap: () => openCountrySelect(
            onCountrySelect: controller.handleCountrySelected),
      );
    });
  }

  _cityInput() {
    return UBBlackContainer(
      child: Obx(
        () => ControlledTextField(
          labelText: LocaleKeys.city.tr,
          text: controller.city.value,
          noBorder: true,
          onChanged: (v) {
            controller.handleCityChange(v);
          },
        ),
      ),
    );
  }

  _postalCodeInput() {
    return UBBlackContainer(
      child: Obx(
        () => ControlledTextField(
          labelText: LocaleKeys.postalCode.tr,
          text: controller.postalCode.value,
          noBorder: true,
          onChanged: (v) {
            controller.handlePostalCodeChange(v);
          },
        ),
      ),
    );
  }

  _addressInput() {
    return UBBlackContainer(
      child: Obx(
        () => ControlledTextField(
          labelText: LocaleKeys.address.tr,
          text: controller.address.value,
          noBorder: true,
          onChanged: (v) {
            controller.handleAddressChange(v);
          },
        ),
      ),
    );
  }

  _submitButton() {
    return Obx(
      () {
        final canSubmit = controller.canSubmit();
        return UBButton(
          disabled: !canSubmit,
          isLodaing: controller.isSubmitting.value,
          onClick: () {
            controller.handleSubmitButtonClick();
          },
          text: LocaleKeys.beginVerification.tr,
        );
      },
    );
  }

  _cancelButton() {
    return Container(
      width: 140,
      child: UBButton(
          onClick: () {
            Get.offNamed(AppPages.ACCOUNT);
          },
          variant: ButtonVariant.TransparentBackground,
          textColor: ColorName.grey80,
          text: LocaleKeys.cancel.tr),
    );
  }
}
