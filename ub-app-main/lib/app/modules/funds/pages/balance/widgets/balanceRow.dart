import 'dart:ui';

import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../../../../generated/colors.gen.dart';
import '../../../../../../utils/commonUtils.dart';
import '../../../../../../utils/extentions/basic.dart';
import '../../../../../../utils/mixins/commonConsts.dart';
import '../../../../../../utils/mixins/formatters.dart';
import '../../../../../common/components/UBCircularImage.dart';
import '../../../../../common/components/UBText.dart';
import '../../../balance_response_model_model.dart';
import '../widgets/balanceBottomSheet.dart';

class BalanceRow extends StatelessWidget with Formatter {
  final Balance balance;
  const BalanceRow({Key key, this.balance}) : super(key: key);
  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: () {
        Get.bottomSheet(
          BackdropFilter(
            child: BalanceBottomSheet(balance: balance),
            filter: ImageFilter.blur(sigmaX: 8, sigmaY: 8),
          ),
          barrierColor: ColorName.black.withOpacity(0.8),
        );
      },
      child: Container(
        height: 49,
        padding: const EdgeInsets.symmetric(horizontal: 12),
        decoration: BoxDecoration(
          border: const Border(
            bottom: const BorderSide(
              width: 1,
              color: ColorName.grey16,
            ),
          ),
        ),
        child: Row(
          children: [
            UBCircularImage(
              imageAddress: balance.image,
              padding: const EdgeInsets.all(0),
            ),
            hspace8,
            Column(
              mainAxisAlignment: MainAxisAlignment.center,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                UBText(
                  size: 14.0,
                  text: balance.code,
                ),
                vspace4,
                UBText(
                  text: balance.name,
                  size: 10,
                  color: ColorName.grey80,
                ),
              ],
            ),
            const Spacer(),
            Column(
              crossAxisAlignment: CrossAxisAlignment.end,
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                UBText(
                  size: 14.0,
                  weight: FontWeight.w700,
                  text: double.parse(balance.totalAmount) == 0
                      ? '0.00'
                      : decimalCoin(
                          value: balance.totalAmount, coinCode: balance.code),
                ),
                vspace4,
                UBText(
                  size: 10.0,
                  weight: FontWeight.w700,
                  color: ColorName.grey97,
                  text: "\$" +
                      balance.equivalentTotalAmount.currencyFormat(
                        centFormat: true,
                        removeInsignificantZeros: true,
                        toFixed: 2,
                      ),
                ),
              ],
            )
          ],
        ),
      ),
    );
  }
}
