import { ContainerState, ContainerActions } from './types';
import ActionTypes, { Themes } from './constants';

// The initial state of the App
export const initialState: ContainerState = {
  loading: false,
  error: false,
  theme: Themes.LIGHT,
  loggedIn: false,
  currencies: [],
  countries: [],
};

// Take this container's state (as a slice of root state), this container's actions and return new state
function appReducer(
  state: ContainerState = initialState,
  action: ContainerActions,
): ContainerState {
  switch (action.type) {
    case ActionTypes.DEFAULT_ACTION:
      return state;
    case ActionTypes.LOGGED_IN_ACTION:
      return { ...state, loggedIn: action.payload };
    case ActionTypes.SET_CURRENCIES_ACTION:
      return { ...state, currencies: action.payload };
    case ActionTypes.SET_COUNTRIES_ACTION:
      return { ...state, countries: action.payload };
    default:
      return state;
  }
}

export default appReducer;
