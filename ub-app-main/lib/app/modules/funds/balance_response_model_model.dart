class BalanceResponseModel {
  String availableSum;
  List<Balance> balances;
  String btcAvailableSum;
  String btcInOrderSum;
  String btcTotalSum;
  String inOrderSum;
  String minimumOfSmallBalances;
  String totalSum;

  BalanceResponseModel(
      {availableSum,
      balances,
      btcAvailableSum,
      btcInOrderSum,
      btcTotalSum,
      inOrderSum,
      minimumOfSmallBalances,
      totalSum});

  BalanceResponseModel.fromJson(Map<String, dynamic> json) {
    availableSum = json['availableSum'];
    if (json['balances'] != null) {
      balances = <Balance>[];
      json['balances'].forEach((v) {
        balances.add(Balance.fromJson(v));
      });
    }
    btcAvailableSum = json['btcAvailableSum'];
    btcInOrderSum = json['btcInOrderSum'];
    btcTotalSum = json['btcTotalSum'];
    inOrderSum = json['inOrderSum'];
    minimumOfSmallBalances = json['minimumOfSmallBalances'];
    totalSum = json['totalSum'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['availableSum'] = availableSum;
    if (balances != null) {
      data['balances'] = balances.map((v) => v.toJson()).toList();
    }
    data['btcAvailableSum'] = btcAvailableSum;
    data['btcInOrderSum'] = btcInOrderSum;
    data['btcTotalSum'] = btcTotalSum;
    data['inOrderSum'] = inOrderSum;
    data['minimumOfSmallBalances'] = minimumOfSmallBalances;
    data['totalSum'] = totalSum;
    return data;
  }
}

class Balance {
  String address;
  String autoExchangeCode;
  String availableAmount;
  String backgroundImage;
  String btcAvailableEquivalentAmount;
  String btcInOrderEquivalentAmount;
  String btcTotalEquivalentAmount;
  String code;
  String equivalentAvailableAmount;
  String equivalentInOrderAmount;
  String equivalentTotalAmount;
  String fee;
  String image;
  String inOrderAmount;
  String minimumWithdraw;
  String name;
  String price;
  int subUnit;
  String totalAmount;

  Balance(
      {address,
      autoExchangeCode,
      availableAmount,
      backgroundImage,
      btcAvailableEquivalentAmount,
      btcInOrderEquivalentAmount,
      btcTotalEquivalentAmount,
      code,
      equivalentAvailableAmount,
      equivalentInOrderAmount,
      equivalentTotalAmount,
      fee,
      image,
      inOrderAmount,
      minimumWithdraw,
      name,
      price,
      subUnit,
      totalAmount});

  Balance.fromJson(Map<String, dynamic> json) {
    address = json['address'];
    autoExchangeCode = json['autoExchangeCode'];
    availableAmount = json['availableAmount'];
    backgroundImage = json['backgroundImage'];
    btcAvailableEquivalentAmount = json['btcAvailableEquivalentAmount'];
    btcInOrderEquivalentAmount = json['btcInOrderEquivalentAmount'];
    btcTotalEquivalentAmount = json['btcTotalEquivalentAmount'];
    code = json['code'];
    equivalentAvailableAmount = json['equivalentAvailableAmount'];
    equivalentInOrderAmount = json['equivalentInOrderAmount'];
    equivalentTotalAmount = json['equivalentTotalAmount'];
    fee = json['fee'];
    image = json['image'];
    inOrderAmount = json['inOrderAmount'];
    minimumWithdraw = json['minimumWithdraw'];
    name = json['name'];
    price = json['price'];
    subUnit = json['subUnit'];
    totalAmount = json['totalAmount'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['address'] = address;
    data['autoExchangeCode'] = autoExchangeCode;
    data['availableAmount'] = availableAmount;
    data['backgroundImage'] = backgroundImage;
    data['btcAvailableEquivalentAmount'] = btcAvailableEquivalentAmount;
    data['btcInOrderEquivalentAmount'] = btcInOrderEquivalentAmount;
    data['btcTotalEquivalentAmount'] = btcTotalEquivalentAmount;
    data['code'] = code;
    data['equivalentAvailableAmount'] = equivalentAvailableAmount;
    data['equivalentInOrderAmount'] = equivalentInOrderAmount;
    data['equivalentTotalAmount'] = equivalentTotalAmount;
    data['fee'] = fee;
    data['image'] = image;
    data['inOrderAmount'] = inOrderAmount;
    data['minimumWithdraw'] = minimumWithdraw;
    data['name'] = name;
    data['price'] = price;
    data['subUnit'] = subUnit;
    data['totalAmount'] = totalAmount;
    return data;
  }
}
