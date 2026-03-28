import '../../../../services/apiService.dart';
import '../../../../services/constants.dart';

class HomePageProvider {
  Future getPairsPrice() async {
    final response = await apiService.get(
      url:
          'currencies/pairs-statistic?pair_currencies=BTC-USDT|ETH-USDT|BCH-USDT|DASH-USDT',
    );
    return response;
  }

  Future getLastNews() async {
    final response = await apiService.rawGet(
      rawUrl: '${Constants.cmsAddress}/ubnews?_sort=date:desc&_limit=5',
    );
    return response;
  }
}
