/*
 *
 * RecapchaContainer actions
 *
 */

import { action } from 'typesafe-actions';
import {} from './types';

import ActionTypes from './constants';

//export const defaultAction = () => action(ActionTypes.DEFAULT_ACTION);
export const getRecapchaAction = () => action(ActionTypes.GET_RECAPCHA_ACTION);

//export const setRecapchaKeyAction = (payload: string) =>
//  action(ActionTypes.SET_RECAPCHA_ACTION, payload);

// export const setRecapchaTokenAction = (payload: string) =>
//   action(ActionTypes.SET_RECAPCHA_TOKEN_ACTION, payload);
// export const resetRecapchaAction = () =>
//   action(ActionTypes.RESET_RECAPCHA_TOKEN_ACTION);
