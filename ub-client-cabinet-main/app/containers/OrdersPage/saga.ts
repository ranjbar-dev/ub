import { takeLatest, put, call, takeEvery } from 'redux-saga/effects';
import ActionTypes, { LMS } from './constants';
import {
  setIsLoadingOpenOrdersAction,
  setOpenOrdersAction,
  setIsLoadingOrderHistoryAction,
  setOrderHistoryAction,
  setIsLoadingTradeHistoryAction,
  setTradeHistoryAction,
  setCurrenciesAction,
  getOpenOrdersAction,
  getCurrencyPairInfoAction,
} from './actions';

import { toast } from 'components/Customized/react-toastify';

import { StandardResponse } from 'services/constants';
import {
  getOpenOrdersAPI,
  getOrderHistoryAPI,
  getTradeHistoryAPI,
  getFilteredOrderHistoryAPI,
  getOrderHistoryDetailAPI,
  getCurrencyPairDetailsAPI,
  createNewOrderAPI,
  cancelOrderAPI,
  getPaginatedOrderHistoryAPI,
  getPaginatedTradeHistoryAPI,
} from 'services/orders_service';
import { MessageNames, MessageService } from 'services/message_service';
import { getCurrenciesAPI } from 'services/address_management_service';
import { OrderHistorySearchModel, FilterModel, Order } from './types';
import { NewOrderModel } from 'containers/TradePage/components/newOrder/types';
import { ToastMessages } from 'services/toastService';

function * getOpenOrders (action: {
  type: string;
  payload?: { silent: boolean };
}) {
  if (!action.payload) {
    yield put(setIsLoadingOpenOrdersAction(true));
    const currenciesResponse: StandardResponse = yield call(getCurrenciesAPI);
    yield put(setCurrenciesAction(currenciesResponse.data.currencies));
  }
  try {
    const response: StandardResponse = yield call(getOpenOrdersAPI);
    if (response.status === false) {
      yield put(setIsLoadingOpenOrdersAction(false));
      toast.error('error getting open orders');
      return;
    } else if (response.status === true) {
      yield put(setOpenOrdersAction(response.data));

      if (action.payload) {
        MessageService.send({
          name: MessageNames.SET_OPEN_ORDERS_DATA,
          payload: response.data,
        });
      }
    }
  } catch (error) {
    yield put(setIsLoadingOpenOrdersAction(false));
    toast.error('error while getting open orders');
  }
}

function * getOrderHistory (action: {
  type: string;
  payload?: { silent: boolean };
}) {
  if (!action.payload) {
    yield put(setIsLoadingOrderHistoryAction(true));
  }
  try {
    const response: StandardResponse = yield call(getOrderHistoryAPI);
    if (response.status === false) {
      if (!action.payload) {
        yield put(setIsLoadingOrderHistoryAction(false));
      }
      toast.error('error getting open orders');
      return;
    } else if (response.status === true) {
      const currectedData = response.data;
      for (let i = 0; i < currectedData.length; i++) {
        currectedData[i].createdAtToFilter = Number(
          currectedData[i]['createdAt'].split(' ')[0].replace(/-/g, ''),
        );
      }

      yield put(setOrderHistoryAction(currectedData));

      if (action.payload) {
        setTimeout(() => {
          MessageService.send({
            name: MessageNames.SET_ORDER_HISTORY_DATA,
            payload: currectedData,
          });
        }, 0);
      }
    }
  } catch (error) {
    if (!action.payload) {
      yield put(setIsLoadingOrderHistoryAction(false));
    }
    toast.error('error while getting order History');
  }
}

function * getPaginatedOrderHistory (action: { type: string; payload: any }) {
  MessageService.send({
    name: MessageNames.SET_PAGE_LOADING_WITH_ID,
    id: 'orderHistory',
    payload: true,
  });
  try {
    const response: StandardResponse = yield call(
      getPaginatedOrderHistoryAPI,
      action.payload,
    );
    if (response.status === false) {
      toast.error('error getting order history');
      return;
    } else if (response.status === true) {
      const currectedData = response.data;
      for (let i = 0; i < currectedData.length; i++) {
        currectedData[i].createdAtToFilter = Number(
          currectedData[i]['createdAt'].split(' ')[0].replace(/-/g, ''),
        );
      }
      setTimeout(() => {
        MessageService.send({
          name: MessageNames.SET_PAGINATED_ORDER_HISTORY_DATA,
          payload: currectedData,
        });
      }, 0);
    }
  } catch (error) {
    toast.error('error while getting order History');
  }
  MessageService.send({
    name: MessageNames.SET_PAGE_LOADING_WITH_ID,
    id: 'orderHistory',
    payload: false,
  });
}

function * getFilteredOrderHistory (action: {
  type: string;
  payload: OrderHistorySearchModel;
}) {
  MessageService.send({ name: MessageNames.SETLOADING, payload: true });
  MessageService.send({
    name: MessageNames.SET_PAGE_LOADING_WITH_ID,
    id: 'orderHistory',
    payload: true,
  });
  try {
    const response: StandardResponse = yield call(
      getFilteredOrderHistoryAPI,
      action.payload,
    );
    if (response.status === false) {
      toast.error('error getting open orders');
      MessageService.send({
        name: MessageNames.SETLOADING,
        payload: false,
      });
      return;
    }
    const currectedData = response.data;
    for (let i = 0; i < currectedData.length; i++) {
      currectedData[i].createdAtToFilter = Number(
        currectedData[i]['createdAt'].split(' ')[0].replace(/-/g, ''),
      );
    }
    yield put(setOrderHistoryAction(currectedData));
    setTimeout(() => {
      MessageService.send({
        name: MessageNames.SET_ORDER_HISTORY_DATA,
        payload: currectedData,
      });
      MessageService.send({ name: MessageNames.SETLOADING, payload: false });
    }, 0);
  } catch (error) {
    MessageService.send({ name: MessageNames.SETLOADING, payload: false });
    toast.error('error while getting order History');
  }
  MessageService.send({
    name: MessageNames.SET_PAGE_LOADING_WITH_ID,
    id: 'orderHistory',
    payload: false,
  });
}

function * getTradeHistory (action: { type: string; payload?: FilterModel }) {
  if (!action.payload) {
    yield put(setIsLoadingTradeHistoryAction(true));
  } else {
    if (!action.payload.silent) {
      MessageService.send({ name: MessageNames.SETLOADING, payload: true });
      MessageService.send({
        name: MessageNames.SET_PAGE_LOADING_WITH_ID,
        id: 'tradeHistory',
        payload: true,
      });
    }
  }

  try {
    const response: StandardResponse = yield call(getTradeHistoryAPI, {
      ...action.payload,
    });
    if (response.status === false) {
      if (!action.payload) {
        yield put(setIsLoadingTradeHistoryAction(false));
      } else {
        if (!action.payload.silent) {
          MessageService.send({
            name: MessageNames.SETLOADING,
            payload: false,
          });
        }
      }
      toast.error('error getting trade hostory orders');
      return;
    }
    const currectedData = response.data;
    for (let i = 0; i < currectedData.length; i++) {
      currectedData[i].createdAtToFilter = Number(
        currectedData[i]['createdAt'].split(' ')[0].replace(/-/g, ''),
      );
      currectedData[i].id = i;
    }
    yield put(setTradeHistoryAction(currectedData));
    if (!action.payload) {
      yield put(setIsLoadingTradeHistoryAction(false));
    } else {
      // setTimeout(() => {
      MessageService.send({
        name: MessageNames.SET_TRADE_HISTORY_DATA,
        payload: response.data,
      });

      MessageService.send({ name: MessageNames.SETLOADING, payload: false });
      // }, 0);
    }
  } catch (error) {
    if (!action.payload) {
      yield put(setIsLoadingTradeHistoryAction(false));
    } else {
      MessageService.send({ name: MessageNames.SETLOADING, payload: false });
    }
    toast.error('error while getting trade History');
  }
  MessageService.send({
    name: MessageNames.SET_PAGE_LOADING_WITH_ID,
    id: 'tradeHistory',
    payload: false,
  });
}
function * getPaginatedTradeHistory (action: { type: string; payload: any }) {
  MessageService.send({
    name: MessageNames.SET_PAGE_LOADING_WITH_ID,
    id: 'tradeHistory',
    payload: true,
  });
  try {
    const response: StandardResponse = yield call(
      getPaginatedTradeHistoryAPI,
      action.payload,
    );
    if (response.status === false) {
      toast.error('error getting trade history');
      return;
    } else if (response.status === true) {
      const currectedData = response.data;
      for (let i = 0; i < currectedData.length; i++) {
        currectedData[i].createdAtToFilter = Number(
          currectedData[i]['createdAt'].split(' ')[0].replace(/-/g, ''),
        );
      }
      setTimeout(() => {
        MessageService.send({
          name: MessageNames.SET_PAGINATED_TRADE_HISTORY_DATA,
          payload: currectedData,
        });
      }, 0);
    }
  } catch (error) {
    toast.error('error while getting trade History');
  }
  MessageService.send({
    name: MessageNames.SET_PAGE_LOADING_WITH_ID,
    id: 'tradeHistory',
    payload: false,
  });
}
function * getOrderDetail (action: {
  type: string;
  payload: { order_id: number; rowId: string };
}) {
  try {
    const response: StandardResponse = yield call(getOrderHistoryDetailAPI, {
      order_id: action.payload.order_id,
    });
    if (response.status === false) {
      toast.error('error getting details');
      return;
    }
    MessageService.send({
      name: MessageNames.SET_ORDER_DETAIL,
      //@ts-ignore
      orderId: action.payload.order_id,
      rowId: action.payload.rowId,
      payload: response.data[0],
    });
  } catch (error) {
    toast.error('error while getting data');
  }
}

function * getCurrencyPairInfo (action: {
  type: string;
  payload: { pair_currency_id: number };
}) {
  try {
    const response: StandardResponse = yield call(
      getCurrencyPairDetailsAPI,
      action.payload,
    );
    if (response.status === false) {
      toast.error('error getting currencyPair Details');
      return;
    } else if (response.status === true) {
      MessageService.send({
        name: MessageNames.SET_CURRENCY_PAIR_DETAILS,
        payload: response.data,
      });
    }
  } catch (error) {
    toast.error('error while getting  currencyPair Details');
  }
}
function * createNewOrder (action: { type: string; payload: NewOrderModel }) {
  const hasToReturn = false;

  MessageService.send({
    name: MessageNames.IS_LOADING_BUY_SELL,
    payload: true,
  });
  try {
    const response: StandardResponse = yield call(createNewOrderAPI, {
      amount: action.payload.amount,
      exchange_type:
        action.payload.exchange_type === LMS.Market
          ? action.payload.exchange_type
          : LMS.Limit,
      pair_currency_id: action.payload.pair_currency_id,
      price: action.payload.price,
      type: action.payload.type,
      stop_point_price: action.payload.stop_point_price,
      user_agent_info: action.payload.user_agent_info,
    });
    if (response.status === false) {
      //toast.error('error creating new Order');

      if (response.data && response.data.length > 0) {
        ToastMessages(response.data);
      } else if (response.message && response.message.length > 0) {
        toast.error(response.message);
      }
      MessageService.send({
        name: MessageNames.IS_LOADING_BUY_SELL,
        payload: false,
      });
      return;
    }

    if (response.status === true) {
      yield put(getOpenOrdersAction({ silent: true }));
    }

    MessageService.send({
      name: MessageNames.IS_LOADING_BUY_SELL,
      payload: false,
    });
  } catch (error) {
    toast.error('error while creating new Order');

    MessageService.send({
      name: MessageNames.IS_LOADING_BUY_SELL,
      payload: false,
    });
  }
}

function * cancelOrder (action: { type: string; payload: Order }) {
  MessageService.send({
    name: MessageNames.IS_CANCELING_ORDER,
    payload: { id: action.payload.id, state: true },
  });
  try {
    const response: StandardResponse = yield call(cancelOrderAPI, {
      order_id: action.payload.id,
      mainType: action.payload.mainType,
    });
    if (response.status === true) {
      //  toast.success('order canceled');

      yield put(getOpenOrdersAction({ silent: true }));
      //  MessageService.send({
      //    name: MessageNames.IS_CANCELING_ORDER,
      //    payload: { id: action.payload.id, state: false },
      //  });
    } else {
      if (response.message) {
        toast.error(response.message);
      } else {
        toast.error('error canceling order');
      }
      MessageService.send({
        name: MessageNames.IS_CANCELING_ORDER,
        payload: { id: action.payload.id, state: false },
      });
    }
  } catch (err) {
    toast.error('error while canceling order');
    MessageService.send({
      name: MessageNames.IS_CANCELING_ORDER,
      payload: { id: action.payload.id, state: false },
    });
  }
}

export default function * ordersPageSaga () {
  yield takeLatest(ActionTypes.GET_OPEN_ORDERS_ACTION, getOpenOrders);
  yield takeLatest(ActionTypes.GET_ORDER_HISTORY_ACTION, getOrderHistory);
  yield takeLatest(
    ActionTypes.GET_PAGINATED_ORDER_HISTORY_ACTION,
    getPaginatedOrderHistory,
  );
  yield takeLatest(
    ActionTypes.GET_PAGINATED_TRADE_HISTORY_ACTION,
    getPaginatedTradeHistory,
  );
  yield takeLatest(
    ActionTypes.GET_FILTERED_ORDER_HISTORY_ACTION,
    getFilteredOrderHistory,
  );
  yield takeLatest(ActionTypes.GET_TRADE_HISTORY_ACTION, getTradeHistory);
  yield takeEvery(ActionTypes.GET_PAYMENT_DETAIL_ACTION, getOrderDetail);
  ////trade new order
  yield takeLatest(ActionTypes.GET_CURRENCY_PAIR_INFO, getCurrencyPairInfo);
  yield takeLatest(ActionTypes.CREATE_NEW_ORDER_ACTION, createNewOrder);
  yield takeEvery(ActionTypes.CANCEL_ORDER_ACTION, cancelOrder);
}
