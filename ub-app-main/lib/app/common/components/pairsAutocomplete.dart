import 'package:flutter/material.dart';
import 'autoCompleteList.dart';
import '../../global/autocompleteModel.dart';
import '../../../services/constants.dart';

class PairsAutoComplete extends StatelessWidget {
  final pairsArray = Constants.pairsAutoCompleteArray();
  final void Function(AutoCompleteItem item) onSelect;
  PairsAutoComplete({Key key, this.onSelect}) : super(key: key);
  @override
  Widget build(BuildContext context) {
    return AutoCompleteList(
      itemHeight: 36,
      title: 'Select Pair',
      itemList: pairsArray,
      onItemSelect: onSelect,
    );
  }
}
