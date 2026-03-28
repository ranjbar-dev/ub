import 'package:dio/dio.dart';
import 'package:file_picker/file_picker.dart';
import 'package:flutter/widgets.dart';
import 'package:get/get.dart' hide MultipartFile, FormData;

import '../../../../utils/logger.dart';
import '../../../../utils/mixins/toast.dart';
import '../providers/identityDocumentsProvider.dart';
import '../user_profile_model.dart';

class IdentityDocumentsController extends GetxController with Toaster {
  final IdentityDocumentsProvider identityDocumentsProvider =
      IdentityDocumentsProvider();
  final canChangeIdentityTypeSelect = true.obs;
  final canChangeResidenceTypeSelect = true.obs;
  final isLoading = true.obs;
  final canSubmit = false.obs;
  final activeIdentitySubTypeindex = 0.obs;
  final activeAddressSubTypeindex = 0.obs;
  final Rx<UserProfileModel> userProfile = UserProfileModel().obs;
  final activeTabIndex = 0.obs;
  List<SubTypes> identitySubTypes;

  final identityFrontImage = UserProfileImages().obs;
  final identityBackImage = UserProfileImages().obs;
  final addressFrontImage = UserProfileImages().obs;
  final addressBackImage = UserProfileImages().obs;

  final identityFrontImageUploadFile = PlatformFile(name: '', size: 0).obs;
  final identityBackImageUploadFile = PlatformFile(name: '', size: 0).obs;
  final addressFrontImageUploadFile = PlatformFile(name: '', size: 0).obs;
  final addressBackImageUploadFile = PlatformFile(name: '', size: 0).obs;

  final hasRejectedIdentity = false.obs;
  final hasRejectedAddress = false.obs;

  final initialRejectedObject = {
    'identity': {'front': false, 'back': false},
    'address': {'front': false, 'back': false}
  }.obs;

  final rejectedObject = {
    'identity': {'front': false, 'back': false},
    'address': {'front': false, 'back': false}
  }.obs;

  final uploadPercent = (-1).obs;
  final cancelToken = CancelToken().obs;
  @override
  void onInit() {
    super.onInit();
    rejectedObject.listen(handleRejectedObjectChanged);
    getUserProfile();
  }

  @override
  void onReady() {
    super.onReady();
  }

  @override
  void onClose() {}
  List<SubTypes> subTypes(type) {
    return userProfile.value.userProfileImagesMetaData.types
        .where((element) => element.name == type)
        .toList()[0]
        .subTypes;
  }

  void getUserProfile({bool silent}) async {
    try {
      if (silent != true) {
        isLoading.value = true;
      }
      final response = await identityDocumentsProvider.getUserProfile();
      if (response['status'] == true) {
        UserProfileModel tmp = UserProfileModel.fromJson(response['data']);
        // userProfile.value = tmp;
        userProfile.value = tmp;
        separateImages(tmp, 'identity');
        separateImages(tmp, 'address');
        uploadPercent.value = -1;
        resetUploads();
      }
    } catch (e) {
      log.e(e.toString());
    } finally {
      isLoading.value = false;
    }
  }

  void handleTabChange(String value) {
    activeTabIndex.value = value == 'poi' ? 0 : 1;
    if (activeTabIndex.value == 0) {
      if (addressBackImageUploadFile.value.size != 0) {
        addressBackImageUploadFile.value = PlatformFile(name: '', size: 0);
      }
      if (addressFrontImageUploadFile.value.size != 0) {
        addressFrontImageUploadFile.value = PlatformFile(name: '', size: 0);
      }
    } else {
      if (identityBackImageUploadFile.value.size != 0) {
        identityBackImageUploadFile.value = PlatformFile(name: '', size: 0);
      }
      if (identityFrontImageUploadFile.value.size != 0) {
        identityFrontImageUploadFile.value = PlatformFile(name: '', size: 0);
      }
    }
    canSubmit.value = false;
    uploadPercent.value = -1;
  }

  void separateImages(UserProfileModel tmp, String type) {
    final profileData = tmp;
    if (profileData.userProfileImages == null) {
      canChangeIdentityTypeSelect.value = true;
      canChangeResidenceTypeSelect.value = true;
    } else {
      final allImages = profileData.userProfileImages ?? [];
      final typeImages = allImages.where((i) => i.type == type).toList();
      //set identity front and back images
      UserProfileImages frontImage;
      UserProfileImages backImage;
      final tempFrontList =
          typeImages.where((element) => element.isBack != true).toList();
      if (tempFrontList.isNotEmpty) {
        frontImage = tempFrontList[0];
      }
      final tempBackList =
          typeImages.where((element) => element.isBack == true).toList();
      if (tempBackList.isNotEmpty) {
        backImage = tempBackList[0];
      }
      //set the active index
      if (frontImage != null) {
        final subtypes = tmp.userProfileImagesMetaData.types
            .where((element) => element.name == type)
            .toList()[0]
            .subTypes;
        for (var i = 0; i < subtypes.length; i++) {
          if (subtypes[i].name == frontImage.subType) {
            if (type == 'identity') {
              activeIdentitySubTypeindex.value = i;
            } else {
              activeAddressSubTypeindex.value = i;
            }
          }
        }
      }
      if (type == 'identity') {
        if (frontImage != null) {
          identityFrontImage.value = frontImage;
        }
        if (backImage != null) {
          identityBackImage.value = backImage;
        }
      } else {
        if (frontImage != null) {
          addressFrontImage.value = frontImage;
        }
        if (backImage != null) {
          addressBackImage.value = backImage;
        }
      }
      //check if we can open identity and residence type select:
      checkIfSelectTypeCanOpen(typeImages, type);
      // checkIfSelectTypeCanOpen(addressImages, 'address');
    }
  }

  void checkIfSelectTypeCanOpen(List<UserProfileImages> images, String type) {
    // ignore: invalid_use_of_protected_member
    var rej = rejectedObject.value;
    var rejChanged = false;
    var count = 0;
    for (var image in images) {
      if (image.status == 'confirmed' ||
          image.status == 'processing' ||
          image.status == 'partially_confirmed') {
        count++;
      } else if (image.status == 'rejected') {
        rejChanged = true;
        rej[type][image.isBack == true ? 'back' : 'front'] = true;
      }
    }
    if (rejChanged) {
      initialRejectedObject.value = rej;
      rejectedObject.value = rej;
      rejectedObject.refresh();
    }
    if (count > 0 && type == 'identity') {
      canChangeIdentityTypeSelect.value = false;
    }
    if (count > 0 && type == 'address') {
      canChangeResidenceTypeSelect.value = false;
    }
  }

  void handleSubtypeChange(int v, String type) {
    if (type == 'identity') {
      activeIdentitySubTypeindex.value = v;
      return;
    }
    activeAddressSubTypeindex.value = v;
  }

  void handleSelectFile({@required String side, @required String type}) async {
    FilePickerResult result = await FilePicker.platform.pickFiles(
      allowMultiple: false,
      type: FileType.custom,
      allowedExtensions: ['jpg', 'png', 'jpeg'],
    );

    final file = result.files.single;

    if (type == 'identity') {
      if (side == 'front') {
        identityFrontImageUploadFile.value = file;
      } else {
        identityBackImageUploadFile.value = file;
      }
    } else {
      if (side == 'front') {
        addressFrontImageUploadFile.value = file;
      } else {
        addressBackImageUploadFile.value = file;
      }
    }
    canSubmit.value = true;
  }

  fileLoader({@required PlatformFile file}) {
    if (!(GetPlatform.isWeb)) {
      return MultipartFile.fromFile(
        file.path,
        filename: file.name,
      );
    } else {
      return MultipartFile.fromBytes(file.bytes, filename: file.name);
    }
  }

  void submitImages() async {
    MultipartFile front;
    MultipartFile back;

    if (identityFrontImageUploadFile.value.size != 0) {
      front = await fileLoader(file: identityFrontImageUploadFile.value);
    }
    if (identityBackImageUploadFile.value.size != 0) {
      back = await fileLoader(file: identityBackImageUploadFile.value);
    }
    if (addressFrontImageUploadFile.value.size != 0) {
      front = await fileLoader(file: addressFrontImageUploadFile.value);
    }
    if (addressBackImageUploadFile.value.size != 0) {
      back = await fileLoader(file: addressBackImageUploadFile.value);
    }
    final type = activeTabIndex.value == 0 ? 'identity' : 'address';

    final additionalMap = {
      'front_image_id': _uploadFrontImageId(type: type),
      'back_image_id': _uploadBackImageId(type: type),
    };

    final uploadForm = FormData.fromMap(
      {
        if (front != null) 'front_image': front,
        //
        if (back != null) 'back_image': back,
        //
        'type': type,
        //
        'sub_type': _uploadImageSubtype(
          type: type,
        ),
        //
        if (additionalMap['front_image_id'] != null)
          'front_image_id': _uploadFrontImageId(type: type),
        //
        if (additionalMap['back_image_id'] != null)
          'back_image_id': _uploadBackImageId(type: type),
      },
    );
    uploadFiles(uploadForm);
    //final uploadTest = UploadTest();
    //uploadTest.upload(stream: uploadPercent);
  }

  void removeAllUploadedImages() {
    if (addressBackImageUploadFile.value.size != 0) {
      addressBackImageUploadFile.value = PlatformFile(name: '', size: 0);
    }
    if (addressFrontImageUploadFile.value.size != 0) {
      addressFrontImageUploadFile.value = PlatformFile(name: '', size: 0);
    }

    if (identityBackImageUploadFile.value.size != 0) {
      identityBackImageUploadFile.value = PlatformFile(name: '', size: 0);
    }
    if (identityFrontImageUploadFile.value.size != 0) {
      identityFrontImageUploadFile.value = PlatformFile(name: '', size: 0);
    }
  }

  void resetUploads() {
    uploadPercent.value = -1;
    removeAllUploadedImages();
    if (canSubmit.value == true) {
      canSubmit.value = false;
    }
  }

  String _uploadImageSubtype({String type}) {
    if (type == 'identity') {
      final identityTypes = subTypes('identity');
      final selectedIdentitySubtypeIndex = activeIdentitySubTypeindex.value;
      return identityTypes[selectedIdentitySubtypeIndex].name;
    } else {
      final addressTypes = subTypes('address');
      final selectedAddressSubtypeIndex = activeAddressSubTypeindex.value;
      return addressTypes[selectedAddressSubtypeIndex].name;
    }
  }

  _uploadFrontImageId({String type}) {
    if (type == 'identity') {
      if (identityFrontImage.value.id != null) {
        return identityFrontImage.value.id;
      }
    } else {
      if (addressFrontImage.value.id != null) {
        return addressFrontImage.value.id;
      }
    }
    return null;
  }

  _uploadBackImageId({String type}) {
    if (type == 'identity') {
      if (identityBackImage.value.id != null) {
        return identityBackImage.value.id;
      }
    } else {
      if (addressBackImage.value.id != null) {
        return addressBackImage.value.id;
      }
    }
    return null;
  }

  void uploadFiles(FormData uploadForm) async {
    try {
      cancelToken.value = CancelToken();
      final response = await identityDocumentsProvider.upload(
          form: uploadForm,
          stream: uploadPercent,
          cancelToken: cancelToken.value);
      if (response['status'] == true) {
        toastSuccess('Document is uploaded');
        getUserProfile(silent: true);
      }
    } catch (e) {
      debugPrint('error while uploading ======> ' + e.toString());
      resetUploads();
      toastError('Error while uploading document!');
    }
  }

  void handleRejectedObjectChanged(Map<String, Map<String, bool>> rejObj) {
    if (rejObj != null) {
      hasRejectedIdentity.value =
          rejObj['identity']['front'] || rejObj['identity']['back'];
      hasRejectedAddress.value =
          rejObj['address']['front'] || rejObj['address']['back'];
    }
  }

  void resetRejectedImages() {
    final type = activeTabIndex.value == 0 ? 'identity' : 'address';
    // ignore: invalid_use_of_protected_member
    var tmp = rejectedObject.value;
    if (activeTabIndex.value == 0) {
      if (tmp[type]['front']) {
        tmp[type]['front'] = false;
        identityFrontImage.value = UserProfileImages();
      }
      if (tmp[type]['back']) {
        tmp[type]['back'] = false;
        identityBackImage.value = UserProfileImages();
      }
    } else {
      if (tmp[type]['front']) {
        tmp[type]['front'] = false;
        addressFrontImage.value = UserProfileImages();
      }
      if (tmp[type]['back']) {
        tmp[type]['back'] = false;
        addressBackImage.value = UserProfileImages();
      }
    }
    rejectedObject.value = tmp;
    rejectedObject.refresh();
  }

  void cancelUpload() {
    uploadPercent.value = -1;
    cancelToken.value.cancel('upload canceled');
  }
}
