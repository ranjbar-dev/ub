/*
 *
 * TradePage reducer
 *
 */

import ActionTypes from './constants';
import { ContainerState, ContainerActions } from './types';

export const initialState: ContainerState = {
  default: null,
  pairsMap: {}
};

function tradePageReducer(
  state: ContainerState = initialState,
  action: ContainerActions,
): ContainerState {
  switch (action.type) {
    case ActionTypes.DEFAULT_ACTION:
      return state;
    case ActionTypes.GET_INITIAL_MARKET_TRADES_DATA:
      return state;
    case ActionTypes.SET_PAIR_MAP:
      return { ...state, pairsMap: action.payload };
    default:
      return state;
  }
}

export default tradePageReducer;
