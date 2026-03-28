import { takeLatest, put, call } from 'redux-saga/effects';
import ActionTypes from './constants';
import {
  setIsLoadingUserProfileDataAction,
  setUserProfileAction,
} from './actions';
import { StandardResponse } from 'services/constants';
import {
  getUserProfileAPI,
  updateUserProfileAPI,
} from 'services/user_acount_service';
import { toast } from 'components/Customized/react-toastify';
import { MessageService, MessageNames } from 'services/message_service';
import { push, replace } from 'redux-first-history';
import { AppPages } from 'containers/App/constants';
import { KycStatus } from 'containers/AcountPage/constants';

function * getUserProfile (action: { type: string }) {
  yield put(setIsLoadingUserProfileDataAction(true));
  try {
    const response: StandardResponse = yield call(getUserProfileAPI);
    if (response.status === false) {
      toast.error('error while getting user profile');
      yield put(setIsLoadingUserProfileDataAction(false));
      return;
    }
    if (response.data.status === KycStatus.CONFIRMED) {
      toast.info('your acount has been verified');
      yield put(setUserProfileAction(response.data));
      yield put(replace(AppPages.AcountPage));
    } else {
      yield put(setUserProfileAction(response.data));
    }
  } catch (error) {
    toast.error('error getting user profile');
    yield put(setIsLoadingUserProfileDataAction(false));
  }
}
function * updateUserProfile (action: { type: string; payload: any }) {
  MessageService.send({ name: MessageNames.SETLOADING, payload: true });
  try {
    const response: StandardResponse = yield call(
      updateUserProfileAPI,
      action.payload,
    );
    if (response.status === false) {
      if (response.message && response.message.length > 0) {
        toast.warn(response.message);
        MessageService.send({ name: MessageNames.SETLOADING, payload: false });
        return;
      }

      toast.error('error while getting user profile');

      MessageService.send({ name: MessageNames.SETLOADING, payload: false });
      return;
    }
    yield put(push(AppPages.DocumentVerification));
    MessageService.send({ name: MessageNames.SETLOADING, payload: false });
  } catch (error) {
    toast.error('error getting user profile');
    MessageService.send({ name: MessageNames.SETLOADING, payload: false });
  }
}

// Individual exports for testing
export default function * changeUserInfoPageSaga () {
  yield takeLatest(ActionTypes.GET_USER_PROFILE, getUserProfile);
  yield takeLatest(ActionTypes.UPDATE_USER_DATA, updateUserProfile);
}
