import { toast } from 'components/Customized/react-toastify';
import { call, put, takeLatest } from 'redux-saga/effects';
import { StandardResponse } from 'services/constants';
import { ToastMessages } from 'services/toastService';
import {
  getCountriesAPI,
  requestSMSAPI,
  verifyCodeAPI,
} from 'services/user_acount_service';

import {
  setCountriesAction,
  setCountriesLoading,
  setisLoadingAction,
  setPhoneNumberAction,
  setStepAction,
} from './actions';
import ActionTypes, { Country, PhoneVerificationSteps } from './constants';
import { MessageService, MessageNames } from 'services/message_service';
import { getUserDataAction } from 'containers/AcountPage/actions';

function * getCountries (action: { type: string }) {
  yield put(setCountriesLoading(true));
  try {
    const response: StandardResponse = yield call(getCountriesAPI);
    yield put(setCountriesAction(response.data));
    yield put(setCountriesLoading(false));
  } catch (error) {
    toast.error('error getting countries');

    yield put(setCountriesAction([]));
    yield put(setCountriesLoading(false));
  }
}

function * getSMS (action: {
  type: string;
  payload: { country: Country; phoneNumber: string };
}) {
  const phone = '+' + action.payload.country.code + action.payload.phoneNumber;
  MessageService.send({ name: MessageNames.SETLOADING, payload: true });
  // yield put(setisLoadingAction(true));
  try {
    const response: StandardResponse = yield call(requestSMSAPI, { phone });
    if (response.status === false) {
      if (response.data) {
        toast.error(response.message);
        ToastMessages(response.data);
      } else if (response.message) {
        toast.error(response.message);
      }
    } else {
      yield put(setPhoneNumberAction(phone));
      yield put(setStepAction(PhoneVerificationSteps.ENTER_CODE));
    }
    MessageService.send({ name: MessageNames.SETLOADING, payload: false });
    // yield put(setisLoadingAction(false));
  } catch (error) {
    toast.error('failed');
    MessageService.send({ name: MessageNames.SETLOADING, payload: false });

    // yield put(setisLoadingAction(false));
  }
}
function * reSendSMS (action: { type: string; payload: { phone: string } }) {
  yield put(setisLoadingAction(true));
  try {
    const response: StandardResponse = yield call(requestSMSAPI, {
      phone: action.payload.phone,
    });
    if (response.status === false) {
      if (response.data) {
        toast.error(response.message);
        ToastMessages(response.data);
      }
    } else {
      toast.info('verification code has been sent');
      // yield put(setStepAction(PhoneVerificationSteps.ENTER_CODE));
    }
    yield put(setisLoadingAction(false));
  } catch (error) {
    toast.error('failed');
    yield put(setisLoadingAction(false));
  }
}
function * verifyCode (action: {
  type: string;
  payload: {
    phone: string;
    code: string;
    '2fa_code'?: string;
    password?: string;
  };
}) {
  MessageService.send({ name: MessageNames.SETLOADING, payload: true });
  try {
    const response: StandardResponse = yield call(verifyCodeAPI, action.payload);
    if (response.status === false) {
      if (response.message.includes('password')) {
        yield put(setStepAction(PhoneVerificationSteps.ENTER_ACOUNT_PASSWORD));
        MessageService.send({ name: MessageNames.SETLOADING, payload: false });
        if (action.payload.password) {
          toast.error(response.message);
        }
        return;
      }
      if (response.data) {
        toast.error(response.message);
        ToastMessages(response.data);
      } else if (response.message) {
        toast.error(response.message);
      }
    } else if (
      response.data &&
      (response.data.need2fa === true || response.data.needEmailCode === true)
    ) {
      yield put(setStepAction(PhoneVerificationSteps.GOOGLE_2FA_STEP));
      return;
    } else if (response.status === true) {
      yield put(setStepAction(PhoneVerificationSteps.DONE_STEP));
      yield put(getUserDataAction());
      // yield put(setStepAction(PhoneVerificationSteps.ENTER_CODE));
    }
    MessageService.send({ name: MessageNames.SETLOADING, payload: false });
  } catch (error) {
    toast.error('failed');
    MessageService.send({ name: MessageNames.SETLOADING, payload: false });
  }
}

// Individual exports for testing
export default function * phoneVerificationPageSaga () {
  yield takeLatest(ActionTypes.GET_COUNTRIES_ACTION, getCountries);
  yield takeLatest(ActionTypes.GET_SMS_ACTION, getSMS);
  yield takeLatest(ActionTypes.RESEND_SMS_ACTION, reSendSMS);
  yield takeLatest(ActionTypes.VERIFY_CODE, verifyCode);
}
