/*
 *
 * AcountPage constants
 *
 */

enum ActionTypes {
  DEFAULT_ACTION = 'app/AcountPage/DEFAULT_ACTION',
  SET_2FA_ENABLED_ACTION = 'app/AcountPage/SET_2FA_ENABLED_ACTION',
  IS_LOADING_ACTION = 'app/AcountPage/IS_LOADING_ACTION',
  SET_USER_DATA_ACTION = 'app/AcountPage/SET_USER_DATA_ACTION',
  LOGGED_IN_ACTION = 'App/LOGGED_IN_ACTION',
  GET_NEW_VERIFICATION_EMAIL_ACTION = 'App/AcountPage/GET_NEW_VERIFICATION_EMAIL_ACTION',
}
export enum SecurityLevel {
  LOW = 'low',
  MEDIUM = 'medium',
  High = 'high',
}
export enum KycStatus {
  INCOMPLETE = 'incomplete',
  PROCESSING = 'processing',
  CONFIRMED = 'confirmed',
  PARTIALLYCONFIRMED = 'partially_confirmed',
  REJECTED = 'rejected',
}
export enum KycLevel {
  NONE = 'incomplete',
  PROCESSING = 'processing',
  CONFIRMED = 'confirmed',
  PARTIALLYCONFIRMED = 'partially_confirmed',
  REJECTED = 'rejected',
}

export default ActionTypes;
