import '../../../../services/apiService.dart';
import '../new_trade_order_model.dart';

class TradeProvider {
  Future getCurrencyPairDetails({int pairId}) async {
    final response = await apiService.get(
      url: "user-balance/pair-balance",
      data: {"pair_currency_id": pairId},
    );
    return response;
  }

  Future createOrder({NewTradeOrderModel model}) async {
    model.userAgentInfo = UserAgentInfo.fromJson(
        {"browser": "Chrome", "device": "web", "os": "Win32"});
    final response = await apiService.post(
      url: "order/create",
      data: model,
    );
    return response;
  }
}
