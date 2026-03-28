import { put, takeLatest } from 'redux-saga/effects';
import { GetCommitionsAPI } from 'services/adminReportsService';
import {
  MessageService,
  MessageNames,
  GridNames,
} from 'services/messageService';
import { UpdateDepositAPI } from 'services/ordersService';
import {
  GetBillingGridDataAPI,
  GetWithdrawDetailAPI,
  UpdateWithdrawAPI,
  AddPaymentCommentAPI,
} from 'services/userManagementService';
import { safeApiCall } from 'utils/sagaUtils';

import { BillingActions } from './slice';
import { DepositSaveData, Payment, PaymentDetails } from './types';


export function* GetBillingGridData(action: {
  type: string;
  payload: { user_id: number; type?: string };
}) {
  const response = yield* safeApiCall(GetBillingGridDataAPI, action.payload);
  if (response) {
    yield put(BillingActions.setBillingData(response.data as Record<string, unknown>));
  }
}
export function* GetBillingDepositsData(action: {
  type: string;
  payload: { user_id: number };
}) {
  const response = yield* safeApiCall(GetBillingGridDataAPI, {
    ...action.payload,
    type: 'deposit',
  });
  if (response) {
    yield put(BillingActions.setBillingDepositsData(response.data as Record<string, unknown>));
  }
}
export function* GetBillingWithdrawalsData(action: {
  type: string;
  payload: { user_id: number };
}) {
  const response = yield* safeApiCall(GetBillingGridDataAPI, {
    ...action.payload,
    type: 'withdraw',
  });
  if (response) {
    yield put(BillingActions.setBillingWithdrawalsData(response.data as Record<string, unknown>));
  }
}
export function* GetBillingAllTransactionsData(action: {
  type: string;
  payload: { user_id: number };
}) {
  const response = yield* safeApiCall(GetBillingGridDataAPI, action.payload);
  if (response) {
    yield put(BillingActions.setBillingAllTransactionsData(response.data as Record<string, unknown>));
  }
}
export const statusSelector = (payload: Record<string, unknown>) => {
  if (payload.admin_status) {
    if (
      payload.admin_status === 'recheck' ||
      payload.admin_status === 'pending'
    ) {
      return 'created';
    } else if (payload.admin_status === 'approved') {
      return 'completed';
    }
  } else if (payload.status) {
    return payload.status;
  }
  return 'created';
};
export function* UpdateBillingWithdraw(action: { type: string; payload: Record<string, unknown> }) {
  const { user_id: rawUserId, buttonId, ...data } = action.payload;
  const user_id = String(rawUserId);
  if (buttonId) {
    MessageService.send({
      name: MessageNames.SET_BUTTON_LOADING,
      loadingId: String(buttonId),
      payload: true,
    });
  }
  try {
    const response = yield* safeApiCall(UpdateWithdrawAPI, data);
    if (response) {
      MessageService.send({
        name: MessageNames.SHOW_WITHDRAW_CONFIRM_POPUP,
        payload: 'withdraw updated ',
        type: 'success',
        userId: user_id + '' + action.payload.id,
      });
      yield put(BillingActions.updateBillingWithdrawRow(response.data as Record<string, unknown>));
      MessageService.send({
        name: MessageNames.UPDATE_GRID_ROW,
        gridName: GridNames.MAIN_WITHDRAWALS,
        rowId: action.payload.id as number,
        payload: {
          status: statusSelector(action.payload),
        },
        userId: user_id as string | number,
      });
      MessageService.send({
        name: MessageNames.UPDATE_GRID_ROW,
        gridName: GridNames.Billing_Withdraw,
        rowId: action.payload.id as number,
        payload: {
          status: statusSelector(action.payload),
        },
        userId: user_id as string | number,
      });
      MessageService.send({
        name: MessageNames.REFRESH_GRID,
        gridName: GridNames.MAIN_WITHDRAWALS,
      });
    }
  } finally {
    if (buttonId) {
      MessageService.send({
        name: MessageNames.SET_BUTTON_LOADING,
        loadingId: buttonId as string,
        payload: false,
      });
      if ((buttonId as string).includes('reject')) {
        MessageService.send({
          name: MessageNames.CLOSE_REJECT_POPUP,
        });
      }
    }
  }
}
export function* GetBillingWithdrawDetails(action: {
  type: string;
  payload: { id: number; user_id: number };
}) {
  MessageService.send({
    name: MessageNames.SET_ROW_LOADING,
    payload: {
      rowId: Number(action.payload.id),
      userId: action.payload.user_id as number,
      gridName: GridNames.Billing_Withdraw,
      value: true,
    },
  });
  try {
    const response = yield* safeApiCall(GetWithdrawDetailAPI, {
      id: action.payload.id,
    });
    if (response) {
      yield put(BillingActions.setWithdrawalItemDetails({ rowData: action.payload as unknown as Payment, details: response.data as PaymentDetails }));
    }
  } finally {
    MessageService.send({
      name: MessageNames.SET_ROW_LOADING,
      payload: {
        rowId: Number(action.payload.id),
        userId: action.payload.user_id as number,
        gridName: GridNames.Billing_Withdraw,
        value: false,
      },
    });
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
        gridName: GridNames.BILLING_DEPOSIT,
        rowId: action.payload.id as number,
        payload: {
          amount: action.payload.amount,
          fromAddress: action.payload.from_address,
          toAddress: action.payload.to_address,
          status: action.payload.status,
          txId: action.payload.tx_id,
        },
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
    MessageService.send({
      name: MessageNames.CLOSE_POPUP,
    });
  }
}
export function* AddPaymentComment(action: { type: string; payload: Record<string, unknown> }) {
  MessageService.send({
    name: MessageNames.SET_BUTTON_LOADING,
    loadingId: 'addWithdrawComment' + action.payload.payment_id,
    payload: true,
  });
  try {
    const response = yield* safeApiCall(AddPaymentCommentAPI, {
      payment_id: action.payload.payment_id,
      comment: action.payload.comment,
    });
    if (!response) {
      MessageService.send({
        name: MessageNames.TOAST,
        type: 'error',
        payload: 'Error While Adding Comment',
        userId: 'withdrawPopup' + action.payload.payment_id,
      });
    }
  } finally {
    MessageService.send({
      name: MessageNames.SET_BUTTON_LOADING,
      loadingId: 'addWithdrawComment' + action.payload.payment_id,
      payload: false,
    });
  }
}
export function* GetCommitions(action: {
  type: string;
  payload: { userId: number };
}) {
  const response = yield* safeApiCall(GetCommitionsAPI, {
    user_id: action.payload.userId,
  });
  if (response) {
    yield put(BillingActions.setCommissionsData(response.data as Record<string, unknown>));
  }
}

export function* billingSaga() {
  yield takeLatest(
    BillingActions.GetBillingGridDataAction.type,
    GetBillingGridData,
  );
  yield takeLatest(
    BillingActions.GetBillingDepositsDataAction.type,
    GetBillingDepositsData,
  );
  yield takeLatest(
    BillingActions.GetBillingWithdrawalsDataAction.type,
    GetBillingWithdrawalsData,
  );
  yield takeLatest(
    BillingActions.GetBillingAllTransactionsDataAction.type,
    GetBillingAllTransactionsData,
  );
  yield takeLatest(
    BillingActions.GetBillingWithdrawDetailsAction.type,
    GetBillingWithdrawDetails,
  );
  yield takeLatest(
    BillingActions.UpdateBillingWithdrawAction.type,
    UpdateBillingWithdraw,
  );
  yield takeLatest(
    BillingActions.UpdateDepositsAction.type,
    UpdateDepositOrder,
  );
  yield takeLatest(
    BillingActions.AddPaymentCommentAction.type,
    AddPaymentComment,
  );
  yield takeLatest(BillingActions.GetCommitionsAction.type, GetCommitions);
}
