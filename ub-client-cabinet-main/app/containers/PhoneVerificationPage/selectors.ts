import { createSelector } from 'reselect';
import { ApplicationRootState } from 'types';
import { initialState } from './reducer';

/**
 * Direct selector to the phoneVerificationPage state domain
 */

const selectPhoneVerificationPageDomain = (state: ApplicationRootState) => {
  return state.phoneVerificationPage || initialState;
};
const selectGlobalDomain = (state: ApplicationRootState) => {
  return state;
};

/**
 * Other specific selectors
 */

/**
 * Default selector used by PhoneVerificationPage
 */

const makeSelectPhoneVerificationPage = () =>
  createSelector(selectPhoneVerificationPageDomain, substate => {
    return substate;
  });
const makeSelectUserData = () =>
  createSelector(selectGlobalDomain, substate => {
    return substate.acountPage ? substate.acountPage.userData : {};
  });

export default makeSelectPhoneVerificationPage;
export { selectPhoneVerificationPageDomain, makeSelectUserData };
