/*
 *
 * RecapchaContainer reducer
 *
 */


import { ContainerState, ContainerActions } from './types';

export const initialState: ContainerState = {
  default: null,
  recapcha: '',
};

function recapchaContainerReducer(
  state: ContainerState = initialState,
  action: ContainerActions,
): ContainerState {
  switch (action.type) {
    //case ActionTypes.DEFAULT_ACTION:
    //  return state;
    //case ActionTypes.SET_RECAPCHA_ACTION:
    //  return { ...state, recapcha: action.payload };
    default:
      return state;
  }
}

export default recapchaContainerReducer;
