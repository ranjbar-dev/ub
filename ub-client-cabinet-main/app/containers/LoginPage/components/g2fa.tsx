import React, { useState, useEffect } from 'react';
import styled from 'styles/styled-components';

import { FormattedMessage } from 'react-intl';
import translate from '../messages';
import { Button } from '@material-ui/core';
import errorIcon from 'images/errorIcon.svg';

import { useDispatch } from 'react-redux';

import { Subscriber, MessageNames } from 'services/message_service';
import IsLoadingWithText from 'components/isLoadingWithText/isLoadingWithText';
// import { loginAction } from 'containers/HomePage/actions';
import UBPinInput from 'components/pinInput';
import { loginAction } from '../actions';
import { LocalStorageKeys } from 'services/constants';
import Mobile2faIcon from 'images/themedIcons/mobile2faIcon';

export default function LoginG2fa (props: {
  username: string;
  password: string;
  message: any;

  onClose: Function;
}) {
  const dispatch = useDispatch();
  const [CanSubmit, setCanSubmit] = useState(false);
  const [EnteredPin, setEnteredPin] = useState('');
  const [ShowError, setShowError] = useState(false);
  const [IsWithdrawing, setIsWithdrawing] = useState(false);
  const handleSubmitClick = (pin?: string) => {
    const sendingData = {
      username: props.username,
      password: props.password,
      recaptcha: localStorage[LocalStorageKeys.RECAPTCHA],
      '2fa_code': pin ?? EnteredPin,
    };
    dispatch(loginAction(sendingData));
  };
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SETLOADING) {
        setIsWithdrawing(message.payload);
      }
      if (message.name === MessageNames.LOGGED_IN) {
        props.onClose();
      }
      if (message.name === MessageNames.AUTH_ERROR_EVENT) {
        setShowError(true);
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);
  return (
    <Wrapper
      className={localStorage[LocalStorageKeys.Theme]}
      onKeyDown={e => {
        if (e.keyCode === 13) {
          if (CanSubmit) {
            handleSubmitClick();
          }
        }
      }}
    >
      <div className='iconWrapper'>
        <Mobile2faIcon />
      </div>
      <div className='messageWrapper'>
        <span>{props.message}</span>
      </div>
      <div className='pinInputWrapper'>
        <UBPinInput
          onComplete={value => {
            setEnteredPin(value);
            setCanSubmit(true);
            handleSubmitClick(value);
          }}
        />
        {ShowError && (
          <div className='errorWrapper'>
            <span className='errorIcon'>
              <img src={errorIcon} alt='' />
            </span>
            <span className='errorText'>2fa Authentication is not valid</span>
          </div>
        )}
      </div>
      <div className='buttonWrapper'>
        <Button
          disabled={!CanSubmit}
          onClick={() => handleSubmitClick()}
          variant='contained'
          color='primary'
        >
          <IsLoadingWithText
            isLoading={IsWithdrawing}
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
  background: var(--white);
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
    position: relative;
  }
  .loadingCircle {
    top: 8px !important;
  }
  .errorWrapper {
    position: absolute;
    bottom: -20px;
    left: 17px;
    color: var(--redText);
    font-size: 11px;
    min-width: 400px;
    display: flex;
    span {
      font-size: 11px;
    }
    img {
      width: 20px;
    }
  }
`;
