import { takeLatest, put } from 'redux-saga/effects';
import { UserAccountsActions } from './slice';
import {
  GetUserAccountsAPI,
  GetInitialUserDataAPI,
} from 'services/userManagementService';
import { MessageService, MessageNames } from 'services/messageService';
import { WindowTypes } from 'app/constants';
import { safeApiCall } from 'utils/sagaUtils';
import { User } from './types';

export function* GetInitialUserAccounts(action: {
  type: string;
  payload: Record<string, unknown>;
}) {
  const response = yield* safeApiCall(GetUserAccountsAPI, action.payload);
  if (response) {
    yield put(UserAccountsActions.setUserAccountsData(response.data as { count: number; users: User[] }));
  }
}
export function* getInitialSingleUserDataAndOpenWindow(action: {
  type: string;
  payload: { id: number | string; windowType: WindowTypes };
}) {
  const element = document.getElementById('loading' + action.payload.id);
  if (element) {
    element.style.display = 'block';
  }
  const response = yield* safeApiCall(GetInitialUserDataAPI, {
    id: action.payload.id,
  });
  if (response) {
    MessageService.send({
      name: MessageNames.OPEN_NEW_WINDOW,
      type: action.payload.windowType,
      payload: response.data,
    });
  }
  if (element) {
    element.style.display = 'none';
  }
}

export function* userAccountsSaga() {
  yield takeLatest(
    UserAccountsActions.GetInitialUserAccountsAction.type,
    GetInitialUserAccounts,
  );
  yield takeLatest(
    UserAccountsActions.getInitialSingleUserDataAndOpenWindowAction.type,
    getInitialSingleUserDataAndOpenWindow,
  );
}
