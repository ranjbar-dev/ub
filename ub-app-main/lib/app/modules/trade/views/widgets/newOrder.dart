import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../../common/components/UBButton.dart';
import '../../../../common/components/UBCountUp.dart';
import '../../../../common/components/UBFlipSwitcher.dart';
import '../../../../common/components/UBPercentSelect.dart';
import '../../../../common/components/UBShimmer.dart';
import '../../../../common/components/UBText.dart';
import '../../../../common/components/UBTopTabs.dart';
import '../../controllers/trade_controller.dart';
import 'tradeInput.dart';
import '../../../../routes/app_pages.dart';
import '../../../../../generated/colors.gen.dart';
import '../../../../../generated/locales.g.dart';
import '../../../../../utils/commonUtils.dart';
import '../../../../../utils/mixins/commonConsts.dart';
import '../../../../../utils/mixins/formatters.dart';
import '../../../../../utils/extentions/basic.dart';

double headerHeight = 34.0;
double topTabsHeight = 32.0;
List<List<List<StatelessWidget>>> tradeInputList = [
  [
    [
      AmountInput(),
      TotalInput(),
    ],
    [
      PercentSelect(),
      TradeInfo(),
    ],
    [
      PriceInput(),
      TradeSubmitButton(),
    ],
  ],
  [
    [
      AmountInput(),
      MarketPrice(),
    ],
    [
      PercentSelect(),
      TradeInfo(),
    ],
    [
      TradeSubmitButton(),
    ],
  ],
  [
    [
      AmountInput(),
      TotalInput(),
    ],
    [
      PriceInput(),
      StopInput(),
    ],
    [
      TradeSubmitButton(),
    ],
  ],
];

class NewOrder extends GetView<TradeController> with Formatter {
  final List<double> sizes;
  NewOrder({this.sizes});

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        Obx(() {
          final mainActiveIndex = controller.mainActiveIndex.value;
          final isLoadingPairBalances = controller.isLoadingPairBalance.value;
          var currency;
          var percision;
          var code;
          if (controller.pairBalanceData.value.sum != null) {
            currency =
                controller.pairBalanceData.value.pairBalances[mainActiveIndex];
            percision = coinPrecision(coinCode: currency.currencyCode);
            code = currency.currencyCode;
          }

          return Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              UBTopTabs(
                tabs: [
                  {
                    "name": LocaleKeys.buy.tr.toUpperCase(),
                    "textColor": ColorName.green
                  },
                  {
                    "name": LocaleKeys.sell.tr.toUpperCase(),
                    "textColor": ColorName.red
                  }
                ],
                tabHeight: topTabsHeight,
                onTabChange: (index) {
                  controller.handleMainTabChange(index);
                },
                activeIndex: mainActiveIndex,
              ),
              if (controller.pairBalanceData.value.sum != null &&
                  controller.globalController.loggedIn.value == true &&
                  !isLoadingPairBalances)
                BalanceValue(
                  balance: currency.balance,
                  code: code,
                  formatter: currencyFormatter,
                  percision: percision,
                )
              else if (isLoadingPairBalances)
                Padding(
                  padding: const EdgeInsets.only(right: 8.0),
                  child: UBShimmer(
                    width: 140,
                    height: 16,
                  ),
                )
              else
                const SizedBox()
            ],
          );
        }),
        Container(
          width: double.infinity,
          color: ColorName.black2c,
          child: Obx(
            () => DefaultTabController(
              length: 3, // length of tabs
              initialIndex: controller.subActiveIndex.value,
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: <Widget>[
                  Transform.translate(
                    offset: const Offset(0, 34.0),
                    child: Container(
                      width: double.infinity,
                      height: 1,
                      color: Colors.black,
                    ),
                  ),
                  Container(
                    height: headerHeight,
                    width: 92 * 3.0,
                    child: TabBar(
                      labelStyle: const TextStyle(
                          fontSize: 14.0, fontWeight: FontWeight.w600),
                      labelPadding: const EdgeInsets.symmetric(
                          horizontal: 12, vertical: 0),
                      isScrollable: true,
                      onTap: (idx) {
                        controller.handleSubTabChange(idx);
                      },
                      labelColor: controller.mainActiveIndex.value == 0
                          ? ColorName.green
                          : ColorName.red,
                      unselectedLabelColor: ColorName.greybf,
                      indicatorColor: controller.mainActiveIndex.value == 0
                          ? ColorName.green
                          : ColorName.red,
                      tabs: [
                        Tab(
                          text: LocaleKeys.limit.tr,
                        ),
                        Tab(text: LocaleKeys.market.tr),
                        Tab(text: LocaleKeys.stoplimit.tr),
                      ],
                    ),
                  ),
                  AnimatedContainer(
                    curve: Curves.easeInOut,
                    padding:
                        const EdgeInsets.symmetric(vertical: 12, horizontal: 6),
                    duration: const Duration(milliseconds: 300),
                    height: sizes[controller.subActiveIndex.value] -
                        headerHeight -
                        34,
                    child: TradeInputs(index: controller.subActiveIndex.value),
                  )
                ],
              ),
            ),
          ),
        ),
      ],
    );
  }
}

class BalanceValue extends StatelessWidget {
  final String balance;
  final int percision;
  final String Function(String) formatter;
  final String code;

  const BalanceValue({
    Key key,
    this.balance,
    this.percision,
    this.formatter,
    this.code,
  }) : super(key: key);
  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(right: 8.0),
      child: Row(
        children: [
          RichText(
            text: TextSpan(
              text: 'Available: ',
              style: const TextStyle(
                fontSize: 14.0,
                fontWeight: FontWeight.w400,
                color: ColorName.greybf,
              ),
            ),
          ),
          UBCountup(
            begin: 0,
            end: double.parse(balance),
            color: ColorName.white,
            duration: const Duration(seconds: 1),
            precision: percision,
          ),
          const SizedBox(width: 4),
          RichText(
            text: TextSpan(
              text: code,
              style: const TextStyle(
                fontSize: 14.0,
                fontWeight: FontWeight.w600,
                color: ColorName.white,
              ),
            ),
          ),
          const SizedBox(width: 4),
        ],
      ),
    );
  }
}

class PriceInput extends GetView<TradeController> {
  const PriceInput({
    Key key,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    int precision = pairPrecision(pairName: controller.currentPairName.value);
    print('precision : $precision');
    return Obx(() => TradeInput(
          focusColor: controller.mainActiveIndex.value == 0
              ? ColorName.green
              : ColorName.red,
          onChange: controller.handlePriceChange,
          value: controller.priceValue.value,
          title: controller.priceInputLabel.value.placeHolder,
          endText: controller.priceInputLabel.value.endLabel,
          precision: precision,
        ));
  }
}

class StopInput extends GetView<TradeController> {
  const StopInput({
    Key key,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Obx(
      () => TradeInput(
        focusColor: controller.mainActiveIndex.value == 0
            ? ColorName.green
            : ColorName.red,
        onChange: controller.handleStopChange,
        value: controller.stopValue.value,
        title: controller.stopPriceInputLabel.value.placeHolder,
        endText: controller.stopPriceInputLabel.value.endLabel,
      ),
    );
  }
}

class TotalInput extends GetView<TradeController> {
  const TotalInput({
    Key key,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Obx(
      () => TradeInput(
        focusColor: controller.mainActiveIndex.value == 0
            ? ColorName.green
            : ColorName.red,
        onChange: controller.handleTotalChange,
        value: controller.totalValue.value,
        title: controller.totalInputLabel.value.placeHolder,
        endText: controller.totalInputLabel.value.endLabel,
      ),
    );
  }
}

class MarketPrice extends GetView<TradeController> {
  const MarketPrice({
    Key key,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      height: 34.0,
      decoration: const BoxDecoration(
        borderRadius: rounded6,
        color: ColorName.black1c,
      ),
      padding: const EdgeInsets.only(left: 12.0),
      child: Align(
        alignment: Alignment.centerLeft,
        child: Obx(() {
          final mainActiveIndex = controller.mainActiveIndex.value;
          return UBText(
            text:
                '${mainActiveIndex == 0 ? "Buy" : "Sell"}' & 'at market price',
            color: ColorName.grey80,
            size: 14.0,
          );
        }),
      ),
    );
  }
}

class AmountInput extends GetView<TradeController> {
  const AmountInput({
    Key key,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Obx(
      () => TradeInput(
        focusColor: controller.mainActiveIndex.value == 0
            ? ColorName.green
            : ColorName.red,
        onChange: controller.handleAmountChange,
        value: controller.amountValue.value,
        title: controller.amountInputLabel.value.placeHolder,
        endText: controller.amountInputLabel.value.endLabel,
      ),
    );
  }
}

class TradeInputs extends StatelessWidget {
  final int index;
  const TradeInputs({Key key, this.index}) : super(key: key);
  @override
  Widget build(BuildContext context) {
    final widgets = tradeInputList[index];
    return Column(
      children: [
        for (var i = 0; i < widgets.length; i++)
          Expanded(
            child: Container(
              child: Row(
                children: [
                  for (var j = 0; j < widgets[i].length; j++)
                    Expanded(
                      child: Container(
                        padding: const EdgeInsets.symmetric(horizontal: 6),
                        child: Center(
                          child: widgets[i][j],
                        ),
                      ),
                    ),
                ],
              ),
            ),
          )
      ],
    );
  }
}

class TradeInfo extends GetView<TradeController> with Formatter {
  const TradeInfo({
    Key key,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.only(top: 4),
      alignment: Alignment.topLeft,
      height: double.infinity,
      width: double.infinity,
      child: Obx(
        () {
          final youGet = controller.youGet.value;
          final fee = controller.tradeFee.value;
          final mainActiveTabIndex = controller.mainActiveIndex.value;
          final currentPairName = controller.currentPairName.value;
          final pairs = currentPairName.split('-');
          final isMarket = controller.subActiveIndex.value == 1;
          return Column(
            children: [
              InfoRow(
                title: 'Trade fee:',
                isMarket: isMarket && fee != '',
                value: fee == ''
                    ? '0.00'
                    : decimalCoin(
                        value: fee, coinCode: pairs[mainActiveTabIndex]),
                appendix: pairs[mainActiveTabIndex],
              ),
              vspace2,
              InfoRow(
                title: 'You will get:',
                isMarket: isMarket && youGet != '',
                value: youGet == ''
                    ? '0.00'
                    : decimalCoin(
                        value: youGet, coinCode: pairs[mainActiveTabIndex]),
                appendix: pairs[mainActiveTabIndex],
              ),
            ],
          );
        },
      ),
    );
  }
}

class InfoRow extends StatelessWidget {
  final String title;
  final String value;
  final String appendix;
  final bool isMarket;
  const InfoRow({
    Key key,
    this.title,
    this.value,
    this.appendix,
    this.isMarket,
  }) : super(key: key);
  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        RichText(text: TextSpan(text: title, style: grey80Bold11)),
        const Spacer(),
        RichText(
            text: TextSpan(
                text: (isMarket ? '~ ' : '') + value, style: whiteBold11)),
        hspace4,
        RichText(text: TextSpan(text: appendix, style: grey80Bold11)),
      ],
    );
  }
}

class PercentSelect extends GetView<TradeController> {
  const PercentSelect({
    Key key,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Obx(() {
      final selectedIndex = controller.selectedPercentIndex.value;
      final numberOfSegments = controller.numberOfPercentSegments;
      final selectedTradeSide = controller.mainActiveIndex.value;
      return UBPercentSelect(
          selectedIndex: selectedIndex,
          numberOfSegments: numberOfSegments,
          selectedColor:
              selectedTradeSide == 0 ? ColorName.green : ColorName.red,
          onPercentClick: (i) {
            controller.handlePercentClick(index: i);
          });
    });
  }
}

class TradeSubmitButton extends GetView<TradeController> {
  @override
  Widget build(BuildContext context) {
    return Obx(() {
      final disabled = controller.pairBalanceData.value.sum == null ||
          controller.isLoadingPairBalance.value;
      final isLoading = controller.isCreatingOrder.value;
      final isBuy = controller.mainActiveIndex.value == 0;
      return controller.globalController.loggedIn.value == true
          ? UBFlipSwitcher(
              child1: UBButton(
                key: ValueKey('front'),
                disabled: disabled,
                isLodaing: isLoading,
                onClick: () {
                  controller.handleSubmitClick();
                },
                text: LocaleKeys.buy.tr.toUpperCase(),
                fontSize: 16,
                buttonColor: ColorName.green,
                height: 32,
              ),
              child2: UBButton(
                key: ValueKey('SellSubmitButton'),
                disabled: disabled,
                isLodaing: isLoading,
                onClick: () {
                  controller.handleSubmitClick();
                },
                text: LocaleKeys.sell.tr.toUpperCase(),
                fontSize: 16,
                buttonColor: ColorName.red,
                height: 32,
              ),
              conditionToShowChild1: isBuy)
          : UBButton(
              onClick: () {
                Get.toNamed(AppPages.LOGIN);
              },
              fontSize: 16,
              text: LocaleKeys.loginToSubmitNewOrder.tr,
              buttonColor: ColorName.primaryBlue,
              height: 32,
            );
    });
  }
}
