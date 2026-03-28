import '../../../../../../services/apiService.dart';
import '../cancel_order_model.dart';

class OpenOrdersProvider {
  Future getOpenOrders() async {
    final response = await apiService.get(
      url: "order/open-orders",
    );
    return response;
  }

  Future cancelOrders({CancelOrderModel model}) async {
    final response = await apiService.post(
      url: "order/cancel",
      data: model,
    );
    return response;
  }
}
