/*
 *
 * OrdersPage actions
 *
 */

import { action } from 'typesafe-actions';
import {
  Order,
  OrderHistorySearchModel,
  OrderDetail,
  FilterModel,
  StreamOrder,
} from './types';

import ActionTypes from './constants';
import { Currency } from 'containers/AddressManagementPage/types';
import { NewOrderModel } from 'containers/TradePage/components/newOrder/types';

export const defaultAction = () => action(ActionTypes.DEFAULT_ACTION);

export const getOpenOrdersAction = (payload?: { silent: boolean }) =>
  action(ActionTypes.GET_OPEN_ORDERS_ACTION, payload);

export const getOrderHistoryAction = (payload?: { silent: boolean }) =>
  action(ActionTypes.GET_ORDER_HISTORY_ACTION, payload);

export const getPaginatedOrderHistoryAction = (payload: any) =>
  action(ActionTypes.GET_PAGINATED_ORDER_HISTORY_ACTION, payload);
export const getPaginatedTradeHistoryAction = (payload: any) =>
  action(ActionTypes.GET_PAGINATED_TRADE_HISTORY_ACTION, payload);
//export const setPaginatedOrderHistoryAction = (payload: any) =>
//  action(ActionTypes.SET_PAGINATED_ORDER_HISTORY_ACTION, payload);

export const getFilteredOrderHistoryAction = (
  payload: OrderHistorySearchModel,
) => action(ActionTypes.GET_FILTERED_ORDER_HISTORY_ACTION, payload);

export const getTradeHistoryAction = (payload?: FilterModel) =>
  action(ActionTypes.GET_TRADE_HISTORY_ACTION, payload);

export const setCurrenciesAction = (payload: Currency[]) =>
  action(ActionTypes.SET_CURRENCIES, payload);

export const setOpenOrdersAction = (payload: Order[]) =>
  action(ActionTypes.SET_OPEN_ORDERS_ACTION, payload);
export const setOrderHistoryAction = (payload: Order[]) =>
  action(ActionTypes.SET_ORDER_HISTORY_ACTION, payload);
export const setTradeHistoryAction = (payload: Order[]) =>
  action(ActionTypes.SET_TRADE_HISTORY_ACTION, payload);

export const setIsLoadingOpenOrdersAction = (payload: boolean) =>
  action(ActionTypes.SET_IS_LOADING_OPEN_ORDERS_ACTION, payload);

export const setIsLoadingOrderHistoryAction = (payload: boolean) =>
  action(ActionTypes.SET_IS_LOADING_ORDER_HISTORY_ACTION, payload);

export const setIsLoadingTradeHistoryAction = (payload: boolean) =>
  action(ActionTypes.SET_IS_LOADING_TRADE_HISTORY_ACTION, payload);

export const getOrderDetailAction = (payLoad: {
  order_id: number;
  rowId: string;
}) => action(ActionTypes.GET_PAYMENT_DETAIL_ACTION, payLoad);
export const setOrderDetailAction = (payLoad: OrderDetail) =>
  action(ActionTypes.SET_PAYMENT_DETAIL_ACTION, payLoad);

////trade new order

export const getCurrencyPairInfoAction = (payload: {
  pair_currency_id: number;
}) => action(ActionTypes.GET_CURRENCY_PAIR_INFO, payload);

export const createNewOrderAction = (payload: NewOrderModel) =>
  action(ActionTypes.CREATE_NEW_ORDER_ACTION, payload);
export const addNewOrderAction = (payload: Order) =>
  action(ActionTypes.ADD_NEW_ORDER_ACTION, payload);
export const cancelOrderAction = (payload: Order) =>
  action(ActionTypes.CANCEL_ORDER_ACTION, payload);
