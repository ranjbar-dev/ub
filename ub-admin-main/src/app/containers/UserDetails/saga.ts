import { takeLatest } from 'redux-saga/effects';
import { MessageService, MessageNames } from 'services/messageService';
import {
  GetUserBalancesAPI,
  GetUserWhiteAddressesAPI,
  UpdateUserDataAPI,
} from 'services/userManagementService';
import { safeApiCall } from 'utils/sagaUtils';

import { UserDetailsActions } from './slice';

export function* GetWallets(action: { type: string; payload: { id: number } }) {
  const response = yield* safeApiCall(GetUserBalancesAPI, {
    id: action.payload.id,
  });
  if (response) {
    MessageService.send({
      name: MessageNames.SET_WALLETS_DATA,
      payload: response.data,
    });
  }
}
export function* GetWhiteAddresses(action: {
  type: string;
  payload: { id: number };
}) {
  const response = yield* safeApiCall(GetUserWhiteAddressesAPI, {
    id: action.payload.id,
  });
  if (response) {
    MessageService.send({
      name: MessageNames.SET_WHITEADDRESSES_DATA,
      payload: response.data,
    });
  }
}

export function* UpdateUserData(action: { type: string; payload: Record<string, unknown> }) {
  MessageService.send({
    name: MessageNames.SET_BUTTON_LOADING,
    loadingId: 'userEdit',
    payload: true,
  });
  try {
    const response = yield* safeApiCall(UpdateUserDataAPI, action.payload);
    if (response) {
      MessageService.send({
        name: MessageNames.TOAST,
        payload: 'user data updated',
        type: 'success',
        userId: String(action.payload.id),
      });
      MessageService.send({
        name: MessageNames.SET_USER_DATA,
        payload: response.data,
      });
    }
  } finally {
    MessageService.send({
      name: MessageNames.SET_BUTTON_LOADING,
      loadingId: 'userEdit',
      payload: false,
    });
  }
}

export function* userDetailsSaga() {
  yield takeLatest(UserDetailsActions.GetWalletsAction.type, GetWallets);
  yield takeLatest(
    UserDetailsActions.GetWhiteAddressesAction.type,
    GetWhiteAddresses,
  );

  yield takeLatest(
    UserDetailsActions.UpdateUserDataAction.type,
    UpdateUserData,
  );
}
