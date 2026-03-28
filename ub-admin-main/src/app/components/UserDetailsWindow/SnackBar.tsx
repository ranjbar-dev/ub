import React, { memo, useEffect, useState } from 'react';
import { Subscriber, MessageNames } from 'services/messageService';
import type { BroadcastMessage } from 'services/messageService';

interface Props {
  userId: number | string;
  priority?: string;
}

/** Notification message shape displayed inside the SnackBar. */
interface ToastMessage {
  value?: string;
  type?: string;
}

/**
 * Toast notification bar scoped to a specific user window.
 * Listens for TOAST messages on the message bus and auto-hides after 3 s.
 *
 * @example
 * ```tsx
 * <SnackBar userId={123} />
 * ```
 */
function SnackBar(props: Props) {
  const { userId, priority } = props;
  //  const { enqueueSnackbar, closeSnackbar } = useSnackbar();
  const [ShowToast, setShowToast] = useState(false);
  const [Message, setMessage] = useState<ToastMessage>({});

  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: BroadcastMessage) => {
      if (message.name === MessageNames.TOAST && message.userId === userId) {
        setMessage({
          value: message.payload as string,
          type: message.type,
        });
        setShowToast(true);
        setTimeout(() => {
          setShowToast(false);
        }, 3000);
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);

  return (
    <div
      role="alert"
      aria-live="polite"
      className={`UbToast ${ShowToast === true ? 'show' : ''} ${
        priority ?? ''
      }`}
    >
      <div
        className={`content ${Message.type}`}
        aria-label={Message.value}
      >
        {Message.value}
      </div>
    </div>
  );
}

export default memo(SnackBar);
