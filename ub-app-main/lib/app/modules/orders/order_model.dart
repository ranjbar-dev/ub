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
    mainType = json['mainType'] as String;
    type = json['type'] as String;
    final rawId = json['id'];
    id = rawId == null ? null : (rawId is int ? rawId : int.tryParse(rawId.toString()));
    pair = json['pair'] as String;
    side = json['side'] as String;
    price = json['price'] as String;
    final rawSubUnit = json['subUnit'];
    subUnit = rawSubUnit == null
        ? null
        : (rawSubUnit is int ? rawSubUnit : int.tryParse(rawSubUnit.toString()));
    averagePrice = json['averagePrice'] as String;
    amount = json['amount'] as String;
    executed = json['executed'] as String;
    total = json['total'] as String;
    createdAt = json['createdAt'] as String;
    triggerCondition = json['triggerCondition'] as String;
    status = json['status'] as String;
    details = json['details'];
    createdAtToFilter = json['createdAtToFilter'] as String;
    final rawIsOpen = json['isDetailsOpen'];
    isDetailsOpen = rawIsOpen == null ? false : (rawIsOpen as bool);
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
