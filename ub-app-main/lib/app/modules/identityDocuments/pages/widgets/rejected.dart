import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../../../generated/assets.gen.dart';
import '../../../../../generated/colors.gen.dart';
import '../../../../../generated/locales.g.dart';
import '../../../../common/components/UBButton.dart';
import '../../../../common/components/UBDottedBorder.dart';
import '../../../../common/components/UBText.dart';
import '../../../../common/custom/rflutter_alert/src/alert.dart';
import '../../../../common/custom/rflutter_alert/src/alert_style.dart';
import '../../../../common/custom/rflutter_alert/src/constants.dart';
import '../../user_profile_model.dart';

class Rejected extends StatelessWidget {
  const Rejected({
    Key key,
    @required this.side,
    @required this.subType,
    this.image,
  }) : super(key: key);
  final String side;
  final SubTypes subType;
  final UserProfileImages image;

  @override
  Widget build(BuildContext context) {
    return UBDottedBorder(
      color: ColorName.red,
      child: Container(
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(12),
        ),
        child: Center(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              //UBText(
              //  text:
              //      "$side side of your ${StringUtils.capitalize(subType.name.replaceAll('_', ' '), allWords: true)} has been rejected!",
              //  color: ColorName.red,
              //),
              //const SizedBox(
              //  height: 12,
              //),
              Assets.images.rejectedIcon.svg(),
              const SizedBox(
                height: 12,
              ),
              UBButton(
                onClick: () {
                  _showReason(context: context, image: image);
                },
                endWidget: Assets.images.envelope.svg(),
                text: 'Reject Reason',
                variant: ButtonVariant.Rounded,
                width: 132,
                height: 24,
                buttonColor: ColorName.grey22,
                textColor: ColorName.red,
              )
            ],
          ),
        ),
      ),
    );
  }

  void _showReason({BuildContext context, UserProfileImages image}) {
    Alert(
      style: AlertStyle(
        animationType: AnimationType.grow,
      ),
      context: context,
      header: Container(
        padding: const EdgeInsets.only(top: 12, left: 8),
        child: Row(
          children: [
            Assets.images.redShield.svg(
              color: ColorName.red,
            ),
            const SizedBox(
              width: 4,
            ),
            UBText(
              text: LocaleKeys.yourDocumentIsNotVerified.tr,
              color: ColorName.red,
            ),
          ],
        ),
      ),
      content: Container(
        child: Column(
          children: [
            Container(
              alignment: Alignment.topLeft,
              child: UBText(
                text: image.rejectionReason,
              ),
            ),
          ],
        ),
      ),
    ).show();
  }
}
