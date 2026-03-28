class TokenModel {
  String token;
  String verificationMessage;
  String refreshToken;
  bool need2fa;
  bool needEmailCode;
  bool isNewDevice;

  TokenModel(
      {this.token,
      this.refreshToken,
      this.verificationMessage,
      this.need2fa,
      this.needEmailCode,
      this.isNewDevice});

  TokenModel.fromJson(Map<String, dynamic> json) {
    token = json['token'];
    refreshToken = json['refreshToken'];
    verificationMessage = json['verificationMessage'];
    need2fa = json['need2fa'];
    needEmailCode = json['needEmailCode'];
    isNewDevice = json['isNewDevice'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['token'] = this.token;
    data['refreshToken'] = this.refreshToken;
    data['verificationMessage'] = this.verificationMessage;
    data['need2fa'] = this.need2fa;
    data['needEmailCode'] = this.needEmailCode;
    data['isNewDevice'] = this.isNewDevice;
    return data;
  }
}
