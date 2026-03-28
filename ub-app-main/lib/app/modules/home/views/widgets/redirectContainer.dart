import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../../../generated/assets.gen.dart';
import '../../../../../generated/colors.gen.dart';
import '../../../../global/controller/globalController.dart';

class RedirectContainer extends GetWidget<GlobalController> {
  const RedirectContainer({Key key}) : super(key: key);

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.only(left: 8, right: 8, top: 1, bottom: 67),
      child: Container(
        padding: const EdgeInsets.all(8),
        decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(10), color: ColorName.black2c),
        child: Row(
          mainAxisAlignment: MainAxisAlignment.start,
          children: [
            Container(
              width: 50,
              height: 50,
              child: Image(
                image: Assets.images.logoPng,
              ),
            ),
            SizedBox(
              width: 16,
            ),
            Column(
              mainAxisAlignment: MainAxisAlignment.center,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  'Download Unitedbit App',
                  style: const TextStyle(
                      fontSize: 14,
                      color: Colors.white,
                      fontWeight: FontWeight.bold),
                ),
                SizedBox(
                  height: 5,
                ),
                Text(
                  'Safe, Diverse & Fast',
                  textAlign: TextAlign.start,
                  style: const TextStyle(
                    fontSize: 13,
                    color: Colors.white,
                  ),
                )
              ],
            ),
            Spacer(),
            CircleAvatar(
              backgroundColor: ColorName.primaryBlue,
              radius: 20,
              child: IconButton(
                padding: EdgeInsets.zero,
                icon: Assets.images.downloadIcon.svg(),
                color: Colors.white,
                onPressed: () => controller.doRedirect(),
              ),
            ),
            Container(
              width: 20,
              child: IconButton(
                onPressed: () {
                  controller.setRedirectConteinerDismissed(true);
                },
                icon: Icon(
                  Icons.close,
                  color: ColorName.grey80,
                  size: 15,
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}
