import { createSelector } from '@reduxjs/toolkit';
import { RootState } from 'types';

import { initialState } from './slice';

const selectDomain = (state: RootState) => state.externalOrders || initialState;

export const selectExternalOrders = createSelector(
  [selectDomain],
  externalOrdersState => externalOrdersState,
);

export const selectExternalOrdersData = createSelector(
  [selectDomain],
  state => state.externalOrdersData,
);

export const selectNetQueueData = createSelector(
  [selectDomain],
  state => state.netQueueData,
);

export const selectAllQueueData = createSelector(
  [selectDomain],
  state => state.allQueueData,
);

export const selectNewQueueDetailList = createSelector(
  [selectDomain],
  state => state.newQueueDetailList,
);
