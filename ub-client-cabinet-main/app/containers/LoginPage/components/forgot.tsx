import React, { useState, useEffect, useLayoutEffect } from 'react';
import translate from '../messages';
import { FormattedMessage } from 'react-intl';
import styled from 'styles/styled-components';
import forgotIcon from 'images/forgotIcon.svg';

import InputWithValidator from 'components/inputWithValidator';
import { EmailValidator } from '../validators/emailValidator';
import { Button } from '@material-ui/core';
import IsLoadingWithText from 'components/isLoadingWithText/isLoadingWithText';
import {
  Subscriber,
  MessageNames,
  MessageService,
} from 'services/message_service';
import { useDispatch } from 'react-redux';
import { forgotPasswordAction } from '../actions';
import { LocalStorageKeys } from 'services/constants';
import { isDevelopment } from 'utils/environment';

let fields = {
  forgotEmail: { isValid: false, value: '' },
};

export default function Forgot () {
  const [CanSubmit, setCanSubmit] = useState(false);
  const [IsLoading, setIsLoading] = useState(false);
  const dispatch = useDispatch();

  useLayoutEffect(() => {
    MessageService.send({
      name: MessageNames.RESET_RECAPTCHA,
    });
    return () => {};
  }, []);
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SETLOADING) {
        setIsLoading(message.payload);
      }
    });
    return () => {
      Subscription.unsubscribe();
      fields = {
        forgotEmail: { isValid: false, value: '' },
      };
    };
  }, []);
  const isFieldValid = (properties: {
    fieldName: string;
    isValid: boolean;
    value: string;
  }) => {
    fields[properties.fieldName].isValid = properties.isValid;
    fields[properties.fieldName].value = properties.value;
    if (fields['forgotEmail'].isValid === true) {
      setCanSubmit(true);
    } else {
      setCanSubmit(false);
    }
  };
  const handleSubmit = () => {
    const isDev = isDevelopment;
    dispatch(
      forgotPasswordAction({
        email: fields.forgotEmail.value,
        recaptcha: isDev
          ? 'recaptcha'
          : localStorage[LocalStorageKeys.RECAPTCHA],
      }),
    );
  };

  return (
    <Wrapper className={localStorage[LocalStorageKeys.Theme]}>
      <div className='titleWrapper'>
        <FormattedMessage {...translate.forgetPassword} />
      </div>
      <div className='iconWrapper'>
        <img src={forgotIcon} />
      </div>
      <div className='messageWrapper'>
        <FormattedMessage {...translate.Pleaseenteryouremailaddress} />
      </div>
      <div className='inputWrapper'>
        <InputWithValidator
          throttleTime={500}
          autoFocus
          label={<FormattedMessage {...translate.email} />}
          onChange={(email: string) => {
            isFieldValid({
              fieldName: 'forgotEmail',
              isValid: EmailValidator({
                uniqueInputId: 'forgotEmail',
                value: email,
              }),
              value: email,
            });
          }}
          uniqueName='forgotEmail'
        />
      </div>
      <div className='buttonWrapper'>
        <Button
          disabled={!CanSubmit}
          onClick={!IsLoading ? handleSubmit : () => {}}
          variant='contained'
          color='primary'
        >
          <IsLoadingWithText
            text={<FormattedMessage {...translate.submit} />}
            isLoading={IsLoading}
          />
        </Button>
      </div>
    </Wrapper>
  );
}
const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 60px 102px;
  text-align: center;
  height: 50vh;
  width: 480px;
  background: var(--white);
  .titleWrapper {
    flex: 5;
    span {
      color: var(--textBlue);
    }
  }
  .iconWrapper {
    flex: 30;
    display: flex;
    img {
      margin-right: -10px;
    }
  }
  .messageWrapper {
    flex: 10;
    display: flex;
    align-items: center;
    width: 226px;
    span {
      font-size: 12px;

      color: var(--textGrey);
    }
  }
  .inputWrapper {
    flex: 15;
    min-width: 278px;
  }
  .buttonWrapper {
    flex: 10;
  }
  .loadingCircle {
    top: 9px !important;
  }
`;
