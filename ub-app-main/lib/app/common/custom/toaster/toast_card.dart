part of 'utopic_toast.dart';

class _ToastCard extends StatelessWidget {
  final String message;
  final ToastAction action;
  final ToastType type;
  final Widget endAction;
  const _ToastCard({
    Key key,
    this.message,
    this.action,
    this.type,
    this.endAction,
  })  : assert(key != null),
        super(key: key);

  @override
  Widget build(BuildContext context) {
    ToastOverlay toastOverlay;
    try {
      toastOverlay = Provider.of<ToastOverlay>(context);
    } on ProviderNotFoundException catch (e) {
      print(e);
    }

    Color backgroundColor =
        toastOverlay?.successfullBackgroundColor ?? ColorName.grey16;
    Color iconContainerColor;
    IconData icon;
    switch (type) {
      case ToastType.error:
        iconContainerColor = ColorName.red;
        icon = Icons.error;
        break;
      case ToastType.warning:
        iconContainerColor = ColorName.orange;
        icon = Icons.warning_amber_rounded;
        break;
      case ToastType.info:
        iconContainerColor = ColorName.primaryBlue;
        icon = Icons.info;
        break;
      case ToastType.success:
        iconContainerColor = ColorName.green;
        icon = Icons.check_circle;
        break;
      default:
        iconContainerColor = ColorName.primaryBlue;
        icon = Icons.info;
    }
    Widget result = Card(
      color: backgroundColor.withOpacity(0.97),
      margin: EdgeInsets.only(top: 12),
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12.0),
      ),
      child: GestureDetector(
        onTap: toastOverlay.enableTapToHide
            ? () => ToastManager()._hideToastByKey(key)
            : null,
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 17, vertical: 10),
          child: Row(
            mainAxisSize: MainAxisSize.min,
            children: [
              Icon(
                icon,
                size: 16,
                color: iconContainerColor,
              ),
              const SizedBox(
                width: 4,
              ),
              Expanded(
                child: Text(
                  message,
                  textAlign: TextAlign.left,
                  style: TextStyle(
                    fontSize: 12,
                    color: ColorName.greybf,
                  ),
                ),
              ),
              endAction ?? const SizedBox(),
              if (endAction != null) hspace4,
              action?.build(
                    context,
                    () => ToastManager()._hideToastByKey(key),
                  ) ??
                  const SizedBox(),
            ],
          ),
        ),
      ),
    );

    if (toastOverlay?.enableSwipeToDismiss != false) {
      result = Dismissible(
        key: key,
        onDismissed: (_) {
          ToastManager()._hideToastByKey(key, showAnim: false);
        },
        child: result,
      );
    }

    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 4, horizontal: 8),
      child: result,
    );
  }
}
