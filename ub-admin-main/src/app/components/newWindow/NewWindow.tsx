import { WindowTypes } from 'app/constants';
import { VerificationWindow } from 'app/containers/VerificationWindow';
import { InitialUserDetails } from 'app/containers/UserAccounts/types';
import { INewWindow } from 'app/NewWindowContainer';
import React, { memo, useState, useEffect } from 'react';
import NewWindow from 'react-new-window';
import {
  Subscriber,
  MessageNames,
  MessageService,
  BroadcastMessage,
} from 'services/messageService';

import UserDetailsWindow from '../UserDetailsWindow/UserDetailsWindow';

/**
 * Message-bus listener that opens UserDetailsWindow or VerificationWindow
 * in a new browser window when an OPEN_NEW_WINDOW event is received.
 * Renders nothing itself — mount once at the app root level.
 *
 * @example
 * ```tsx
 * <UBWindow />
 * ```
 */
function UBWindow() {
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: BroadcastMessage) => {
      if (message.name === MessageNames.OPEN_NEW_WINDOW) {
        ////
        const payload = message.payload as InitialUserDetails;
        let id =
          message.type === WindowTypes.User
            ? 'UD #' + payload.id
            : 'UV #' + payload.id;
        ////
        const messageData: INewWindow = {
          name: MessageNames.OPEN_WINDOW,
          payload: {
            windowHeight: 745,
            windowWidth: 1175,
            component:
              message.type === WindowTypes.User ? (
                <UserDetailsWindow initialData={payload} />
              ) : (
                <VerificationWindow initialData={payload} />
              ),
            fullTitle: id,
            id,
          },
        };
        MessageService.send(messageData);
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);

  return <></>;
}

export default memo(UBWindow);
