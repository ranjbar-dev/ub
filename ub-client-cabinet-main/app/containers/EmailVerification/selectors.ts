import { createSelector } from 'reselect';
import { ApplicationRootState } from 'types';
import { initialState } from './reducer';

/**
 * Direct selector to the emailAuthentication state domain
 */

const selectEmailAuthenticationDomain = (state: ApplicationRootState) => {
  return state || initialState;
};

/**
 * Other specific selectors
 */

/**
 * Default selector used by EmailAuthentication
 */

const makeSelectEmailAuthentication = () =>
  createSelector(selectEmailAuthenticationDomain, substate => {
    return substate;
  });
const makeSelectLocation = () =>
  createSelector(selectEmailAuthenticationDomain, substate => {
    return substate.router.location;
  });

export default makeSelectEmailAuthentication;
export { selectEmailAuthenticationDomain, makeSelectLocation };
