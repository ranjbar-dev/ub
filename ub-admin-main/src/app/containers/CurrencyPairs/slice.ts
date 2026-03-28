import { PayloadAction } from '@reduxjs/toolkit';
import { createSlice } from 'utils/@reduxjs/toolkit';

import { ContainerState } from './types';

// The initial state of the CurrencyPairs container
export const initialState: ContainerState = { currencyPairsData: null, isLoading: false, error: null };

const currencyPairsSlice = createSlice({
  name: 'currencyPairs',
  initialState,
  reducers: {
    GetCurrencyPairsAction(state, action: PayloadAction<Record<string, unknown>>) {},
    UpdateCurrencyPairAction(state, action: PayloadAction<Record<string, unknown>>) {},
    setCurrencyPairsData(state, action: PayloadAction<Record<string, unknown>>) {
      state.currencyPairsData = action.payload;
    },
  },
});

export const {
  actions: CurrencyPairsActions,
  reducer: CurrencyPairsReducer,
  name: sliceKey,
} = currencyPairsSlice;
