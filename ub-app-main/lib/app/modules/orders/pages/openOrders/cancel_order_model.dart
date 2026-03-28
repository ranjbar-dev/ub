class CancelOrderModel {
  int orderId;

  CancelOrderModel({this.orderId});

  CancelOrderModel.fromJson(Map<String, dynamic> json) {
    orderId = json['order_id'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['order_id'] = orderId;
    return data;
  }
}
