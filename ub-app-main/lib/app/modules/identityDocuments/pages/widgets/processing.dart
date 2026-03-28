import 'package:flutter/material.dart';
import 'package:transparent_image/transparent_image.dart';
import '../../../../common/components/UBDottedBorder.dart';
import '../../../../common/components/UBRoundedButton.dart';
import '../../../../common/components/UBText.dart';
import '../../../../common/custom/rflutter_alert/rflutter_alert.dart';
import '../../user_profile_model.dart';
import '../../../../../generated/assets.gen.dart';
import '../../../../../generated/colors.gen.dart';
import '../../../../../generated/locales.g.dart';
import 'package:get/get.dart';

class Processing extends StatelessWidget {
  const Processing({
    Key key,
    @required this.image,
  }) : super(key: key);
  final UserProfileImages image;
  @override
  Widget build(BuildContext context) {
    return UBDottedBorder(
      child: Stack(
        children: [
          Container(
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(12),
              color: ColorName.green.withOpacity(0.15),
            ),
            child: Container(
              decoration: BoxDecoration(
                borderRadius: BorderRadius.circular(
                  11,
                ),
              ),
              clipBehavior: Clip.antiAlias,
              width: double.infinity,
              height: double.infinity,
              child: FadeInImage.memoryNetwork(
                fit: BoxFit.cover,
                placeholder: kTransparentImage,
                image: image.image,
                fadeInDuration: const Duration(
                  milliseconds: 300,
                ),
              ),
            ),
          ),
          Container(
            decoration: BoxDecoration(
              color: ColorName.primaryBlue.withOpacity(0.8),
              borderRadius: BorderRadius.circular(
                11,
              ),
            ),
            clipBehavior: Clip.antiAlias,
            child: Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  UBText(
                    text: LocaleKeys.documentReceived.tr,
                    color: ColorName.white,
                    // weight: FontWeight.bold,
                  ),
                  UBText(
                      text: LocaleKeys.willBeReviewed.tr,
                      color: ColorName.white),
                ],
              ),
            ),
          ),
          Positioned(
            right: 12,
            bottom: 12,
            child: UBRoundButton(
              child: Assets.images.expand.svg(),
              color: Colors.transparent,
              onClick: () {
                Alert(
                  style: AlertStyle(
                    animationType: AnimationType.grow,
                  ),
                  context: context,
                  content: Container(
                    child: InteractiveViewer(
                      panEnabled: true, // Set it to false to prevent panning.
                      boundaryMargin: const EdgeInsets.all(80),
                      minScale: 0.5,
                      maxScale: 4,
                      child: FadeInImage.memoryNetwork(
                        fit: BoxFit.cover,
                        placeholder: kTransparentImage,
                        image: image.image,
                        fadeInDuration: const Duration(
                          milliseconds: 300,
                        ),
                      ),
                    ),
                  ),
                ).show();
              },
            ),
          )
        ],
      ),
    );
  }
}
