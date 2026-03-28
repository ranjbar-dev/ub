import { createSelector } from 'reselect';
import { ApplicationRootState } from 'types';
import { initialState } from './reducer';

/**
 * Direct selector to the signupPage state domain
 */

const selectSignupPageDomain = (state: ApplicationRootState) => {
  return state || initialState;
};

/**
 * Other specific selectors
 */

/**
 * Default selector used by SignupPage
 */

const makeSelectSignupPage = () =>
  createSelector(selectSignupPageDomain, substate => {
    return substate;
  });

export default makeSelectSignupPage;
export { selectSignupPageDomain };
