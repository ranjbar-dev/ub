import { createSelector } from '@reduxjs/toolkit';
import { RootState } from 'types';

import { initialState } from './slice';

const selectDomain = (state: RootState) => state.marketTicks || initialState;

export const selectMarketTicks = createSelector(
  [selectDomain],
  marketTicksState => marketTicksState,
);

export const selectMarketTicksData = createSelector(
  [selectDomain],
  marketTicksState => marketTicksState.marketTicksData,
);

export const selectSyncListData = createSelector(
  [selectDomain],
  marketTicksState => marketTicksState.syncListData,
);
