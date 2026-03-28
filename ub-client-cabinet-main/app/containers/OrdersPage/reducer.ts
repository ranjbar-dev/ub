/*
 *
 * OrdersPage reducer
 *
 */

import ActionTypes from './constants';
import { ContainerState, ContainerActions } from './types';

export const initialState: ContainerState = {
  default: null,

  openOrders: [],
  orderHistory: [],
  tradeHistory: [],

  currencies: [],
  isLoadingOpenOrders: true,
  isLoadingOrderHistory: true,
  isLoadingTradeHistory: false,
};

function ordersPageReducer(
  state: ContainerState = initialState,
  action: ContainerActions,
): ContainerState {
  switch (action.type) {
    case ActionTypes.DEFAULT_ACTION:
      return state;

    case ActionTypes.SET_OPEN_ORDERS_ACTION:
      return {
        ...state,
        openOrders: action.payload,
        isLoadingOpenOrders: false,
      };
    case ActionTypes.SET_ORDER_HISTORY_ACTION:
      return {
        ...state,
        orderHistory: action.payload,
        isLoadingOrderHistory: false,
      };
    //case ActionTypes.SET_PAGINATED_ORDER_HISTORY_ACTION:
    //  return {
    //    ...state,
    //    orderHistory: [state.orderHistory, ...action.payload],
    //    isLoadingOrderHistory: false,
    //  };
    case ActionTypes.SET_TRADE_HISTORY_ACTION:
      return {
        ...state,
        tradeHistory: action.payload,
        isLoadingTradeHistory: false,
      };
    ////set loadings
    case ActionTypes.SET_IS_LOADING_OPEN_ORDERS_ACTION:
      return { ...state, isLoadingOpenOrders: action.payload };
    case ActionTypes.SET_IS_LOADING_ORDER_HISTORY_ACTION:
      return { ...state, isLoadingOrderHistory: action.payload };
    case ActionTypes.SET_IS_LOADING_TRADE_HISTORY_ACTION:
      return { ...state, isLoadingTradeHistory: action.payload };
    case ActionTypes.ADD_NEW_ORDER_ACTION:
      return {
        ...state,
        orderHistory:
          state.orderHistory.length > 0
            ? [action.payload, ...state.orderHistory]
            : [],
      };
    /////
    case ActionTypes.SET_CURRENCIES:
      return { ...state, currencies: action.payload };

    default:
      return state;
  }
}

export default ordersPageReducer;
