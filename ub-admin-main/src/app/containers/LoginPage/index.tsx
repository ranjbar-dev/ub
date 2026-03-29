/**
 *
 * LoginPage
 *
 */

import Button from '@material-ui/core/Button';
import InputWithValidator from 'app/components/inputWithValidator';
import IsLoadingWithText from 'app/components/isLoadingWithText/isLoadingWithText';
import { Buttons } from 'app/constants';
import bgImage2 from 'images/staticImages/bg-pattern-2.png';
import bgImage from 'images/staticImages/bg-pattern.png';
import { translations } from 'locales/i18n';
import React, { memo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useSelector, useDispatch } from 'react-redux';
import styled from 'styled-components/macro';
import { StyleConstants } from 'styles/StyleConstants';
import { useInjectReducer, useInjectSaga } from 'utils/redux-injectors';

import { loginPageSaga } from './saga';
import { selectLoading, selectLoginError } from './selectors';
import { reducer, sliceKey, actions } from './slice';
import { EmailValidator } from './validators/emailValidator';
import { PasswordValidator } from './validators/passwordValidator';

interface Props {}
let fields: Record<string, { isValid: boolean; value: string }> = {
  loginEmail: { isValid: false, value: '' },
  loginPassword: { isValid: false, value: '' },
};

export const LoginPage = memo((props: Props) => {
  useInjectReducer({ key: sliceKey, reducer: reducer });
  useInjectSaga({ key: sliceKey, saga: loginPageSaga });
  const [CanSubmit, setCanSubmit] = useState(false);

  const isLoading = useSelector(selectLoading);
  const loginError = useSelector(selectLoginError);
  const dispatch = useDispatch();

  const { t } = useTranslation();
  const isFieldValid = (properties: {
    fieldName: string;
    isValid: boolean;
    value: string;
  }) => {
    fields[properties.fieldName].isValid = properties.isValid;
    fields[properties.fieldName].value = properties.value;
    if (
      fields['loginEmail'].isValid === true &&
      fields['loginPassword'].isValid === true
    ) {
      setCanSubmit(true);
    } else {
      setCanSubmit(false);
    }
  };
  const handleEnterKeyDown = () => {
    if (
      fields['loginEmail'].isValid === true &&
      fields['loginPassword'].isValid === true
    ) {
      loginClick();
    }
  };

  const checkForValidation = () => {
    const inp = document.querySelector<HTMLInputElement>('#loginEmail');
    if (inp) {
      setTimeout(() => {
        isFieldValid({
          fieldName: 'loginEmail',
          isValid: EmailValidator({
            uniqueInputId: 'loginEmail',
            value: inp.value,
            errors: [
              t(translations.Errors.emailIsRequired()),
              t(translations.Errors.emailIsNotValid()),
            ],
          }),
          value: inp.value,
        });
      }, 50);
    }
  };
  const loginClick = () => {
    dispatch(
      actions.LoginAction({
        username: fields['loginEmail'].value,
        password: fields['loginPassword'].value,
      }),
    );
  };

  return (
    <LoginWrapper>
      <div className="UbTitle">UNITEDBIT ADMIN</div>
      <div className="inputsWrapper">
        <InputWithValidator
          throttleTime={500}
          inputType="text"
          label={t(translations.LoginPage.Email())}
          onChange={(email: string) => {
            isFieldValid({
              fieldName: 'loginEmail',
              isValid: EmailValidator({
                uniqueInputId: 'loginEmail',
                value: email,
                errors: [
                  t(translations.Errors.emailIsRequired()),
                  t(translations.Errors.emailIsNotValid()),
                ],
              }),
              value: email,
            });
          }}
          uniqueName="loginEmail"
        />
        <InputWithValidator
          throttleTime={500}
          inputType="password"
          className="lastInput"
          isPickable={true}
          onEnter={handleEnterKeyDown}
          label={t(translations.LoginPage.password())}
          onChange={(password: string) => {
            isFieldValid({
              fieldName: 'loginPassword',
              isValid: PasswordValidator({
                uniqueInputId: 'loginPassword',
                value: password,
                errors: [
                  t(translations.Errors.passwordIsRequired()),
                  t(translations.Errors.minimum8Character()),
                ],
              }),
              value: password,
            });
            checkForValidation();
          }}
          uniqueName="loginPassword"
        />
        {loginError && (
          <p style={{ color: 'red', fontSize: 12, marginTop: -10 }}>{loginError}</p>
        )}
        <Button
          color="primary"
          onClick={() => {
            loginClick();
          }}
          disabled={!CanSubmit}
          className={Buttons.SubmitButton}
          variant="contained"
        >
          <IsLoadingWithText
            isLoading={isLoading}
            text={t(translations.LoginPage.Login())}
          />
        </Button>
      </div>
    </LoginWrapper>
  );
});

const LoginWrapper = styled.div`
  .UbTitle {
    margin-top: -130px;
    color: white;
    text-shadow: 0 0 10px #000000;
  }
  display: grid;
  align-content: center;
  justify-items: center;
  background-image: url(${bgImage});
  background-color: #55545f;
  background-size: cover;
  height: 100%;
  width: 100%;
  padding-top: ${StyleConstants.NAV_BAR_HEIGHT};
  .inputsWrapper {
    width: 354px;
    height: 300px;
    padding: 35px;
    background: ${p => p.theme.background};
    background-image: url(${bgImage2});
    border-radius: 5px;
    display: flex;
    background-size: cover;
    flex-direction: column;
    margin-top: -100px;
    .MuiFormLabel-root {
      font-size: 14px !important;
    }
  }
  .lastInput {
    margin-bottom: 20px;
  }
`;
