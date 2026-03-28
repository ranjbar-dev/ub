/*
 * UpdatePasswordPage Messages
 *
 * This contains all the text for the UpdatePasswordPage container.
 */

import { defineMessages } from 'react-intl';
import { GlobalTranslateScope } from 'containers/App/constants';

export const scope = 'containers.UpdatePasswordPage';

export default defineMessages({
  ResetPassword: {
    id: `${scope}.ResetPassword`,
    defaultMessage: 'ET.ResetPassword',
  },
  NewPassword: {
    id: `${scope}.NewPassword`,
    defaultMessage: 'ET.NewPassword',
  },
  ConfirmNewPassword: {
    id: `${scope}.ConfirmNewPassword`,
    defaultMessage: 'ET.ConfirmNewPassword',
  },
  GoToLoginPage: {
    id: `${scope}.GoToLoginPage`,
    defaultMessage: 'ET.GoToLoginPage',
  },
  Yourpasswordhasbeenchanged: {
    id: `${scope}.Yourpasswordhasbeenchanged`,
    defaultMessage: 'ET.Yourpasswordhasbeenchanged',
  },
  cancel: {
    id: `${GlobalTranslateScope}.cancel`,
    defaultMessage: 'ET.cancel',
  },
  GoToHome: {
    id: `${GlobalTranslateScope}.GoToHome`,
    defaultMessage: 'ET.GoToHome',
  },
});
