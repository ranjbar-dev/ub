import { createSelector } from '@reduxjs/toolkit';

import { RootState } from 'types';
import { initialState } from './slice';

const selectDomain = (state: RootState) => state.openOrders || initialState;

export const selectOpenOrders = createSelector(
  [selectDomain],
  openOrdersState => openOrdersState,
);

export const selectOpenOrdersData = createSelector(
  [selectDomain],
  openOrdersState => openOrdersState.openOrdersData,
);

export const selectOpenOrdersIsLoading = createSelector(
  [selectDomain],
  openOrdersState => openOrdersState.isLoading,
);