/*
 *
 * SignupPage actions
 *
 */

import { action } from 'typesafe-actions';
import { RegisterModel } from './types';

import ActionTypes from './constants';

export const defaultAction = () => action(ActionTypes.DEFAULT_ACTION);
export const registerAction = (payload: RegisterModel) =>
  action(ActionTypes.REGISTER_ACTION, payload);
