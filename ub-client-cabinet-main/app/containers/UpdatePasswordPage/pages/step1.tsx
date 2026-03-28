import React, { useState, useEffect, useRef } from 'react';
import styled from 'styles/styled-components';
import resetPasswordIcon from 'images/resetPasswordIcon.svg';

import { replace } from 'redux-first-history';
import InputWithValidator from 'components/inputWithValidator';
import validationTranslate from 'containers/SignupPage/validators/messages';

import { FormattedMessage } from 'react-intl';
import { useDispatch } from 'react-redux';

import translate from '../messages';

import { AppPages, Buttons } from 'containers/App/constants';

import {
  Subscriber,
  MessageNames,
  MessageService,
} from 'services/message_service';
import { UpdatePasswordModel } from '../types';
import { Button } from '@material-ui/core';
import IsLoadingWithText from 'components/isLoadingWithText/isLoadingWithText';
import { Card } from '@material-ui/core';
import { resetPasswordAction } from '../actions';
import MainAppIcon from 'images/themedIcons/mainAppIcon';
import {
  passwordContentValidator,
  PasswordErrors,
} from 'containers/SignupPage/validators/passwordContentValidator';

let fields = {
  newPassword: { isValid: false, value: '' },
  confirmNewPassword: { isValid: false, value: '' },
};
enum Ids {
  newPasswordId = 'newPassword',
  ConfirmPasswordId = 'confirmPassword',
}

export default function Step1 (props: { email: string; code: string }) {
  const [CanSubmit, setCanSubmit] = useState(false);
  const [IsLoading, setIsLoading] = useState(false);

  const validatedFields = useRef<{
    [Ids.newPasswordId]: boolean;
    [Ids.ConfirmPasswordId]: boolean;
  }>({
    [Ids.newPasswordId]: false,
    [Ids.ConfirmPasswordId]: false,
  });
  const enteredFields = useRef<{
    [Ids.newPasswordId]: string;
    [Ids.ConfirmPasswordId]: string;
  }>({
    [Ids.newPasswordId]: '',
    [Ids.ConfirmPasswordId]: '',
  });

  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SETLOADING) {
        setIsLoading(message.payload);
      }
    });
    return () => {
      fields = {
        newPassword: { isValid: false, value: '' },
        confirmNewPassword: { isValid: false, value: '' },
      };
      Subscription.unsubscribe();
    };
  }, []);
  //custom hooks

  const dispatch = useDispatch();
  //items wrapper

  //button actions
  const handleSubmit = () => {
    const data: UpdatePasswordModel = {
      password: enteredFields.current[Ids.newPasswordId],
      confirmed: enteredFields.current[Ids.ConfirmPasswordId],
      email: props.email,
      code: props.code,
    };
    dispatch(resetPasswordAction(data));
  };
  const handleCancelButton = () => {
    dispatch(replace(AppPages.LoginPage));
  };

  const checkForOverallValidation = () => {
    let isAllGood = true;

    if (
      validatedFields.current[Ids.newPasswordId] === false ||
      validatedFields.current[Ids.ConfirmPasswordId] === false
    ) {
      isAllGood = false;
    }
    if (isAllGood === true) {
      if (CanSubmit === false) {
        setCanSubmit(true);
      }
    } else {
      if (CanSubmit === true) {
        setCanSubmit(false);
      }
    }
  };

  const validator = ({
    fieldId,
    value,
    compareWithFieldId,
  }: {
    fieldId: string;
    value: string;
    compareWithFieldId: string;
  }) => {
    enteredFields.current[fieldId] = value;

    const passError = passwordContentValidator({ password: value });
    let err;

    switch (passError) {
      case PasswordErrors.Min8:
        err = <FormattedMessage {...validationTranslate.minimum8Character} />;
        break;
      case PasswordErrors.NoNumber:
        err = <FormattedMessage {...validationTranslate.number} />;
        break;
      case PasswordErrors.NoUppercase:
        err = <FormattedMessage {...validationTranslate.upper} />;
        break;
      case PasswordErrors.NoSpecialCharacter:
        err = <FormattedMessage {...validationTranslate.special} />;
        break;

      default:
        err = null;
        break;
    }
    if (
      !err &&
      compareWithFieldId &&
      enteredFields.current[compareWithFieldId] !== ''
    ) {
      if (value !== enteredFields.current[compareWithFieldId]) {
        err = <FormattedMessage {...validationTranslate.oldAndNewPassword} />;
      } else {
        validatedFields.current[fieldId] = true;
        validatedFields.current[compareWithFieldId] = true;
        err = null;
        MessageService.send({
          name: MessageNames.SET_INPUT_ERROR,
          value: fieldId,
          payload: err,
        });
        MessageService.send({
          name: MessageNames.SET_INPUT_ERROR,
          value: compareWithFieldId,
          payload: err,
        });
      }
    }
    MessageService.send({
      name: MessageNames.SET_INPUT_ERROR,
      value: fieldId,
      payload: err,
    });
    if (err) {
      validatedFields.current[fieldId] = false;
    }
    checkForOverallValidation();
  };

  return (
    <Wrapper>
      <div className='logo'>
        <MainAppIcon />
      </div>
      <div className='titleWrapper blue'>
        <FormattedMessage {...translate.ResetPassword} />
      </div>
      <div className='mainIcon'>
        <img src={resetPasswordIcon} alt='' />
      </div>
      <div className='inputsWrapper'>
        <div className='inputWrapper'>
          <InputWithValidator
            inputType='password'
            isPickable
            throttleTime={500}
            label={<FormattedMessage {...translate.NewPassword} />}
            onChange={(password: string) => {
              validator({
                fieldId: Ids.newPasswordId,
                value: password,
                compareWithFieldId: Ids.ConfirmPasswordId,
              });
            }}
            uniqueName={Ids.newPasswordId}
          />
        </div>
        <div className='inputWrapper'>
          <InputWithValidator
            inputType='password'
            throttleTime={500}
            isPickable
            label={<FormattedMessage {...translate.ConfirmNewPassword} />}
            onChange={(confirmPassword: string) => {
              validator({
                fieldId: Ids.ConfirmPasswordId,
                value: confirmPassword,
                compareWithFieldId: Ids.newPasswordId,
              });
            }}
            uniqueName={Ids.ConfirmPasswordId}
          />
        </div>
      </div>
      <div className='buttonsWrapper'>
        <Button
          disabled={!CanSubmit}
          onClick={() => (!IsLoading ? handleSubmit() : null)}
          className='ubButton'
          fullWidth
          variant='contained'
          color='primary'
        >
          <IsLoadingWithText
            text={<FormattedMessage {...translate.ResetPassword} />}
            isLoading={IsLoading}
          />
        </Button>
        <Button className={Buttons.CancelButton} onClick={handleCancelButton}>
          <FormattedMessage {...translate.cancel} />
        </Button>
      </div>
    </Wrapper>
  );
}
const Wrapper = styled(Card)`
  border-radius: 10px !important;
  box-shadow: none !important;
  height: 100vh;
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 131px 0;
  .logo {
    flex: 1;
    svg {
      min-width: 14.5vw;
      min-height: 55px;
    }
  }
  .titleWrapper {
    flex: 1;
    span {
      font-size: 22px;
    }
  }
  .mainIcon {
    padding: 0vh 0 1vh 0;
    flex: 4;
  }
  .inputsWrapper {
    min-width: 17vw;
    flex: 1;
    display: flex;
    flex-direction: column;
    place-content: space-around;
    margin-bottom: 2vh;
  }
  .buttonsWrapper {
    min-width: 12vw;
    flex: 4;
    display: flex;
    flex-direction: column;
    align-items: center;
    .MuiButton-root {
      margin: 0.5vh 0;
    }
  }
`;
