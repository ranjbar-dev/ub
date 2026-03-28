// import { take, call, put, select } from 'redux-saga/effects';

import { takeLatest, put, call } from 'redux-saga/effects';
import ActionTypes from './constants';
import { setIsLoadingAction, set2faQrCodeAction } from './actions';
import { toast } from 'components/Customized/react-toastify';
import { StandardResponse } from 'services/constants';
import { get2faQrcodeAPIAPI } from 'services/user_acount_service';
import { SetG2FaModel } from './types';
import { MessageService, MessageNames } from 'services/message_service';
import { set2FaAPI } from 'services/security_service';
import { ToastMessages } from 'services/toastService';
import { AppPages } from 'containers/App/constants';
import { replace } from 'redux-first-history';
import {
  set2faEnabledAction,
} from 'containers/AcountPage/actions';
import { cookieConfig, CookieKeys, cookies } from 'services/cookie';

function* get2faQrCode(action: { type: string }) {
  yield put(setIsLoadingAction(true));
  try {
    const response: StandardResponse = yield call(get2faQrcodeAPIAPI);
    if (response.status === false) {
      yield put(setIsLoadingAction(false));
      toast.error('error while getting qr code');
      return;
    }
    yield put(set2faQrCodeAction(response.data));
    yield put(setIsLoadingAction(false));
  } catch (error) {
    yield put(setIsLoadingAction(false));
    toast.error('error getting qr code');
  }
}
function* toggle2Fa(action: { type: string; payload: SetG2FaModel }) {
  MessageService.send({ name: MessageNames.SETLOADING, payload: true });
  try {
    const response: StandardResponse = yield call(set2FaAPI, action.payload);
    if (response.status === false) {
      toast.error('error while updating 2fa');
      if (response.message && response.message.length > 0) {
        toast.error(response.message);
      }
      ToastMessages(response.data);
      MessageService.send({ name: MessageNames.SETLOADING, payload: false });
      return;
    }
    MessageService.send({ name: MessageNames.SETLOADING, payload: false });

    if (response.data.token) {
      cookies.set(CookieKeys.Token, response.data.token, cookieConfig());
    }
    if (action.payload.setEnable === true) {

      MessageService.send({ name: MessageNames.SET_STEP, payload: 2 });
    } else {
      toast.success('google 2fa authentication disabled');
      yield put(set2faEnabledAction(false));
      yield put(replace(AppPages.AcountPage));
    }
  } catch (error) {
    toast.error('error updating 2fa');
    MessageService.send({ name: MessageNames.SETLOADING, payload: false });
  }
}
// Individual exports for testing
export default function* googleAuthenticationPageSaga() {
  yield takeLatest(ActionTypes.GET_2FA_QRCODE, get2faQrCode);
  yield takeLatest(ActionTypes.TOGGLE_2FA, toggle2Fa);
}
