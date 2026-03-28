class CurrencyPairsModel {
  int id;
  String code;
  String name;
  int subUnit;
  List<Pairs> pairs;

  CurrencyPairsModel({id, code, name, subUnit, pairs});

  CurrencyPairsModel.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    code = json['code'];
    name = json['name'];
    subUnit = json['subUnit'];
    if (json['pairs'] != null) {
      pairs = <Pairs>[];
      json['pairs'].forEach((v) {
        pairs.add(Pairs.fromJson(v));
      });
    }
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['id'] = id;
    data['code'] = code;
    data['name'] = name;
    data['subUnit'] = subUnit;
    if (pairs != null) {
      data['pairs'] = pairs.map((v) => v.toJson()).toList();
    }
    return data;
  }
}

class Pairs {
  int pairId;
  String pairName;
  String dependentCode;
  String dependentName;
  int dependentId;
  int subUnit;
  int basisSubUnit;
  String basisCode;

  bool isFavorite;
  bool isMain;
  String percent;
  String image;
  double makerFee;
  double takerFee;
  String equivalentPrice;
  String price;
  String volume;
  String formattedEquivalentPrice;
  String formattedPrice;
  String formattedVolume;
  int showDigits;

  Pairs({
    this.pairId,
    this.pairName,
    this.dependentCode,
    this.dependentName,
    this.subUnit,
    this.basisSubUnit,
    this.basisCode,
    this.dependentId,
    this.isFavorite,
    this.isMain,
    this.price,
    this.equivalentPrice,
    this.percent,
    this.image,
    this.makerFee,
    this.takerFee,
    this.volume,
    this.showDigits,
    this.formattedEquivalentPrice,
    this.formattedPrice,
    this.formattedVolume,
  });

  Pairs.fromJson(Map<String, dynamic> json) {
    pairId = json['pairId'];
    pairName = json['pairName'];
    dependentCode = json['dependentCode'];
    dependentName = json['dependentName'];
    subUnit = json['subUnit'];
    basisSubUnit = json['basisSubUnit'];
    basisCode = json['basisCode'];
    dependentId = json['dependentId'];
    isFavorite = json['isFavorite'];
    isMain = json['isMain'];
    price = json['price'];
    equivalentPrice = json['equivalentPrice'];
    volume = json['volume'];
    percent = json['percent'];
    image = json['image'];
    makerFee = json['makerFee'];
    takerFee = json['takerFee'];
    volume = json['volume'];
    showDigits = json['showDigits'];
    formattedEquivalentPrice = json['formattedEquivalentPrice'];
    formattedPrice = json['formattedPrice'];
    formattedVolume = json['formattedVolume'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['pairId'] = pairId;
    data['pairName'] = pairName;
    data['dependentCode'] = dependentCode;
    data['dependentName'] = dependentName;
    data['subUnit'] = subUnit;
    data['basisSubUnit'] = basisSubUnit;
    data['basisCode'] = basisCode;
    data['dependentId'] = dependentId;
    data['isFavorite'] = isFavorite;
    data['isMain'] = isMain;
    data['price'] = price;
    data['equivalentPrice'] = equivalentPrice;
    data['percent'] = percent;
    data['image'] = image;
    data['makerFee'] = makerFee;
    data['takerFee'] = takerFee;
    data['showDigits'] = showDigits;
    data['formattedEquivalentPrice'] = formattedEquivalentPrice;
    data['formattedPrice'] = formattedPrice;
    data['formattedVolume'] = formattedVolume;
    return data;
  }
}
