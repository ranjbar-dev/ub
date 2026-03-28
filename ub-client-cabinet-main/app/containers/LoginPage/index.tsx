/*
 *
 * LoginPage
 *
 */

import React from 'react';
import { Helmet } from 'react-helmet';

import { FullPageWrapper } from 'components/wrappers/fullPageWrapper';
//import LocaleToggle from 'containers/LocaleToggle';
import LoginBody from './loginBody';
import { useDispatch } from 'react-redux';
import { replace } from 'redux-first-history';
import { AppPages } from 'containers/App/constants';
import { cookies, CookieKeys } from 'services/cookie';
//import { createStructuredSelector } from 'reselect';
//import { makeSelectLoggedIn } from './selectors';
//import { MessageNames, MessageService } from 'services/message_service';
//import { toast } from 'components/Customized/react-toastify';

//const stateSelector = createStructuredSelector({
//  loggedIn: makeSelectLoggedIn(),
//});

function LoginPage (props: { isPopup?: boolean }) {
  const dispatch = useDispatch();
  const { isPopup } = props;

  const token = cookies.get(CookieKeys.Token);
  if (token && !isPopup) {
    dispatch(replace(AppPages.TradePage));
    return <></>;
  }
  //  const { loggedIn } = useSelector(stateSelector);

  //  useEffect(() => {
  //    const token = cookies.get(CookieKeys.Token);
  //    if (token && !isPopup && loggedIn) {
  //      dispatch(replace(AppPages.AcountPage));
  //    } else if (token && isPopup && loggedIn) {
  //      MessageService.send({
  //        name: MessageNames.CLOSE_MODAL,
  //      });
  //      toast.success('Successfully Logged In ');
  //    }

  //    return () => {};
  //  }, [loggedIn]);
  return (
    <>
      {!isPopup ? (
        <FullPageWrapper>
          <Helmet>
            <title>Login </title>
            <meta name='description' content='Description of LoginPage' />
          </Helmet>

          <div className='head darkTheme WhiteHeader'>
            {/*<LocaleToggle />*/}
          </div>
          <div className='body'>{<LoginBody />}</div>
        </FullPageWrapper>
      ) : (
        <LoginBody {...props} />
      )}
    </>
  );
}

export default LoginPage;
