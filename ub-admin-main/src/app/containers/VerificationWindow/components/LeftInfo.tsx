import { InitialUserDetails } from 'app/containers/UserAccounts/types';
import { translations } from 'locales/i18n';
import React, { memo, useCallback } from 'react';
import { useTranslation } from 'react-i18next';
import styled from 'styled-components/macro';
import './imageWrapperStyle.scss';
import EditDropDown from 'app/containers/UserDetails/components/EditDropDown';

import { IdentityTypes } from '../constants';

import { MessageService, MessageNames } from 'services/messageService';

import { ProfileImageData } from '../types';

import RawInput from 'app/components/RawInput/RawInput';
interface Props {
  data: InitialUserDetails;
  type: IdentityTypes;
  selectedImage: ProfileImageData;
  subTypes: { identity: { name: string }[]; address: { name: string }[] };
}

function LeftInfo(props: Props) {
  const { data, type, selectedImage, subTypes } = props;
  const { t } = useTranslation();

  const RowData = useCallback(
    (props: { title: string; value: React.ReactNode; isLast?: boolean }) => {
      return (
        <>
          <div className="rowDataContainer">
            <div className="title" style={{ minWidth: '118px' }}>
              {props.title} :{' '}
            </div>
            <div className="value">{props.value}</div>
          </div>
          {!props.isLast && <div className="hr"></div>}
        </>
      );
    },
    [type, selectedImage],
  );
  const handleDocumentIdChange = (e: string) => {
    MessageService.send({
      name: MessageNames.DATASEND,
      payload: {
        userId: data.id,
        documentId: e,
      },
    });
  };
  const handleTypeChange = (e: string) => {
    MessageService.send({
      name: MessageNames.DATASEND,
      payload: {
        userId: data.id,
        sub_type: e,
      },
    });
  };
  return (
    <Wrapper className="leftInfo">
      <RowData
        title={t(translations.CommonTitles.Name())}
        value={data.fullName}
      />
      <RowData
        title={t(translations.CommonTitles.Email())}
        value={data.email}
      />
      <RowData
        title={t(translations.CommonTitles.Gender())}
        value={data.gender}
      />
      <RowData
        title={t(translations.CommonTitles.Birthdate())}
        value={data.birthDate}
      />
      <RowData
        title={t(translations.CommonTitles.Mobile())}
        value={data.mobile}
      />
      <RowData
        title={t(translations.CommonTitles.RegistrationDate())}
        value={data.registrationDate}
      />
      <RowData
        title={t(translations.CommonTitles.RegisteredIP())}
        value={data.registeredIp}
      />
      <RowData
        title={t(translations.CommonTitles.Country())}
        value={data.country}
      />
      <RowData title={t(translations.CommonTitles.City())} value={data.city} />
      <RowData
        title={t(translations.CommonTitles.PostalCode())}
        value={data.postalCode}
      />
      <RowData
        title={t(translations.CommonTitles.Address())}
        value={data.address}
      />
      <RowData
        title={t(translations.CommonTitles.SystemID())}
        value={data.systemId}
      />
      {/*<RowData
        title={t(translations.CommonTitles.ReferralID())}
        value={data.referralId}
      />
      <RowData
        title={t(translations.CommonTitles.ReferKey())}
        value={data.referKey}
      />*/}
      <RowData
        title={t(translations.CommonTitles.DocumentId())}
        value={
          <RawInput
            className="documentIdInput"
            onChange={e => handleDocumentIdChange(e)}
            initialValue={(selectedImage && selectedImage.idCardCode) ?? ''}
          />
        }
      />
      {selectedImage && selectedImage.updatedAt && (
        <RowData
          title={t(translations.CommonTitles.DocumentType())}
          value={
            <EditDropDown
              className="wideDrop"
              onSelect={handleTypeChange}
              initialValue={{
                name: selectedImage.subType.toLowerCase().replace('_', ' '),
                id: selectedImage.subType.toLowerCase(),
              }}
              options={subTypes[type].map((item, index) => {
                return {
                  name: item.name.toLowerCase().replace('_', ' '),
                  id: item.name,
                };
              })}
            />
          }
          isLast={true}
        />
      )}
    </Wrapper>
  );
}

export default memo(LeftInfo);
const Wrapper = styled.div`
  .rowDataContainer {
    display: flex;
  }
  .hr {
    width: 370px;
    margin-left: -1px;
    height: 1px;
    background: #e8e8e8;
  }
`;
