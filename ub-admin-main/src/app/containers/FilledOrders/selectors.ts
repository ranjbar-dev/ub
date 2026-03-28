import { createSelector } from '@reduxjs/toolkit';
import { RootState } from 'types';

import { initialState } from './slice';

const selectDomain = (state: RootState) => state.filledOrders || initialState;

export const selectFilledOrders = createSelector(
  [selectDomain],
  filledOrdersState => filledOrdersState,
);

export const selectFilledOrdersData = createSelector(
  [selectDomain],
  filledOrdersState => filledOrdersState.filledOrdersData,
);

export const selectFilledOrdersIsLoading = createSelector(
  [selectDomain],
  filledOrdersState => filledOrdersState.isLoading,
);
