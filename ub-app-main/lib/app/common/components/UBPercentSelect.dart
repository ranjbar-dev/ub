import 'package:flutter/material.dart';

import '../../../generated/colors.gen.dart';

class UBPercentSelect extends StatelessWidget {
  final int selectedIndex;
  final int numberOfSegments;
  final Function(int) onPercentClick;
  final Color selectedColor;
  const UBPercentSelect({
    Key key,
    @required this.selectedIndex,
    this.numberOfSegments = 4,
    @required this.onPercentClick,
    this.selectedColor = ColorName.primaryBlue,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.only(left: 1, top: 8),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          for (var i = 0; i < numberOfSegments; i++)
            GestureDetector(
              onTap: () {
                onPercentClick(i);
                return;
              },
              child: ConstrainedBox(
                constraints: BoxConstraints(maxWidth: 38),
                child: Container(
                  color: ColorName.black2c,
                  child: Column(
                    children: [
                      SizedBox(
                        height: 8,
                        child: Container(
                          height: i <= selectedIndex ? 8.0 : 2.0,
                          decoration: BoxDecoration(
                            borderRadius: const BorderRadius.all(
                              const Radius.circular(2),
                            ),
                            color: i <= selectedIndex
                                ? (selectedColor)
                                : ColorName.black,
                          ),
                        ),
                      ),
                      const SizedBox(
                        height: 2,
                      ),
                      RichText(
                        text: TextSpan(
                            text:
                                '${((100 / numberOfSegments) * (i + 1)).toInt().toString()}%',
                            style: TextStyle(
                              fontSize: 10,
                              color: i <= selectedIndex
                                  ? (selectedColor)
                                  : ColorName.grey80,
                            )),
                      )
                    ],
                  ),
                ),
              ),
            )
        ],
      ),
    );
  }
}
