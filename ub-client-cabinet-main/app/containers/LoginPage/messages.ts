/*
 * LoginPage Messages
 *
 * This contains all the text for the LoginPage container.
 */

import { defineMessages } from 'react-intl';
import { GlobalTranslateScope } from 'containers/App/constants';

export const scope = 'containers.LoginPage';

export default defineMessages({
  header: {
    id: `${scope}.header`,
    defaultMessage: 'This is the LoginPage container!',
  },

  loginMessage: {
    id: `${scope}.loginMessage`,
    defaultMessage: 'ET_loginMessage',
  },
  password: {
    id: `${scope}.password`,
    defaultMessage: 'ET_password',
  },
  email: {
    id: `${scope}.email`,
    defaultMessage: 'ET_email',
  },
  login: {
    id: `${scope}.login`,
    defaultMessage: 'ET_login',
  },
  dont_have_acount: {
    id: `${scope}.dont_have_acount`,
    defaultMessage: 'ET_dont_have_acount',
  },
  forget_password: {
    id: `${scope}.forget_password`,
    defaultMessage: 'ET_forget_password',
  },

  signup: {
    id: `${scope}.signup`,
    defaultMessage: 'ET_signup',
  },
  go_back_home: {
    id: `${scope}.go_back_home`,
    defaultMessage: 'ET_signup',
  },
  forgetPassword: {
    id: `${scope}.forgetPassword`,
    defaultMessage: 'ET.forgetPassword',
  },
  invalidUsernameOrPassword: {
    id: `${scope}.invalidUsernameOrPassword`,
    defaultMessage: 'ET.invalidUsernameOrPassword',
  },
  Pleaseenteryouremailaddress: {
    id: `${scope}.Pleaseenteryouremailaddressandcheckmailboxtoresetyourpassword`,
    defaultMessage:
      'ET.Pleaseenteryouremailaddressandcheckmailboxtoresetyourpassword',
  },

  submit: {
    id: `${GlobalTranslateScope}.submit`,
    defaultMessage: 'ET.submit',
  },
});
