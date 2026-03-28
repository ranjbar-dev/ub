/*
 *
 * ChangeUserInfoPage actions
 *
 */

import { action } from 'typesafe-actions';
import { UserProfileData } from './types';

import ActionTypes from './constants';

export const defaultAction = () => action(ActionTypes.DEFAULT_ACTION);
export const getUserProfileAction = () => action(ActionTypes.GET_USER_PROFILE);
export const setUserProfileAction = (payload: UserProfileData) =>
  action(ActionTypes.SET_USER_PROFILE, payload);
export const setIsLoadingUserProfileDataAction = (payload: boolean) =>
  action(ActionTypes.SET_IS_LOADING, payload);

export const updateUserProfileDataAction = (payload: any) =>
  action(ActionTypes.UPDATE_USER_DATA, payload);
