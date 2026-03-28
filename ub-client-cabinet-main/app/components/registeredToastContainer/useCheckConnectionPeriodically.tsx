import { useOnlineStatus } from 'hooks/onlineStatusHook/provider';
import { useEffect, useRef } from 'react';
import { EventMessageService, MessageNames } from 'services/message_service';

export const useCheckConnectionPeriodically=()=>{

  const isConnected = useOnlineStatus();
  const initialized = useRef(false);
  useEffect(() => {
    if (isConnected === true && initialized.current === true) {
      if (process.env.NODE_ENV !== 'production') {
        console.log(
          `%csending reconnect signal ${new Date().toISOString()}`,
          'color:green;font-size:11px;',
        );
      }
      EventMessageService.send({
        name: MessageNames.RECONNECT_EVENT,
      });
    }
    if (isConnected === true) {
      initialized.current = true;
    }
    return () => {};
  }, [isConnected]);

};
