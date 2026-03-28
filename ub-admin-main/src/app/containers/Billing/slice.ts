import { PayloadAction } from '@reduxjs/toolkit';
import { createSlice } from 'utils/@reduxjs/toolkit';

import { ContainerState, Payment, PaymentDetails, DepositSaveData } from './types';

// The initial state of the Billing container
export const initialState: ContainerState = {
  billingData: null,
  depositsData: null,
  withdrawalsData: null,
  allTransactionsData: null,
  selectedPaymentDetails: null,
  commissions: null,
  isLoading: false,
  error: null,
};

const billingSlice = createSlice({
  name: 'billing',
  initialState,
  reducers: {
    GetBillingGridDataAction(
      state,
      action: PayloadAction<{ user_id: number; type?: string }>,
    ) {},
    GetBillingDepositsDataAction(
      state,
      action: PayloadAction<{ user_id: number }>,
    ) {},
    GetBillingWithdrawalsDataAction(
      state,
      action: PayloadAction<{ user_id: number }>,
    ) {},
    GetBillingAllTransactionsDataAction(
      state,
      action: PayloadAction<{ user_id: number }>,
    ) {},
    GetBillingWithdrawDetailsAction(
      state,
      action: PayloadAction<{ id: number; user_id: number }>,
    ) {},
    AddPaymentCommentAction(
      state,
      action: PayloadAction<{ comment: string; payment_id: number }>,
    ) {},
    UpdateBillingWithdrawAction(state, action: PayloadAction<Record<string, unknown>>) {},
    UpdateDepositsAction(state, action: PayloadAction<DepositSaveData>) {},
    GetCommitionsAction(state, action: PayloadAction<Record<string, unknown>>) {},

    // Data reducers – populated by sagas via yield put()
    setBillingData(state, action: PayloadAction<Record<string, unknown>>) {
      state.billingData = action.payload;
    },
    setBillingDepositsData(state, action: PayloadAction<Record<string, unknown>>) {
      state.depositsData = action.payload;
    },
    setBillingWithdrawalsData(state, action: PayloadAction<Record<string, unknown>>) {
      state.withdrawalsData = action.payload;
    },
    setBillingAllTransactionsData(state, action: PayloadAction<Record<string, unknown>>) {
      state.allTransactionsData = action.payload;
    },
    setCommissionsData(state, action: PayloadAction<Record<string, unknown>>) {
      state.commissions = action.payload;
    },
    updateBillingWithdrawRow(state, action: PayloadAction<Record<string, unknown>>) {
      if (state.withdrawalsData) {
        const payments = state.withdrawalsData.payments as Record<string, unknown>[] | undefined;
        if (Array.isArray(payments)) {
          const index = payments.findIndex(p => p.id === action.payload.id);
          if (index !== -1) {
            payments[index] = { ...payments[index], ...action.payload };
          }
        }
      }
    },
    setWithdrawalItemDetails(
      state,
      action: PayloadAction<{ rowData: Payment; details: PaymentDetails }>,
    ) {
      state.selectedPaymentDetails = action.payload;
    },
  },
});

export const {
  actions: BillingActions,
  reducer: BillingReducer,
  name: sliceKey,
} = billingSlice;
