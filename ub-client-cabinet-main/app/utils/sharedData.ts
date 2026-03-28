import { Currency } from 'containers/AddressManagementPage/types';
import { LocalStorageKeys } from 'services/constants';
import { storage } from './storage';

export const savedPairName = (): string => {
  return (
    storage.read(LocalStorageKeys.SAVED_TRADE_PAIR, 'BTC-USDT')?.name ??
    'BTC-USDT'
  );
};
export const savedPairID = (): number => {
  return storage.read(LocalStorageKeys.SAVED_TRADE_PAIR, 1)?.id ?? 1;
};
let map;
export const currencyMap = ({ code }): Currency | undefined => {
  if (!map) {
    map = storage.read(LocalStorageKeys.CURRENCY_MAP);
  }
  if (map) {
    return map[code];
  } else {
    return undefined;
  }
};
