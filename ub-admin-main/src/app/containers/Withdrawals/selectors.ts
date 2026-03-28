import { createSelector } from '@reduxjs/toolkit';

import { RootState } from 'types';
import { initialState } from './slice';

const selectDomain = (state: RootState) => state.withdrawals || initialState;

export const selectWithdrawals = createSelector(
  [selectDomain],
  withdrawalsState => withdrawalsState,
);

export const selectWithdrawalsData = createSelector(
  [selectDomain],
  withdrawalsState => withdrawalsState.withdrawalsData,
);

export const selectWithdrawalsIsLoading = createSelector(
  [selectDomain],
  withdrawalsState => withdrawalsState.isLoading,
);