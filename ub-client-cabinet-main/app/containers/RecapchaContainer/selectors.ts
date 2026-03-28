import { createSelector } from 'reselect';
import { ApplicationRootState } from 'types';
import { initialState } from './reducer';

/**
 * Direct selector to the recapchaContainer state domain
 */

const selectRecapchaContainerDomain = (state: ApplicationRootState) => {
  return state || initialState;
};

/**
 * Other specific selectors
 */

/**
 * Default selector used by RecapchaContainer
 */

const makeSelectRecapchaContainer = () =>
  createSelector(selectRecapchaContainerDomain, substate => {
    return substate;
  });
const makeSelectRecapcha = () =>
  createSelector(selectRecapchaContainerDomain, substate => {
    return substate.recapchaContainer
      ? substate.recapchaContainer.recapcha
      : '';
  });

export default makeSelectRecapchaContainer;
export { selectRecapchaContainerDomain, makeSelectRecapcha };
