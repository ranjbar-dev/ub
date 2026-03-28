import { PayloadAction } from '@reduxjs/toolkit';
import { createSlice } from 'utils/@reduxjs/toolkit';
import { ContainerState } from './types';

// The initial state of the OpenOrders container
export const initialState: ContainerState = { openOrdersData: null, isLoading: false, error: null };

const openOrdersSlice = createSlice({
  name: 'openOrders',
  initialState,
  reducers: {
    GetOpenOrdersAction(state, action: PayloadAction<Record<string, unknown>>) {
      state.isLoading = true;
    },
    setOpenOrdersData(state, action: PayloadAction<Record<string, unknown>>) {
      state.openOrdersData = action.payload;
      state.isLoading = false;
    },
    CancelOpenOrderAction(state, action: PayloadAction<Record<string, unknown>>) {},
    FullFillOpenOrderAction(state, action: PayloadAction<Record<string, unknown>>) {},
  },
});

export const {
  actions: OpenOrdersActions,
  reducer: OpenOrdersReducer,
  name: sliceKey,
} = openOrdersSlice;