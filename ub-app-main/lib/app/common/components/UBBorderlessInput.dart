import 'package:flutter/material.dart';
import '../../../generated/colors.gen.dart';

class UBBorderlessInput extends StatelessWidget {
  final String placeholder;
  final TextInputType type;
  final bool autoFocus;
  final void Function(String) onChange;
  final TextEditingController controller;
  final double fontSize;

  const UBBorderlessInput({
    Key key,
    this.placeholder,
    this.type = TextInputType.text,
    @required this.onChange,
    this.controller,
    this.fontSize = 13.0,
    this.autoFocus = false,
  }) : super(key: key);
  @override
  Widget build(BuildContext context) {
    return TextFormField(
      controller: controller,
      autofocus: autoFocus,
      onChanged: onChange,
      keyboardType: type,
      style: TextStyle(
        fontSize: fontSize,
        color: ColorName.lightText,
        fontWeight: FontWeight.w600,
      ),
      decoration: InputDecoration(
        focusedBorder: const OutlineInputBorder(
          borderSide: const BorderSide(
            color: Colors.transparent,
            width: 0.0,
          ),
        ),
        enabledBorder: const OutlineInputBorder(
          borderSide: const BorderSide(
            color: Colors.transparent,
            width: 0.0,
          ),
        ),
        contentPadding: const EdgeInsets.only(
          bottom: 8.0,
          top: 8.0,
        ),
        hintText: placeholder,
      ),
    );
  }
}
