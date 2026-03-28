import { createSelector } from 'reselect';
import { ApplicationRootState } from 'types';
import { initialState } from './reducer';

/**
 * Direct selector to the changePassword state domain
 */

const selectChangePasswordDomain = (state: ApplicationRootState) => {
  return state || initialState;
};

/**
 * Other specific selectors
 */

/**
 * Default selector used by ChangePassword
 */

const makeSelectChangePassword = () =>
  createSelector(
    selectChangePasswordDomain,
    substate => {
      return substate;
    },
  );
const makeSelectTheme = () =>
  createSelector(
    selectChangePasswordDomain,
    substate => {
      return substate.global.theme;
    },
  );
const makeSelectIsLoading = () =>
  createSelector(
    selectChangePasswordDomain,
    substate => {
      return (
        substate.changePassword && substate.changePassword.isChangingPassword
      );
    },
  );

export default makeSelectChangePassword;
export { selectChangePasswordDomain, makeSelectTheme, makeSelectIsLoading };
