import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../../common/components/UBShimmer.dart';
import '../../../../common/components/UBTooltip.dart';
import '../../models/price_model.dart';
import '../../../../../utils/commonUtils.dart';
import '../../../../../utils/mixins/formatters.dart';
import '../../../../../utils/mixins/popups.dart';
import '../../../../common/components/UBDropdown.dart';
import '../../../../common/components/UBText.dart';
import '../../controllers/trade_controller.dart';
import '../../../../../generated/colors.gen.dart';
import 'favoriteStar.dart';
import '../../../../../utils/extentions/basic.dart';

double lastPrice = 0.0;

var lastValidPrice;

class TradeHead extends GetView<TradeController> with Formatter, Popups {
  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        Container(
            padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 4),
            // height: 42,
            decoration: const BoxDecoration(
              color: ColorName.black2c,
              border: const Border(
                bottom: const BorderSide(
                  color: ColorName.black,
                  width: 1,
                ),
              ),
            ),
            child: Column(
              children: [
                Row(
                  children: [
                    //UBButton(
                    //  onClick: () {
                    //    //controller.globalController.handleLoggedOut();
                    //    //Get.put(AfterSplashController(test: true));
                    //    final GlobalController globalController = Get.find();
                    //    globalController.getVersion();
                    //  },
                    //  text: 'lo',
                    //),
                    //  //     {
                    //  //   openVerificationPopup(
                    //  //     onSubmit: () {},
                    //  //     need2fa: true,
                    //  //     needEmailCode: true,
                    //  //     isNewDevice: true,
                    //  //   )
                    //  // },
                    PairAndPercent(),
                    LastPrice(),
                    const Spacer(),
                    ToggleCharts(),
                    FavButton(),
                  ],
                ),
                Obx(() {
                  final isOrderBook =
                      controller.activeChart.value == TradeTopCharts.OrderBook;
                  final currentPrice = controller.currentPairPrice.value;
                  if (currentPrice.high != null) {
                    lastValidPrice = currentPrice;
                  }
                  if (isOrderBook && lastValidPrice != null) {
                    return HeaderWidget(
                      currentPrice: currentPrice,
                    );
                  }
                  return const SizedBox();
                })
              ],
            )),
        Obx(() {
          final currentPrice = controller.currentPairPrice.value;
          if ((currentPrice.high == null)) {
            return UBShimmer(
              width: Get.width,
              height: 46,
              opacity: 0.3,
            );
          }
          return const SizedBox();
        })
      ],
    );
  }
}

class PairAndPercent extends GetView<TradeController> {
  @override
  Widget build(BuildContext context) {
    return Container(
      alignment: Alignment.topLeft,
      margin: const EdgeInsets.only(right: 12),
      child: Obx(
        () {
          final ddOptions =
              // ignore: invalid_use_of_protected_member
              controller.globalController.currencyPairsArray.value;
          final pairName = controller.currentPairName.value;
          return Column(
            mainAxisAlignment: MainAxisAlignment.start,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              UBDropDown(
                // expanded: true,
                // ignore: invalid_use_of_protected_member
                dense: true,
                options: ddOptions,
                value: pairName,
                onChange: (v) {
                  controller.handlePairChange(v);
                  return;
                },
              ),
              if (controller.currentPairPrice.value.high != null)
                Percent(percent: controller.currentPairPrice.value.percentage)
              else
                const SizedBox(
                  height: 12,
                )
            ],
          );
        },
      ),
    );
  }
}

class FavButton extends GetView<TradeController> {
  @override
  Widget build(BuildContext context) {
    return Obx(
      () {
        final currentPair = controller.currentPairName.value;
        return FavoriteStar(pairName: currentPair);
      },
    );
  }
}

class LastPrice extends GetView<TradeController> with Formatter {
  @override
  Widget build(BuildContext context) {
    return Container(
      alignment: Alignment.topLeft,
      child: Obx(
        () {
          if (controller.currentPairPrice.value.high != null) {
            final priceNumber =
                double.parse(controller.currentPairPrice.value.price);
            final isOhlc = controller.activeChart.value == TradeTopCharts.OHLC;
            final isRising = priceNumber >= lastPrice;
            final color = isRising ? ColorName.green : ColorName.red;
            final price = controller.currentPairPrice.value.price;
            final String fractionCoinName =
                controller.currentPairName.value.split('-')[1];
            lastPrice = priceNumber;
            return Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Padding(
                  padding: const EdgeInsets.only(bottom: 2.0, top: 2.0),
                  child: SizedBox(
                    height: 18,
                    child: RichText(
                      text: TextSpan(
                        text: decimalCoin(
                            value: price, coinCode: fractionCoinName),
                        style: TextStyle(
                          fontWeight: FontWeight.bold,
                          fontSize: 18,
                          height: isRising ? 1.1 : 1.14,
                          color: color,
                        ),
                      ),
                    ),
                  ),
                ),
                if (isOhlc)
                  Volume(
                    volume: controller.currentPairPrice.value.volume,
                  )
                else
                  const SizedBox(
                    height: 14,
                  )
              ],
            );
          }
          return const SizedBox();
        },
      ),
    );
  }
}

class ToggleCharts extends GetView<TradeController> {
  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.only(top: 3),
      alignment: Alignment.topCenter,
      color: ColorName.black2c,
      child: GestureDetector(
        onTap: controller.toggleCharts,
        child: Obx(() {
          final isOrderBook =
              controller.activeChart.value == TradeTopCharts.OrderBook;
          return UBToolTip(
            message: isOrderBook ? 'Show Trade chart' : 'Show Order Book',
            child: SizedBox(
              width: 30,
              height: 30,
              child: Icon(
                Icons.list_rounded,
                color: isOrderBook ? ColorName.primaryBlue : ColorName.greybf,
              ),
            ),
          );
        }),
      ),
    );
  }
}

class Volume extends StatelessWidget {
  final String volume;

  const Volume({Key key, this.volume}) : super(key: key);
  @override
  Widget build(BuildContext context) {
    return UBText(
      text:
          volume.currencyFormat(removeInsignificantZeros: true, compact: true),
    );
  }
}

class Percent extends StatelessWidget {
  final String percent;

  const Percent({Key key, this.percent}) : super(key: key);
  @override
  Widget build(BuildContext context) {
    final number = double.parse(percent);
    final color = number > 0 ? ColorName.green : ColorName.red;
    final prefix = number > 0 ? '+' : '';
    return Row(
      children: [
        UBText(
          text: 'C: ',
          size: 11,
        ),
        UBText(
          size: 11,
          text: "$prefix$percent%",
          color: color,
        )
      ],
    );
  }
}

class HeaderWidget extends StatelessWidget with Formatter {
  final PriceModel currentPrice;

  const HeaderWidget({Key key, this.currentPrice}) : super(key: key);
  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        Container(
          child: Row(
            children: [
              SizedBox(
                child: ColumnText(
                  title: 'Low',
                  value: currentPrice.low ?? lastValidPrice.low,
                  formatter: currencyFormatter,
                ),
              ),
              const SizedBox(width: 8),
              ColumnText(
                title: 'High',
                value: currentPrice.high ?? lastValidPrice.low,
                formatter: currencyFormatter,
              ),
              const SizedBox(width: 8),
              ColumnText(
                title: 'Vol(USDT)',
                value: currentPrice.volume ?? lastValidPrice.volume,
                formatter: currencyFormatter,
              )
            ],
          ),
        ),
        if (currentPrice.high == null)
          UBShimmer(
            width: Get.width,
            height: 28,
            opacity: 0.3,
          )
      ],
    );
  }
}

class ColumnText extends StatelessWidget {
  final String title, value;
  final String Function(String) formatter;
  final Color greyColor;

  const ColumnText(
      {Key key, this.title, this.value, this.formatter, this.greyColor})
      : super(key: key);
  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.only(top: 4),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          RichText(
            text: TextSpan(
              text: title,
              style: const TextStyle(
                fontSize: 9,
                color: ColorName.greybf,
              ),
            ),
          ),
          RichText(
            text: TextSpan(
              text: formatter(double.parse(value).toStringAsFixed(2)),
              style: const TextStyle(
                fontSize: 11,
                color: ColorName.white,
              ),
            ),
          ),
        ],
      ),
    );
  }
}
