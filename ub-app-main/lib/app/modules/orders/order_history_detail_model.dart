class OrderHistoryDetailModel {
  String createdAt;
  String updatedAt;
  String pair;
  String type;
  dynamic subUnit;
  String price;
  String executed;
  String fee;
  String amount;

  OrderHistoryDetailModel(
      {createdAt, pair, type, subUnit, price, executed, fee, amount});

  OrderHistoryDetailModel.fromJson(Map<String, dynamic> json) {
    createdAt = json['createdAt'];
    updatedAt = json['updatedAt'];
    pair = json['pair'];
    type = json['type'];
    subUnit = json['subUnit'];
    price = json['price'];
    executed = json['executed'];
    fee = json['fee'];
    amount = json['amount'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['createdAt'] = createdAt;
    data['updatedAt'] = updatedAt;
    data['pair'] = pair;
    data['type'] = type;
    data['subUnit'] = subUnit;
    data['price'] = price;
    data['executed'] = executed;
    data['fee'] = fee;
    data['amount'] = amount;
    return data;
  }
}
