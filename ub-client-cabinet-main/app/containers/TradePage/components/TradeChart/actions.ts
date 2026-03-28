/*
 *
 * TradeChart actions
 *
 */

import { action } from 'typesafe-actions';
import {} from './types';

import ActionTypes from './constants';

export const defaultAction = () => action(ActionTypes.DEFAULT_ACTION);
export const getChartConfigAction = () => action(ActionTypes.GET_CHART_CONFIG);
export const setChartConfigAction = (payload: any) =>
  action(ActionTypes.SET_CHART_CONFIG, payload);
