class CurrencyModel {
  String backgroundImage;
  String code;
  int id;
  String image;
  String mainNetwork;
  String completedNetworkName;
  String name;
  int showDigits;
  List<dynamic> otherBlockChainNetworks;

  CurrencyModel(
      {backgroundImage,
      code,
      id,
      image,
      mainNetwork,
      name,
      completedNetworkName,
      showDigits,
      otherBlockChainNetworks});

  CurrencyModel.fromJson(Map<String, dynamic> json) {
    backgroundImage = json['backgroundImage'];
    code = json['code'];
    id = json['id'];
    image = json['image'];
    completedNetworkName = json['completedNetworkName'];
    showDigits = json['showDigits'];
    mainNetwork = json['mainNetwork'];
    name = json['name'];
    otherBlockChainNetworks = json['otherBlockChainNetworks'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['backgroundImage'] = backgroundImage;
    data['code'] = code;
    data['id'] = id;
    data['showDigits'] = showDigits;
    data['image'] = image;
    data['mainNetwork'] = mainNetwork;
    data['name'] = name;
    data['completedNetworkName'] = completedNetworkName;
    data['otherBlockChainNetworks'] = otherBlockChainNetworks;
    return data;
  }
}
