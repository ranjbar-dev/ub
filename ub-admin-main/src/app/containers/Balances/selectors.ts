import { createSelector } from '@reduxjs/toolkit';
import { RootState } from 'types';

import { initialState } from './slice';

const selectDomain = (state: RootState) => state.balances || initialState;

export const selectBalances = createSelector(
  [selectDomain],
  balancesState => balancesState,
);

export const selectBalancesData = createSelector(
  [selectDomain],
  balancesState => balancesState.balances,
);

export const selectTransferModalBalancesData = createSelector(
  [selectDomain],
  balancesState => balancesState.transferModalBalances,
);

export const selectBalancesHistoryData = createSelector(
  [selectDomain],
  balancesState => balancesState.balancesHistory,
);
