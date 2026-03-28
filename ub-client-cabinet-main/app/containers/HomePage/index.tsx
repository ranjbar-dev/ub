/*
 * HomePage
 *
 * This is the first thing users see of our App, at the '/' route
 */

import React, { useEffect } from 'react';

import { useDispatch } from 'react-redux';

import { replace } from 'redux-first-history';
import { AppPages } from 'containers/App/constants';
import { CookieKeys, cookies } from 'services/cookie';

export default function HomePage () {
  const dispatch = useDispatch();
  useEffect(() => {
    if (!cookies.get(CookieKeys.Token)) {
      dispatch(replace(AppPages.LoginPage));
    } else {
      dispatch(replace(AppPages.TradePage));
    }

    return () => {};
  }, []);

  return <></>;
}
