import { createSelector } from '@reduxjs/toolkit';
import { RootState } from 'types';

import { initialState } from './slice';

const selectDomain = (state: RootState) => state.billing || initialState;

export const selectBilling = createSelector(
  [selectDomain],
  billingState => billingState,
);

export const selectBillingData = createSelector(
  [selectDomain],
  billingState => billingState.billingData,
);

export const selectBillingDepositsData = createSelector(
  [selectDomain],
  billingState => billingState.depositsData,
);

export const selectBillingWithdrawalsData = createSelector(
  [selectDomain],
  billingState => billingState.withdrawalsData,
);

export const selectBillingAllTransactionsData = createSelector(
  [selectDomain],
  billingState => billingState.allTransactionsData,
);

export const selectBillingCommissionsData = createSelector(
  [selectDomain],
  billingState => billingState.commissions,
);

export const selectWithdrawalItemDetails = createSelector(
  [selectDomain],
  billingState => billingState.selectedPaymentDetails,
);
