import 'package:get/get.dart' show GetPlatform;

class NewTradeOrderModel {
  String amount;
  String exchangeType;
  int pairCurrencyId;
  String price;
  String stopPointPrice;
  String type;
  UserAgentInfo userAgentInfo = _buildUserAgentInfo();
  bool isFastExchange;

  static UserAgentInfo _buildUserAgentInfo() {
    String browser;
    String device;
    String os;

    if (GetPlatform.isAndroid) {
      browser = 'WebView';
      device = 'mobile';
      os = 'Android';
    } else if (GetPlatform.isIOS) {
      browser = 'WebView';
      device = 'mobile';
      os = 'iOS';
    } else if (GetPlatform.isWeb) {
      browser = 'Browser';
      device = 'web';
      os = 'Web';
    } else if (GetPlatform.isDesktop) {
      browser = 'App';
      device = 'desktop';
      os = GetPlatform.isWindows
          ? 'Windows'
          : GetPlatform.isMacOS
              ? 'macOS'
              : 'Linux';
    } else {
      browser = 'Unknown';
      device = 'unknown';
      os = 'Unknown';
    }
    return UserAgentInfo(browser: browser, device: device, os: os);
  }

  NewTradeOrderModel(
      {this.amount,
      this.exchangeType,
      this.pairCurrencyId,
      this.price,
      this.stopPointPrice,
      this.type,
      this.userAgentInfo,
      this.isFastExchange = false});

  NewTradeOrderModel.fromJson(Map<String, dynamic> json) {
    this.amount = json['amount'];
    this.exchangeType = json['exchange_type'];
    this.pairCurrencyId = json['pair_currency_id'];
    if (json['price'] != null) {
      this.price = json['price'];
    }
    if (json['stop_point_price'] != null) {
      this.stopPointPrice = json['stop_point_price'];
    }
    this.type = json['type'];
    this.userAgentInfo = _buildUserAgentInfo();
    this.isFastExchange = json['is_fast_exchange'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['amount'] = this.amount;
    data['exchange_type'] = this.exchangeType;
    data['pair_currency_id'] = this.pairCurrencyId;

    if (this.price != null) {
      data['price'] = this.price;
    }
    if (this.stopPointPrice != null) {
      data['stop_point_price'] = this.stopPointPrice;
    }

    data['type'] = this.type;
    if (this.userAgentInfo != null) {
      data['user_agent_info'] = this.userAgentInfo.toJson();
    }
    data['is_fast_exchange'] = this.isFastExchange;
    return data;
  }
}

class UserAgentInfo {
  String browser;
  String device;
  String os;

  UserAgentInfo({this.browser, this.device, this.os});

  UserAgentInfo.fromJson(Map<String, dynamic> json) {
    this.browser = json['browser'];
    this.device = json['device'];
    this.os = json['os'];
  }

  Map<String, dynamic> toJson() {
    final data = <String, dynamic>{};
    data['browser'] = this.browser;
    data['device'] = this.device;
    data['os'] = this.os;
    return data;
  }
}
