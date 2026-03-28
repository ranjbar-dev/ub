import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:get_storage/get_storage.dart';
import 'UBSection.dart';
import 'UBText.dart';
import 'UBWrappedButtons.dart';
import '../../global/autocompleteModel.dart';
import '../../../generated/assets.gen.dart';
import '../../../generated/colors.gen.dart';
import '../../../utils/mixins/popups.dart';

class CoinSearchHistory extends StatelessWidget with Popups {
  final RxList<AutoCompleteItem> stream;
  final String storageKey;
  final Function(AutoCompleteItem coin) onCoinClick;
  const CoinSearchHistory(
      {Key key, this.stream, this.storageKey, this.onCoinClick})
      : super(key: key);
  @override
  Widget build(BuildContext context) {
    // ignore: invalid_use_of_protected_member
    final searchedCoins = stream.value;
    return searchedCoins.length > 0
        ? UBSection(
            hTitlePadding: 12.0,
            title: 'Search History',
            titleEndWidget: GestureDetector(
              onTap: _clearCoinSearchHistory,
              child: Container(
                width: 32.0,
                height: 32.0,
                color: Colors.transparent,
                child: Assets.images.trashIcon.svg(),
              ),
            ),
            child: Padding(
              padding: const EdgeInsets.symmetric(horizontal: 12.0),
              child: UBWrappedButtons(
                otherNoMarginLeftIndexes: [4],
                unselectedButtonTextColor: ColorName.greybf,
                minButtonWidth: (Get.width / 4) - 20,
                buttons: [
                  for (var item in searchedCoins)
                    WrappedButtonModel(text: item.code)
                ],
                onButtonClick: (index) {
                  // ignore: invalid_use_of_protected_member
                  onCoinClick(stream.value[index]);
                },
                selectedIndex: -1,
              ),
            ),
          )
        : const SizedBox();
  }

  void _clearCoinSearchHistory() {
    final storage = GetStorage();
    openConfirmation(
        onConfirm: () {
          storage.remove(storageKey);
          stream.assignAll([]);
        },
        titleWidget: UBText(
          text: 'Clear search history?',
        ),
        confirmText: 'Clear');
  }
}
