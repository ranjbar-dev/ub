/**
 *
 * AppHeader
 *
 */
import React, { memo, useEffect, useMemo, useCallback } from 'react';
import { Paper, Tab, Tabs, makeStyles } from '@material-ui/core';

// import styled from 'styles/styled-components';

import { FormattedMessage } from 'react-intl';
import translate from './messages';

import TradeIcon from 'images/selectableImages/tradeIcon';
import FundsIcon from 'images/selectableImages/fundsIcon';
import OrdersIcon from 'images/selectableImages/ordersIcon';
import AcountIcon from 'images/selectableImages/acountIcon';
import styled from 'styles/styled-components';
import { useDispatch, useSelector } from 'react-redux';
import { push, replace } from 'redux-first-history';
import { AppPages } from 'containers/App/constants';
// import LocaleToggle from 'containers/LocaleToggle';
import ThemeToggler from './themeToggler';
import MainAppIcon from 'images/themedIcons/mainAppIcon';

import {
  Subscriber,
  MessageNames,
  MessageService,
} from 'services/message_service';
import LayoutSelect from './layoutSelect';
import { createStructuredSelector } from 'reselect';
import { makeSelectLoggedIn } from './selectors';
import ExitButton from './ExitButton';
import { LandingPageAddress } from 'services/constants';
import { Link } from 'react-router-dom';

const stateSelector = createStructuredSelector({
  loggedIn: makeSelectLoggedIn(),
});

const useStyles = makeStyles({
  root: {
    flexGrow: 1,
  },
});

const HeaderTab = (data: { icon: any; title: any; to }) => {
  return (
    <Link
      style={{ textDecoration: 'none' }}
      to={data.to}
      onClick={(e) => {
        e.preventDefault();
      }}
    >
      <span>
        {data.icon}
        {data.title}
      </span>
    </Link>
  );
};

const AppHeader = () => {
  const classes = useStyles();
  const [activeIndex, setactiveIndex] = React.useState(3);
  const { loggedIn } = useSelector(stateSelector);

  const dispatch = useDispatch();
  const openPage = useCallback((index: number) => {
    switch (index) {
      case 3:
        dispatch(push(AppPages.AcountPage));
        break;
      case 2:
        dispatch(push(AppPages.Orders));
        break;
      case 1:
        dispatch(push(AppPages.Funds));
        break;
      case 0:
        dispatch(push(AppPages.TradePage));
        break;

      default:
        dispatch(push(AppPages.HomePage));
        break;
    }
  }, []);

  const handleChange = useCallback((event, newactiveIndex) => {
    setactiveIndex(newactiveIndex);
    setTimeout(() => {
      openPage(newactiveIndex);
    }, 130);
  }, []);

  const handleExit = useCallback(() => {
    if (
      activeIndex === 0 ||
      window.location.href.includes(AppPages.TradePage)
    ) {
      MessageService.send({ name: MessageNames.OPEN_LOGIN_POPUP });
      MessageService.send({ name: MessageNames.RESET_RECAPTCHA });
    } else {
      dispatch(replace(AppPages.LoginPage));
    }
  }, [activeIndex]);

  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SET_TAB) {
        if (
          message.payload.includes(AppPages.AcountPage) &&
          activeIndex !== 3
        ) {
          setactiveIndex(3);
        } else if (
          message.payload.includes(AppPages.Orders) &&
          activeIndex !== 2
        ) {
          setactiveIndex(2);
        } else if (
          message.payload.includes(AppPages.Funds) &&
          activeIndex !== 1
        ) {
          setactiveIndex(1);
        } else if (
          message.payload.includes(AppPages.TradePage) &&
          activeIndex !== 0
        ) {
          setactiveIndex(0);
        }
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, [activeIndex]);

  return useMemo(
    () => (
      <>
        <TopTabsWrapper className="darkTheme">
          <Paper
            className={classes.root}
            style={{ borderRadius: 0 }}
            elevation={0}
          >
            <span className="mainIcon">
              <a href={LandingPageAddress}>
                <MainAppIcon />
              </a>
            </span>
            <div className="mainBarWrapper">
              <div className="tabsWrpr">
                <Tabs
                  value={activeIndex}
                  onChange={handleChange}
                  indicatorColor="primary"
                  textColor="primary"
                >
                  <Tab
                    className="first"
                    disableRipple
                    label={
                      <HeaderTab
                        to={AppPages.TradePage}
                        icon={<TradeIcon isSelected={activeIndex === 0} />}
                        title={<FormattedMessage {...translate.trade} />}
                      />
                    }
                  />
                  <Tab
                    disableRipple
                    label={
                      <HeaderTab
                        to={AppPages.Funds}
                        icon={<FundsIcon isSelected={activeIndex === 1} />}
                        title={<FormattedMessage {...translate.funds} />}
                      />
                    }
                  />
                  <Tab
                    disableRipple
                    label={
                      <HeaderTab
                        to={AppPages.Orders}
                        icon={<OrdersIcon isSelected={activeIndex === 2} />}
                        title={<FormattedMessage {...translate.orders} />}
                      />
                    }
                  />
                  <Tab
                    disableRipple
                    label={
                      <HeaderTab
                        to={AppPages.AcountPage}
                        icon={<AcountIcon isSelected={activeIndex === 3} />}
                        title={<FormattedMessage {...translate.acount} />}
                      />
                    }
                  />
                </Tabs>
              </div>
              <div className="headerIconsWrapper">
                <ExitButton
                  loggedIn={loggedIn}
                  handleLoginButtonClick={handleExit}
                />
                {activeIndex === 0 && <LayoutSelect />}
                <ThemeToggler />
                {/*<IconButton className='headerButton' size='small'>
                  <BellIcon />
                </IconButton>*/}

                {/* <LocaleToggle /> */}
              </div>
            </div>
          </Paper>
        </TopTabsWrapper>
      </>
    ),
    [activeIndex, loggedIn],
  );
};
export default memo(AppHeader);
const TopTabsWrapper = styled.div`
  width: 100vw;
  position: fixed;
  top: 0;
  z-index: 2;
  .MuiButtonBase-root {
    /* min-width: 130px !important; */
    span {
      padding: 0;
      max-width: fit-content;
      padding-right: 1px !important;
      min-width: max-content;
      svg {
        margin-right: 4px;
      }
    }
  }
  .MuiTab-root {
    min-width: unset !important;
    max-width: fit-content;
    padding: 0;
    margin: 0 24px;
    min-height: 60px;
    &.first {
      margin-left: 48px;
    }
    span {
      color: var(--textGrey);
      line-height: 20px;
    }
  }
  .mainIcon {
    position: relative;
    float: left;
    margin-top: 18px;
  }
  span {
    padding: 0 10px;
  }
  .mainBarWrapper {
    display: flex;
    justify-content: space-between;
    align-items: center;
    .headerButton {
      max-width: 40px !important;
      max-height: 40px;
      min-width: 40px !important;
    }
    max-height: 60px;
    min-height: 60px;
  }
  span {
    text-transform: none;
  }
  .headerIconsWrapper {
    display: flex;
    height: 55px;
    align-items: center;
    svg {
      min-width: 25px;
    }
  }
  .MuiTabs-indicator {
    transition: left 60ms linear;
    background: var(--textBlue) !important;
  }
  .Mui-selected {
    span {
      font-weight: 600;
      padding-right: 1px !important;
      min-width: max-content;
      color: var(--textBlue);
    }
    path {
      fill: var(--textBlue);
    }
  }
  .layoutIcon {
    path {
      fill: var(--blackText);
    }
  }
  .tabsWrpr {
    min-width: 600px;
  }
`;
