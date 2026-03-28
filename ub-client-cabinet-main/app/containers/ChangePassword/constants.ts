

/*
 *
 * ChangePassword constants
 *
 */

enum ActionTypes {
  DEFAULT_ACTION = 'app/ChangePassword/DEFAULT_ACTION',
  CHANGE_PASSWORD_ACTION = 'app/ChangePassword/CHANGE_PASSWORD_ACTION',
  IS_CHANGING_PASSWORD_ACTION = 'app/ChangePassword/IS_CHANGING_PASSWORD_ACTION',
}
export type ChangePasswordData = {
  oldPassword?: string;
  newPassword?: string;
  confirmNewPassword?: string;
};

export default ActionTypes;
