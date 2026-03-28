/*
 * PhoneVerificationPage Messages
 *
 * This contains all the text for the PhoneVerificationPage container.
 */

import { defineMessages } from 'react-intl';
import { GlobalTranslateScope } from 'containers/App/constants';

export const scope = 'containers.PhoneVerificationPage';

export default defineMessages({
  header: {
    id: `${scope}.header`,
    defaultMessage: 'This is the PhoneVerificationPage container!',
  },
  step1: {
    id: `${scope}.step1`,
    defaultMessage: 'ET_step1',
  },
  step2: {
    id: `${scope}.step2`,
    defaultMessage: 'ET_step2',
  },
  step3: {
    id: `${scope}.step3`,
    defaultMessage: 'ET_step3',
  },
  enterPhoneNumber: {
    id: `${scope}.enterPhoneNumber`,
    defaultMessage: 'ET_enterPhoneNumber',
  },
  enterCode: {
    id: `${scope}.enterCode`,
    defaultMessage: 'ET_enterCode',
  },
  phoneNumberIsRequired: {
    id: `${scope}.phoneNumberIsRequired`,
    defaultMessage: 'ET.phoneNumberIsRequired',
  },
  phoneNumberIsNotValid: {
    id: `${scope}.phoneNumberIsNotValid`,
    defaultMessage: 'ET.phoneNumberIsNotValid',
  },
  enterAcountPassword: {
    id: `${scope}.enterAcountPassword`,
    defaultMessage: 'ET_enterAcountPassword',
  },
  selectCountry: {
    id: `${GlobalTranslateScope}.selectCountry`,
    defaultMessage: 'ET.selectCountry',
  },
  submit: {
    id: `${GlobalTranslateScope}.submit`,
    defaultMessage: 'ET.submit',
  },
});
