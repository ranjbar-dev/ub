import { PayloadAction } from '@reduxjs/toolkit';
import { createSlice } from 'utils/@reduxjs/toolkit';

import { ContainerState } from './types';

// The initial state of the Orders container
export const initialState: ContainerState = {
  openOrdersData: null,
  orderHistoryData: null,
  tradeHistoryData: null,
  orderHistory: null,
  tradeHistory: null,
  isLoading: false,
  error: null,
};

const ordersSlice = createSlice({
  name: 'orders',
  initialState,
  reducers: {
    GetOpenOrdersAction(state, action: PayloadAction<Record<string, unknown>>) {},
    GetOrderHistoryAction(state, action: PayloadAction<Record<string, unknown>>) {},
    GetTradeHistoryAction(state, action: PayloadAction<Record<string, unknown>>) {},
    setOpenOrdersData(state, action: PayloadAction<Record<string, unknown>>) {
      state.openOrdersData = action.payload;
    },
    setOrderHistoryData(state, action: PayloadAction<Record<string, unknown>>) {
      state.orderHistoryData = action.payload;
    },
    setTradeHistoryData(state, action: PayloadAction<Record<string, unknown>>) {
      state.tradeHistoryData = action.payload;
    },
  },
});

export const {
  actions: OrdersActions,
  reducer: OrdersReducer,
  name: sliceKey,
} = ordersSlice;
