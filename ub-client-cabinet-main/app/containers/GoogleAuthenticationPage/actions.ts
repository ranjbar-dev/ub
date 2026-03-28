/*
 *
 * GoogleAuthenticationPage actions
 *
 */

import { action } from 'typesafe-actions';
import { QrCode, SetG2FaModel } from './types';

import ActionTypes from './constants';

export const defaultAction = () => action(ActionTypes.DEFAULT_ACTION);
export const get2faQrCodeAction = () => action(ActionTypes.GET_2FA_QRCODE);

export const setIsLoadingAction = (payload: boolean) =>
  action(ActionTypes.SET_IS_LOADING, payload);
export const set2faQrCodeAction = (payload: QrCode) =>
  action(ActionTypes.SET_2FA_QRCODE, payload);
export const toggle2FaAction = (payload: SetG2FaModel) =>
  action(ActionTypes.TOGGLE_2FA, payload);
