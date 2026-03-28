/*
 *
 * ContactUsPage actions
 *
 */

import { action } from 'typesafe-actions';
import {} from './types';

import ActionTypes from './constants';

export const defaultAction = () => action(ActionTypes.DEFAULT_ACTION);
export const incrementAction = () => action(ActionTypes.INCREMENT);
export const decrementAction = () => action(ActionTypes.DECREMENT);
export const getNumberFromApi = () => action(ActionTypes.GET_NUMBER_FROM_API);
export const recieveNumberFromApi = (payload: number) =>
  action(ActionTypes.RECIEVE_NUMBER_FROM_API, payload);

export const changeInputValueAction = (payload: string) =>
  action(ActionTypes.CHANGE_INPUT_VALUE, payload);

export const addValueByInputAction = (payload: number) =>
  action(ActionTypes.ADD_BY_INPUT_VALUE, payload);

export const subtractValueByInputAction = (payload: number) =>
  action(ActionTypes.SUBTRACT_BY_INPUT_VALUE, payload);
