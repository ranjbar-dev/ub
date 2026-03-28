import { PayloadAction } from '@reduxjs/toolkit';
import { createSlice } from 'utils/@reduxjs/toolkit';

import { ContainerState } from './types';

// The initial state of the HomePage container
export const initialState: ContainerState = { isLoading: false, error: null };

const homePageSlice = createSlice({
  name: 'homePage',
  initialState,
  reducers: {
    someAction(state, action: PayloadAction<Record<string, unknown>>) {},
    getUserByIdAction(state, action: PayloadAction<Record<string, unknown>>) {},
    getWithdrawalByIdAction(state, action: PayloadAction<Record<string, unknown>>) {},
  },
});

export const {
  actions: HomePageActions,
  reducer: HomePageReducer,
  name: sliceKey,
} = homePageSlice;
