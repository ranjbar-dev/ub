import { createSelector } from '@reduxjs/toolkit';
import { RootState } from 'types';

import { initialState } from './slice';

const selectDomain = (state: RootState) => state.userAccounts || initialState;
const selectAppState = (state: RootState) => state;

export const selectUserAccounts = createSelector(
  [selectDomain],
  userAccountsState => userAccountsState.userAccountsData,
);
export const selectUserAccountsData = createSelector(
  [selectDomain],
  userAccountsState => userAccountsState.userAccountsData,
);
export const selectIsLoading = createSelector(
  [selectDomain],
  userAccountsState => userAccountsState.isLoading,
);
export const selectRouter = createSelector(
  [selectAppState],
  state => state.router,
);
