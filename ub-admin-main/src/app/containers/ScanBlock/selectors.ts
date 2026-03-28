import { createSelector } from '@reduxjs/toolkit';
import { RootState } from 'types';

import { initialState } from './slice';

const selectDomain = (state: RootState) => state.scanBlock || initialState;

export const selectScanBlock = createSelector(
  [selectDomain],
  scanBlockState => scanBlockState,
);
