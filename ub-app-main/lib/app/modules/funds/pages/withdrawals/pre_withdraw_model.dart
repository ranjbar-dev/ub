class PreWithdrawModel {
  String code;
  String amount;
  String address;
  String network;
  String g2faCode;
  String emailCode;

  PreWithdrawModel(
      {this.code,
      this.amount,
      this.address,
      this.network,
      this.g2faCode,
      this.emailCode});

  PreWithdrawModel.fromJson(Map<String, dynamic> json) {
    code = json['code'];
    amount = json['amount'];
    address = json['address'];
    network = json['network'];
    g2faCode = json['g2faCode'];
    emailCode = json['emailCode'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['code'] = code;
    data['amount'] = amount;
    data['address'] = address;
    data['network'] = network;
    data['2fa_code'] = g2faCode;
    data['email_code'] = emailCode;
    return data;
  }
}
