import 'package:flutter/material.dart';

class UBCheckBox extends StatefulWidget {
  const UBCheckBox({
    Key key,
  }) : super(key: key);

  @override
  _UBCheckBoxState createState() => _UBCheckBoxState();
}

class _UBCheckBoxState extends State<UBCheckBox> {
  @override
  Widget build(BuildContext context) {
    return Checkbox(
      value: true,
      onChanged: (e) => {},
    );
  }
}
