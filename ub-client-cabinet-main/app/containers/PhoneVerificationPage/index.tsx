/*
 *
 * PhoneVerificationPage
 *
 */

import React, { memo } from 'react';

import BreadCrumb from 'components/BreadCrumb';
import { MaxWidthWrapper } from 'components/wrappers/maxWidthWrapper';
import { replace, push } from 'redux-first-history';
import { AppPages } from 'containers/App/constants';
import { Helmet } from 'react-helmet';
import { FormattedMessage } from 'react-intl';
import { useDispatch, useSelector } from 'react-redux';
import { createStructuredSelector } from 'reselect';
import styled from 'styles/styled-components';
import { useInjectReducer } from 'utils/injectReducer';
import { useInjectSaga } from 'utils/injectSaga';

import { Card } from '@material-ui/core';

import {
  getSMSAction,
  resendSMSAction,
  setStepAction,
  verifyCodeAction,
} from './actions';
import StepIndicator from './components/StepIndicator';
import { PhoneVerificationSteps } from './constants';
import translate from './messages';
import reducer from './reducer';
import saga from './saga';
import makeSelectPhoneVerificationPage, {
  makeSelectUserData,
} from './selectors';
import Step1 from './steps/step1';
import Step2 from './steps/step2';
import { Step } from './types';
import Step3 from './steps/step3';
import { UserData } from 'containers/AcountPage/types';
import G2faStep from './steps/goggle2faStep';
import DoneStep from './steps/doneStep';
import { MaxContainer } from 'components/wrappers/maxContainer';
import { LocalStorageKeys } from 'services/constants';

const stateSelector = createStructuredSelector({
  phoneVerificationPage: makeSelectPhoneVerificationPage(),
  userData: makeSelectUserData(),
});

interface Props {}
let code = '';
function PhoneVerificationPage (props: Props) {
  // Warning: Add your key to RootState in types/index.d.ts file
  useInjectReducer({ key: 'phoneVerificationPage', reducer: reducer });
  useInjectSaga({ key: 'phoneVerificationPage', saga: saga });

  const { phoneVerificationPage, userData } = useSelector(stateSelector);
  const data: any = userData;
  const userdata: UserData = data;
  const selectedStep = phoneVerificationPage.activeStep;

  const dispatch = useDispatch();

  const isLoading = phoneVerificationPage.isLoading;
  const stepSelector = () => {
    switch (selectedStep) {
      case 0:
        return (
          <Step1
            countries={JSON.parse(localStorage[LocalStorageKeys.COUNTRIES])}
            onCancel={() => {
              dispatch(replace(AppPages.AcountPage));
            }}
            onSubmit={data => {
              dispatch(getSMSAction(data));
            }}
            submitIsLoading={isLoading}
          />
        );
      case 1:
        return (
          <Step2
            isLoading={isLoading}
            onGoToStep={(step: PhoneVerificationSteps) => {
              dispatch(setStepAction(step));
            }}
            onResend={() =>
              dispatch(
                resendSMSAction({
                  phone: phoneVerificationPage.enteredPhoneNumber,
                }),
              )
            }
            phone={phoneVerificationPage.enteredPhoneNumber}
            onCancel={() => {
              dispatch(replace(AppPages.AcountPage));
              dispatch(
                setStepAction(PhoneVerificationSteps.ENTER_PHONE_NUMBER),
              );
            }}
            onSubmit={data => {
              code = data;
              dispatch(
                verifyCodeAction({
                  code: data,
                  phone: phoneVerificationPage.enteredPhoneNumber,
                }),
              );
            }}
            submitIsLoading={isLoading}
          />
        );
      case 2:
        return (
          <Step3
            userData={userdata}
            on2fa={() => {
              dispatch(push(AppPages.GoogleAuthentication));
              // dispatch(setStepAction(PhoneVerificationSteps.GOOGLE_2FA_STEP));
            }}
            onCancel={() => {
              dispatch(replace(AppPages.AcountPage));
              dispatch(
                setStepAction(PhoneVerificationSteps.ENTER_PHONE_NUMBER),
              );
            }}
            onSubmit={data => {
              dispatch(
                verifyCodeAction({
                  code: code,
                  password: data,
                  phone: phoneVerificationPage.enteredPhoneNumber,
                }),
              );
            }}
            submitIsLoading={isLoading}
          />
        );
      case 3:
        return (
          <G2faStep
            onCancel={() => {
              dispatch(setStepAction(PhoneVerificationSteps.ENTER_CODE));
            }}
            onSubmit={data => {
              dispatch(
                verifyCodeAction({
                  '2fa_code': data,
                  code: code,
                  phone: phoneVerificationPage.enteredPhoneNumber,
                }),
              );
            }}
            code={code}
          />
        );
      case 4:
        return (
          <DoneStep
            onCancel={() => {
              dispatch(replace(AppPages.AcountPage));
              dispatch(
                setStepAction(PhoneVerificationSteps.ENTER_PHONE_NUMBER),
              );
            }}
            onSubmit={data => {}}
            submitIsLoading={isLoading}
          />
        );
      default:
        return (
          <Step1
            countries={phoneVerificationPage.countries}
            onCancel={() => dispatch(replace(AppPages.AcountPage))}
            onSubmit={data => {
              console.log(data);
            }}
            submitIsLoading={isLoading}
          />
        );
    }
  };
  const steps: Step[] = [
    {
      title: <FormattedMessage {...translate.step1} />,
      description: <FormattedMessage {...translate.enterPhoneNumber} />,
    },
    {
      title: <FormattedMessage {...translate.step2} />,
      description: <FormattedMessage {...translate.enterCode} />,
    },
    {
      title: <FormattedMessage {...translate.step3} />,
      description: <FormattedMessage {...translate.enterAcountPassword} />,
    },
  ];
  return (
    <div>
      <Helmet>
        <title>PhoneVerification</title>
        <meta
          name='description'
          content='Description of PhoneVerificationPage'
        />
      </Helmet>
      <MaxWidthWrapper>
        <BreadCrumb
          links={[
            { pageName: 'home', pageLink: AppPages.HomePage },
            {
              pageName: 'acountAndSecurity',
              pageLink: AppPages.AcountPage,
            },
            {
              pageName: 'phoneVerification',
              pageLink: AppPages.PhoneVerification,
              last: true,
            },
          ]}
        />
        <Wrapper>
          <MaxContainer>
            <StepIndicator steps={steps} selectesStep={selectedStep} />
            {stepSelector()}
          </MaxContainer>
        </Wrapper>
      </MaxWidthWrapper>
    </div>
  );
}

export default memo(PhoneVerificationPage);
const Wrapper = styled(Card)`
  background: white;
  border-radius: 10px !important;
  box-shadow: none !important;
  height: calc(98vh - 115px);
  min-height: 590px;
  display: flex;
  align-items: center;
  .autoComplete {
    input {
      margin-top: -5px;
      font-size: 14px;
    }
  }
  .selectedFlag {
    position: absolute;
    top: 67px;
    left: 10px;
    img {
      width: 25px;
      border-radius: 3px;
    }
  }
`;
