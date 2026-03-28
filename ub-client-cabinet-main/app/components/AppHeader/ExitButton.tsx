import React, { memo, useState, useCallback } from 'react';
import {
  IconButton,
  Tooltip,
  ClickAwayListener,
  Divider,
} from '@material-ui/core';
import PersonIcon from 'images/themedIcons/person';
import translate from 'components/AppHeader/messages';
import ExitIcon from 'images/themedIcons/exitIcon';
import styled from 'styles/styled-components';
import { FormattedMessage } from 'react-intl';
import { LocalStorageKeys } from 'services/constants';
import { toast } from 'components/Customized/react-toastify';
//import { MessageService, MessageNames } from 'services/message_service';
import { useDispatch } from 'react-redux';
import { loggedInAction } from 'containers/App/actions';
import { AppPages } from 'containers/App/constants';
import { replace, push } from 'redux-first-history';
import MenuListIcon from './icons/menuListIcon';
import MenuDepositIcon from './icons/menuDepositIcon';
import MenuWithdrawIcon from './icons/menuWithdrawIcon';
import MenuLogOutIcon from './icons/menuLogOutIcon';
import { cookieConfig, CookieKeys, cookies } from 'services/cookie';
import { FundsPages } from 'containers/FundsPage/constants';
import { MessageNames, MessageService } from 'services/message_service';

interface Props {
  loggedIn: boolean;
  handleLoginButtonClick: () => void;
}

function ExitButton(props: Props) {
  const { loggedIn, handleLoginButtonClick } = props;

  const [IsTooltipOpen, setIsTooltipOpen] = useState(false);
  const dispatch = useDispatch();
  const handleTooltipClose = useCallback(() => {
    setIsTooltipOpen(false);
  }, [loggedIn]);
  const handleExitIconClick = useCallback(() => {
    if (loggedIn === true) {
      setIsTooltipOpen(true);
      return;
    }
    handleLoginButtonClick();
  }, [loggedIn]);
  const handleSignoutClick = () => {
    const theme = localStorage[LocalStorageKeys.Theme];
    const countries = localStorage[LocalStorageKeys.COUNTRIES];
    localStorage.clear();
    localStorage[LocalStorageKeys.COUNTRIES] = countries;
    cookies.remove(CookieKeys.Token, {
      path: cookieConfig().path,
      domain: cookieConfig().domain,
    });
    cookies.remove(CookieKeys.Email, {
      path: cookieConfig().path,
      domain: cookieConfig().domain,
    });
    localStorage[LocalStorageKeys.Theme] = theme;
    dispatch(loggedInAction(false));
    if (!window.location.href.includes(AppPages.TradePage)) {
      dispatch(replace(AppPages.LoginPage));
      return;
    }
    toast.info('successfully logged out');
    MessageService.send({ name: MessageNames.LOGGED_OUT });
  };
  const handleEmailClick = useCallback(() => {
    if (!window.location.href.includes(AppPages.AcountPage)) {
      dispatch(push(AppPages.AcountPage));
    }
  }, []);
  const handleAddressManagementClick = () => {
    dispatch(push(AppPages.AddressManagement));
    requestAnimationFrame(() => {
      MessageService.send({
        name: MessageNames.SET_TAB,
        payload: AppPages.AcountPage,
      });
    });
  };
  const handleDepositClick = () => {
    dispatch(push(AppPages.Funds + '/' + FundsPages.DEPOSIT));
  };
  const handleWithdrawalsClick = () => {
    dispatch(push(AppPages.Funds + '/' + FundsPages.WITHDRAWALS));
  };
  return (
    <ClickAwayListener onClickAway={handleTooltipClose}>
      <Tooltip
        open={IsTooltipOpen}
        interactive
        className='signoutTooltip'
        placement='bottom-end'
        title={
          loggedIn === true ? (
            <Wrapper>
              <div>
                <TooltipElement onClick={handleEmailClick}>
                  {<span>{cookies.get(CookieKeys.Email)}</span>}
                </TooltipElement>

                <Divider />

                <TooltipElement onClick={handleAddressManagementClick}>
                  <MenuListIcon />
                  <FormattedMessage {...translate.AddressManagement} />
                </TooltipElement>

                <Divider />

                <TooltipElement onClick={handleDepositClick}>
                  <MenuDepositIcon />
                  <FormattedMessage {...translate.Deposit} />
                </TooltipElement>

                <TooltipElement onClick={handleWithdrawalsClick}>
                  <MenuWithdrawIcon />
                  <FormattedMessage {...translate.Withdrawal} />
                </TooltipElement>

                <Divider />

                <TooltipElement onClick={handleSignoutClick}>
                  <MenuLogOutIcon />
                  <FormattedMessage {...translate.LogOut} />
                </TooltipElement>
              </div>
            </Wrapper>
          ) : (
            ''
          )
        }
        arrow
      >
        <MainWrapper>
          <IconButton
            onClick={handleExitIconClick}
            className='headerButton'
            size='small'
          >
            {loggedIn === true ? <PersonIcon /> : <ExitIcon />}
          </IconButton>
        </MainWrapper>
      </Tooltip>
    </ClickAwayListener>
  );
}

export default memo(ExitButton);

const MainWrapper = styled.div`
  margin-right: 8px;
`;

const TooltipElement = styled.div`
  cursor: pointer;
  padding: 16px 4px;
  display: flex;
  align-items: center;
  span {
    color: white !important;
    font-size: 12px !important;
  }
`;

const Wrapper = styled.div`
  width: 220px;
  padding: 0 8px;

  hr {
    background-color: rgb(64 69 78) !important;
  }
`;
