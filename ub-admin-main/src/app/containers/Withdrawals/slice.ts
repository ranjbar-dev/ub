import { PayloadAction } from '@reduxjs/toolkit';
import { createSlice } from 'utils/@reduxjs/toolkit';
import { ContainerState } from './types';

// The initial state of the Withdrawals container
export const initialState: ContainerState = { withdrawalsData: null, isLoading: false, error: null };

const withdrawalsSlice = createSlice({
  name: 'withdrawals',
  initialState,
  reducers: {
    GetWithdrawals(state, action: PayloadAction<Record<string, unknown>>) {
      state.isLoading = true;
    },
    setWithdrawalsData(state, action: PayloadAction<Record<string, unknown>>) {
      state.withdrawalsData = action.payload;
      state.isLoading = false;
    },
    GetWithdrawalDetailAction(state, action: PayloadAction<Record<string, unknown>>) {},
  },
});

export const {
  actions: WithdrawalsActions,
  reducer: WithdrawalsReducer,
  name: sliceKey,
} = withdrawalsSlice;