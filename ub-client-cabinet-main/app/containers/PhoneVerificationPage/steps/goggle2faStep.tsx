import React, { useState, useEffect } from 'react';
import { MainIconWrapper } from 'components/wrappers/mainIconWrapper';
import { CenterInputsWrapper } from 'components/wrappers/centerInputsWrapper';
import { CenterButtonsWrapper } from 'components/wrappers/centerButtonsWrapper';
import icon from 'images/securityWithMobileIcon.svg';

import translate from './messages';
import { FormattedMessage } from 'react-intl';
import { Button } from '@material-ui/core';

import { useDispatch } from 'react-redux';
import IsLoadingWithText from 'components/isLoadingWithText/isLoadingWithText';

import UBPinInput from 'components/pinInput';
import { Buttons } from 'containers/App/constants';
import { Subscriber, MessageNames } from 'services/message_service';
let code2fa = '';
export default function G2faStep(props: {
  onCancel: Function;
  onSubmit: Function;
  code: string;
}) {
  const [IsLoading, setIsLoading] = useState(false);
  const [CanSubmit, setCanSubmit] = useState(false);
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
  const dispatch = useDispatch();

  const handleSubmit = () => {
    props.onSubmit(code2fa);
  };
  const handleCancelButton = () => {
    props.onCancel();
  };
  return (
    <>
      <MainIconWrapper>
        <img src={icon} alt="" />
      </MainIconWrapper>
      <CenterInputsWrapper style={{ flex: 3 }}>
        <div className="centerHor">
          <FormattedMessage
            {...translate.PleaseopenGoogleAuthenticatorappinyourphoneandenter}
          />
          <span className="bold p5">
            <FormattedMessage {...translate.g2FaCode} />
          </span>
        </div>
        <UBPinInput
          onComplete={value => {
            code2fa = value;
          }}
          onChange={value => {
            if (value.length === 6) {
              setCanSubmit(true);
            }
          }}
        />
      </CenterInputsWrapper>
      <CenterButtonsWrapper>
        <Button
          // style={{ minWidth: '105px' }}
          disabled={!CanSubmit}
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
        {/* <div className="centerHor">
          <FormattedMessage
            {...translate.Forsecurityreasonswerecommendtoenableyour2Fa}
          />
          <Button color="primary" className={Buttons.TransParentRoundButton}>
            <FormattedMessage {...translate.Goto2Fapage} />
          </Button>
        </div> */}
      </CenterButtonsWrapper>
    </>
  );
}
