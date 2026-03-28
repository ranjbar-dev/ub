import { apiService } from './api_service';
import { RequestTypes } from './constants';
import { OrderHistorySearchModel } from 'containers/OrdersPage/types';
import { NewOrderModel } from 'containers/TradePage/components/newOrder/types';
export const getOpenOrdersAPI = () => {
  return apiService.fetchData({
    data: {},
    url: 'order/open-orders',
    requestType: RequestTypes.GET,
  });
};
export const getOrderHistoryAPI = () => {
  return apiService.fetchData({
    data: {},
    url: 'order/full-history',
    requestType: RequestTypes.GET,
  });
};
export const getPaginatedOrderHistoryAPI = (params: any) => {
  return apiService.fetchData({
    data: params,
    url: 'order/full-history',
    requestType: RequestTypes.GET,
  });
};
export const getPaginatedTradeHistoryAPI = (params: any) => {
  return apiService.fetchData({
    data: params,
    url: 'trade/full-history',
    requestType: RequestTypes.GET,
  });
};

export const getFilteredOrderHistoryAPI = (data: OrderHistorySearchModel) => {
  return apiService.fetchData({
    data: data,
    url: 'order/full-history',
    requestType: RequestTypes.GET,
  });
};
export const getTradeHistoryAPI = (data: any) => {
  return apiService.fetchData({
    data: data,
    url: 'trade/full-history',
    requestType: RequestTypes.GET,
  });
};
export const getOrderHistoryDetailAPI = (data: { order_id: number }) => {
  return apiService.fetchData({
    data: data,
    url: 'order/detail',
    requestType: RequestTypes.GET,
  });
};
export const getCurrencyPairDetailsAPI = (data: {
  pair_currency_id: number;
}) => {
  return apiService.fetchData({
    data: data,
    url: 'user-balance/pair-balance',
    requestType: RequestTypes.GET,
  });
};
export const createNewOrderAPI = (data: NewOrderModel) => {
  return apiService.fetchData({
    data: data,
    url: 'order/create',
    requestType: RequestTypes.POST,
  });
};
export const createNewStopOrderAPI = (data: NewOrderModel) => {
  return apiService.fetchData({
    data: data,
    url: 'order/stop-order-create',
    requestType: RequestTypes.POST,
  });
};
export const cancelOrderAPI = (data: {
  order_id: number;
  mainType: string;
}) => {
  const sendingData = { order_id: data.order_id };
  return apiService.fetchData({
    data: sendingData,
    url: 'order/cancel',
    requestType: RequestTypes.POST,
  });
};
