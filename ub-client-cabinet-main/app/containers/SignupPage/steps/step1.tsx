import React, {
  useState,
  useEffect,
  memo,
  useLayoutEffect,
  useRef,
} from 'react';
import styled from 'styles/styled-components';
import { AppPages, Buttons } from 'containers/App/constants';

import { FormattedMessage } from 'react-intl';
import { useDispatch } from 'react-redux';

import { useInjectSaga } from 'utils/injectSaga';
import { useInjectReducer } from 'utils/injectReducer';
import reducer from '../reducer';
import saga from '../saga';
import translate from '../messages';
import validationTranslate from '../validators/messages';
import InputWithValidator from 'components/inputWithValidator';

import signupMainIcon from 'images/signupMainIcon.svg';
import { replace } from 'redux-first-history';
import { registerAction } from '../actions';
import IsLoadingWithText from 'components/isLoadingWithText/isLoadingWithText';
import {
  Subscriber,
  MessageNames,
  MessageService,
} from 'services/message_service';
import { makeStyles } from '@material-ui/core';
import { Button } from '@material-ui/core';
import MainAppIcon from 'images/themedIcons/mainAppIcon';
import CenterWrapp from 'components/wrappers/centerWrapp';
import { LocalStorageKeys } from 'services/constants';
import { emailValidator } from 'utils/validators/inputValidators/emailValidator';
import {
  passwordContentValidator,
  PasswordErrors,
} from '../validators/passwordContentValidator';
import {
  hBiggerThan700,
  hSmallerThan730,
  hSmallerThan800,
} from 'styles/mediaQueries';
import { isDevelopment } from 'utils/environment';

const materialClasses = makeStyles({
  loadingIndicator: {
    color: 'white',
  },
  loginButton: {
    minHeight: '40px',
    marginBottom: '24px',
  },
});

let recaptcha = '';

enum Ids {
  EmailId = 'signupEmail',
  PasswordId = 'signupPassword',
  ConfirmPasswordId = 'confirmPassword',
}

const Step1 = () => {
  const isDev = isDevelopment;

  useInjectReducer({ key: 'signupPage', reducer: reducer });
  useInjectSaga({ key: 'signupPage', saga: saga });

  useLayoutEffect(() => {
    MessageService.send({
      name: MessageNames.RESET_RECAPTCHA,
    });
    return () => {};
  }, []);

  const [CanSubmit, setCanSubmit] = useState(false);
  const [IsRegistering, setIsRegistering] = useState(false);
  const classes = materialClasses();

  const dispatch = useDispatch();

  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SETLOADING) {
        setIsRegistering(message.payload);
      } else if (message.name === MessageNames.SET_RECAPTCHA) {
        recaptcha = message.payload;
      }
    });

    return () => {
      Subscription.unsubscribe();
      recaptcha = '';
    };
  }, []);
  const handleLoginClickButton = () => {
    dispatch(replace(AppPages.LoginPage));
  };
  const handleSignupClickButton = () => {
    dispatch(
      registerAction({
        email: enteredFields.current[Ids.EmailId],
        password: enteredFields.current[Ids.PasswordId],
        recaptcha: isDev
          ? 'recaptcha'
          : recaptcha
          ? recaptcha
          : localStorage[LocalStorageKeys.RECAPTCHA],
      }),
    );
  };

  const validatedFields = useRef<{
    [Ids.EmailId]: boolean;
    [Ids.PasswordId]: boolean;
    [Ids.ConfirmPasswordId]: boolean;
  }>({
    [Ids.EmailId]: false,
    [Ids.PasswordId]: false,
    [Ids.ConfirmPasswordId]: false,
  });
  const enteredFields = useRef<{
    [Ids.EmailId]: string;
    [Ids.PasswordId]: string;
    [Ids.ConfirmPasswordId]: string;
  }>({
    [Ids.EmailId]: '',
    [Ids.PasswordId]: '',
    [Ids.ConfirmPasswordId]: '',
  });

  const checkForOverallValidation = () => {
    let isAllGood = true;
    if (validatedFields.current[Ids.EmailId] === false) {
      isAllGood = false;
      MessageService.send({
        name: MessageNames.SET_INPUT_ERROR,
        value: Ids.EmailId,
        payload: <FormattedMessage {...validationTranslate.emailIsNotValid} />,
      });
    } else {
      MessageService.send({
        name: MessageNames.SET_INPUT_ERROR,
        value: Ids.EmailId,
        payload: null,
      });
    }
    if (
      validatedFields.current[Ids.PasswordId] === false ||
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
    fieldId: Ids;
    value: string;
    compareWithFieldId?: string;
  }) => {
    enteredFields.current[fieldId] = value;
    if (fieldId === Ids.EmailId) {
      validatedFields.current[Ids.EmailId] = emailValidator({ email: value });
    } else {
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
    }

    checkForOverallValidation();
  };
  const throttleTime = 500;
  return (
    <CenterWrapp>
      <SignupBodyWrapper>
        <div className="logo">
          <MainAppIcon />
        </div>
        <div className="titleWrapper blue">
          <FormattedMessage {...translate.CreateAccount} />
        </div>
        <div className="message">
          <FormattedMessage {...translate.Singuptocontinue} />
        </div>
        <div className="icon">
          <img src={signupMainIcon} alt="" />
        </div>
        <div className="mainContent">
          <div className="inputsWrapper">
            <InputWithValidator
              {...{ throttleTime }}
              autoComplete="false"
              label={<FormattedMessage {...translate.email} />}
              onChange={(email: string) => {
                validator({ fieldId: Ids.EmailId, value: email });
              }}
              uniqueName="signupEmail"
            />
            <InputWithValidator
              {...{ throttleTime }}
              inputType="password"
              isPickable
              label={<FormattedMessage {...translate.password} />}
              onChange={(password: string) => {
                validator({
                  fieldId: Ids.PasswordId,
                  value: password,
                  compareWithFieldId: Ids.ConfirmPasswordId,
                });
              }}
              uniqueName="signupPassword"
            />
            <InputWithValidator
              inputType="password"
              {...{ throttleTime }}
              isPickable
              className="lastInput"
              label={<FormattedMessage {...translate.confirmNewPassword} />}
              onChange={(confirmPassword: string) => {
                validator({
                  fieldId: Ids.ConfirmPasswordId,
                  value: confirmPassword,
                  compareWithFieldId: Ids.PasswordId,
                });
              }}
              uniqueName="confirmPassword"
            />
          </div>
          <Button
            disabled={!CanSubmit}
            variant="contained"
            fullWidth
            onClick={!IsRegistering ? handleSignupClickButton : () => {}}
            color="primary"
            className={classes.loginButton + ' SubmitB'}
          >
            <IsLoadingWithText
              isLoading={IsRegistering}
              text={<FormattedMessage {...translate.CreateAccount} />}
            />
          </Button>

          <div className="centerHor aWrapper">
            <div className="grey">
              <FormattedMessage {...translate.haveAcount} />
            </div>

            <Button
              onClick={handleLoginClickButton}
              className={`button blue ${Buttons.Underlined} `}
            >
              <FormattedMessage {...translate.login} />
            </Button>
          </div>
          {/*<div className="centerHor">
          <Button
            onClick={() => console.log('go home')}
            className={`button blue shadedButton ${Buttons.SimpleRoundButton}`}
          >
            <FormattedMessage {...translate.go_back_home} />
          </Button>
        </div>*/}
        </div>
      </SignupBodyWrapper>
    </CenterWrapp>
  );
};
export default memo(Step1);
const SignupBodyWrapper = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  ${hSmallerThan730} {
    padding-top: 25px;
  }
  overflow: auto;
  ${hBiggerThan700} {
    overflow: hidden;
  }
  .SubmitB {
    ${hSmallerThan730} {
      margin-bottom: 12px;
    }
  }
  .titleWrapper {
    span {
      font-size: 39px;
    }
  }
  .logo {
    display: flex;
    align-items: center;
    svg {
      min-width: 282px;
      min-height: 55px;
    }
  }
  .message {
    color: var(--textGrey);
    margin-bottom: 64px;
    ${hSmallerThan800} {
      margin-bottom: 32px;
    }
    ${hSmallerThan730} {
      margin-bottom: 12px;
    }
  }
  .blue {
    color: var(--textBlue);
    border-radius: 40px !important;
    padding: 5px 15px;
  }
  .icon {
    margin-bottom: 64px;
    ${hSmallerThan800} {
      margin-bottom: 32px;
    }
    ${hSmallerThan730} {
      display: none;
    }
  }
  .lastInput {
    margin-bottom: 48px;
    ${hSmallerThan730} {
      margin-bottom: 24px;
    }
    ${hSmallerThan800} {
      margin-bottom: 16px;
    }
  }
  .inputsWrapper {
    max-width: 280px;
    min-width: 280px;
    display: flex;
    flex-direction: column;
    justify-content: space-evenly;
  }
  .centerHor {
    .grey {
      span {
        color: var(--textGrey);
      }
    }
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
  .loadingCircle {
    top: 10px !important;
  }
  .aWrapper {
    padding-top: 0;
  }
  .MuiInputLabel-outlined.MuiInputLabel-marginDense {
    margin-top: -1px !important;
  }
`;
