import { createSelector } from 'reselect';
import { ApplicationRootState } from 'types';
import { initialState } from './reducer';

/**
 * Direct selector to the googleAuthenticationPage state domain
 */

const selectGoogleAuthenticationPageDomain = (state: ApplicationRootState) => {
  return state || initialState;
};

/**
 * Other specific selectors
 *
 * /

/**
 * Default selector used by GoogleAuthenticationPage
 */

const makeSelectGoogleAuthenticationPage = () =>
  createSelector(selectGoogleAuthenticationPageDomain, (substate) => {
    return substate;
  });
const makeSelectQrCode = () =>
  createSelector(selectGoogleAuthenticationPageDomain, (substate) => {
    return substate.googleAuthenticationPage
      ? substate.googleAuthenticationPage.qrCode
      : {};
  });
const makeSelectIsLoading = () =>
  createSelector(selectGoogleAuthenticationPageDomain, (substate) => {
    return substate.googleAuthenticationPage
      ? substate.googleAuthenticationPage.isLoading
      : true;
  });
const makeSelectUserData = () =>
  createSelector(selectGoogleAuthenticationPageDomain, (substate) => {
    return substate.acountPage ? substate.acountPage.userData : {};
  });

export default makeSelectGoogleAuthenticationPage;
export {
  selectGoogleAuthenticationPageDomain,
  makeSelectQrCode,
  makeSelectIsLoading,
  makeSelectUserData,
};
