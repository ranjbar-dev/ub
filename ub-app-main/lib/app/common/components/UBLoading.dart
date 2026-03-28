import 'package:flutter/material.dart';
import 'package:rive/rive.dart';

import '../../../generated/colors.gen.dart';

class UBLoading extends StatelessWidget {
  final bool useUnitedbitLogo;

  const UBLoading({Key key, this.useUnitedbitLogo = true}) : super(key: key);
  @override
  Widget build(BuildContext context) {
    return useUnitedbitLogo
        ? SizedBox(
            width: 37.0,
            child: RiveAnimation.asset(
              'assets/rive/loading.riv',
            ),
          )
        : CircularProgressIndicator(
            valueColor: AlwaysStoppedAnimation<Color>(
              ColorName.primaryBlue,
            ),
            strokeWidth: 3,
            backgroundColor: Colors.transparent,
          );
  }
}
