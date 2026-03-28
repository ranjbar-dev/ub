import { createSelector } from '@reduxjs/toolkit';
import { RootState } from 'types';

import { initialState } from './slice';

const selectDomain = (state: RootState) => state.currencyPairs || initialState;

export const selectCurrencyPairs = createSelector(
  [selectDomain],
  currencyPairsState => currencyPairsState,
);

export const selectCurrencyPairsData = createSelector(
  [selectDomain],
  currencyPairsState => currencyPairsState.currencyPairsData,
);
