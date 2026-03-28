import { createSelector } from 'reselect';
import { ApplicationRootState } from 'types';
import { initialState } from './reducer';

/**
 * Direct selector to the fundsPage state domain
 */

const selectFundsPageDomain = (state: ApplicationRootState) => {
  return state || initialState;
};

/**
 * Other specific selectors
 */

/**
 * Default selector used by FundsPage
 */

const makeSelectFundsPage = () =>
  createSelector(selectFundsPageDomain, substate => {
    return substate;
  });
const makeSelectBalancePageData = () =>
  createSelector(selectFundsPageDomain, substate => {
    return substate.fundsPage ? substate.fundsPage.balancePageData : {};
  });
const makeSelectIsLoadingBalancePageData = () =>
  createSelector(selectFundsPageDomain, substate => {
    return substate.fundsPage
      ? substate.fundsPage.isLoadingBalancePageData
      : true;
  });
const makeSelectIsLoadingdepositAndWithDrawData = () =>
  createSelector(selectFundsPageDomain, substate => {
    return substate.fundsPage
      ? substate.fundsPage.isLoadingDepositeAndWithdraw
      : true;
  });
const makeSelectdepositAndWithDrawData = () =>
  createSelector(selectFundsPageDomain, substate => {
    return substate.fundsPage ? substate.fundsPage.depositAndWithDrawData : {};
  });
const makeSelectFormerWithdrawAddresses = () =>
  createSelector(selectFundsPageDomain, substate => {
    return substate.fundsPage ? substate.fundsPage.formerWithdrawAddresses : [];
  });

const makeSelectTransactionHistoryPageData = () =>
  createSelector(selectFundsPageDomain, substate => {
    return substate.fundsPage ? substate.fundsPage.transactionHistoryData : [];
  });
const makeSelectIsLoadingTransactionHistoryPageData = () =>
  createSelector(selectFundsPageDomain, substate => {
    return substate.fundsPage
      ? substate.fundsPage.isLoadingTransactionHistory
      : true;
  });
const makeSelectUserData = () =>
  createSelector(selectFundsPageDomain, substate => {
    return substate.fundsPage?.userData;
  });
const makeSelectLocation = () =>
  createSelector(selectFundsPageDomain, substate => {
    return substate.router ? substate.router.location : {};
  });

export default makeSelectFundsPage;
export {
  selectFundsPageDomain,
  makeSelectBalancePageData,
  makeSelectIsLoadingBalancePageData,
  makeSelectdepositAndWithDrawData,
  makeSelectIsLoadingdepositAndWithDrawData,
  makeSelectTransactionHistoryPageData,
  makeSelectIsLoadingTransactionHistoryPageData,
  makeSelectFormerWithdrawAddresses,
  makeSelectLocation,
  makeSelectUserData,
};
