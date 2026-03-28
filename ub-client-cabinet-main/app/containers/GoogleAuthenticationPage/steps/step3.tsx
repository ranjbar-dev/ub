import React, { useEffect } from 'react';
import { MainIconWrapper } from 'components/wrappers/mainIconWrapper';
import { CenterInputsWrapper } from 'components/wrappers/centerInputsWrapper';
import { CenterButtonsWrapper } from 'components/wrappers/centerButtonsWrapper';
import icon from 'images/doneGoogleAuthentication.svg';

import translate from '../messages';
import { FormattedMessage } from 'react-intl';
import { Button, Card } from '@material-ui/core';

import { useDispatch } from 'react-redux';
import styled from 'styles/styled-components';
import { Buttons, AppPages } from 'containers/App/constants';
import { replace } from 'redux-first-history';
import { getUserDataAction } from 'containers/AcountPage/actions';
import { MaxContainer } from 'components/wrappers/maxContainer';

export default function Step3 (props: {}) {
  const dispatch = useDispatch();
  const handleCancelButton = () => {
    dispatch(replace(AppPages.AcountPage));
  };
  useEffect(() => {
    dispatch(getUserDataAction());
    return () => {};
  }, []);
  return (
    <Wrapper>
      <MaxContainer>
        <CenterInputsWrapper style={{ flex: 3, padding: 0 }}>
          <div className='centerHor'>
            <FormattedMessage {...translate.done} />
          </div>
        </CenterInputsWrapper>
        <MainIconWrapper className='noPadding'>
          <img src={icon} alt='' />
        </MainIconWrapper>
        <div className='open2fa'>
          <FormattedMessage
            {...translate.OpenTowFactorAuthenticationWithGoogle}
          />
          <span className='blue'>
            <FormattedMessage {...translate.enabled} />
          </span>
        </div>{' '}
        <YellowText>
          For security reasons, you will not be able to withdraw for next 24
          hours
        </YellowText>
        <CenterButtonsWrapper style={{ flex: 5, padding: 0 }}>
          <Button
            className={Buttons.SimpleRoundButton}
            onClick={handleCancelButton}
          >
            <FormattedMessage {...translate.backToDashboard} />
          </Button>
        </CenterButtonsWrapper>
      </MaxContainer>
    </Wrapper>
  );
}
const YellowText = styled.p`
  color: var(--orange) !important;
  font-size: 14px;
  margin-top: -12px !important;
  font-weight: 600;
`;
const Wrapper = styled(Card)`
  width: 100%;
  background: white;
  border-radius: 10px !important;
  box-shadow: none !important;
  height: calc(98vh - 115px);
  min-height: 740px;
  padding: 77px;

  .noPadding {
    padding: 0 !important;
    flex: 10 !important;
  }
  .centerHor {
    span {
      font-size: 25px;
      font-weight: 700;
      color: var(--textBlue);
    }
  }
  .open2fa {
    flex: 4;
    display: flex;
    align-items: center;
    span {
      color: var(--blackText);
    }
    .blue {
      span {
        color: var(--textBlue);
      }
    }
  }
  .blue {
    padding: 0 6px;
  }
`;
