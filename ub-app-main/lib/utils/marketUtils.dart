import '../app/modules/trade/models/price_model.dart';
import 'commonUtils.dart';
import 'extentions/basic.dart';

class MarketUtils {
  static priceJson(PriceModel lastPrice) {
    final json = {
      "pairName": lastPrice.name,
      "percent": lastPrice.percentage,
      "price": lastPrice.price,
      "pairId": lastPrice.id,
      "volume": lastPrice.volume,
      "equivalentPrice": lastPrice.equivalentPrice,
      "formattedVolume": lastPrice.volume.currencyFormat(compact: true),
      "formattedPrice":
          decimalPair(value: lastPrice.price, pairName: lastPrice.name),
      "formattedEquivalentPrice":
          decimalCoin(value: lastPrice.equivalentPrice, coinCode: "USDT"),
    };
    return json;
  }
}
