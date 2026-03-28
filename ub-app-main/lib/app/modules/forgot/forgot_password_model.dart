class ForgotPasswordModel {
  String email;
  String recaptcha;
  String device;

  ForgotPasswordModel({this.email, this.recaptcha, this.device = 'mobile'});

  ForgotPasswordModel.fromJson(Map<String, dynamic> json) {
    this.email = json['email'];
    this.recaptcha = json['recaptcha'];
    this.device = json['device'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['email'] = this.email;
    data['recaptcha'] = this.recaptcha;
    data['device'] = this.device;
    return data;
  }
}
