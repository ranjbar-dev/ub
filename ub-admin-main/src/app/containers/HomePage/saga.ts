import { WindowTypes } from 'app/constants';
import { takeLatest } from 'redux-saga/effects';
import { MessageService, MessageNames } from 'services/messageService';
import {
  GetInitialUserDataAPI,
  GetWithdrawDetailAPI,
} from 'services/userManagementService';
import { safeApiCall } from 'utils/sagaUtils';

import { HomePageActions } from './slice';

export function* getUserInfoById(action: {
  type: string;
  payload: { id?: string; email?: string };
}) {
  const response = yield* safeApiCall(GetInitialUserDataAPI, {
    ...(action.payload.id && { id: action.payload.id }),
    ...(action.payload.email && { email: action.payload.email }),
  });
  if (response) {
    MessageService.send({
      name: MessageNames.OPEN_NEW_WINDOW,
      type: WindowTypes.User,
      payload: response.data,
    });
  }
}
export function* getWithdrawalById(action: {
  type: string;
  payload: { id: string };
}) {
  const response = yield* safeApiCall(GetWithdrawDetailAPI, {
    id: action.payload.id,
  });
  if (response) {
    MessageService.send({
      name: MessageNames.SET_MAIN_WITHDRAWALS_ITEM_DETAILS,
      payload: {
        rowData: { ...action.payload, userId: (response.data as Record<string, unknown>).userId },
        details: response.data,
      },
    });
  }
}

export function* homePageSaga() {
  yield takeLatest(HomePageActions.getUserByIdAction.type, getUserInfoById);
  yield takeLatest(
    HomePageActions.getWithdrawalByIdAction.type,
    getWithdrawalById,
  );
}
