import 'package:flutter/material.dart';
import 'UBCircularImage.dart';
import 'UBScrollBar.dart';
import 'UBSection.dart';
import 'UBText.dart';
import '../../global/autocompleteModel.dart';
import '../../../generated/colors.gen.dart';
import '../../../services/constants.dart';
import '../../../utils/mixins/commonConsts.dart';

class CoinsList extends StatelessWidget {
  final Function(AutoCompleteItem coin) onSelect;
  final coins = Constants.currencyArray();
  CoinsList({Key key, this.onSelect}) : super(key: key);
  @override
  Widget build(BuildContext context) {
    final ScrollController scrollController = ScrollController();
    return Expanded(
      child: UBSection(
        hTitlePadding: 12.0,
        title: 'Coin List',
        child: Expanded(
          child: UBScrollBar(
              itemCount: coins.length,
              builder: (context, index) {
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
              },
              scrollController: scrollController),
        ),
      ),
    );
  }
}
