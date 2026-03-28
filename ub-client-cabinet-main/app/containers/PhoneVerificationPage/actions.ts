/*
 *
 * PhoneVerificationPage actions
 *
 */

import { action } from 'typesafe-actions';

import ActionTypes, { PhoneVerificationSteps } from './constants';

export const defaultAction = () => action(ActionTypes.DEFAULT_ACTION);

export const getCountriesAction = () =>
  action(ActionTypes.GET_COUNTRIES_ACTION);

export const setCountriesAction = (payload: any) =>
  action(ActionTypes.SET_COUNTRIES_ACTION, payload);

export const setCountriesLoading = (payload: boolean) =>
  action(ActionTypes.SET_COUNTRIES_LOADING, payload);

export const setisLoadingAction = (payload: boolean) =>
  action(ActionTypes.SET_IS_SENDING_SMS, payload);

export const getSMSAction = (payload: any) =>
  action(ActionTypes.GET_SMS_ACTION, payload);

export const resendSMSAction = (payload: { phone: string }) =>
  action(ActionTypes.RESEND_SMS_ACTION, payload);

export const setStepAction = (payload: PhoneVerificationSteps) =>
  action(ActionTypes.SET_STEP_ACTION, payload);

export const setPhoneNumberAction = (payload: string) =>
  action(ActionTypes.SET_PHONE_NUMBER_ACTION, payload);

export const verifyCodeAction = (payload: {
  code: string;
  phone: string;
  '2fa_code'?: string;
  password?: string;
}) => action(ActionTypes.VERIFY_CODE, payload);
