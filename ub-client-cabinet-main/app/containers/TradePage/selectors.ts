import { createSelector } from 'reselect';
import { ApplicationRootState } from 'types';
import { initialState } from './reducer';

/**
 * Direct selector to the tradePage state domain
 */

const selectTradePageDomain = (state: ApplicationRootState) => {
  return state || initialState;
};

const selectPairMapDomain = (state: ApplicationRootState) => {
  return state?.tradePage?.pairsMap || {};
};

/**
 * Other specific selectors
 */

/**
 * Default selector used by TradePage
 */

const makeSelectTradePage = () =>
  createSelector(selectTradePageDomain, substate => {
    return substate;
  });

const makeSelectPairMap = () =>
  createSelector([selectPairMapDomain], pairMap => pairMap);


export default makeSelectTradePage;
export { selectTradePageDomain, makeSelectPairMap };
