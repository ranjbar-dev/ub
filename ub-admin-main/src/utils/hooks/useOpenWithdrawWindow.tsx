import WithdrawModal from 'app/containers/Billing/components/WithdrawModal';
import { Payment, PaymentDetails } from 'app/containers/Billing/types';
import { INewWindow } from 'app/NewWindowContainer';
import React, { useEffect } from 'react';
import {
  Subscriber,
  MessageNames,
  MessageService,
  BroadcastMessage,
} from 'services/messageService';

interface WithdrawModalData {
  rowData: Payment;
  details: PaymentDetails;
}

interface Props {}

function useOpenWithdrawWindow(props: Props) {
  useEffect(() => {
    const Subscription = Subscriber.subscribe((message: BroadcastMessage) => {
      if (message.name === MessageNames.SET_MAIN_WITHDRAWALS_ITEM_DETAILS) {
        const payload = message.payload as WithdrawModalData;
        const messageData: INewWindow = {
          name: MessageNames.OPEN_WINDOW,
          payload: {
            windowHeight: 560,
            windowWidth: 1000,
            component: <WithdrawModal initialModalData={payload} />,
            fullTitle: 'WD #' + payload.details.id,
            id: 'WD' + payload.details.id,
          },
        };
        MessageService.send(messageData);
      }
    });
    return () => {
      Subscription.unsubscribe();
    };
  }, []);
}

export default useOpenWithdrawWindow;
