/*
 *
 * ChangeUserInfoPage
 *
 */

import React, { memo, useEffect } from 'react';
import { Helmet } from 'react-helmet';
import { useSelector, useDispatch } from 'react-redux';
import { createStructuredSelector } from 'reselect';

import { useInjectSaga } from 'utils/injectSaga';
import { useInjectReducer } from 'utils/injectReducer';
import {
  makeSelectUserProfileData,
  makeSelectIsLoadingUserProfileData,
} from './selectors';
import reducer from './reducer';
import saga from './saga';
import { MaxWidthWrapper } from 'components/wrappers/maxWidthWrapper';
import BreadCrumb from 'components/BreadCrumb';
import { AppPages } from 'containers/App/constants';

import AnimateChildren from 'components/AnimateChildren';
import { getUserProfileAction } from './actions';
import ComponentsWrapper from './components';

const stateSelector = createStructuredSelector({
  userProfileData: makeSelectUserProfileData(),
  isLoading: makeSelectIsLoadingUserProfileData(),
});

interface Props {}
function ChangeUserInfoPage(props: Props) {
  // Warning: Add your key to RootState in types/index.d.ts file
  useInjectReducer({ key: 'changeUserInfoPage', reducer: reducer });
  useInjectSaga({ key: 'changeUserInfoPage', saga: saga });

  const { userProfileData, isLoading } = useSelector(stateSelector);
  const dispatch = useDispatch();
  useEffect(() => {
    //if (!userProfileData.updatedAt) {
    dispatch(getUserProfileAction());
    //}
    return () => {};
  }, []);
  return (
    <>
      <Helmet>
        <title>change user info</title>
        <meta name="description" content="Description of change user info" />
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
              pageName: 'ChangeInfo',
              pageLink: AppPages.UserInfo,
              last: true,
            },
          ]}
        />
        <AnimateChildren isLoading={isLoading}>
          <ComponentsWrapper data={userProfileData} />
        </AnimateChildren>
      </MaxWidthWrapper>
    </>
  );
}

export default memo(ChangeUserInfoPage);
