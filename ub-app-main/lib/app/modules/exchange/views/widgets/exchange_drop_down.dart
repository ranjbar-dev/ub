import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../../../generated/colors.gen.dart';
import '../../../../../utils/mixins/commonConsts.dart';
import '../../../../common/components/UBCircularImage.dart';
import '../../../../common/components/UBText.dart';

class ExchangeDropDown extends StatelessWidget {
  const ExchangeDropDown(
      {Key key,
      this.name,
      this.imageUrl,
      this.isFrom,
      this.desc,
      this.autoExchangeCode,
      this.isDisabled = false})
      : super(key: key);

  final String name;
  final String desc;
  final String imageUrl;
  final bool isFrom;
  final bool isDisabled;
  final String autoExchangeCode;

  @override
  Widget build(BuildContext context) {
    return Container(
      height: 40,
      width: (MediaQuery.of(context).size.width / 5) * 2,
      decoration: BoxDecoration(
        borderRadius: BorderRadius.circular(6),
        color: isDisabled ? ColorName.black1c : ColorName.black,
      ),
      child: InkWell(
        child: Padding(
          padding: const EdgeInsets.symmetric(horizontal: 4.0),
          child: Row(
            mainAxisAlignment: MainAxisAlignment.spaceEvenly,
            children: [
              GetPlatform.isWeb
                  ? UBCircularImage(
                      imageAddress: imageUrl,
                    )
                  : isDisabled
                      ? ColorFiltered(
                          colorFilter: ColorFilter.mode(
                              ColorName.black1c, BlendMode.saturation),
                          child: UBCircularImage(
                            imageAddress: imageUrl,
                          ),
                        )
                      : UBCircularImage(
                          imageAddress: imageUrl,
                        ),
              hspace2,
              UBText(
                text: name,
                size: 13,
                color: ColorName.white,
              ),
              hspace4,
              Flexible(
                fit: FlexFit.loose,
                child: Text(
                  '($desc)',
                  softWrap: false,
                  overflow: TextOverflow.clip,
                  style: TextStyle(
                    fontSize: 9,
                    color: ColorName.grey80,
                  ),
                ),
                // child: UBText(
                //   size: 9,
                //   text: '($desc)',
                //   color: ColorName.grey80,
                // ),
              ),
              Spacer(),
              Icon(
                Icons.keyboard_arrow_down,
                color: ColorName.grey97,
                size: 22,
              )
            ],
          ),
        ),
        onTap: () => isDisabled
            ? null
            : Get.toNamed("/exchange-search",
                arguments: {
                  'isFrom': isFrom,
                  'autoExchangeCode': autoExchangeCode
                },
                preventDuplicates: false),
      ),
    );
  }
}
