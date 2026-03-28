class Toggle2faModel {
  String code;
  String password;

  Toggle2faModel({this.code, this.password});

  Toggle2faModel.fromJson(Map<String, dynamic> json) {
    code = json['code'];
    password = json['password'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['code'] = this.code;
    data['password'] = this.password;
    return data;
  }
}
