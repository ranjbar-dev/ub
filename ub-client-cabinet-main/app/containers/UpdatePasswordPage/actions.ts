/*
 *
 * UpdatePasswordPage actions
 *
 */

import { action } from 'typesafe-actions';
import { UpdatePasswordModel } from './types';

import ActionTypes from './constants';

export const defaultAction = () => action(ActionTypes.DEFAULT_ACTION);
export const resetPasswordAction = (payload: UpdatePasswordModel) =>
  action(ActionTypes.RESET_PASSWORD_ACTION, payload);
