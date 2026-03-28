import ActionTypes from './constants';
import { takeLatest, put, call } from 'redux-saga/effects';
import { setBalancePageDataAction } from './actions';

import { toast } from 'components/Customized/react-toastify';
import { StandardResponse } from 'services/constants';
import { getBalancesAPI } from 'services/funds_services';

export const infinitePageSize = window.innerHeight > 1100 ? 60 : 30;
function* getBalances(action: { type: string }) {
  try {
    const response: StandardResponse = yield call(getBalancesAPI);
    if (response.status === false) {
      toast.error('error getting balances');
      return;
    }
    yield put(setBalancePageDataAction(response.data));
  } catch (error) {
    toast.error('error while getting balances');
  }
}

// Individual exports for testing
export default function* fundsPageSaga() {
  yield takeLatest(ActionTypes.GET_BALANCE_DATA_ACTION, getBalances);
}
