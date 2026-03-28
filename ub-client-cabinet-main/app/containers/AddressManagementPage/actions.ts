/*
 *
 * AddressManagementPage actions
 *
 */

import { action } from 'typesafe-actions';
import { Currency, WithdrawAddress } from './types';

import ActionTypes from './constants';

export const initialAction = () => action(ActionTypes.INITIAL_ACTION);

export const setIsLoadingAction = (payload: boolean) =>
  action(ActionTypes.IS_LOADING_ACTION, payload);

export const setCurrenciesAction = (payload: Currency[]) =>
  action(ActionTypes.SET_CURRENCIES_ACTION, payload);

export const setWithdrawAddressesAction = (payload: WithdrawAddress[]) =>
  action(ActionTypes.SET_WITHDRAW_ADDRESS_ACTION, payload);

export const addOneToWithdrawAddressesAction = (payload: WithdrawAddress) =>
  action(ActionTypes.ADD_ONE_TO_ADDRESSES_ACTION, payload);

export const deleteAddressAction = (payload: {
  data: WithdrawAddress;
  rowIndex: number;
}) => action(ActionTypes.DELETE_ADDRESS_ACTION, payload);

export const applyDeleteAddressAction = (payload: {
  data: WithdrawAddress;
  rowIndex: number;
}) => action(ActionTypes.APPLY_DELETE_ADDRESS_ACTION, payload);

export const favoriteAddressAction = (payload: {
  data: {
    action: string;
    id: number;
  };
  rowIndex: number;
}) => action(ActionTypes.FAVORITE_ADDRESS_ACTION, payload);

export const applyfavoriteAddressAction = (payload: {
  data: {
    action: string;
    id: number;
  };
  rowIndex: number;
}) => action(ActionTypes.APPLY_FAVORITE_ADDRESS_ACTION, payload);

export const addNewAddressAction = (payload: {
  address: string;
  code: string;
  label: string;
}) => action(ActionTypes.ADD_NEW_ADDRESS_ACTION, payload);
