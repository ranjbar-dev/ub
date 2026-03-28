/*
 *
 * GoogleAuthenticationPage reducer
 *
 */

import ActionTypes from './constants';
import { ContainerState, ContainerActions } from './types';

export const initialState: ContainerState = {
  default: null,
  isLoading: true,
  qrCode: {},
};

function googleAuthenticationPageReducer (
  state: ContainerState = initialState,
  action: ContainerActions,
): ContainerState {
  switch (action.type) {
    case ActionTypes.DEFAULT_ACTION:
      return state;
    case ActionTypes.SET_IS_LOADING:
      return { ...state, isLoading: action.payload };
    case ActionTypes.SET_2FA_QRCODE:
      return { ...state, qrCode: action.payload };
    default:
      return state;
  }
}

export default googleAuthenticationPageReducer;
