import 'package:flutter/material.dart';
import '../../../common/components/roundedCard.dart';
import '../components/accountRow.dart';

import '../accountRowModel.dart';

class AccountCard extends StatelessWidget {
  final List<AccountCardRowModel> rows;
  const AccountCard({
    Key key,
    this.rows,
  }) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return RoundedCard(
      child: Column(
        children: [
          for (final item in rows)
            AccountRow(
              icon: item.icon.svg(
                fit: BoxFit.fitHeight,
              ),
              title: item.title,
              value: item.value,
              endButtonTitle: item.endButtonTitle,
              inactiveEndButtonTitle: item.inactiveEndButtonTitle,
              endButtonActive: item.endButtonActive,
              isLast: item.isLast,
              endWidget: item.endWidget,
              endButtonOnClick: item.onEndButtonClick,
            ),
        ],
      ),
    );
  }
}
