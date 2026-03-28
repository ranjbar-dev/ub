/*
 *
 * DocumentVerificationPage constants
 *
 */

enum ActionTypes {
  DEFAULT_ACTION = 'app/DocumentVerificationPage/DEFAULT_ACTION',
  GET_USER_PROFILE = 'app/DocumentVerificationPage/GET_USER_PROFILE',
  SET_USER_PROFILE = 'app/DocumentVerificationPage/SET_USER_PROFILE',
  SET_IS_LOADING = 'app/DocumentVerificationPage/SET_IS_LOADING',
  UPLOAD_FILE = 'app/DocumentVerificationPage/UPLOAD_FILE',
  UPLOAD_MULTI_FILE = 'app/DocumentVerificationPage/UPLOAD_MULTI_FILE',
  DELETE_FILE = 'app/DocumentVerificationPage/DELETE_FILE',
  DELETE_USER_IMAGE_ACTION = 'app/DocumentVerificationPage/DELETE_USER_IMAGE_ACTION',
  SET_UPLOADED_FILE = 'app/DocumentVerificationPage/SET_UPLOADED_FILE',
}
export enum UploadState {
  READY = 'ready',
  UPLOADING = 'uploading',
  ERROR = 'error',
  PREVIEW = 'preview',
  REJECTED = 'rejected',
  UPLOADED = 'uploaded',
  BLOCKED = 'blocked',
  CONFIRMED = 'confirmed',
  PROCESSING = 'processing',
  HOVER = 'hover',
}
export enum IdentityDocumentTypes {
  NATIONAL_ID = 'Identity Card',
  DRIVER_LICENSE = 'Driver License',
  PASSPORT = 'Passport',
}
export enum ResidenceDocumentTypes {
  BANK_STATEMENT = 'Bank Statement',
}
export enum ImageTypes {
  IDENTITY = 'identity',
  ADDRESS = 'address',
  IDENTITY_BACK = 'identityBack',
}
export enum ProfileImageStatus {
  PROCCESSING = 'processing',
}
export default ActionTypes;
