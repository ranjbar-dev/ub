import '../../../../../../services/apiService.dart';

class OrderHistoryProvider {
  Future getOrderHistory({Map<String, dynamic> data}) async {
    final response =
        await apiService.get(url: "order/full-history", data: data);
    return response;
  }

  Future getOrderDetails(int orderId) async {
    final response = await apiService.get(
      url: "order/detail",
      data: {
        "order_id": orderId.toString(),
      },
    );
    return response;
  }
}
