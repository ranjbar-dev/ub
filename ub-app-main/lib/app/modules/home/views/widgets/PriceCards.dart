import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../../common/components/UBNetworkIcon.dart';
import '../../../../common/components/UBShimmer.dart';
import '../../../../common/components/UBText.dart';
import '../../../../common/custom/sparkline/src/sparkline.dart';
//import 'package:unitedbit/app/common/custom/sparkline/src/sparkline.dart';
import 'package:unitedbit/app/modules/home/controllers/home_controller.dart';
import 'package:unitedbit/generated/assets.gen.dart';
import 'package:unitedbit/generated/colors.gen.dart';
import 'package:unitedbit/utils/mixins/commonConsts.dart';
import 'package:supercharged/supercharged.dart';
import 'package:unitedbit/utils/extentions/basic.dart';

class PriceCards extends GetView<HomeController> {
  const PriceCards({Key key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    final cardWidth = 142.0;
    return Container(
      height: 140.0,
      child: Obx(() {
        // ignore: invalid_use_of_protected_member
        final pairList = controller.sparkLinePairs.value;
        final hasError =
            !controller.isLoadingSparkLine.value && pairList.isEmpty;
        return hasError
            ? UBText(text: 'something went wrong :(')
            : pairList.isEmpty
                ? UBShimmer(
                    width: Get.width - 24,
                    height: 140.0,
                  )
                : ListView.builder(
                    scrollDirection: Axis.horizontal,
                    itemCount: pairList.length,
                    itemBuilder: (ctx, index) {
                      final data = pairList[index];
                      bool isRaised;
                      if (data.trendData.last.change != null) {
                        isRaised = !(data.trendData.last.change.contains('-'));
                      }
                      final lastData = data.trendData.last;
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

                      return GestureDetector(
                        onTap: () => controller.handlePairClick(
                            data.dependentCode + '-' + data.basisCode),
                        child: Container(
                          margin:
                              const EdgeInsets.only(right: 12.0, bottom: 5.0),
                          clipBehavior: Clip.antiAlias,
                          decoration: BoxDecoration(
                              // color: ColorName.black2c,
                              borderRadius: rounded_big,
                              border: Border.all(
                                width: 1.0,
                                color: '#3E3E3E'.toColor(),
                              )),
                          height: 140.0,
                          width: cardWidth,
                          child: Stack(
                            children: [
                              Positioned(
                                bottom: 10.0,
                                child: SizedBox(
                                  height: 80.0,
                                  width: cardWidth + 13.0,
                                  child: Sparkline(
                                    lineWidth: 1,
                                    data: newList,
                                    //data: chartData,
                                    fillGradient: LinearGradient(
                                      begin: Alignment.topCenter,
                                      end: Alignment.bottomCenter,
                                      colors: [
                                        '#444D67'.toColor(),
                                        '#000000'.toColor()
                                      ],
                                    ),
                                  ),
                                ),
                              ),
                              Column(
                                crossAxisAlignment: CrossAxisAlignment.center,
                                children: [
                                  vspace12,
                                  Row(
                                    children: [
                                      fill,
                                      UBNetworkIcon(
                                          imageAddress:
                                              data.secondImage ?? data.image),
                                      hspace4,
                                      UBText(
                                        text: data.dependentCode +
                                            '/' +
                                            data.basisCode,
                                        size: 12.0,
                                        color: ColorName.white,
                                        weight: FontWeight.w700,
                                      ),
                                      hspace12,
                                      fill
                                    ],
                                  ),
                                  fill,
                                  Row(
                                    mainAxisAlignment: MainAxisAlignment.center,
                                    children: [
                                      if (isRaised == null || isRaised == true)
                                        Assets.images.triangleUp.svg()
                                      else
                                        Assets.images.triangleDown.svg(),
                                      hspace2,
                                      UBText(
                                        text: (lastData.change ?? ' 0.000 %') &
                                            '%',
                                        color: ColorName.white,
                                        weight: FontWeight.w700,
                                        size: 9.0,
                                      ),
                                    ],
                                  ),
                                  vspace4,
                                  UBText(
                                    text: lastData.price.currencyFormat(
                                        removeInsignificantZeros: true,
                                        centFormat: true),
                                    color: isRaised == null
                                        ? ColorName.green
                                        : isRaised
                                            ? ColorName.green
                                            : ColorName.red,
                                    weight: FontWeight.w700,
                                    size: 15.0,
                                  ),
                                  vspace8
                                ],
                              )
                            ],
                          ),
                        ),
                      );
                    },
                  );
      }),
    );
  }
}
