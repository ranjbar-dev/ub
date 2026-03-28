import 'package:flutter/material.dart';
import '../../../utils/mixins/toast.dart';
import '../../../utils/throttle.dart';

final toastThrottle =
    new Throttling(duration: const Duration(milliseconds: 4000));

class UBToastOnTap extends StatelessWidget with Toaster {
  final Widget child;
  final Widget toastEnd;
  final bool active;
  final String toastText;
  final Function onTap;
  const UBToastOnTap(
      {Key key,
      @required this.child,
      this.toastEnd,
      this.active = false,
      this.toastText,
      this.onTap})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: () {
        if (active) {
          if (onTap == null) {
            if (toastEnd != null) {
              toastThrottle.throttle(() {
                toastAction(toastText, toastEnd);
              });
            } else {
              toastThrottle.throttle(() {
                toastInfo(
                  toastText,
                );
              });
            }
          } else {
            toastThrottle.throttle(() {
              onTap();
            });
          }
        }
      },
      child: AbsorbPointer(
        absorbing: active,
        child: child,
      ),
    );
  }
}
