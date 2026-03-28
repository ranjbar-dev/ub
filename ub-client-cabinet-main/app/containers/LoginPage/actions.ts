/*
 *
 * LoginPage actions
 *
 */

import { action } from 'typesafe-actions';
import {} from './types';

import ActionTypes, { LoginData } from './constants';

export const loginAction = (payload: LoginData) =>
  action(ActionTypes.LOGIN_ACTION, payload);
export const isLoggingInAction = (payload: boolean) =>
  action(ActionTypes.IS_LOGGING_IN, payload);

export const forgotPasswordAction = (payload: {
  email: string;
  recaptcha: string;
}) => action(ActionTypes.FORGOT_PASSWORD_ACTION, payload);
