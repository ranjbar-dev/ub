import 'package:carousel_slider/carousel_controller.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:intl/intl.dart';
import 'package:transparent_image/transparent_image.dart';

import '../../../../../generated/assets.gen.dart';
import '../../../../../generated/colors.gen.dart';
import '../../../../../services/constants.dart';
import '../../../../../utils/commonUtils.dart';
import '../../../../../utils/mixins/commonConsts.dart';
import '../../../../common/components/UBCarousel.dart';
import '../../../../common/components/UBShimmer.dart';
import '../../../../common/components/UBText.dart';
import '../../controllers/home_controller.dart';
import '../../news_model.dart';
import 'HomaPageTitle.dart';

class LastNewsCards extends GetView<HomeController> {
  const LastNewsCards({Key key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    CarouselController carouselController = CarouselController();

    return Obx(() {
      // ignore: invalid_use_of_protected_member
      final news = controller.latestNews.value;
      final hasError = !controller.isLoadingNews.value && news.isEmpty;
      return hasError
          ? const UBText(text: 'something went wrong :(')
          : Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    const HomePageTitle(text: 'Top News'),
                    fill,
                    GestureDetector(
                      onTap: () {
                        launchURL('${Constants.webLandingAddress}/news');
                      },
                      child: const UBText(text: 'All'),
                    )
                  ],
                ),
                news.isEmpty
                    ? UBShimmer(
                        height: 85.0,
                        width: Get.width - 24,
                      )
                    : Stack(
                        children: [
                          Container(
                            transform:
                                Matrix4.translationValues(0.0, -4.0, 0.0),
                            child: Container(
                              height: 85.0,
                              padding: const EdgeInsets.only(left: 6.0),
                              decoration: const BoxDecoration(
                                borderRadius: rounded7,
                                color: ColorName.black2c,
                              ),
                              child: UBCarousel(
                                //disabledrag: true,
                                controller: carouselController,
                                scrollDirection: Axis.horizontal,
                                showNavigationArrows: false,
                                height: 95.0,
                                fraction: 1,
                                autoplayInterval: 5000,
                                onChange: (i) {},
                                autoPlay: true,
                                infiniteScroll: true,
                                slides: [
                                  for (var item in news)
                                    CarouselCard(
                                      item: item,
                                      key: ValueKey(
                                        item.id.toString(),
                                      ),
                                    )
                                ],
                              ),
                            ),
                          ),
                          Positioned(
                            bottom: 12.0,
                            right: 12.0,
                            child: Assets.images.flame.svg(),
                          )
                        ],
                      ),
              ],
            );
    });
  }
}

class CarouselCard extends StatelessWidget {
  final NewsModel item;
  const CarouselCard({Key key, this.item}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    String title = item.title;
    if (title.length > 70) {
      title = title.substring(0, 70) + '...';
    }
    return GestureDetector(
      onTap: () {
        launchURL(Constants.newsAddress + item.slug);
        //print(item.title.length);
        //Get.to(() => WebViewPageView(
        //      url: Constants.newsAddress + item.id.toString(),
        //      title: item.title,
        //    ));
      },
      child: Container(
        color: Colors.transparent,
        child: Row(
          children: [
            Container(
              height: 72.0,
              width: 72.0,
              clipBehavior: Clip.antiAlias,
              decoration: const BoxDecoration(
                borderRadius: rounded7,
              ),
              child: FadeInImage.memoryNetwork(
                fit: BoxFit.cover,
                placeholder: kTransparentImage,
                image:
                    Constants.cmsAddress + item.mainImage.formats.thumbnail.url,
                fadeInDuration: const Duration(
                  milliseconds: 200,
                ),
              ),
            ),
            Padding(
              padding:
                  const EdgeInsets.only(left: 12.0, top: 12.0, bottom: 6.0),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Container(
                    width: Get.width - 130.0,
                    child: UBText(
                      text: title,
                      lineHeight: 1.3,
                      color: ColorName.white,
                      size: 13.0,
                    ),
                  ),
                  UBText(
                    text: DateFormat().format(DateTime.parse(item.date)),
                    size: 9.0,
                    color: ColorName.grey97,
                  ),
                ],
              ),
            )
          ],
        ),
      ),
    );
  }
}
