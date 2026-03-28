import React, { useState, useEffect, useRef, useLayoutEffect } from 'react';
import { CircularProgress, Button } from '@material-ui/core';

import LoginG2fa from './components/g2fa';
import {
  Subscriber,
  MessageNames,
  MessageService,
} from 'services/message_service';
import MainAppIcon from 'images/themedIcons/mainAppIcon';
import styled from 'styles/styled-components';

import LoginMainIcon from 'images/themedIcons/loginMainIcon';
import { Themes, AppPages, Buttons } from 'containers/App/constants';
import { loginAction } from './actions';

import InputWithValidator from 'components/inputWithValidator';
import { EmailValidator } from './validators/emailValidator';
import { PasswordValidator } from './validators/passwordValidator';
import Forgot from './components/forgot';
import PopupModal from 'components/materialModal/modal';
import { push } from 'redux-first-history';
import { FormattedMessage } from 'react-intl';
import { useSelector, useDispatch } from 'react-redux';
import { createStructuredSelector } from 'reselect';

import { useInjectSaga } from 'utils/injectSaga';
import { useInjectReducer } from 'utils/injectReducer';
import makeSelectLoginPage, { makeSelectIsLoadingLoginPage } from './selectors';
import reducer from './reducer';
import saga from './saga';
import translate from './messages';

import { makeStyles } from '@material-ui/core';
import CenterWrapp from 'components/wrappers/centerWrapp';
import { LandingPageAddress, LocalStorageKeys } from 'services/constants';
import { hSmallerThan730 } from 'styles/mediaQueries';
import { isDevelopment } from 'utils/environment';
const stateSelector = createStructuredSelector({
  loginPage: makeSelectLoginPage(),
  isLoggingIn: makeSelectIsLoadingLoginPage(),
});
const materialClasses = makeStyles({
  loadingIndicator: {
    color: 'white',
  },
  loginButton: {
    minHeight: '40px',
  },
});

export default function LoginBody(props: { isPopup?: boolean }) {
  const isDev = isDevelopment;

  useInjectReducer({ key: 'loginPage', reducer: reducer });
  useInjectSaga({ key: 'loginPage', saga: saga });
  useLayoutEffect(() => {
    MessageService.send({
      name: MessageNames.RESET_RECAPTCHA,
    });
    return () => {};
  }, []);
  const [RememberMe] = useState(true);
  const [Recaptcha, setRecaptcha] = useState('');
  const [CanSubmit] = useState(true);
  const [IsForgotOpen, setIsForgotOpen] = useState(false);
  const [G2FaData, setG2FaData] = useState({
    username: '',
    password: '',
    message: '',
  });
  const InputValues = useRef({
    loginEmail: { isValid: false, value: '' },
    loginPassword: { isValid: false, value: '' },
  });
  const [IsG2faModalOpen, setIsG2faModalOpen] = useState(false);
  //////

  const loginClick = () => {
    if (!Recaptcha && !isDev) {
      return;
    }
    const emailInp: any = document.querySelector('#loginEmail');
    if (emailInp) {
      InputValues.current.loginEmail.value = emailInp.value;
    }
    const passwordInp: any = document.querySelector('#loginPassword');
    if (passwordInp) {
      InputValues.current.loginPassword.value = passwordInp.value;
    }
    const { loginEmail, loginPassword } = InputValues.current;
    ////validators are outside the component because they may needed to be externally executed
    loginEmail.isValid = EmailValidator({
      uniqueInputId: 'loginEmail',
      value: loginEmail.value,
    });
    loginPassword.isValid = PasswordValidator({
      uniqueInputId: 'loginPassword',
      value: loginPassword.value,
    });
    if (loginEmail.isValid && loginPassword.isValid) {
      dispatch(
        loginAction({
          username: loginEmail.value,
          password: loginPassword.value,
          remember: RememberMe,
          recaptcha: isDev
            ? 'recaptcha'
            : localStorage[LocalStorageKeys.RECAPTCHA],
          ...(props.isPopup && { fromPopup: true }),
        }),
      );
    }
  };

  const classes = materialClasses();

  const isLoggingIn = () => {
    if (isDev) {
      return false;
    }
    if (!loginPage) {
      return (
        <span>
          <FormattedMessage {...translate.login} />
        </span>
      );
    }
    return loginPage.isLoggingIn === true || Recaptcha == '' ? (
      <CircularProgress size={14} className={classes.loadingIndicator} />
    ) : (
      <span>
        <FormattedMessage {...translate.login} />
      </span>
    );
  };
  const { loginPage } = useSelector(stateSelector);
  const dispatch = useDispatch();

  const setInputValues = (data: { fieldName: string; value: string }) => {
    const { fieldName, value } = data;
    InputValues.current[fieldName].value = value;
  };

  useEffect(() => {
    const timeout = setTimeout(() => {
      if (!localStorage[LocalStorageKeys.RECAPTCHA]) {
        MessageService.send({
          name: MessageNames.RESET_RECAPTCHA,
        });
      }
    }, 10000);
    return () => {
      clearTimeout(timeout);
    };
  }, []);

  useEffect(() => {
    const timeout = setInterval(() => {
      MessageService.send({
        name: MessageNames.RESET_RECAPTCHA,
      });
    }, 110000);
    return () => {
      clearInterval(timeout);
    };
  }, []);

  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name == MessageNames.OPEN_G2FA) {
        setG2FaData(message.payload);
        setIsG2faModalOpen(true);
      } else if (message.name === MessageNames.CLOSE_MODAL) {
        setIsForgotOpen(false);
      } else if (message.name === MessageNames.AUTH_ERROR_EVENT) {
        if (!IsG2faModalOpen) {
          // setCanSubmit(false);
          MessageService.send({
            name: MessageNames.SET_INPUT_ERROR,
            value: 'loginPassword',
            payload: (
              <FormattedMessage {...translate.invalidUsernameOrPassword} />
            ),
          });
        }
      } else if (message.name === MessageNames.SET_RECAPTCHA) {
        setRecaptcha(message.payload);
      } else if (message.name === MessageNames.RESET_RECAPTCHA) {
        setRecaptcha('');
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, [IsG2faModalOpen]);
  ////clean up useEffect

  const handleForgotClick = () => {
    setIsForgotOpen(true);
  };
  const handleSignupClick = () => {
    dispatch(push(AppPages.SignupPage));
  };
  const handleEnterKeyDown = () => {
    loginClick();
  };

  return (
    <CenterWrapp className={`${props.isPopup ? 'popup' : ''}`}>
      <LoginBodyWrapper className={`${props.isPopup ? 'popup' : ''}`}>
        <PopupModal
          isOpen={IsG2faModalOpen}
          onClose={() => {
            setIsG2faModalOpen(false);
          }}
        >
          <LoginG2fa
            username={G2FaData.username}
            password={G2FaData.password}
            message={G2FaData.message}
            onClose={() => {
              setIsG2faModalOpen(false);
              // setCanSubmit(true);
            }}
          />
        </PopupModal>
        <PopupModal
          isOpen={IsForgotOpen}
          onClose={() => {
            setIsForgotOpen(false);
          }}
        >
          {IsForgotOpen === true && <Forgot />}
        </PopupModal>

        {!props.isPopup && (
          <div className="logo">
            <MainAppIcon />
          </div>
        )}
        {!props.isPopup && (
          <div className="message">
            <FormattedMessage {...translate.loginMessage} />
          </div>
        )}
        <div className="icon">
          <LoginMainIcon theme={Themes.LIGHT} />
        </div>
        <div className="mainContent">
          <div className="inputsWrapper">
            <InputWithValidator
              throttleTime={5}
              inputType="text"
              label={<FormattedMessage {...translate.email} />}
              onChange={(email: string) => {
                setInputValues({ fieldName: 'loginEmail', value: email });
              }}
              uniqueName="loginEmail"
            />
            <InputWithValidator
              throttleTime={5}
              inputType="password"
              className="lastInput"
              isPickable={true}
              onEnter={handleEnterKeyDown}
              label={<FormattedMessage {...translate.password} />}
              onChange={(password: string) => {
                setInputValues({ fieldName: 'loginPassword', value: password });
              }}
              uniqueName="loginPassword"
            />
          </div>
          <Button
            disabled={!isDev && !CanSubmit}
            variant="contained"
            fullWidth
            onClick={
              isDev
                ? loginClick
                : (loginPage && loginPage.isLoggingIn) ||
                  Recaptcha === '' ||
                  !CanSubmit
                ? () => {}
                : loginClick
            }
            color="primary"
            className={classes.loginButton}
          >
            {isLoggingIn()}
          </Button>
          <div className="centerHor">
            <Button
              onClick={handleForgotClick}
              className={`button blue ${Buttons.Underlined} `}
            >
              <FormattedMessage {...translate.forgetPassword} />
            </Button>
            <span className="black">|</span>
            <Button
              onClick={handleSignupClick}
              className={`button blue ${Buttons.Underlined} `}
            >
              <FormattedMessage {...translate.signup} />
            </Button>
          </div>

          {!props.isPopup && (
            <div className="centerHor">
              <Button
                onClick={() => location.replace(LandingPageAddress)}
                className={`button blue shadedButton ${Buttons.SimpleRoundButton}`}
              >
                <FormattedMessage {...translate.go_back_home} />
              </Button>
            </div>
          )}
        </div>
      </LoginBodyWrapper>
    </CenterWrapp>
  );
}
const LoginBodyWrapper = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;

  min-width: 450px;
  .logo {
    svg {
      min-width: 282px;
      min-height: 55px;
    }
  }
  .message {
    color: var(--textBlue);
    margin-bottom: 48px;
    ${hSmallerThan730} {
      margin-bottom: 12px;
    }
  }
  .icon {
    margin-bottom: 48px;
    ${hSmallerThan730} {
      margin-bottom: 32px;
      max-height: 135px;
      svg {
        max-width: 320px;
      }
    }
  }
  .mainContent {
  }
  .lastInput {
    margin-bottom: 34px;
    ${hSmallerThan730} {
      margin-bottom: 24px;
    }
  }
  .info {
    font-size: 0.8em;
    margin: 0;
    span {
      color: var(--textBlue);
    }
  }
  .blue {
    color: var(--textBlue);
    border-radius: 40px !important;
    padding: 5px 15px;
  }
  .inputsWrapper {
    padding: 0vh 0px 2vh 0px;
    max-width: 280px;
    min-width: 280px;
    min-height: 140px;
    display: flex;
    height: 15vh;
    flex-direction: column;
    justify-content: space-evenly;
  }
  .centerHor {
    padding: 1vh 0px;
  }
  .shadedButton {
    background: #f9fafe;
    border-radius: 40px !important;
    padding: 5px 15px;
  }
  label {
    font-size: 14px;
    margin-top: 1px;
  }
  .haveAcountWrapper {
    display: flex;
    align-items: center;
    place-content: center;
    padding-top: 0;
    .grey {
      span {
        color: var(--textGrey);
      }
    }
  }
  .MuiInputLabel-outlined.MuiInputLabel-marginDense {
    margin-top: -1px !important;
  }
  &.popup {
    padding-bottom: 5vh;
    max-height: 525px;
    .icon {
      svg {
        max-width: 335px;
      }
    }
  }
`;
