// ignore_for_file: invalid_use_of_protected_member

import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:unitedbit/app/common/components/CenterUBLoading.dart';

import '../../../../../generated/colors.gen.dart';
import '../../../../common/components/UBText.dart';
import '../../controllers/exchange_controller.dart';
import 'top_market_row.dart';

class TopMarkets extends GetView<ExchangeController> {
  const TopMarkets({Key key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Container(
          padding: const EdgeInsets.only(top: 20, bottom: 6),
          color: ColorName.black,
          child: Align(
            alignment: AlignmentDirectional.centerStart,
            child: UBText(
                color: ColorName.white,
                text: 'Top Markets',
                align: TextAlign.start,
                size: 19),
          ),
        ),
        Obx(
          () {
            final popularPairs = controller.sparkLinePairs.value;
            return controller.isLoadingSparkLine.value
                ? Container(
                    height: MediaQuery.of(context).size.height / 4,
                    child: CenterUbLoading())
                : ListView.builder(
                    physics: const NeverScrollableScrollPhysics(),
                    shrinkWrap: true,
                    scrollDirection: Axis.vertical,
                    itemCount: popularPairs.length,
                    itemBuilder: (ctx, index) {
                      final data = popularPairs[index];
                      // bool isRaised;
                      // if (data.trendData.last.change != null) {
                      //   isRaised = !(data.trendData.last.change.contains('-'));
                      // }

                      return TopMarketsRow(
                        data: data,
                      );
                    },
                  );
          },
        )
      ],
    );
  }
}
