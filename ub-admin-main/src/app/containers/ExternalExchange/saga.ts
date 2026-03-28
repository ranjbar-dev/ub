// import { take, call, put, select, takeLatest } from 'redux-saga/effects';
// import { actions } from './slice';

import { takeLatest, put } from 'redux-saga/effects';
import { GetExternalExchangeAPI } from 'services/userManagementService';
import { safeApiCall } from 'utils/sagaUtils';

import { ExternalExchangeActions } from './slice';

export function* GetExternalExchange(action: { type: string; payload: Record<string, unknown> }) {
  const response = yield* safeApiCall(GetExternalExchangeAPI, action.payload);
  if (response) {
    yield put(ExternalExchangeActions.setExternalExchangeData(response.data as Record<string, unknown>));
  }
}

export function* externalExchangeSaga() {
  yield takeLatest(
    ExternalExchangeActions.GetExternalExchange.type,
    GetExternalExchange,
  );
}
