import 'package:flutter/material.dart';

import '../../../generated/colors.gen.dart';

class AppBarWithBottomBorder extends StatelessWidget {
  final Widget title;

  const AppBarWithBottomBorder({Key key, this.title}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return AppBar(
      title: title,
      centerTitle: true,
      bottom: PreferredSize(
          child: Container(
            color: ColorName.greybf,
            height: 1.0,
          ),
          preferredSize: Size.fromHeight(10)),
    );
  }
}
