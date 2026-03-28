/*
 *
 * ContactUsPage reducer
 *
 */

import ActionTypes from './constants';
import { ContainerState, ContainerActions } from './types';

export const initialState: ContainerState = {
  default: null,
  counterValue: 0,
  inputValue: '',
};

function contactUsPageReducer(
  state: ContainerState = initialState,
  action: ContainerActions,
): ContainerState {
  switch (action.type) {
    case ActionTypes.DEFAULT_ACTION:
      return state;

    case ActionTypes.INCREMENT: {
      const currentValue = state.counterValue;
      const nextValue = currentValue + 1;
      return { ...state, counterValue: nextValue };
    }

    case ActionTypes.DECREMENT:
      const currentValue = state.counterValue;
      const nextValue = currentValue - 1;
      return { ...state, counterValue: nextValue };

    case ActionTypes.CHANGE_INPUT_VALUE:
      return { ...state, inputValue: action.payload };

    case ActionTypes.ADD_BY_INPUT_VALUE:
      const currentValueToAdd = state.counterValue;
      const newValueToSubtract = currentValueToAdd + action.payload;
      return { ...state, counterValue: newValueToSubtract };

    case ActionTypes.SUBTRACT_BY_INPUT_VALUE:
      const currentValueToSubtract = state.counterValue;
      const newValue = currentValueToSubtract - action.payload;
      return { ...state, counterValue: newValue };

    case ActionTypes.RECIEVE_NUMBER_FROM_API:
      return { ...state, counterValue: action.payload };

    default:
      return state;
  }
}

export default contactUsPageReducer;
