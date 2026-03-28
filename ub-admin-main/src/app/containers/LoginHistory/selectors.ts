import { createSelector } from '@reduxjs/toolkit';
import { RootState } from 'types';

import { initialState } from './slice';

const selectDomain = (state: RootState) => state.loginHistory || initialState;

export const selectLoginHistory = createSelector(
  [selectDomain],
  loginHistoryState => loginHistoryState,
);

export const selectLoginHistoryData = createSelector(
  [selectDomain],
  loginHistoryState => loginHistoryState.loginHistoryData,
);
