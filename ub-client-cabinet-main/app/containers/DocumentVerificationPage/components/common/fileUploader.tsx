import React, { useState, useRef, useEffect, useCallback } from 'react';
import DropToUpload from 'react-drop-to-upload';
import { UploadState } from 'containers/DocumentVerificationPage/constants';
import translate from '../../messages';
import { FormattedMessage } from 'react-intl';
import styled from 'styles/styled-components';
import {
  Subscriber,
  MessageNames,
  MessageService,
} from 'services/message_service';
import { IUserProfileImage } from 'containers/DocumentVerificationPage/types';

import Ready from './ready';
import Hover from './hover';
import Preview from './preview';
import Uploading from './uploading';
import Processing from './processing';
import Rejected from './rejected';
import Verified from './verified';
import { UploadModel } from 'services/upload_service';
import { toast } from 'components/Customized/react-toastify';

export interface ISetUploaderStateMessage {
  name: MessageNames;
  payload: {
    uploaderId: string;
    setTo: UploadState;
  };
}
export default function FileUploader (props: {
  initialState: UploadState;
  isBack: boolean;
  type: string;
  frontId?: number;
  uploadedImage: Partial<IUserProfileImage>;
  uploaderId: string;
  subtype: { name: string; hasBack: boolean; nameToShow: string };
}) {
  const {
    initialState,
    isBack,
    subtype,
    type,
    frontId,
    uploaderId,
    uploadedImage,
  } = props;
  const [UploaderState, setUploaderState] = useState(UploadState.READY);
  const [PreviewFile, setPreviewFile] = useState(null);
  const [UploadedImage, setUploadedImage] = useState(uploadedImage);
  const input: any = useRef(null);
  const readyToUploadDescs = {
    fileSubtype: subtype?.nameToShow,
    boldTitle:
      isBack === true ? (
        <FormattedMessage {...translate.BackSide} />
      ) : (
        <FormattedMessage {...translate.FrontSide} />
      ),
  };
  const isFileValid = useCallback((fileType: string) => {
    let isValid = false;
    const type = fileType.split('/')[1];
    switch (type) {
      case 'jpeg':
        isValid = true;
        break;
      case 'jpg':
        isValid = true;
        break;
      case 'png':
        isValid = true;
        break;
      case 'bmp':
        isValid = true;
        break;
      case 'pdf':
        isValid = true;
        break;
      default:
        break;
    }
    return isValid;
  }, []);
  const handleBrowsButtonClick = useCallback(() => {
    input.current.value = null;
    input.current?.click();
  }, [input.current]);
  const handleDropHover = () => {
    if (UploaderState === UploadState.READY) {
      setUploaderState(UploadState.HOVER);
    }
  };

  const onFileLoaded = (file: File) => {
    if (file.size >= 4000000) {
      toast.warn('Maximum file size is 4MB');
      return;
    }
    if (!isFileValid(file.type)) {
      toast.warn('Only JPEG, JPG, PNG, BMP or PDF formats are valid');
      return;
    }
    //@ts-ignore
    setPreviewFile(file);
    setUploaderState(UploadState.PREVIEW);
    const toSend: UploadModel = {
      file,
      isBack,
      subtype: subtype.name,
      type,
      ...(isBack === true && { mainImageId: frontId }),
    };
    MessageService.send({
      name: MessageNames.PROFILE_FILE_LOADED,
      payload: toSend,
    });
  };
  const handleLeaveHover = () => {
    setUploaderState(initialState);
  };
  const handleDeleteClick = () => {
    input.current.value = null;
    setUploaderState(UploadState.READY);
    MessageService.send({
      name: MessageNames.TOGGLE_SEND_IMAGE_BUTTON,
      payload: false,
      additional: type,
      value: isBack !== true ? 'front_image' : 'back_image',
    });
  };
  useEffect(() => {
    setUploaderState(initialState);
    const Subscription = Subscriber.subscribe((message: any) => {
      if (
        message.name === MessageNames.SET_UPLOADER_STATE &&
        message.payload.uploaderId === uploaderId
      ) {
        setUploaderState(message.payload.setTo);
      }
      if (
        message.name === MessageNames.RESET_IMAGES &&
        message.payload.uploaderId === uploaderId
      ) {
        setUploaderState(message.payload.newState);
      }
      if (
        message.name === MessageNames.SET_UPLOADED_IMAGE &&
        message.payload.uploaderId === uploaderId
      ) {
        setUploadedImage(message.payload);
        setUploaderState(UploadState.PROCESSING);
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, [initialState, subtype, uploaderId]);
  const stateSelector = () => {
    switch (UploaderState) {
      case UploadState.READY:
        return (
          <Ready
            onBrowsClick={handleBrowsButtonClick}
            readyToUploadDescs={readyToUploadDescs}
          />
        );
      case UploadState.HOVER:
        return <Hover />;
      case UploadState.UPLOADING:
        return (
          <Uploading
            subtype={subtype.name}
            uploaderId={uploaderId}
            isBack={isBack}
          />
        );
      case UploadState.PREVIEW:
        return (
          <Preview
            //@ts-ignore
            previewFile={PreviewFile}
            deletePreview={handleDeleteClick}
          />
        );
      case UploadState.BLOCKED:
        return (
          <Ready
            onBrowsClick={handleBrowsButtonClick}
            readyToUploadDescs={readyToUploadDescs}
          />
        );
      case UploadState.PROCESSING:
        //@ts-ignore
        return <Processing image={UploadedImage.image ?? ''} />;
      case UploadState.REJECTED:
        //@ts-ignore
        return <Rejected uploadedImage={uploadedImage} />;
      case UploadState.CONFIRMED:
        //@ts-ignore
        return <Verified uploadedImage={uploadedImage} />;
      default:
        return (
          <Ready
            onBrowsClick={handleBrowsButtonClick}
            readyToUploadDescs={readyToUploadDescs}
          />
        );
    }
  };
  return (
    <Wrapper>
      <input
        type='file'
        ref={input}
        id={subtype.name + isBack}
        onChange={(e: any) => {
          e.persist();
          if (e.target.files[0]) onFileLoaded(e.target.files[0]);
        }}
        style={{ display: 'none' }}
      />

      <div className={`uploader front ${UploaderState}`}>
        <DropToUpload
          className='dropper'
          onOver={e => {
            handleDropHover();
          }}
          onLeave={e => {
            handleLeaveHover();
          }}
          onDrop={(e: any[]) => {
            onFileLoaded(e[0]);
          }}
        >
          {stateSelector()}
        </DropToUpload>
        {/* <input
          type="file"
          id={'props.uniqueId'}
          onChange={(e: any) => {
            e.persist();
            if (e.target.files[0]) onFileLoaded(e.target.files[0]);
          }}
          style={{ display: 'none' }}
        /> */}
      </div>
    </Wrapper>
  );
}

const Wrapper = styled.div`
  flex: 1;
  width: calc(100% - 96px);
  display: flex;
  justify-content: space-between;
  align-items: center;
  .contentWrapper {
    height: 100%;
    display: flex;
    flex: 10;
    flex-direction: column;
    place-items: center;
    justify-content: center;
    max-height: 322.5px;
  }
  .dropper {
    height: 100%;
    width: 100%;
  }
  .decripts {
    color: var(--textBlue);
  }
  .uploader {
    width: calc(100% - 6px);
    height: 90%;
    max-height: 35vh;
  }
  .blocked {
    filter: grayscale(1);
    pointer-events: none;
    opacity: 0.3;
  }
  .ready,
  .blocked,
  .hover,
  .uploading,
  .uploaded,
  .preview {
    border: 1px dashed var(--textBlue);
    border-radius: 10px;
    background: var(--lightBlue);
    transition: box-shadow 0.5s;
    box-shadow: 0px 0px 0px 0px #0000001f;
  }
  .confirmed {
    border: 1px dashed var(--greenText);
    border-radius: 10px;
    background: var(--lightBlue);
  }
  .rejected {
    border: 1px dashed var(--redText);
    border-radius: 10px;
    background: var(--lightBlue);
  }
  .hover,
  .uploading {
    box-shadow: 5px 3px 5px 2px #0000001f;
    .proccessing,
    .expandImage,
    .delete {
      display: none;
    }
  }
  .blocked {
    /* border: none; */
  }

  .uploadedWrapper {
    max-width: 100%;
    width: 100%;
    img {
      min-width: calc(100% + 60px);
      max-width: calc(100% + 60px);
      margin-left: -30px;
      max-height: 35vh;
    }
  }
  .withIconWrapper {
    display: flex;
    flex-direction: column;
    align-items: center;
    .button {
      margin-top: 24px;
    }
  }
  .confirmedWrapper {
    span {
      color: var(--greenText);
      font-size: 15px;
    }
  }
  .rejectedWrapper {
    span {
      color: var(--redText);
      font-size: 15px;
    }
  }
  .previewWrapper {
    img.previewImage {
      max-width: calc(100% + 80px);
      min-width: calc(100% + 80px);
      margin-left: -40px;
      border-radius: 7px;
      max-height: 245px;
    }
    .deletePreview {
      position: absolute;
      z-index: 1;
      left: -40px;
    }
  }
  .blueItalicBold {
    span {
      font-weight: bold;
      font-style: italic;
      color: var(--textBlue);
    }
  }
  .uploadReadyContainer {
    display: flex;
    flex-direction: column;
    height: 100%;
    align-items: center;
    text-align: center;
    padding: 0 40px;
    position: relative;
    overflow: hidden;
  }
  .uploadIcon {
  }
  .decripts {
    display: flex;
    align-items: center;
  }
  .browsButtonWrapper {
    .browsButton {
      max-height: 24px;
      padding: 0;
      span {
        font-size: 12px;
      }
    }
  }
  .spacer {
    flex: 4;
  }

  .delete {
    position: absolute;
    bottom: 0;
    left: 0;
    z-index: 2;
    .MuiIconButton-root {
      background: var(--white);
    }
  }
  .deletePreview {
    position: absolute;
    z-index: 1;
    left: 10px;
    bottom: 10px;
  }
  .expandImage {
    position: absolute;
    margin: auto;
    z-index: 5;
    top: calc(50% - 23px);
    .MuiIconButton-root {
      background: var(--darkTrans);
      &:hover {
        background: var(--darkTrans) !important;
        opacity: 0.9;
      }
    }
  }
`;
