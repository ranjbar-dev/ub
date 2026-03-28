/*
 *
 * EmailAuthentication actions
 *
 */

import { action } from 'typesafe-actions';
import {} from './types';

import ActionTypes from './constants';

export const defaultAction = () => action(ActionTypes.DEFAULT_ACTION);
export const acountActivationAction = (payload: { code: string }) =>
  action(ActionTypes.ACOUNT_ACTIVATION_ACTION, payload);
