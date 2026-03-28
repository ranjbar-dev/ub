/*
 * OrdersPage Messages
 *
 * This contains all the text for the OrdersPage container.
 */

import { defineMessages } from 'react-intl';
import { GlobalTranslateScope } from 'containers/App/constants';

export const scope = 'containers.OrdersPage';

export default defineMessages({
  header: {
    id: `${scope}.header`,
    defaultMessage: 'This is the OrdersPage container!',
  },
  history: {
    id: `${scope}.history`,
    defaultMessage: 'ET_history',
  },
  openOrder: {
    id: `${scope}.openOrder`,
    defaultMessage: 'ET_openOrder',
  },
  orderHistory: {
    id: `${scope}.orderHistory`,
    defaultMessage: 'ET_orderHistory',
  },
  tradeHistory: {
    id: `${scope}.tradeHistory`,
    defaultMessage: 'ET_tradeHistory',
  },

  Gototrade: {
    id: `${scope}.Gototrade`,
    defaultMessage: 'ET_Gototrade',
  },
  timePeriod: {
    id: `${scope}.timePeriod`,
    defaultMessage: 'ET_timePeriod',
  },
  Hidecancelledorders: {
    id: `${scope}.Hidecancelledorders`,
    defaultMessage: 'ET_Hidecancelledorders',
  },
  week1: {
    id: `${scope}.week1`,
    defaultMessage: 'ET_week1',
  },
  month1: {
    id: `${scope}.month1`,
    defaultMessage: 'ET_month1',
  },
  month3: {
    id: `${scope}.month3`,
    defaultMessage: 'ET_month3',
  },
  /////////////
  search: {
    id: `${GlobalTranslateScope}.search`,
    defaultMessage: 'ET_search',
  },
  cancel: {
    id: `${GlobalTranslateScope}.cancel`,
    defaultMessage: 'Cancel',
  },
  Reset: {
    id: `${GlobalTranslateScope}.reset`,
    defaultMessage: 'Reset',
  },
  coin: {
    id: `${GlobalTranslateScope}.coin`,
    defaultMessage: 'ET_Coin',
  },
  type: {
    id: `${GlobalTranslateScope}.type`,
    defaultMessage: 'ET.type',
  },
  all: {
    id: `${GlobalTranslateScope}.all`,
    defaultMessage: 'ET_all',
  },
  buy: {
    id: `${GlobalTranslateScope}.buy`,
    defaultMessage: 'ET_buy',
  },
  sell: {
    id: `${GlobalTranslateScope}.sell`,
    defaultMessage: 'ET_sell',
  },
});
