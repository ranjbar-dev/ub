/*
 *
 * ChangeUserInfoPage reducer
 *
 */

import ActionTypes from './constants';
import { ContainerState, ContainerActions } from './types';

export const initialState: ContainerState = {
  default: null,
  userProfileData: {},
  isLoadingData: true,
};

function changeUserInfoPageReducer(
  state: ContainerState = initialState,
  action: ContainerActions,
): ContainerState {
  switch (action.type) {
    case ActionTypes.DEFAULT_ACTION:
      return state;
    case ActionTypes.SET_USER_PROFILE:
      return {
        ...state,
        userProfileData: action.payload,
        isLoadingData: false,
      };
    case ActionTypes.SET_IS_LOADING:
      return { ...state, isLoadingData: action.payload };
    default:
      return state;
  }
}

export default changeUserInfoPageReducer;
