import { PayloadAction } from '@reduxjs/toolkit';
import { createSlice } from 'utils/@reduxjs/toolkit';

import { ContainerState } from './types';

// The initial state of the Admins container
export const initialState: ContainerState = {
  adminsData: null,
  isLoading: false,
  error: null,
};

const adminsSlice = createSlice({
  name: 'admins',
  initialState,
  reducers: {
    someAction(state, action: PayloadAction<Record<string, unknown>>) {},
  },
});

export const {
  actions: AdminsActions,
  reducer: AdminsReducer,
  name: sliceKey,
} = adminsSlice;
