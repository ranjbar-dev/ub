import React, { useEffect, useRef, useState } from 'react';

import { CenterButtonsWrapper } from 'components/wrappers/centerButtonsWrapper';
import { CenterInputsWrapper } from 'components/wrappers/centerInputsWrapper';
import { MainIconWrapper } from 'components/wrappers/mainIconWrapper';
import { Buttons } from 'containers/App/constants';
import RefreshIcon from 'images/themedIcons/refreshIcon';
import { FormattedMessage } from 'react-intl';
import { useDispatch } from 'react-redux';

import { Button } from '@material-ui/core';
import { TextField } from '@material-ui/core';

import { PhoneVerificationSteps } from '../constants';
import translate from './messages';
import PhoneVStep1Icon from 'images/themedIcons/phoneVStep1Icon';
import StreamLoadingButton from 'components/streamLoadingButton';

export default function Step2 (props: {
  onCancel: Function;
  onSubmit: Function;
  isLoading: boolean;
  onGoToStep: Function;
  onResend: Function;
  phone: string;
  submitIsLoading: boolean;
}) {
  // const INITIAL_INPUTS_STATE = {};
  // const {
  //   values,
  //   handleChange,
  //   handleBlur,
  //   errors,
  //   hasError,
  // } = UseFormValidation(INITIAL_INPUTS_STATE, PhoneNumberValidator);
  const dispatch = useDispatch();

  const handleSubmit = () => {
    props.onSubmit(codeValue);
  };

  const editPhone = () => {
    props.onGoToStep(PhoneVerificationSteps.ENTER_PHONE_NUMBER);
  };

  const handleCancelButton = () => {
    props.onCancel();
  };
  const [CountDownNumber, setCountDownNumber] = useState(59);
  const [CanResend, setCanResend] = useState(false);
  const [codeValue, setCodeValue] = useState('');
  const [canSubmit, setCanSubmit] = useState(false);
  const codeInputChange = e => {
    if (e.target.value.length === 7) {
      return;
    }
    if (e.target.value.length < 7) {
      setCodeValue(e.target.value);
    }
    if (e.target.value.length === 6) {
      setCanSubmit(true);
    } else {
      setCanSubmit(false);
    }
    // if (e.target.value.length == 6) {
    //   props.onSubmit(e.target.value);
    // }
  };

  useEffect(() => {
    const timer = setInterval(() => {
      if (CountDownNumber > 1) {
        setCountDownNumber(CountDownNumber - 1);
      } else {
        clearInterval(timer);
        setCanResend(true);
      }
    }, 1050);

    return () => {
      clearInterval(timer);
    };
  });
  return (
    <>
      <MainIconWrapper className='fl fl9 minimized'>
        <PhoneVStep1Icon />
      </MainIconWrapper>
      <CenterInputsWrapper className='noPadd pt3 even' style={{ flex: 6 }}>
        <div className='flexSpacer1'></div>
        <div className='fl1'>
          <span>
            <FormattedMessage {...translate.weAreSendingCodeTo} />
          </span>
          <span className='phoneNumber'>{props.phone}</span>
        </div>
        <div className='fl1'>
          <FormattedMessage {...translate.pleaseCheckYourPhoneAndEnter} />

          <span className='boldGrey'>
            <FormattedMessage {...translate.AuthenticationCode} />
          </span>
        </div>
        <div className='restrictedInputWrapper fl3 spaceAround'>
          <TextField
            variant='outlined'
            autoFocus
            margin='dense'
            onChange={codeInputChange}
            value={codeValue}
            className='restrictedInput'
            label={<FormattedMessage {...translate.enterCodeHere} />}
          />
          <StreamLoadingButton
            disabled={!canSubmit || codeValue.length < 6}
            onClick={handleSubmit}
            className='ubButton nmt'
            variant='contained'
            color='primary'
            text={<FormattedMessage {...translate.submit} />}
          />
          <div>
            {CanResend && (
              <span style={{ fontStyle: 'italic' }}>
                <FormattedMessage {...translate.Dontreceivecode} />
              </span>
            )}
            <Button
              className={Buttons.TransParentRoundButton}
              onClick={editPhone}
              color='primary'
            >
              <FormattedMessage {...translate.editPhoneNumber} />
            </Button>
            {CanResend && <FormattedMessage {...translate.or} />}
            {!CanResend && (
              <span>
                00:
                {CountDownNumber >= 10
                  ? CountDownNumber
                  : '0' + CountDownNumber}
              </span>
            )}
            {CanResend && (
              <Button
                color='primary'
                onClick={() => {
                  setCanResend(false);
                  setCountDownNumber(59);
                  props.onResend();
                }}
                className={Buttons.TransParentRoundButton}
              >
                <FormattedMessage {...translate.resend} />
                <RefreshIcon />
              </Button>
            )}
          </div>
        </div>
      </CenterInputsWrapper>
      <CenterButtonsWrapper>
        <Button
          className={Buttons.SimpleRoundButton}
          onClick={handleCancelButton}
        >
          <FormattedMessage {...translate.backToDashboard} />
        </Button>
      </CenterButtonsWrapper>
    </>
  );
}
