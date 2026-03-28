/*
 *
 * FundsPage reducer
 *
 */

import ActionTypes from './constants';
import { ContainerState, ContainerActions } from './types';

export const initialState: ContainerState = {
  default: null,
  balancePageData: {},
};

function fundsPageReducer(
  state: ContainerState = initialState,
  action: ContainerActions,
): ContainerState {
  switch (action.type) {
    case ActionTypes.DEFAULT_ACTION:
      return state;
    case ActionTypes.SET_BALANCE_DATA_ACTION:
      return {
        ...state,
        balancePageData: action.payload,
      };
    default:
      return state;
  }
}

export default fundsPageReducer;
