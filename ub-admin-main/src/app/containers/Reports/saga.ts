// import { take, call, put, select, takeLatest } from 'redux-saga/effects';
// import { actions } from './slice';

import { takeLatest, put } from 'redux-saga/effects';
import { GetWithdrawalCommentsAPI } from 'services/adminReportsService';
import {
  AddAdminCommentAPI,
  DeleteAdminCommentAPI,
  EditAdminCommentAPI,
} from 'services/adminReportsService';
import { MessageService, MessageNames } from 'services/messageService';
import { GetAdminReportsAPI } from 'services/userManagementService';
import { safeApiCall } from 'utils/sagaUtils';

import { Report } from './types';
import { ReportsActions } from './slice';

export function* GetAdminReports(action: {
  type: string;
  payload: { id: number };
}) {
  const response = yield* safeApiCall(GetAdminReportsAPI, action.payload);
  if (response) {
    yield put(ReportsActions.setAdminReportsData(response.data as { comments: Report[] }));
  }
}
export function* AddAdminComment(action: {
  type: string;
  payload: {
    user_id: number;
    comment: string;
  };
}) {
  MessageService.send({
    name: MessageNames.SET_BUTTON_LOADING,
    loadingId: 'userReportsSubmitButton' + action.payload.user_id,
    payload: true,
  });
  try {
    const response = yield* safeApiCall(AddAdminCommentAPI, action.payload);
    if (response) {
      yield put(
        ReportsActions.GetAdminReportsAction({
          id: Number(action.payload.user_id),
        }),
      );
    }
  } finally {
    MessageService.send({
      name: MessageNames.SET_BUTTON_LOADING,
      loadingId: 'userReportsSubmitButton' + action.payload.user_id,
      payload: false,
    });
  }
}

export function* DeleteAdminComment(action: { type: string; payload: Record<string, unknown> }) {
  MessageService.send({
    name: MessageNames.SET_BUTTON_LOADING,
    loadingId: 'deleteComentButton' + action.payload.user_id,
    payload: true,
  });
  try {
    const response = yield* safeApiCall(DeleteAdminCommentAPI, {
      id: action.payload.id,
    });
    if (response) {
      yield put(
        ReportsActions.GetAdminReportsAction({
          id: Number(action.payload.user_id),
        }),
      );
    }
  } finally {
    MessageService.send({
      name: MessageNames.SET_BUTTON_LOADING,
      loadingId: 'deleteComentButton' + action.payload.user_id,
      payload: false,
    });
    MessageService.send({
      name: MessageNames.CLOSE_POPUP,
    });
  }
}

export function* EditAdminComment(action: { type: string; payload: Record<string, unknown> }) {
  MessageService.send({
    name: MessageNames.SET_BUTTON_LOADING,
    loadingId: 'editComentButton' + action.payload.user_id,
    payload: true,
  });
  try {
    const response = yield* safeApiCall(EditAdminCommentAPI, {
      id: action.payload.id,
      comment: action.payload.comment,
    });
    if (response) {
      yield put(
        ReportsActions.GetAdminReportsAction({
          id: Number(action.payload.user_id),
        }),
      );
    }
  } finally {
    MessageService.send({
      name: MessageNames.SET_BUTTON_LOADING,
      loadingId: 'editComentButton' + action.payload.user_id,
      payload: false,
    });
    MessageService.send({
      name: MessageNames.CLOSE_POPUP,
    });
  }
}
export function* GetWithdrawalComments(action: { type: string; payload: Record<string, unknown> }) {
  const response = yield* safeApiCall(GetWithdrawalCommentsAPI, action.payload);
  if (response) {
    yield put(ReportsActions.setWithdrawalComments(response.data as Record<string, unknown>));
  }
}
export function* reportsSaga() {
  yield takeLatest(ReportsActions.GetAdminReportsAction.type, GetAdminReports);
  yield takeLatest(ReportsActions.AddAdmiCommentAction.type, AddAdminComment);
  yield takeLatest(
    ReportsActions.DeleteAdminCommentAction.type,
    DeleteAdminComment,
  );
  yield takeLatest(
    ReportsActions.EditAdminCommentAction.type,
    EditAdminComment,
  );
  yield takeLatest(
    ReportsActions.GetWithdrawalCommentsAction.type,
    GetWithdrawalComments,
  );
}
