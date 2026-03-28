// import { take, call, put, select } from 'redux-saga/effects';

import {takeLatest,call} from 'redux-saga/effects';
import ActionTypes from './constants';
import {RegisterModel} from './types';
import {MessageService,MessageNames} from 'services/message_service';
import {toast} from 'components/Customized/react-toastify';
import {StandardResponse} from 'services/constants';
import {registerAPI} from 'services/security_service';
import {ToastMessages} from 'services/toastService';

function* register(action: {type: string; payload: RegisterModel}) {
	action.payload.email=action.payload.email.toLowerCase().replace(/ /g,'');
	MessageService.send({name: MessageNames.SETLOADING,payload: true});
	try {
		const response: StandardResponse=yield call(registerAPI,action.payload);

		if(response.status===false) {
			MessageService.send({name: MessageNames.SETLOADING,payload: false});
			if(response.message&&response.message.length>0) {
				toast.warn(response.message);
			}
			ToastMessages(response.data);
			MessageService.send({
				name: MessageNames.RESET_RECAPTCHA,
			});
			return;
		}
		MessageService.send({name: MessageNames.SETLOADING,payload: false});
		toast.success('Please check your email account');
		MessageService.send({name: MessageNames.SET_STEP,payload: 1});
	} catch(error) {
		MessageService.send({name: MessageNames.SETLOADING,payload: false});
		toast.error('register error!');
	}
	MessageService.send({
		name: MessageNames.RESET_RECAPTCHA,
	});
}

export default function* signupPageSaga() {
	yield takeLatest(ActionTypes.REGISTER_ACTION,register);
}
