import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../../../generated/colors.gen.dart';
import '../../../../common/components/UBText.dart';
import '../../controllers/exchange_controller.dart';

class ExchangeInput extends GetView<ExchangeController> {
  const ExchangeInput({
    Key key,
    this.hasBadge = false,
    this.value,
    this.precision,
    this.onChange,
  }) : super(key: key);

  final bool hasBadge;
  final String value;
  final int precision;
  final Function onChange;

  @override
  Widget build(BuildContext context) {
    return Obx(
      () => Container(
        height: 40,
        child: Theme(
          data: Theme.of(context).copyWith(splashColor: Colors.transparent),
          child: Stack(
            children: [
              TextField(
                enabled: hasBadge,
                autofocus: false,
                onChanged: (text) =>
                    hasBadge ? controller.calcHowMuchYouWillGet() : null,
                controller: hasBadge
                    ? controller.inputControllerFrom.value
                    : controller.inputControllerTo.value,
                onTap: () => hasBadge
                    ? controller.inputControllerFrom.value.selection =
                        TextSelection(
                            baseOffset: 0,
                            extentOffset: controller
                                .inputControllerFrom.value.text.length)
                    : null,
                keyboardType: TextInputType.numberWithOptions(decimal: true),
                style: TextStyle(
                    fontSize: 16.0,
                    color: hasBadge ? ColorName.white : ColorName.grey80,
                    fontWeight: FontWeight.bold),
                decoration: InputDecoration(
                  contentPadding: hasBadge
                      ? EdgeInsets.only(
                          bottom: 10, left: 6 // HERE THE IMPORTANT PART
                          )
                      : EdgeInsets.only(
                          bottom: 10, left: 24 // HERE THE IMPORTANT PART
                          ),
                  filled: true,
                  fillColor: ColorName.black,
                  enabledBorder: UnderlineInputBorder(
                    borderSide: BorderSide(color: ColorName.black),
                    borderRadius: BorderRadius.circular(6),
                  ),
                  focusedBorder: OutlineInputBorder(
                    borderSide: BorderSide(color: ColorName.primaryBlue),
                    borderRadius: BorderRadius.circular(6),
                  ),
                  border: UnderlineInputBorder(
                    borderSide: BorderSide(color: ColorName.black),
                    borderRadius: BorderRadius.circular(6),
                  ),
                ),
              ),
              Visibility(
                visible: hasBadge,
                child: Positioned(
                  right: 24,
                  top: 12,
                  child: InkWell(
                    onTap: () => controller.allIn(),
                    child: Container(
                      height: 16,
                      width: 30,
                      decoration: BoxDecoration(
                        color: ColorName.grey42,
                        borderRadius: BorderRadius.circular(6),
                      ),
                      child: Center(
                        child: UBText(
                          text: 'All',
                          size: 11,
                          color: ColorName.grey97,
                        ),
                      ),
                    ),
                  ),
                ),
              )
            ],
          ),
        ),
      ),
    );
  }
}
