// import { take, call, put, select } from 'redux-saga/effects';

import { takeLatest, call } from 'redux-saga/effects';
import ActionTypes from './constants';
import { StandardResponse, SessionStorageKeys } from 'services/constants';
import { getRecapchaKeyAPI } from 'services/security_service';

import { toast } from 'components/Customized/react-toastify';
import { MessageService, MessageNames } from 'services/message_service';
//import {setRecapchaKeyAction} from './actions';

export function* getRecapcha(action: { type: string }) {
  try {
    const response: StandardResponse = yield call(getRecapchaKeyAPI);
    if (response.status === true) {
      //  localStorage[LocalStorageKeys.SITEKEY] = response.data.recaptchaSiteKey;
      sessionStorage[SessionStorageKeys.SITE_KEY] =
        response.data.recaptchaSiteKey;
      //  yield put(setRecapchaKeyAction(response.data.recaptchaSiteKey));
      MessageService.send({
        name: MessageNames.SET_SITE_KEY,
        payload: response.data.recaptchaSiteKey,
      });
    } else {
      toast.error('connection error');
      MessageService.send({ name: MessageNames.CLOSE_MODAL });
    }
  } catch (error) {
    toast.error('connection error!');
    MessageService.send({ name: MessageNames.CLOSE_MODAL });
  }
}

// Individual exports for testing
export default function* recapchaContainerSaga() {
  yield takeLatest(ActionTypes.GET_RECAPCHA_ACTION, getRecapcha);
}
