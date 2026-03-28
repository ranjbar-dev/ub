import { PayloadAction } from '@reduxjs/toolkit';
import { createSlice } from 'utils/@reduxjs/toolkit';
import { ContainerState } from './types';
import { DepositSaveData } from '../Billing/types';

// The initial state of the Deposits container
export const initialState: ContainerState = { depositsData: null, isLoading: false, error: null };

const depositsSlice = createSlice({
  name: 'deposits',
  initialState,
  reducers: {
    GetDepositsAction(state, action: PayloadAction<Record<string, unknown>>) {
      state.isLoading = true;
    },
    setDepositsData(state, action: PayloadAction<Record<string, unknown>>) {
      state.depositsData = action.payload;
      state.isLoading = false;
    },
    UpdateDepositsAction(state, action: PayloadAction<DepositSaveData>) {},
  },
});

export const {
  actions: DepositsActions,
  reducer: DepositsReducer,
  name: sliceKey,
} = depositsSlice;
