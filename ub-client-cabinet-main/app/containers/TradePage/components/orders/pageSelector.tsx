import React, { useCallback } from 'react';
import { OrderPage } from 'containers/TradePage/constants';
import OpenOrder from 'containers/OrdersPage/pages/openOrder/Loadable';
import OrderHistory from 'containers/OrdersPage/pages/orderHistory/Loadable';
import TradeHistory from 'containers/OrdersPage/pages/tradeHistory/Loadable';
import { OrdersPageState } from 'containers/OrdersPage/types';
import translate from 'containers/TradePage/messages';
import { Button } from '@material-ui/core';
import { FormattedMessage } from 'react-intl';
import { MessageService, MessageNames } from 'services/message_service';
import { Buttons } from 'containers/App/constants';

export default function OrderPageSelector(props: {
  loggedIn: boolean;
  selectedStep: OrderPage;
  ordersData: OrdersPageState;
}) {
  const pageSelector = useCallback(
    (page: OrderPage) => {
      switch (page) {
        case OrderPage.OpenOrders:
          return (
            <OpenOrder
              isMini
              isLoadingOpenOrders={props.ordersData?.isLoadingOpenOrders}
              openOrders={props.ordersData?.openOrders}
            />
          );
        case OrderPage.OrderHistory:
          return (
            <OrderHistory
              isMini
              isLoadingOrderHistory={props.ordersData?.isLoadingOrderHistory}
              orderHistory={props.ordersData?.orderHistory}
            />
          );
        case OrderPage.TradeHistory:
          return (
            <TradeHistory
              isMini
              isLoadingTradeHistory={props.ordersData?.isLoadingTradeHistory}
              tradeHistory={props.ordersData?.tradeHistory}
            />
          );
        default:
          return (
            <OpenOrder
              isMini
              isLoadingOpenOrders={props.ordersData.isLoadingOpenOrders}
              openOrders={props.ordersData.openOrders}
            />
          );
      }
    },
    [props.ordersData],
  );

  return props.loggedIn == undefined || props.loggedIn === false ? (
    <div className="registration">
      <Button
        className={Buttons.SimpleRoundButton}
        onClick={() => {
          MessageService.send({ name: MessageNames.OPEN_LOGIN_POPUP });
        }}
        color="primary"
      >
        <FormattedMessage {...translate.LoginToViewOrders} />
      </Button>
    </div>
  ) : (
    <div className="pageSelectorWrapper">
      {pageSelector(props.selectedStep)}
    </div>
  );
}
