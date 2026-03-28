import 'package:get/get.dart';
import '../../../../exchange/controllers/exchange_controller.dart';
import '../controllers/auto_exchange_controller.dart';

import '../../balance/controllers/balance_controller.dart';

class AutoExchangeBinding extends Bindings {
  @override
  void dependencies() {
    Get.put<AutoExchangeController>(AutoExchangeController(), permanent: true);
    Get.put<BalanceController>(BalanceController(), permanent: true);
    Get.put<ExchangeController>(ExchangeController(), permanent: true);
  }
}
