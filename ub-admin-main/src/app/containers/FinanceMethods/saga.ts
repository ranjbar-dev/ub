// import { take, call, put, select, takeLatest } from 'redux-saga/effects';
// import { actions } from './slice';

import { put, takeLatest } from 'redux-saga/effects';
import { UpdateFinancialMethodAPI } from 'services/adminReportsService';
import {
  MessageService,
  MessageNames,
  GridNames,
} from 'services/messageService';
import { GetFinanceMethodsAPI } from 'services/userManagementService';
import { safeApiCall } from 'utils/sagaUtils';

import { FinanceMethodsActions } from './slice';

export function* GetFinanceMethods(action: { type: string; payload: Record<string, unknown> }) {
  const response = yield* safeApiCall(GetFinanceMethodsAPI, action.payload);
  if (response) {
    yield put(FinanceMethodsActions.setFinanceMethodsData({ data: response.data, count: 1 }));
  }
}

export function* UpdateFinanceMethods(action: { type: string; payload: Record<string, unknown> }) {
  MessageService.send({
    name: MessageNames.SET_BUTTON_LOADING,
    loadingId: 'ConstructiveModalSubmit',
    payload: true,
  });
  try {
    const response = yield* safeApiCall(UpdateFinancialMethodAPI, action.payload);
    if (response) {
      MessageService.send({
        name: MessageNames.REFRESH_GRID,
        gridName: GridNames.FINANCE_METHODS_PAGE,
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

export function* financeMethodsSaga() {
  yield takeLatest(
    FinanceMethodsActions.GetFinanceMethods.type,
    GetFinanceMethods,
  );
  yield takeLatest(
    FinanceMethodsActions.UpdateFinanceMethod.type,
    UpdateFinanceMethods,
  );
}