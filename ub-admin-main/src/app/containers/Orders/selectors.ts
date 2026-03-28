import { createSelector } from '@reduxjs/toolkit';
import { RootState } from 'types';

import { initialState } from './slice';

const selectDomain = (state: RootState) => state.orders || initialState;

export const selectOrders = createSelector(
  [selectDomain],
  ordersState => ordersState,
);

export const selectOpenOrdersData = createSelector(
  [selectDomain],
  ordersState => ordersState.openOrdersData,
);

export const selectOrderHistoryData = createSelector(
  [selectDomain],
  ordersState => ordersState.orderHistoryData,
);

export const selectTradeHistoryData = createSelector(
  [selectDomain],
  ordersState => ordersState.tradeHistoryData,
);
