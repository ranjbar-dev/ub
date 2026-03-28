import 'package:flutter/material.dart';
import '../../../common/components/UBCircularImage.dart';
import '../../../common/components/UBText.dart';
import '../withdraw_address_model.dart';
import '../../../../generated/assets.gen.dart';
import '../../../../generated/colors.gen.dart';
import '../../../../utils/commonUtils.dart';

class WithdrawAddressRow extends StatelessWidget {
  final bool onlySelectable;
  final WithdrawAddressModel item;
  final void Function() onDeleteClick;
  const WithdrawAddressRow({
    Key key,
    @required this.item,
    this.onlySelectable,
    this.onDeleteClick,
  }) : super(key: key);
  @override
  Widget build(BuildContext context) {
    return Container(
      height: 68,
      decoration: const BoxDecoration(
        border: const Border(
          top: const BorderSide(
            width: 1,
            color: ColorName.grey16,
          ),
        ),
      ),
      padding: const EdgeInsets.only(
        top: 9,
        bottom: 12,
        left: 12,
        right: 12,
      ),
      child: Row(
        children: [
          Expanded(
            child: Container(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Container(
                    child: Row(children: [
                      UBCircularImage(
                        size: 28,
                        imageAddress: item.icon,
                      ),
                      UBText(
                        text: item.label,
                        color: ColorName.greybf,
                      ),
                      if (item.network != null && item.network != '')
                        UBText(
                          text: " (Network:${item.network})",
                          color: ColorName.grey97,
                          size: 10.0,
                        ),
                    ]),
                  ),
                  Padding(
                    padding: const EdgeInsets.only(left: 8.0),
                    child: UBText(
                      text: censorAddress(address: item.address),
                      color: ColorName.grey80,
                    ),
                  )
                ],
              ),
            ),
          ),
          if (onlySelectable != true)
            GestureDetector(
              onTap: () => onDeleteClick(),
              child: Container(
                alignment: Alignment.centerRight,
                width: 24,
                height: 24,
                color: Colors.transparent,
                child: Assets.images.trashIcon.svg(),
              ),
            )
          else
            const SizedBox(width: 24, height: 24)
        ],
      ),
    );
  }
}
