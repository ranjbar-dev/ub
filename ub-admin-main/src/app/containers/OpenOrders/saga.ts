import { takeLatest, put } from 'redux-saga/effects';
import { OpenOrdersActions } from './slice';
import { GetOpenOrdersAPI } from 'services/userManagementService';
import { MessageService, MessageNames } from 'services/messageService';
import { CancelOrderAPI, FullFillOrderAPI } from 'services/ordersService';
import { safeApiCall } from 'utils/sagaUtils';

export function* GetOpenOrders(action: { type: string; payload: Record<string, unknown> }) {
  const response = yield* safeApiCall(GetOpenOrdersAPI, action.payload);
  if (response) {
    yield put(OpenOrdersActions.setOpenOrdersData(response.data as Record<string, unknown>));
  }
}
export function* CancelOpenOrder(action: { type: string; payload: Record<string, unknown> }) {
  MessageService.send({
    name: MessageNames.SET_BUTTON_LOADING,
    loadingId: 'cancelButton' + action.payload.id,
    payload: true,
  });
  try {
    const response = yield* safeApiCall(CancelOrderAPI, action.payload);
    if (response) {
      yield put(OpenOrdersActions.GetOpenOrdersAction({}));
    }
  } finally {
    MessageService.send({
      name: MessageNames.SET_BUTTON_LOADING,
      loadingId: 'cancelButton' + action.payload.id,
      payload: false,
    });
  }
}
export function* FullFillOpenOrder(action: { type: string; payload: Record<string, unknown> }) {
  MessageService.send({
    name: MessageNames.SET_BUTTON_LOADING,
    loadingId: 'fullFillButton' + action.payload.id,
    payload: true,
  });
  try {
    const response = yield* safeApiCall(FullFillOrderAPI, action.payload);
    if (response) {
      yield put(OpenOrdersActions.GetOpenOrdersAction({}));
    }
  } finally {
    MessageService.send({
      name: MessageNames.SET_BUTTON_LOADING,
      loadingId: 'fullFillButton' + action.payload.id,
      payload: false,
    });
  }
}

export function* openOrdersSaga() {
  yield takeLatest(OpenOrdersActions.GetOpenOrdersAction.type, GetOpenOrders);
  yield takeLatest(
    OpenOrdersActions.CancelOpenOrderAction.type,
    CancelOpenOrder,
  );
  yield takeLatest(
    OpenOrdersActions.FullFillOpenOrderAction.type,
    FullFillOpenOrder,
  );
}