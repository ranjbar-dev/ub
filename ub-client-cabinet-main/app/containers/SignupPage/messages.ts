/*
 * SignupPage Messages
 *
 * This contains all the text for the SignupPage container.
 */

import { defineMessages } from 'react-intl';

export const scope = 'containers.SignupPage';

export default defineMessages({
  header: {
    id: `${scope}.header`,
    defaultMessage: 'This is the SignupPage container!',
  },
  Singuptocontinue: {
    id: `${scope}.Singuptocontinue`,
    defaultMessage: 'ET.Singuptocontinue',
  },
  CreateAccount: {
    id: `${scope}.CreateAccount`,
    defaultMessage: 'ET.CreateAccount',
  },
  email: {
    id: `${scope}.email`,
    defaultMessage: 'ET.email',
  },
  password: {
    id: `${scope}.password`,
    defaultMessage: 'ET.password',
  },
  confirmNewPassword: {
    id: `${scope}.confirmNewPassword`,
    defaultMessage: 'ET.confirmNewPassword',
  },
  haveAcount: {
    id: `${scope}.haveAcount`,
    defaultMessage: 'ET.haveAcount',
  },
  login: {
    id: `${scope}.login`,
    defaultMessage: 'ET.login',
  },
  go_back_home: {
    id: `${scope}.go_back_home`,
    defaultMessage: 'ET.go_back_home',
  },
  Youraccounthasbeencreated: {
    id: `${scope}.Youraccounthasbeencreated`,
    defaultMessage: 'ET.Youraccounthasbeencreated',
  },
  Pleasecheckyourinboxtoconfirmyouremailaccount: {
    id: `${scope}.Pleasecheckyourinboxtoconfirmyouremailaccount`,
    defaultMessage: 'ET.Pleasecheckyourinboxtoconfirmyouremailaccount',
  },
  GoToLoginPage: {
    id: `${scope}.GoToLoginPage`,
    defaultMessage: 'ET.GoToLoginPage',
  },
});
