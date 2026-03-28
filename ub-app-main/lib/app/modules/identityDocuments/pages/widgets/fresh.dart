import 'package:basic_utils/basic_utils.dart';
import 'package:file_picker/file_picker.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';
import '../../../../common/components/UBDottedBorder.dart';
import '../../../../common/components/UBRoundedButton.dart';
import '../../../../common/components/UBText.dart';
import '../../../../common/custom/rflutter_alert/src/alert.dart';
import '../../../../common/custom/rflutter_alert/src/alert_style.dart';
import '../../../../common/custom/rflutter_alert/src/constants.dart';
import '../../controllers/identity_documents_controller.dart';
import '../../user_profile_model.dart';
import '../../../../../generated/assets.gen.dart';
import '../../../../../generated/colors.gen.dart';
import 'package:percent_indicator/percent_indicator.dart';

class Fresh extends GetView<IdentityDocumentsController> {
  const Fresh({
    Key key,
    @required this.side,
    @required this.subType,
    @required this.type,
  }) : super(key: key);
  final String side;
  final SubTypes subType;
  final String type;
  @override
  Widget build(BuildContext context) {
    return UBDottedBorder(
      child: GestureDetector(
        onTap: () {
          controller.handleSelectFile(
            type: type,
            side: side,
          );
        },
        child: Container(
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(12),
            color: ColorName.primaryBlue.withOpacity(0.15),
          ),
          child: Obx(
            () {
              final image = _freshImage(type: type, side: side);
              return image.size == 0
                  ? Center(
                      child: Column(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          UBText(
                            text:
                                "Tap to upload $side side of your ${StringUtils.capitalize(subType.name.replaceAll('_', ' '), allWords: true)}",
                            color: ColorName.primaryBlue,
                          ),
                          const SizedBox(
                            height: 12,
                          ),
                        ],
                      ),
                    )
                  : Container(
                      child: ClipRRect(
                        borderRadius: BorderRadius.circular(11),
                        child: Stack(
                          children: [
                            Container(
                              width: double.infinity,
                              height: double.infinity,
                              child: Image.memory(
                                image.bytes,
                                fit: BoxFit.cover,
                              ),
                            ),
                            Positioned(
                              bottom: 12,
                              right: 12,
                              child: UBRoundButton(
                                child: Assets.images.expand.svg(),
                                onClick: () {
                                  _openImage(image: image, context: context);
                                },
                              ),
                            ),
                            Positioned(
                              bottom: 12,
                              left: 12,
                              child: UBRoundButton(
                                child: Assets.images.trash.svg(),
                                onClick: () {
                                  _deleteImage(side: side, type: type);
                                },
                              ),
                            ),
                            if (controller.uploadPercent.value != -1)
                              Container(
                                width: double.infinity,
                                height: double.infinity,
                                color: ColorName.black.withOpacity(0.8),
                                child: Column(
                                  mainAxisAlignment: MainAxisAlignment.center,
                                  children: [
                                    CircularPercentIndicator(
                                      radius: 64.0,
                                      lineWidth: 4.0,
                                      percent:
                                          controller.uploadPercent.value / 100,
                                      center: UBRoundButton(
                                        size: 24,
                                        child: const Icon(
                                          Icons.close,
                                          size: 12.0,
                                          color: ColorName.grey80,
                                        ),
                                        onClick: () {
                                          controller.cancelUpload();
                                        },
                                      ),
                                      circularStrokeCap:
                                          CircularStrokeCap.round,
                                      backgroundColor: ColorName.black2c,
                                      progressColor: ColorName.primaryBlue,
                                    ),
                                    const SizedBox(
                                      height: 8,
                                    ),
                                    UBText(
                                      text: "${controller.uploadPercent} %",
                                      color: ColorName.grey80,
                                    )
                                  ],
                                ),
                              )
                          ],
                        ),
                      ),
                    );
            },
          ),
        ),
      ),
    );
  }

  PlatformFile _freshImage({String type, String side}) {
    if (type == 'identity') {
      if (side == 'front') {
        return controller.identityFrontImageUploadFile.value;
      } else {
        return controller.identityBackImageUploadFile.value;
      }
    } else {
      if (side == 'front') {
        return controller.addressFrontImageUploadFile.value;
      } else {
        return controller.addressBackImageUploadFile.value;
      }
    }
  }

  void _openImage({PlatformFile image, BuildContext context}) {
    Alert(
      style: AlertStyle(
        animationType: AnimationType.grow,
      ),
      context: context,
      content: Container(
        child: InteractiveViewer(
          panEnabled: true, // Set it to false to prevent panning.
          boundaryMargin: const EdgeInsets.all(80),
          minScale: 0.5,
          maxScale: 4,
          child: Image.memory(
            image.bytes,
            fit: BoxFit.cover,
          ),
        ),
      ),
    ).show();
  }

  void _deleteImage({String side, String type}) {
    if (type == 'identity') {
      if (side == 'front') {
        controller.identityFrontImageUploadFile.value =
            PlatformFile(name: '', size: 0);
        if (controller.identityBackImageUploadFile.value.path == null) {
          controller.canSubmit.value = false;
        }
      } else {
        controller.identityBackImageUploadFile.value =
            PlatformFile(name: '', size: 0);
        if (controller.identityFrontImageUploadFile.value.path == null) {
          controller.canSubmit.value = false;
        }
      }
    } else {
      if (side == 'front') {
        controller.addressFrontImageUploadFile.value =
            PlatformFile(name: '', size: 0);
        if (controller.addressBackImageUploadFile.value.path == null) {
          controller.canSubmit.value = false;
        }
      } else {
        controller.addressBackImageUploadFile.value =
            PlatformFile(name: '', size: 0);
        if (controller.addressFrontImageUploadFile.value.path == null) {
          controller.canSubmit.value = false;
        }
      }
    }
  }
}
