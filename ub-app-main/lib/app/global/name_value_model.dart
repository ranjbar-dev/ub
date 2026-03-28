class NameValueModel {
  String name;
  String value;

  NameValueModel({name, value});

  NameValueModel.fromJson(Map<String, dynamic> json) {
    name = json['name'];
    value = json['value'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['name'] = name;
    data['value'] = value;
    return data;
  }
}
