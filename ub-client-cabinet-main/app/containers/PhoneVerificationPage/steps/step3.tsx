import React, { useState, useEffect } from 'react';
import { MainIconWrapper } from 'components/wrappers/mainIconWrapper';
import { CenterInputsWrapper } from 'components/wrappers/centerInputsWrapper';
import { CenterButtonsWrapper } from 'components/wrappers/centerButtonsWrapper';

import translate from './messages';
import { FormattedMessage } from 'react-intl';
import { Button } from '@material-ui/core';

import { useDispatch } from 'react-redux';
import IsLoadingWithText from 'components/isLoadingWithText/isLoadingWithText';
import { Buttons } from 'containers/App/constants';
import { Subscriber, MessageNames } from 'services/message_service';
import ChangePasswordMainIcon from 'images/themedIcons/changePasswordMainIcon';
import InputWithValidator from 'components/inputWithValidator';
let password = '';
export default function Step3(props: {
  onCancel: Function;
  onSubmit: Function;
  on2fa: Function;
  userData: any;
  submitIsLoading: boolean;
}) {
  const dispatch = useDispatch();
  const [IsLoading, setIsLoading] = useState(false);
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SETLOADING) {
        setIsLoading(message.payload);
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);

  const handleSubmit = (e) => {
    props.onSubmit(password);
  };
  const handleCancelButton = () => {
    props.onCancel();
  };

  return (
    <>
      <MainIconWrapper className="fl fl9 minimized">
        <ChangePasswordMainIcon />
      </MainIconWrapper>
      <CenterInputsWrapper className="noPadd pt3 start" style={{ flex: 4 }}>
        <div className="flexSpacer1  mh40"></div>
        <div className="centerHor fl1">
          <FormattedMessage {...translate.Pleaseenteryouraccountpassword} />
        </div>
        <div className="inputWrapper">
          <InputWithValidator
            throttleTime={0}
            inputType="password"
            isPickable={true}
            autoFocus={true}
            label={<FormattedMessage {...translate.password} />}
            onChange={(pass: string) => {
              password = pass;
              //isFieldValid({
              //  fieldName: 'signupPassword',
              //  isValid: PasswordValidator({
              //    uniqueInputId: 'signupPassword',
              //    value: password,
              //  }),
              //  value: password,
              //});
            }}
            uniqueName="signupPassword"
          />
          {/*<TextField
            fullWidth
            type="password"
            autoComplete="off"
            name="password"
            onChange={(e) => {
              password = e.target.value;
            }}
            label={<FormattedMessage {...translate.Password} />}
            variant="outlined"
            margin="dense"
          />*/}
        </div>
      </CenterInputsWrapper>
      <CenterButtonsWrapper style={{ flex: 6 }}>
        <Button
          // style={{ minWidth: '105px' }}
          // disabled={hasError || SelectedIndex === -1}
          onClick={handleSubmit}
          className="ubButton"
          variant="contained"
          color="primary"
        >
          <IsLoadingWithText
            isLoading={IsLoading}
            text={<FormattedMessage {...translate.submit} />}
          />
        </Button>
        <Button className={Buttons.CancelButton} onClick={handleCancelButton}>
          <FormattedMessage {...translate.cancel} />
        </Button>
        {props.userData.google2faEnabled === false && (
          <div className="centerHor">
            <div className="black">
              <FormattedMessage
                {...translate.Forsecurityreasonswerecommendtoenableyour2Fa}
              />
            </div>

            <Button
              onClick={() => props.on2fa()}
              color="primary"
              className={Buttons.TransParentRoundButton}
            >
              <FormattedMessage {...translate.Goto2Fapage} />
            </Button>
          </div>
        )}
      </CenterButtonsWrapper>
    </>
  );
}
