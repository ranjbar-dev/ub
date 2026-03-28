import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../controllers/exchange_controller.dart';

import '../../../../../generated/colors.gen.dart';
import '../../../../../utils/mixins/commonConsts.dart';
import '../../../../common/components/UBCircularImage.dart';
import '../../../../common/components/UBScrollBar.dart';
import '../../../../common/components/UBSection.dart';
import '../../../../common/components/UBText.dart';
import '../../../../global/autocompleteModel.dart';

class ExchangeDropDownResults extends GetView<ExchangeController> {
  final Function(AutoCompleteItem coin) onSelect;
  final isFrom;

  ExchangeDropDownResults({Key key, this.isFrom, this.onSelect})
      : super(key: key);
  @override
  Widget build(BuildContext context) {
    final ScrollController scrollController = ScrollController();
    final coins = isFrom
        ? controller.coinsList
        : (controller.pairLocalInfo.possiblePairs ?? []);

    return Expanded(
      child: UBSection(
        hTitlePadding: 12.0,
        title: 'Coin List',
        child: Expanded(
          child: UBScrollBar(
              pullToRefreshConfig: PullToRefreshConfig(isLoading: true),
              itemCount: coins.length,
              builder: (context, index) {
                if (index != null) {
                  final coin = coins[index];
                  return GestureDetector(
                    onTap: () => onSelect(coin),
                    child: Container(
                        padding: px12,
                        decoration: const BoxDecoration(
                          border: borderBottomBlack2c,
                          color: ColorName.black,
                        ),
                        height: 49,
                        child: Row(children: [
                          UBCircularImage(
                            imageAddress: coin.image,
                            padding: const EdgeInsets.all(0),
                          ),
                          hspace8,
                          Column(
                            mainAxisAlignment: MainAxisAlignment.center,
                            crossAxisAlignment: CrossAxisAlignment.start,
                            children: [
                              UBText(
                                size: 14.0,
                                text: coin.code,
                              ),
                              vspace4,
                              UBText(
                                text: coin.desc,
                                size: 10,
                                color: ColorName.grey80,
                              ),
                            ],
                          ),
                          const Spacer(),
                        ])),
                  );
                } else
                  return Text('What?');
              },
              scrollController: scrollController),
        ),
      ),
    );
  }
}
