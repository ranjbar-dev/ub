class HomePagePairPriceModel {
  int pairId;
  String pairName;
  String dependentCode;
  String dependentName;
  int dependentId;
  String subUnit;
  int basisSubUnit;
  String basisCode;
  bool isFavorite;
  bool isMain;

  String equivalentPrice;
  String image;
  String secondImage;

  String lastUpdate;
  double makerFee;

  String percent;
  String price;
  int showDigits;

  double takerFee;
  List<TrendData> trendData;

/*

18:"trendData" -> List (100 items)

 */

  HomePagePairPriceModel(
      {this.basisCode,
      this.basisSubUnit,
      this.dependentCode,
      this.dependentId,
      this.dependentName,
      this.equivalentPrice,
      this.image,
      this.secondImage,
      this.isFavorite,
      this.isMain,
      this.lastUpdate,
      this.makerFee,
      this.pairId,
      this.pairName,
      this.percent,
      this.price,
      this.showDigits,
      this.subUnit,
      this.takerFee,
      this.trendData});

  HomePagePairPriceModel.fromJson(Map<String, dynamic> json) {
    basisCode = json['basisCode'];
    basisSubUnit = json['basisSubUnit'];
    dependentCode = json['dependentCode'];
    dependentId = json['dependentId'];
    dependentName = json['dependentName'];
    equivalentPrice = json['equivalentPrice'];
    image = json['image'];
    secondImage = json['secondImage'];
    isFavorite = json['isFavorite'];
    isMain = json['isMain'];
    lastUpdate = json['lastUpdate'];
    makerFee = json['makerFee'];
    pairId = json['pairId'];
    pairName = json['pairName'];
    percent = json['percent'];
    price = json['price'];
    showDigits = json['showDigits'];
    subUnit = json['subUnit'].toString();
    takerFee = json['takerFee'];
    if (json['trendData'] != null) {
      trendData = <TrendData>[];
      json['trendData'].forEach((v) {
        trendData.add(TrendData.fromJson(v));
      });
    }
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['basisCode'] = basisCode;
    data['basisSubUnit'] = basisSubUnit;
    data['dependentCode'] = dependentCode;
    data['dependentId'] = dependentId;
    data['dependentName'] = dependentName;
    data['equivalentPrice'] = equivalentPrice;
    data['image'] = image;
    data['secondImage'] = secondImage;
    data['isFavorite'] = isFavorite;
    data['isMain'] = isMain;
    data['lastUpdate'] = lastUpdate;
    data['makerFee'] = makerFee;
    data['pairId'] = pairId;
    data['pairName'] = pairName;
    data['percent'] = percent;
    data['price'] = price;
    data['showDigits'] = showDigits;
    data['subUnit'] = subUnit;
    data['takerFee'] = takerFee;
    if (trendData != null) {
      data['trendData'] = trendData.map((v) => v.toJson()).toList();
    }
    return data;
  }
}

class TrendData {
  String price;
  String time;
  String change;

  TrendData({this.price, this.time, this.change});

  TrendData.fromJson(Map<String, dynamic> json) {
    price = json['price'];
    time = json['time'];
    change = json['change'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['price'] = price;
    data['time'] = time;
    data['change'] = change;
    return data;
  }
}
