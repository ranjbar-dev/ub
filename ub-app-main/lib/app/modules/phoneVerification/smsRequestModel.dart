class SMSRequestModel {
  String phone;
  SMSRequestModel({this.phone});

  SMSRequestModel.fromJson(Map<String, dynamic> json) {
    phone = json['phone'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['phone'] = phone;
    return data;
  }
}
