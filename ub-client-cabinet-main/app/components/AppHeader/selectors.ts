import { createSelector } from 'reselect';
import { ApplicationRootState } from 'types';
import {CookieKeys, cookies} from 'services/cookie';

/**
 * Direct selector to the signupPage state domain
 */

const AppHeaderDomain = (state: ApplicationRootState) => {
  return state;
};

/**
 * Other specific selectors
 */

/**
 * Default selector used by SignupPage
 */

const makeSelectLocation = () =>
  createSelector(AppHeaderDomain, (substate) => {
    return substate.router.location;
  });
const makeSelectLoggedIn = () =>
  createSelector(AppHeaderDomain, (substate) => {
    return (
      substate.global?.loggedIn === true ||
      cookies.get(CookieKeys.Token)!=null
    );
  });

export { AppHeaderDomain, makeSelectLocation, makeSelectLoggedIn };
