/*
 *
 * FundsPage actions
 *
 */

import { action } from 'typesafe-actions';
import {
  BalancePageData,
  depositAndWithDrawData,
  Transaction,
  OrderDetail,
  WithdrawModel,
  DWData,
  InfiniteDwModel,
} from './types';

import ActionTypes from './constants';
import { WithdrawAddress } from 'containers/AddressManagementPage/types';
import { UserData } from 'containers/AcountPage/types';
import { FilterModel } from 'containers/OrdersPage/types';

export const defaultAction = () => action(ActionTypes.DEFAULT_ACTION);

export const getBalancePageDataAction = (payload: { isSilent?: boolean }) =>
  action(ActionTypes.GET_BALANCE_PAGE_DATA_ACTION, payload);
export const setBalancePageDataAction = (payLoad: BalancePageData) =>
  action(ActionTypes.SET_BALANCE_PAGE_DATA_ACTION, payLoad);
export const setIsLoadingBalancePageDataAction = (payLoad: boolean) =>
  action(ActionTypes.SET_IS_LOADING_BALANCE_PAGE_DATA_ACTION, payLoad);

export const getDepositAndWithDrawDataAction = (payload: {
  code: string;
  type: string;
}) => action(ActionTypes.GET_DEPOSITE_AND_WITHDRAWS_DATA_ACTION, payload);

export const getRawdepositAndWithDrawDataAction = (payload: {
  code: string;
  type: string;
  fromCoinChange?: boolean;
}) => action(ActionTypes.GET_RAW_DEPOSITE_AND_WITHDRAWS_DATA_ACTION, payload);

export const setDepositAndWithDrawDataAction = (
  payLoad: depositAndWithDrawData,
) => action(ActionTypes.SET_DEPOSITE_AND_WITHDRAWS_DATA_ACTION, payLoad);
export const setUserDataAction = (payLoad: UserData) =>
  action(ActionTypes.SET_USER_DATA_ACTION, payLoad);

export const getUserDataAction = () => action(ActionTypes.GET_USER_DATA_ACTION);

export const addWithdrawDataAction = (payLoad: DWData) =>
  action(ActionTypes.ADD_WITHDRAWS_DATA_ACTION, payLoad);

export const setIsLoadingdepositAndWithDrawDataAction = (payLoad: boolean) =>
  action(
    ActionTypes.SET_IS_LOADING_DEPOSITE_AND_WITHDRAWS_DATA_ACTION,
    payLoad,
  );

export const getTransactionHistoryPageDataAction = (payload?: FilterModel) =>
  action(ActionTypes.GET_TRANSACTION_HISTORY_PAGE_DATA_ACTION, payload);

export const getInfiniteDWAction = (payload: InfiniteDwModel) =>
  action(ActionTypes.GET_INFINITE_DW_ACTION, payload);

export const setTransactionHistoryPageDataAction = (payLoad: Transaction[]) =>
  action(ActionTypes.SET_TRANSACTION_HISTORY_PAGE_DATA_ACTION, payLoad);
export const setIsLoadingTransactionHistoryPageDataAction = (
  payLoad: boolean,
) =>
  action(
    ActionTypes.SET_IS_LOADING_TRANSACTION_HISTORY_PAGE_DATA_ACTION,
    payLoad,
  );

export const getOrderDetailAction = (payLoad: { id: number; rowId: string }) =>
  action(ActionTypes.GET_PAYMENT_DETAIL_ACTION, payLoad);
export const setOrderDetailAction = (payLoad: OrderDetail) =>
  action(ActionTypes.SET_PAYMENT_DETAIL_ACTION, payLoad);

export const getFormerWithdrawAddressesAction = (payload: { code: string }) =>
  action(ActionTypes.GET_FORMER_WITHDRAW_ADDRESSES, payload);
export const setFormerWithdrawAddressesAction = (payload: WithdrawAddress[]) =>
  action(ActionTypes.SET_FORMER_WITHDRAW_ADDRESSES, payload);
export const addFormerWithdrawAddressesAction = (payload: WithdrawAddress) =>
  action(ActionTypes.ADD_FORMER_WITHDRAW_ADDRESSES, payload);

export const withdrawAction = (payload: WithdrawModel) =>
  action(ActionTypes.WITHDRAW_ACTION, payload);

export const preWithdrawAction = (payload: WithdrawModel) =>
  action(ActionTypes.PRE_WITHDRAW_ACTION, payload);

export const addNewAddressAction = (payload: {
  address: string;
  code: string;
  label: string;
}) => action(ActionTypes.ADD_NEW_ADDRESS_ACTION, payload);
