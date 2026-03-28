import {takeLatest,put,call} from 'redux-saga/effects';
import ActionTypes from './constants';
// import { login } from 'containers/LoginPage/saga';
import {isLoadingAction,setUserDataAction} from './actions';
import {
	getUserDataAPI,
	getNewVerificationEmailAPI,
} from 'services/security_service';
import {StandardResponse} from 'services/constants';
import {toast} from 'components/Customized/react-toastify';

// import { take, call, put, select } from 'redux-saga/effects';
export function* getUserData(action) {
	// yield put(isLoadingAction(true));
	try {
		const response: StandardResponse=yield call(getUserDataAPI);
		if(response.status===true) {
			yield put(setUserDataAction(response.data));
		} else {
			yield put(isLoadingAction(false));
		}
	} catch(error) {
		yield put(isLoadingAction(false));
	}
}
export function* getNewVerificationEmail(action) {
	try {
		const response: StandardResponse=yield call(getNewVerificationEmailAPI);
		if(response.status===true) {
		} else {
			toast.error('Please Try Again Later');
		}
	} catch(error) { }
}

// Individual exports for testing
export default function* acountPageSaga() {
	// See example in containers/HomePage/saga.js
	// yield takeLatest(ActionTypes.LOGIN_ACTION, login);
	yield takeLatest(ActionTypes.DEFAULT_ACTION,getUserData);
	yield takeLatest(
		ActionTypes.GET_NEW_VERIFICATION_EMAIL_ACTION,
		getNewVerificationEmail,
	);
}
