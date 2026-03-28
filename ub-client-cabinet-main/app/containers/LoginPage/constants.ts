/*
 *
 * LoginPage constants
 *
 */

enum ActionTypes {
  DEFAULT_ACTION = 'app/LoginPage/DEFAULT_ACTION',
  LOGIN_ACTION = 'app/LoginPage/LOGIN_ACTION',
  IS_LOGGING_IN = 'app/LoginPage/IS_LOGGING_IN',
  FORGOT_PASSWORD_ACTION = 'app/LoginPage/FORGOT_PASSWORD_ACTION',
}
type LoginData = {
  username: string;
  password: string;
  remember?: boolean;
  '2fa_code'?: string;
  fromPopup?: boolean;
  recaptcha: string;
};
export default ActionTypes;
export { LoginData };
