import React from 'react';
import { MainIconWrapper } from 'components/wrappers/mainIconWrapper';
import { CenterInputsWrapper } from 'components/wrappers/centerInputsWrapper';
import { CenterButtonsWrapper } from 'components/wrappers/centerButtonsWrapper';
import icon from 'images/doneSMSIcon.svg';

import translate from './messages';
import { FormattedMessage } from 'react-intl';
import { Button } from '@material-ui/core';

import { useDispatch } from 'react-redux';
import { Buttons } from 'containers/App/constants';

export default function DoneStep(props: {
  onCancel: Function;
  onSubmit: Function;
  submitIsLoading: boolean;
}) {
  const dispatch = useDispatch();

  const handleSubmit = () => {};
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
          <FormattedMessage {...translate.SMSAuthenticator} />
          <span className="blue p5">
            <FormattedMessage {...translate.Enabled} />
          </span>
        </div>
      </CenterInputsWrapper>
      <CenterButtonsWrapper>
        <Button
          className={Buttons.SimpleRoundButton}
          onClick={handleCancelButton}
        >
          <FormattedMessage {...translate.backToDashboard} />
        </Button>
      </CenterButtonsWrapper>
    </>
  );
}
