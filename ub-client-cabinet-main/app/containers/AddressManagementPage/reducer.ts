/*
 *
 * AddressManagementPage reducer
 *
 */

import ActionTypes from './constants';
import { ContainerState, ContainerActions, WithdrawAddress } from './types';

export const initialState: ContainerState = {
  default: null,
  isLoading: false,
  currencies: [],
  withDrawArrdesses: [],
  // isAddingAddress: false,
};

function addressManagementPageReducer(
  state: ContainerState = initialState,
  action: ContainerActions,
): ContainerState {
  switch (action.type) {
    case ActionTypes.INITIAL_ACTION:
      return state;
    case ActionTypes.IS_LOADING_ACTION:
      return { ...state, isLoading: action.payload };

    case ActionTypes.SET_CURRENCIES_ACTION:
      return { ...state, currencies: action.payload };

    case ActionTypes.SET_WITHDRAW_ADDRESS_ACTION:
      return { ...state, withDrawArrdesses: action.payload, isLoading: false };

    case ActionTypes.ADD_ONE_TO_ADDRESSES_ACTION:
      return {
        ...state,
        withDrawArrdesses: [action.payload, ...state.withDrawArrdesses],
      };
    case ActionTypes.APPLY_DELETE_ADDRESS_ACTION:
      const newList: WithdrawAddress[] = [];
      for (let i = 0; i < state.withDrawArrdesses.length; i++) {
        if (action.payload.data.id != state.withDrawArrdesses[i].id) {
          newList.push(state.withDrawArrdesses[i]);
        }
      }
      return {
        ...state,
        withDrawArrdesses: newList,
      };
    case ActionTypes.APPLY_FAVORITE_ADDRESS_ACTION:
      const myList: WithdrawAddress[] = state.withDrawArrdesses;
      for (let i = 0; i < state.withDrawArrdesses.length; i++) {
        if (action.payload.data.id === state.withDrawArrdesses[i].id) {
          // myList.push(state.withDrawArrdesses[i]);
          myList[i].isFavorite =
            action.payload.data.action == 'add' ? true : false;
        }
      }
      return {
        ...state,
        withDrawArrdesses: myList,
      };
    default:
      return state;
  }
}

export default addressManagementPageReducer;
