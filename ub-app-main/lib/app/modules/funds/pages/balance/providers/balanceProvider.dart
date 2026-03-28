import '../../../../../../services/apiService.dart';

class BalanceProvider {
  Future getBalances() async {
    final response = await apiService.get(
      url: "user-balance/balance?sort=desc",
    );
    return response;
  }
}
