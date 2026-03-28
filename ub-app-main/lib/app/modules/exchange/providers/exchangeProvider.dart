import '../../trade/new_trade_order_model.dart';

import '../../../../../../services/apiService.dart';

class ExchangeProvider {
  Future getCurrencyPairDetails({int pairId}) async {
    final response = await apiService.get(
      url: "user-balance/pair-balance",
      data: {"pair_currency_id": pairId},
    );
    return response;
  }

  Future createExchange({NewTradeOrderModel model}) async {
    model.userAgentInfo = UserAgentInfo.fromJson(
        {"browser": "Chrome", "device": "web", "os": "Win32"});
    final response = await apiService.post(
      url: "order/create",
      data: model,
    );
    return response;
  }

  Future getPairsPrice() async {
    final response = await apiService.get(
      url:
          'currencies/pairs-statistic?pair_currencies=BTC-USDT|ETH-USDT|BCH-USDT|DASH-USDT|DOGE-USDT|MKR-USDT|LTC-USDT|ETH-BTC|TRX-USDT',
    );
    return response;
  }
}
