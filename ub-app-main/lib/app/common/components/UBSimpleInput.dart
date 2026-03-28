import 'package:flutter/material.dart';

import '../../../generated/assets.gen.dart';
import '../../../generated/colors.gen.dart';
import '../../../utils/mixins/commonConsts.dart';

class UBSimpleInput extends StatefulWidget {
  final double errorTopPosition;
  final String placeHolder;
  final Function(String v) onChange;
  final bool isSecure;
  final String error;
  final bool isPickable;
  final double height;
  final TextInputType type;
  const UBSimpleInput({
    Key key,
    @required this.placeHolder,
    @required this.onChange,
    this.isSecure,
    this.error,
    this.isPickable,
    this.height = 42.0,
    this.errorTopPosition,
    this.type = TextInputType.text,
  }) : super(key: key);

  @override
  _UBSimpleInputState createState() => _UBSimpleInputState();
}

class _UBSimpleInputState extends State<UBSimpleInput> {
  bool isFocused = false;
  bool isVisible = true;

  @override
  void initState() {
    if (widget.isSecure == true) {
      setState(() {
        isVisible = false;
      });
    }
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    final errorTopPosition = widget.height + 1.0;
    var focusColor = (widget.error != null && widget.error != '')
        ? ColorName.red
        : ColorName.primaryBlue;
    var normalColor = (widget.error != null && widget.error != '')
        ? ColorName.red
        : ColorName.black2c;
    return Stack(
      clipBehavior: Clip.none,
      children: [
        Focus(
          onFocusChange: (e) => {_handleOnFocusChange(e)},
          child: Container(
            height: widget.height,
            child: TextFormField(
              onChanged: widget.onChange,
              obscureText: !isVisible,
              style: const TextStyle(
                fontSize: 16.0,
                fontWeight: FontWeight.w600,
                color: ColorName.lightText,
              ),
              decoration: InputDecoration(
                border: const OutlineInputBorder(
                  borderRadius: const BorderRadius.all(
                    const Radius.circular(6.0),
                  ),
                ),
                suffixIcon: widget.isPickable == true
                    ? GestureDetector(
                        onTap: () => {
                          setState(() {
                            isVisible = !isVisible;
                          })
                        },
                        child: Container(
                          child: Assets.images.eye.svg(
                            fit: BoxFit.scaleDown,
                            color: isVisible
                                ? ColorName.primaryBlue
                                : ColorName.greybf,
                          ),
                        ),
                      )
                    : null,
                filled: true,
                fillColor: ColorName.black,
                contentPadding: const EdgeInsets.only(
                  left: 14.0,
                  bottom: 6.0,
                  top: 8.0,
                ),
                hintText: widget.placeHolder,
                labelStyle: TextStyle(
                  color: isFocused ? focusColor : ColorName.greybf,
                ),
                focusedBorder: OutlineInputBorder(
                  borderRadius: rounded6,
                  borderSide: BorderSide(
                    color: focusColor,
                    width: 1.0,
                  ),
                ),
                enabledBorder: OutlineInputBorder(
                  borderRadius: rounded6,
                  borderSide: BorderSide(
                    color: normalColor,
                    width: 1.0,
                  ),
                ),
              ),
              keyboardType: widget.type,
            ),
          ),
        ),
        ...((widget.error != null && widget.error != '')
            ? [
                Positioned(
                  top: errorTopPosition,
                  left: 4,
                  child: Text(
                    widget.error,
                    style: const TextStyle(
                      color: ColorName.red,
                      fontSize: 10,
                    ),
                  ),
                )
              ]
            : [])
      ],
    );
  }

  _handleOnFocusChange(bool e) {
    if (e) {
      setState(() {
        isFocused = true;
      });
    } else {
      setState(() {
        isFocused = false;
      });
    }
  }
}
