import { createSelector } from 'reselect';
import { ApplicationRootState } from 'types';
import { initialState } from './reducer';

/**
 * Direct selector to the contactUsPage state domain
 */

const selectContactUsPageDomain = (state: ApplicationRootState) => {
  return state.contactUsPage || initialState;
};

/**
 * Other specific selectors
 */

/**
 * Default selector used by ContactUsPage
 */

const makeSelectContactUsPage = () =>
  createSelector(selectContactUsPageDomain, (substate) => {
    return substate;
  });
const makeSelectCounter = () =>
  createSelector(selectContactUsPageDomain, (contactusPageState) => {
    return contactusPageState.counterValue;
  });
const makeSelectInputValue = () =>
  createSelector(selectContactUsPageDomain, (contactusPageState) => {
    return contactusPageState.inputValue;
  });

export default makeSelectContactUsPage;
export { selectContactUsPageDomain, makeSelectCounter, makeSelectInputValue };
