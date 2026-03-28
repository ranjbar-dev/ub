// import { take, call, put, select } from 'redux-saga/effects';

import { toast } from 'components/Customized/react-toastify';
import { takeLatest, delay, put } from 'redux-saga/effects';
import { MessageNames, Subscriber } from 'services/message_service';
import { recieveNumberFromApi } from './actions';
import ActionTypes from './constants';

function* increment(action: { type: string }) {}
function* decement(action: { type: string }) {}
function* addByInputValue(action: { type: string; payload: number }) {}

function* subtractByInputValue(action: { type: string; payload: number }) {
  console.log(action.payload);
}
function* getNumberFromApi(action: { type: string; payload: number }) {
  Subscriber.next({
    name: MessageNames.SET_LOADING_TEST,
  });

  yield delay(2000);

  const apiData = Math.floor(Math.random() * (100 - 1));

  yield put(recieveNumberFromApi(apiData));

  Subscriber.next({
    name: MessageNames.SET_LOADING_END,
  });

  toast.success('Number replaced.');
}

// Individual exports for testing
export default function* contactUsPageSaga() {
  yield takeLatest(ActionTypes.INCREMENT, increment);
  yield takeLatest(ActionTypes.DECREMENT, decement);
  yield takeLatest(ActionTypes.ADD_BY_INPUT_VALUE, addByInputValue);
  yield takeLatest(ActionTypes.SUBTRACT_BY_INPUT_VALUE, subtractByInputValue);
  yield takeLatest(ActionTypes.GET_NUMBER_FROM_API, getNumberFromApi);
}
