class UserProfileModel {
  String address;
  dynamic adminComment;
  int country;
  String countryName;
  String dateOfBirth;
  String firstName;
  String gender;
  int id;
  String lastName;
  String postalCode;
  String regionAndCity;
  String status;
  String updatedAt;
  List<UserProfileImages> userProfileImages;
  UserProfileImagesMetaData userProfileImagesMetaData;

  UserProfileModel(
      {address,
      adminComment,
      country,
      countryName,
      dateOfBirth,
      firstName,
      gender,
      id,
      lastName,
      postalCode,
      regionAndCity,
      status,
      updatedAt,
      userProfileImages,
      userProfileImagesMetaData});

  UserProfileModel.fromJson(Map<String, dynamic> json) {
    address = json['address'];
    adminComment = json['adminComment'];
    country = json['country'];
    countryName = json['countryName'];
    dateOfBirth = json['dateOfBirth'];
    firstName = json['firstName'];
    gender = json['gender'];
    id = json['id'];
    lastName = json['lastName'];
    postalCode = json['postalCode'];
    regionAndCity = json['regionAndCity'];
    status = json['status'];
    updatedAt = json['updatedAt'];
    if (json['userProfileImages'] != null) {
      userProfileImages = <UserProfileImages>[];
      json['userProfileImages'].forEach((v) {
        userProfileImages.add(UserProfileImages.fromJson(v));
      });
    }
    userProfileImagesMetaData = json['userProfileImagesMetaData'] != null
        ? UserProfileImagesMetaData.fromJson(json['userProfileImagesMetaData'])
        : null;
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['address'] = address;
    data['adminComment'] = adminComment;
    data['country'] = country;
    data['countryName'] = countryName;
    data['dateOfBirth'] = dateOfBirth;
    data['firstName'] = firstName;
    data['gender'] = gender;
    data['id'] = id;
    data['lastName'] = lastName;
    data['postalCode'] = postalCode;
    data['regionAndCity'] = regionAndCity;
    data['status'] = status;
    data['updatedAt'] = updatedAt;
    if (userProfileImages != null) {
      data['userProfileImages'] =
          userProfileImages.map((v) => v.toJson()).toList();
    }
    if (userProfileImagesMetaData != null) {
      data['userProfileImagesMetaData'] = userProfileImagesMetaData.toJson();
    }
    return data;
  }
}

class UserProfileImages {
  String createdAt;
  int id;
  String idCardCode;
  String image;
  int imageId;
  bool isBack;
  int mainImageId;
  String rejectionReason;
  String status;
  String subType;
  String type;

  UserProfileImages(
      {createdAt,
      id,
      idCardCode,
      image,
      imageId,
      isBack,
      mainImageId,
      rejectionReason,
      status,
      subType,
      type});

  UserProfileImages.fromJson(Map<String, dynamic> json) {
    createdAt = json['createdAt'];
    id = json['id'];
    idCardCode = json['idCardCode'];
    image = json['image'];
    imageId = json['imageId'];
    isBack = json['isBack'];
    mainImageId = json['mainImageId'];
    rejectionReason = json['rejectionReason'];
    status = json['status'];
    subType = json['subType'];
    type = json['type'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['createdAt'] = createdAt;
    data['id'] = id;
    data['idCardCode'] = idCardCode;
    data['image'] = image;
    data['imageId'] = imageId;
    data['isBack'] = isBack;
    data['mainImageId'] = mainImageId;
    data['rejectionReason'] = rejectionReason;
    data['status'] = status;
    data['subType'] = subType;
    data['type'] = type;
    return data;
  }
}

class UserProfileImagesMetaData {
  List<Types> types;

  UserProfileImagesMetaData({types});

  UserProfileImagesMetaData.fromJson(Map<String, dynamic> json) {
    if (json['types'] != null) {
      types = <Types>[];
      json['types'].forEach((v) {
        types.add(Types.fromJson(v));
      });
    }
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    if (types != null) {
      data['types'] = types.map((v) => v.toJson()).toList();
    }
    return data;
  }
}

class Types {
  String name;
  List<SubTypes> subTypes;

  Types({name, subTypes});

  Types.fromJson(Map<String, dynamic> json) {
    name = json['name'];
    if (json['subTypes'] != null) {
      subTypes = <SubTypes>[];
      json['subTypes'].forEach((v) {
        subTypes.add(SubTypes.fromJson(v));
      });
    }
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['name'] = name;
    if (subTypes != null) {
      data['subTypes'] = subTypes.map((v) => v.toJson()).toList();
    }
    return data;
  }
}

class SubTypes {
  String name;
  bool hasBack;

  SubTypes({name, hasBack});

  SubTypes.fromJson(Map<String, dynamic> json) {
    name = json['name'];
    hasBack = json['hasBack'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['name'] = name;
    data['hasBack'] = hasBack;
    return data;
  }
}
