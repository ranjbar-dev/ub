/*
 *
 * GoogleAuthenticationPage
 *
 */

import React, { memo, useEffect } from 'react';
import { Helmet } from 'react-helmet';
import { useSelector, useDispatch } from 'react-redux';
import { createStructuredSelector } from 'reselect';

import { useInjectSaga } from 'utils/injectSaga';
import { useInjectReducer } from 'utils/injectReducer';
import {
  makeSelectQrCode,
  makeSelectIsLoading,
  makeSelectUserData,
} from './selectors';
import reducer from './reducer';
import saga from './saga';
import { MaxWidthWrapper } from 'components/wrappers/maxWidthWrapper';
import BreadCrumb from 'components/BreadCrumb';
import { AppPages } from 'containers/App/constants';
import AnimateChildren from 'components/AnimateChildren';

import { get2faQrCodeAction } from './actions';
import StepSelector from './steps/stepSelector';
import { QrCode } from './types';
import { replace } from 'redux-first-history';

const stateSelector = createStructuredSelector({
  qrCode: makeSelectQrCode(),
  isLoading: makeSelectIsLoading(),
  userData: makeSelectUserData(),
});

interface Props {}

function GoogleAuthenticationPage () {
  useInjectReducer({ key: 'googleAuthenticationPage', reducer: reducer });
  useInjectSaga({ key: 'googleAuthenticationPage', saga: saga });
  let {
    qrCode,
    isLoading,
    userData,
  }: {
    qrCode: QrCode;
    isLoading: boolean;
    userData: any;
  } = useSelector(stateSelector);
  const dispatch = useDispatch();
  useEffect(() => {
    if (userData.google2faEnabled === false) {
      dispatch(get2faQrCodeAction());
    }
    if (userData.google2faEnabled === undefined) {
      dispatch(replace(AppPages.AcountPage));
    }
    return () => {};
  }, []);
  if (userData.google2faEnabled === true) {
    isLoading = false;
  }
  return (
    <div>
      <Helmet>
        <title> Google Authentication </title>
        <meta
          name='description'
          content='Description of GoogleAuthenticationPage'
        />
      </Helmet>
      <MaxWidthWrapper>
        <BreadCrumb
          links={[
            { pageName: 'home', pageLink: AppPages.HomePage },
            {
              pageName: 'acountAndSecurity',
              pageLink: AppPages.AcountPage,
            },
            {
              pageName:
                userData && userData.google2faEnabled === true
                  ? 'DisableGoogleAuthenticator'
                  : 'TwoFA',
              pageLink: AppPages.GoogleAuthentication,
              last: true,
            },
          ]}
        />
        <AnimateChildren isLoading={isLoading}>
          <StepSelector
            qrCode={qrCode}
            userData={userData}
            isAuthenticated={userData.google2faEnabled}
          />
        </AnimateChildren>
      </MaxWidthWrapper>
    </div>
  );
}

export default memo(GoogleAuthenticationPage);
