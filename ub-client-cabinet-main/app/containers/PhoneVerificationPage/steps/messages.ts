/*
 * StepIndicator Messages
 *
 * This contains all the text for the StepIndicator component.
 */

import { GlobalTranslateScope } from 'containers/App/constants';
import { defineMessages } from 'react-intl';

export const scope = 'containers.PhoneVerificationPage.steps';

export default defineMessages({
  selectCountry: {
    id: `${scope}.selectCountry`,
    defaultMessage: 'ET_selectCountry',
  },
  getSMS: {
    id: `${scope}.getSMS`,
    defaultMessage: 'ET_getSMS',
  },
  weAreSendingCodeTo: {
    id: `${scope}.weAreSendingCodeTo`,
    defaultMessage: 'ET_weAreSendingCodeTo',
  },
  editPhoneNumber: {
    id: `${scope}.editPhoneNumber`,
    defaultMessage: 'ET_editPhoneNumber',
  },
  resend: {
    id: `${scope}.resend`,
    defaultMessage: 'ET_resend',
  },
  Dontreceivecode: {
    id: `${scope}.Dontreceivecode`,
    defaultMessage: 'ET_Dontreceivecode',
  },
  pleaseCheckYourPhoneAndEnter: {
    id: `${scope}.pleaseCheckYourPhoneAndEnter`,
    defaultMessage: 'ET_pleaseCheckYourPhoneAndEnter',
  },
  AuthenticationCode: {
    id: `${scope}.AuthenticationCode`,
    defaultMessage: 'ET_AuthenticationCode',
  },
  enterCodeHere: {
    id: `${scope}.enterCodeHere`,
    defaultMessage: 'ET_enterCodeHere',
  },
  Pleaseenteryouraccountpassword: {
    id: `${scope}.Pleaseenteryouraccountpassword`,
    defaultMessage: 'ET_Pleaseenteryouraccountpassword',
  },
  Forsecurityreasonswerecommendtoenableyour2Fa: {
    id: `${scope}.Forsecurityreasonswerecommendtoenableyour2Fa`,
    defaultMessage: 'ET_Forsecurityreasonswerecommendtoenableyour2Fa',
  },
  Goto2Fapage: {
    id: `${scope}.Goto2Fapage`,
    defaultMessage: 'ET_Goto2Fapage',
  },
  PleaseopenGoogleAuthenticatorappinyourphoneandenter: {
    id: `${scope}.PleaseopenGoogleAuthenticatorappinyourphoneandenter`,
    defaultMessage: 'ET_PleaseopenGoogleAuthenticatorappinyourphoneandenter',
  },
  g2FaCode: {
    id: `${scope}.g2FaCode`,
    defaultMessage: 'ET_g2FaCode',
  },
  SMSAuthenticator: {
    id: `${scope}.SMSAuthenticator`,
    defaultMessage: 'ET_SMSAuthenticator',
  },
  Enabled: {
    id: `${scope}.Enabled`,
    defaultMessage: 'ET_Enabled',
  },
  //////////////////////////////////////////
  cancel: {
    id: `${GlobalTranslateScope}.cancel`,
    defaultMessage: 'ET_cancel',
  },
  submit: {
    id: `${GlobalTranslateScope}.submit`,
    defaultMessage: 'ET_submit',
  },
  password: {
    id: `${GlobalTranslateScope}.password`,
    defaultMessage: 'ET_password',
  },
  Password: {
    id: `${GlobalTranslateScope}.Password`,
    defaultMessage: 'ET_password',
  },
  or: {
    id: `${GlobalTranslateScope}.or`,
    defaultMessage: 'ET_or',
  },
  goToDashboard: {
    id: `${GlobalTranslateScope}.goToDashboard`,
    defaultMessage: 'ET_goToDashboard',
  },
  backToDashboard: {
    id: `${GlobalTranslateScope}.backToDashboard`,
    defaultMessage: 'ET_backToDashboard',
  },

  enterPhoneNumber: {
    id: `${GlobalTranslateScope}.enterPhoneNumber`,
    defaultMessage: 'ET_enterPhoneNumber',
  },
});
