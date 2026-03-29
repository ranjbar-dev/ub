import { takeLatest, put } from 'redux-saga/effects';
import { DepositsActions } from './slice';
import { GetPaymentAPI } from 'services/userManagementService';
import {
  MessageService,
  MessageNames,
  GridNames,
} from 'services/messageService';
import { UpdateDepositAPI } from 'services/ordersService';
import { DepositSaveData } from '../Billing/types';
import { safeApiCall } from 'utils/sagaUtils';

export function* GetDepositOrders(action: { type: string; payload: Record<string, unknown> }) {
  const response = yield* safeApiCall(GetPaymentAPI, action.payload);
  if (response) {
    yield put(DepositsActions.setDepositsData(response.data as Record<string, unknown>));
  }
}
export function* UpdateDepositOrder(action: {
  type: string;
  payload: DepositSaveData;
}) {
  MessageService.send({
    name: MessageNames.SET_BUTTON_LOADING,
    loadingId:
      action.payload.should_deposit === true
        ? 'DepositModalSaveAndDepositButton'
        : 'DepositModalSaveButton',
    payload: true,
  });
  try {
    const response = yield* safeApiCall(UpdateDepositAPI, action.payload);
    if (response) {
      MessageService.send({
        name: MessageNames.UPDATE_GRID_ROW,
        gridName: GridNames.DEPOSITS_PAGE,
        rowId: Number(action.payload.id),
        payload: {
          amount: action.payload.amount,
          fromAddress: action.payload.from_address,
          toAddress: action.payload.to_address,
          status: action.payload.status,
          txId: action.payload.tx_id,
        },
      });
      MessageService.send({
        name: MessageNames.CLOSE_POPUP,
      });
    }
  } finally {
    MessageService.send({
      name: MessageNames.SET_BUTTON_LOADING,
      loadingId:
        action.payload.should_deposit === true
          ? 'DepositModalSaveAndDepositButton'
          : 'DepositModalSaveButton',
      payload: false,
    });
  }
}

export function* depositsSaga() {
  yield takeLatest(DepositsActions.GetDepositsAction.type, GetDepositOrders);
  yield takeLatest(
    DepositsActions.UpdateDepositsAction.type,
    UpdateDepositOrder,
  );
}
