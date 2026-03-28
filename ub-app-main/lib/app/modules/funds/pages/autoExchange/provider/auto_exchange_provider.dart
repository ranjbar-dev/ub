import 'package:meta/meta.dart' show required;

import '../../../../../../services/apiService.dart';
import '../controllers/auto_exchange_controller.dart';

class AutoExchangeProvider {
  Future submitAutoExchange({@required AutoExchangeModel model}) async {
    final response =
        await apiService.post(url: "user-balance/auto-exchange", data: {
      "code": model.code.value,
      "auto_exchange_code": model.autoExchangeCode.value,
      "mode": model.mode.value
    });
    return response;
  }
}
