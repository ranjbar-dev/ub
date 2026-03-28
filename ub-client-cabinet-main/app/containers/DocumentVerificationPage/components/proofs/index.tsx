import React, { memo, useState, useEffect, useRef } from 'react';
import translate from '../../messages';
import styled from 'styles/styled-components';
import identityIcon from 'images/identityIcon.svg';
import { FormattedMessage } from 'react-intl';
import { NarrowInputs } from 'global-styles';
import { Select, MenuItem } from '@material-ui/core';
import IconAndTitle from './components/iconAndTitle';
import {
  IUserProfileMetaData,
  UserProfileData,
  IUserProfileImage,
} from 'containers/DocumentVerificationPage/types';
import ExpandMore from 'images/themedIcons/expandMore';
import {
  Subscriber,
  MessageNames,
  MessageService,
} from 'services/message_service';
import { UpperFirstLetters } from 'utils/formatters';
import Descriptions from './components/descriptions';
import ResidenceDescriptions from './components/residence_descriptions';
import FileUploader from '../common/fileUploader';
import { UploadState } from 'containers/DocumentVerificationPage/constants';
import SubmitButton from '../common/submitButton';
import residenceIcon from 'images/residenceIcon.svg';

const ProofOfIdentity = (props: {
  userProfileData: UserProfileData;
  documentType: string;
}) => {
  const { documentType } = props;
  const [SelectedIdentityTypeName, setSelectedIdentityTypeName] = useState('');
  const [IsSelectSubTypeLocked, setIsSelectSubTypeLocked] = useState(false);
  const TemporaryUploadedFiles: any = useRef({});
  const [SubTypes, setSubTypes]: [
    { name?: string; hasBack?: boolean; nameToShow?: string }[],
    any,
  ] = useState([{ name: 'passport', nameToShow: 'Passport', hasBack: false }]);
  const SubtypesRef: any = useRef({});
  const { current: subtypesObj } = SubtypesRef;
  const disableBackRef = useRef(false);
  const frontImageId = useRef(-1);

  const [DocumentImages, setDocumentImages]: [
    {
      frontImage?: Partial<IUserProfileImage>;
      backImage?: Partial<IUserProfileImage>;
    },
    any,
  ] = useState({});

  const handleDocumentTypeChange = (e: any) => {
    MessageService.send({
      name: MessageNames.TOGGLE_SEND_IMAGE_BUTTON,
      payload: false,
      additional: documentType,
      value: 'ALL',
    });

    setSelectedIdentityTypeName(e.target.value);
    setTimeout(() => {
      MessageService.send({
        name: MessageNames.SET_UPLOADER_STATE,
        payload: {
          uploaderId: documentType + e.target.value + true,
          setTo: UploadState.READY,
        },
      });
      MessageService.send({
        name: MessageNames.SET_UPLOADER_STATE,
        payload: {
          uploaderId: documentType + e.target.value + false,
          setTo: UploadState.READY,
        },
      });
    }, 0);
  };
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SET_DOCUMENT_IMAGES) {
        const {
          userProfileImages,
          userProfileImagesMetaData,
        }: {
          userProfileImages: IUserProfileImage[];
          userProfileImagesMetaData: IUserProfileMetaData;
        } = message.payload;
        ///////////set subtypes
        const [subtypes] = userProfileImagesMetaData.types
          .filter(item => {
            return item.name === documentType;
          })
          .map(item => {
            return item.subTypes.map(item => {
              return {
                name: item.name,
                hasBack: item.hasBack ?? false,
                nameToShow: UpperFirstLetters(item.name.replace('_', ' ')),
              };
            });
          });
        subtypes.forEach(item => {
          SubtypesRef.current[item.name] = item;
        });

        ///////////
        /////////// set document images
        const images: {
          frontImage?: Partial<IUserProfileImage>;
          backImage?: Partial<IUserProfileImage>;
        } = {};

        userProfileImages
          .filter(item => item.type === documentType)
          .forEach(item => {
            if (item.isBack === true) {
              images['backImage'] = item;
            } else {
              images['frontImage'] = item;
            }
          });
        if (images.backImage) {
          if (images.backImage.mainImageId !== images.frontImage?.imageId) {
            //@ts-ignore
            images.backImage = null;
          }
        }
        if (images.backImage || images.frontImage) {
          setIsSelectSubTypeLocked(true);
        }
        if (
          !images.frontImage ||
          images.frontImage.status === UploadState.REJECTED
        ) {
          disableBackRef.current = true;
        } else {
          //@ts-ignore
          frontImageId.current = images.frontImage.id;
          disableBackRef.current = false;
        }
        setDocumentImages(images);
        setSubTypes(subtypes);
        ////select default subtype if image uploaded before
        if (images.frontImage) {
          setSelectedIdentityTypeName(images.frontImage.subType ?? '');
        } else {
          setSelectedIdentityTypeName(subtypes[0].name ?? '');
        }
        //////////////
      }
      if (
        message.name === MessageNames.UNLOCK_SUBTYPE_SELECT &&
        message.payload.type === documentType
      ) {
        setIsSelectSubTypeLocked(false);
      }

      if (message.name === MessageNames.PROFILE_FILE_LOADED) {
        TemporaryUploadedFiles.current[message.payload.type] = message.payload;
      }
      if (
        message.name === MessageNames.SET_UPLOADED_IMAGE &&
        message.payload.disable === documentType
      ) {
        setIsSelectSubTypeLocked(true);
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);

  return (
    <Wrapper>
      {documentType === 'identity' ? (
        <IconAndTitle
          icon={identityIcon}
          title={
            <FormattedMessage {...translate.PleaseUploadProofOfIdentity} />
          }
        />
      ) : (
        <IconAndTitle
          icon={residenceIcon}
          title={
            <FormattedMessage {...translate.PleaseUploadProofOfResidence} />
          }
        />
      )}
      <div className='selectDocument'>
        <Select
          fullWidth
          disabled={IsSelectSubTypeLocked}
          IconComponent={ExpandMore}
          MenuProps={{
            getContentAnchorEl: null,
            anchorOrigin: {
              vertical: 'bottom',
              horizontal: 'left',
            },
          }}
          variant='outlined'
          value={SelectedIdentityTypeName}
          onChange={handleDocumentTypeChange}
        >
          {SubTypes.map((item, index: number) => {
            return (
              <MenuItem key={item.name} value={item.name}>
                {item.nameToShow}
              </MenuItem>
            );
          })}
        </Select>
      </div>
      <div className='uploadWrapper'>
        {!SelectedIdentityTypeName && (
          <div>please select a document type to start uploading</div>
        )}
        {SelectedIdentityTypeName && (
          <>
            <FileUploader
              isBack={false}
              uploaderId={
                documentType +
                subtypesObj[SelectedIdentityTypeName].name +
                false
              }
              type={documentType}
              //@ts-ignore
              initialState={
                DocumentImages.frontImage
                  ? DocumentImages.frontImage.status
                  : UploadState.READY
              }
              //@ts-ignore
              uploadedImage={DocumentImages.frontImage}
              subtype={subtypesObj[SelectedIdentityTypeName]}
            />
            {subtypesObj[SelectedIdentityTypeName]['hasBack'] === true && (
              <FileUploader
                isBack={true}
                uploaderId={
                  documentType +
                  subtypesObj[SelectedIdentityTypeName].name +
                  true
                }
                frontId={frontImageId.current}
                type={documentType}
                //@ts-ignore
                uploadedImage={DocumentImages.backImage}
                //@ts-ignore
                initialState={
                  DocumentImages.backImage
                    ? DocumentImages.backImage.status
                    : disableBackRef.current === false
                    ? UploadState.READY
                    : UploadState.READY
                }
                subtype={subtypesObj[SelectedIdentityTypeName]}
              />
            )}
            {/* <FileUploader isBack={true} initialState={UploadState.READY} /> */}
          </>
        )}
      </div>

      {documentType === 'identity' ? (
        <Descriptions />
      ) : (
        <ResidenceDescriptions />
      )}
      <div className='submitButton'>
        {subtypesObj[SelectedIdentityTypeName] && (
          <SubmitButton
            type={documentType}
            subtype={subtypesObj[SelectedIdentityTypeName].name}
            DocumentImages={DocumentImages}
          />
        )}

        {/* {useMemo(
          () => (
            <SubmitButton
              type={documentType}
              subtype={subtypesObj[SelectedIdentityTypeName]?.name}
              DocumentImages={DocumentImages}
            />
          ),
          [subtypesObj[SelectedIdentityTypeName]],
        )} */}
      </div>
    </Wrapper>
  );
};
export default memo(ProofOfIdentity);
const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  height: 100%;
  align-items: center;

  ${NarrowInputs}

  .iconWrapper {
    flex: 13;
    display: flex;
    flex-direction: column;
    justify-content: flex-end;
    min-height: 95px;
  }
  .title {
    flex: 5;
    display: flex;
    flex-direction: column;
    justify-content: center;
    span {
      color: var(--blackText);
    }
  }
  .selectDocument {
    flex: 5;
    width: calc(100% - 358px);
  }

  .acceptableDescs {
    background: orange;
    flex: 15;
    min-height: 170px;
  }
  .submitButton {
    flex: 10;
    display: flex;
    flex-direction: column;
    justify-content: center;
  }
  .uploadWrapper {
    flex: 40;
    display: flex;
    width: calc(100% - 96px);
  }
`;
