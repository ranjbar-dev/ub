import 'package:meta/meta.dart' show required;

import '../../../../../../services/apiService.dart';

class DepositsProvider {
  Future getUserDepositData({@required String code}) async {
    final response = await apiService.get(
      url: "user-balance/withdraw-deposit",
      data: {"code": code},
    );
    return response;
  }
}
