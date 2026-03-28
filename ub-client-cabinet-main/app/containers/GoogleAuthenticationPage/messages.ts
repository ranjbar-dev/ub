/*
 * GoogleAuthenticationPage Messages
 *
 * This contains all the text for the GoogleAuthenticationPage container.
 */

import { defineMessages } from 'react-intl';
import { GlobalTranslateScope } from 'containers/App/constants';

export const scope = 'containers.GoogleAuthenticationPage';

export default defineMessages({
  DOWNLOADANDINSTALL: {
    id: `${scope}.DOWNLOADANDINSTALL`,
    defaultMessage: 'ET.DOWNLOADANDINSTALL',
  },
  ENTERPROVIDEDKEYORScanQRCode: {
    id: `${scope}.ENTERPROVIDEDKEYORScanQRCode`,
    defaultMessage: 'ET.ENTERPROVIDEDKEYORScanQRCode',
  },
  GetAuthenticationCodeFromAppAndEnterHere: {
    id: `${scope}.GetAuthenticationCodeFromAppAndEnterHere`,
    defaultMessage: 'ET.GetAuthenticationCodeFromAppAndEnterHere',
  },
  Providedkey: {
    id: `${scope}.Providedkey`,
    defaultMessage: 'ET.Providedkey',
  },
  codeCopiedToClipboard: {
    id: `${scope}.codeCopiedToClipboard`,
    defaultMessage: 'ET.codeCopiedToClipboard',
  },
  Pleaseenteryouraccountpassword: {
    id: `${scope}.Pleaseenteryouraccountpassword`,
    defaultMessage: 'ET.Pleaseenteryouraccountpassword',
  },
  password: {
    id: `${scope}.password`,
    defaultMessage: 'ET.password',
  },
  done: {
    id: `${scope}.done`,
    defaultMessage: 'ET.done',
  },
  OpenTowFactorAuthenticationWithGoogle: {
    id: `${scope}.OpenTowFactorAuthenticationWithGoogle`,
    defaultMessage: 'ET.OpenTowFactorAuthenticationWithGoogle',
  },
  PleaseopenGoogleAuthenticatorappinyourphoneandenter2FaCode: {
    id: `${scope}.PleaseopenGoogleAuthenticatorappinyourphoneandenter2FaCode`,
    defaultMessage:
      'ET.PleaseopenGoogleAuthenticatorappinyourphoneandenter2FaCode',
  },
  g2faWraning: {
    id: `${scope}.g2faWraning`,
    defaultMessage: 'ET.g2faWraning',
  },

  ////////////////
  backToDashboard: {
    id: `${GlobalTranslateScope}.backToDashboard`,
    defaultMessage: 'ET.backToDashboard',
  },
  next: {
    id: `${GlobalTranslateScope}.next`,
    defaultMessage: 'ET.next',
  },
  submit: {
    id: `${GlobalTranslateScope}.submit`,
    defaultMessage: 'ET.submit',
  },
  cancel: {
    id: `${GlobalTranslateScope}.cancel`,
    defaultMessage: 'ET.cancel',
  },
  enabled: {
    id: `${GlobalTranslateScope}.enabled`,
    defaultMessage: 'ET.enabled',
  },
});
