import React, { useState, useEffect } from 'react';
import { MainIconWrapper } from 'components/wrappers/mainIconWrapper';
import { CenterInputsWrapper } from 'components/wrappers/centerInputsWrapper';
import { CenterButtonsWrapper } from 'components/wrappers/centerButtonsWrapper';

import translate from '../messages';
import { FormattedMessage } from 'react-intl';

import { useDispatch } from 'react-redux';
import IsLoadingWithText from 'components/isLoadingWithText/isLoadingWithText';
import styled from 'styles/styled-components';
import { Button } from '@material-ui/core';
import { Card } from '@material-ui/core';
import { UserData } from 'containers/AcountPage/types';
import { SetG2FaModel } from '../types';
import { toggle2FaAction } from '../actions';
import { Subscriber, MessageNames } from 'services/message_service';
import { Buttons } from 'containers/App/constants';
import ChangePasswordMainIcon from 'images/themedIcons/changePasswordMainIcon';
import { MaxContainer } from 'components/wrappers/maxContainer';
import InputWithValidator from 'components/inputWithValidator';

export default function Step2 (props: {
  onCancel: Function;
  onSubmit: Function;
  code: string;
  userData: UserData;
}) {
  const [PasswordValue, setPasswordValue] = useState('');
  const [IsValidating, setIsValidating] = useState(false);
  const dispatch = useDispatch();

  const handleSubmit = () => {
    const sendingData: SetG2FaModel = {
      code: props.code,
      password: PasswordValue,
      setEnable: !props.userData.google2faEnabled,
    };
    dispatch(toggle2FaAction(sendingData));
  };
  const handleCancelButton = () => {
    props.onCancel();
  };
  const handlePasswordChange = (e: any) => {
    setPasswordValue(e);
  };
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SETLOADING) {
        setIsValidating(message.payload);
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);
  return (
    <Wrapper>
      <MaxContainer className='max740'>
        <MainIconWrapper className='noPadding'>
          <ChangePasswordMainIcon />
        </MainIconWrapper>
        <CenterInputsWrapper style={{ flex: 3, padding: 0 }}>
          <div className='centerHor'>
            <FormattedMessage {...translate.Pleaseenteryouraccountpassword} />
          </div>
          <InputWithValidator
            throttleTime={0}
            inputType='password'
            isPickable={true}
            label={<FormattedMessage {...translate.password} />}
            onChange={(password: string) => {
              handlePasswordChange(password);
              //isFieldValid({
              //  fieldName: 'changeGauth',
              //  isValid: PasswordValidator({
              //    uniqueInputId: 'changeGauth',
              //    value: password,
              //  }),
              //  value: password,
              //});
            }}
            uniqueName='changeGauth'
          />
        </CenterInputsWrapper>
        <CenterButtonsWrapper style={{ flex: 5, padding: 0 }}>
          <Button
            onClick={handleSubmit}
            className='ubButton'
            variant='contained'
            color='primary'
          >
            <IsLoadingWithText
              isLoading={IsValidating}
              text={<FormattedMessage {...translate.submit} />}
            />
          </Button>
          <Button className={Buttons.CancelButton} onClick={handleCancelButton}>
            <FormattedMessage {...translate.cancel} />
          </Button>
        </CenterButtonsWrapper>
      </MaxContainer>
    </Wrapper>
  );
}
const Wrapper = styled(Card)`
  width: 100%;
  background: white;
  border-radius: 10px !important;
  box-shadow: none !important;
  height: calc(98vh - 115px);
  padding: 77px;

  .noPadding {
    padding: 0 !important;
    flex: 10 !important;
  }
  .loadingCircle {
    top: 10px !important;
  }
`;
