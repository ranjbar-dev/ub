import ActionTypes from './constants';
import { takeLatest, call, put } from 'redux-saga/effects';
import { StandardResponse } from 'services/constants';
import { getChartConfigAPI } from 'services/trade_chart_service';
import { setChartConfigAction } from './actions';
import { toast } from 'components/Customized/react-toastify';

function* getChartConfig(action: { type: string }) {
  try {
    const response: StandardResponse = yield call(getChartConfigAPI);
    if (response.status === true) {
      yield put(setChartConfigAction(response.data));
    } else {
      toast.error('error getting chart config');
    }
  } catch (error) {
    toast.error('error while getting chart config');
  }
}

// Individual exports for testing
export default function* tradeChartSaga() {
  yield takeLatest(ActionTypes.GET_CHART_CONFIG, getChartConfig);
}
