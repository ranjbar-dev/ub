import React from 'react';
import styled from 'styles/styled-components';
import translate from './messages';
import { FormattedMessage } from 'react-intl';
import icon from 'images/notFounIcon.svg';
import { Button } from '@material-ui/core';
import { AppPages, Buttons } from 'containers/App/constants';
import { useDispatch } from 'react-redux';
import { replace } from 'redux-first-history';
import { FullPageWrapper } from 'components/wrappers/fullPageWrapper';
//import LocaleToggle from 'containers/LocaleToggle';
import MainAppIcon from 'images/themedIcons/mainAppIcon';
export default function NotFound () {
  const dispatch = useDispatch();
  return (
    <FullPageWrapper>
      <div className='head darkTheme'>{/*<LocaleToggle />*/}</div>

      <Wrapper>
        <div className='logo'>
          <MainAppIcon />
        </div>
        <div className='iconWrapper'>
          <img src={icon} alt='' />
        </div>
        <div className='err'>404 Error</div>
        <div className='message'>
          <FormattedMessage {...translate.Pagenotfound} />
        </div>
        <div className='buttonWrapper'>
          <Button
            onClick={() => dispatch(replace(AppPages.LoginPage))}
            className={Buttons.SimpleRoundButton}
          >
            <FormattedMessage {...translate.GoToHome} />
          </Button>
        </div>
      </Wrapper>
    </FullPageWrapper>
  );
}
const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 94vh;
  justify-content: center;
  .titleWrapper {
    span {
      font-size: 39px;
    }
  }
  .logo {
    display: flex;
    align-items: center;
    margin-bottom: 64px;
    @media screen and (max-height: 730px) {
      margin-bottom: 32px;
    }
    svg {
      min-width: 282px;
      min-height: 55px;
    }
  }
  .iconWrapper {
    margin-left: -75px;
    margin-bottom: 48px;
    @media screen and (max-height: 730px) {
      margin-bottom: 24px;
    }
  }
  .message {
    span {
      color: var(--textGrey);
    }
  }
  .err {
    font-size: 37px;
    color: var(--textGrey);
  }
  .buttonWrapper {
    margin-top: 48px;
    @media screen and (max-height: 730px) {
      margin-top: 24px;
    }
  }
  .blue {
    color: var(--textGrey);
    border-radius: 40px !important;
    padding: 5px 15px;
  }

  .centerHor {
  }
  .shadedButton {
    background: #f9fafe;
    border-radius: 40px !important;
    padding: 5px 15px;
  }
`;
