/*
 *
 * ChangePassword
 *
 */

import React, { memo } from 'react';
import { Helmet } from 'react-helmet';

import { MaxWidthWrapper } from 'components/wrappers/maxWidthWrapper';
import BreadCrumb from 'components/BreadCrumb';
import { AppPages } from 'containers/App/constants';

import StepSelector from './steps/stepSelector';

interface Props {}

function ChangePassword(props: Props) {
  return (
    <div>
      <Helmet>
        <title>Change Password</title>
        <meta name="description" content="Description of ChangePassword" />
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
              pageName: 'changePassword',
              pageLink: AppPages.ChangePassword,
              last: true,
            },
          ]}
        />
        <StepSelector />
      </MaxWidthWrapper>
    </div>
  );
}

export default memo(ChangePassword);
