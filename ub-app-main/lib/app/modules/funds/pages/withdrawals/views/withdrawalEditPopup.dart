import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../../../common/components/CenterUBLoading.dart';
import '../../../../../common/components/UBButton.dart';
import '../../../../../common/components/UBCountUp.dart';
import '../../../../../common/components/UBGreyContainer.dart';
import '../../../../../common/components/UBLi.dart';
import '../../../../../common/components/UBPercentSelect.dart';
import '../../../../../common/components/UBShimmer.dart';
import '../../../../../common/components/UBText.dart';
import '../../../../../common/components/UBToastOnTap.dart';
import '../../../../../common/components/UBTooltip.dart';
import '../../../../../common/components/UBTwoPartText.dart';
import '../../../../../common/components/UBWrappedButtons.dart';
import '../../../../../common/components/controlledInput.dart';
import '../controllers/withdrawals_controller.dart';
import '../../../../../../generated/assets.gen.dart';
import '../../../../../../generated/colors.gen.dart';
import '../../../../../../generated/locales.g.dart';
import '../../../../../../utils/commonUtils.dart';
import '../../../../../../utils/mixins/commonConsts.dart';
import '../../../../../../utils/mixins/formatters.dart';
import '../../../../../../utils/mixins/toast.dart';
import '../../../../../../utils/extentions/basic.dart';

class WithdrawalEditPopup extends GetView<WithdrawalsController>
    with Toaster, Formatter {
  @override
  Widget build(BuildContext context) {
    return Container(
        height: 580,
        width: Get.width,
        decoration: const BoxDecoration(
            color: ColorName.black2c, borderRadius: roundedTop_big),
        child: Stack(
          children: [
            Column(
              children: [
                vspace24,
                _title(),
                vspace8,
                _networkSelector(),
                vspace12,
                _addressInput(),
                vspace12,
                _withdrawWarning(),
                _amountInputAndBelowColumns(),
                const Spacer(),
                _submitButton(),
                vspace24
              ],
            ),
            Obx(() {
              final isLoading = controller.showLoadingInsidePopup.value;
              if (isLoading == true) {
                return Container(
                  color: ColorName.black.withOpacity(0.4),
                  child: CenterUbLoading(),
                );
              }
              return const SizedBox();
            })
          ],
        ));
  }

  _title() {
    return Obx(() {
      final selectedCoin = controller.selectedCoin.value;
      final selectedNetwork = controller.selectedNetwork.value;
      final data = controller.withdrawAndDepositData.value;
      return Container(
        padding: px12,
        width: Get.width,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            UBText(
              size: 15.0,
              text: 'Withdraw' & selectedCoin.name,
            ),
            if (selectedNetwork.completedNetworkName != '') vspace4,
            if (selectedNetwork.completedNetworkName != '')
              UBText(
                  size: 15.0,
                  text:
                      'on' & '${selectedNetwork.completedNetworkName} network')
            else
              UBText(
                  size: 15.0,
                  text: 'on' & '${data.completedNetworkName} network'),
          ],
        ),
      );
    });
  }

  UBGreyContainer _addressInput() {
    return UBGreyContainer(
      color: ColorName.black,
      margin: const EdgeInsets.symmetric(horizontal: 12),
      child: Row(
        children: [
          SizedBox(
            width: Get.width - 118,
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
          UBToolTip(
            message: 'Select from saved Addresses',
            preferBelow: false,
            child: GestureDetector(
              onTap: () {
                FocusScope.of(Get.context).unfocus();
                controller.handleSearchAddressClick();
              },
              child: Padding(
                padding: const EdgeInsets.symmetric(horizontal: 4),
                child: Assets.images.addressManagement.svg(),
              ),
            ),
          ),
          UBToolTip(
            message: 'Scan or select from gallery',
            preferBelow: false,
            child: GestureDetector(
              onTap: () {
                FocusScope.of(Get.context).unfocus();
                _scan();
              },
              child: Container(
                color: Colors.transparent,
                padding: const EdgeInsets.symmetric(horizontal: 4),
                child: SizedBox(
                  width: 26,
                  height: 26,
                  child: Assets.images.qrScanIcon.svg(),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }

  Obx _amountInputAndBelowColumns() {
    return Obx(() {
      final data = controller.withdrawAndDepositData.value;
      final isLoading = controller.isLoadingWithdrawAndDepositData.value;
      final selectedCoin = controller.selectedCoin.value;
      if (data.balance == null) {
        return const SizedBox();
      }
      final balance = data.balance.availableAmount;
      String fee = controller.fee.value;

      final amount = controller.amount.value;
      var youGet;
      if (amount != '' && fee != null) {
        youGet = double.parse(amount.replaceAll(',', '')) - double.parse(fee);
      }
      final hasBalance = double.parse(balance) > 0.00;
      return (data.balance == null)
          ? CenterUbLoading()
          : Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                UBToastOnTap(
                  active: !hasBalance,
                  onTap: controller.handleEmptyBalanceClick,
                  child: UBGreyContainer(
                    color: ColorName.black,
                    margin: EdgeInsets.only(left: 12, right: 12),
                    child: Row(
                      children: [
                        Expanded(
                          child: Container(
                            color: ColorName.black,
                            child: Obx(
                              () => ControlledTextField(
                                type: TextInputType.number,
                                labelText: LocaleKeys.amount.tr,
                                text: controller.amount.value,
                                isCurrencyInput: true,
                                noBorder: true,
                                onChanged: (v) {
                                  controller.handleAmountChange(v);
                                },
                              ),
                            ),
                          ),
                        ),
                        const Spacer(),
                        Container(
                          child: Row(
                            children: [
                              if (isLoading)
                                UBShimmer(
                                  width: 40.0,
                                  height: 20.0,
                                )
                              else
                                UBCountup(
                                  begin: 0,
                                  end: double.parse(balance),
                                  suffix: " ${selectedCoin.code}",
                                ),
                              hspace4,
                              // if (hasBalance)
                              //   UBButton(
                              //     onClick: () {
                              //       controller.handleAllAmountClick();
                              //     },
                              //     text: LocaleKeys.all.tr,
                              //     height: 20,
                              //     width: 40,
                              //     fontSize: 11,
                              //     textColor: ColorName.primaryBlue,
                              //     buttonColor: ColorName.grey23,
                              //   )
                            ],
                          ),
                        )
                      ],
                    ),
                  ),
                )
                //if (!hasBalance)
                //  GestureDetector(
                //    onTap: controller.handleEmptyBalanceClick,
                //    child: UBGreyContainer(
                //      color: ColorName.black,
                //      margin: EdgeInsets.only(left: 12, right: 12),
                //      child: Row(
                //        children: [
                //          UBText(text: LocaleKeys.amount.tr),
                //          const Spacer(),
                //          if (isLoading)
                //            UBShimmer(
                //              width: 40.0,
                //              height: 20.0,
                //            )
                //          else
                //            UBCountup(
                //              begin: 0,
                //              end: double.parse(balance),
                //              suffix: " ${selectedCoin.code}",
                //            ),
                //        ],
                //      ),
                //    ),
                //  )
                //else
                ,
                if (hasBalance)
                  Container(
                    padding: const EdgeInsets.only(left: 12.0),
                    width: 190.0,
                    height: 30,
                    child: UBPercentSelect(
                      onPercentClick: controller.handlePercentClick,
                      selectedIndex: controller.selectedPercentIndex.value,
                      numberOfSegments: controller.numberOfPercentSegments,
                    ),
                  ),
                vspace24,
                UBTwoPartText(
                  title: LocaleKeys.minimumWithdrawal.tr,
                  valueWidget: isLoading
                      ? UBShimmer()
                      : UBText(text: data.balance.minimumWithdraw),
                  size: 11,
                  padding: EdgeInsets.only(left: 12, right: 12, bottom: 18),
                ),
                if (fee != null)
                  UBTwoPartText(
                    title: LocaleKeys.transactionFee.tr,
                    valueWidget: isLoading
                        ? UBShimmer()
                        : UBText(
                            text: decimalCoin(value: fee & selectedCoin.code)),
                  ),
                if (youGet != null)
                  UBTwoPartText(
                    title: LocaleKeys.youWillGet.tr,
                    valueWidget: isLoading
                        ? UBShimmer()
                        : UBText(
                            text: decimalCoin(
                                value: youGet.toString() & selectedCoin.code),
                          ),
                  ),
                if (data.withdrawComments != null)
                  for (var item in data.withdrawComments)
                    UBLi(
                      text: item,
                      dotColor: ColorName.orange,
                    )
              ],
            );
    });
  }

  Padding _withdrawWarning() {
    return Padding(
      padding: EdgeInsets.only(left: 12, right: 12, bottom: 48),
      child: Row(
        children: [
          Icon(
            Icons.warning_amber_rounded,
            color: ColorName.orange,
            size: 24,
          ),
          const SizedBox(
            width: 4,
          ),
          UBText(
            wrapped: true,
            color: ColorName.grey80,
            text: LocaleKeys.dontWithdrawDirectly.tr,
          )
        ],
      ),
    );
  }

  Padding _submitButton() {
    return Padding(
      padding: const EdgeInsets.only(bottom: 12, left: 12, right: 12),
      child: Obx(() {
        final isLoading = controller.isLoadingWithdrawAndDepositData.value;
        final selectedCoin = controller.selectedCoin.value;
        final address = controller.address.value;
        final amount = controller.amount.value;
        final isSubmitting = controller.isSubmitting.value;
        return UBButton(
          isLodaing: isSubmitting,
          onClick: () {
            controller.handleSubmitClick();
          },
          text: LocaleKeys.submit.tr,
          disabled: (isLoading ||
              selectedCoin.id == null ||
              address.length < 10 ||
              amount == ''),
        );
      }),
    );
  }

  Obx _networkSelector() {
    return Obx(
      () {
        final data = controller.withdrawAndDepositData.value;
        final isLoading = controller.isLoadingWithdrawAndDepositData.value;
        return (data.networksConfigsAndAddresses == null || isLoading == true)
            ? const SizedBox()
            : Container(
                padding: const EdgeInsets.only(right: 12.0, left: 6.0),
                child: Row(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Container(
                      alignment: Alignment.centerLeft,
                      padding: const EdgeInsets.symmetric(horizontal: 8),
                      height: 36,
                      child: UBText(
                        text: LocaleKeys.selectNetwork.tr,
                        size: 11,
                        color: ColorName.greybf,
                      ),
                    ),
                    Expanded(
                      child: const SizedBox(),
                    ),
                    UBWrappedButtons(
                      buttonBackground: ColorName.black1c,
                      selectedIndex: controller.selectedNetworkIndex.value,
                      buttons: [
                        for (var item in data.networksConfigsAndAddresses)
                          WrappedButtonModel(
                            text: '${item.code}',
                          ),
                      ],
                      onButtonClick: (i) {
                        controller.handleNetworkChange(i);
                        return;
                      },
                    )
                  ],
                ),
              );
      },
    );
  }

  void _scan() {
    controller.handleScanButtonTap();
  }
}
