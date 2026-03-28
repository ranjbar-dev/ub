import { useEffect } from 'react';
import { LocalStorageKeys } from 'services/constants';
import { RegisteredUserSubscriber } from 'services/message_service';
import { storage } from 'utils/storage';

export const useNewPaymentEvent = ({
  toRunAfterNewEvent,
  dependencies,
}: {
  toRunAfterNewEvent: () => void;
  dependencies?: any[];
}) => {
  useEffect(() => {
    const RegisteredUserSubscription = RegisteredUserSubscriber.subscribe(
      (message: any) => {
        if (
          message.name &&
          localStorage[LocalStorageKeys.CHANNEL] &&
          message.name.includes(storage.read(LocalStorageKeys.CHANNEL))
        ) {
          toRunAfterNewEvent();
        }
      },
    );

    return () => {
      RegisteredUserSubscription.unsubscribe();
    };
  }, dependencies ?? []);
};
