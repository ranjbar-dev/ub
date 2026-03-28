import 'package:flutter/material.dart';
import 'countryAutocomplete.dart';
import '../../global/autocompleteModel.dart';
import '../../../generated/colors.gen.dart';

class SelectCountryBottomSheet extends StatelessWidget {
  final void Function(AutoCompleteItem item) onSelect;
  final bool closeOnSelect;
  const SelectCountryBottomSheet({
    Key key,
    this.onSelect,
    this.closeOnSelect,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      child: Container(
        margin: EdgeInsets.only(top: 24, left: 12, right: 12),
        decoration: BoxDecoration(
          color: ColorName.grey16,
          borderRadius: const BorderRadius.only(
            topLeft: const Radius.circular(12),
            topRight: const Radius.circular(12),
          ),
        ),
        child: Column(
          children: <Widget>[
            Expanded(
              child: CountryAutoComplete(
                onSelect: (item) {
                  onSelect(item);
                  if (closeOnSelect == true) {
                    Navigator.pop(context);
                  }
                  return;
                },
              ),
            )
          ],
        ),
      ),
    );
  }
}
