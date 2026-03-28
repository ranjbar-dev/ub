/*
 * ChangeUserInfoPage Messages
 *
 * This contains all the text for the ChangeUserInfoPage container.
 */

import { defineMessages } from 'react-intl';
import { GlobalTranslateScope } from 'containers/App/constants';

export const scope = 'containers.ChangeUserInfoPage';

export default defineMessages({
  TheinformationyoufillinmustbeconsistentwiththeinformationinyourIDdocuments: {
    id: `${scope}.TheinformationyoufillinmustbeconsistentwiththeinformationinyourIDdocuments`,
    defaultMessage:
      'ET.TheinformationyoufillinmustbeconsistentwiththeinformationinyourIDdocuments',
  },
  Basicinfo: {
    id: `${scope}.Basicinfo`,
    defaultMessage: 'ET.Basicinfo',
  },
  Residentialaddress: {
    id: `${scope}.Residentialaddress`,
    defaultMessage: 'ET.Residentialaddress',
  },
  BeginVerification: {
    id: `${scope}.BeginVerification`,
    defaultMessage: 'ET.BeginVerification',
  },
  FirstName: {
    id: `${scope}.FirstName`,
    defaultMessage: 'ET.FirstName',
  },
  LastName: {
    id: `${scope}.LastName`,
    defaultMessage: 'ET.LastName',
  },
  male: {
    id: `${scope}.male`,
    defaultMessage: 'ET.male',
  },
  female: {
    id: `${scope}.female`,
    defaultMessage: 'ET.female',
  },
  postalCode: {
    id: `${scope}.postalCode`,
    defaultMessage: 'ET.postalCode',
  },
  city: {
    id: `${scope}.city`,
    defaultMessage: 'ET.city',
  },
  address: {
    id: `${scope}.address`,
    defaultMessage: 'ET.address',
  },
  cancel: {
    id: `${GlobalTranslateScope}.cancel`,
    defaultMessage: 'ET.cancel',
  },
  Required: {
    id: `${GlobalTranslateScope}.Required`,
    defaultMessage: 'ET.Required',
  },
  next: {
    id: `${GlobalTranslateScope}.next`,
    defaultMessage: 'ET.next',
  },
});
