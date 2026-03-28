import '../../../../../../services/apiService.dart';

class TransactionHistoryProvider {
  Future getTransactionHistory() async {
    final response = await apiService.get(
      url: "crypto-payment",
    );
    return response;
  }

  Future cancelWithdraw({int id}) async {
    final response =
        await apiService.post(url: "crypto-payment/cancel", data: {"id": id});
    return response;
  }
}
