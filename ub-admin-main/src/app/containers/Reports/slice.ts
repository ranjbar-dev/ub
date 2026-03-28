import { PayloadAction } from '@reduxjs/toolkit';
import { createSlice } from 'utils/@reduxjs/toolkit';

import { ContainerState, Report } from './types';

// The initial state of the Reports container
export const initialState: ContainerState = { adminReports: null, withdrawalComments: null, isLoading: false, error: null };

const reportsSlice = createSlice({
  name: 'reports',
  initialState,
  reducers: {
    GetAdminReportsAction(state, action: PayloadAction<{ id: number }>) {
      state.adminReports = null;
    },
    GetWithdrawalCommentsAction(
      state,
      action: PayloadAction<Record<string, unknown>>,
    ) {},
    AddAdmiCommentAction(state, action: PayloadAction<Record<string, unknown>>) {},
    DeleteAdminCommentAction(state, action: PayloadAction<Record<string, unknown>>) {},
    EditAdminCommentAction(state, action: PayloadAction<Record<string, unknown>>) {},
    setAdminReportsData(state, action: PayloadAction<{ comments: Report[] }>) {
      state.adminReports = action.payload.comments;
    },
    setWithdrawalComments(state, action: PayloadAction<Record<string, unknown>>) {
      state.withdrawalComments = action.payload;
    },
  },
});

export const {
  actions: ReportsActions,
  reducer: ReportsReducer,
  name: sliceKey,
} = reportsSlice;
