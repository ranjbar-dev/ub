import 'package:basic_utils/basic_utils.dart';
import 'package:flutter/material.dart';

import '../../../../../generated/colors.gen.dart';
import '../../../../common/components/UBDottedBorder.dart';
import '../../../../common/components/UBText.dart';
import '../../user_profile_model.dart';

class Confirmed extends StatelessWidget {
  const Confirmed({
    Key key,
    @required this.side,
    @required this.subType,
  }) : super(key: key);
  final String side;
  final SubTypes subType;

  @override
  Widget build(BuildContext context) {
    return UBDottedBorder(
      color: ColorName.green,
      child: Container(
        decoration: BoxDecoration(
          borderRadius: BorderRadius.circular(12),
          color: ColorName.green.withOpacity(0.15),
        ),
        child: Center(
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              UBText(
                text:
                    "$side side of your ${StringUtils.capitalize(subType.name.replaceAll('_', ' '), allWords: true)} is accepted!",
                color: ColorName.green,
              ),
              const SizedBox(
                height: 12,
              ),
              const Icon(
                Icons.check_circle_sharp,
                color: ColorName.green,
              )
            ],
          ),
        ),
      ),
    );
  }
}
