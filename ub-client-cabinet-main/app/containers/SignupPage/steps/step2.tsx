import React from 'react';
import styled from 'styles/styled-components';
import translate from '../messages';
import { FormattedMessage } from 'react-intl';
import icon from 'images/signedUpIcon.svg';
import { AppPages, Buttons } from 'containers/App/constants';
import { useDispatch } from 'react-redux';
import { replace } from 'redux-first-history';
import { Button } from '@material-ui/core';
import MainAppIcon from 'images/themedIcons/mainAppIcon';
export default function Step2 () {
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
        <FormattedMessage {...translate.Youraccounthasbeencreated} />
      </div>
      <div className='orange'>
        <FormattedMessage
          {...translate.Pleasecheckyourinboxtoconfirmyouremailaccount}
        />
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
  overflow: auto;
  height: 100vh;
  justify-content: center;
  .titleWrapper {
    margin-bottom: 48px;
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
  .iconWrapper {
    margin-bottom: 48px;
    margin-left: -40px;
  }
  .message {
    margin-bottom: 8px;
    span {
      color: var(--blackText);
    }
  }
  .orange {
    margin-bottom: 24px;
    span {
      color: var(--orange);
    }
  }
  .buttonWrapper {
  }
  .blue {
    color: var(--textBlue);
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
  .MuiInputLabel-outlined.MuiInputLabel-marginDense {
    margin-top: -1px !important;
  }
`;
