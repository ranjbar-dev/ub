/*
 *
 * FundsPage reducer
 *
 */

import ActionTypes from './constants';
import { ContainerState, ContainerActions } from './types';

export const initialState: ContainerState = {
  default: null,
  balancePageData: {},
  userData: {},
  isLoadingBalancePageData: true,
  isLoadingDepositeAndWithdraw: true,
  isLoadingTransactionHistory: true,
  depositAndWithDrawData: {},
  transactionHistoryData: [],
  formerWithdrawAddresses: [],
};

function fundsPageReducer (
  state: ContainerState = initialState,
  action: ContainerActions,
): ContainerState {
  switch (action.type) {
    case ActionTypes.DEFAULT_ACTION:
      return state;

    case ActionTypes.SET_BALANCE_PAGE_DATA_ACTION:
      return {
        ...state,
        balancePageData: action.payload,
        isLoadingBalancePageData: false,
      };
    case ActionTypes.SET_IS_LOADING_BALANCE_PAGE_DATA_ACTION:
      return { ...state, isLoadingBalancePageData: action.payload };

    case ActionTypes.SET_DEPOSITE_AND_WITHDRAWS_DATA_ACTION:
      return {
        ...state,
        depositAndWithDrawData: action.payload,
        isLoadingDepositeAndWithdraw: false,
      };
    case ActionTypes.SET_IS_LOADING_DEPOSITE_AND_WITHDRAWS_DATA_ACTION:
      return { ...state, isLoadingDepositeAndWithdraw: action.payload };

    case ActionTypes.SET_TRANSACTION_HISTORY_PAGE_DATA_ACTION:
      return {
        ...state,
        transactionHistoryData: action.payload,
        isLoadingTransactionHistory: false,
      };
    case ActionTypes.SET_IS_LOADING_TRANSACTION_HISTORY_PAGE_DATA_ACTION:
      return { ...state, isLoadingTransactionHistory: action.payload };

    case ActionTypes.SET_FORMER_WITHDRAW_ADDRESSES:
      return { ...state, formerWithdrawAddresses: action.payload };
    case ActionTypes.ADD_FORMER_WITHDRAW_ADDRESSES:
      const former = state.formerWithdrawAddresses;
      former.unshift(action.payload);
      return { ...state, formerWithdrawAddresses: former };

    case ActionTypes.SET_USER_DATA_ACTION:
      const pageData = { ...state };

      return { ...pageData, userData: action.payload };

    case ActionTypes.ADD_WITHDRAWS_DATA_ACTION:
      const withdraws = state.depositAndWithDrawData
        ? state.depositAndWithDrawData.withdrawTransactions
        : [];
      const depositAndWithDrawData = state.depositAndWithDrawData;
      if (withdraws) {
        withdraws.unshift(action.payload);
      }
      depositAndWithDrawData.withdrawTransactions = withdraws;
      return {
        ...state,
        depositAndWithDrawData,
        transactionHistoryData:
          state.transactionHistoryData.length !== 0
            ? [action.payload, ...state.transactionHistoryData]
            : [],
      };
    default:
      return state;
  }
}

export default fundsPageReducer;
