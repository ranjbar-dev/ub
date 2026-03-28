import { createSelector } from '@reduxjs/toolkit';
import { RootState } from 'types';

const selectDomain = (state: RootState) => state;

export const selectRouter = createSelector(
  [selectDomain],
  state => state.router!,
);
