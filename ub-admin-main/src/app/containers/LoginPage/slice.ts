import { PayloadAction } from '@reduxjs/toolkit';
import { createSlice } from 'utils/@reduxjs/toolkit';

import { ContainerState } from './types';

// The initial state of the LoginPage container
export const initialState: ContainerState = {
  isLoading: false,
};

const loginPageSlice = createSlice({
  name: 'loginPage',
  initialState,
  reducers: {
    setIsLoadingAction(state, action: PayloadAction<boolean>) {
      state.isLoading = action.payload;
    },
    LoginAction(
      state,
      action: PayloadAction<{ username: string; password: string }>,
    ) {},
  },
});

export const { actions, reducer, name: sliceKey } = loginPageSlice;
