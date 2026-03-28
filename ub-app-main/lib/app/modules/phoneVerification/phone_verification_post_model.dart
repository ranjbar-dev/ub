class PhoneVerificationPostModel {
  String phone;
  String code;
  String password;
  String s2faCode;
  String emailCode;

  PhoneVerificationPostModel({
    this.phone,
    this.code,
    this.password,
    this.s2faCode,
    this.emailCode,
  });

  PhoneVerificationPostModel.fromJson(Map<String, dynamic> json) {
    phone = json['phone'];
    code = json['code'];
    s2faCode = json['2fa_code'];
    password = json['password'];
    s2faCode = json['2fa_code'];
    emailCode = json['emailCode'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['phone'] = phone;
    data['code'] = code;
    data['password'] = password;
    if (this.s2faCode != null) {
      data['2fa_code'] = this.s2faCode;
    }
    if (this.emailCode != null) {
      data['emailCode'] = this.emailCode;
    }
    return data;
  }
}
