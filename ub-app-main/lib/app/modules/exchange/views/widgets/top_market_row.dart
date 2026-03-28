import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';
import 'package:get/get.dart';
import 'package:unitedbit/app/common/components/UBShimmer.dart';

import '../../../../../generated/assets.gen.dart';
import '../../../../../generated/colors.gen.dart';
import '../../../../../utils/extentions/basic.dart';
import '../../../../../utils/mixins/commonConsts.dart';
import '../../../../../utils/pairAndCurrencyUtils.dart';
import '../../../../common/components/UBCircularImage.dart';
import '../../../../common/components/UBText.dart';
import '../../../../common/custom/sparkline/src/sparkline.dart';
import '../../../home/home_page_pair_price_model.dart';
import '../../controllers/exchange_controller.dart';

final green = ColorName.green;
final red = ColorName.red;

class TopMarketsRow extends GetView<ExchangeController> {
  final HomePagePairPriceModel data;

  const TopMarketsRow({
    Key key,
    this.data,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final lastData = data.trendData.last;
    bool isRaised;
    if (data.trendData.last.change != null) {
      isRaised = !(data.trendData.last.change.contains('-'));
    }
    final chartData = data.trendData
        .map(
          (e) => double.parse(e.price),
        )
        .toList();
    final List<double> newList = [];
    for (var i = 0; i < chartData.length; i++) {
      if (i % 9 == 0) {
        newList.add(chartData[i]);
      }
    }
    return isRaised == null
        ? ClipRRect(
            borderRadius: BorderRadius.all(
              Radius.circular(12),
            ),
            child: UBShimmer(
              width: MediaQuery.of(context).size.width,
              height: 50,
            ),
          )
        : GestureDetector(
            onTap: () => controller.handlePairSelect(data),
            child: Padding(
              padding: const EdgeInsets.symmetric(vertical: 4),
              child: Container(
                height: 50,
                decoration: BoxDecoration(
                    borderRadius: const BorderRadius.all(
                      const Radius.circular(12),
                    ),
                    color: ColorName.black2c),
                child: Padding(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  child: Row(
                    //mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      Expanded(
                        flex: 35,
                        child: Align(
                          alignment: Alignment.centerLeft,
                          child: Row(
                            mainAxisAlignment: MainAxisAlignment.start,
                            children: [
                              UBCircularImage(
                                imageAddress: data.image != null
                                    ? data.image
                                    : PairAndCurrencyUtils.findCoinImageByCode(
                                        data.pairName.split('-').first),
                              ),
                              UBText(
                                text: data.pairName
                                    .replaceFirst(RegExp('-'), '/'),
                                size: 13,
                                color: ColorName.white,
                                weight: FontWeight.bold,
                              ),
                            ],
                          ),
                        ),
                      ),
                      Expanded(
                        flex: 20,
                        child: Center(
                          child: Row(
                            mainAxisAlignment: MainAxisAlignment.center,
                            children: [
                              !isRaised
                                  ? new RotationTransition(
                                      turns:
                                          new AlwaysStoppedAnimation(180 / 360),
                                      child: Assets.images.polygonImage.svg(
                                          color: ColorName.red,
                                          width: 9,
                                          height: 9),
                                    )
                                  : Assets.images.polygonImage
                                      .svg(width: 9, height: 9),
                              hspace2,
                              RichText(
                                text: TextSpan(
                                  text: lastData.change,
                                  style: TextStyle(
                                    color: ColorName.white,
                                    fontWeight: FontWeight.w600,
                                    fontSize: 12,
                                  ),
                                ),
                              ),
                            ],
                          ),
                        ),
                      ),
                      Expanded(
                        flex: 25,
                        child: Center(
                          child: Padding(
                            padding: const EdgeInsets.only(
                                top: 10, bottom: 10, left: 10),
                            child: Sparkline(
                              lineColor:
                                  !isRaised ? ColorName.red : ColorName.green,
                              lineWidth: 1,
                              data: newList,
                              fillGradient: LinearGradient(
                                begin: Alignment.topCenter,
                                end: Alignment.bottomCenter,
                                colors: [
                                  !isRaised
                                      ? ColorName.red.withOpacity(0.0)
                                      : ColorName.green.withOpacity(0.0),
                                  ColorName.black2c
                                ],
                              ),
                            ),
                          ),
                        ),
                      ),
                      Expanded(
                        flex: 30,
                        child: Align(
                          alignment: Alignment.centerRight,
                          child: RichText(
                            text: TextSpan(
                              text: lastData.price.currencyFormat(
                                  removeInsignificantZeros: true,
                                  centFormat: true),
                              style: TextStyle(
                                  color: !isRaised ? red : green,
                                  fontWeight: FontWeight.bold,
                                  fontSize: 16),
                            ),
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ),
          );
  }
}
