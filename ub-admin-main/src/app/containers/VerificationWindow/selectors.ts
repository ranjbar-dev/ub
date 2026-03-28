import { createSelector } from '@reduxjs/toolkit';
import { RootState } from 'types';

import { initialState } from './slice';

const selectDomain = (state: RootState) =>
  state.verificationWindow || initialState;

export const selectVerificationWindow = createSelector(
  [selectDomain],
  verificationWindowState => verificationWindowState,
);

export const selectUserImages = createSelector(
  [selectDomain],
  state => state.userImages,
);

export const selectPermissionsData = createSelector(
  [selectDomain],
  state => state.permissionsData,
);
