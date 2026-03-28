import { PayloadAction } from '@reduxjs/toolkit';
import { createSlice } from 'utils/@reduxjs/toolkit';

import { ContainerState } from './types';

// The initial state of the ExternalExchange container
export const initialState: ContainerState = { externalExchangeData: null, isLoading: false, error: null };

const externalExchangeSlice = createSlice({
  name: 'externalExchange',
  initialState,
  reducers: {
    GetExternalExchange(state, action: PayloadAction<Record<string, unknown>>) {},
    setExternalExchangeData(state, action: PayloadAction<Record<string, unknown>>) {
      state.externalExchangeData = action.payload;
    },
  },
});

export const {
  actions: ExternalExchangeActions,
  reducer: ExternalExchangeReducer,
  name: sliceKey,
} = externalExchangeSlice;
