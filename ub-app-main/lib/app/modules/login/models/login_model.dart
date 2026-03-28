import 'package:meta/meta.dart' show required;

class LoginModel {
  String username;
  String password;
  String s2faCode;
  String recaptcha;
  String emailCode;

  LoginModel(
      {@required this.username,
      @required this.password,
      this.s2faCode,
      this.recaptcha,
      this.emailCode});

  LoginModel.fromJson(Map<String, dynamic> json) {
    username = json['username'];
    password = json['password'];
    s2faCode = json['2fa_code'];
    emailCode = json['emailCode'];
    recaptcha = json['recaptcha'];
  }

  Map<String, dynamic> toJson() {
    final Map<String, dynamic> data = <String, dynamic>{};
    data['username'] = this.username;
    data['password'] = this.password;
    if (this.recaptcha != null) {
      data['recaptcha'] = this.recaptcha;
    }
    if (this.s2faCode != null) {
      data['2fa_code'] = this.s2faCode;
    }
    if (this.emailCode != null) {
      data['emailCode'] = this.emailCode;
    }
    return data;
  }
}
