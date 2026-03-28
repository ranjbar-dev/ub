class OhlcModel {
  String ohlcStartTime;
  String ohlcCloseTime;
  dynamic openPrice;
  dynamic closePrice;
  dynamic highPrice;
  dynamic lowPrice;
  dynamic baseVolume;
  dynamic quoteVolume;
  dynamic takerBuyBaseVolume;
  dynamic takerBuyQuoteVolume;

  OhlcModel(
      {ohlcStartTime,
      ohlcCloseTime,
      openPrice,
      closePrice,
      highPrice,
      lowPrice,
      baseVolume,
      quoteVolume,
      takerBuyBaseVolume,
      takerBuyQuoteVolume});

  OhlcModel.fromJson(Map<String, dynamic> json) {
    ohlcStartTime = json['ohlcStartTime'];
    ohlcCloseTime = json['ohlcCloseTime'];
    openPrice = json['openPrice'];
    closePrice = json['closePrice'];
    highPrice = json['highPrice'];
    lowPrice = json['lowPrice'];
    baseVolume = json['baseVolume'];
    quoteVolume = json['quoteVolume'];
    takerBuyBaseVolume = json['takerBuyBaseVolume'];
    takerBuyQuoteVolume = json['takerBuyQuoteVolume'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['ohlcStartTime'] = ohlcStartTime;
    data['ohlcCloseTime'] = ohlcCloseTime;
    data['openPrice'] = openPrice;
    data['closePrice'] = closePrice;
    data['highPrice'] = highPrice;
    data['lowPrice'] = lowPrice;
    data['baseVolume'] = baseVolume;
    data['quoteVolume'] = quoteVolume;
    data['takerBuyBaseVolume'] = takerBuyBaseVolume;
    data['takerBuyQuoteVolume'] = takerBuyQuoteVolume;
    return data;
  }
}
