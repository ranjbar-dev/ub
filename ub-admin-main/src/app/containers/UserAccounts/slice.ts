import { PayloadAction } from '@reduxjs/toolkit';
import { WindowTypes } from 'app/constants';
import { createSlice } from 'utils/@reduxjs/toolkit';
import { useInjectReducer, useInjectSaga } from "utils/redux-injectors";

import { userAccountsSaga } from "./saga";
import { ContainerState, User } from './types';


// The initial state of the UserAccounts container
export const initialState: ContainerState = {
  isLoading: true,
  userAccountsData: {
    count: 0,
    users: [],
  },
};

const userAccountsSlice = createSlice({
  name: 'userAccounts',
  initialState,
  reducers: {
    GetInitialUserAccountsAction(state, action: PayloadAction<Record<string, unknown>>) {
      state.isLoading = true;
    },

    setUserAccountsData(state, action: PayloadAction<{ count: number; users: User[] }>) {
      state.userAccountsData = action.payload;
      state.isLoading = false;
    },

    getInitialSingleUserDataAndOpenWindowAction(
      state,
      action: PayloadAction<{ id: number; windowType: WindowTypes }>,
    ) {
    },
  },
});

export const {
  actions: UserAccountsActions,
  reducer: UserAccountsReducer,
  name: sliceKey,
} = userAccountsSlice;

export const useUserAccountsSlice = () => {
  useInjectReducer({ key: sliceKey, reducer: UserAccountsReducer });
  useInjectSaga({ key: sliceKey, saga: userAccountsSaga });
  return { UserAccountsActions }
}