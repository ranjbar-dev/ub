import { createSelector } from 'reselect';
import { ApplicationRootState } from 'types';
import { initialState } from './reducer';

/**
 * Direct selector to the updatePasswordPage state domain
 */

const selectUpdatePasswordPageDomain = (state: ApplicationRootState) => {
  return state || initialState;
};

/**
 * Other specific selectors
 */

/**
 * Default selector used by UpdatePasswordPage
 */

const makeSelectUpdatePasswordPage = () =>
  createSelector(selectUpdatePasswordPageDomain, substate => {
    return substate;
  });
const makeSelectLocation = () =>
  createSelector(selectUpdatePasswordPageDomain, substate => {
    return substate.router.location;
  });

export default makeSelectUpdatePasswordPage;
export { selectUpdatePasswordPageDomain, makeSelectLocation };
