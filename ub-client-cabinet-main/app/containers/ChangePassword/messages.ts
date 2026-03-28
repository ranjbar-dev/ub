/*
 * ChangePassword Messages
 *
 * This contains all the text for the ChangePassword container.
 */

import { defineMessages } from 'react-intl';
import { GlobalTranslateScope } from 'containers/App/constants';

export const scope = 'containers.ChangePassword';
export const globalScope = 'app.globalTitles';
export default defineMessages({
  cancel: {
    id: `${globalScope}.cancel`,
    defaultMessage: 'ET_changePassword cancel',
  },

  header: {
    id: `${scope}.header`,
    defaultMessage: 'This is the ChangePassword container!',
  },
  changePassword: {
    id: `${scope}.changePassword`,
    defaultMessage: 'ET_changePassword',
  },
  oldPassword: {
    id: `${scope}.oldPassword`,
    defaultMessage: 'ET_oldPassword',
  },
  newPassword: {
    id: `${scope}.newPassword`,
    defaultMessage: 'ET_newPassword',
  },
  confirmNewPassword: {
    id: `${scope}.confirmNewPassword`,
    defaultMessage: 'ET_confirmNewPassword',
  },
  Yourpasswordhasbeenchanged: {
    id: `${scope}.Yourpasswordhasbeenchanged`,
    defaultMessage: 'ET.Yourpasswordhasbeenchanged',
  },
  backToDashboard: {
    id: `${GlobalTranslateScope}.backToDashboard`,
    defaultMessage: 'ET.backToDashboard',
  },
  submit: {
    id: `${GlobalTranslateScope}.submit`,
    defaultMessage: 'ET.submit',
  },
});
