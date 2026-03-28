import { PayloadAction, createSelector } from '@reduxjs/toolkit';
import { createSlice } from 'utils/@reduxjs/toolkit';
import { RootState } from 'types';
import { LocalStorageKeys } from 'services/constants';

export interface GlobalState {
  loggedIn: boolean;
}

// The initial state of the global slice
export const initialState: GlobalState = {
  loggedIn: false,
};

const globalSlice = createSlice({
  name: 'global',
  initialState,
  reducers: {
    setIsLoggedIn(state, action: PayloadAction<boolean>) {
      state.loggedIn = action.payload;
      if (action.payload === false) {
        Object.values(LocalStorageKeys).forEach((key) => {
          localStorage.removeItem(key);
        });
      }
    },
  },
});

export const {
  actions: globalActions,
  reducer: globalReducer,
  name: sliceKey,
} = globalSlice;

const selectGlobalDomain = (state: RootState) => state.global || initialState;

export const selectLoggedIn = createSelector(
  [selectGlobalDomain],
  (globalState) => globalState.loggedIn,
);
