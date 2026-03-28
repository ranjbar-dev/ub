/*
 *
 * FundsPage actions
 *
 */

import { action } from 'typesafe-actions';

import ActionTypes from './constants';

import { BalancePageData } from 'containers/FundsPage/types';

export const defaultAction = () => action(ActionTypes.DEFAULT_ACTION);

export const getBalancePageDataAction = () =>
  action(ActionTypes.GET_BALANCE_DATA_ACTION);
export const setBalancePageDataAction = (payLoad: BalancePageData) =>
  action(ActionTypes.SET_BALANCE_DATA_ACTION, payLoad);
