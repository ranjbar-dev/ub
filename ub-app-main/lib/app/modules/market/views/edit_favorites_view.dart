import 'package:flutter/material.dart';

import 'package:get/get.dart';
import '../../../common/components/appbarTextTitle.dart';
import '../../../common/components/pageContainer.dart';
import '../../../global/currency_pairs_model.dart';
import '../controllers/market_controller.dart';
import '../widgets/editFavsHead.dart';
import '../../../../generated/assets.gen.dart';
import '../../../../generated/colors.gen.dart';
import '../../../../utils/mixins/commonConsts.dart';

class EditFavoritesView extends GetView<MarketController> {
  @override
  Widget build(BuildContext context) {
    // final scrollController = controller.favsScrollController;

    return PageContainer(
      appbarTitle: AppBarTextTitle(title: 'Edit Favorites'),
      child: Column(
        children: [
          EditFavsHead(),
          Expanded(
            child: Obx(
              () {
                final pairsList = controller.orderedPairs.isNotEmpty
                    ? controller.orderedPairs
                    : controller.pairs;
                return RawScrollbar(
                  controller: controller.favsScrollController,
                  fadeDuration: 200.milliseconds,
                  radius: const Radius.circular(12),
                  thumbColor: ColorName.grey23,
                  thickness: 4.0,
                  child: Theme(
                    data: ThemeData(
                      canvasColor: ColorName.black2c,
                    ),
                    child: ReorderableListView.builder(
                      itemCount: pairsList.length,
                      scrollController: controller.favsScrollController,
                      onReorder: controller.handleFavPairsReorder,
                      itemBuilder: (BuildContext context, int index) {
                        return EditOrdersRow(
                          onFavCheckboxClick: (pair) =>
                              controller.handleFavCheckboxClick(pair: pair),
                          onMovePairToTop: (pair) =>
                              controller.movePairToTopOfOrderedList(pair: pair),
                          pair: pairsList[index],
                          index: index,
                          key: Key('$index'),
                        );
                      },
                    ),
                  ),
                );
              },
            ),
          ),
        ],
      ),
    );
  }
}

class EditOrdersRow extends StatelessWidget {
  final Pairs pair;
  final Function(Pairs) onFavCheckboxClick;
  final Function(Pairs) onMovePairToTop;
  final int index;

  const EditOrdersRow({
    Key key,
    this.pair,
    this.onFavCheckboxClick,
    this.index,
    this.onMovePairToTop,
  }) : super(key: key);
  @override
  Widget build(BuildContext context) {
    final splitted = pair.pairName.split('-');
    final pairName1 = splitted[0];
    final pairName2 = splitted[1];
    return Container(
      key: key,
      padding: px12,
      height: 40,
      decoration: rowDecoration,
      child: Row(
        children: [
          GestureDetector(
            onTap: () {
              onFavCheckboxClick(pair);
            },
            child: pair.isFavorite == true
                ? Assets.images.filledCheckbox.svg()
                : Assets.images.emptyCheckbox.svg(),
          ),
          hspace4,
          RichText(
            text: TextSpan(
              children: <TextSpan>[
                TextSpan(
                  text: pairName1,
                  style: whiteBold14,
                ),
                TextSpan(
                  text: '/',
                  style: whiteBold10,
                ),
                TextSpan(
                  text: pairName2,
                  style: whiteBold10,
                ),
              ],
            ),
          ),
          const Spacer(),
          GestureDetector(
            onTap: () {
              onMovePairToTop(pair);
            },
            child: SizedBox(
              width: 24,
              height: 24,
              child: Assets.images.toTop.svg(),
            ),
          ),
          hspace20,
          ReorderableDragStartListener(
            index: index,
            child: SizedBox(
              width: 24,
              height: 24,
              child: Assets.images.threeLine.svg(),
            ),
          ),
        ],
      ),
    );
  }
}
