class WithdrawAddressModel {
  int id;
  String address;
  String label;
  bool isFavorite;
  String code;
  String name;
  String network;
  String icon;

  WithdrawAddressModel(
      {id, address, label, isFavorite, code, name, network, icon});

  WithdrawAddressModel.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    address = json['address'];
    label = json['label'];
    isFavorite = json['isFavorite'];
    code = json['code'];
    name = json['name'];
    network = json['network'];
    icon = json['icon'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['id'] = id;
    data['address'] = address;
    data['label'] = label;
    data['isFavorite'] = isFavorite;
    data['code'] = code;
    data['name'] = name;
    data['network'] = network;
    data['icon'] = icon;
    return data;
  }
}
