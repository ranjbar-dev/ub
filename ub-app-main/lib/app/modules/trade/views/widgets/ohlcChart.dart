import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../../common/components/CenterUBLoading.dart';
import '../../../../common/components/UBScaleSwitcher.dart';
import '../../../../common/custom/k_chart/k_chart_widget.dart';
import '../../controllers/ohlcChart_controller.dart';
import '../../../../../generated/colors.gen.dart';

import 'TimeFrameSelector.dart';

class OHLCChart extends GetView<OHLCChartController> {
  @override
  Widget build(BuildContext context) {
    return Obx(
      () {
        // ignore: invalid_use_of_protected_member
        final chartData = controller.chartData.value;
        final isLine = controller.isLine.value;
        final mainState = controller.mainState.value;
        final secondaryState = controller.secondaryState.value;
        final isLoading = controller.isLoadingOhlc.value;
        final isDetailsOpen = controller.isOhlcDetailsOpen.value;
        return isLoading
            ? CenterUbLoading()
            : Stack(
                children: [
                  KChartWidget(
                    chartData,
                    isLine: isLine,
                    mainState: mainState,
                    // volHidden: true,
                    secondaryState: secondaryState,
                    onDetailsOpenChange: (isOpen) {
                      controller.handleIsOhlcDetailsOpen(isOpen: isOpen);
                    },
                    fixedLength: 8,
                    timeFormat: TimeFormat.YEAR_MONTH_DAY_WITH_HOUR,
                    bgColor: [
                      ColorName.black,
                      ColorName.black,
                      ColorName.black
                    ],
                  ),
                  Positioned(
                    top: 20.0,
                    left: 14.0,
                    child: UBScaleSwitcher(
                      conditionToShowChild1: isDetailsOpen,
                      child1: AbsorbPointer(
                        child: TimeFrameContainer(
                          key: ValueKey('timeFrame'),
                          color: Colors.transparent,
                        ),
                      ),
                      child2: TimeFrameSelector(),
                    ),
                  )
                ],
              );
      },
    );
  }

  // Widget buildButtons() {
  //   return Column(
  //     children: [
  //       Wrap(
  //         runSpacing: 6.0,
  //         spacing: 6.0,
  //         alignment: WrapAlignment.center,
  //         children: [
  //           button(
  //             "Line",
  //             onPressed: () => isLine = true,
  //             selected: isLine,
  //           ),
  //           button(
  //             "Bars",
  //             onPressed: () => isLine = false,
  //             selected: !isLine,
  //           ),
  //         ],
  //       ),
  //       Padding(
  //         padding: EdgeInsets.only(
  //           top: 6.0,
  //         ),
  //       ),
  //       Wrap(
  //         runSpacing: 6.0,
  //         spacing: 6.0,
  //         alignment: WrapAlignment.center,
  //         children: [
  //           button(
  //             "MACD",
  //             onPressed: () => _secondaryState = SecondaryState.MACD,
  //             selected: _mainState == MainState.MA,
  //           ),
  //           button(
  //             "KDJ",
  //             onPressed: () => _secondaryState = SecondaryState.KDJ,
  //             selected: _secondaryState == SecondaryState.KDJ,
  //           ),
  //           button(
  //             "RSI",
  //             onPressed: () => _secondaryState = SecondaryState.RSI,
  //             selected: _secondaryState == SecondaryState.RSI,
  //           ),
  //           button(
  //             "WR",
  //             onPressed: () => _secondaryState = SecondaryState.WR,
  //             selected: _secondaryState == SecondaryState.WR,
  //           ),
  //           button(
  //             "NONE",
  //             onPressed: () => _secondaryState = SecondaryState.NONE,
  //             selected: _secondaryState == SecondaryState.NONE,
  //           ),
  //         ],
  //       ),
  //       Padding(
  //         padding: EdgeInsets.only(
  //           top: 6.0,
  //         ),
  //       ),
  //       Wrap(
  //         runSpacing: 6.0,
  //         spacing: 6.0,
  //         alignment: WrapAlignment.center,
  //         children: [
  //           button(
  //             "MA",
  //             onPressed: () => _mainState = MainState.MA,
  //             selected: _mainState == MainState.MA,
  //           ),
  //           button(
  //             "BOLL",
  //             onPressed: () => _mainState = MainState.BOLL,
  //             selected: _mainState == MainState.BOLL,
  //           ),
  //           button(
  //             "NONE",
  //             onPressed: () => _mainState = MainState.NONE,
  //             selected: _mainState == MainState.NONE,
  //           ),
  //         ],
  //       ),
  //       Padding(
  //         padding: EdgeInsets.only(
  //           top: 6.0,
  //         ),
  //       ),
  //     ],
  //   );
  // }

  // Widget button(String text, {VoidCallback onPressed, bool selected = false}) {
  //   return SizedBox(
  //     width: 50.0,
  //     height: 30.0,
  //     child: FlatButton(
  //       padding: EdgeInsets.all(0.0),
  //       onPressed: () {
  //         if (onPressed != null) {
  //           onPressed();
  //           setState(() {});
  //         }
  //       },
  //       child: Text(
  //         text,
  //         style: TextStyle(
  //           fontSize: 12.0,
  //         ),
  //       ),
  //       color: selected ? Colors.blue : Colors.blue.withOpacity(0.6),
  //     ),
  //   );
  // }

  void getData(String period) {
/*
list of {

"id" : 1619712000
"open" : 53709.7
"close" : 54746.93
"low" : 52336.36
"high" : 54786.0
"amount" : 15234.622333590716
"vol" : 814797176.240874
"count" : 653355
}
 */
    // final result = chartData;
    // Map parseJson = json.decode(result);
    // List list = parseJson['data'];
    // datas = list
    //     .map((item) => KLineEntity.fromJson(item))
    //     .toList()
    //     .reversed
    //     .toList()
    //     .cast<KLineEntity>();
    // DataUtil.calculate(datas);
    // showLoading = false;
    // // setState(() {});
  }
}
