class UserModel {
  String email;
  String ubId;
  String phone;
  String kycLevel;
  String kycStatus;
  String kycLevelMessage;
  String securityLevel;
  String securityLevelMessage;
  String profileStatus;
  bool google2faEnabled;
  bool has2fa;
  bool isAccountVerified;
  String channelName;
  int themeId;

  UserModel(
      {email,
      ubId,
      phone,
      kycLevel,
      kycStatus,
      kycLevelMessage,
      securityLevel,
      securityLevelMessage,
      profileStatus,
      google2faEnabled,
      has2fa,
      isAccountVerified,
      channelName,
      themeId});

  UserModel.fromJson(Map<String, dynamic> json) {
    email = json['email'];
    ubId = json['ubId'];
    phone = json['phone'];
    kycLevel = json['kycLevel'];
    kycStatus = json['kycStatus'];
    kycLevelMessage = json['kycLevelMessage'];
    securityLevel = json['securityLevel'];
    securityLevelMessage = json['securityLevelMessage'];
    profileStatus = json['profileStatus'];
    google2faEnabled = json['google2faEnabled'];
    has2fa = json['has2fa'];
    isAccountVerified = json['isAccountVerified'];
    channelName = json['channelName'];
    themeId = json['themeId'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['email'] = email;
    data['ubId'] = ubId;
    data['phone'] = phone;
    data['kycLevel'] = kycLevel;
    data['kycStatus'] = kycStatus;
    data['kycLevelMessage'] = kycLevelMessage;
    data['securityLevel'] = securityLevel;
    data['securityLevelMessage'] = securityLevelMessage;
    data['profileStatus'] = profileStatus;
    data['google2faEnabled'] = google2faEnabled;
    data['has2fa'] = has2fa;
    data['isAccountVerified'] = isAccountVerified;
    data['channelName'] = channelName;
    data['themeId'] = themeId;
    return data;
  }
}
