import '../../../../services/apiService.dart';

class MarketProvider {
  Future getTransactionHistory() async {
    final response = await apiService.get(
      url: "crypto-payment",
    );
    return response;
  }
}
