import { createSelector } from '@reduxjs/toolkit';
import { RootState } from 'types';

import { initialState } from './slice';

const selectDomain = (state: RootState) => state.liquidityOrders || initialState;

export const selectLiquidityOrders = createSelector(
  [selectDomain],
  liquidityOrdersState => liquidityOrdersState,
);

export const selectLiquidityOrdersData = createSelector(
  [selectDomain],
  liquidityOrdersState => liquidityOrdersState.liquidityOrders,
);
