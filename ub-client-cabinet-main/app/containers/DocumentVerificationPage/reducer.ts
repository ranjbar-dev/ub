/*
 *
 * DocumentVerificationPage reducer
 *
 */

import ActionTypes from './constants';
import { ContainerState, ContainerActions } from './types';

export const initialState: ContainerState = {
  default: null,
  userProfileData: {},
  isLoadingData: true,
};

function documentVerificationPageReducer(
  state: ContainerState = initialState,
  action: ContainerActions,
): ContainerState {
  switch (action.type) {
    case ActionTypes.DEFAULT_ACTION:
      return state;
    case ActionTypes.SET_USER_PROFILE:
      return {
        ...state,
        userProfileData: action.payload,
        isLoadingData: false,
      };
    case ActionTypes.SET_IS_LOADING:
      return { ...state, isLoadingData: action.payload };
    case ActionTypes.SET_UPLOADED_FILE:
      const images = state.userProfileData.userProfileImages
        ? state.userProfileData.userProfileImages
        : [];
      let found = false;
      for (let i = 0; i < images.length; i++) {
        if (images[i].type === action.payload.type) {
          found = true;
          images[i].image = action.payload.image;
          images[i].id = action.payload.id;
        }
      }
      if (found === false) {
        images.unshift({});
        images[0].image = action.payload.image;
        images[0].type = action.payload.type;
        images[0].id = action.payload.id;
      }
      const profileData = state.userProfileData;
      return {
        ...state,
        userProfileData: { ...profileData, userProfileImages: images },
      };

    default:
      return state;
  }
}

export default documentVerificationPageReducer;
