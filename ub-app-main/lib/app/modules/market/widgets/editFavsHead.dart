import 'package:flutter/material.dart';
import 'package:get/state_manager.dart';
import '../controllers/market_controller.dart';
import '../../../../generated/colors.gen.dart';
import '../../../../utils/mixins/commonConsts.dart';

class EditFavsHead extends GetView<MarketController> {
  @override
  Widget build(BuildContext context) {
    return Container(
        height: 24,
        color: ColorName.grey16,
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 12.0),
          child: Row(
            children: [
              SizedBox(
                  width: thirdWidthPlus24,
                  child: _header(
                    title: 'Coin',
                  )),
              const Spacer(),
              _header(title: 'Top'),
              hspace24,
              _header(
                title: 'Sort',
              )
            ],
          ),
        ));
  }

  _header({String title}) {
    return Text(
      title,
      style: grey80Bold13,
    );
  }
}
