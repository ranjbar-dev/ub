import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../../generated/assets.gen.dart';
import '../../../../generated/colors.gen.dart';
import '../../../../utils/pairAndCurrencyUtils.dart';
import '../../../common/components/UBText.dart';
import '../../../common/components/pageContainer.dart';
import '../../../global/autocompleteModel.dart';
import '../controllers/exchange_controller.dart';
import '../models/pairLocalInfoModel.dart';
import 'widgets/exchange_body.dart';

class ExchangeView extends GetView<ExchangeController> {
  @override
  Widget build(BuildContext context) {
    controller.pairLocalInfo = PairLocalInfoModel(
        activePairID: 1.obs,
        activePairName: 'BTC-USDT'.obs,
        pairPrecision: 8.obs,
        type: 'sell'.obs,
        basisCoin: AutoCompleteItem(
          name: 'BTC-USDT',
          code: 'BTC',
          desc: 'Bitcoin',
          image: PairAndCurrencyUtils.findCoinImageByCode('BTC'),
        ).obs,
        dependantCoin: AutoCompleteItem(
          name: 'BTC-USDT',
          code: 'USDT',
          desc: 'Tether',
          image: PairAndCurrencyUtils.findCoinImageByCode('USDT'),
        ).obs,
        basisBalance: 0.0000.obs,
        possiblePairs: [
          AutoCompleteItem(
            name: 'BTC-USDT',
            code: 'USDT',
            desc: 'Tether',
            image: PairAndCurrencyUtils.findCoinImageByCode('USDT'),
          )
        ].obs,
        dependentBalance: 0.0000.obs);
    return PageContainer(
      activeBottomNavIndex: 3,
      child: Column(
        children: [
          Container(
            height: 60.0,
            padding: const EdgeInsets.symmetric(horizontal: 12.0),
            child: Row(
              children: [
                Container(
                  width: 120.0,
                  child: Hero(
                    tag: 'logo',
                    child: Assets.images.logoSvg.svg(),
                  ),
                ),
                const Spacer(),
                GestureDetector(
                  onTap: () => Get.toNamed("/order-history",
                      arguments: {'fullScreen': true, 'isFromExchange': true}),
                  child: Container(
                    height: 24,
                    width: 92,
                    decoration: BoxDecoration(
                        borderRadius: BorderRadius.circular(8.0),
                        color: ColorName.black2c),
                    child: Padding(
                      padding: const EdgeInsets.symmetric(horizontal: 12),
                      child: Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          UBText(
                            text: 'History',
                            color: ColorName.greyd8,
                            size: 14,
                          ),
                          Assets.images.exchangeHistory.svg()
                        ],
                      ),
                    ),
                  ),
                )
              ],
            ),
          ),
          Expanded(
            child: ExchangeBody(),
          ),
        ],
      ),
    );
  }
}
