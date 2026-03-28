import { PayloadAction } from '@reduxjs/toolkit';
import { createSlice } from 'utils/@reduxjs/toolkit';

import { ContainerState } from './types';

// The initial state of the FilledOrders container
export const initialState: ContainerState = { filledOrdersData: null, isLoading: false, error: null };

const filledOrdersSlice = createSlice({
  name: 'filledOrders',
  initialState,
  reducers: {
    GetFilledOrdersAction(state, action: PayloadAction<Record<string, unknown>>) {
      state.isLoading = true;
    },
    setFilledOrdersData(state, action: PayloadAction<Record<string, unknown>>) {
      state.filledOrdersData = action.payload;
      state.isLoading = false;
    },
  },
});

export const {
  actions: FilledOrdersActions,
  reducer: FilledOrdersReducer,
  name: sliceKey,
} = filledOrdersSlice;
