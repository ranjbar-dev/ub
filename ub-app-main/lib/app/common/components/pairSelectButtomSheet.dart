import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../generated/colors.gen.dart';
import '../../global/autocompleteModel.dart';
import 'pairsAutocomplete.dart';

class SelectPairBottomSheet extends StatelessWidget {
  final void Function(AutoCompleteItem item) onSelect;
  final bool closeOnSelect;
  const SelectPairBottomSheet({
    Key key,
    this.onSelect,
    this.closeOnSelect,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      child: Container(
        margin: const EdgeInsets.only(top: 24, left: 12, right: 12),
        decoration: const BoxDecoration(
          color: ColorName.grey16,
          borderRadius: const BorderRadius.only(
            topLeft: const Radius.circular(12),
            topRight: const Radius.circular(12),
          ),
        ),
        child: Column(
          children: <Widget>[
            Expanded(
              child: PairsAutoComplete(
                onSelect: (item) {
                  onSelect(item);
                  if (closeOnSelect == true) {
                    Get.back();
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
