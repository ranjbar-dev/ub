import { useEffect, useRef, useState } from 'react';
import { LocalStorageKeys } from 'services/constants';
import { CookieKeys, cookies } from 'services/cookie';
import { CentrifugoAuthService } from 'services/CentrifugoAuthService';
import { storage } from 'utils/storage';

export const useConnectToCentrifugoAuth = () => {
  const serviceRef = useRef<CentrifugoAuthService | null>(null);
  const [token, setToken] = useState(() => cookies.get(CookieKeys.Token));

  const currentToken = cookies.get(CookieKeys.Token);
  if (currentToken !== token) {
    setToken(currentToken);
  }

  useEffect(() => {
    let channel: string | undefined;
    if (token) {
      serviceRef.current = CentrifugoAuthService.getInstance(token);
      channel = storage.read(LocalStorageKeys.CHANNEL);
      serviceRef.current.ConnectToSubject({
        subject: `user:${channel}:open-orders`,
      });

      serviceRef.current.ConnectToSubject({
        subject: `user:${channel}:crypto-payments`,
      });
    }
    return () => {
      serviceRef.current &&
        serviceRef.current.DisconnectFromSubject({
          subject: `user:${channel}:open-orders`,
        });

      serviceRef.current &&
        serviceRef.current.DisconnectFromSubject({
          subject: `user:${channel}:crypto-payments`,
        });

      serviceRef.current && (serviceRef.current = null);
    };
  }, [token]);
};
