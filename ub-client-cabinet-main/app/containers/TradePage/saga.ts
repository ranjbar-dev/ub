// import { take, call, put, select } from 'redux-saga/effects';

import { toast } from 'components/Customized/react-toastify';
import { Currency } from 'containers/AddressManagementPage/types';
import { LocalStorageKeys, StandardResponse } from 'services/constants';
import { getMarketTradesAPI } from 'services/marketTrade_services';
import {
  DataInjectMessageNames,
  DataInjectMessageService,
} from 'services/message_service';
import { all, call, put, takeEvery, takeLatest } from 'redux-saga/effects';
import { getCurrenciesAPI } from 'services/address_management_service';
import { addRemoveFavPairAPI, getPairsListAPI } from 'services/pairs_service';
import { storage } from 'utils/storage';
import { setPairMapAction } from './actions';
import ActionTypes from './constants';
import { PairItem } from './types';

function* addRemoveFavoritePair(action: {
  type: string;
  payload: { pair_currency_id: number; action: 'add' | 'remove' };
}) {
  try {
    yield call(addRemoveFavPairAPI, action.payload);
  } catch (error) {
    console.log(error);
  }
}
function* getCurrencies(action: { type: string }) {
  try {


    const [currenciesResponse, pairsResponse] = yield all([getCurrenciesAPI(), getPairsListAPI()]);

    if (currenciesResponse.status === true) {
      storage.write(LocalStorageKeys.CURRENCIES, currenciesResponse.data.currencies);
      const tmp = {};
      const currencies = storage.read(LocalStorageKeys.CURRENCIES);
      currencies.forEach((item: Currency) => {
        tmp[item.code] = item;
      });
      storage.write(LocalStorageKeys.CURRENCY_MAP, tmp);
    }
    //set pairsMap
    let pairs: PairItem[] = [];
    pairsResponse.data.forEach((item: any) => {
      pairs = [...pairs, ...item.pairs];
    });
    const pairsMap = {};
    pairs.forEach((item: PairItem) => {
      pairsMap[item.pairName] = item;
    });
    storage.write(LocalStorageKeys.PAIRS_MAP, pairsMap);
    yield put(setPairMapAction(pairsMap));

  } catch (error) {
    console.log(error);
  }
}

function* getInitialMarketTradesData(action: {
  type: string;
  payload: { pairName: string };
}) {
  try {
    const response: StandardResponse = yield call(getMarketTradesAPI, {
      pair: action.payload.pairName,
    });
    DataInjectMessageService.send({
      name: DataInjectMessageNames.MARKET_TRADES_INITIAL_DATA,
      data: response.data,
    });
  } catch (error) {
    toast.error('error fetching market trades');
  }
}

// Individual exports for testing
export default function* tradePageSaga() {
  yield takeEvery(ActionTypes.ADD_REMOVE_FAVORITE_PAIR, addRemoveFavoritePair);
  yield takeEvery(ActionTypes.GET_CURRENCIES, getCurrencies);
  yield takeLatest(
    ActionTypes.GET_INITIAL_MARKET_TRADES_DATA,
    getInitialMarketTradesData,
  );
}
