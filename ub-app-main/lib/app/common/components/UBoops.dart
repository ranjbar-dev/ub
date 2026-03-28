import 'package:flutter/material.dart';
import 'UBColumnAnimator.dart';
import 'UBText.dart';
import '../../../generated/assets.gen.dart';
import '../../../generated/colors.gen.dart';
import '../../../generated/locales.g.dart';
import 'package:get/get.dart';
import '../../../utils/mixins/commonConsts.dart';

enum OopsVariant {
  OpenOrderOops,
  OrderHistoryOops,
  NoFilterResultsOops,
  AddressManagementOops,
  NoFavPairsOops,
  NoTransactionHistory,
  NoBalancesOops,
  SearchOops,
  ErrorOops
}

class UBoops extends StatelessWidget {
  final OopsVariant variant;
  const UBoops({Key key, this.variant}) : super(key: key);
  @override
  Widget build(BuildContext context) {
    switch (variant) {
      case OopsVariant.OpenOrderOops:
        return OopsAnimator(
          children: [
            Assets.images.noOpenOrderIcon.svg(),
            vspace10,
            oopsText(
              text: LocaleKeys.youDontHaveAnyOpenOrder.tr,
            )
          ],
        );
        break;
      case OopsVariant.NoBalancesOops:
        return OopsAnimator(
          children: [
            Assets.images.noOpenOrderIcon.svg(),
            vspace12,
            oopsText(
              text: "You don't have any balance yet",
            )
          ],
        );
        break;
      case OopsVariant.SearchOops:
        return OopsAnimator(
          children: [
            Assets.images.noSearchResults.svg(),
            vspace12,
            oopsText(
              text: "No Results!",
            )
          ],
        );
        break;
      case OopsVariant.NoTransactionHistory:
        return OopsAnimator(children: [
          Assets.images.noOpenOrderIcon.svg(),
          vspace12,
          oopsText(
            text: 'You dont have any transaction yet',
          )
        ]);
        break;
      case OopsVariant.ErrorOops:
        return OopsAnimator(
          children: [
            SizedBox(
              height: 80,
              width: 80,
              child: Assets.images.close.image(),
            ),
            oopsText(
              text: 'Sorry, Something Went Wrong :(',
            )
          ],
        );
        break;
      case OopsVariant.NoFavPairsOops:
        return OopsAnimator(children: [
          SizedBox(
            width: 64,
            height: 64,
            child: Assets.images.filledStar.svg(color: const Color(0xff3C3C47)),
          ),
          oopsText(
            text: 'Your favorite list is empty',
          )
        ]);
        break;
      case OopsVariant.OrderHistoryOops:
        return OopsAnimator(children: [
          Assets.images.noOpenOrderIcon.svg(),
          vspace10,
          oopsText(
            text: LocaleKeys.emptyOrderHistory.tr,
          )
        ]);
        break;
      case OopsVariant.NoFilterResultsOops:
        return OopsAnimator(children: [
          Assets.images.noFilterResult.svg(),
          vspace12,
          oopsText(
            text: 'There are no results for your filters',
          )
        ]);
        break;
      case OopsVariant.AddressManagementOops:
        return OopsAnimator(children: [
          Assets.images.addressesOpsImage.svg(),
          const SizedBox(
            height: 24,
          ),
          oopsText(
            text: LocaleKeys.noWithdrawAddressLine1.tr,
          ),
        ]);
        break;

      default:
        return Container(
          child: Center(
            child: Container(
              color: Colors.black87,
            ),
          ),
        );
    }
  }
}

oopsText({String text}) {
  return UBText(
    size: 14.0,
    color: ColorName.grey80,
    text: text,
  );
}

class OopsAnimator extends StatelessWidget {
  final List<Widget> children;
  const OopsAnimator({Key key, @required this.children}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      child: Center(
        child: UBColumnScaleAnimator(
          children: children,
        ),
      ),
    );
  }
}
