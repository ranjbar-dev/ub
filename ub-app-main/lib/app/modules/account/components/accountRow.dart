import 'package:flutter/material.dart';

import '../../../../generated/colors.gen.dart';
import '../../../common/components/UBButton.dart';

class AccountRow extends StatelessWidget {
  final icon;
  final String title;
  final String value;
  final String endButtonTitle;
  final String inactiveEndButtonTitle;
  final Function endButtonOnClick;
  final bool endButtonActive;
  final bool isLast;
  final Widget endWidget;
  const AccountRow({
    Key key,
    @required this.icon,
    this.title,
    @required this.value,
    this.endButtonTitle,
    this.endButtonOnClick,
    this.endButtonActive,
    this.isLast,
    this.inactiveEndButtonTitle,
    this.endWidget,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Container(
      decoration: BoxDecoration(
        border: isLast != true
            ? Border(
                bottom: BorderSide(
                  color: ColorName.grey42,
                  width: 1,
                ),
              )
            : null,
      ),
      width: double.infinity,
      height: 43,
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          Container(
            child: Row(
              children: [
                Container(
                  height: 24,
                  width: 24,
                  margin: const EdgeInsets.only(
                    right: 4,
                  ),
                  child: icon,
                ),
                if (title != null)
                  Container(
                    margin: const EdgeInsets.only(
                      right: 4,
                    ),
                    child: RichText(
                      text: TextSpan(
                        text: title,
                        style: const TextStyle(
                          color: ColorName.grey80,
                          fontSize: 13.0,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                    ),
                  ),
                Container(
                  child: RichText(
                    text: TextSpan(
                      text: value,
                      style: const TextStyle(
                        color: ColorName.white,
                        fontSize: 13.0,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ),
          if (endWidget != null)
            endWidget
          else
            Container(
              child: endButtonTitle != null
                  ? AbsorbPointer(
                      absorbing: endButtonActive == false,
                      child: UBButton(
                        onClick: endButtonOnClick,
                        height: 24,
                        width: 70,
                        text: endButtonActive != false
                            ? endButtonTitle
                            : inactiveEndButtonTitle ?? endButtonTitle,
                        buttonColor: ColorName.black,
                        textColor: endButtonActive != false
                            ? ColorName.white
                            : ColorName.grey80,
                        fontSize: 12.0,
                        padding: const EdgeInsets.symmetric(
                          horizontal: 10,
                        ),
                      ),
                    )
                  : null,
            ),
        ],
      ),
    );
  }
}
