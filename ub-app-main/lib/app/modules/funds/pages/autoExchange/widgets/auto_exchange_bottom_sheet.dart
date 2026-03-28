import 'package:flutter/material.dart';
import 'package:flutter_switch/flutter_switch.dart';
import 'package:get/get.dart';

import '../../../../../../generated/assets.gen.dart';
import '../../../../../../generated/colors.gen.dart';
import '../../../../../../services/constants.dart';
import '../../../../../../utils/mixins/commonConsts.dart';
import '../../../../../../utils/mixins/toast.dart';
import '../../../../../../utils/throttle.dart';
import '../../../../../common/components/UBButton.dart';
import '../../../../../common/components/UBCircularImage.dart';
import '../../../../../common/components/UBText.dart';
import '../../../../exchange/controllers/exchange_controller.dart';
import '../../../../exchange/views/widgets/exchange_drop_down.dart';
import '../../../balance_response_model_model.dart';
import '../controllers/auto_exchange_controller.dart';

final thr = new Throttling(duration: const Duration(milliseconds: 4000));

class AutoExchangeBottomSheet extends GetView<AutoExchangeController>
    with Toaster {
  AutoExchangeBottomSheet({
    Key key,
    @required this.balance,
  }) : super(key: key);

  final Balance balance;
  final coins = Constants.currencyArray();

  @override
  Widget build(BuildContext context) {
    ExchangeController exchangeController = Get.find<ExchangeController>();
    return Container(
      height: /*MediaQuery.of(context).size.height * 0.43*/ 300,
      decoration: BoxDecoration(
          borderRadius: BorderRadius.vertical(
            top: Radius.circular(20),
          ),
          color: ColorName.black2c),
      child: Column(
        children: [
          Padding(
            padding: const EdgeInsets.symmetric(vertical: 10.0),
            child: Container(
              height: 30,
              child: Stack(
                children: [
                  Center(
                    child: UBText(
                      text: 'Auto Ex',
                      color: ColorName.white,
                      size: 14,
                      weight: FontWeight.w600,
                    ),
                  ),
                  Positioned(
                    right: 16,
                    top: 3,
                    child: InkWell(
                      onTap: () {
                        Get.back();
                      },
                      child: Container(
                        height: 24,
                        width: 24,
                        //color: Colors.red,
                        child: Icon(
                          Icons.close,
                          color: ColorName.grey80,
                          size: 16,
                        ),
                      ),
                    ),
                  )
                ],
              ),
            ),
          ),
          Divider(
            color: ColorName.grey80,
            thickness: 0.5,
            height: 0,
          ),
          vspace24,
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 24),
            child: Row(
              children: [
                Assets.images.bigInfoIcon.svg(),
                hspace8,
                UBText(
                    color: ColorName.greyd8,
                    size: 12,
                    wrapped: true,
                    weight: FontWeight.w400,
                    text:
                        'Lorem ipsum dolor sit amet, consectetur adipiscing elit. Amet.')
              ],
            ),
          ),
          vspace24,
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 24),
            child: Row(
              children: [
                UBText(
                  text: 'Auto Ex',
                  size: 12,
                  color: ColorName.white,
                  weight: FontWeight.w600,
                ),
                hspace8,
                Obx(
                  () => Container(
                    width: 32,
                    height: 16,
                    child: FlutterSwitch(
                      padding: 2,
                      toggleSize: 15,
                      inactiveColor: ColorName.black1c,
                      inactiveToggleColor: ColorName.grey80,
                      activeToggleColor: ColorName.white,
                      activeColor: ColorName.primaryBlue,
                      value: controller.switchValue.value,
                      onToggle: (val) {
                        controller.toggleAutoExchange(balance: balance);
                      },
                    ),
                  ),
                ),
                hspace8,
                Obx(
                  () => UBText(
                    text: controller.switchValue.value ? 'On' : 'Off',
                    size: 12,
                    color: ColorName.greybf,
                    weight: FontWeight.w600,
                  ),
                ),
              ],
            ),
          ),
          vspace24,
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 24),
            child: Row(
              //mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Expanded(
                  flex: 46,
                  child: Stack(
                    alignment: Alignment.center,
                    children: [
                      Container(
                        height: 40,
                        decoration: BoxDecoration(
                          color: ColorName.black1c,
                          borderRadius: BorderRadius.circular(6),
                        ),
                      ),
                      Row(
                        mainAxisAlignment: MainAxisAlignment.start,
                        crossAxisAlignment: CrossAxisAlignment.center,
                        children: [
                          GetPlatform.isWeb
                              ? UBCircularImage(
                                  imageAddress: balance.image,
                                )
                              : Obx(
                                  () => controller.switchValue.value
                                      ? UBCircularImage(
                                          imageAddress: balance.image,
                                        )
                                      : ColorFiltered(
                                          colorFilter: ColorFilter.mode(
                                              ColorName.black1c,
                                              BlendMode.saturation),
                                          child: UBCircularImage(
                                            imageAddress: balance.image,
                                          ),
                                        ),
                                ),
                          UBText(
                            text: balance.code,
                            color: ColorName.grey80,
                            size: 12,
                            weight: FontWeight.w400,
                          ),
                          hspace6,
                          UBText(
                            text: '( ' + balance.name + ' )',
                            color: ColorName.grey80,
                            size: 8,
                            weight: FontWeight.w400,
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
                Expanded(
                  flex: 8,
                  child: Center(
                    child: UBText(
                      text: 'To',
                      size: 13,
                      color: ColorName.grey97,
                    ),
                  ),
                ),
                Obx(
                  () => Expanded(
                    flex: 46,
                    child: ExchangeDropDown(
                      autoExchangeCode:
                          controller.currentBalance.autoExchangeCode,
                      name: exchangeController
                          .pairLocalInfo.dependantCoin.value.code,
                      desc: exchangeController
                          .pairLocalInfo.dependantCoin.value.desc,
                      imageUrl: exchangeController
                          .pairLocalInfo.dependantCoin.value.image,
                      isFrom: false,
                      isDisabled: !controller.switchValue.value,
                    ),
                  ),
                ),
              ],
            ),
          ),
          Padding(
            padding: const EdgeInsets.all(24.0),
            child: Obx(
              () => UBButton(
                onClick: () {
                  controller.submitAutoExchange();
                },
                disabled: !controller.canSubmitAutoExchange.value,
                height: 40,
                isLodaing: controller.isSubmitLoading.value,
                text: 'Apply',
                fontSize: 18,
                borderRadius: 6.0,
              ),
            ),
          ),
        ],
      ),
    );
  }
}
