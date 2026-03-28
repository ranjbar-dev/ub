import React, { useState, useEffect } from 'react';

import { replace } from 'redux-first-history';
import { changePasswordAction, isChangingPasswordAction } from '../actions';
import InputWithValidator from 'components/inputWithValidator';
import { NewPasswordValidator } from '../validators/newPasswordValidator';
import { OldPasswordValidator } from '../validators/oldPasswordValidator';
import { ConfirmNewPasswordValidator } from '../validators/confirmNewPasswordValidator';

import { FormattedMessage } from 'react-intl';
import { useSelector, useDispatch } from 'react-redux';
import { createStructuredSelector } from 'reselect';

import { useInjectSaga } from 'utils/injectSaga';
import { useInjectReducer } from 'utils/injectReducer';
import makeSelectChangePassword, {
  makeSelectTheme,
  makeSelectIsLoading,
} from '../selectors';
import reducer from '../reducer';
import saga from '../saga';
import translate from '../messages';

import { AppPages, Buttons } from 'containers/App/constants';

import { ChangePasswordModel } from '../types';
import ChangePasswordMainIcon from 'images/themedIcons/changePasswordMainIcon';
import { makeStyles, CircularProgress, Button } from '@material-ui/core';
import { MaxContainer } from 'components/wrappers/maxContainer';
import { StepsWrapper } from '../stepsWrapper';
import { TwofaAndEmailVerificationPopup } from 'components/twoFaAndVerificationPopup/twofaAndEmailVerificationPopup';

const materialClasses = makeStyles({
  loadingIndicator: {
    color: 'white',
  },
});
let fields = {
  oldPassword: { isValid: false, value: '' },
  newPassword: { isValid: false, value: '' },
  confirmNewPassword: { isValid: false, value: '' },
};
const stateSelector = createStructuredSelector({
  changePassword: makeSelectChangePassword(),
  theme: makeSelectTheme(),
  isLoading: makeSelectIsLoading(),
});

export default function Step1 () {
  useInjectReducer({ key: 'changePassword', reducer: reducer });
  useInjectSaga({ key: 'changePassword', saga: saga });
  const [CanSubmit, setCanSubmit] = useState(false);

  const isFieldValid = (properties: {
    fieldName: string;
    isValid: boolean;
    value: string;
  }) => {
    const { fieldName, isValid, value } = properties;
    fields[fieldName].isValid = isValid;
    fields[fieldName].value = value;
    if (fieldName === 'confirmNewPassword' && isValid === true) {
      if (fields['newPassword'].value === value) {
        fields['newPassword'].isValid = true;
      }
    }
    if (fieldName === 'newPassword' && isValid === true) {
      if (fields['confirmNewPassword'].value === value) {
        fields['confirmNewPassword'].isValid = true;
      }
    }
    if (
      fields['oldPassword'].isValid === true &&
      fields['newPassword'].isValid === true &&
      fields['confirmNewPassword'].isValid === true
    ) {
      setCanSubmit(true);
    } else {
      setCanSubmit(false);
    }
  };

  useEffect(() => {
    return () => {
      fields = {
        oldPassword: { isValid: false, value: '' },
        newPassword: { isValid: false, value: '' },
        confirmNewPassword: { isValid: false, value: '' },
      };
    };
  }, []);
  const classes = materialClasses();
  //custom hooks

  const { isLoading } = useSelector(stateSelector);
  const dispatch = useDispatch();

  //button actions
  const handleSubmit = (additionalData?: any) => {
    const data: ChangePasswordModel = {
      old_password: fields.oldPassword.value,
      new_password: fields.newPassword.value,
      confirmed: fields.confirmNewPassword.value,
      ...additionalData,
    };
    dispatch(changePasswordAction(data));
  };
  const handleCancelButton = () => {
    dispatch(replace(AppPages.AcountPage));
  };

  const isChangingPassword = () => {
    if (!isLoading) {
      return (
        <span>
          <FormattedMessage {...translate.changePassword} />
        </span>
      );
    }
    return isLoading ? (
      <CircularProgress size={14} className={classes.loadingIndicator} />
    ) : (
      <span>
        <FormattedMessage {...translate.changePassword} />
      </span>
    );
  };

  return (
    <StepsWrapper className='flexCenter'>
      <MaxContainer className='mh1'>
        <TwofaAndEmailVerificationPopup
          onSubmit={handleSubmit}
          onClose={() => dispatch(isChangingPasswordAction(false))}
        />
        <div className='mainIcon middleIcon'>
          <ChangePasswordMainIcon />
        </div>
        <div className='inputsWrapper'>
          <div className='inputWrapper'>
            <InputWithValidator
              inputType='password'
              isPickable={true}
              throttleTime={500}
              label={<FormattedMessage {...translate.oldPassword} />}
              onChange={(oldPassword: string) => {
                isFieldValid({
                  fieldName: 'oldPassword',
                  isValid: OldPasswordValidator({
                    uniqueInputId: 'oldPassword',
                    value: oldPassword,
                  }),
                  value: oldPassword,
                });
              }}
              uniqueName='oldPassword'
            />
          </div>
          <div className='inputWrapper'>
            <InputWithValidator
              inputType='password'
              isPickable={true}
              throttleTime={500}
              label={<FormattedMessage {...translate.newPassword} />}
              onChange={(newPassword: string) => {
                isFieldValid({
                  fieldName: 'newPassword',
                  isValid: NewPasswordValidator({
                    uniqueInputId: 'newPassword',
                    value: newPassword,
                    compareTo: fields['confirmNewPassword'].value,
                    compareId: 'confirmNewPassword',
                  }),
                  value: newPassword,
                });
              }}
              uniqueName='newPassword'
            />
          </div>
          <div className='inputWrapper last'>
            <InputWithValidator
              inputType='password'
              isPickable={true}
              throttleTime={500}
              label={<FormattedMessage {...translate.confirmNewPassword} />}
              onChange={(repeatPassword: string) => {
                isFieldValid({
                  fieldName: 'confirmNewPassword',
                  isValid: NewPasswordValidator({
                    uniqueInputId: 'confirmNewPassword',
                    value: repeatPassword,
                    compareTo: fields['newPassword'].value,
                    compareId: 'newPassword',
                  }),

                  value: repeatPassword,
                });
              }}
              uniqueName='confirmNewPassword'
            />
          </div>
        </div>
        <div className='buttonsWrapper'>
          <Button
            disabled={!CanSubmit}
            onClick={() => (!isLoading ? handleSubmit() : null)}
            className='ubButton'
            fullWidth
            variant='contained'
            color='primary'
          >
            {isChangingPassword()}
          </Button>
          <Button className={Buttons.CancelButton} onClick={handleCancelButton}>
            <FormattedMessage {...translate.cancel} />
          </Button>
        </div>
      </MaxContainer>
    </StepsWrapper>
  );
}
