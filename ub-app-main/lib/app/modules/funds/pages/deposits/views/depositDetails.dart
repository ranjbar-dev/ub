import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:get/get.dart';
import 'package:qr_flutter/qr_flutter.dart';
import '../../../../../common/components/CenterUBLoading.dart';
import '../../../../../common/components/UBCarousel.dart';
import '../../../../../common/components/UBLi.dart';
import '../../../../../common/components/UBSection.dart';
import '../../../../../common/components/UBText.dart';
import '../../../../../common/components/UBWarningRow.dart';
import '../../../../../common/components/UBWrappedButtons.dart';
import '../../../../../common/custom/toaster/utopic_toast.dart';
import '../../../../../global/autocompleteModel.dart';
import '../controllers/deposits_controller.dart';
import '../../../../../../generated/assets.gen.dart';
import '../../../../../../generated/colors.gen.dart';
import '../../../../../../generated/locales.g.dart';
import '../../../../../../utils/mixins/commonConsts.dart';
import '../../../../../../utils/mixins/popups.dart';
import 'package:carousel_slider/carousel_slider.dart';

class DepostDetailsView extends GetView<DepositsController> with Popups {
  final AutoCompleteItem coin;

  DepostDetailsView({this.coin});
  @override
  Widget build(BuildContext context) {
    CarouselController carouselController = CarouselController();
    return Container(
      decoration: const BoxDecoration(
        color: ColorName.black2c,
        borderRadius: rounded_big,
      ),
      child: Column(
        children: [
          vspace16,
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16.0),
            child: Row(
              children: [
                UBText(text: "Deposit ${coin.desc}"),
                const Spacer(),
                GestureDetector(
                    onTap: () => Get.back(),
                    child: Container(
                      color: Colors.transparent,
                      child: SizedBox(
                        width: 24.0,
                        height: 24.0,
                        child: Assets.images.closeIcon.svg(color: Colors.white),
                      ),
                    ))
              ],
            ),
          ),
          UBWarningRow(
            text:
                'If you have deposited, please pay attention to the text messages, site letters and emails we send to you.',
            background: ColorName.black1c,
            separatorColor: ColorName.black2c,
            padding: const EdgeInsets.symmetric(
              horizontal: 12.0,
            ),
            margin: const EdgeInsets.only(
              left: 12.0,
              right: 12.0,
              top: 24.0,
            ),
          ),
          vspace24,
          Obx(
            () {
              final data = controller.withdrawAndDepositData.value;
              final isLoading =
                  controller.isLoadingWithdrawAndDepositData.value;
              if (isLoading == true) {
                return Container(
                  height: Get.height / 2,
                  child: CenterUbLoading(),
                );
              }
              return (data.networksConfigsAndAddresses == null ||
                      isLoading == true)
                  ? const Center()
                  : Container(
                      child: UBSection(
                        title: LocaleKeys.selectNetwork.tr,
                        child: Padding(
                          padding: const EdgeInsets.symmetric(horizontal: 12.0),
                          child: UBWrappedButtons(
                            buttonBackground: ColorName.black,
                            minButtonWidth: 70.0,
                            selectedIndex:
                                controller.selectedNetworkIndex.value,
                            buttons: [
                              for (var item in data.networksConfigsAndAddresses)
                                WrappedButtonModel(
                                  text: item.code,
                                ),
                            ],
                            onButtonClick: (i) {
                              controller.handleNetworkChange(i);
                              carouselController.animateToPage(i);
                              return;
                            },
                          ),
                        ),
                      ),
                    );
            },
          ),
          vspace24,
          Obx(() {
            final data = controller.withdrawAndDepositData.value;
            final selectedNetwork = controller.selectedNetwork.value;
            final isLoading = controller.isLoadingWithdrawAndDepositData.value;
            return selectedNetwork.code == null || isLoading == true
                ? const SizedBox()
                : Container(
                    child: Column(
                      children: [
                        UBCarousel(
                            controller: carouselController,
                            slides: [
                              for (var item in data.networksConfigsAndAddresses)
                                Container(
                                  alignment: Alignment.center,
                                  width: Get.width - 150,
                                  decoration: BoxDecoration(
                                    color: ColorName.white,
                                    borderRadius: BorderRadius.circular(
                                      12,
                                    ),
                                  ),
                                  padding: const EdgeInsets.all(9),
                                  child: QrImage(
                                    data: item.address,
                                    version: QrVersions.auto,
                                    size: Get.width - 150,
                                  ),
                                ),
                            ],
                            height: Get.width - 150,
                            selectedIndex:
                                controller.selectedNetworkIndex.value,
                            onChange: (index) {
                              controller.handleNetworkChange(index);
                              return;
                            }),
                        //vspace48,
                        Container(
                          height: 36,
                          padding: px12,
                          child: Row(
                            mainAxisAlignment: MainAxisAlignment.spaceBetween,
                            children: [
                              UBText(
                                text: LocaleKeys.depositAddress.tr,
                              ),
                              GestureDetector(
                                onTap: () {
                                  Clipboard.setData(
                                    ClipboardData(
                                      text: selectedNetwork.address,
                                    ),
                                  );
                                  ToastManager().showToast(
                                    LocaleKeys.addressCopiedToClipboard.tr,
                                    type: ToastType.info,
                                    action: ToastAction(
                                      onPressed: (hideToastFn) {
                                        hideToastFn();
                                      },
                                    ),
                                  );
                                },
                                child: Container(
                                  height: 32,
                                  width: 32,
                                  decoration: const BoxDecoration(
                                    color: ColorName.grey16,
                                    borderRadius: rounded6,
                                  ),
                                  child: Assets.images.copyIcon.svg(),
                                ),
                              )
                            ],
                          ),
                        ),
                        Container(
                          padding: px12,
                          alignment: Alignment.topLeft,
                          child: UBText(
                            text: selectedNetwork.address,
                            color: ColorName.white,
                          ),
                        ),
                        if (data.depositComments != null)
                          for (var item in data.depositComments)
                            UBLi(
                              text: item,
                              dotColor: ColorName.orange,
                            )
                      ],
                    ),
                  );
          }),
          //vspace24,
          fill,
          Obx(
            () {
              final data = controller.withdrawAndDepositData.value;
              final isLoading =
                  controller.isLoadingWithdrawAndDepositData.value;
              return data.currencyExtraInfo == null || isLoading == true
                  ? const SizedBox()
                  : Container(
                      padding: const EdgeInsets.symmetric(
                        horizontal: 12,
                      ),
                      alignment: Alignment.topLeft,
                      child: Column(
                        crossAxisAlignment: CrossAxisAlignment.start,
                        children: [
                          Row(
                            children: [
                              Assets.images.warningTriangle.svg(),
                              SizedBox(
                                width: 4,
                              ),
                              Container(
                                child: Row(
                                  children: [
                                    UBText(
                                      text: controller.selectedCoin.value.name,
                                      color: ColorName.greybf,
                                    ),
                                    const SizedBox(width: 4),
                                    UBText(
                                      text: LocaleKeys.depositInfo.tr
                                          .toUpperCase(),
                                      color: ColorName.greybf,
                                    ),
                                  ],
                                ),
                              )
                            ],
                          ),
                          UBText(
                            text: data.currencyExtraInfo.description,
                            color: ColorName.grey80,
                            size: 11,
                          )
                        ],
                      ),
                    );
            },
          ),
          vspace24
        ],
      ),
    );
  }
}
