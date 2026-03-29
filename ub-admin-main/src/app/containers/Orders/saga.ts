import { put, takeLatest } from 'redux-saga/effects';
import {
  GetOpenOrdersAPI,
  GetOrderHistoryAPI,
  GetTradeHistoryAPI,
} from 'services/userManagementService';
import { safeApiCall } from 'utils/sagaUtils';

import { OrdersActions } from './slice';

export function* GetOpenOrders(action: { type: string; payload: Record<string, unknown> }) {
  const response = yield* safeApiCall(GetOpenOrdersAPI, action.payload);
  if (response) {
    yield put(OrdersActions.setOpenOrdersData(response.data as Record<string, unknown>));
  }
}
export function* GetOrderHistory(action: { type: string; payload: Record<string, unknown> }) {
  const response = yield* safeApiCall(GetOrderHistoryAPI, action.payload);
  if (response) {
    yield put(OrdersActions.setOrderHistoryData(response.data as Record<string, unknown>));
  }
}
export function* GetTradeHistory(action: { type: string; payload: Record<string, unknown> }) {
  const response = yield* safeApiCall(GetTradeHistoryAPI, action.payload);
  if (response) {
    yield put(OrdersActions.setTradeHistoryData(response.data as Record<string, unknown>));
  }
}

export function* ordersSaga() {
  yield takeLatest(OrdersActions.GetOpenOrdersAction.type, GetOpenOrders);
  yield takeLatest(OrdersActions.GetOrderHistoryAction.type, GetOrderHistory);
  yield takeLatest(OrdersActions.GetTradeHistoryAction.type, GetTradeHistory);
}
