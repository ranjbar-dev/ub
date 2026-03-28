class SignupModel {
  String email;
  String password;
  String recaptcha;

  SignupModel({this.email, this.password, this.recaptcha = 'palangMalang'});

  SignupModel.fromJson(Map<String, dynamic> json) {
    this.email = json['email'];
    this.password = json['password'];
    this.recaptcha = json['recaptcha'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['email'] = this.email;
    data['password'] = this.password;
    data['recaptcha'] = this.recaptcha;
    return data;
  }
}
