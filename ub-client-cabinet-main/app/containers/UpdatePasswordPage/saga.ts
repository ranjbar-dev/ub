import ActionTypes, { UpdatePasswordPages } from './constants';
import { takeLatest, call } from 'redux-saga/effects';
import { StandardResponse } from 'services/constants';
import { resetPasswordAPI } from 'services/security_service';
import { toast } from 'components/Customized/react-toastify';
import { MessageService, MessageNames } from 'services/message_service';
import { ToastMessages } from 'services/toastService';
import { UpdatePasswordModel } from './types';

export function* resetPassword(action: {
  type: string;
  payload: UpdatePasswordModel;
}) {
  MessageService.send({ name: MessageNames.SETLOADING, payload: true });
  try {
    const response: StandardResponse = yield call(
      resetPasswordAPI,
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
        payload: UpdatePasswordPages.UpdatedPage,
      });
    }
    MessageService.send({ name: MessageNames.SETLOADING, payload: false });
  } catch (error) {
    MessageService.send({ name: MessageNames.SETLOADING, payload: false });
    toast.error('verification Error');
  }
}

// Individual exports for testing
export default function* updatePasswordPageSaga() {
  yield takeLatest(ActionTypes.RESET_PASSWORD_ACTION, resetPassword);
}
