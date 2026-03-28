import 'package:flutter/material.dart';
import '../../../generated/assets.gen.dart';
import '../../../generated/colors.gen.dart';
import '../../../utils/mixins/popups.dart';

class UBDatePicker extends StatelessWidget with Popups {
  final double width;
  final double height;
  final String placeHolder;
  final String value;
  final Function(String date) onDateSelect;
  final Color backgroundColor;
  final Function() onClearClick;
  final Color filledDateBorderColor;

  const UBDatePicker({
    Key key,
    @required this.width,
    this.height = 36.0,
    @required this.placeHolder,
    this.value = '',
    this.backgroundColor = ColorName.black1c,
    this.filledDateBorderColor = ColorName.primaryBlue,
    @required this.onDateSelect,
    this.onClearClick,
  }) : super(key: key);
  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: () {
        openDatePickerPopup(onDateSelect: onDateSelect);
      },
      child: Container(
        height: 36,
        width: width,
        decoration: BoxDecoration(
          color: ColorName.black1c,
          borderRadius: const BorderRadius.all(
            const Radius.circular(
              4,
            ),
          ),
          border: Border.all(
            color: value == '' ? ColorName.black1c : ColorName.primaryBlue,
          ),
        ),
        padding: const EdgeInsets.symmetric(
          horizontal: 12.0,
          vertical: 6.0,
        ),
        child: Row(
          children: [
            RichText(
              text: TextSpan(
                  text: value == '' ? placeHolder : value,
                  style: const TextStyle(
                      color: ColorName.greyd8,
                      fontSize: 13,
                      fontWeight: FontWeight.w600)),
            ),
            const Spacer(),
            if (onClearClick != null && value != '')
              GestureDetector(
                onTap: onClearClick,
                child: Assets.images.closeIcon.svg(),
              ),
            Assets.images.calendar.svg(),
          ],
        ),
      ),
    );
  }
}
