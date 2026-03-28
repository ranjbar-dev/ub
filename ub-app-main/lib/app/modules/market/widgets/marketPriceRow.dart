import 'package:flutter/widgets.dart';
import '../../../global/currency_pairs_model.dart';
import '../../../../generated/colors.gen.dart';
import '../../../../utils/mixins/commonConsts.dart';

final green = ColorName.green;
final red = ColorName.red;
final grey = ColorName.greybf;
final percentRed = red.withOpacity(0.3);
final percentGreen = green.withOpacity(0.3);

class MarketPriceRow extends StatelessWidget {
  final Pairs data;
  final String price;
  final String equivalentPrice;
  final String volume;
  final Function(String) onClick;
  final double firstColWidth;

  const MarketPriceRow(
      {Key key,
      this.data,
      this.price,
      this.equivalentPrice,
      this.volume,
      this.onClick,
      this.firstColWidth})
      : super(key: key);

  @override
  Widget build(BuildContext context) {
    final percent = data.percent;
    final isRed = percent.contains('-');
    final splitted = data.pairName.split('-');
    final pairName1 = splitted[0];
    final pairName2 = splitted[1];
    return GestureDetector(
      onTap: () => onClick(data.pairName),
      child: Container(
        height: 41,
        decoration: rowDecoration,
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 12.0),
          child: Row(
            children: [
              SizedBox(
                width: firstColWidth ?? thirdWidthPlus24,
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    RichText(
                      text: TextSpan(
                        children: <TextSpan>[
                          TextSpan(
                            text: pairName1,
                            style: whiteBold14,
                          ),
                          TextSpan(
                            text: ' / ',
                            style: whiteBold10,
                          ),
                          TextSpan(
                            text: pairName2,
                            style: grey97Bold10,
                          ),
                        ],
                      ),
                    ),
                    vspace2,
                    RichText(
                      text: TextSpan(
                        children: <TextSpan>[
                          TextSpan(
                            text: 'Vol ',
                            style: grey97Bold10,
                          ),
                          TextSpan(
                            text: volume,
                            style: grey97Bold10,
                          ),
                        ],
                      ),
                    )
                  ],
                ),
              ),
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  RichText(
                    text: TextSpan(
                      text: price,
                      style: TextStyle(
                          color: isRed ? red : green,
                          fontWeight: FontWeight.w700,
                          fontSize: 14),
                    ),
                  ),
                  vspace2,
                  RichText(
                    text: TextSpan(
                      children: <TextSpan>[
                        TextSpan(
                          text: "\$$equivalentPrice",
                          style: grey97Bold10,
                        ),
                      ],
                    ),
                  )
                ],
              ),
              const Spacer(),
              Container(
                width: 60,
                height: 24,
                alignment: Alignment.center,
                decoration: BoxDecoration(
                  color: percent == '0.0'
                      ? ColorName.grey97.withOpacity(0.2)
                      : isRed
                          ? percentRed
                          : percentGreen,
                  borderRadius: const BorderRadius.all(
                    const Radius.circular(4),
                  ),
                ),
                child: RichText(
                  text: TextSpan(
                    text: "$percent%",
                    style: TextStyle(
                      color: percent == '0.0'
                          ? ColorName.grey97
                          : isRed
                              ? red
                              : green,
                      fontWeight: FontWeight.w600,
                      fontSize: 12,
                    ),
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}
