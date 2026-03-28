import { createSelector } from '@reduxjs/toolkit';

import { RootState } from 'types';
import { initialState } from './slice';

const selectDomain = (state: RootState) => state.deposits || initialState;

export const selectDeposits = createSelector(
  [selectDomain],
  depositsState => depositsState,
);

export const selectDepositsData = createSelector(
  [selectDomain],
  depositsState => depositsState.depositsData,
);

export const selectDepositsIsLoading = createSelector(
  [selectDomain],
  depositsState => depositsState.isLoading,
);
