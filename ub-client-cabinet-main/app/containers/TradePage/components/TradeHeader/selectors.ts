import { createSelector } from 'reselect';
import { ApplicationRootState } from 'types';
import { initialState } from 'containers/FundsPage/reducer';
import {cookies,CookieKeys} from 'services/cookie';

/**
 * Direct selector to the tradeHeader state domain
 */

const selectTradeHeaderDomain = (state: ApplicationRootState) => {
  return state || initialState;
};

/**
 * Other specific selectors
 */

/**
 * Default selector used by TradeHeader
 */

const makeSelectTradeHeader = () =>
  createSelector(selectTradeHeaderDomain, (substate) => {
    return substate;
  });
const makeSelectLoggedIn = () =>
  createSelector(selectTradeHeaderDomain, (substate) => {
    return (
      substate.global?.loggedIn === true ||cookies.get(CookieKeys.Token)!=null
    );
  });
const makeSelectBalances = () =>
  createSelector(selectTradeHeaderDomain, (substate) => {
    return substate.tradeHeader ? substate.tradeHeader.balancePageData : {};
  });

export default makeSelectTradeHeader;
export { selectTradeHeaderDomain, makeSelectLoggedIn, makeSelectBalances };
