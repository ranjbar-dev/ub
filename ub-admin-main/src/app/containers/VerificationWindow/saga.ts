// import { take, call, put, select, takeLatest } from 'redux-saga/effects';
// import { actions } from './slice';

import { takeLatest, put } from 'redux-saga/effects';
import {
  MessageService,
  MessageNames,
  GridNames,
} from 'services/messageService';
import { UpdateProfileImageStatusAPI } from 'services/profileImageService';
import {
  GetUserImagesAPI,
  GetUserPermissionsAPI,
  UpdateUserPermissionsAPI,
} from 'services/userManagementService';
import { safeApiCall } from 'utils/sagaUtils';

import { VerificationWindowActions } from './slice';
import { UserImagesData } from './types';
import { Permission } from '../UserDetails/types';

export function* GetUserImages(action: { type: string; payload: Record<string, unknown> }) {
  const response = yield* safeApiCall(GetUserImagesAPI, action.payload);
  if (response) {
    yield put(
      VerificationWindowActions.setUserImages({
        userId: Number(action.payload.user_id),
        data: response.data as UserImagesData,
      }),
    );
  }
}
export function* UpdateProfileImageStatus(action: {
  type: string;
  payload: {
    user_id: number;
    id: number;
    confirmation_status: string;
    type: string;
    sub_type?: string;
    newType?: string;
    loadingButtonId: string;
    rejection_reason?: string;
    id_card_code?: string;
  };
}) {
  MessageService.send({
    name: MessageNames.SET_BUTTON_LOADING,
    loadingId: action.payload.loadingButtonId,
    payload: true,
  });
  try {
    const response = yield* safeApiCall(UpdateProfileImageStatusAPI, {
      id: action.payload.id,
      confirmation_status: action.payload.confirmation_status,
      ...(action.payload.newType && { type: action.payload.newType }),
      ...(action.payload.sub_type && { sub_type: action.payload.sub_type }),
      ...(action.payload.id_card_code && {
        id_card_code: action.payload.id_card_code,
      }),
      ...(action.payload.confirmation_status === 'rejected' && {
        rejection_reason: action.payload.rejection_reason,
      }),
    });
    if (response) {
      yield put(
        VerificationWindowActions.GetUserImagesAction({
          user_id: Number(action.payload.user_id),
          type: action.payload.type,
        }),
      );
      if (action.payload.confirmation_status === 'rejected') {
        MessageService.send({
          name: MessageNames.CLOSE_POPUP,
        });
      }
      MessageService.send({
        name: MessageNames.REFRESH_GRID,
        gridName: GridNames.USER_VERIFICATION,
      });
    }
  } finally {
    MessageService.send({
      name: MessageNames.SET_BUTTON_LOADING,
      loadingId: action.payload.loadingButtonId,
      payload: false,
    });
  }
}
export function* GetPermissions(action: {
  type: string;
  payload: { id: number };
}) {
  const response = yield* safeApiCall(GetUserPermissionsAPI, {
    id: action.payload.id,
  });
  if (response) {
    yield put(
      VerificationWindowActions.setPermissionsData({
        userId: action.payload.id,
        data: response.data as Permission[],
      }),
    );
  }
}
export function* UpdatePermissions(action: {
  type: string;
  payload: { id: number; permissions: number[] };
}) {
  MessageService.send({
    name: MessageNames.SET_BUTTON_LOADING,
    loadingId: 'PermissionsButton' + action.payload.id,
    payload: true,
  });
  try {
    const response = yield* safeApiCall(
      UpdateUserPermissionsAPI,
      action.payload,
    );
    if (response) {
      MessageService.send({
        name: MessageNames.TOAST,
        payload: 'Permissions Updated ',
        type: 'success',
        userId: String(action.payload.id),
      });
    }
  } finally {
    MessageService.send({
      name: MessageNames.SET_BUTTON_LOADING,
      loadingId: 'PermissionsButton' + action.payload.id,
      payload: false,
    });
  }
}

export function* verificationWindowSaga() {
  yield takeLatest(
    VerificationWindowActions.GetUserImagesAction.type,
    GetUserImages,
  );
  yield takeLatest(
    VerificationWindowActions.UpdateProfileImageStatusAction.type,
    UpdateProfileImageStatus,
  );
  yield takeLatest(
    VerificationWindowActions.GetPermissionsAction.type,
    GetPermissions,
  );
  yield takeLatest(
    VerificationWindowActions.UpdatePermissionsAction.type,
    UpdatePermissions,
  );
}
