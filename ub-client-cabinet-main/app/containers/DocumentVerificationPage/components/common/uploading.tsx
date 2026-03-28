import React, { useState, useEffect } from 'react';
import { FormattedMessage } from 'react-intl';
import translate from '../../messages';
import { Subscriber, MessageNames } from 'services/message_service';
import styled from 'styles/styled-components';
import { useDispatch } from 'react-redux';

interface Props {
  subtype: string;
  isBack: boolean;
  uploaderId: string;
}
export interface IUploadPercentMessage {
  name: MessageNames;
  payload: {
    uploaderId: string;
    percent: number;
  };
}
function Uploading (props: Props) {
  const { subtype, isBack, uploaderId } = props;
  const [Percent, setPercent] = useState(0);
  const dispatch = useDispatch();
  useEffect(() => {
    const Subscription = Subscriber.subscribe(
      (message: IUploadPercentMessage) => {
        if (
          message.name === MessageNames.UPLOAD_PERCENTAGE &&
          message.payload.uploaderId === uploaderId
        ) {
          setPercent(message.payload.percent);
        }
      },
    );
    return () => {
      Subscription.unsubscribe();
    };
  }, []);
  return (
    <Wrapper>
      <div className='uploadingWrapper'>
        <div className='per'>
          <span>
            <FormattedMessage {...translate.Uploaded} /> {Percent}
            {' %'}
          </span>
        </div>
        <div className='progg'>
          <div className='progPercent' style={{ width: `${Percent}%` }}></div>
        </div>
      </div>
    </Wrapper>
  );
}

export default Uploading;
const Wrapper = styled.div`
  position: relative;
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  .uploadingWrapper {
    span {
      color: var(--textBlue);
    }
    .progg {
      width: calc(100% + 65px);
      height: 6px;
      border: 1px solid #d8d8d8;
      position: relative;
      border-radius: 10px;
      margin-left: -31px;
      .progPercent {
        height: 6px;
        background: var(--textBlue);
        margin-top: -1px;
        border-radius: 10px;
        transition: width 0.3s;
      }
    }
  }
`;
