import { useEffect, useRef } from 'react';
import { LocalStorageKeys } from 'services/constants';
import { MessageNames, MessageService } from 'services/message_service';
import { Themes } from '../constants';

export const useResizeAndThemeHook = () => {
  const timeOut = useRef<any>();
  useEffect(() => {
    window.onresize = (e: any) => {
      clearTimeout(timeOut.current);
      timeOut.current = setTimeout(() => {
        MessageService.send({
          name: MessageNames.RESIZE,
          payload: e.target.innerWidth,
        });
      }, 100);
    };
    setTimeout(() => {
      if (
        localStorage[LocalStorageKeys.Theme] &&
        localStorage[LocalStorageKeys.Theme] != ''
      ) {
        const app = document.querySelector('body');
        const hTml = document.querySelector('html');
        if (app && localStorage[LocalStorageKeys.Theme] === Themes.LIGHT) {
          hTml?.classList.remove('htmldark');
          app.classList.remove(Themes.DARK);
          app.classList.add(localStorage[LocalStorageKeys.Theme]);
        } else {
          hTml?.classList.add('htmldark');
        }
      }
    }, 0);

    return () => {};
  }, []);
};
