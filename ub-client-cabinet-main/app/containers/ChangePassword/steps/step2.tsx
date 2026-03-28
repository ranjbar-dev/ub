import React from 'react';
import { MainIconWrapper } from 'components/wrappers/mainIconWrapper';
import { CenterButtonsWrapper } from 'components/wrappers/centerButtonsWrapper';
import greenShield from 'images/greenShieldIcon.svg';

import translate from '../messages';
import { FormattedMessage } from 'react-intl';
import { Button, Card } from '@material-ui/core';

import { useDispatch } from 'react-redux';
import styled from 'styles/styled-components';
import { Buttons, AppPages } from 'containers/App/constants';
import { replace } from 'redux-first-history';
import { MaxContainer } from 'components/wrappers/maxContainer';
import { PasswordChangedIcon } from 'images/themedIcons/passwordChanged';

export default function Step2 (props: {}) {
  const dispatch = useDispatch();
  const handleCancelButton = () => {
    dispatch(replace(AppPages.AcountPage));
  };
  return (
    <Wrapper>
      <MaxContainer>
        <MainIconWrapper className='noPadding'>
          {/* <img src={icon} alt='' /> */}
          <PasswordChangedIcon />
        </MainIconWrapper>
        <div className='centerHor'>
          <span className='icon'>
            <img src={greenShield} alt='' />
          </span>{' '}
          <FormattedMessage {...translate.Yourpasswordhasbeenchanged} />
        </div>
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
const Wrapper = styled(Card)`
  background: white;
  border-radius: 10px !important;
  box-shadow: none !important;
  height: calc(98vh - 115px);
  /*min-height: 735px;*/
  padding: 77px;
  display: flex;
  flex-direction: column;
  align-items: center;
  .noPadding {
    flex: 10 !important;
    display: flex;
    align-items: flex-end;
    padding-bottom: 4vh;
  }
  .centerHor {
    span {
      color: var(--greenText);
    }
    flex: 3;
  }
  .open2fa {
    flex: 4;
    display: flex;
    align-items: center;
  }
  .blue {
    padding: 0 6px;
  }
`;
