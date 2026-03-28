import { createSelector } from '@reduxjs/toolkit';
import { RootState } from 'types';

import { initialState } from './slice';

const selectDomain = (state: RootState) => state.financeMethods || initialState;

export const selectFinanceMethods = createSelector(
  [selectDomain],
  financeMethodsState => financeMethodsState,
);

export const selectFinanceMethodsData = createSelector(
  [selectDomain],
  financeMethodsState => financeMethodsState.financeMethodsData,
);
