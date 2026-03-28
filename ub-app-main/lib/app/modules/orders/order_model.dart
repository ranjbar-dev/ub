class OrderModel {
  String mainType;
  String type;
  int id;
  String pair;
  String side;
  String price;
  int subUnit;
  String averagePrice;
  String amount;
  String executed;
  String total;
  String createdAt;
  String triggerCondition;
  String status;
  dynamic details;
  String createdAtToFilter;
  bool isDetailsOpen;

  OrderModel(
      {mainType,
      type,
      id,
      pair,
      side,
      price,
      subUnit,
      averagePrice,
      amount,
      executed,
      total,
      createdAt,
      triggerCondition,
      status,
      details,
      createdAtToFilter,
      isDetailsOpen});

  OrderModel.fromJson(Map<String, dynamic> json) {
    mainType = json['mainType'];
    type = json['type'];
    id = json['id'];
    pair = json['pair'];
    side = json['side'];
    price = json['price'];
    subUnit = json['subUnit'];
    averagePrice = json['averagePrice'];
    amount = json['amount'];
    executed = json['executed'];
    total = json['total'];
    createdAt = json['createdAt'];
    triggerCondition = json['triggerCondition'];
    status = json['status'];
    details = json['details'];
    createdAtToFilter = json['createdAtToFilter'];
    isDetailsOpen = json['isDetailsOpen'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['mainType'] = mainType;
    data['type'] = type;
    data['id'] = id;
    data['pair'] = pair;
    data['side'] = side;
    data['price'] = price;
    data['subUnit'] = subUnit;
    data['averagePrice'] = averagePrice;
    data['amount'] = amount;
    data['executed'] = executed;
    data['total'] = total;
    data['createdAt'] = createdAt;
    data['triggerCondition'] = triggerCondition;
    data['status'] = status;
    data['details'] = details;
    data['createdAtToFilter'] = createdAtToFilter;
    data['isDetailsOpen'] = isDetailsOpen;
    return data;
  }
}
