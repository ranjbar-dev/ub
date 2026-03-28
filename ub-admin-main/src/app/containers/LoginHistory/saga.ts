import { put, takeLatest } from 'redux-saga/effects';
import { GetLoginHistoryAPI } from 'services/userManagementService';
import { safeApiCall } from 'utils/sagaUtils';

import { LoginHistoryActions } from './slice';

export function* GetLoginHistory(action: { type: string; payload: Record<string, unknown> }) {
  const response = yield* safeApiCall(GetLoginHistoryAPI, action.payload);
  if (response) {
    yield put(LoginHistoryActions.setLoginHistoryData(response.data as Record<string, unknown>));
  }
}

export function* loginHistorySaga() {
  yield takeLatest(LoginHistoryActions.GetLoginHistory.type, GetLoginHistory);
}
