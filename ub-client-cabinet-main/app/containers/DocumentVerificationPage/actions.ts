/*
 *
 * DocumentVerificationPage actions
 *
 */

import { action } from 'typesafe-actions';
import { UserProfileData } from './types';

import ActionTypes from './constants';
import { UploadModel } from 'services/upload_service';

export const defaultAction = () => action(ActionTypes.DEFAULT_ACTION);
export const getUserProfileAction = (payload: { silent?: boolean }) =>
  action(ActionTypes.GET_USER_PROFILE, payload);
export const setUserProfileAction = (payload: UserProfileData) =>
  action(ActionTypes.SET_USER_PROFILE, payload);
export const setIsLoadingUserProfileDataAction = (payload: boolean) =>
  action(ActionTypes.SET_IS_LOADING, payload);

export const uploadFileAction = (payload: UploadModel) =>
  action(ActionTypes.UPLOAD_FILE, payload);
export const uploadMultiFileAction = (payload: {
  frontImage: File;
  backImage: File;
  type: string;
  subtype: string;
  front_image_id?: number | string;
}) => action(ActionTypes.UPLOAD_MULTI_FILE, payload);
export const uploadImageAction = (payload: UploadModel) =>
  action(ActionTypes.UPLOAD_FILE, payload);

export const deleteFileAction = (payload: { id: number; type: string }) =>
  action(ActionTypes.DELETE_FILE, payload);
export const deleteUserImageAction = (payload: { id: number }) =>
  action(ActionTypes.DELETE_USER_IMAGE_ACTION, payload);

export const setUploadedFileAction = (payload: {
  type: string;
  image: string;
  id: number;
  isBack?: boolean;
}) => action(ActionTypes.SET_UPLOADED_FILE, payload);
