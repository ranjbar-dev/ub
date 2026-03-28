/*
 *
 * ChangePassword actions
 *
 */

import { action } from 'typesafe-actions';
import { ChangePasswordModel } from './types';

import ActionTypes from './constants';

export const defaultAction = () => action(ActionTypes.DEFAULT_ACTION);
export const changePasswordAction = (payload: ChangePasswordModel) =>
  action(ActionTypes.CHANGE_PASSWORD_ACTION, payload);
export const isChangingPasswordAction = (payload: boolean) =>
  action(ActionTypes.IS_CHANGING_PASSWORD_ACTION, payload);
