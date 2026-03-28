import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../../generated/assets.gen.dart';
import '../../../../generated/colors.gen.dart';
import '../../../../generated/locales.g.dart';
import '../../../../utils/mixins/commonConsts.dart';
import '../../../../utils/mixins/popups.dart';
import '../../../common/components/UBBorderlessInput.dart';
import '../../../common/components/UBBottomSheetContainer.dart';
import '../../../common/components/UBButton.dart';
import '../../../common/components/UBDDMockButton.dart';
import '../../../common/components/UBGreyContainer.dart';
import '../../../common/components/UBWrappedButtons.dart';
import '../../../common/components/appbarTextTitle.dart';
import '../../../common/components/controlledInput.dart';
import '../../../common/components/pageContainer.dart';
import '../../funds/withdraw_deposit_data_model.dart';
import '../controllers/add_new_address_controller.dart';

class AddNewAddressView extends GetView<AddNewAddressController> with Popups {
  @override
  Widget build(BuildContext context) {
    return PageContainer(
      appbarTitle: AppBarTextTitle(
        title: LocaleKeys.addNewAddress.tr,
      ),
      child: Column(
        children: [
          fill,
          Container(
            height: 400,
            decoration: const BoxDecoration(
              borderRadius: roundedTop_big,
              color: ColorName.black2c,
            ),
            child: Column(
              children: [
                vspace24,
                _coinSelector(),
                _networkConditionalSpace(),
                _networkSelector(),
                vspace24,
                _walletLabelInput(),
                vspace24,
                _addressInput(),
                const Spacer(),
                Padding(
                  padding: const EdgeInsets.symmetric(
                    horizontal: 12,
                  ),
                  child: Obx(
                    () {
                      final isLoading = controller.isAddingNewAddress.value;
                      final address = controller.address.value;
                      final label = controller.newAddressLabel.value;
                      // ignore: invalid_use_of_protected_member
                      final networks = controller.networks.value;
                      final selectedNetwork = controller.selectedNetwork.value;
                      final shouldSelectNetwork =
                          networks.length > 0 && selectedNetwork.code == null;

                      return UBButton(
                          isLodaing: isLoading,
                          disabled: address.length < 15 ||
                              label.length < 3 ||
                              shouldSelectNetwork,
                          onClick: () {
                            controller.handleCreateClick();
                          },
                          text: LocaleKeys.create.tr);
                    },
                  ),
                ),
                Container(
                  width: 120,
                  padding:
                      const EdgeInsets.symmetric(horizontal: 12, vertical: 24),
                  child: UBButton(
                    onClick: () {
                      Get.back();
                    },
                    variant: ButtonVariant.TransparentBackground,
                    textColor: ColorName.grey80,
                    text: LocaleKeys.cancel.tr,
                  ),
                )
              ],
            ),
          ),
        ],
      ),
    );
  }

  Obx _coinSelector() {
    return Obx(
      () {
        final selectedCoin = controller.selectedCoin.value;
        return UBDDMockButton(
          title: selectedCoin.id == null
              ? LocaleKeys.pleaseselectanycoin.tr
              : selectedCoin.name,
          titleAppendix: selectedCoin.desc,
          iconAddress: selectedCoin.image,
          onTap: () =>
              openCoinSelectPopup(onCoinSelect: controller.handleCoinSelected),
        );
      },
    );
  }

  Obx _networkConditionalSpace() {
    return Obx(
      () {
        // ignore: invalid_use_of_protected_member
        final networks = controller.networks.value;
        return networks.length > 0 ? vspace24 : const SizedBox();
      },
    );
  }

  Obx _networkSelector() {
    return Obx(
      () {
        final selectedNetwork = controller.selectedNetwork.value;
        // ignore: invalid_use_of_protected_member
        final networks = controller.networks.value;
        return networks.length > 0
            ? UBDDMockButton(
                title: selectedNetwork.code == null
                    ? LocaleKeys.selectNetwork.tr
                    : selectedNetwork.completedNetworkName,
                onTap: () {
                  _openNetworkSelect(
                      networks: networks,
                      selectedNetworkCode: selectedNetwork.code);
                })
            : const SizedBox();
      },
    );
  }

  void _openNetworkSelect(
      {@required String selectedNetworkCode,
      @required List<OtherNetworksConfigsAndAddresses> networks}) {
    Get.bottomSheet(
      UBButtomSheetContainer(
        title: LocaleKeys.selectNetwork.tr,
        child: Container(
          child: UBWrappedButtons(
            buttons: [
              for (var item in networks)
                WrappedButtonModel(
                  text: item.completedNetworkName,
                ),
            ],
            onButtonClick: (v) {
              Get.back();
              controller.selectedNetwork.value = networks[v];
            },
            selectedIndex: networks
                .indexWhere((element) => element.code == selectedNetworkCode),
          ),
        ),
      ),
    );
  }

  UBGreyContainer _walletLabelInput() {
    return UBGreyContainer(
      margin: const EdgeInsets.symmetric(horizontal: 12),
      child: UBBorderlessInput(
        placeholder: LocaleKeys.walletLabel.tr,
        onChange: (v) {
          controller.handleNewAddressLabelChange(v);
        },
      ),
    );
  }

  UBGreyContainer _addressInput() {
    return UBGreyContainer(
      margin: const EdgeInsets.symmetric(horizontal: 12),
      child: Row(
        children: [
          SizedBox(
            width: Get.width - 75,
            child: Obx(() {
              final value = controller.address.value;
              return ControlledTextField(
                labelText: LocaleKeys.address.tr,
                text: value,
                noBorder: true,
                onChanged: (v) {
                  controller.handleAddressChange(v);
                },
              );
            }),
          ),
          GestureDetector(
              onTap: () {
                FocusScope.of(Get.context).unfocus();
                controller.handleScanButtonTap();
              },
              child: SizedBox(
                width: 26,
                height: 26,
                child: Assets.images.qrScanIcon.svg(),
              )),
        ],
      ),
    );
  }
}
