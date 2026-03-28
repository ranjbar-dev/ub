import 'package:flutter/material.dart';
import '../../../generated/colors.gen.dart';

class UBRawInput extends StatelessWidget {
  final String placeHolder;
  final TextInputType type;
  const UBRawInput({Key key, this.placeHolder, this.type}) : super(key: key);
  @override
  Widget build(BuildContext context) {
    return TextFormField(
      keyboardType: type ?? TextInputType.text,
      style: const TextStyle(
        fontSize: 12,
        color: ColorName.lightText,
      ),
      decoration: InputDecoration(
          focusedBorder: const OutlineInputBorder(
            borderSide: const BorderSide(
              color: Colors.transparent,
              width: 1.0,
            ),
          ),
          enabledBorder: const OutlineInputBorder(
            borderSide: const BorderSide(
              color: Colors.transparent,
              width: 1.0,
            ),
          ),
          contentPadding: const EdgeInsets.only(
            bottom: 8.0,
            top: 8.0,
          ),
          hintText: placeHolder),
    );
  }
}
