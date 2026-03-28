import { action } from 'typesafe-actions';
import ActionTypes from './constants';
import { Currency } from 'containers/AddressManagementPage/types';
import { Country } from 'containers/PhoneVerificationPage/constants';

export const defaultAction = (data: any) =>
  action(ActionTypes.DEFAULT_ACTION, data);
export const loggedInAction = (payload: boolean) =>
  action(ActionTypes.LOGGED_IN_ACTION, payload);
export const setGlobalCurrenciesAction = (payload: Currency[]) =>
  action(ActionTypes.SET_CURRENCIES_ACTION, payload);
export const setGlobalCountriesAction = (payload: Country[]) =>
  action(ActionTypes.SET_COUNTRIES_ACTION, payload);
