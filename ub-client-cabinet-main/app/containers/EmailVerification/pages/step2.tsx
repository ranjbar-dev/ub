import React from 'react';
import styled from 'styles/styled-components';
import translate from '../messages';
import { FormattedMessage } from 'react-intl';
import icon from 'images/activatedIcon.svg';
import { AppPages, Buttons } from 'containers/App/constants';
import { useDispatch } from 'react-redux';
import { replace } from 'redux-first-history';
import { Button } from '@material-ui/core';
import MainAppIcon from 'images/themedIcons/mainAppIcon';
export default function EmailVerified () {
  const dispatch = useDispatch();
  return (
    <Wrapper>
      <div className='logo'>
        <MainAppIcon />
      </div>
      <div className='titleWrapper blue'>
        <FormattedMessage {...translate.CreateAccount} />
      </div>
      <div className='iconWrapper'>
        <img src={icon} alt='' />
      </div>
      <div className='message'>
        <FormattedMessage {...translate.Youraccounthasbeenactivated} />
      </div>

      <div className='buttonWrapper'>
        <Button
          onClick={() => dispatch(replace(AppPages.LoginPage))}
          className={`button blue shadedButton ${Buttons.SimpleRoundButton}`}
        >
          <FormattedMessage {...translate.GoToLoginPage} />
        </Button>
      </div>
    </Wrapper>
  );
}
const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  height: 100vh;
  padding-top: 15vh;
  .titleWrapper {
    flex: 10;
    span {
      font-size: 39px;
    }
  }
  .logo {
    flex: 7;
    display: flex;
    align-items: center;
    svg {
      min-width: 350px;
      min-height: 55px;
    }
  }
  .iconWrapper {
    flex: 20;
  }
  .message {
    flex: 3;
    span {
      color: var(--greenText);
    }
  }
  .orange {
    flex: 5;
    span {
      color: var(--orange);
    }
  }
  .buttonWrapper {
    flex: 20;
  }
  .blue {
    color: var(--textBlue);
    border-radius: 40px !important;
    padding: 5px 15px;
  }

  .centerHor {
    padding: 1vh 0px;
  }
  .shadedButton {
    background: #f9fafe;
    border-radius: 40px !important;
    padding: 5px 15px;
  }
`;
