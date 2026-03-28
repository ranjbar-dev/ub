import 'package:flutter/widgets.dart';

import '../../../generated/assets.gen.dart';

class AccountCardRowModel {
  final String title;
  final String value;
  final String endButtonTitle;
  final String inactiveEndButtonTitle;
  final bool endButtonActive;
  final bool isLast;
  final Function onEndButtonClick;
  final SvgGenImage icon;
  final Widget endWidget;

  AccountCardRowModel({
    this.title,
    this.value,
    this.endButtonTitle,
    this.inactiveEndButtonTitle,
    this.endButtonActive,
    this.isLast,
    this.onEndButtonClick,
    this.endWidget,
    this.icon,
  });
}
