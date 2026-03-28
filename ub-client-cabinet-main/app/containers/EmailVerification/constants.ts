/*
 *
 * EmailAuthentication constants
 *
 */

enum ActionTypes {
  DEFAULT_ACTION = 'app/EmailAuthentication/DEFAULT_ACTION',
  ACOUNT_ACTIVATION_ACTION = 'app/EmailAuthentication/ACOUNT_ACTIVATION_ACTION',
}
export enum EmailVerificationPages {
  Loading = 'loading',
  Error = 'error',
  Verified = 'verified',
}
export default ActionTypes;
