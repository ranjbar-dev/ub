/*
 *
 * AcountPage reducer
 *
 */

import ActionTypes from './constants';
import { ContainerState, ContainerActions } from './types';

export const initialState: ContainerState = {
  default: null,
  isLoading: true,
  userData: null,
};

function acountPageReducer(
  state: ContainerState = initialState,
  action: ContainerActions,
): ContainerState {
  switch (action.type) {
    case ActionTypes.IS_LOADING_ACTION:
      return { ...state, isLoading: action.payload };
    // case ActionTypes.LOGGED_IN_ACTION:
    //   console.log(action);
    //   if (action.payload === false) {
    //     return { ...state, isLoading: true, default: null, userData: null };
    //   }
    //   return state;
    case ActionTypes.SET_USER_DATA_ACTION:
      return { ...state, userData: action.payload, isLoading: false };
    case ActionTypes.SET_2FA_ENABLED_ACTION:
      const data = state.userData;
      if (data) {
        data.google2faEnabled = action.payload;
      }
      return { ...state, userData: data };
    default:
      return state;
  }
}

export default acountPageReducer;
