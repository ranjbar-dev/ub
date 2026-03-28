import { PayloadAction } from '@reduxjs/toolkit';
import { createSlice } from 'utils/@reduxjs/toolkit';

import { ContainerState } from './types';

// The initial state of the LoginHistory container
export const initialState: ContainerState = { loginHistoryData: null, isLoading: false, error: null };

const loginHistorySlice = createSlice({
  name: 'loginHistory',
  initialState,
  reducers: {
    GetLoginHistory(state, action: PayloadAction<Record<string, unknown>>) {},
    setLoginHistoryData(state, action: PayloadAction<Record<string, unknown>>) {
      state.loginHistoryData = action.payload;
    },
  },
});

export const {
  actions: LoginHistoryActions,
  reducer: LoginHistoryReducer,
  name: sliceKey,
} = loginHistorySlice;
