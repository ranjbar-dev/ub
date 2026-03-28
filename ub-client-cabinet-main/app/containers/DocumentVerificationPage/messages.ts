/*
 * DocumentVerificationPage Messages
 *
 * This contains all the text for the DocumentVerificationPage container.
 */

import { defineMessages } from 'react-intl';
import { GlobalTranslateScope } from 'containers/App/constants';

export const scope = 'containers.DocumentVerificationPage';

export default defineMessages({
  PleaseUploadProofOfIdentity: {
    id: `${scope}.PleaseUploadProofOfIdentity`,
    defaultMessage: 'ET_PleaseUploadProofOfIdentity',
  },
  PleaseUploadProofOfResidence: {
    id: `${scope}.PleaseUploadProofOfResidence`,
    defaultMessage: 'ET.PleaseUploadProofOfResidence',
  },
  AcceptabledocumentsNationalIDDriverLicensePassport: {
    id: `${scope}.AcceptabledocumentsNationalIDDriverLicensePassport`,
    defaultMessage: 'ET_AcceptabledocumentsNationalIDDriverLicensePassport',
  },
  AcceptabledocumentsBankStatementUtilityBill: {
    id: `${scope}.AcceptabledocumentsBankStatementUtilityBill`,
    defaultMessage: 'ET.AcceptabledocumentsBankStatementUtilityBill',
  },
  Thedocumentmustbeissuedwithinthelastsixmonths: {
    id: `${scope}.Thedocumentmustbeissuedwithinthelastsixmonths`,
    defaultMessage: 'ET_Thedocumentmustbeissuedwithinthelastsixmonths',
  },
  Theresolutionoftheuploadeddocumentsmustbeatleast300dpi: {
    id: `${scope}.Theresolutionoftheuploadeddocumentsmustbeatleast300dpi`,
    defaultMessage: 'ET_Theresolutionoftheuploadeddocumentsmustbeatleast300dpi',
  },
  AcceptablefileformatsGIFJPEGJPGPNGBMPorPDF: {
    id: `${scope}.AcceptablefileformatsGIFJPEGJPGPNGBMPorPDF`,
    defaultMessage: 'ET_AcceptablefileformatsGIFJPEGJPGPNGBMPorPDF',
  },
  Electronicdocumentsarenotacceptable: {
    id: `${scope}.Electronicdocumentsarenotacceptable`,
    defaultMessage: 'ET_Electronicdocumentsarenotacceptable',
  },
  Dropfiletouploadorbrowse: {
    id: `${scope}.Dropfiletouploadorbrowse`,
    defaultMessage: 'ET.Dropfiletouploadorbrowse',
  },
  FrontSide: {
    id: `${scope}.FrontSide`,
    defaultMessage: 'ET.FrontSide',
  },
  Of: {
    id: `${scope}.Of`,
    defaultMessage: 'ET.Of',
  },
  BackSide: {
    id: `${scope}.BackSide`,
    defaultMessage: 'ET.BackSide',
  },
  your: {
    id: `${scope}.your`,
    defaultMessage: 'ET.Your',
  },
  browse: {
    id: `${scope}.browse`,
    defaultMessage: 'ET.browse',
  },
  YourDocumentHasBeenReceivedItWillBeReviewedSoon: {
    id: `${scope}.YourDocumentHasBeenReceivedItWillBeReviewedSoon`,
    defaultMessage: 'ET.YourDocumentHasBeenReceivedItWillBeReviewedSoon',
  },
  YourDocumentHasBeenVerified: {
    id: `${scope}.YourDocumentHasBeenVerified`,
    defaultMessage: 'ET.YourDocumentHasBeenVerified',
  },
  YourDocumentHasBeenRejected: {
    id: `${scope}.YourDocumentHasBeenRejected`,
    defaultMessage: 'ET.YourDocumentHasBeenRejected',
  },

  RejectReason: {
    id: `${scope}.RejectReason`,
    defaultMessage: 'ET.RejectReason',
  },
  YourDocumentIsNotVerified: {
    id: `${scope}.YourDocumentIsNotVerified`,
    defaultMessage: 'ET.YourDocumentIsNotVerified',
  },
  ////////////
  submit: {
    id: `${GlobalTranslateScope}.submit`,
    defaultMessage: 'ET.submit',
  },
  Uploaded: {
    id: `${GlobalTranslateScope}.Uploaded`,
    defaultMessage: 'ET.Uploaded',
  },
  tryAgain: {
    id: `${GlobalTranslateScope}.tryAgain`,
    defaultMessage: 'ET.tryAgain',
  },
});
