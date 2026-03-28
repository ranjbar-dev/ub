import { useEffect, useRef } from 'react';
import { LocalStorageKeys } from 'services/constants';
import { CookieKeys, cookies } from 'services/cookie';
import { RegisteredMqttService } from 'services/RegisteredMqttService';
import { storage } from 'utils/storage';

export const useConnectToAuthorizedMqtt2 = () => {
  const mqtt2 = useRef<any>(null);
  useEffect(() => {
    let channel;
    if (cookies.get(CookieKeys.Token)) {
      mqtt2.current = RegisteredMqttService.getInstance(
        cookies.get(CookieKeys.Token),
      );
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
  }, [cookies.get(CookieKeys.Token)]);
};
