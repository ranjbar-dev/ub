import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:unitedbit/app/modules/funds/pages/autoExchange/controllers/auto_exchange_controller.dart';

import '../../../../../../generated/colors.gen.dart';
import '../../../../../../utils/mixins/commonConsts.dart';
import '../../../../../../utils/mixins/formatters.dart';
import '../../../../../generated/assets.gen.dart';
import '../../../../../utils/mixins/popups.dart';
import '../../../../common/components/UBCircularImage.dart';
import '../../../../common/components/UBText.dart';
import '../../balance_response_model_model.dart';

class AutoExchangeRow extends GetView<AutoExchangeController>
    with Formatter, Popups {
  final Balance balance;
  const AutoExchangeRow({Key key, this.balance}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return InkWell(
      onTap: () {
        controller.setCurrentBalance(balance);
        openAutoExchangePopup(balance: balance);
      },
      child: Container(
        height: 49,
        //padding: const EdgeInsets.symmetric(horizontal: 12),
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
            Container(
              height: 24,
              width: 65,
              decoration: BoxDecoration(
                  borderRadius: BorderRadius.circular(6),
                  color: ColorName.black1c),
              child: Center(
                child: balance.autoExchangeCode == null ||
                        balance.autoExchangeCode == ""
                    ? UBText(
                        text: 'Auto Off',
                        color: ColorName.grey97,
                        size: 12,
                        weight: FontWeight.w600,
                      )
                    : Padding(
                        padding: const EdgeInsets.symmetric(horizontal: 10),
                        child: Row(
                          children: [
                            UBText(
                              text: balance.autoExchangeCode,
                              size: 12,
                              color: ColorName.white,
                            ),
                            Spacer(),
                            Assets.images.autoExchangeOn.svg(),
                          ],
                        ),
                      ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
