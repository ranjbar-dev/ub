// import { take, call, put, select, takeLatest } from 'redux-saga/effects';
// import { actions } from './slice';

import { toast } from 'app/components/Customized/react-toastify';
import { call, put, takeLatest } from 'redux-saga/effects';
import { MessageService, MessageNames } from 'services/messageService';
import {
  GetMarketTicksAPI,
  GetCurrencyPairsAPI,
  SyncTicksAPI,
  GetSyncListAPI,
} from 'services/userManagementService';
import { safeApiCall } from 'utils/sagaUtils';

import { MarketTicksActions } from './slice';

export function* GetMarketTicks(action: { type: string; payload: Record<string, unknown> }) {
  const response = yield* safeApiCall(GetMarketTicksAPI, action.payload);
  if (response) {
    yield put(MarketTicksActions.setMarketTicksData({ data: response.data, count: 200 }));
  }
}
export function* GetCurrencyPairs(action: { type: string; payload: Record<string, unknown> }) {
  const response = yield* safeApiCall(GetCurrencyPairsAPI, {});
  if (response) {
    MessageService.send({
      name: MessageNames.SET_CURRENCY_PAIRS,
      payload: response.data,
    });
  }
}
export function* GetSyncList(action: { type: string; payload: Record<string, unknown> }) {
  const response = yield* safeApiCall(GetSyncListAPI, action.payload);
  if (response) {
    yield put(MarketTicksActions.setSyncListData({ data: response.data, count: 200 }));
  }
}
export function* SyncTicks(action: { type: string; payload: Record<string, unknown> }) {
  MessageService.send({
    name: MessageNames.SET_BUTTON_LOADING,
    loadingId: 'syncButton',
    payload: true,
  });
  try {
    const response = yield* safeApiCall(SyncTicksAPI, action.payload);
    if (response) {
      toast.success('started syncing ticks on server');
    }
  } finally {
    MessageService.send({
      name: MessageNames.SET_BUTTON_LOADING,
      loadingId: 'syncButton',
      payload: false,
    });
  }
}

export function* marketTicksSaga() {
  yield takeLatest(
    MarketTicksActions.GetMarketTicksAction.type,
    GetMarketTicks,
  );
  yield takeLatest(
    MarketTicksActions.GetCurrencyPairsAction.type,
    GetCurrencyPairs,
  );
  yield takeLatest(
    MarketTicksActions.GetSyncListAction.type,
    GetSyncList,
  );
  yield takeLatest(MarketTicksActions.SyncTicksAction.type, SyncTicks);
}