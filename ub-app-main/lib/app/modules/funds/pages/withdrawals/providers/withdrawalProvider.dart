import 'package:meta/meta.dart' show required;
import '../pre_withdraw_model.dart';
import '../../../../../../services/apiService.dart';

class WithdrawalProvider {
  Future preWithdraw({@required PreWithdrawModel data}) async {
    final response = await apiService.post(
      url: 'crypto-payment/pre-withdraw',
      data: data,
    );
    return response;
  }

  Future withdraw({@required PreWithdrawModel data}) async {
    final response = await apiService.post(
      url: 'crypto-payment/withdraw',
      data: data,
    );
    return response;
  }
}
