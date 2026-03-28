import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'autoCompleteList.dart';
import '../../global/autocompleteModel.dart';
import '../../../generated/locales.g.dart';
import '../../../services/constants.dart';

class CountryAutoComplete extends StatelessWidget {
  final countryArray = Constants.countriesArray();
  final void Function(AutoCompleteItem item) onSelect;
  CountryAutoComplete({Key key, this.onSelect}) : super(key: key);
  @override
  Widget build(BuildContext context) {
    return AutoCompleteList(
      itemHeight: 36,
      title: LocaleKeys.selectCountry.tr,
      itemList: countryArray,
      onItemSelect: onSelect,
    );
  }
}
