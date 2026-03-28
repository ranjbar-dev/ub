import 'package:flutter/material.dart';
import 'package:basic_utils/basic_utils.dart';
import 'package:get/get.dart';
import '../../../../common/components/UBShimmer.dart';
import '../../order_model.dart';
import '../../../../../generated/assets.gen.dart';
import '../../../../../generated/colors.gen.dart';
import '../../../../../generated/locales.g.dart';
import '../../../../../utils/extentions/basic.dart';
import '../../../../../utils/mixins/commonConsts.dart';

final firstColWidth = (Get.width - 24.0) / 3.0;
final rowHeight = 53.0;

class OrderHistoryRow extends StatelessWidget {
  final bool isLoading;
  final OrderModel order;
  final String Function(String) formatter;
  final void Function(OrderModel data) onDetailsClick;
  const OrderHistoryRow(
      {Key key,
      @required this.order,
      this.onDetailsClick,
      this.isLoading = false,
      this.formatter})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    final canceledColor = order.status == 'canceled' ? ColorName.grey80 : null;
    final isCancelled = canceledColor != null;
    return GestureDetector(
      onTap: () {
        if (!isCancelled) {
          onDetailsClick(order);
        }
      },
      child: Stack(
        children: [
          Container(
            height: rowHeight,
            padding: const EdgeInsets.only(right: 12, left: 12, top: 4),
            decoration: const BoxDecoration(
              color: ColorName.black,
              border: const Border(
                bottom: const BorderSide(
                  color: ColorName.black2c,
                  width: 1,
                ),
              ),
            ),
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                Row(
                  children: [
                    SizedBox(
                      width: firstColWidth,
                      child: RichText(
                        text: TextSpan(
                          children: [
                            TextSpan(
                                text: order.pair.replaceAll('-', ''),
                                style: TextStyle(
                                    fontSize: 13,
                                    fontWeight: FontWeight.w600,
                                    color: isCancelled
                                        ? canceledColor
                                        : ColorName.white))
                          ],
                        ),
                      ),
                    ),
                    RichText(
                      text: TextSpan(
                        children: [
                          TextSpan(
                            text: LocaleKeys.amount.tr + ': ',
                            style: TextStyle(
                              fontSize: 10,
                              fontWeight: FontWeight.w600,
                              color: isCancelled
                                  ? canceledColor
                                  : ColorName.greybf,
                            ),
                          ),
                          TextSpan(
                            text: formatter(order.amount.split(' ')[0]),
                            style: TextStyle(
                              fontSize: 11,
                              fontWeight: FontWeight.w600,
                              color:
                                  isCancelled ? canceledColor : ColorName.white,
                            ),
                          ),
                        ],
                      ),
                    ),
                    const Spacer(),
                    if (isCancelled)
                      Container(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 8,
                          vertical: 2,
                        ),
                        transform: Matrix4.translationValues(0.0, 1.5, 0.0),
                        decoration: const BoxDecoration(
                            color: ColorName.grey36,
                            borderRadius: const BorderRadius.all(
                              Radius.circular(2.0),
                            )),
                        child: Text(
                          'Cancelled',
                          style: const TextStyle(
                            color: ColorName.greybf,
                            fontSize: 10,
                            fontWeight: FontWeight.w600,
                          ),
                        ),
                      )
                    else
                      Container(
                        padding: const EdgeInsets.symmetric(
                          horizontal: 8,
                          vertical: 2,
                        ),
                        transform: Matrix4.translationValues(0.0, 1.5, 0.0),
                        decoration: const BoxDecoration(
                            color: ColorName.grey36,
                            borderRadius: const BorderRadius.all(
                              Radius.circular(2.0),
                            )),
                        child: Text(
                          'Filled',
                          style: const TextStyle(
                              color: ColorName.green,
                              fontSize: 10,
                              fontWeight: FontWeight.w600),
                        ),
                      )
                  ],
                ),
                Row(
                  children: [
                    SizedBox(
                      width: firstColWidth,
                      child: RichText(
                        text: TextSpan(
                          children: [
                            TextSpan(
                              text:
                                  "${order.side.toUpperCase()} ${StringUtils.capitalize(order.type)}",
                              style: TextStyle(
                                fontSize: 11,
                                fontWeight: FontWeight.w600,
                                color: canceledColor != null
                                    ? canceledColor
                                    : order.side == 'buy'
                                        ? ColorName.green
                                        : ColorName.red,
                              ),
                            ),
                          ],
                        ),
                      ),
                    ),
                    RichText(
                      text: TextSpan(
                        children: [
                          TextSpan(
                            text: LocaleKeys.price.tr + ': ',
                            style: TextStyle(
                              fontSize: 10,
                              fontWeight: FontWeight.w600,
                              color: isCancelled
                                  ? canceledColor
                                  : ColorName.greybf,
                            ),
                          ),
                          TextSpan(
                            text: order.price == ''
                                ? 'Market'
                                : order.price.currencyFormat(
                                    removeInsignificantZeros: true,
                                    centFormat: true),
                            style: TextStyle(
                              fontSize: 11,
                              fontWeight: FontWeight.w600,
                              color:
                                  isCancelled ? canceledColor : ColorName.white,
                            ),
                          ),
                        ],
                      ),
                    ),
                    const Spacer(),
                    Padding(
                      padding: const EdgeInsets.only(
                        top: 6,
                        bottom: 6,
                        left: 6,
                      ),
                      child: RowDate(
                        date: order.createdAt,
                        isCancelled: isCancelled,
                        canceledColor: canceledColor,
                      ),
                    ),
                  ],
                )
              ],
            ),
          ),
          if (isLoading)
            UBShimmer(
              height: rowHeight,
              width: Get.width,
              opacity: 0.7,
            ),
        ],
      ),
    );
  }
}

class RowDate extends StatelessWidget {
  final String date;
  final bool isCancelled;
  final Color canceledColor;

  const RowDate({
    Key key,
    this.date,
    this.isCancelled,
    this.canceledColor,
  }) : super(key: key);
  @override
  Widget build(BuildContext context) {
    final splitted = date.split(' ');
    return Row(
      children: [
        RichText(
          text: TextSpan(
            text: splitted[1] + ' ' + splitted[0],
            style: TextStyle(
                fontSize: 10,
                fontWeight: FontWeight.w600,
                color: isCancelled ? canceledColor : ColorName.greybf),
          ),
        ),
        if (isCancelled != true) hspace4,
        if (isCancelled != true) Assets.images.keyboardArrowRight.svg()
      ],
    );
  }
}
