import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../common/components/CenterUBLoading.dart';
import '../../../common/components/UBButton.dart';
import '../../../common/components/UBScrollBar.dart';
import '../../../common/components/UBoops.dart';
import '../controllers/market_controller.dart';
import '../widgets/marketPriceRow.dart';
import '../../../routes/app_pages.dart';
import '../../../../generated/assets.gen.dart';
import '../../../../generated/colors.gen.dart';
import '../../../../utils/mixins/commonConsts.dart';
import '../../../../utils/mixins/formatters.dart';

class FavoritesList extends GetView<MarketController> with Formatter {
  @override
  Widget build(BuildContext context) {
    final ScrollController scrollController = ScrollController();
    return Obx(() {
      final pairsArray = controller.pairs;
      final favList = controller.favorites;
      if (favList.length == 0) {
        if (pairsArray.length == 0) {
          return CenterUbLoading();
        }
        return Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            UBoops(
              variant: OopsVariant.NoFavPairsOops,
            ),
            vspace12,
            SizedBox(
              width: 155.0,
              child: UBButton(
                onClick: () {
                  Get.toNamed(AppPages.EDIT_FAVORITES);
                },
                height: 32,
                buttonColor: ColorName.black2c,
                fontSize: 13.0,
                textColor: ColorName.primaryBlue,
                variant: ButtonVariant.Rounded,
                text: '  Edit Favorites list',
                endWidget: Assets.images.editWithUnderline.svg(
                  color: ColorName.primaryBlue,
                ),
              ),
            )
          ],
        );
      }
      return UBScrollBar(
          scrollController: scrollController,
          itemCount: favList.length,
          builder: (BuildContext context, int index) {
            final data = favList[index];
            return MarketPriceRow(
              onClick: controller.handlePriceRowClick,
              price: data.formattedPrice,
              volume: data.formattedVolume,
              equivalentPrice: data.formattedEquivalentPrice,
              data: data,
            );
          });
    });
  }
}
