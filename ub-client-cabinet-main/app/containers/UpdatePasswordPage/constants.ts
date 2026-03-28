/*
 *
 * UpdatePasswordPage constants
 *
 */

enum ActionTypes {
  DEFAULT_ACTION = 'app/UpdatePasswordPage/DEFAULT_ACTION',
  VERIFY_CODE_ACTION = 'app/UpdatePasswordPage/VERIFY_CODE_ACTION',
  RESET_PASSWORD_ACTION = 'app/UpdatePasswordPage/RESET_PASSWORD_ACTION',
}
export enum UpdatePasswordPages {
  ResetPage = 'ResetPage',
  Error = 'error',
  UpdatedPage = 'updatedPage',
}
export default ActionTypes;
