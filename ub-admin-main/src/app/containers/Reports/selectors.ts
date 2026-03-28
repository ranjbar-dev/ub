import { createSelector } from '@reduxjs/toolkit';
import { RootState } from 'types';

import { initialState } from './slice';

const selectDomain = (state: RootState) => state.reports || initialState;

export const selectReports = createSelector(
  [selectDomain],
  reportsState => reportsState,
);

export const selectAdminReportsData = createSelector(
  [selectDomain],
  reportsState => reportsState.adminReports,
);

export const selectWithdrawalComments = createSelector(
  [selectDomain],
  reportsState => reportsState.withdrawalComments,
);
