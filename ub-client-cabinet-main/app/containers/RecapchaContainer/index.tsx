/*
 *
 * RecapchaContainer
 *
 */

import React, { useEffect, useState } from 'react';
import { useDispatch } from 'react-redux';

import { useInjectSaga } from 'utils/injectSaga';
import { useInjectReducer } from 'utils/injectReducer';
import reducer from './reducer';
import saga from './saga';
import { getRecapchaAction } from './actions';
import { GoogleReCaptchaProvider } from 'react-google-recaptcha-v3';
import { GoogleRecaptchaComponent } from './recapchaComponent';
import { MessageNames, Subscriber } from 'services/message_service';
import { SessionStorageKeys } from 'services/constants';

function RecaptchaContainer() {
  useInjectReducer({ key: 'recapchaContainer', reducer: reducer });
  useInjectSaga({ key: 'recapchaContainer', saga: saga });
  //  const { recapcha } = useSelector(stateSelector);
  const [SiteKEy, setSiteKEy] = useState(
    sessionStorage[SessionStorageKeys.SITE_KEY] ?? '',
  );
  const dispatch = useDispatch();
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.RESET_SITE_KEY) {
        // setSiteKEy('');
        dispatch(getRecapchaAction());
      }
      if (message.name === MessageNames.SET_SITE_KEY) {
        setSiteKEy(message.payload);
      }
    });

    if (!SiteKEy || !sessionStorage[SessionStorageKeys.SITE_KEY]) {
      dispatch(getRecapchaAction());
    }

    return () => {
      Subscription.unsubscribe();
      //  console.log('recapcha dismounted');
    };
  }, [SiteKEy]);

  return (
    <div className="recapcha">
      {SiteKEy && (
        <GoogleReCaptchaProvider reCaptchaKey={SiteKEy}>
          <GoogleRecaptchaComponent />
        </GoogleReCaptchaProvider>
      )}
    </div>
  );
}

export default RecaptchaContainer;
