import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:supercharged/supercharged.dart';

import '../../../../../generated/colors.gen.dart';
import '../../../../../generated/locales.g.dart';
import '../../../../../utils/extentions/basic.dart';
import '../../../../../utils/mixins/formatters.dart';
import '../../../../common/components/CenterUBLoading.dart';
import '../../../../common/components/UBRectangle.dart';
import '../../controllers/trade_controller.dart';

final rowHeight = 17.2;

class OrderBook extends GetView<TradeController> {
  @override
  Widget build(BuildContext context) {
    final maxWidth = (Get.width - 24) / 2;
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 12),
      color: ColorName.black2c,
      child: Column(
        children: [
          Container(
            height: 35,
            child: Obx(() {
              final pairName = controller.currentPairName.value;
              return Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  OrderBookHeaderText(
                      text:
                          "${LocaleKeys.amount.tr}(${pairName.split('-')[0]})"),
                  OrderBookHeaderText(
                      text:
                          "${LocaleKeys.price.tr}(${pairName.split('-')[1]})"),
                  OrderBookHeaderText(
                      text:
                          "${LocaleKeys.amount.tr}(${pairName.split('-')[0]})")
                ],
              );
            }),
          ),
          Expanded(child: Obx(
            () {
              final data = controller.orderBookData;
              return data['asks'] == null
                  ? Container(
                      color: ColorName.black2c,
                      child: CenterUbLoading(),
                    )
                  : ChartRender(
                      data: data,
                      maxWidth: maxWidth,
                      onAskClick: controller.handleAskClick,
                      onBidClick: controller.handleBidClick,
                    );
            },
          )),
        ],
      ),
    );
  }
}

class ChartRender extends StatelessWidget {
  final dynamic data;
  final double maxWidth;
  final Function onAskClick;
  final Function onBidClick;
  const ChartRender(
      {Key key, this.data, this.maxWidth, this.onAskClick, this.onBidClick})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    final List asks = data['asks'];
    List bids = data['bids'];
    final asksMax = asks[asks.length - 1]['sum'];
    final bidsMax = bids[bids.length - 1]['sum'];
    return Column(
      children: [
        for (var index = 0; index < asks.length; index++)
          SizedBox(
            height: rowHeight,
            child: Row(
              children: [
                Expanded(
                  child: GestureDetector(
                    onTap: () => onAskClick(asks[index]),
                    child: OrderBookElement(
                      data: asks[index],
                      max: asksMax,
                      maxWidth: maxWidth,
                    ),
                  ),
                ),
                const SizedBox(
                  width: 4,
                ),
                Expanded(
                  child: GestureDetector(
                    onTap: () => onBidClick(bids[index]),
                    child: OrderBookElement(
                      data: bids[index],
                      max: bidsMax,
                      maxWidth: maxWidth,
                    ),
                  ),
                ),
              ],
            ),
          )
      ],
    );
  }
}

class OrderBookElement extends StatelessWidget with Formatter {
  final data;
  final String max;
  final double maxWidth;

/*
  {
      "price": "57920.780000",
      "amount": "0.002041",
      "value": "118.216312",
      "percentage": "0.00",
      "sum": "0.002041",
      "type": "ask"
    }*/

  const OrderBookElement({
    Key key,
    @required this.data,
    @required this.max,
    this.maxWidth,
  }) : super(key: key);
  @override
  Widget build(BuildContext context) {
    final isAsk = data['type'] == 'ask';
    final String sum = data['sum'];
    final String whiteValue = data['amount'];
    final String price = data['price'];
    final percent = (double.parse(sum) / double.parse(max)) * 100;

    return SizedBox(
      child: Stack(
        fit: StackFit.expand,
        children: [
          Positioned(
              right: isAsk ? 0 : null,
              left: isAsk ? null : 0,
              child: UBRectangle(
                color: isAsk ? '#3A4D45'.toColor() : '#341D1D'.toColor(),
                width: (percent * maxWidth) / 100,
                height: rowHeight,
              )),
          Row(
            textDirection: isAsk ? TextDirection.ltr : TextDirection.rtl,
            children: [
              ChartText(
                text: whiteValue.currencyFormat(formatSmall: true),
                color: ColorName.white,
              ),
              const Spacer(),
              isAsk
                  ? ChartText(
                      text: price.currencyFormat(formatSmall: true),
                      color: ColorName.green,
                    )
                  : ChartText(
                      text: price.currencyFormat(formatSmall: true),
                      color: ColorName.red,
                    )
            ],
          ),
        ],
      ),
    );
  }
}

class ChartText extends StatelessWidget {
  final String text;
  final Color color;

  const ChartText({Key key, this.text, this.color}) : super(key: key);
  @override
  Widget build(BuildContext context) {
    return RichText(
      text: TextSpan(
        text: text,
        style: TextStyle(
          color: color,
          fontSize: 11.0,
          fontWeight: FontWeight.w600,
        ),
      ),
    );
  }
}

class OrderBookHeaderText extends StatelessWidget {
  final String text;

  const OrderBookHeaderText({Key key, this.text}) : super(key: key);
  @override
  Widget build(BuildContext context) {
    return Text(
      text,
      style: const TextStyle(color: ColorName.grey80, fontSize: 11.0),
    );
  }
}
