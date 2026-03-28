import 'package:basic_utils/basic_utils.dart';
import 'package:flutter/material.dart';
import '../../../../../common/components/UBText.dart';
import '../transaction_history_model.dart';
import '../../../../../../generated/assets.gen.dart';
import '../../../../../../generated/colors.gen.dart';
import '../../../../../../utils/mixins/commonConsts.dart';
//import 'package:unitedbit/utils/mixins/commonConsts.dart';

class TransactionHistoryRow extends StatelessWidget {
  final Payments data;
  final Function(Payments data) onTap;

  const TransactionHistoryRow({Key key, this.data, this.onTap})
      : super(key: key);
  @override
  Widget build(BuildContext context) {
    final isCompleted = data.status == 'completed';
    final isRejected = data.status == 'rejected';
    final isPending = data.status == 'pending';
    final statusColor = isCompleted
        ? ColorName.green
        : isRejected
            ? ColorName.red
            : isPending
                ? ColorName.orange
                : ColorName.grey80;
    //final tappable = data.txIdExplorerUrl != null && data.txIdExplorerUrl != '';
    return GestureDetector(
      onTap: () {
        onTap(data);
      },
      child: Column(
        children: [
          Container(
            height: 46,
            padding: const EdgeInsets.symmetric(
              horizontal: 12,
            ),
            decoration: const BoxDecoration(
              border: const Border(
                bottom: const BorderSide(
                  width: 1,
                  color: ColorName.grey16,
                ),
              ),
            ),
            child: Row(
              children: [
                Container(
                  child: Column(
                    mainAxisAlignment: MainAxisAlignment.center,
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          SizedBox(
                            width: 16,
                            child: data.type == 'deposit'
                                ? Assets.images.greenDownArrow.svg()
                                : Assets.images.redUpArrow.svg(),
                          ),
                          hspace4,
                          UBText(
                            text: "${data.code} ${data.amount}",
                            color: ColorName.white,
                          ),
                        ],
                      ),
                      vspace4,
                      Row(
                        children: [
                          hspace16,
                          hspace4,
                          UBText(
                            text: "${data.createdAt}",
                            color: ColorName.grey80,
                            size: 9,
                          ),
                        ],
                      ),
                    ],
                  ),
                ),
                fill,
                Container(
                  child: Row(
                    children: [
                      if (isPending)
                        SizedBox(
                          width: 9,
                          height: 9,
                          child: CircularProgressIndicator(
                            valueColor:
                                AlwaysStoppedAnimation<Color>(ColorName.grey42),
                            strokeWidth: 1,
                            backgroundColor: Colors.transparent,
                          ),
                        ),
                      hspace8,
                      UBText(
                        text: StringUtils.capitalize(data.status),
                        size: 11,
                        color: statusColor,
                      ),
                      //if (tappable)
                      Container(
                          alignment: Alignment.center,
                          width: 24,
                          height: 30,
                          child: Icon(
                            Icons.keyboard_arrow_right_rounded,
                            color: ColorName.grey80,
                            size: 16.0,
                          ))
                      //else
                      //  hspace24
                    ],
                  ),
                ),
              ],
            ),
          )
        ],
      ),
    );
  }
}
