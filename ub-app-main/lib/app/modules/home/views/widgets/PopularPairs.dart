import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../../../generated/colors.gen.dart';
import '../../../../../utils/mixins/commonConsts.dart';
import '../../../market/widgets/marketPriceRow.dart';
import '../../controllers/home_controller.dart';

class PopularPairs extends GetView<HomeController> {
  const PopularPairs({Key key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    double firstColWidth = (Get.width / 3) + 12.0;
    return Column(
      children: [
        Container(
          padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
          color: ColorName.grey16,
          child: Row(
            children: [
              SizedBox(
                width: firstColWidth,
                child: Text(
                  'Coin',
                  style: grey80Bold13,
                ),
              ),
              Text(
                'Last Price',
                style: grey80Bold13,
              ),
              const Spacer(),
              Text(
                'Chg%',
                style: grey80Bold13,
              ),
            ],
          ),
        ),
        Obx(
          () {
            // ignore: invalid_use_of_protected_member
            final popularPairs = controller.popularPairs.value;
            return Container(
              height: popularPairs.length * 41.0,
              child: ListView.builder(
                physics: const NeverScrollableScrollPhysics(),
                itemCount: popularPairs.length,
                itemBuilder: (ctx, index) {
                  final data = popularPairs[index];
                  return MarketPriceRow(
                    firstColWidth: firstColWidth,
                    onClick: (s) => controller.handlePairClick(data.pairName),
                    data: data,
                    equivalentPrice: data.formattedEquivalentPrice,
                    price: data.formattedPrice,
                    volume: data.formattedVolume,
                  );
                },
              ),
            );
          },
        )
      ],
    );
  }
}
