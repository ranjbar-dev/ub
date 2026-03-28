/*
 *
 * TradePage actions
 *
 */

import { action } from 'typesafe-actions';
import { PairItem } from './types';

import ActionTypes from './constants';

export const defaultAction = () => action(ActionTypes.DEFAULT_ACTION);

export const getCurrenciesAction = () => action(ActionTypes.GET_CURRENCIES);

export const getInitialMarketTradeDataAction = (payload: {
  pairName: string;
}) => action(ActionTypes.GET_INITIAL_MARKET_TRADES_DATA, payload);

export const setPairMapAction = (payload: { [key: string]: PairItem }) => action(ActionTypes.SET_PAIR_MAP, payload);

export const AddRemoveFavoritePair = (payload: {
  pair_currency_id: number;
  action: 'add' | 'remove';
}) => action(ActionTypes.ADD_REMOVE_FAVORITE_PAIR, payload);
