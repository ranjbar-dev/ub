/*
 *
 * AcountPage actions
 *
 */

import {action} from 'typesafe-actions';
import {UserData} from './types';

import ActionTypes from './constants';
// import { LoginData } from 'containers/LoginPage/constants';

export const getUserDataAction=() => action(ActionTypes.DEFAULT_ACTION);
export const isLoadingAction=(payload: boolean) =>
	action(ActionTypes.IS_LOADING_ACTION,payload);
export const setUserDataAction=(payload: UserData) =>
	action(ActionTypes.SET_USER_DATA_ACTION,payload);
export const set2faEnabledAction=(payload: boolean) =>
	action(ActionTypes.SET_2FA_ENABLED_ACTION,payload);
export const loginAction=(payload: boolean) =>
	action(ActionTypes.LOGGED_IN_ACTION,payload);
export const getNewVerificationEmailAction=() =>
	action(ActionTypes.GET_NEW_VERIFICATION_EMAIL_ACTION);
