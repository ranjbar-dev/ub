import ActionTypes from './constants';
import { takeLatest, put, call } from 'redux-saga/effects';
import { isChangingPasswordAction } from './actions';
import { StandardResponse } from 'services/constants';
import { changePasswordAPI } from 'services/security_service';
import { toast } from 'components/Customized/react-toastify';
import { ToastMessages } from 'services/toastService';
import { MessageService, MessageNames } from 'services/message_service';
import { ChangePasswordModel } from './types';
import { apiService } from 'services/api_service';
import { cookieConfig, CookieKeys, cookies } from 'services/cookie';

export function* changePassword(action: {
  type: string;
  payload: ChangePasswordModel;
}) {
  yield put(isChangingPasswordAction(true));
  MessageService.send({
    name: MessageNames.SETLOADING,
    payload: true,
  });
  try {
    const response: StandardResponse = yield call(
      changePasswordAPI,
      action.payload,
    );
    if (response.status == false) {
      if (response.message === 'please insert your 2fa code') {
        MessageService.send({
          name: MessageNames.OPEN_G2FA,
          payload: {
            data: action.payload,
            message: 'please insert your 2fa code',
          },
        });
      } else if (response.message && response.message.length > 0) {
        if (response.message !== 'validation failed') {
          toast.error(response.message);
        }
        MessageService.send({
          name: MessageNames.SETLOADING,
          payload: false,
        });
      }
      ToastMessages(response.data);
      yield put(isChangingPasswordAction(false));
    }

    else if (response.data && (response.data.need2fa === true || response.data.needEmailCode === true)) {
      MessageService.send({
        name: MessageNames.OPEN_TWOFA_AND_EMAILCODE_POPUP,
        payload: response.data
      });
    }
    else {
      yield put(isChangingPasswordAction(false));
      toast.success('password successfully changed!');
      cookies.set(CookieKeys.Token, response.data.token, cookieConfig());
      MessageService.send({ name: MessageNames.SET_STEP, payload: 1 });
      MessageService.send({
        name: MessageNames.CLOSE_MODAL,
      });
    }
  } catch (error) {
    put(isChangingPasswordAction(false));
  }
}
// Individual exports for testing
export default function* changePasswordSaga() {
  // See example in containers/HomePage/saga.js
  yield takeLatest(ActionTypes.CHANGE_PASSWORD_ACTION, changePassword);
}
