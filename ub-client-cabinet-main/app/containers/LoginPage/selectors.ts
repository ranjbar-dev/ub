import { createSelector } from 'reselect';
import { ApplicationRootState } from 'types';
import { initialState } from './reducer';

/**
 * Direct selector to the loginPage state domain
 */

const selectLoginPageDomain = (state: ApplicationRootState) => {
  return state.loginPage || initialState;
};
const selectAppDomain = (state: ApplicationRootState) => {
  return state || initialState;
};

/**
 * Other specific selectors
 */

/**
 * Default selector used by LoginPage
 */

const makeSelectLoginPage = () =>
  createSelector(selectLoginPageDomain, (substate) => {
    return substate;
  });

const makeSelectIsLoadingLoginPage = () =>
  createSelector(selectLoginPageDomain, (loginPageState) => {
    return loginPageState.isLoggingIn;
  });

const makeSelectLoggedIn = () =>
  createSelector(selectAppDomain, (substate) => {
    return substate.global.loggedIn;
  });

export default makeSelectLoginPage;
export {
  selectLoginPageDomain,
  makeSelectLoggedIn,
  makeSelectIsLoadingLoginPage,
};
