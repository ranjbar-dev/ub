/*
 * BreadCrumb Messages
 *
 * This contains all the text for the BreadCrumb component.
 */

import { defineMessages } from 'react-intl';

export const scope = 'pageNames';

export default defineMessages({
  home: {
    id: `${scope}.Home`,
    defaultMessage: 'ET_HOME',
  },
  acountAndSecurity: {
    id: `${scope}.AcountAndSecurity`,
    defaultMessage: 'ET_AcountAndSecurity',
  },
  changePassword: {
    id: `${scope}.ChangePassword`,
    defaultMessage: 'ET_ChangePassword',
  },
  DisableGoogleAuthenticator: {
    id: `${scope}.DisableGoogleAuthenticator`,
    defaultMessage: 'ET.DisableGoogleAuthenticator',
  },
  phoneVerification: {
    id: `${scope}.phoneVerification`,
    defaultMessage: 'ET_phoneVerification',
  },
  AddressManagement: {
    id: `${scope}.AddressManagement`,
    defaultMessage: 'ET_AddressManagement',
  },
  DocumentVerification: {
    id: `${scope}.DocumentVerification`,
    defaultMessage: 'ET.DocumentVerification',
  },
  ChangeInfo: {
    id: `${scope}.ChangeInfo`,
    defaultMessage: 'ET.ChangeInfo',
  },
  TwoFA: {
    id: `${scope}.TwoFA`,
    defaultMessage: 'ET.TwoFA',
  },
});
