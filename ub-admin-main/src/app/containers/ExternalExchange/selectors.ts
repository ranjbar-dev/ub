import { createSelector } from '@reduxjs/toolkit';
import { RootState } from 'types';

import { initialState } from './slice';

const selectDomain = (state: RootState) =>
  state.externalExchange || initialState;

export const selectExternalExchange = createSelector(
  [selectDomain],
  externalExchangeState => externalExchangeState,
);

export const selectExternalExchangeData = createSelector(
  [selectDomain],
  externalExchangeState => externalExchangeState.externalExchangeData,
);
