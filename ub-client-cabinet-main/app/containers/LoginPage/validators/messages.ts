/*
 * validator Messages
 *
 * This contains all the text for the LoginPage container.
 */

import { defineMessages } from 'react-intl';

export const scope = 'validationMessages';

export default defineMessages({
  emailIsRequired: {
    id: `${scope}.emailIsRequired`,
    defaultMessage: 'ET.emailIsRequired',
  },
  emailIsNotValid: {
    id: `${scope}.emailIsNotValid`,
    defaultMessage: 'ET.emailIsNotValid',
  },
  passwordIsRequired: {
    id: `${scope}.passwordIsRequired`,
    defaultMessage: 'ET.passwordIsRequired',
  },
  minimum8Character: {
    id: `${scope}.minimum8Character`,
    defaultMessage: 'ET.minimum8Character',
  },
});
