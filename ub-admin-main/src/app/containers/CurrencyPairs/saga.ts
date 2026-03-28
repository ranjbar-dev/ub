// import { take, call, put, select, takeLatest } from 'redux-saga/effects';
// import { actions } from './slice';

import { takeLatest, put } from 'redux-saga/effects';
import { UpdateCurrencyPairAPI } from 'services/adminReportsService';
import {
  MessageService,
  MessageNames,
  GridNames,
} from 'services/messageService';
import { GetCurrencyPairsAPI } from 'services/userManagementService';
import { safeApiCall } from 'utils/sagaUtils';

import { CurrencyPairsActions } from './slice';

export function* GetCurrencyPairs(action: { type: string; payload: Record<string, unknown> }) {
  const response = yield* safeApiCall(GetCurrencyPairsAPI, action.payload);
  if (response) {
    yield put(CurrencyPairsActions.setCurrencyPairsData({ data: response.data, count: 1 }));
  }
}

export function* UpdateCurrencyPair(action: { type: string; payload: Record<string, unknown> }) {
  MessageService.send({
    name: MessageNames.SET_BUTTON_LOADING,
    loadingId: 'ConstructiveModalSubmit',
    payload: true,
  });
  let { payload } = action;
  if (payload.is_active === true) {
    payload.is_active = 1;
  } else if (payload.is_active === false) {
    payload.is_active = 0;
  }
  try {
    const response = yield* safeApiCall(UpdateCurrencyPairAPI, payload);
    if (response) {
      MessageService.send({
        name: MessageNames.REFRESH_GRID,
        gridName: GridNames.CURRENCY_PAIRS_PAGE,
      });
    }
  } finally {
    MessageService.send({
      name: MessageNames.SET_BUTTON_LOADING,
      loadingId: 'ConstructiveModalSubmit',
      payload: false,
    });
    MessageService.send({
      name: MessageNames.CLOSE_POPUP,
    });
  }
}

export function* currencyPairsSaga() {
  yield takeLatest(
    CurrencyPairsActions.GetCurrencyPairsAction.type,
    GetCurrencyPairs,
  );
  yield takeLatest(
    CurrencyPairsActions.UpdateCurrencyPairAction.type,
    UpdateCurrencyPair,
  );
}