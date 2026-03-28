import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../../../../generated/colors.gen.dart';
import '../../../../../../services/constants.dart';
import '../../../../../../utils/commonUtils.dart';
import '../../../../../../utils/extentions/basic.dart';
import '../../../../../../utils/mixins/commonConsts.dart';
import '../../../../../../utils/mixins/toast.dart';
import '../../../../../../utils/throttle.dart';
import '../../../../../common/components/UBButton.dart';
import '../../../../../common/components/UBCircularImage.dart';
import '../../../../../common/components/UBText.dart';
import '../../../../../global/autocompleteModel.dart';
import '../../../balance_response_model_model.dart';
import '../../../controllers/funds_controller.dart';
import '../../deposits/controllers/deposits_controller.dart';
import '../../deposits/views/depositDetails.dart';
import '../../withdrawals/controllers/withdrawals_controller.dart';

final thr = new Throttling(duration: const Duration(milliseconds: 4000));

class BalanceBottomSheet extends GetView<FundsController> with Toaster {
  BalanceBottomSheet({
    Key key,
    @required this.balance,
  }) : super(key: key);

  final Balance balance;
  final coins = Constants.currencyArray();

  @override
  Widget build(BuildContext context) {
    var currentCoin =
        coins.lastWhere((element) => element.code == balance.code);
    DepositsController depositsController = Get.find();
    WithdrawalsController withdawalsController = Get.find();
    return Container(
      height: /*MediaQuery.of(context).size.height * 0.43*/ 325,
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
                    child: Row(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        UBCircularImage(
                          imageAddress: balance.image,
                          padding: const EdgeInsets.all(0),
                        ),
                        hspace8,
                        UBText(
                          size: 15.0,
                          text: balance.code,
                        ),
                        hspace8,
                        UBText(
                          text: balance.name,
                          size: 11,
                          color: ColorName.grey97,
                        ),
                      ],
                    ),
                  ),
                  Positioned(
                    right: 16,
                    top: 8,
                    child: InkWell(
                      onTap: () => Get.back(),
                      child: Icon(
                        Icons.close,
                        color: ColorName.grey80,
                        size: 16,
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
          Padding(
            padding: const EdgeInsets.only(top: 12.0, left: 8.0, right: 0.8),
            child: Container(
              height: 85,
              child: Column(
                children: [
                  Align(
                    alignment: AlignmentDirectional.centerStart,
                    child: UBText(
                      text: "Total",
                      color: ColorName.greybf,
                      size: 13.0,
                    ),
                  ),
                  vspace10,
                  Align(
                    alignment: AlignmentDirectional.centerStart,
                    child: UBText(
                      size: 24.0,
                      color: ColorName.white,
                      weight: FontWeight.w700,
                      text: double.parse(balance.totalAmount) == 0
                          ? '0.00'
                          : decimalCoin(
                              value: balance.totalAmount,
                              coinCode: balance.code),
                    ),
                  ),
                  vspace10,
                  Align(
                    alignment: AlignmentDirectional.centerStart,
                    child: UBText(
                      size: 13.0,
                      weight: FontWeight.w700,
                      color: ColorName.greybf,
                      text: "\$" +
                          balance.equivalentTotalAmount.currencyFormat(
                            centFormat: true,
                            removeInsignificantZeros: true,
                            toFixed: 2,
                          ),
                    ),
                  ),
                ],
              ),
            ),
          ),
          Divider(
            color: ColorName.grey80,
            thickness: 0.5,
          ),
          IntrinsicHeight(
            child: Container(
              height: 90,
              child: Stack(
                children: [
                  Row(
                    children: <Widget>[
                      Expanded(
                        flex: 49,
                        child: Padding(
                          padding: const EdgeInsets.symmetric(
                              horizontal: 8.0, vertical: 10.0),
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              UBText(
                                text: "Available Balance",
                                color: ColorName.greybf,
                                size: 12.0,
                              ),
                              vspace10,
                              UBText(
                                size: 16.0,
                                weight: FontWeight.w700,
                                text: double.parse(balance.availableAmount) == 0
                                    ? '0.00'
                                    : decimalCoin(
                                        value: balance.totalAmount,
                                        coinCode: balance.code),
                                color: ColorName.white,
                              ),
                              vspace10,
                              UBText(
                                size: 13.0,
                                weight: FontWeight.w700,
                                color: ColorName.greybf,
                                text: "\$" +
                                    balance.equivalentAvailableAmount
                                        .currencyFormat(
                                      centFormat: true,
                                      removeInsignificantZeros: true,
                                      toFixed: 2,
                                    ),
                              ),
                            ],
                          ),
                        ),
                      ),
                      Expanded(
                        flex: 2,
                        child: VerticalDivider(
                          color: ColorName.grey80,
                          width: 0,
                        ),
                      ),
                      Expanded(
                        flex: 49,
                        child: Padding(
                          padding: const EdgeInsets.symmetric(
                              horizontal: 8.0, vertical: 10.0),
                          child: Column(
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              UBText(
                                text: "In Order",
                                color: ColorName.greybf,
                                size: 12.0,
                              ),
                              vspace10,
                              UBText(
                                size: 16.0,
                                weight: FontWeight.w700,
                                color: ColorName.white,
                                text: double.parse(balance.inOrderAmount) == 0
                                    ? '0.00'
                                    : decimalCoin(
                                        value: balance.totalAmount,
                                        coinCode: balance.code),
                              ),
                              vspace10,
                              UBText(
                                size: 13.0,
                                weight: FontWeight.w700,
                                color: ColorName.greybf,
                                text: "\$" +
                                    balance.equivalentInOrderAmount
                                        .currencyFormat(
                                      centFormat: true,
                                      removeInsignificantZeros: true,
                                      toFixed: 2,
                                    ),
                              ),
                            ],
                          ),
                        ),
                      )
                    ],
                  ),
                ],
              ),
            ),
          ),
          Divider(
            color: ColorName.grey80,
            thickness: 0.5,
          ),
          Container(
            height: 50,
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                UBButton(
                  height: 32.0,
                  width: (Get.width / 2) - 24.0,
                  onClick: () {
                    if (controller.isUserVerified.value == true) {
                      depositsController.handleCoinSelected(currentCoin);
                      _openDepositDetailsPopup(currentCoin);
                    } else {
                      toastToVerifyEmail();
                    }
                  },
                  text: 'Deposit ' + balance.code,
                ),
                hspace12,
                UBButton(
                  height: 32.0,
                  width: (Get.width / 2) - 24.0,
                  buttonColor: ColorName.black1c,
                  borderColor: ColorName.black1c,
                  variant: ButtonVariant.Filled,
                  onClick: () {
                    if (controller.isUserVerified.value == true) {
                      withdawalsController.handleCoinSelected(currentCoin);
                    } else {
                      toastToVerifyEmail();
                    }
                  },
                  text: 'Withdraw ' + balance.code,
                ),
              ],
            ),
          )
        ],
      ),
    );
  }

  toastToVerifyEmail() {
    thr.throttle(() {
      toastWarning(
          "we've sent you a verification email, please verify your email first");
    });
  }

  _openDepositDetailsPopup(AutoCompleteItem coin) {
    Get.dialog(
      Container(
        margin: const EdgeInsets.all(12.0),
        decoration: BoxDecoration(
          borderRadius: rounded_big,
          border: Border.all(color: ColorName.black2c, width: 1),
        ),
        width: Get.width,
        child: DepostDetailsView(coin: coin),
      ),
    );
  }
}
