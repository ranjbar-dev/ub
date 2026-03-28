/**
 * Homepage selectors
 */

import { createSelector } from 'reselect';
import { initialState } from './reducer';
import { ApplicationRootState } from 'types';

const selectHome = (state: ApplicationRootState) => {
  return state.home || initialState;
};
const selectGlobal = (state: ApplicationRootState) => {
  return state.global || initialState;
};
const makeSelectLoggedIn = () =>
  createSelector(
    selectGlobal,
    substate => {
      return substate.loggedIn;
    },
  );

export { selectHome, makeSelectLoggedIn };
