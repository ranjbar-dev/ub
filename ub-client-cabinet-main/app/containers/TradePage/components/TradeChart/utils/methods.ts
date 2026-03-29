import { LanguageCode } from '../charting_library/charting_library.min';
import { CentrifugoChannels } from 'containers/App/constants';

export const getLanguageFromURL = (): LanguageCode | null => {
  const regex = new RegExp('[\\?&]lang=([^&#]*)');
  const results = regex.exec(location.search);
  return results === null
    ? null
    : (decodeURIComponent(results[1].replace(/\+/g, ' ')) as LanguageCode);
};
const intervalVariable: string = '1';
const localPair = localStorage.getItem('pair');
const symbol =
  localPair === null
    ? 'BTC/USDT'
    : JSON.parse(localPair).pairName.replace('-', '/');
export const prepareTopic = () => {
  let itervalStr;

  switch (intervalVariable) {
    case '1':
    case '3':
    case '5':
    case '15':
    case '30':
    case '45':
      itervalStr = '1minute';
      break;
    case '60':
    case '120':
    case '180':
    case '240':
      itervalStr = '1hour';
      break;
    case '1D':
    case '1W':
    case '1M':
      itervalStr = '1day';
      break;
    default:
      itervalStr = '1minute';
  }
  return `${CentrifugoChannels.TradeChartPrefix}${itervalStr}:${symbol.replace(
    '/',
    '-',
  )}`;
};
export const prepareSymbolName = (symbol: string) => {
  return symbol.split('/');
};
