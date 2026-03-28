import 'package:basic_utils/basic_utils.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../../../generated/colors.gen.dart';
import '../../../../../generated/locales.g.dart';
import '../../../../../utils/extentions/basic.dart';
import '../../../../../utils/mixins/commonConsts.dart';
import '../../../../../utils/mixins/formatters.dart';
import '../../../../common/components/UBButton.dart';
import '../../../../common/components/UBText.dart';
import '../../order_model.dart';

class OpenOrderRow extends StatelessWidget with Formatter {
  final OrderModel order;
  final bool isCancelLoading;
  final void Function(int) onCancelClick;
  const OpenOrderRow({
    Key key,
    @required this.order,
    this.isCancelLoading,
    this.onCancelClick,
  }) : super(key: key);
  @override
  Widget build(BuildContext context) {
    return Container(
      height: 54,
      padding: const EdgeInsets.only(
        left: 12,
        top: 2,
        bottom: 2,
        right: 12,
      ),
      decoration: BoxDecoration(
        color: ColorName.black,
        border: const Border(
          bottom: const BorderSide(
            color: ColorName.black2c,
            width: 1,
          ),
        ),
      ),
      child: Row(
        children: [
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            mainAxisAlignment: MainAxisAlignment.spaceEvenly,
            children: [
              UBText(
                text: order.pair.replaceAll('-', ''),
                size: 11,
                weight: FontWeight.bold,
                color: ColorName.white,
              ),
              UBText(
                text:
                    "${order.side.toUpperCase()} ${StringUtils.capitalize(order.type)}",
                size: 12,
                weight: FontWeight.w600,
                color: order.side == 'buy' ? ColorName.green : ColorName.red,
              ),
            ],
          ),
          hspace24,
          hspace24,
          hspace12,
          Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            mainAxisAlignment: MainAxisAlignment.spaceEvenly,
            children: [
              Row(
                children: [
                  UBText(
                    text: LocaleKeys.amount.tr + ': ',
                    size: 10,
                    weight: FontWeight.w600,
                    color: ColorName.greybf,
                  ),
                  UBText(
                    text: (order.amount.split(' ')[0]).currencyFormat(
                      removeInsignificantZeros: true,
                      centFormat: true,
                    ),
                    size: 11,
                    weight: FontWeight.w600,
                    color: ColorName.white,
                  ),
                ],
              ),
              Row(
                children: [
                  UBText(
                    text: LocaleKeys.price.tr + ': ',
                    size: 10,
                    weight: FontWeight.w600,
                    color: ColorName.greybf,
                  ),
                  UBText(
                    text: order.price.currencyFormat(
                      removeInsignificantZeros: true,
                      centFormat: true,
                    ),
                    size: 11,
                    weight: FontWeight.w600,
                    color: ColorName.white,
                  )
                ],
              ),
            ],
          ),
          const Spacer(),
          Column(
            mainAxisAlignment: MainAxisAlignment.spaceAround,
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              UBText(
                text: order.createdAt,
                size: 10,
                weight: FontWeight.w600,
                color: ColorName.greybf,
              ),
              Container(
                child: UBButton(
                  isLodaing: isCancelLoading,
                  height: 16.0,
                  width: 48.0,
                  smallLoading: true,
                  fontSize: 10.0,
                  borderRadius: 2.0,
                  textColor: ColorName.greybf,
                  variant: ButtonVariant.Filled,
                  buttonColor: ColorName.grey36,
                  onClick: () => {onCancelClick(order.id)},
                  text: LocaleKeys.cancel.tr,
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}
