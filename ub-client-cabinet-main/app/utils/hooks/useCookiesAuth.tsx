import { replace } from 'redux-first-history';
import { AppPages } from 'containers/App/constants';
import { useEffect, useState } from 'react';
import { useDispatch } from 'react-redux';
import { CookieKeys, cookies } from 'services/cookie';

export default function useCookiesAuth () {
  const [Authed, setAuthed] = useState<boolean>(false);

  const dispatch = useDispatch();
  useEffect(() => {
    const token =
      cookies.get(CookieKeys.Token) &&
      //@ts-ignore
      cookies.get(CookieKeys.Token)?.length > 0;
    if (!token && !window.location.href.includes(AppPages.TradePage)) {
      dispatch(replace(AppPages.LoginPage));
      return;
    }
    setAuthed(true);
    return () => {};
  }, [Authed]);
  return Authed;
}
