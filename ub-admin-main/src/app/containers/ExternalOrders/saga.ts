// import { take, call, put, select, takeLatest } from 'redux-saga/effects';
// import { actions } from './slice';

import { toast } from 'app/components/Customized/react-toastify';
import { takeLatest, put } from 'redux-saga/effects';
import {
  CancelNetQueueAPI,
  ChangeNetQueueStatusAPI,
  GetAllQueueAPI,
  GetExternalOrdersAPI,
  GetNetQueueAPI,
  SubmitNetQueueAPI,
} from 'services/externalOrdersService';
import { MessageService, MessageNames } from 'services/messageService';
import { safeApiCall } from 'utils/sagaUtils';

import { ExternalOrdersActions } from './slice';

export function* GetExternalOrder(action: { type: string; payload: Record<string, unknown> }) {
  let callObj: Record<string, unknown> = {};
  for (const key in action.payload) {
    if (action.payload[key] === 'removeFilter') {
      continue;
    }
    callObj[key] = action.payload[key];
  }

  const response = yield* safeApiCall(GetExternalOrdersAPI, callObj);
  if (response) {
    yield put(ExternalOrdersActions.setExternalOrdersData(response.data as Record<string, unknown>));
  }
}
export function* GetNetQueue(action: { type: string; payload: Record<string, unknown> }) {
  const response = yield* safeApiCall(GetNetQueueAPI, action.payload);
  if (response) {
    yield put(ExternalOrdersActions.setNetQueueData(response.data as Record<string, unknown>));
  }
}
export function* GetAllQueue(action: { type: string; payload: Record<string, unknown> }) {
  const response = yield* safeApiCall(GetAllQueueAPI, action.payload);
  if (response) {
    yield put(ExternalOrdersActions.setAllQueueData(response.data as Record<string, unknown>));
  }
}

export function* SubmitNetQueue(action: { type: string; payload: Record<string, unknown> }) {
  MessageService.send({
    name: MessageNames.SET_BUTTON_LOADING,
    loadingId: 'SubmitNetQueueRowButton' + action.payload.pair_currency_id,
    payload: true,
  });
  try {
    const response = yield* safeApiCall(SubmitNetQueueAPI, action.payload);
    if (response) {
      toast.success('Order submitted');
      yield put(ExternalOrdersActions.GetNetQueueAction({}));
    }
  } finally {
    MessageService.send({
      name: MessageNames.SET_BUTTON_LOADING,
      loadingId: 'SubmitNetQueueRowButton' + action.payload.pair_currency_id,
      payload: false,
    });
  }
}

export function* CancelNetQueue(action: { type: string; payload: Record<string, unknown> }) {
  MessageService.send({
    name: MessageNames.SET_BUTTON_LOADING,
    loadingId: 'ResetNetQueueRowButton' + action.payload.pair_currency_id,
    payload: true,
  });
  try {
    const response = yield* safeApiCall(CancelNetQueueAPI, action.payload);
    if (response) {
      toast.success('Order cancelled');
      yield put(ExternalOrdersActions.GetNetQueueAction({}));
    }
  } finally {
    MessageService.send({
      name: MessageNames.SET_BUTTON_LOADING,
      loadingId: 'ResetNetQueueRowButton' + action.payload.pair_currency_id,
      payload: false,
    });
  }
}
export function* GetNewQueueList(action: { type: string; payload: Record<string, unknown> }) {
  MessageService.send({
    name: MessageNames.SET_BUTTON_LOADING,
    loadingId: 'NetQueueRowShowListButton' + action.payload.pair_currency_id,
    payload: true,
  });
  try {
    const response = yield* safeApiCall(GetAllQueueAPI, action.payload);
    if (response) {
      yield put(ExternalOrdersActions.setNewQueueDetailList(response.data as Record<string, unknown>));
    }
  } finally {
    MessageService.send({
      name: MessageNames.SET_BUTTON_LOADING,
      loadingId: 'NetQueueRowShowListButton' + action.payload.pair_currency_id,
      payload: false,
    });
  }
}
export function* ChangeNetQueueStatus(action: { type: string; payload: Record<string, unknown> }) {
  const response = yield* safeApiCall(ChangeNetQueueStatusAPI, action.payload);
  if (response) {
    toast.success('status updated');
  }
}

export function* externalOrdersSaga() {
  yield takeLatest(
    ExternalOrdersActions.GetExternalOrderAction.type,
    GetExternalOrder,
  );
  yield takeLatest(
    ExternalOrdersActions.GetNetQueueAction.type,
    GetNetQueue,
  );
  yield takeLatest(
    ExternalOrdersActions.GetAllQueueAction.type,
    GetAllQueue,
  );
  yield takeLatest(
    ExternalOrdersActions.SubmitNetQueueAction.type,
    SubmitNetQueue,
  );
  yield takeLatest(
    ExternalOrdersActions.CancelNetQueueAction.type,
    CancelNetQueue,
  );
  yield takeLatest(
    ExternalOrdersActions.GetListNetQueueAction.type,
    GetNewQueueList,
  );
  yield takeLatest(
    ExternalOrdersActions.ChangeNetQueueStatus.type,
    ChangeNetQueueStatus,
  );
}