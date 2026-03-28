import { createSelector } from 'reselect';
import { ApplicationRootState } from 'types';
import { initialState } from './reducer';
import {cookies,CookieKeys} from 'services/cookie';

/**
 * Direct selector to the ordersPage state domain
 */

const selectOrdersPageDomain = (state: ApplicationRootState) => {
  return state || initialState;
};

/**
 * Other specific selectors
 */

/**
 * Default selector used by OrdersPage
 */

const makeSelectOrdersPage = () =>
  createSelector(selectOrdersPageDomain, substate => {
    return substate.ordersPage;
  });
const makeSelectLoggedIn = () =>
  createSelector(selectOrdersPageDomain, substate => {
    return (
      substate.global?.loggedIn === true ||
      cookies.get(CookieKeys.Token)!=null
    );
  });
const makeSelectOpenOrders = () =>
  createSelector(selectOrdersPageDomain, substate => {
    return substate.ordersPage ? substate.ordersPage.openOrders : [];
  });
const makeSelectOrderHistory = () =>
  createSelector(selectOrdersPageDomain, substate => {
    return substate.ordersPage ? substate.ordersPage.orderHistory : [];
  });
const makeSelectTradeHistory = () =>
  createSelector(selectOrdersPageDomain, substate => {
    return substate.ordersPage ? substate.ordersPage.tradeHistory : [];
  });
const makeSelectIsLoadingOpenOrders = () =>
  createSelector(selectOrdersPageDomain, substate => {
    return substate.ordersPage ? substate.ordersPage.isLoadingOpenOrders : true;
  });
const makeSelectIsLoadingOrderHistory = () =>
  createSelector(selectOrdersPageDomain, substate => {
    return substate.ordersPage
      ? substate.ordersPage.isLoadingOrderHistory
      : true;
  });
const makeSelectIsLoadingTradeHistory = () =>
  createSelector(selectOrdersPageDomain, substate => {
    return substate.ordersPage
      ? substate.ordersPage.isLoadingTradeHistory
      : true;
  });
const makeSelectCurrencies = () =>
  createSelector(selectOrdersPageDomain, substate => {
    return substate.ordersPage ? substate.ordersPage.currencies : [];
  });

export default makeSelectOrdersPage;
export {
  selectOrdersPageDomain,
  makeSelectIsLoadingOpenOrders,
  makeSelectIsLoadingOrderHistory,
  makeSelectIsLoadingTradeHistory,
  makeSelectOpenOrders,
  makeSelectLoggedIn,
  makeSelectOrderHistory,
  makeSelectTradeHistory,
  makeSelectCurrencies,
};
