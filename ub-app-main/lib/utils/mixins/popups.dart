import 'package:basic_utils/basic_utils.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_datetime_picker/flutter_datetime_picker.dart';
import 'package:get/get.dart';
import 'package:get_storage/get_storage.dart';
import 'package:unitedbit/app/modules/funds/pages/autoExchange/widgets/auto_exchange_bottom_sheet.dart';
import 'package:url_launcher/url_launcher.dart';

import '../../app/common/components/UBBorderlessInput.dart';
import '../../app/common/components/UBButton.dart';
import '../../app/common/components/UBCircularImage.dart';
import '../../app/common/components/UBGreyContainer.dart';
import '../../app/common/components/UBHorizontalDivider.dart';
import '../../app/common/components/UBScrollColumnExpandable.dart';
import '../../app/common/components/UBText.dart';
import '../../app/common/components/countrySelectBottomSheet.dart';
import '../../app/common/components/currencySelectButtomSheet.dart';
import '../../app/common/components/pairSelectButtomSheet.dart';
import '../../app/common/custom/rflutter_alert/rflutter_alert.dart';
import '../../app/common/custom/toaster/utopic_toast.dart';
import '../../app/global/autocompleteModel.dart';
import '../../app/modules/funds/balance_response_model_model.dart';
import '../../app/modules/funds/pages/autoExchange/controllers/auto_exchange_controller.dart';
import '../../app/modules/funds/pages/transactionHistory/transaction_history_model.dart';
import '../../app/modules/orders/order_history_detail_model.dart';
import '../../app/modules/orders/order_model.dart';
import '../../app/popups/controllers/twofaPopupController.dart';
import '../../app/popups/veiws/twoFaPopupView.dart';
import '../../generated/assets.gen.dart';
import '../../generated/colors.gen.dart';
import '../../generated/locales.g.dart';
import '../../services/storageKeys.dart';
import '../commonUtils.dart';
import '../extentions/basic.dart';
import 'commonConsts.dart';

mixin Popups {
  openDatePickerPopup({Function(String date) onDateSelect}) {
    DatePicker.showDatePicker(Get.context,
        showTitleActions: true,
        minTime: DateTime(1900, 1, 1),
        maxTime: DateTime.now(),
        theme: DatePickerTheme(
          headerColor: ColorName.black2c,
          backgroundColor: ColorName.grey16,
          itemStyle: TextStyle(
            color: Colors.white,
            fontWeight: FontWeight.bold,
            fontSize: 18,
          ),
          doneStyle: TextStyle(
            color: Colors.white,
            fontSize: 16,
          ),
          cancelStyle: TextStyle(
            color: Colors.white,
            fontSize: 16,
          ),
        ), onConfirm: (date) {
      onDateSelect(date.toString());
    }, currentTime: DateTime.now(), locale: LocaleType.en);
  }

  openCoinSelectPopup(
      {@required Function(AutoCompleteItem) onCoinSelect,
      Function afterClose}) {
    Get.bottomSheet(
      SelectCurrencyBottomSheet(
        onSelect: (item) {
          onCoinSelect(item);
          Get.back();
          if (afterClose != null) {
            afterClose();
          }
          return;
        },
      ),
      isScrollControlled: true,
      ignoreSafeArea: false,
    );
  }

  openPairSelectPopup(
      {@required Function(AutoCompleteItem) onSelect, Function afterClose}) {
    Get.bottomSheet(
      SelectPairBottomSheet(
        onSelect: (item) {
          onSelect(item);
          Get.back();
          if (afterClose != null) {
            afterClose();
          }
          return;
        },
      ),
      isScrollControlled: true,
      ignoreSafeArea: false,
    );
  }

  openCountrySelect({@required Function(AutoCompleteItem) onCountrySelect}) {
    Get.bottomSheet(
      SelectCountryBottomSheet(
        onSelect: (item) {
          onCountrySelect(item);
          Get.back();
        },
      ),
      isScrollControlled: true,
      ignoreSafeArea: false,
    );
  }

  openTwofaInputPopup({@required Function(String) onSubmit}) {
    String value = '';
    Get.bottomSheet(
      Container(
        height: Get.height,
        color: ColorName.black,
        padding: const EdgeInsets.symmetric(horizontal: 12),
        child: UBScrollColumnExpandable(
          children: [
            const Spacer(),
            Container(
              child: Column(
                children: [
                  Container(
                    child: Assets.images.twofaicon.svg(),
                  ),
                  vspace24,
                  SizedBox(
                    width: 300,
                    child: UBText(
                      text: LocaleKeys.pleaseOpen2Fa.tr,
                      align: TextAlign.center,
                      lineHeight: 1.6,
                    ),
                  )
                ],
              ),
            ),
            const Spacer(),
            UBGreyContainer(
              child: UBBorderlessInput(
                type: TextInputType.number,
                onChange: (v) => {value = v},
              ),
            ),
            vspace24,
            UBButton(
                onClick: () {
                  onSubmit(value);
                  Get.back();
                },
                text: LocaleKeys.submit.tr),
            vspace24,
            SizedBox(
                width: 120,
                child: UBButton(
                  variant: ButtonVariant.TransparentBackground,
                  onClick: () {
                    Get.back();
                  },
                  text: LocaleKeys.cancel.tr,
                )),
            vspace24,
          ],
        ),
      ),
      isScrollControlled: true,
      ignoreSafeArea: false,
    );
  }

  Future openVerificationPopup({
    @required bool need2fa,
    @required bool needEmailCode,
    @required bool isNewDevice,
    @required Function onSubmit,
  }) async {
    Get.put(
      TwoFaController(),
    );
    await Get.bottomSheet(
      TwoFaPopupView(
        onSubmit: onSubmit,
        isNewDevice: isNewDevice,
        need2fa: need2fa,
        needEmailCode: needEmailCode,
      ),
      isScrollControlled: true,
      ignoreSafeArea: false,
    );
  }

  openAutoExchangePopup({Balance balance}) {
    Get.put(
      AutoExchangeController(),
    );
    // Get.put(
    //   ExchangeController(),
    // );
    Get.bottomSheet(
      AutoExchangeBottomSheet(balance: balance),
    );
  }

  openExchangeSubmitPopup(
      {double spentAmount,
      String coinCode,
      dynamic model,
      Function backTapped}) {
    Get.bottomSheet(
      Container(
        height: 450,
        margin: const EdgeInsets.symmetric(horizontal: 12),
        decoration: BoxDecoration(
          color: ColorName.black2c,
          borderRadius: const BorderRadius.only(
            topLeft: const Radius.circular(16),
            topRight: const Radius.circular(16),
          ),
        ),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.start,
          children: [
            vspace48,
            UBText(
              text: 'Successful',
              color: ColorName.green,
              size: 20,
            ),
            vspace8,
            Stack(
              alignment: AlignmentDirectional.center,
              children: [
                Container(
                  width: 65.0,
                  height: 65.0,
                  decoration: BoxDecoration(
                      color: ColorName.black3c, shape: BoxShape.circle),
                  child: Padding(
                    padding: const EdgeInsets.all(15.0),
                    child: Assets.images.exchangeSuccess.svg(),
                  ),
                ),
                Icon(
                  Icons.done,
                  color: ColorName.black2c,
                )
              ],
            ),
            vspace16,
            UBText(
              text: 'You Spent',
              size: 12,
              weight: FontWeight.bold,
              color: ColorName.grey80,
            ),
            vspace8,
            Row(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                UBText(
                  text: spentAmount.toString().currencyFormat(),
                  color: ColorName.white,
                  size: 24,
                ),
                hspace12,
                UBText(
                  text: coinCode,
                  color: ColorName.grey97,
                  size: 24,
                )
              ],
            ),
            vspace16,
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 12.0),
              child: Divider(
                color: ColorName.grey97,
              ),
            ),
            vspace16,
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 12.0),
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      UBText(
                        text: 'Converted',
                        color: ColorName.grey97,
                        size: 16,
                      ),
                      UBText(
                        text: model['amount'].toString(),
                        color: ColorName.white,
                        size: 16,
                      ),
                    ],
                  ),
                  vspace16,
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      UBText(
                        text: 'Price',
                        color: ColorName.grey97,
                        size: 16,
                      ),
                      UBText(
                        text: model['price'].toString(),
                        color: ColorName.white,
                        size: 16,
                      ),
                    ],
                  ),
                  vspace16,
                  Row(
                    mainAxisAlignment: MainAxisAlignment.spaceBetween,
                    children: [
                      UBText(
                        text: 'Inverse Price',
                        color: ColorName.grey97,
                        size: 16,
                      ),
                      UBText(
                        text: model['fee'].toString(),
                        color: ColorName.white,
                        size: 16,
                      ),
                    ],
                  ),
                ],
              ),
            ),
            Spacer(),
            Padding(
              padding: const EdgeInsets.all(12.0),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Expanded(
                    flex: 48,
                    child: UBButton(
                      onClick: backTapped,
                      text: 'Back',
                      fontSize: 15,
                      buttonColor: ColorName.black1c,
                      borderRadius: 8.0,
                    ),
                  ),
                  Expanded(
                    flex: 4,
                    child: SizedBox(
                      width: 6,
                    ),
                  ),
                  Expanded(
                    flex: 48,
                    child: UBButton(
                      onClick: () => Get.toNamed("/order-history", arguments: {
                        'fullScreen': true,
                        'isFromExchange': true
                      }),
                      text: 'View Status',
                      fontSize: 15,
                      borderRadius: 8.0,
                    ),
                  )
                ],
              ),
            )
          ],
        ),
      ),
    );
  }

  openWithdrawSubmitPopup({
    @required AutoCompleteItem coin,
    @required Function onSubmit,
    @required String youWillGetAmount,
    @required String address,
    @required String network,
    @required String amount,
    @required String transactionFee,
  }) {
    Get.bottomSheet(
      Container(
        padding: const EdgeInsets.symmetric(horizontal: 12.0),
        height: 550.0,
        decoration: const BoxDecoration(
            color: ColorName.black2c, borderRadius: roundedTop_big),
        child: Column(
          children: [
            vspace48,
            SizedBox(
              child: Center(
                child: UBText(
                  text: 'Confirm Withdraw',
                  size: 18.0,
                  weight: FontWeight.w600,
                  color: ColorName.greyd8,
                ),
              ),
            ),
            vspace12,
            UBCircularImage(
              size: 36,
              imageAddress: coin.image,
              padding: const EdgeInsets.all(0),
            ),
            vspace12,
            UBText(
              text: 'You Will Get:',
              size: 11.0,
              color: ColorName.grey80,
            ),
            vspace12,
            RichText(
              text: TextSpan(
                children: [
                  TextSpan(text: youWillGetAmount + ' ', style: whiteBold24),
                  TextSpan(text: coin.code, style: grey97Bold24),
                ],
              ),
            ),
            UBHorizontalDivider(),
            TwoPartText(
              title: 'Coin',
              value: coin.code,
            ),
            TwoPartText(
              title: 'Address',
              value: address,
            ),
            TwoPartText(
              title: 'Network',
              value: network,
            ),
            UBHorizontalDivider(),
            TwoPartText(
              title: 'Amount',
              value: amount,
            ),
            TwoPartText(
              title: 'Transaction Fee',
              value: transactionFee,
            ),
            fill,
            UBButton(
              onClick: () {
                Get.back();
                onSubmit();
              },
              text: 'Confirm',
            ),
            vspace24
          ],
        ),
      ),
      isScrollControlled: true,
      ignoreSafeArea: false,
    );
  }

  void openConfirmation({
    @required Function onConfirm,
    @required Widget titleWidget,
    double titleDistanceFromTop = -20.0,
    String cancelText = 'Cancel',
    String confirmText,
    bool autoBackAfterConfirm = true,
    Color cancelTextColor = ColorName.white,
    Color confirmTextColor = ColorName.red,
  }) {
    Alert(
      style: AlertStyle(
        animationType: AnimationType.grow,
      ),
      context: Get.context,
      content: Container(
        child: Column(
          children: [
            Container(
              transform: Matrix4.translationValues(0, titleDistanceFromTop, 0),
              child: titleWidget,
            ),
            vspace24,
            Container(
              child: Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  GestureDetector(
                    child: UBText(
                      text: cancelText,
                      color: cancelTextColor,
                    ),
                    onTap: () => Get.back(),
                  ),
                  if (confirmText != null)
                    Container(
                      height: 24,
                      width: 1,
                      margin: const EdgeInsets.symmetric(horizontal: 24),
                      color: ColorName.white,
                    ),
                  if (confirmText != null)
                    GestureDetector(
                      child: UBText(
                        text: confirmText,
                        color: confirmTextColor,
                      ),
                      onTap: () {
                        onConfirm();
                        if (autoBackAfterConfirm) {
                          Get.back();
                        }
                      },
                    ),
                ],
              ),
            )
          ],
        ),
      ),
    ).show();
  }
}

openUpdatePopup({
  @required bool forceUpdate,
  @required List<String> features,
  @required List<String> bugFixes,
  @required String version,
  @required String url,
}) {
  final shouldUpdate = forceUpdate;
  Alert(
    onClose: () {
      final storage = GetStorage();
      storage.write(
        StorageKeys.lastCancelUpdate,
        {
          'date': DateTime.now().toString(),
          'version': version,
        },
      );
    },
    style: AlertStyle(
      animationType: AnimationType.grow,
      isCloseButton: !shouldUpdate,
      isOverlayTapDismiss: !shouldUpdate,
    ),
    context: Get.context,
    content: Container(
      child: Column(
        children: [
          Assets.images.warningInCircle.svg(),
          vspace12,
          UBText(text: 'Unitedbit new version release'),
          vspace24,
          if (features.isNotEmpty)
            FeatureList(
              title: 'New Features',
              features: features,
            ),
          if (features.isNotEmpty) vspace12,
          if (bugFixes.isNotEmpty)
            FeatureList(
              title: 'Bug Fixes',
              features: bugFixes,
            ),
          if (features.isNotEmpty || bugFixes.isNotEmpty) vspace24,
          SizedBox(
            width: 140,
            child: UBButton(
              onClick: () {
                launchURL(url);
              },
              text: 'Update to' & version,
            ),
          ),
          if (shouldUpdate != true) vspace8,
          if (shouldUpdate != true)
            SizedBox(
              width: 70,
              child: UBButton(
                height: 24,
                variant: ButtonVariant.Rounded,
                textColor: ColorName.greyd8,
                fontSize: 12.0,
                buttonColor: ColorName.black1c,
                onClick: () {
                  Get.back();
                },
                text: 'Ok, got it',
              ),
            )
        ],
      ),
    ),
  ).show();
}

openTransactionDetailsPopup(
    {@required Payments data, Function(int id) onCancelClick}) {
  final statusColor = data.status == 'completed'
      ? ColorName.green
      : data.status == 'rejected'
          ? ColorName.red
          : ColorName.grey80;
  final isPending = data.status == 'pending';
  final hasAddress = data.address != '' && data.address != null;
  final hasTxID = data.txId != '' && data.txId != null;
  final canExplore =
      (data.addressExplorerUrl != null && data.addressExplorerUrl != '') ||
          data.txIdExplorerUrl != null && data.txIdExplorerUrl != '';

  Get.bottomSheet(
    Container(
      height: 460,
      decoration: const BoxDecoration(
          color: ColorName.black2c, borderRadius: roundedTop_big),
      child: Column(
        children: [
          headerWithCloseButton(title: 'Order Details'),
          vspace12,
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              Text(
                data.code,
                style: whiteBold14,
              ),
              hspace4,
              Container(
                padding: const EdgeInsets.symmetric(
                  horizontal: 4,
                ),
                decoration: BoxDecoration(
                  borderRadius: rounded2,
                  color:
                      data.type == 'deposit' ? ColorName.green : ColorName.red,
                ),
                child: data.type == 'deposit'
                    ? Assets.images.greenDownArrow.svg(color: ColorName.white)
                    : Assets.images.redUpArrow.svg(color: ColorName.white),
              )
            ],
          ),
          Text(
            StringUtils.capitalize(data.type),
            style: grey80Bold10,
          ),
          vspace24,
          TwoPartText(
              title: 'Status',
              value: StringUtils.capitalize(data.status),
              valueColor: statusColor),
          TwoPartText(
            title: 'Coin',
            value: data.code,
          ),
          TwoPartText(
            title: 'Amount',
            value: data.amount.currencyFormat(
                removeInsignificantZeros: true, centFormat: true),
          ),
          TwoPartText(
            title: 'Date',
            value: data.createdAt,
          ),
          Padding(
            padding: px12,
            child: dividerGrey42,
          ),
          vspace12,
          SizedBox(
            width: Get.width - 24,
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                if (hasAddress)
                  Text(
                    'Address',
                    style: grey80Bold12,
                  ),
                if (hasAddress) vspace12,
                if (hasAddress)
                  SelectableText(
                    data.address,
                    style: whiteBold12,
                  ),
                if (hasTxID) vspace24,
                if (hasTxID)
                  Text(
                    'Transaction ID',
                    style: grey80Bold12,
                  ),
                if (hasTxID) vspace12,
                if (hasTxID)
                  SelectableText(
                    data.txId,
                    style: whiteBold12,
                  ),
              ],
            ),
          ),
          fill,
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  if (isPending)
                    UBButton(
                      height: 32.0,
                      width: (Get.width / 2) - 24.0,
                      borderColor: ColorName.red,
                      variant: ButtonVariant.Outline,
                      textColor: ColorName.red,
                      buttonColor: ColorName.black1c,
                      onClick: () {
                        if (onCancelClick != null) {
                          onCancelClick(data.id);
                        }
                      },
                      text: 'Cancel',
                    )
                  else if (hasTxID)
                    UBButton(
                      height: 32.0,
                      width: (Get.width / 2) - 24.0,
                      buttonColor: ColorName.black1c,
                      onClick: () {
                        final toastManager = ToastManager();

                        Clipboard.setData(
                          ClipboardData(
                            text: data.txId,
                          ),
                        );
                        toastManager.showToast(
                          'Transaction ID copied to clipboard',
                          type: ToastType.info,
                          action: ToastAction(
                            onPressed: (hideToastFn) {
                              hideToastFn();
                            },
                          ),
                        );
                      },
                      text: 'Copy TX ID',
                    ),
                  if (canExplore) hspace12,
                  if (canExplore)
                    UBButton(
                      height: 32.0,
                      width: (Get.width / 2) - 24.0,
                      onClick: () async {
                        await canLaunchUrl(Uri.parse(data.txIdExplorerUrl == ''
                                ? data.addressExplorerUrl
                                : data.txIdExplorerUrl))
                            ? await launchUrl(Uri.parse(
                                data.txIdExplorerUrl == ''
                                    ? data.addressExplorerUrl
                                    : data.txIdExplorerUrl))
                            : throw 'Could not launch ${data.txIdExplorerUrl == '' ? data.addressExplorerUrl : data.txIdExplorerUrl}';
                      },
                      text: "Check Explorer",
                    ),
                ],
              ),
            ],
          ),
          vspace24
        ],
      ),
    ),
    isScrollControlled: true,
    ignoreSafeArea: false,
  );
}

openOrderDetailsPopup(
    {@required OrderHistoryDetailModel details, OrderModel originalData}) {
  String filled = details.executed;
  if (details.executed != null && details.executed.contains('%')) {
    String tmp = details.executed.split('%')[0];
    if (double.parse(tmp) > 99.9999) {
      filled = '100%';
    }
  }
  Get.bottomSheet(
    Container(
      height: 320,
      decoration: const BoxDecoration(
          color: ColorName.black2c, borderRadius: roundedTop_big),
      child: Column(
        children: [
          headerWithCloseButton(title: 'Order Details'),
          vspace12,
          Row(
            mainAxisAlignment: MainAxisAlignment.center,
            crossAxisAlignment: CrossAxisAlignment.center,
            children: [
              Text(
                details.pair.replaceFirst('-', '/'),
                style: whiteBold14,
              ),
              hspace4,
              Container(
                transform: Matrix4.translationValues(0, 1, 0),
                padding: const EdgeInsets.symmetric(horizontal: 4, vertical: 2),
                decoration: BoxDecoration(
                  borderRadius: rounded2,
                  color: originalData.side == 'buy'
                      ? ColorName.green
                      : ColorName.red,
                ),
                child: Text(
                  originalData.side == 'buy' ? 'BUY' : 'Sell',
                  style: whiteBold8,
                ),
              )
            ],
          ),
          vspace8,
          Text(
            "Filled ${decimalCoin(value: filled)}",
            style: grey80Bold10,
          ),
          vspace24,
          TwoPartText(
            title: 'Type',
            value: StringUtils.capitalize(originalData.type),
          ),
          TwoPartText(
            title: 'Filled / Amount',
            value:
                "${decimalCoin(value: filled)} / ${decimalCoin(value: details.amount)}",
          ),
          Padding(
            padding: px12,
            child: dividerGrey42,
          ),
          vspace12,
          TwoPartText(
            title: 'Fee',
            value: decimalCoin(value: details.fee),
          ),
          TwoPartText(
            title: 'Total',
            value: decimalCoin(value: originalData.total),
          ),
          Padding(
            padding: px12,
            child: dividerGrey42,
          ),
          vspace12,
          ThreePartText(
            part1: 'Date',
            part2: 'Trading Price',
            part3: 'Filled',
            color: ColorName.grey80,
          ),
          vspace12,
          ThreePartText(
            part1: details.updatedAt ?? '-',
            part2: decimalCoin(
              value: details.price + " " + details.amount.split(" ")[1],
            ),
            part3: decimalCoin(value: details.executed),
            color: ColorName.white,
          ),
        ],
      ),
    ),
    isScrollControlled: true,
    ignoreSafeArea: false,
  );
}

class ThreePartText extends StatelessWidget {
  final String part1;
  final String part2;
  final String part3;
  final Color color;

  const ThreePartText({Key key, this.part1, this.part2, this.part3, this.color})
      : super(key: key);
  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: px12,
      child: Row(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Container(
            width: 2 * (Get.width / 5),
            child: Text(
              part1,
              style: TextStyle(
                color: color,
                fontSize: 12,
                fontWeight: FontWeight.w600,
              ),
            ),
          ),
          Container(
            child: Text(
              part2,
              style: TextStyle(
                color: color,
                fontSize: 12,
                fontWeight: FontWeight.w600,
              ),
            ),
          ),
          const Spacer(),
          Container(
            alignment: Alignment.centerRight,
            child: Text(
              part3,
              style: TextStyle(
                color: color,
                fontSize: 12,
                fontWeight: FontWeight.w600,
              ),
            ),
          ),
        ],
      ),
    );
  }
}

class TwoPartText extends StatelessWidget {
  final String title;
  final String value;
  final Color valueColor;

  const TwoPartText(
      {Key key, this.title, this.value, this.valueColor = ColorName.white})
      : super(key: key);
  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(
        left: 12,
        right: 12,
        bottom: 12,
      ),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Text(
            title,
            style: grey80Bold12,
          ),
          Text(
            value,
            style: TextStyle(
              color: valueColor,
              fontSize: 12.0,
              fontWeight: FontWeight.w600,
            ),
          ),
        ],
      ),
    );
  }
}

class FeatureList extends StatelessWidget {
  final List<String> features;
  final String title;
  const FeatureList({Key key, this.features, this.title}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        UBText(
          text: title,
          color: ColorName.greyd8,
        ),
        vspace8,
        for (var item in features)
          Row(
            children: [
              Container(
                height: 13.0,
                padding: const EdgeInsets.symmetric(
                  horizontal: 8.0,
                ),
                child: Center(
                  child: Container(
                    width: 3.0,
                    height: 3.0,
                    decoration: const BoxDecoration(
                        borderRadius: rounded_big, color: ColorName.greybf),
                  ),
                ),
              ),
              UBText(
                text: item,
                color: ColorName.greybf,
                weight: FontWeight.w400,
              )
            ],
          )
      ],
    );
  }
}
