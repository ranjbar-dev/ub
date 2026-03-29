/**
 *
 * App
 *
 * This component is the skeleton around the actual pages, and should only
 * contain code that should be seen on all pages. (e.g. navigation bar)
 */

import React from 'react';
import { Helmet } from 'react-helmet';
import styled from 'styles/styled-components';
import { Switch, Route } from 'react-router-dom';
// import 'styles/ag-grid.min.css';
// import 'styles/ag-theme-balham.min.css';
import HomePage from 'containers/HomePage';
import LoginPage from 'containers/LoginPage/Loadable';
import SignupPage from 'containers/SignupPage/Loadable';
import AcountPage from 'containers/AcountPage/Loadable';
import ChangePasswordPage from 'containers/ChangePassword/Loadable';
import NotFoundPage from 'containers/NotFoundPage/Loadable';
import PhoneVerificationPage from 'containers/PhoneVerificationPage/Loadable';
import AddressManagementPage from 'containers/AddressManagementPage/Loadable';
import OrdersPage from 'containers/OrdersPage/Loadable';
import TradePage from 'containers/TradePage/Loadable';
import FundsPage from 'containers/FundsPage/Loadable';
import DocumentVerificationPage from 'containers/DocumentVerificationPage/Loadable';
import ChangeUserInfoPage from 'containers/ChangeUserInfoPage/Loadable';
import GoogleAuthenticationPage from 'containers/GoogleAuthenticationPage/Loadable';
import EmailVerificationPage from 'containers/EmailVerification/Loadable';
import UpdatePasswordPage from 'containers/UpdatePasswordPage/Loadable';

import GlobalStyle from '../../global-styles';

import { create } from 'jss';
import rtl from 'jss-rtl';

import { createStructuredSelector } from 'reselect';
import { useSelector, useDispatch } from 'react-redux';
import {
  makeSelectAppState,
  makeSelectLoggedIn,
  makeSelectLanguage,
  makeSelectLocation,
  makeSelectTheme,
} from './selectors';
import AppHeader from 'components/AppHeader';
import { AppPages } from './constants';
import {
  jssPreset,
  ThemeProvider,
  createMuiTheme,
  StylesProvider,
} from '@material-ui/core';

import Footer from 'components/Footer';
import { useMemo, useRef, useEffect } from 'react';
import { MessageService, MessageNames } from 'services/message_service';
import { replace } from 'redux-first-history';
import RecaptchaContainer from 'containers/RecapchaContainer';
import { loggedInAction } from './actions';
import { CookieKeys, cookies } from 'services/cookie';
import PrivateRoute from './PrivateRoute';
import RegisteredToastContainer from 'components/registeredToastContainer/registeredToastContainer';
import { useResizeAndThemeHook } from './hooks/resizeAndThemeHook';
import { useConnectToCentrifugoAuth } from './hooks/connectToMqtt2';
import { useCheckForAuthErrors } from './hooks/checkForAuthErrors';
// import { OnlineStatusProvider } from 'hooks/onlineStatusHook/provider';
import { useGetInitialData } from './hooks/useGetInitialData';
import { initGA, logPageView } from '../../utils/analytics';
import history from 'utils/history';
import ContactUsPage from 'containers/ContactUsPage/Loadable';

//import Tawk from 'components/tawk/tawk';
// const dotenv = require('dotenv');
// const dotenvExpand = require('dotenv-expand');
// dotenvExpand(dotenv.config());

// @ts-ignore
// window.ubVersion = process.env.VERSION;

const jss = create({ plugins: [...jssPreset().plugins, rtl()] });

const stateSelector = createStructuredSelector({
  appState: makeSelectAppState(),
  language: makeSelectLanguage(),
  location: makeSelectLocation(),
  loggedIn: makeSelectLoggedIn(),
  theme: makeSelectTheme(),
});

const mainTheme = createMuiTheme({
  typography: {
    fontFamily: 'Open Sans',
    // fontFamily: language === 'ar' ? 'ar' : 'Open Sans',
  },

  // direction: language === 'ar' ? 'rtl' : 'ltr',
  direction: 'ltr',
  palette: {
    secondary: {
      main: '#d20e42',
    },
    primary: {
      main: '#396DE0',
    },
  },
});

const registeredToastContainer = () => {
  if (cookies.get(CookieKeys.Token)) {
    return <RegisteredToastContainer />;
  }
  return <></>;
};

export default function App() {
  const dispatch = useDispatch();

  useEffect(() => {
    if (typeof window === undefined) return;
    const handleRouteChange = (url: string) => {
      logPageView(url);
    };
    //  @ts-ignore
    if (!window.GA_INITIALIZED) {
      initGA();
      //  @ts-ignore
      window.GA_INITIALIZED = true;
    }

    history.listen((e) => handleRouteChange(e.pathname));

    return () => {};
  }, []);

  useResizeAndThemeHook();
  useConnectToCentrifugoAuth();
  useCheckForAuthErrors({ dispatch, loggedInAction, replace });
  useGetInitialData();
  const firstLoad = useRef<boolean>(false);
  const { appState, language, loggedIn, location, theme } = useSelector(
    stateSelector,
  );
  if (
    location &&
    (!location.pathname.includes(AppPages.LoginPage) ||
    !location.pathname.includes(AppPages.SignupPage) ||
    !location.pathname.includes(AppPages.AcountPage))
  ) {
    if (firstLoad.current === false) {
      firstLoad.current = true;
      require('styles/ag-grid.min.css');
      require('styles/ag-theme-balham.min.css');
    }
  }

  if (
    location &&
    (location.pathname.includes(AppPages.AcountPage) ||
    location.pathname.includes(AppPages.Orders) ||
    location.pathname.includes(AppPages.Funds) ||
    location.pathname.includes(AppPages.TradePage))
  ) {
    requestAnimationFrame(() => {
      MessageService.send({
        name: MessageNames.SET_TAB,
        payload: location.pathname,
      });
    });
  }

  const pathname = location?.pathname ?? '';

  const IsUserLoggedIn = (): boolean => {
    if (
      loggedIn ||
      (cookies.get(CookieKeys.Token) &&
        !pathname.includes('login') &&
        !pathname.includes('signup') &&
        !pathname.includes('auth/verify') &&
        !pathname.includes('auth/forgot-password/update')) ||
      pathname.includes('trade')
    ) {
      return true;
    }
    return false;
  };

  return (
    <AppWrapper
      id="appWrapper"
      className={`${IsUserLoggedIn() === true ? `loggedIn` : ''}`}
    >
      <Helmet titleTemplate="%s | UnitedBit" defaultTitle="UnitedBit">
        <meta name="description" content="UnitedBit" />
      </Helmet>
      {/*{useMemo(
        () => (
          <Tawk />
        ),
        [],
      )}*/}
      {useMemo(
        () =>
          process.env.NODE_ENV === 'production' &&
          process.env.IS_DEV_BUILD !== 'true' && <RecaptchaContainer />,
        [],
      )}
      {useMemo(() => registeredToastContainer(), [
        cookies.get(CookieKeys.Token),
      ])}
      <StylesProvider jss={jss}>
        <ThemeProvider theme={mainTheme}>
          {useMemo(() => IsUserLoggedIn() && <AppHeader></AppHeader>, [
            IsUserLoggedIn(),
          ])}
          <Switch>
            <Route exact path={AppPages.HomePage} component={HomePage} />
            <Route path={AppPages.LoginPage} component={LoginPage} />
            <PrivateRoute
              reversePrivate
              path={AppPages.SignupPage}
              component={SignupPage}
            />
            <PrivateRoute path={AppPages.AcountPage} component={AcountPage} />
            <PrivateRoute
              path={AppPages.ChangePassword}
              component={ChangePasswordPage}
            />
            <PrivateRoute
              path={AppPages.PhoneVerification}
              component={PhoneVerificationPage}
            />
            <PrivateRoute
              path={AppPages.AddressManagement}
              component={AddressManagementPage}
            />
            <PrivateRoute
              path={AppPages.UserInfo}
              component={ChangeUserInfoPage}
            />
            <PrivateRoute
              path={AppPages.DocumentVerification}
              component={DocumentVerificationPage}
            />
            <PrivateRoute
              path={AppPages.GoogleAuthentication}
              component={GoogleAuthenticationPage}
            />
            <Route path={AppPages.Orders} component={OrdersPage} />
            <PrivateRoute path={AppPages.Funds} component={FundsPage} />
            <Route
              path={AppPages.VerifyEmail}
              component={EmailVerificationPage}
            />
            <Route path={AppPages.TradePage} component={TradePage} />
            <Route
              path={AppPages.UpdatePassword}
              component={UpdatePasswordPage}
            />
            <Route path={AppPages.ContactUs} component={ContactUsPage} />

            <Route path={AppPages.NotFoundPage} component={NotFoundPage} />
          </Switch>
        </ThemeProvider>
      </StylesProvider>
      {!pathname.includes(AppPages.TradePage) && (
        <Footer className={`${IsUserLoggedIn() === true ? `loggedIn` : ''}`} />
      )}

      <GlobalStyle />
    </AppWrapper>
  );
}

const AppWrapper = styled.div`
  max-width: 100%;
  margin: 0 auto;
  display: flex;
  min-height: 100%;
  flex-direction: column;
  background-color: var(--greyBackground);
  &.loggedIn {
    padding-top: 60px;
  }
`;
