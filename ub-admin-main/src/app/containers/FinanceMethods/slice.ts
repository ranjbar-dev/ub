import { PayloadAction } from '@reduxjs/toolkit';
import { createSlice } from 'utils/@reduxjs/toolkit';

import { ContainerState } from './types';

// The initial state of the FinanceMethods container
export const initialState: ContainerState = { financeMethodsData: null, isLoading: false, error: null };

const financeMethodsSlice = createSlice({
  name: 'financeMethods',
  initialState,
  reducers: {
    GetFinanceMethods(state, action: PayloadAction<Record<string, unknown>>) {},
    UpdateFinanceMethod(state, action: PayloadAction<Record<string, unknown>>) {},
    setFinanceMethodsData(state, action: PayloadAction<Record<string, unknown>>) {
      state.financeMethodsData = action.payload;
    },
  },
});

export const {
  actions: FinanceMethodsActions,
  reducer: FinanceMethodsReducer,
  name: sliceKey,
} = financeMethodsSlice;
