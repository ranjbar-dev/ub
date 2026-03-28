import { createSelector } from 'reselect';
import { ApplicationRootState } from 'types';
import { initialState } from './reducer';

/**
 * Direct selector to the changeUserInfoPage state domain
 */

const selectChangeUserInfoPageDomain = (state: ApplicationRootState) => {
  return state || initialState;
};

/**
 * Other specific selectors
 */

/**
 * Default selector used by ChangeUserInfoPage
 */
const makeSelectUserProfileData = () =>
  createSelector(
    selectChangeUserInfoPageDomain,
    substate => {
      return substate.changeUserInfoPage
        ? substate.changeUserInfoPage.userProfileData
        : {};
    },
  );
const makeSelectIsLoadingUserProfileData = () =>
  createSelector(
    selectChangeUserInfoPageDomain,
    substate => {
      return substate.changeUserInfoPage
        ? substate.changeUserInfoPage.isLoadingData
        : true;
    },
  );

const makeSelectChangeUserInfoPage = () =>
  createSelector(
    selectChangeUserInfoPageDomain,
    substate => {
      return substate;
    },
  );

export default makeSelectChangeUserInfoPage;
export {
  selectChangeUserInfoPageDomain,
  makeSelectUserProfileData,
  makeSelectIsLoadingUserProfileData,
};
