import { createSelector } from 'reselect';
import { ApplicationRootState } from 'types';
import { initialState } from './reducer';

/**
 * Direct selector to the documentVerificationPage state domain
 */

const selectDocumentVerificationPageDomain = (state: ApplicationRootState) => {
  return state || initialState;
};

/**
 * Other specific selectors
 */

/**
 * Default selector used by DocumentVerificationPage
 */

const makeSelectDocumentVerificationPage = () =>
  createSelector(
    selectDocumentVerificationPageDomain,
    substate => {
      return substate;
    },
  );
const makeSelectUserProfileData = () =>
  createSelector(
    selectDocumentVerificationPageDomain,
    substate => {
      return substate.documentVerificationPage
        ? substate.documentVerificationPage.userProfileData
        : {};
    },
  );
const makeSelectIsLoadingUserProfileData = () =>
  createSelector(
    selectDocumentVerificationPageDomain,
    substate => {
      return substate.documentVerificationPage
        ? substate.documentVerificationPage.isLoadingData
        : true;
    },
  );

export default makeSelectDocumentVerificationPage;
export {
  selectDocumentVerificationPageDomain,
  makeSelectUserProfileData,
  makeSelectIsLoadingUserProfileData,
};
