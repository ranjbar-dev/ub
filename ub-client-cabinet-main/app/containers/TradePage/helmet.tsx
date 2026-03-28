import React, { useEffect } from 'react';
import {
  SideSubscriber,
  MessageNames,
  MarketWatchSubscriber,
} from 'services/message_service';
import { savedPairName } from 'utils/sharedData';
import { PairFormat, CurrencyFormater } from 'utils/formatters';
export default function TradeHelmet () {
  let pairName = savedPairName();
  useEffect(() => {
    const BSubscription = SideSubscriber.subscribe((message: any) => {
      if (message.name === MessageNames.SET_TRADE_PAGE_CURRENCY_PAIR) {
        document.title = 'loading...';
        pairName = message.payload.name;
      }
    });
    const MarketWatchSubscription = MarketWatchSubscriber.subscribe(
      (message: any) => {
        if (message.payload.name === pairName) {
          document.title = message.payload.price
            ? CurrencyFormater(message.payload.price + '') +
              ' | ' +
              PairFormat(pairName)
            : 'loading...';
        }
      },
    );
    return () => {
      MarketWatchSubscription.unsubscribe();
      BSubscription.unsubscribe();
      pairName = savedPairName();
    };
  }, []);

  return <></>;
}
