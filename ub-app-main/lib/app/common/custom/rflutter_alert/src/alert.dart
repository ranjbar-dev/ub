import 'dart:async';

import 'package:flutter/material.dart';

import 'alert_style.dart';
import 'animation_transition.dart';
import 'constants.dart';
import '../../../../../generated/colors.gen.dart';

/// Main class to create alerts.
///
/// You must call the "show()" method to view the alert class you have defined.
class Alert {
  final String id;
  final BuildContext context;
  final AlertType type;
  final AlertStyle style;
  final Widget image;
  final String title;
  final String desc;
  final Widget content;
  final Widget header;
  final Function closeFunction;
  final Widget closeIcon;
  final bool onWillPopActive;
  final bool useRootNavigator;
  final Function onClose;
  final AlertAnimation alertAnimation;

  /// Alert constructor
  ///
  /// [context], [title] are required.
  Alert({
    @required this.context,
    this.onClose,
    this.id,
    this.type,
    this.style = const AlertStyle(),
    this.image,
    this.title,
    this.desc,
    this.content = const SizedBox(),
    this.closeFunction,
    this.closeIcon,
    this.onWillPopActive = false,
    this.alertAnimation,
    this.useRootNavigator = true,
    this.header,
  });

  /// Displays defined alert window
  Future<bool> show() async {
    return await showGeneralDialog(
        context: context,
        pageBuilder: (BuildContext buildContext, Animation<double> animation,
            Animation<double> secondaryAnimation) {
          return _buildDialog();
        },
        barrierDismissible: style.isOverlayTapDismiss,
        barrierLabel:
            MaterialLocalizations.of(context).modalBarrierDismissLabel,
        barrierColor: style.overlayColor,
        useRootNavigator: useRootNavigator,
        transitionDuration: style.animationDuration,
        transitionBuilder: (
          BuildContext context,
          Animation<double> animation,
          Animation<double> secondaryAnimation,
          Widget child,
        ) =>
            alertAnimation == null
                ? _showAnimation(animation, secondaryAnimation, child)
                : alertAnimation(
                    context, animation, secondaryAnimation, child));
  }

  /// Dismisses the alert dialog.
  Future<void> dismiss() async {
    Navigator.of(context, rootNavigator: useRootNavigator).pop();
  }

  /// Alert dialog content widget
  Widget _buildDialog() {
    final Widget _child = ConstrainedBox(
      constraints: style.constraints ??
          const BoxConstraints.expand(
            width: double.infinity,
            height: double.infinity,
          ),
      child: WillPopScope(
        onWillPop: () {
          final canClose = style.isOverlayTapDismiss;
          if (canClose && onClose != null) {
            onClose();
          }
          return Future.value(canClose);
        },
        child: Align(
          alignment: style.alertAlignment,
          child: SingleChildScrollView(
            child: AlertDialog(
              key: id == null ? null : Key(id),
              backgroundColor: ColorName.black2c,
              shape: style.alertBorder ?? _defaultShape(),
              insetPadding: style.alertPadding,
              elevation: style.alertElevation,
              titlePadding: const EdgeInsets.all(0.0),
              title: Container(
                child: Center(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: <Widget>[
                      if (header != null)
                        Row(
                          mainAxisAlignment: MainAxisAlignment.spaceBetween,
                          crossAxisAlignment: CrossAxisAlignment.center,
                          children: [header, _getCloseButton()],
                        )
                      else
                        _getCloseButton(),
                      Padding(
                        padding: EdgeInsets.fromLTRB(
                            12, (style.isCloseButton ? 0 : 10), 12, 12),
                        child: Column(
                          children: <Widget>[
                            const SizedBox(
                              height: 5,
                            ),
                            if (title != null)
                              Text(
                                title,
                                style: style.titleStyle,
                                textAlign: style.titleTextAlign,
                              ),
                            const SizedBox(height: 5),
                            desc == null
                                ? Container()
                                : Text(
                                    desc,
                                    style: style.descStyle,
                                    textAlign: style.descTextAlign,
                                  ),
                            content,
                          ],
                        ),
                      )
                    ],
                  ),
                ),
              ),
              contentPadding: style.buttonAreaPadding,
            ),
          ),
        ),
      ),
    );
    return onWillPopActive
        ? WillPopScope(onWillPop: () async => false, child: _child)
        : _child;
  }

  /// Returns the close button on the top right
  Widget _getCloseButton() {
    return style.isCloseButton
        ? Padding(
            padding: const EdgeInsets.fromLTRB(0, 6, 10, 0),
            child: GestureDetector(
              onTap: () {
                if (closeFunction == null) {
                  Navigator.of(context, rootNavigator: useRootNavigator).pop();
                } else {
                  closeFunction();
                }
              },
              child: Container(
                alignment: FractionalOffset.topRight,
                child: this.closeIcon != null
                    ? Container(child: this.closeIcon)
                    : Container(
                        width: 20,
                        height: 20,
                        child: Icon(
                          Icons.close,
                          color: ColorName.grey80,
                        ),
                      ),
              ),
            ),
          )
        : Container();
  }

  /// Returns alert default border style
  ShapeBorder _defaultShape() {
    return RoundedRectangleBorder(
      borderRadius: BorderRadius.circular(10.0),
    );
  }

  /// Shows alert with selected animation
  _showAnimation(animation, secondaryAnimation, child) {
    switch (style.animationType) {
      case AnimationType.fromRight:
        return AnimationTransition.fromRight(
            animation, secondaryAnimation, child);
      case AnimationType.fromLeft:
        return AnimationTransition.fromLeft(
            animation, secondaryAnimation, child);
      case AnimationType.fromBottom:
        return AnimationTransition.fromBottom(
            animation, secondaryAnimation, child);
      case AnimationType.grow:
        return AnimationTransition.grow(animation, secondaryAnimation, child);
      case AnimationType.shrink:
        return AnimationTransition.shrink(animation, secondaryAnimation, child);
      case AnimationType.fromTop:
        return AnimationTransition.fromTop(
            animation, secondaryAnimation, child);
    }
  }
}
