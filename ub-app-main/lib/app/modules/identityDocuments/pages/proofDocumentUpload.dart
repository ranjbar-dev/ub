import 'package:basic_utils/basic_utils.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';

import '../../../../generated/colors.gen.dart';
import '../../../../generated/locales.g.dart';
import '../../../../utils/mixins/commonConsts.dart';
import '../../../../utils/mixins/toast.dart';
import '../../../common/components/UBBottomSheetContainer.dart';
import '../../../common/components/UBButton.dart';
import '../../../common/components/UBSelectButton.dart';
import '../../../common/components/UBText.dart';
import '../../../common/components/UBToastOnTap.dart';
import '../../../common/components/UBWrappedButtons.dart';
import '../controllers/identity_documents_controller.dart';
import '../user_profile_model.dart';
import 'widgets/confirmed.dart';
import 'widgets/fresh.dart';
import 'widgets/processing.dart';
import 'widgets/rejected.dart';

class ProofDocumentUpload extends GetView<IdentityDocumentsController>
    with Toaster {
  final String mainType;

  ProofDocumentUpload({@required this.mainType});
  @override
  Widget build(BuildContext context) {
    return Expanded(
      child: Column(
        children: [
          Padding(
            padding: const EdgeInsets.only(
              top: 24,
              right: 12,
              bottom: 24,
              left: 12,
            ),
            child: UBText(
              text: "Choose the document that you want to be verified.",

              // mainType == 'identity'
              //? LocaleKeys.proofOfIdentityTopText.tr
              //: LocaleKeys.proofOfResidenceTopText.tr,
              color: ColorName.grey80,
              lineHeight: 1.4,
              size: 13,
            ),
          ),
          Padding(
            padding: const EdgeInsets.only(bottom: 24),
            child: Obx(
              () {
                final canSelect = _canSelect(mainType);
                final types = controller.subTypes(mainType);
                final selectedSubtypeIndex = _selectedSubTypeIndex(mainType);

                final dropDownTitle = StringUtils.capitalize(
                    types[selectedSubtypeIndex].name.replaceAll('_', ' '),
                    allWords: true);
                return UBToastOnTap(
                  active: !canSelect,
                  toastText: "Type can't be changed after upload",
                  child: UBSelectButton(
                      valueText: dropDownTitle,
                      onClick: () {
                        Get.bottomSheet(
                          UBButtomSheetContainer(
                            title: LocaleKeys.documentType.tr,
                            child: Container(
                              child: UBWrappedButtons(
                                buttons: [
                                  for (var item in types)
                                    WrappedButtonModel(
                                      text: StringUtils.capitalize(
                                          item.name.replaceAll('_', ' '),
                                          allWords: true),
                                    )
                                ],
                                onButtonClick: (v) {
                                  Get.back();
                                  controller.handleSubtypeChange(
                                    v,
                                    mainType,
                                  );
                                },
                                selectedIndex: _selectedSubTypeIndex(mainType),
                              ),
                            ),
                          ),
                        );
                      }),
                );
              },
            ),
          ),
          Padding(
            padding: const EdgeInsets.only(bottom: 24),
            child: Obx(
              () {
                final types = controller.subTypes(mainType);
                final selectedSubtypeIndex = _selectedSubTypeIndex(mainType);
                final hasBack = types[selectedSubtypeIndex].hasBack;
                return AspectRatio(
                  aspectRatio: hasBack == true ? (336 / 122) : (336 / 268),
                  child: _imageState(
                      type: mainType,
                      image: _frontImage(mainType: mainType),
                      subType: types[selectedSubtypeIndex],
                      side: "front"),
                );
              },
            ),
          ),
          Obx(
            () {
              final types = controller.subTypes(mainType);
              final selectedSubtypeIndex = _selectedSubTypeIndex(mainType);
              final hasBack = types[selectedSubtypeIndex].hasBack;
              return hasBack == true
                  ? Padding(
                      padding: const EdgeInsets.only(bottom: 24),
                      child: AspectRatio(
                        aspectRatio: (336 / 122),
                        child: _imageState(
                            type: mainType,
                            image: _backImage(mainType: mainType),
                            subType: types[selectedSubtypeIndex],
                            side: 'back'),
                      ),
                    )
                  : const SizedBox();
            },
          ),
          fill,
          Obx(
            () {
              final canSubmit = controller.canSubmit.value;
              final tabIndex = controller.activeTabIndex.value;
              var hasRejected = false;
              if (tabIndex == 0 && controller.hasRejectedIdentity.value) {
                hasRejected = true;
              }
              if (tabIndex == 1 && controller.hasRejectedAddress.value) {
                hasRejected = true;
              }
              if (hasRejected) {
                return UBButton(
                  onClick: () {
                    controller.resetRejectedImages();
                  },
                  variant: ButtonVariant.Outline,
                  text: LocaleKeys.tryAgain.tr,
                );
              }
              return UBButton(
                onClick: () {
                  if (controller.uploadPercent.value == -1) {
                    controller.submitImages();
                  }
                },
                disabled: !canSubmit,
                text: LocaleKeys.submit.tr,
              );
            },
          ),
          vspace12
        ],
      ),
    );
  }

  Widget _imageState(
      {UserProfileImages image,
      @required SubTypes subType,
      @required String side,
      @required String type}) {
    final state = image.status ?? 'new';
    switch (state) {
      case 'new':
        return Fresh(
          side: side,
          subType: subType,
          type: type,
        );
      case 'processing':
        return Processing(
          image: image,
        );

      case 'confirmed':
      case 'partially_confirmed':
        return Confirmed(
          side: side,
          subType: subType,
        );
      case 'rejected':
        return Rejected(side: side, subType: subType, image: image);

        break;

      default:
        return Container(
          color: Colors.red,
        );
    }
  }

  bool _canSelect(String mainType) {
    if (mainType == 'identity') {
      return controller.canChangeIdentityTypeSelect.value;
    }
    return controller.canChangeResidenceTypeSelect.value;
  }

  int _selectedSubTypeIndex(String mainType) {
    if (mainType == 'identity') {
      return controller.activeIdentitySubTypeindex.value;
    }
    return controller.activeAddressSubTypeindex.value;
  }

  UserProfileImages _frontImage({String mainType}) {
    if (mainType == 'identity') {
      return controller.identityFrontImage.value;
    }
    return controller.addressFrontImage.value;
  }

  UserProfileImages _backImage({String mainType}) {
    if (mainType == 'identity') {
      return controller.identityBackImage.value;
    }
    return controller.addressBackImage.value;
  }
}
