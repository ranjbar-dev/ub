import React, { memo, useState, useEffect, useRef } from 'react';
import { Button } from '@material-ui/core';
import translate from '../../messages';
import { FormattedMessage } from 'react-intl';
import ReplayIcon from 'images/themedIcons/replayIcon';
import { IUserProfileImage } from 'containers/DocumentVerificationPage/types';
import {
  Subscriber,
  MessageNames,
  MessageService,
} from 'services/message_service';
import { useDispatch } from 'react-redux';
import { UploadState } from 'containers/DocumentVerificationPage/constants';
import Anime from 'react-anime';
import { uploadMultiFileAction } from 'containers/DocumentVerificationPage/actions';

interface SubmitButtonProps {
  type: string;
  subtype: string;
  DocumentImages: {
    frontImage?: Partial<IUserProfileImage>;
    backImage?: Partial<IUserProfileImage>;
  };
}
export interface IResetImagesMessage {
  name: MessageNames;
  payload: {
    uploaderId: string;
    newState: UploadState;
  };
}
const SubmitButton: React.FC<SubmitButtonProps> = (
  props: SubmitButtonProps,
) => {
  const dispatch = useDispatch();
  const { DocumentImages, type, subtype } = props;

  const [FrontUploaded, setFrontUploaded] = useState(false);
  const [CanSubmit, setCanSubmit] = useState(false);
  const [HasRejectedImage, setHasRejectedImage] = useState(false);
  const UploadedImages = useRef({});
  const frontId = useRef(-1);
  const backId = useRef(-1);
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (
        message.name === MessageNames.TOGGLE_SEND_IMAGE_BUTTON &&
        message.additional === type
      ) {
        if (message.value === 'front_image') {
          delete UploadedImages.current['front_image'];
          setCanSubmit(message.payload);
        } else if (message.value === 'back_image') {
          delete UploadedImages.current['back_image'];
          if (!UploadedImages.current['front_image']) {
            setCanSubmit(message.payload);
          }
        } else if (message.value === 'ALL') {
          delete UploadedImages.current['front_image'];
          delete UploadedImages.current['back_image'];
          setCanSubmit(message.payload);
        }
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (
        message.name === MessageNames.PROFILE_FILE_LOADED &&
        message.payload.type === type
      ) {
        if (message.payload.isBack !== true) {
          UploadedImages.current['front_image'] = message.payload;
        } else if (message.payload.isBack === true) {
          UploadedImages.current['back_image'] = message.payload;
        }
        if (
          message.payload.isBack !== true ||
          (DocumentImages.frontImage &&
            DocumentImages.frontImage.status !== UploadState.REJECTED) ||
          FrontUploaded === true ||
          UploadedImages.current['back_image']
        ) {
          setCanSubmit(true);
        }
      }
      if (message.name === MessageNames.SET_UPLOADED_IMAGE) {
        if (message.payload.uploaderId === type + subtype + false) {
          delete UploadedImages.current['front_image'];
          frontId.current = message.payload.imageId;
          setFrontUploaded(true);
        } else if (message.payload.uploaderId === type + subtype + true) {
          delete UploadedImages.current['back_image'];
          if (message.payload.imageId) {
            backId.current = message.payload.imageId;
          }
          setFrontUploaded(false);
        }
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, [DocumentImages, FrontUploaded]);
  useEffect(() => {
    if (
      (DocumentImages.frontImage &&
        DocumentImages.frontImage.status === UploadState.REJECTED) ||
      (DocumentImages.backImage &&
        DocumentImages.backImage.status === UploadState.REJECTED)
    ) {
      setHasRejectedImage(true);
    }
    if (
      DocumentImages.frontImage &&
      DocumentImages.frontImage.status !== UploadState.REJECTED
    ) {
      //@ts-ignore
      frontId.current = DocumentImages.frontImage.id;
    }
    if (
      DocumentImages.backImage &&
      DocumentImages.backImage.status !== UploadState.REJECTED
    ) {
      //@ts-ignore
      backId.current = DocumentImages.backImage.id;
    }

    return () => {};
  }, []);
  const reset = () => {
    if (
      DocumentImages.frontImage &&
      DocumentImages.frontImage.status === UploadState.REJECTED
    ) {
      const id =
        //@ts-ignore
        DocumentImages.frontImage.type +
        //@ts-ignore
        DocumentImages.frontImage.subType +
        false;
      const tosend: IResetImagesMessage = {
        name: MessageNames.RESET_IMAGES,
        payload: {
          uploaderId: id,
          newState: UploadState.READY,
        },
      };
      MessageService.send(tosend);
    }
    if (
      DocumentImages.backImage &&
      DocumentImages.backImage.status === UploadState.REJECTED
    ) {
      const id =
        //@ts-ignore
        DocumentImages.backImage.type +
        //@ts-ignore
        DocumentImages.backImage.subType +
        true;
      const tosend: IResetImagesMessage = {
        name: MessageNames.RESET_IMAGES,
        payload: {
          uploaderId: id,
          newState: UploadState.READY,
        },
      };
      MessageService.send(tosend);
    }
    if (
      (DocumentImages.frontImage &&
        DocumentImages.frontImage.status === UploadState.REJECTED &&
        DocumentImages.backImage &&
        DocumentImages.backImage.status === UploadState.REJECTED) ||
      (DocumentImages.frontImage &&
        DocumentImages.frontImage.status === UploadState.REJECTED &&
        !DocumentImages.backImage)
    ) {
      MessageService.send({
        name: MessageNames.UNLOCK_SUBTYPE_SELECT,
        payload: { type: type },
      });
    }
    setHasRejectedImage(false);
  };
  const handleSubmit = () => {
    //@ts-ignore
    //console.log(backId.current);
    dispatch(
      uploadMultiFileAction({
        frontImage: UploadedImages.current['front_image']?.file,
        backImage: UploadedImages.current['back_image']?.file,
        type,
        subtype,
        ...(frontId.current !== -1 && { front_image_id: frontId.current }),
        ...(backId.current !== -1 && { back_image_id: backId.current }),
      }),
    );
    setCanSubmit(false);
  };
  return (
    <>
      {HasRejectedImage === false ? (
        <Anime
          className='animated'
          duration={400}
          delay={600}
          easing='easeInOutElastic'
          scale={[0, 1]}
          opacity={[0, 1]}
        >
          {' '}
          <Button
            disabled={!CanSubmit}
            onClick={handleSubmit}
            variant='contained'
            color='primary'
          >
            <FormattedMessage {...translate.submit} />
          </Button>
        </Anime>
      ) : (
        <Button
          onClick={reset}
          variant='outlined'
          color='primary'
          endIcon={<ReplayIcon />}
        >
          <FormattedMessage {...translate.tryAgain} />
        </Button>
      )}
    </>
  );
};

export default memo(SubmitButton);
