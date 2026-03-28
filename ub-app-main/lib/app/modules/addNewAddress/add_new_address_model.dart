class AddNewAddressModel {
  String address;
  String code;
  String label;
  String network;

  AddNewAddressModel({this.address, this.code, this.label, this.network});

  AddNewAddressModel.fromJson(Map<String, dynamic> json) {
    this.address = json['address'];
    this.code = json['code'];
    this.label = json['label'];
    this.network = json['network'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['address'] = this.address;
    data['code'] = this.code;
    data['label'] = this.label;
    if (this.network != null) data['network'] = this.network;
    return data;
  }
}
