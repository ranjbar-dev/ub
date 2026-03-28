import React, { useState } from 'react';
import BlueTitle from '../components/blueTitle';
import { FormattedMessage } from 'react-intl';
import translate from '../messages';
import { Button, Card } from '@material-ui/core';
import { QrCode } from '../types';
import CodeButton from '../components/codeButton';
import { replace } from 'redux-first-history';
import { AppPages, Buttons } from 'containers/App/constants';
import styled from 'styles/styled-components';
import { useDispatch } from 'react-redux';
import enAppstore from 'images/ENappstore.png';
import enGoogleplay from 'images/ENgoogleplay.png';
import UBPinInput from 'components/pinInput';
import { MaxContainer } from 'components/wrappers/maxContainer';
export default function Step1 (props: {
  qrCode: QrCode;
  onNextClick: Function;
}) {
  const qrCode = props.qrCode;
  const dispatch = useDispatch();
  const [CanSubmit, setCanSubmit] = useState(false);
  const [G2faCode, setG2faCode] = useState('');
  const handleNextClick = () => {
    props.onNextClick(G2faCode);
  };
  const handleGooglePlayClick = () => {
    window.open(
      'https://play.google.com/store/apps/details?id=com.google.android.apps.authenticator2&hl=en_GB',
    );
  };
  const handleAppStoreClick = () => {
    window.open(
      'https://apps.apple.com/gb/app/google-authenticator/id388497605',
    );
  };
  return (
    <Wrapper>
      <MaxContainer>
        <div className='downloadButtonsWrapper'>
          <BlueTitle
            number={1}
            title={<FormattedMessage {...translate.DOWNLOADANDINSTALL} />}
          />
          <div className='buttons'>
            <Button onClick={handleAppStoreClick}>
              <img src={enAppstore} alt='' />
            </Button>
            <Button onClick={handleGooglePlayClick}>
              <img src={enGoogleplay} alt='' />
            </Button>
          </div>
        </div>
        <div className='qrCodeWrapper'>
          <BlueTitle
            className='mb12'
            number={2}
            title={
              <FormattedMessage {...translate.ENTERPROVIDEDKEYORScanQRCode} />
            }
          />
          <div className='qrCodeImageWrapper'>
            <img src={qrCode.url} alt='' />
          </div>
          <div className='key'>
            <div className='mb12'>
              <FormattedMessage {...translate.Providedkey} />
            </div>
            <CodeButton qrCode={qrCode} />
          </div>
        </div>
        <div className='codeInputWrapper'>
          <BlueTitle
            number={3}
            title={
              <FormattedMessage
                {...translate.GetAuthenticationCodeFromAppAndEnterHere}
              />
            }
          />

          <UBPinInput
            onEnter={(val: string) => {
              props.onNextClick(val);
            }}
            onComplete={value => {
              setCanSubmit(true);
              setG2faCode(value);
            }}
          />
          <YellowText>
            After enabling 2fa, For security reasons, you will not be able to
            withdraw for next 24 hours
          </YellowText>
        </div>
        <div className='buttonsWrapper'>
          <Button
            disabled={!CanSubmit}
            onClick={handleNextClick}
            className='submitButton'
            variant='contained'
            color='primary'
          >
            <FormattedMessage {...translate.next} />
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
}
const YellowText = styled.p`
  color: var(--orange) !important;
  font-size: 14px;
  margin-top: 12px !important;
  font-weight: 600;
  margin-bottom: -10px;
`;
const Wrapper = styled(Card)`
  width: 100%;
  background: white;
  border-radius: 10px !important;
  box-shadow: none !important;
  height: calc(98vh - 115px);
  min-height: 655px;

  .downloadButtonsWrapper {
    flex: 18;
    text-align: center;
    display: flex;
    flex-direction: column;
    justify-content: center;
    padding-top: 3vh;
    @media screen and (max-height: 750px) {
      flex: 8;
      padding-top: 1vh;
    }
    .MuiButtonBase-root {
      margin: 12px 8px;
      padding: 0 !important;
      img {
        width: 134px;
      }
    }
  }
  .qrCodeWrapper {
    flex: 45;
    text-align: center;
    display: flex;
    flex-direction: column;
    justify-content: space-evenly;
    span {
      font-size: 13px;
    }
    img {
      min-width: 250px;
    }
    .key {
      display: flex;
      flex-direction: column;
      place-content: center;
      span {
        color: var(--textGrey);
      }
    }
    .qrCodeImageWrapper {
      min-height: 250px;
      @media screen and (max-height: 750px) {
        min-height: 200px;
        img {
          min-width: 200px !important;
        }
      }
    }
  }
  .codeInputWrapper {
    flex: 15;
    text-align: center;
    display: flex;
    flex-direction: column;
    justify-content: flex-end;
  }
  .buttonsWrapper {
    flex: 20;
    display: flex;
    flex-direction: column;
    justify-content: center;
    @media screen and (max-height: 750px) {
      justify-content: flex-start;
      margin-top: 15px;
    }
    .submitButton {
      margin-bottom: 1.1vh;
    }
  }
`;
