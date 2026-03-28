/*
 *
 * PhoneVerificationPage constants
 *
 */

enum ActionTypes {
  DEFAULT_ACTION = 'app/PhoneVerificationPage/DEFAULT_ACTION',
  GET_COUNTRIES_ACTION = 'app/PhoneVerificationPage/GET_COUNTRIES_ACTION',
  GET_SMS_ACTION = 'app/PhoneVerificationPage/GET_SMS_ACTION',
  RESEND_SMS_ACTION = 'app/PhoneVerificationPage/RESEND_SMS_ACTION',
  SET_COUNTRIES_ACTION = 'app/PhoneVerificationPage/SET_COUNTRIES_ACTION',
  SET_COUNTRIES_LOADING = 'app/PhoneVerificationPage/SET_COUNTRIES_LOADING',
  SET_STEP_ACTION = 'app/PhoneVerificationPage/SET_STEP_ACTION',
  SET_PHONE_NUMBER_ACTION = 'app/PhoneVerificationPage/SET_PHONE_NUMBER',
  SET_IS_SENDING_SMS = 'app/PhoneVerificationPage/SET_IS_SENDING_SMS',
  VERIFY_CODE = 'app/PhoneVerificationPage/VERIFY_CODE',
}
export enum PhoneVerificationSteps {
  ENTER_PHONE_NUMBER = 0,
  ENTER_CODE = 1,
  ENTER_ACOUNT_PASSWORD = 2,
  GOOGLE_2FA_STEP = 3,
  DONE_STEP = 4,
}
interface Country {
  id: number;
  name: string;
  fullName: string;
  code: string;
  image: string;
}
export default ActionTypes;
export { Country };
