// import { take, call, put, select, takeLatest } from 'redux-saga/effects';
// import { actions } from './slice';

import { takeLatest, put } from 'redux-saga/effects';
import { GetOpenOrdersAPI } from 'services/userManagementService';
import { safeApiCall } from 'utils/sagaUtils';

import { FilledOrdersActions } from './slice';

export function* GetFilledOrders(action: { type: string; payload: Record<string, unknown> }) {
  const response = yield* safeApiCall(GetOpenOrdersAPI, action.payload);
  if (response) {
    yield put(FilledOrdersActions.setFilledOrdersData(response.data as Record<string, unknown>));
  }
}

export function* filledOrdersSaga() {
  yield takeLatest(
    FilledOrdersActions.GetFilledOrdersAction.type,
    GetFilledOrders,
  );
}
