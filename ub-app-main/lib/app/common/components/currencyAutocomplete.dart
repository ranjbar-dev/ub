import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'autoCompleteList.dart';
import '../../global/autocompleteModel.dart';
import '../../../generated/locales.g.dart';
import '../../../services/constants.dart';

class CurrencyAutoComplete extends StatelessWidget {
  final currencyArray = Constants.currencyArray();
  final void Function(AutoCompleteItem item) onSelect;
  CurrencyAutoComplete({Key key, this.onSelect}) : super(key: key);
  @override
  Widget build(BuildContext context) {
    return AutoCompleteList(
      title: LocaleKeys.selectCoin.tr,
      itemList: currencyArray,
      onItemSelect: onSelect,
    );
  }
}
