import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../../generated/colors.gen.dart';
import '../../../../generated/locales.g.dart';
import '../../../../utils/mixins/commonConsts.dart';
import '../../../../utils/mixins/popups.dart';
import '../../../common/components/UBBorderlessInput.dart';
import '../../../common/components/UBButton.dart';
import '../../../common/components/UBDDMockButton.dart';
import '../../../common/components/UBGreyContainer.dart';
import '../../../common/components/UBText.dart';
import '../controllers/phone_verification_controller.dart';

class EnterPhoneNumber extends GetView<PhoneVerificationController>
    with Popups {
  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        fill,
        Container(
          decoration: const BoxDecoration(
            color: ColorName.black2c,
            borderRadius: roundedTop_big,
          ),
          height: 252.0,
          child: Column(
            children: [
              vspace24,
              vspace12,
              Obx(() {
                final selectedCountry = controller.selectedCountry.value;
                return UBDDMockButton(
                  title: selectedCountry.id == null
                      ? LocaleKeys.selectCountry.tr
                      : selectedCountry.name,
                  onTap: () => openCountrySelect(
                      onCountrySelect: controller.handleCountrySelected),
                );
              }),
              vspace24,
              UBGreyContainer(
                color: ColorName.black,
                margin: const EdgeInsets.symmetric(
                  horizontal: 12,
                ),
                child: Obx(
                  () {
                    final selectedCountry = controller.selectedCountry.value;
                    return Row(
                      children: [
                        if (selectedCountry.id != null)
                          Container(
                            height: 25,
                            //width: 25,
                            child: Center(
                              child: UBText(
                                  text: "+ ${selectedCountry.code.toString()}"),
                            ),
                          ),
                        if (selectedCountry.id != null)
                          const SizedBox(
                            width: 8,
                          ),
                        Expanded(
                          child: UBBorderlessInput(
                            type: TextInputType.number,
                            placeholder: LocaleKeys.enterPhoneNumber.tr,
                            onChange: controller.handlePhoneNumberChange,
                          ),
                        ),
                      ],
                    );
                  },
                ),
              ),
              vspace24,
              Padding(
                padding: const EdgeInsets.symmetric(horizontal: 12),
                child: Obx(
                  () {
                    final phoneNumber = controller.phoneNumber.value;
                    final selectedCountry = controller.selectedCountry.value;
                    return UBButton(
                      isLodaing: controller.isRequestingForSms.value,
                      disabled: (!(phoneNumber.length > 9) ||
                          selectedCountry.id == null),
                      onClick: controller.phoneNumberSubmitted,
                      text: LocaleKeys.beginVerification.tr,
                    );
                  },
                ),
              ),
              vspace12,
              Container(
                width: 100.0,
                height: 24.0,
                child: UBButton(
                  variant: ButtonVariant.TransparentBackground,
                  onClick: () {
                    Get.back();
                  },
                  text: LocaleKeys.cancel.tr,
                ),
              ),
            ],
          ),
        ),
      ],
    );
  }
}
