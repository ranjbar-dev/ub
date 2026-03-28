class PairBalanceModel {
  String sum;
  PairData pairData;
  List<PairBalances> pairBalances;
  Fee fee;
  List<Charts> charts;

  PairBalanceModel({sum, pairData, pairBalances, fee, charts});

  PairBalanceModel.fromJson(Map<String, dynamic> json) {
    sum = json['sum'];
    pairData =
        json['pairData'] != null ? PairData.fromJson(json['pairData']) : null;
    if (json['pairBalances'] != null) {
      pairBalances = <PairBalances>[];
      json['pairBalances'].forEach((v) {
        pairBalances.add(PairBalances.fromJson(v));
      });
    }
    fee = json['fee'] != null ? Fee.fromJson(json['fee']) : null;
    if (json['charts'] != null) {
      charts = <Charts>[];
      json['charts'].forEach((v) {
        charts.add(Charts.fromJson(v));
      });
    }
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['sum'] = sum;
    if (pairData != null) {
      data['pairData'] = pairData.toJson();
    }
    if (pairBalances != null) {
      data['pairBalances'] = pairBalances.map((v) => v.toJson()).toList();
    }
    if (fee != null) {
      data['fee'] = fee.toJson();
    }
    if (charts != null) {
      data['charts'] = charts.map((v) => v.toJson()).toList();
    }
    return data;
  }
}

class PairData {
  int id;
  String minimumOrderAmount;
  String name;

  PairData({id, minimumOrderAmount, name});

  PairData.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    minimumOrderAmount = json['minimumOrderAmount'];
    name = json['name'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['id'] = id;
    data['minimumOrderAmount'] = minimumOrderAmount;
    data['name'] = name;
    return data;
  }
}

class PairBalances {
  String balance;
  String currencyCode;
  int currencyId;
  String currencyName;

  PairBalances({balance, currencyCode, currencyId, currencyName});

  PairBalances.fromJson(Map<String, dynamic> json) {
    balance = json['balance'];
    currencyCode = json['currencyCode'];
    currencyId = json['currencyId'];
    currencyName = json['currencyName'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['balance'] = balance;
    data['currencyCode'] = currencyCode;
    data['currencyId'] = currencyId;
    data['currencyName'] = currencyName;
    return data;
  }
}

class Fee {
  double makerFee;
  double takerFee;

  Fee({makerFee, takerFee});

  Fee.fromJson(Map<String, dynamic> json) {
    makerFee = json['makerFee'];
    takerFee = json['takerFee'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['makerFee'] = makerFee;
    data['takerFee'] = takerFee;
    return data;
  }
}

class Charts {
  String amount;
  String equivalentAmount;
  String name;
  String percent;

  Charts({amount, equivalentAmount, name, percent});

  Charts.fromJson(Map<String, dynamic> json) {
    amount = json['amount'];
    equivalentAmount = json['equivalentAmount'];
    name = json['name'];
    percent = json['percent'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['amount'] = amount;
    data['equivalentAmount'] = equivalentAmount;
    data['name'] = name;
    data['percent'] = percent;
    return data;
  }
}
