class CountryModel {
  int id;
  String name;
  String fullName;
  String code;
  String image;

  CountryModel({id, name, fullName, code, image});

  CountryModel.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    name = json['name'];
    fullName = json['fullName'];
    code = json['code'];
    image = json['image'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['id'] = id;
    data['name'] = name;
    data['fullName'] = fullName;
    data['code'] = code;
    data['image'] = image;
    return data;
  }
}
