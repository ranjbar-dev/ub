class ChangePasswordModel {
  String oldPassword;
  String newPassword;
  String confirmed;
  String s2faCode;
  String emailCode;

  ChangePasswordModel({
    this.oldPassword,
    this.newPassword,
    this.confirmed,
    this.emailCode,
    this.s2faCode,
  });

  ChangePasswordModel.fromJson(Map<String, dynamic> json) {
    oldPassword = json['old_password'];
    newPassword = json['new_password'];
    confirmed = json['confirmed'];
    s2faCode = json['2fa_code'];
    emailCode = json['emailCode'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['old_password'] = oldPassword;
    data['new_password'] = newPassword;
    data['confirmed'] = confirmed;
    if (this.s2faCode != null) {
      data['2fa_code'] = this.s2faCode;
    }
    if (this.emailCode != null) {
      data['emailCode'] = this.emailCode;
    }
    return data;
  }
}
