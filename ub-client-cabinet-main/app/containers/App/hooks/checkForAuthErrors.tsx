import { useEffect } from 'react';
import { LocalStorageKeys } from 'services/constants';
import { MessageNames, Subscriber } from 'services/message_service';
import { AppPages } from '../constants';

export const useCheckForAuthErrors=({dispatch,loggedInAction,replace})=>{
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: any) => {
      if (message.name === MessageNames.AUTH_ERROR_EVENT) {
        if (location.pathname.includes(AppPages.TradePage)) {
          const theme = localStorage[LocalStorageKeys.Theme];
          const layOut = localStorage[LocalStorageKeys.LAYOUT_NAME];
          const tradeLayout = localStorage[LocalStorageKeys.TRADELAYOUT];
          const countries = localStorage[LocalStorageKeys.COUNTRIES];
          localStorage.clear();
          localStorage[LocalStorageKeys.Theme] = theme;
          localStorage[LocalStorageKeys.COUNTRIES] = countries;
          if (layOut) {
            localStorage[LocalStorageKeys.LAYOUT_NAME] = layOut;
          }
          if (tradeLayout) {
            localStorage[LocalStorageKeys.TRADELAYOUT] = tradeLayout;
          }
          dispatch(loggedInAction(false));
          return;
        }
        if (!location.pathname.includes('login')) {
          dispatch(replace(AppPages.LoginPage));
        }
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, [location.pathname]);
};
