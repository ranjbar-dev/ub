import { PayloadAction } from '@reduxjs/toolkit';
import { Permission } from 'app/containers/UserDetails/types';
import { createSlice } from 'utils/@reduxjs/toolkit';

import { ContainerState, UserImagesData } from './types';

// The initial state of the VerificationWindow container
export const initialState: ContainerState = {
  userImages: { userId: null, data: null },
  permissionsData: { userId: null, data: null },
};

const verificationWindowSlice = createSlice({
  name: 'verificationWindow',
  initialState,
  reducers: {
    GetUserImagesAction(state, action: PayloadAction<Record<string, unknown>>) {},
    GetPermissionsAction(state, action: PayloadAction<Record<string, unknown>>) {},
    UpdatePermissionsAction(state, action: PayloadAction<Record<string, unknown>>) {},
    UpdateProfileImageStatusAction(
      state,
      action: PayloadAction<{
        user_id: number;
        id: number;
        confirmation_status: string;
        type: string;
        newType?: string;
        loadingButtonId: string;
        rejection_reason?: string;
        id_card_code?: string;
      }>,
    ) {},
    setUserImages(
      state,
      action: PayloadAction<{ userId: number; data: UserImagesData }>,
    ) {
      state.userImages = action.payload;
    },
    setPermissionsData(
      state,
      action: PayloadAction<{ userId: number; data: Permission[] }>,
    ) {
      state.permissionsData = action.payload;
    },
  },
});

export const {
  actions: VerificationWindowActions,
  reducer: VerificationWindowReducer,
  name: sliceKey,
} = verificationWindowSlice;
