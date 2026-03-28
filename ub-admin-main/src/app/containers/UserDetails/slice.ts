import { PayloadAction } from '@reduxjs/toolkit';
import { createSlice } from 'utils/@reduxjs/toolkit';

import { ContainerState } from './types';

// The initial state of the UserDetails container
export const initialState: ContainerState = {};

const userDetailsSlice = createSlice({
  name: 'userDetails',
  initialState,
  reducers: {
    GetWalletsAction(state, action: PayloadAction<{ id: number }>) {},
    GetWhiteAddressesAction(state, action: PayloadAction<{ id: number }>) {},
    GetPermissionsAction(state, action: PayloadAction<{ id: number }>) {},
    UpdateUserDataAction(state, action: PayloadAction<Record<string, unknown>>) {},
  },
});

export const {
  actions: UserDetailsActions,
  reducer: UserDetailsReducer,
  name: sliceKey,
} = userDetailsSlice;
