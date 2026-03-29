import { useEffect, useRef, useState } from 'react';
import { LocalStorageKeys } from 'services/constants';
import { CookieKeys, cookies } from 'services/cookie';
import { RegisteredMqttService } from 'services/RegisteredMqttService';
import { storage } from 'utils/storage';

export const useConnectToAuthorizedMqtt2 = () => {
  const mqtt2 = useRef<any>(null);
  const [token, setToken] = useState(() => cookies.get(CookieKeys.Token));

  const currentToken = cookies.get(CookieKeys.Token);
  if (currentToken !== token) {
    setToken(currentToken);
  }

  useEffect(() => {
    let channel;
    if (token) {
      mqtt2.current = RegisteredMqttService.getInstance(token);
      channel = storage.read(LocalStorageKeys.CHANNEL);
      mqtt2.current.ConnectToSubject({
        subject: `main/trade/user/${channel}/open-orders/`,
      });

      mqtt2.current.ConnectToSubject({
        subject: `main/trade/user/${channel}/crypto-payments/`,
      });
    }
    return () => {
      mqtt2.current &&
        mqtt2.current.DisconnectFromSubject({
          subject: `main/trade/user/${channel}/open-orders/`,
        });

      mqtt2.current &&
        mqtt2.current.DisconnectFromSubject({
          subject: `main/trade/user/${channel}/crypto-payments/`,
        });

      mqtt2.current && (mqtt2.current = null);
    };
  }, [token]);
};
