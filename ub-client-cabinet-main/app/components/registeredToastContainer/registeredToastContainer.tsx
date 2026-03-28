import { toast } from 'components/Customized/react-toastify';
import { OrderPages } from 'containers/OrdersPage/constants';
import useOnlineStatus from 'hooks/onlineStatusHook/provider';
import React, { useEffect, useRef } from 'react';
import { LocalStorageKeys } from 'services/constants';
import { CookieKeys, cookies } from 'services/cookie';
import {
  EventMessageNames,
  EventMessageService,
  MessageNames,
  MessageService,
  RegisteredUserSubscriber,
} from 'services/message_service';
import { CurrencyFormater } from 'utils/formatters';
import { storage } from 'utils/storage';
import { useCheckConnectionPeriodically } from './useCheckConnectionPeriodically';

const toastMessage = payload => {
  const coin = payload.pairCurrency.split('-')[0].toString();
  const otherCoin = payload.pairCurrency.split('-')[1].toString();

  if (payload.status === 'open') {
    payload.status = 'Placed';
  }
  if (payload.status === 'filled') {
    payload.status = 'Filled';
  }
  if (payload.status === 'canceled') {
    payload.status = 'Canceled';
  }

  const formattedAmount =
    payload.amount != '' && payload.amount != null
      ? CurrencyFormater(payload.amount) + ''
      : '';

  const formattedPrice =
    payload.price != '' && payload.price != null
      ? CurrencyFormater(payload.price) + ' ' + otherCoin + ' |'
      : 'Market price' + ' | ';

  const t = [
    payload.type == 'buy' ? 'BUY |' : 'SELL |',
    `${formattedAmount} ${coin} ${payload.amount != '' ? ' ON' : ' '}`,
    `${formattedPrice} [${payload.status}]`,
  ];

  const toastText = t.join(' ');
  switch (payload.status) {
    case 'Placed':
      toast.info(toastText);
      break;
    case 'Canceled':
      toast.warn(toastText);
      break;
    case 'Filled':
      toast.success(toastText);
      break;
    default:
      break;
  }
};

const RegisteredToastContainer = () => {
  useCheckConnectionPeriodically();

  useEffect(() => {
    let RegisteredUserSubscription;
    if (cookies.get(CookieKeys.Token)) {
      RegisteredUserSubscription = RegisteredUserSubscriber.subscribe(
        (message: any) => {
          if (
            message.name &&
            localStorage[LocalStorageKeys.CHANNEL] &&
            message.name.includes(storage.read(LocalStorageKeys.CHANNEL))
          ) {
            if (message.name.includes('crypto-payments')) {
              toast.info('Balance updated');
              MessageService.send({
                name: MessageNames.NEW_ORDER_NOTIFICATION,
              });
              return;
            }

            const payload: any = message.payload;
            MessageService.send({
              name: MessageNames.NEW_ORDER_NOTIFICATION,
              payload,
            });
            toastMessage(payload);
            const eventPayload = {
              name: EventMessageNames.REFRESH_ORDER_GRID,
              ...((payload.status === 'open' ||
                payload.status === 'Placed') && { id: OrderPages.OPEN_ORDER }),
            };
            EventMessageService.send(eventPayload);
          }
        },
      );
      if (process.env.NODE_ENV !== 'production') {
        console.log('%cregistered toaster connected', 'color:green');
      }
    }
    return () => {
      RegisteredUserSubscription && RegisteredUserSubscription.unsubscribe();
      if (process.env.NODE_ENV !== 'production') {
        console.log('%cregistered toaster disconnected', 'color:red;');
      }
    };
  }, [cookies.get(CookieKeys.Token)]);

  return <></>;
};

export default RegisteredToastContainer;
