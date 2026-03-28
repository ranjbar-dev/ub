class PriceModel {
  String price;
  String percentage;
  int id;
  String name;
  String equivalentPrice;
  String volume;
  String high;
  String low;

  PriceModel({price, percentage, id, name, equivalentPrice, volume, high, low});

  PriceModel.fromJson(Map<String, dynamic> json) {
    price = json['price'];
    percentage = json['percentage'];
    id = json['id'];
    name = json['name'];
    equivalentPrice = json['equivalentPrice'];
    volume = json['volume'];
    high = json['high'];
    low = json['low'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['price'] = price;
    data['percentage'] = percentage;
    data['id'] = id;
    data['name'] = name;
    data['equivalent_price'] = equivalentPrice;
    data['volume'] = volume;
    data['high'] = high;
    data['low'] = low;
    return data;
  }
}
