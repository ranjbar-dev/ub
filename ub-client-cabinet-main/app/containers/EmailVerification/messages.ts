/*
 * EmailAuthentication Messages
 *
 * This contains all the text for the EmailAuthentication container.
 */

import { defineMessages } from 'react-intl';

export const scope = 'containers.EmailVerification';

export default defineMessages({
  header: {
    id: `${scope}.header`,
    defaultMessage: 'This is the EmailAuthentication container!',
  },
  CreateAccount: {
    id: `${scope}.CreateAccount`,
    defaultMessage: 'ET.CreateAccount',
  },
  Youraccounthasbeenactivated: {
    id: `${scope}.Youraccounthasbeenactivated`,
    defaultMessage: 'ET.Youraccounthasbeenactivated',
  },
  GoToLoginPage: {
    id: `${scope}.GoToLoginPage`,
    defaultMessage: 'ET.GoToLoginPage',
  },
});
