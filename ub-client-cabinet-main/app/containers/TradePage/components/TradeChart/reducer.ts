/*
 *
 * TradeChart reducer
 *
 */

import ActionTypes from './constants';
import { ContainerState, ContainerActions } from './types';

export const initialState: ContainerState = {
  default: null,
  chartConfig: {},
};

function tradeChartReducer(
  state: ContainerState = initialState,
  action: ContainerActions,
): ContainerState {
  switch (action.type) {
    case ActionTypes.DEFAULT_ACTION:
      return state;
    case ActionTypes.SET_CHART_CONFIG:
      return { ...state, chartConfig: action.payload };
    default:
      return state;
  }
}

export default tradeChartReducer;
