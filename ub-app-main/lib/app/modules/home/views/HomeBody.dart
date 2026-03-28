import 'package:flutter/material.dart';
import 'package:get/state_manager.dart';
import 'package:pull_to_refresh/pull_to_refresh.dart';

import '../../../../utils/logger.dart';
import '../../../../utils/mixins/commonConsts.dart';
import '../../../common/components/UBColumnAnimator.dart';
import '../../../common/components/UBPullToRefresh.dart';
import '../../funds/widgets/DepositWithdrawButtons.dart';
import '../controllers/home_controller.dart';
import 'widgets/HomaPageTitle.dart';
import 'widgets/LastNewsCard.dart';
import 'widgets/PopularPairs.dart';
import 'widgets/PriceCards.dart';
import 'widgets/banners.dart';

class HomeBody extends GetView<HomeController> {
  @override
  Widget build(BuildContext context) {
    final RefreshController _refreshController = RefreshController();
    return UBPullToRefresh(
      showUpdateText: true,
      releaseToUpdateText: 'Release to update',
      afterUpdateText: 'Updated',
      beforeUpdateText: 'Pull to update',
      updatingText: 'Updating...',
      withUpdateDate: true,
      controller: _refreshController,
      onRefresh: () async {
        try {
          await controller.refreshHomePage();
        } catch (e) {
          log.e(e.toString());
        } finally {
          _refreshController.refreshCompleted();
        }
      },
      child: SingleChildScrollView(
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 12.0),
          child: UBColumnSlideAnimator(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              //vspace12,
              HomePageBanners(),
              vspace12,
              Obx(
                () => DepositWithdrawButtons(
                  isUserVerified: controller.isUserVerified.value,
                ),
              ),
              HomePageTitle(text: 'Top Markets'),
              PriceCards(),
              vspace8,
              LastNewsCards(
                key: ValueKey('latest news'),
              ),
              HomePageTitle(text: 'Popular Pairs'),
              PopularPairs(),
              SizedBox(
                height: 70,
              )
            ],
          ),
        ),
      ),
    );
  }
}
