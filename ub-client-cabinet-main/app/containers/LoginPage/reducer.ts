/*
 *
 * LoginPage reducer
 *
 */

import ActionTypes from './constants';
import { ContainerState, ContainerActions } from './types';

export const initialState: ContainerState = {
  default: null,
  isLoggingIn: false,
};

function loginPageReducer(
  state: ContainerState = initialState,
  action: ContainerActions,
): ContainerState {
  switch (action.type) {
    case ActionTypes.IS_LOGGING_IN:
      return { ...state, isLoggingIn: action.payload };
    default:
      return state;
  }
}

export default loginPageReducer;
