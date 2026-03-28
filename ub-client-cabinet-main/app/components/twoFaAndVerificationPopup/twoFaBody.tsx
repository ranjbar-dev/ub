import { Button, TextField } from '@material-ui/core';
import IsLoadingWithText from 'components/isLoadingWithText/isLoadingWithText';
import UBPinInput from 'components/pinInput';
import NewEmailDesktopIcon from 'containers/FundsPage/pages/withdrawalsPage/components/NewEmailDesktopIcon';
import Mobile2faIcon from 'images/themedIcons/mobile2faIcon';
import React, { FC, useRef, useState } from 'react';
import styled, { css } from 'styles/styled-components';
import { ISubmitTwoFaAndEmailCode } from './types';

interface BodyProps {
  requiredData: { need2fa: boolean; needEmailCode: boolean };
  onFinalSubmit: (data: ISubmitTwoFaAndEmailCode) => void;
}
export const TwoFaAndEmailBody: FC<BodyProps> = (props) => {
  const{requiredData,onFinalSubmit}=props;
  const { need2fa, needEmailCode }=requiredData;
  const [CanSubmit, setCanSubmit] = useState(false);
  const EnteredPin = useRef('');
  const EnteredEmailCode = useRef('');
  const [Step, setStep] = useState<number>(
    (needEmailCode && need2fa) || needEmailCode ? 1 : 2,
  );

  const handleSubmitClick = () => {
    if (whatSteps() === 'emailAnd2fa' && Step === 1) {
      setStep(2);
      setCanSubmit(false);
      return;
    }
    const sendingData = {
      ...(need2fa && { '2fa_code': EnteredPin.current }),
      ...(needEmailCode && { email_code: EnteredEmailCode.current }),
    };
    onFinalSubmit(sendingData);
  };


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


  const whatSteps = (): 'onlyEmail' | 'only2fa' | 'emailAnd2fa' => {
    if (needEmailCode && need2fa) {
      return 'emailAnd2fa';
    } else if (need2fa) {
      return 'only2fa';
    } else return 'onlyEmail';
  };

  return (
    <Wrapper step={Step}>
      <StepsWrapper>
        <StepWrapper moved={Step === 2}>
          <div className='iconWrapper'>
            <NewEmailDesktopIcon />
          </div>
          <div className='messageWrapper'>
            <span>Enter the verification code from email</span>
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
            label="Enter code here"
          />
        </StepWrapper>

        <StepWrapper moved={Step === 2}>
          <div className='iconWrapper'>
            <Mobile2faIcon />
          </div>
          <div className='messageWrapper'>
              <span>Please Open Google Authenticator App And Enter 2FA Code</span>

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
            isLoading={false}
            text={
              whatSteps() === 'emailAnd2fa' && Step === 1
              ? 'Next'
              : 'Submit'
            }
          />
        </Button>
      </div>
    </Wrapper>
  );
};
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
    color: var(--blackText) !important;
  }
`;
