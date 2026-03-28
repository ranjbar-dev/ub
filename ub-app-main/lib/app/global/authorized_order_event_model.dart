class AuthorizedOrderEventModel {
  dynamic id;
  String amount;
  String price;
  String status;
  String type;
  String pairCurrency;

  AuthorizedOrderEventModel({id, amount, price, status, type, pairCurrency});

  AuthorizedOrderEventModel.fromJson(Map<String, dynamic> json) {
    id = json['id'];
    amount = json['amount'];
    price = json['price'];
    status = json['status'];
    type = json['type'];
    pairCurrency = json['pairCurrency'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['id'] = id;
    data['amount'] = amount;
    data['price'] = price;
    data['status'] = status;
    data['type'] = type;
    data['pairCurrency'] = pairCurrency;
    return data;
  }
}
