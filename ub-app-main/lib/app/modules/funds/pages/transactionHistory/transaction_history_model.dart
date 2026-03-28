class TransactionHistoryModel {
  List<Payments> payments;

  TransactionHistoryModel({payments});

  TransactionHistoryModel.fromJson(Map<String, dynamic> json) {
    if (json['payments'] != null) {
      payments = <Payments>[];
      json['payments'].forEach((v) {
        payments.add(Payments.fromJson(v));
      });
    }
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    if (payments != null) {
      data['payments'] = payments.map((v) => v.toJson()).toList();
    }
    return data;
  }
}

class Payments {
  String address;
  String addressExplorerUrl;
  String amount;
  String code;
  String createdAt;
  int id;
  String status;
  String txId;
  String txIdExplorerUrl;
  String type;

  Payments(
      {address,
      addressExplorerUrl,
      amount,
      code,
      createdAt,
      id,
      status,
      txId,
      txIdExplorerUrl,
      type});

  Payments.fromJson(Map<String, dynamic> json) {
    address = json['address'];
    addressExplorerUrl = json['addressExplorerUrl'];
    amount = json['amount'];
    code = json['code'];
    createdAt = json['createdAt'];
    id = json['id'];
    status = json['status'];
    txId = json['txId'];
    txIdExplorerUrl = json['txIdExplorerUrl'];
    type = json['type'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['address'] = address;
    data['addressExplorerUrl'] = addressExplorerUrl;
    data['amount'] = amount;
    data['code'] = code;
    data['createdAt'] = createdAt;
    data['id'] = id;
    data['status'] = status;
    data['txId'] = txId;
    data['txIdExplorerUrl'] = txIdExplorerUrl;
    data['type'] = type;
    return data;
  }
}
