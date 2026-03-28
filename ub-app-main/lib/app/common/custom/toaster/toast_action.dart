part of 'utopic_toast.dart';

class ToastAction {
  final bool withClose;
  final Color textColor;
  final Color disabledTextColor;
  final void Function(void Function()) onPressed;

  const ToastAction({
    Key key,
    this.onPressed,
    this.disabledTextColor,
    this.withClose = true,
    this.textColor,
  });

  Widget build(BuildContext context, void Function() hideToast) {
    return Padding(
      padding: const EdgeInsets.only(left: 8.0),
      child: withClose == false
          ? const SizedBox()
          : GestureDetector(
              child: Container(
                color: Colors.transparent,
                padding: const EdgeInsets.symmetric(
                  horizontal: 12,
                  vertical: 10,
                ),
                child: const Icon(
                  Icons.close,
                  color: ColorName.grey80,
                  size: 12,
                ),
              ),
              onTap: () => onPressed(hideToast),
            ),
    );
  }
}
