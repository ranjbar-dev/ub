import { createSelector } from 'reselect';
import { ApplicationRootState } from 'types';
import { initialState } from './reducer';

/**
 * Direct selector to the acountPage state domain
 */

const selectAcountPageDomain = (state: ApplicationRootState) => {
  return state || initialState;
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

/**
 * Default selector used by AcountPage
 */

const makeSelectAcountPage = () =>
  createSelector(
    selectAcountPageDomain,
    substate => {
      return substate.acountPage;
    },
  );

export default makeSelectAcountPage;
export { selectAcountPageDomain, makeSelectLoggedIn };
