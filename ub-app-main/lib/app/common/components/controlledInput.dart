import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:get/get.dart';
import '../../../generated/colors.gen.dart';

class ControlledTextField extends StatefulWidget {
  final String text;
  final String labelText;
  final String placeholder;
  final void Function(String text) onChanged;
  final TextInputType type;
  final bool noBorder;
  final bool autoFocus;
  final Function(bool) onFocusChanged;
  final bool isCurrencyInput;
  final List<TextInputFormatter> formatters;
  final TextStyle textStyle;
  final int maxLength;

  const ControlledTextField({
    Key key,
    this.text,
    this.labelText,
    this.onChanged,
    this.noBorder,
    this.type,
    this.textStyle = const TextStyle(
        fontSize: 13.0,
        color: ColorName.lightText,
        fontWeight: FontWeight.w600),
    this.placeholder,
    this.isCurrencyInput = false,
    this.formatters,
    this.autoFocus = false,
    this.onFocusChanged,
    this.maxLength,
  }) : super(key: key);

  @override
  _ControlledTextFieldState createState() => _ControlledTextFieldState();
}

class _ControlledTextFieldState extends State<ControlledTextField> {
  final _focusNode = FocusNode();
  bool isFocused = false;
  final isWeb = GetPlatform.isWeb;
  final _textEditingController = TextEditingController();

  @override
  void initState() {
    super.initState();
    _textEditingController.value =
        _textEditingController.value.copyWith(text: widget.text);
    _focusNode.addListener(() {
      if (widget.onFocusChanged != null) {
        widget.onFocusChanged(_focusNode.hasFocus);
      }
      setState(() {
        isFocused = _focusNode.hasFocus;
      });
    });
  }

  @override
  void dispose() {
    _textEditingController.dispose();
    _focusNode.dispose();

    super.dispose();
  }

  @override
  void didUpdateWidget(ControlledTextField oldWidget) {
    if (oldWidget != this.widget) {
      if (this.widget.text != null &&
          this.widget.text != this._textEditingController.value.text) {
        this._textEditingController.text = this.widget.text;

        _textEditingController.selection =
            TextSelection.collapsed(offset: this.widget.text.length);
      }
    }
    super.didUpdateWidget(oldWidget);
  }

  @override
  Widget build(BuildContext context) {
    return TextField(
      autofocus: this.widget.autoFocus,
      textAlignVertical: TextAlignVertical.center,
      controller: _textEditingController,
      focusNode: _focusNode,
      onChanged: (v) {
        this.widget.onChanged(v);
      },
      keyboardType: this.widget.type,
      inputFormatters: this.widget.formatters,
      style: this.widget.textStyle,
      maxLength: this.widget.maxLength,
      buildCounter: (BuildContext context,
              {int currentLength, int maxLength, bool isFocused}) =>
          null,
      decoration: widget.noBorder != true
          ? InputDecoration(
              counter: const Offstage(),
              hintText: (isFocused && isWeb) ? '' : this.widget.labelText,
            )
          : InputDecoration.collapsed(
              hintText: (isFocused && isWeb) ? '' : widget.labelText),
    );
  }
}
