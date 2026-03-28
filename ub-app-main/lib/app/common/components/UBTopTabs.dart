import 'package:flutter/material.dart';

import '../../../generated/colors.gen.dart';
import 'UBText.dart';

class UBTopTabs extends StatelessWidget {
  final tabs;
  final double tabHeight;
  final Function(int index) onTabChange;
  final int activeIndex;
  const UBTopTabs(
      {Key key,
      this.tabs,
      this.onTabChange,
      this.activeIndex,
      this.tabHeight = 32.0})
      : super(key: key);

  Widget build(BuildContext context) {
    return Container(
      height: tabHeight,
      child: Row(
        children: [
          for (var i = 0; i < tabs.length; i++)
            GestureDetector(
              onTapDown: (e) {
                onTabChange(i);
              },
              child: Container(
                width: 60,
                decoration: BoxDecoration(
                  color: activeIndex == i ? ColorName.black2c : ColorName.black,
                  borderRadius: const BorderRadius.only(
                    topLeft: const Radius.circular(4),
                    topRight: const Radius.circular(4),
                  ),
                ),
                child: Align(
                  child: UBText(
                    size: 16,
                    weight:
                        activeIndex != i ? FontWeight.normal : FontWeight.w600,
                    text: tabs[i]["name"],
                    color: activeIndex == i
                        ? tabs[i]["textColor"]
                        : ColorName.white,
                  ),
                ),
              ),
            ),
        ],
      ),
    );
  }
}
