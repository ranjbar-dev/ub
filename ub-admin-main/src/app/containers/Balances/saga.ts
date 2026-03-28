import { toast } from 'app/components/Customized/react-toastify';
import { take, call, put, select, takeLatest } from 'redux-saga/effects';
import { MessageService, MessageNames } from 'services/messageService';
import {
  GetBalanceHistoryAPI,
  GetBalancesAPI,
  InternalTransferAPI,
  UpdateAllBalancesAPI,
} from 'services/ordersService';
import { safeApiCall } from 'utils/sagaUtils';

import { BalancesActions } from './slice';
import { IWallet, WalletTypes } from './types';


export function* GetBalances(action: { type: string; payload: Record<string, unknown> }) {
  const response = yield* safeApiCall(GetBalancesAPI, action.payload);
  if (response) {
    yield put(BalancesActions.setBalancesData({ balances: response.data, count: 1 }));
  }
}
export function* GetBalancesForTransferModal(action: {
  type: string;
  payload: Record<string, unknown>;
}) {
  const response = yield* safeApiCall(GetBalancesAPI, action.payload);
  if (response) {
    yield put(BalancesActions.setTransferModalBalancesData({ balances: response.data as IWallet[], type: action.payload.type as string }));
  }
}
export function* UpdateAllBalances(action: { type: string; payload: Record<string, unknown> }) {
  const { loaderId, code } = action.payload;
  delete action.payload.loaderId;
  MessageService.send({
    name: MessageNames.SET_BUTTON_LOADING,
    loadingId: String(loaderId),
    payload: true,
  });
  try {
    const response = yield* safeApiCall(UpdateAllBalancesAPI, action.payload);
    if (response) {
      yield put(BalancesActions.GetBalancesAction({ type: action.payload.type }));
      toast.success(
        `${action.payload.type as string}  Wallet Updated ${code ? ',Coin : ' + code : ''}`,
      );
    }
  } finally {
    MessageService.send({
      name: MessageNames.SET_BUTTON_LOADING,
      loadingId: String(loaderId),
      payload: false,
    });
  }
}
export function* InternalTransfer(action: {
  type: string;
  payload: {
    loaderId: string;
    code: string;
    from: WalletTypes;
    fee: string;
    to: WalletTypes;
    to_custom_address?: WalletTypes;
    amount: string;
    network?: string;
  };
}) {
  const { loaderId } = action.payload;
  // @ts-expect-error — deleting non-optional property to strip it before API call
  delete action.payload.loaderId;
  // @ts-expect-error — deleting non-optional property to strip it before API call
  delete action.payload.fee;
  if (action.payload.to_custom_address) {
    // @ts-expect-error — deleting non-optional property to strip it before API call
    delete action.payload.to;
  }

  MessageService.send({
    name: MessageNames.SET_BUTTON_LOADING,
    loadingId: String(loaderId),
    payload: true,
  });
  try {
    const response = yield* safeApiCall(InternalTransferAPI, action.payload);
    if (response) {
      yield put(BalancesActions.GetBalancesAction({ type: action.payload.from }));
      toast.success('transfer in progress');
      MessageService.send({
        name: MessageNames.CLOSE_POPUP,
      });
    }
  } finally {
    MessageService.send({
      name: MessageNames.SET_BUTTON_LOADING,
      loadingId: String(loaderId),
      payload: false,
    });
  }
}

export function* GetBalanceHistory(action: { type: string; payload: Record<string, unknown> }) {
  const response = yield* safeApiCall(GetBalanceHistoryAPI, action.payload);
  if (response) {
    yield put(BalancesActions.setBalancesHistoryData(response.data as Record<string, unknown>));
  }
}

export function* balancesSaga() {
  yield takeLatest(BalancesActions.GetBalancesAction.type, GetBalances);
  yield takeLatest(
    BalancesActions.GetBalancesForTransferModalAction.type,
    GetBalancesForTransferModal,
  );
  yield takeLatest(BalancesActions.UpdateAllBalancesAction.type, UpdateAllBalances);
  yield takeLatest(BalancesActions.InternalTransferAction.type, InternalTransfer);
  yield takeLatest(BalancesActions.GetBalanceHistoryAction.type, GetBalanceHistory);
}
