/*
 *
 * ChangePassword reducer
 *
 */

import ActionTypes from './constants';
import { ContainerState, ContainerActions } from './types';

export const initialState: ContainerState = {
  default: null,
  isChangingPassword: false,
};

function changePasswordReducer(
  state: ContainerState = initialState,
  action: ContainerActions,
): ContainerState {
  switch (action.type) {
    case ActionTypes.DEFAULT_ACTION:
      return state;
    case ActionTypes.IS_CHANGING_PASSWORD_ACTION:
      return { ...state, isChangingPassword: action.payload };
    default:
      return state;
  }
}

export default changePasswordReducer;
