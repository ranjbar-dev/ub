import 'package:flutter/material.dart';

import 'package:get/get.dart';
import '../../../common/components/CenterUBLoading.dart';
import '../../../common/components/appbarTextTitle.dart';
import '../../../common/components/pageContainer.dart';
import '../../../common/custom/flutter_advanced_segment/flutter_advanced_segment.dart';
import '../pages/proofDocumentUpload.dart';
import '../../../../generated/locales.g.dart';

import '../controllers/identity_documents_controller.dart';

class IdentityDocumentsView extends GetView<IdentityDocumentsController> {
  final segmentController = AdvancedSegmentController('poi');
  @override
  Widget build(BuildContext context) {
    if (controller.isLoading.value == false) {
      controller.getUserProfile(silent: true);
    }

    return PageContainer(
      appbarTitle: AppBarTextTitle(
        title: LocaleKeys.identityVerification.tr,
      ),
      child: Container(
        width: Get.width,
        padding: const EdgeInsets.only(left: 12, right: 12, top: 12),
        child: Obx(() {
          final isLoading = controller.isLoading.value;
          return isLoading == true
              ? Container(
                  child: CenterUbLoading(),
                  height: Get.height - 60,
                )
              : Column(
                  children: [
                    Container(
                      height: 36,
                      alignment: Alignment.center,
                      child: AdvancedSegment(
                        itemWidth: Get.width / 2 - 12,
                        onChange: controller.handleTabChange,
                        segments: {
                          'poi': LocaleKeys.poi.tr,
                          'por': LocaleKeys.por.tr,
                        },
                        controller: segmentController,
                      ),
                    ),
                    Obx(
                      () => controller.activeTabIndex.value == 0
                          ? ProofDocumentUpload(
                              mainType: 'identity',
                            )
                          : ProofDocumentUpload(
                              mainType: 'address',
                            ),
                    )
                  ],
                );
        }),
      ),
    );
  }
}
