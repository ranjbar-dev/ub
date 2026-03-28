import React, { memo } from 'react';
import styled from 'styles/styled-components';
import { IUserProfileImage } from 'containers/DocumentVerificationPage/types';
import RejectedResidence from 'images/themedIcons/rejectedResidence';
import { FormattedMessage } from 'react-intl';
import RejectedIdentity from 'images/themedIcons/rejectedIdentity';
import translate from '../../messages';
import rejectEnvelopeIcon from 'images/rejectEnvelopeIcon.svg';
import redShieldIcon from 'images/redShieldIcon.svg';

import { Button } from '@material-ui/core';
import { Buttons } from 'containers/App/constants';
import { MessageService, MessageNames } from 'services/message_service';
import Anime from 'react-anime';

interface Props {
  uploadedImage: IUserProfileImage;
}

function Rejected (props: Props) {
  const { uploadedImage } = props;
  const showRejectReason = (reason: string) => {
    // confirmAlert({
    //   type: 'success',
    //   message: (

    //   ),
    //   buttons: [],
    // });
    MessageService.send({
      name: MessageNames.OPEN_ALERT,
      payload: (
        <div
          style={{
            display: 'flex',
            flexDirection: 'column',
            alignItems: 'flex-start',
            minWidth: '650px',
            minHeight: '240px',
            padding: '48px',
          }}
        >
          <div>
            <span>
              <img src={redShieldIcon} />
            </span>
            <span style={{ color: '#E64141' }}>
              {' '}
              <FormattedMessage {...translate.YourDocumentIsNotVerified} />
            </span>
          </div>
          <span
            style={{
              margin: '3px 0px',
              marginTop: '32px',
              color: 'var(--blackText)',
            }}
          >
            {reason}
          </span>
        </div>
      ),
    });
  };
  return (
    <Wrapper>
      <Anime
        className='animated'
        duration={600}
        delay={200}
        easing='easeInCirc'
        opacity={[0, 1]}
      >
        <div className='withIconWrapper rejectedWrapper'>
          {uploadedImage.type === 'address' ? (
            <RejectedResidence />
          ) : (
            <RejectedIdentity />
          )}
          <div className='message'>
            <FormattedMessage
              style={{ pointerEvents: 'none' }}
              {...translate.YourDocumentHasBeenRejected}
            />
          </div>
          <Button
            onClick={() => {
              showRejectReason(
                uploadedImage.rejectionReason
                  ? uploadedImage.rejectionReason
                  : '',
              );
            }}
            endIcon={<img src={rejectEnvelopeIcon} />}
            className={`button ${Buttons.RoundedRedButton}`}
          >
            <FormattedMessage {...translate.RejectReason} />
          </Button>
        </div>
      </Anime>
    </Wrapper>
  );
}

export default memo(Rejected);
const Wrapper = styled.div`
  height: 100%;
  .rejectedWrapper {
    height: 80%;
    margin-top: 55px;
    justify-content: center;
    text-align: center;
  }
  .message {
    span {
      font-size: 14px !important;
      padding: 0px 10px;
    }
  }
`;
