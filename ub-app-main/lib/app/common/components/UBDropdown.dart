import 'package:flutter/material.dart';

import '../../../generated/colors.gen.dart';
import 'UBText.dart';

class UBDropDown extends StatelessWidget {
  final String value;
  final List<dynamic> options;
  final Function onChange;
  final bool expanded;
  final bool dense;

  const UBDropDown(
      {Key key,
      this.value,
      this.options,
      this.onChange,
      this.expanded,
      this.dense})
      : super(key: key);
  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        DropdownButton<String>(
          dropdownColor: ColorName.black2c,
          isDense: dense ?? false,
          isExpanded: expanded ?? false,
          value: value,
          icon: null,
          iconSize: 0,
          elevation: 16,
          style: const TextStyle(color: ColorName.white, fontSize: 15),
          underline: const SizedBox(
            height: 0,
          ),
          onChanged: onChange,
          items: options.map<DropdownMenuItem<String>>((dynamic item) {
            return DropdownMenuItem<String>(
              value: item['value'],
              child: Container(
                child: UBText(
                  text: item['name'],
                  color: item['name'] == value
                      ? ColorName.primaryBlue
                      : ColorName.white,
                ),
              ),
            );
          }).toList(),
        ),
        IgnorePointer(
          child: Container(
            padding: const EdgeInsets.only(top: 2.0),
            color: ColorName.black2c,
            child: Row(
              children: [
                UBText(
                  text: value,
                  color: ColorName.white,
                  size: 15.0,
                ),
                const Icon(
                  Icons.keyboard_arrow_down,
                  size: 18.0,
                  color: ColorName.white,
                )
              ],
            ),
          ),
        )
      ],
    );
  }
}
