import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../../common/components/UBWrappedButtons.dart';
import '../../controllers/ohlcChart_controller.dart';
import '../../../../../generated/colors.gen.dart';

class TimeFrameSelector extends GetView<OHLCChartController> {
  @override
  Widget build(BuildContext context) {
    return Obx(() {
      final timeFrameButtons = controller.timeFrameButtons;
      final selectedTimeFrameIndex =
          controller.selectedTimeFrameButtonIndex.value;
      final selectedText =
          controller.timeFrameButtons[selectedTimeFrameIndex].text;
      final isTimeFrameOpoupOpen = controller.isTimeFramePoupOpen.value;
      return Container(
        height: isTimeFrameOpoupOpen ? 200 : 24,
        child: Stack(
          clipBehavior: Clip.none,
          children: [
            for (var i = 0; i < timeFrameButtons.length; i++)
              AnimatedPositioned(
                curve: Curves.ease,
                top: isTimeFrameOpoupOpen ? (((i + 1).toDouble() * 24 + 4)) : 0,
                duration: const Duration(milliseconds: 200),
                child: SmallButton(
                  borderColor: selectedTimeFrameIndex == i
                      ? ColorName.primaryBlue
                      : Colors.transparent,
                  minWidth: 80,
                  text: timeFrameButtons[i].text,
                  onClick: () => controller.handleTimeFrameChange(i),
                ),
              ),
            GestureDetector(
              onTap: () {
                controller.isTimeFramePoupOpen.toggle();
              },
              child: TimeFrameContainer(
                color: ColorName.black2c,
                child: Align(
                  child: RichText(
                    text: TextSpan(
                      text: selectedText,
                      style: TextStyle(
                        fontSize: 11.0,
                        fontWeight: FontWeight.w600,
                        color: ColorName.greyd8,
                      ),
                    ),
                  ),
                ),
              ),
            ),
          ],
        ),
      );
    });
  }
}

class TimeFrameContainer extends StatelessWidget {
  final Color color;
  final Widget child;

  const TimeFrameContainer({
    Key key,
    this.color,
    this.child,
  }) : super(key: key);
  @override
  Widget build(BuildContext context) {
    return Container(
      width: 80.0,
      height: 24.0,
      color: color,
      child: child,
    );
  }
}
