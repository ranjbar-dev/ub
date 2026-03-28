import React, { useState, useEffect } from 'react';
import { FormattedMessage, injectIntl } from 'react-intl';
import translate from '../messages';
import { Button, Card } from '@material-ui/core';

import { replace } from 'redux-first-history';
import { AppPages, Buttons } from 'containers/App/constants';
import styled from 'styles/styled-components';
import { useDispatch } from 'react-redux';
import ChangePasswordMainIcon from 'images/themedIcons/changePasswordMainIcon';

import UBPinInput from 'components/pinInput';
import { MainIconWrapper } from 'components/wrappers/mainIconWrapper';
import InputWithValidator from 'components/inputWithValidator';
import { toggle2FaAction } from '../actions';
import PopupModal from 'components/materialModal/modal';
import warningIcon from 'images/warning.svg';
import { Subscriber, MessageNames } from 'services/message_service';
import { MaxContainer } from 'components/wrappers/maxContainer';
const DisableStep = (props: { intl }) => {
  const intl = props.intl;
  const dispatch = useDispatch();
  const [Password, setPassword] = useState('');
  const [G2faCode, setG2faCode] = useState('');
  const [IsLoading, setIsLoading] = useState(false);
  const [DataToSend, setDataToSend]: [any, any] = useState({});
  const [IsAlertOpen, setIsAlertOpen] = useState(false);

  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SETLOADING) {
        setIsLoading(message.payload);
        if (message.payload === false) {
          setIsAlertOpen(false);
        }
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);

  const handleNextClick = () => {
    const dataTosend = {
      code: G2faCode,
      password: Password,
      setEnable: false,
    };
    setDataToSend(dataTosend);
    setIsAlertOpen(true);
  };
  const sendData = () => {
    dispatch(toggle2FaAction(DataToSend));
  };
  const onPasswordChange = (e: string) => {
    setPassword(e);
  };
  return (
    <Wrapper>
      <MaxContainer>
        <PopupModal
          isOpen={IsAlertOpen}
          onClose={() => {
            setIsAlertOpen(false);
          }}
        >
          <div className='alertWrapper alertConfirmWrapper'>
            <span>
              {intl.formatMessage({
                id: 'containers.GoogleAuthenticationPage.Areyousure',
                defaultMessage: 'ET.Label',
              })}
              <span className='red'>
                {' '}
                {intl.formatMessage({
                  id:
                    'containers.GoogleAuthenticationPage.DisableGoogleAuthenticator',
                  defaultMessage: 'ET.Label',
                })}
              </span>
              {intl.formatMessage({
                id: 'containers.GoogleAuthenticationPage.question',
                defaultMessage: 'ET.Label',
              })}
            </span>
            <YellowText>
              For security reasons, you will not be able to withdraw for next 24
              hours
            </YellowText>
          </div>
          <div className='alertButtonsWrapper'>
            <Button
              onClick={() => {
                setIsAlertOpen(false);
              }}
            >
              {intl.formatMessage({
                id: 'containers.GoogleAuthenticationPage.cancel',
                defaultMessage: 'ET.Label',
              })}
            </Button>
            <div className='separator'></div>
            <Button onClick={sendData}>
              {intl.formatMessage({
                id: 'containers.GoogleAuthenticationPage.allow',
                defaultMessage: 'ET.Label',
              })}
            </Button>
          </div>
        </PopupModal>
        <div className='warningFlexWrapper'>
          <div className='warningWrapper'>
            <div className='iconWrapper'>
              <img src={warningIcon} alt='' />
            </div>
            <div>
              <FormattedMessage {...translate.g2faWraning} />
            </div>
          </div>
        </div>
        <MainIconWrapper className='noPadding'>
          <ChangePasswordMainIcon width={313} />
        </MainIconWrapper>
        <div className='passwordInputWrapper'>
          <div className='message'>
            <FormattedMessage {...translate.Pleaseenteryouraccountpassword} />
          </div>
          <div className='passwordWrapper'>
            <InputWithValidator
              inputType='password'
              uniqueName='2faInput'
              label={<FormattedMessage {...translate.password} />}
              isPickable={true}
              onChange={onPasswordChange}
            />
          </div>
        </div>
        <div className='codeInputWrapper'>
          <div className='message'>
            <FormattedMessage
              {...translate.PleaseopenGoogleAuthenticatorappinyourphoneandenter2FaCode}
            />
          </div>
          <UBPinInput
            onComplete={value => {
              setG2faCode(value);
            }}
          />
        </div>
        <div className='buttonsWrapper'>
          <Button
            disabled={Password.length < 8 || G2faCode === ''}
            onClick={handleNextClick}
            variant='contained'
            color='primary'
          >
            <FormattedMessage {...translate.submit} />
          </Button>
          <Button
            onClick={() => {
              dispatch(replace(AppPages.AcountPage));
            }}
            className={Buttons.CancelButton}
            color='primary'
          >
            <FormattedMessage {...translate.cancel} />
          </Button>
        </div>
      </MaxContainer>
    </Wrapper>
  );
};
export default injectIntl(DisableStep);
const YellowText = styled.p`
  color: var(--orange) !important;
  font-size: 14px;
  margin-top: 40px !important;
  font-weight: 600;
`;
const Wrapper = styled(Card)`
  background: white;
  border-radius: 10px !important;
  box-shadow: none !important;
  height: calc(98vh - 115px);
  padding: 4vh 77px;
  min-height: 654px;
  overflow-y: auto;

  .codeInputWrapper,
  .passwordInputWrapper {
    text-align: center;
    .message {
      span {
        color: var(--textGrey);
      }
    }
  }
  .passwordInputWrapper {
    display: flex;
    flex-direction: column;
    justify-content: center;
    flex: 20;
    margin-bottom: 1vh;
  }
  .passwordWrapper {
    width: 276px;
  }
  .buttonsWrapper {
    flex: 15;
    display: flex;
    flex-direction: column;
    justify-content: space-evenly;
  }
  .warningFlexWrapper {
    flex: 20;
    display: flex;
    flex-direction: column;
    justify-content: center;
    .warningWrapper {
      width: 426px;
      background: var(--oddRows);
      padding: 12px;

      border-radius: 7px;
      display: flex;
      align-items: center;
      span {
        color: var(--orange);
      }
      .iconWrapper {
        width: 90px;
        display: flex;
        place-content: center;
        margin-right: 15px;
        border-right: 1px solid var(--lightGrey);
        min-height: 45px;
        padding-right: 6px;
      }
    }
  }
`;
