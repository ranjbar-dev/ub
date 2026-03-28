import React, { memo, useMemo } from 'react';
import { Helmet } from 'react-helmet';
import { FullPageWrapper } from 'components/wrappers/fullPageWrapper';
//import LocaleToggle from 'containers/LocaleToggle';
import StepSelector from './steps/stepSelector';
function SignupPage () {
  return useMemo(
    () => (
      <FullPageWrapper>
        <Helmet>
          <title>Signup </title>
          <meta name='description' content='Description of Signup Page' />
        </Helmet>

        <div
          className='head darkTheme WhiteHeader'
          style={{ position: 'fixed', zIndex: 1 }}
        >
          {/*<LocaleToggle />*/}
        </div>
        <StepSelector />
      </FullPageWrapper>
    ),
    [],
  );
}

export default memo(SignupPage);
