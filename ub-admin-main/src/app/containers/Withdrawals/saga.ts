import { takeLatest, put } from 'redux-saga/effects';
import { WithdrawalsActions } from './slice';
import {
  GetPaymentAPI,
  GetWithdrawDetailAPI,
} from 'services/userManagementService';
import { MessageNames, MessageService } from 'services/messageService';
import { safeApiCall } from 'utils/sagaUtils';

export function* GetWithdrawals(action: { type: string; payload: Record<string, unknown> }) {
  const response = yield* safeApiCall(GetPaymentAPI, action.payload);
  if (response) {
    yield put(WithdrawalsActions.setWithdrawalsData(response.data as Record<string, unknown>));
  }
}
export function* GetWithdrawalDetails(action: { type: string; payload: Record<string, unknown> }) {
  const element = document.getElementById(
    'loading_main_withdrawals' + action.payload.id,
  );
  if (element) {
    element.style.display = 'block';
  }
  const response = yield* safeApiCall(GetWithdrawDetailAPI, {
    id: action.payload.id,
  });
  if (element) {
    element.style.display = 'none';
  }

  if (response) {
    MessageService.send({
      name: MessageNames.SET_MAIN_WITHDRAWALS_ITEM_DETAILS,
      payload: { rowData: action.payload, details: response.data },
    });
  }
}
export function* withdrawalsSaga() {
  yield takeLatest(WithdrawalsActions.GetWithdrawals.type, GetWithdrawals);
  yield takeLatest(
    WithdrawalsActions.GetWithdrawalDetailAction.type,
    GetWithdrawalDetails,
  );
}