// import { take, call, put, select } from 'redux-saga/effects';

import { takeLatest, call } from 'redux-saga/effects';
import ActionTypes, { EmailVerificationPages } from './constants';
import { acountActivationAPI } from 'services/security_service';
import { StandardResponse } from 'services/constants';
import { toast } from 'components/Customized/react-toastify';
import { ToastMessages } from 'services/toastService';

import { MessageService, MessageNames } from 'services/message_service';
function* activateAcount(action: { type: string; payload: { code: string } }) {
  try {
    const response: StandardResponse = yield call(
      acountActivationAPI,
      action.payload,
    );
    if (response.status === false) {
      if (response.message && response.message.length > 0) {
        toast.warn(response.message);
      }
      ToastMessages(response.data);
      return;
    } else if (response.status === true) {
      MessageService.send({
        name: MessageNames.SET_STEP,
        payload: EmailVerificationPages.Verified,
      });
    }
  } catch (error) {
    toast.error('authentication Error');
  }
}
// Individual exports for testing
export default function* emailAuthenticationSaga() {
  yield takeLatest(ActionTypes.ACOUNT_ACTIVATION_ACTION, activateAcount);
}
