import 'package:get/get.dart';
import '../../market/controllers/market_controller.dart';
import '../controllers/ohlcChart_controller.dart';
import '../controllers/trade_controller.dart';

class TradeBinding extends Bindings {
  @override
  void dependencies() {
    Get.put<TradeController>(TradeController(), permanent: true);
    Get.put<OHLCChartController>(OHLCChartController(), permanent: true);
    Get.put<MarketController>(MarketController(), permanent: true);
  }
}
