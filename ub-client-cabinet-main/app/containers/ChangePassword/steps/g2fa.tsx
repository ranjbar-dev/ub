import React, { useState, useEffect } from 'react';
import styled from 'styles/styled-components';
import mobile2fa from 'images/mobile2faIcon.svg';
import { FormattedMessage } from 'react-intl';
import translate from '../messages';
import { Button } from '@material-ui/core';

import { useDispatch } from 'react-redux';
import { Subscriber, MessageNames } from 'services/message_service';
import IsLoadingWithText from 'components/isLoadingWithText/isLoadingWithText';
import { ChangePasswordModel } from '../types';
import { changePasswordAction } from '../actions';
import UBPinInput from 'components/pinInput';

export default function G2fa (props: {
  data: ChangePasswordModel;
  message: any;
  onClose: Function;
}) {
  const dispatch = useDispatch();
  const [CanSubmit, setCanSubmit] = useState(false);
  const [EnteredPin, setEnteredPin] = useState('');
  const [IsChangingPassword, setIsChangingPassword] = useState(false);
  const handleSubmitClick = () => {
    const sendingData = {
      old_password: props.data.old_password,
      new_password: props.data.new_password,
      confirmed: props.data.confirmed,
      '2fa_code': EnteredPin,
    };
    dispatch(changePasswordAction(sendingData));
  };
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SETLOADING) {
        setIsChangingPassword(message.payload);
      }
      if (message.name === MessageNames.CLOSE_MODAL) {
        props.onClose();
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);
  return (
    <Wrapper>
      <div className='iconWrapper'>
        <img src={mobile2fa} alt='' />
      </div>
      <div className='messageWrapper'>
        <span>{props.message}</span>
      </div>
      <div className='pinInputWrapper'>
        <UBPinInput
          onComplete={value => {
            setEnteredPin(value);
            setCanSubmit(true);
          }}
        />
      </div>
      <div className='buttonWrapper'>
        <Button
          disabled={!CanSubmit}
          onClick={handleSubmitClick}
          variant='contained'
          color='primary'
        >
          <IsLoadingWithText
            isLoading={IsChangingPassword}
            text={<FormattedMessage {...translate.submit} />}
          />
        </Button>
      </div>
    </Wrapper>
  );
}
const Wrapper = styled.div`
  width: 646px;
  height: 572px;
  display: flex;
  flex-direction: column;
  padding: 48px 80px;
  text-align: center;
  .iconWrapper {
    margin-bottom: 40px;
  }
  .messageWrapper {
    margin-bottom: 40px;
    span {
      color: var(--textGrey);
    }
  }
  .pinInputWrapper {
    margin-bottom: 40px;
  }
  .loadingCircle {
    top: 8px !important;
  }
`;
