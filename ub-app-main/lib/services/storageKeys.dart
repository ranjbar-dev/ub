class StorageKeys {
  static final StorageKeys _singleton = StorageKeys._internal();
  factory StorageKeys() {
    return _singleton;
  }
  StorageKeys._internal();

  static final token = 'token';
  static final lastLoginDate = 'lastLoginDate';
  static final countries = 'countries';
  static final currencies = 'currencies';
  static final favPairs = 'favPairs';
  static final darkMode = 'darkMode';
  static final channel = 'channel';
  static final refresh = 'refresh';
  static final pairs = 'pairs';
  static final pairsHashMap = 'pairsHashMap';
  static final currencyPairsHashMap = 'currencyPairsHashMap';
  static final coinsHashMap = 'coinsHashMap';
  static final selectedPair = 'selectedPair';
  static final loggedInOnce = 'loggedInOnce';
  static final selectedTimeFrame = 'selectedTimeFrame';
  static final orderedPairs = 'orderedPairs';
  static final activeMarketTabIndex = 'activeMarketTabIndex';
  static final savedDepositCoins = 'savedDepositCoins';
  static final savedWithdrawalCoins = 'savedDepositCoins';
  static final lastCancelUpdate = 'lastCancelUpdate';
  static final biometricsActivated = 'biometricsActivated';
}

class SecureStorageKeys {
  static final SecureStorageKeys _singleton = SecureStorageKeys._internal();
  factory SecureStorageKeys() {
    return _singleton;
  }
  SecureStorageKeys._internal();

  static final email = 'se';
  static final password = 'sp';
}
