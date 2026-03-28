/*
 *
 * PhoneVerificationPage reducer
 *
 */

import ActionTypes from './constants';
import { ContainerState, ContainerActions } from './types';

export const initialState: ContainerState = {
  default: null,
  isCountriesLoading: false,
  countries: [],
  isLoading: false,
  enteredPhoneNumber: '',
  activeStep: 0,
};

function phoneVerificationPageReducer(
  state: ContainerState = initialState,
  action: ContainerActions,
): ContainerState {
  switch (action.type) {
    case ActionTypes.DEFAULT_ACTION:
      return state;
    case ActionTypes.SET_COUNTRIES_LOADING:
      return { ...state, isCountriesLoading: action.payload };
    case ActionTypes.SET_COUNTRIES_ACTION:
      return { ...state, countries: action.payload };
    case ActionTypes.SET_IS_SENDING_SMS:
      return { ...state, isLoading: action.payload };
    case ActionTypes.SET_STEP_ACTION:
      return { ...state, activeStep: action.payload };
    case ActionTypes.SET_PHONE_NUMBER_ACTION:
      return { ...state, enteredPhoneNumber: action.payload };
    default:
      return state;
  }
}

export default phoneVerificationPageReducer;
