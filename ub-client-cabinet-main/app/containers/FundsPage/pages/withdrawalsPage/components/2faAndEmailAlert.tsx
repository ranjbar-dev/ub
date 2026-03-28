import React, { useState, useEffect, useRef } from 'react';
import styled, { css } from 'styles/styled-components';

import { FormattedMessage } from 'react-intl';
import translate from '../../../messages';
import { Button, TextField } from '@material-ui/core';
import { useDispatch } from 'react-redux';
import { WithdrawModel } from 'containers/FundsPage/types';
import { withdrawAction } from 'containers/FundsPage/actions';
import { Subscriber, MessageNames } from 'services/message_service';
import IsLoadingWithText from 'components/isLoadingWithText/isLoadingWithText';
import UBPinInput from 'components/pinInput';
import Mobile2faIcon from 'images/themedIcons/mobile2faIcon';
import NewEmailDesktopIcon from './NewEmailDesktopIcon';

export default function G2fa (props: {
  code: string;
  amount: string;
  address: string;
  label: string;
  intl: any;
  network?: string;
  onClose: () => void;
  requiredData: { need2fa: boolean; needEmailCode: boolean };
}) {
  const {
    code,
    amount,
    address,
    label,
    intl,
    network,
    requiredData: { need2fa, needEmailCode },
  } = props;

  const whatSteps = (): 'onlyEmail' | 'only2fa' | 'emailAnd2fa' => {
    if (needEmailCode && need2fa) {
      return 'emailAnd2fa';
    } else if (need2fa) {
      return 'only2fa';
    } else return 'onlyEmail';
  };

  const dispatch = useDispatch();
  const [CanSubmit, setCanSubmit] = useState(false);
  const EnteredPin = useRef('');
  const EnteredEmailCode = useRef('');
  const [IsWithdrawing, setIsWithdrawing] = useState(false);

  const [Step, setStep] = useState<number>(
    (needEmailCode && need2fa) || needEmailCode ? 1 : 2,
  );

  const handleSubmitClick = () => {
    if (whatSteps() === 'emailAnd2fa' && Step === 1) {
      setStep(2);
      setCanSubmit(false);
      return;
    }
    const sendingData: WithdrawModel = {
      code,
      amount,
      address,
      label,
      ...(network && { network }),
      ...(need2fa && { G2fa_code: EnteredPin.current }),
      ...(needEmailCode && { email_code: EnteredEmailCode.current }),
    };
    dispatch(withdrawAction(sendingData));
  };

  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SETLOADING) {
        setIsWithdrawing(message.payload);
      }
      if (message.name === MessageNames.ADD_DATA_ROW_TO_WITHDRAWS) {
        props.onClose();
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);
  const checkValidation = () => {
    if (whatSteps() === 'emailAnd2fa') {
      setCanSubmit(true);
    } else if (whatSteps() === 'only2fa') {
      if (EnteredPin.current) {
        setCanSubmit(true);
      }
    } else if (whatSteps() === 'onlyEmail') {
      if (EnteredEmailCode) {
        setCanSubmit(true);
      }
    }
  };
  return (
    <Wrapper step={Step}>
      <StepsWrapper>
        <StepWrapper moved={Step === 2}>
          <div className='iconWrapper'>
            <NewEmailDesktopIcon />
          </div>
          <div className='messageWrapper'>
            <FormattedMessage {...translate.EnterVerificationEmail} />
          </div>
          <TextField
            variant='outlined'
            margin='dense'
            fullWidth
            className='codeInput'
            onChange={({ target: { value } }) => {
              EnteredEmailCode.current = value;
              checkValidation();
            }}
            label={intl.formatMessage({
              id: 'app.globalTitles.EnterCodeHere',
              defaultMessage: 'ET.Code',
            })}
          />
        </StepWrapper>

        <StepWrapper moved={Step === 2}>
          <div className='iconWrapper'>
            <Mobile2faIcon />
          </div>
          <div className='messageWrapper'>
            <FormattedMessage
              {...translate.PleaseopenGoogleAuthenticatorappinyourphoneandenter2FaCode}
            />
          </div>
          <div className='pinInputWrapper'>
            <UBPinInput
              onComplete={value => {
                EnteredPin.current = value;
                checkValidation();
              }}
            />
          </div>
        </StepWrapper>
      </StepsWrapper>
      <div className='buttonWrapper'>
        <Button
          disabled={!CanSubmit}
          onClick={handleSubmitClick}
          variant='contained'
          color='primary'
        >
          <IsLoadingWithText
            isLoading={IsWithdrawing}
            text={
              <FormattedMessage
                {...translate[
                  whatSteps() === 'emailAnd2fa' && Step === 1
                    ? 'next'
                    : 'submit'
                ]}
              />
            }
          />
        </Button>
      </div>
    </Wrapper>
  );
}
const StepWrapper = styled.div<{ moved: boolean }>`
  transition: transform 0.5s cubic-bezier(0.4, 0, 0.2, 1);
  min-width: 100%;
  ${({ moved }) =>
    moved &&
    css`
      transform: translateX(-100%);
    `}
`;
const StepsWrapper = styled.div`
  display: flex;
  overflow: hidden;
`;
const Wrapper = styled.div<{ step: number }>`
  width: 646px;
  height: 572px;
  transition: height 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  ${({ step }) =>
    step === 1 &&
    css`
      height: 475px;
    `}
  display: flex;
  flex-direction: column;
  padding: 48px 80px;
  text-align: center;

  .messageWrapper {
    span {
      color: var(--textGrey);
    }
  }
  .iconWrapper,
  .messageWrapper,
  .pinInputWrapper,
  .codeInput {
    margin-bottom: 40px;
  }

  .loadingCircle {
    top: 8px !important;
  }
  .MuiFormLabel-root {
    color: var(--placeHolderColor) !important;
    font-size: 13px;
  }
`;
