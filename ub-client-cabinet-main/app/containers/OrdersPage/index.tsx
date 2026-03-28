/*
 *
 * OrdersPage
 *
 */

import React, { memo, useState, useEffect } from 'react';
import { Helmet } from 'react-helmet';
import { FormattedMessage } from 'react-intl';
import { useSelector, useDispatch } from 'react-redux';
import { createStructuredSelector } from 'reselect';

import { useInjectSaga } from 'utils/injectSaga';
import { useInjectReducer } from 'utils/injectReducer';
import {
  makeSelectOpenOrders,
  makeSelectOrderHistory,
  makeSelectTradeHistory,
  makeSelectIsLoadingOpenOrders,
  makeSelectIsLoadingOrderHistory,
  makeSelectIsLoadingTradeHistory,
} from './selectors';
import reducer from './reducer';
import saga from './saga';
import translate from './messages';
import { MaxWidthWrapper1600 } from 'components/wrappers/maxWidthWrapper1600';
import { Card } from '@material-ui/core';
import styled from 'styles/styled-components';
import SegmentSelector from './components/segmentSelector';
import TitledComponent from 'components/titled';
import { OrderPages } from './constants';

import OpenOrderPage from './pages/openOrder';
import OrderHistoryPage from './pages/orderHistory';
import TradeHistoryPage from './pages/tradeHistory';
import { Order } from './types';
import { OverflowWrapper } from 'components/wrappers/overflowWrapper';
import { LocalStorageKeys } from 'services/constants';

import { RegisteredUserSubscriber } from 'services/message_service';

const stateSelector = createStructuredSelector({
  openOrders: makeSelectOpenOrders(),
  orderHistory: makeSelectOrderHistory(),
  tradeHistory: makeSelectTradeHistory(),
  isLoadingOpenOrders: makeSelectIsLoadingOpenOrders(),
  isLoadingOrderHistory: makeSelectIsLoadingOrderHistory(),
  isLoadingTradeHistory: makeSelectIsLoadingTradeHistory(),
});

interface Props {}
function OrdersPage (props: Props) {
  useInjectReducer({ key: 'ordersPage', reducer: reducer });
  useInjectSaga({ key: 'ordersPage', saga: saga });
  const {
    openOrders,
    orderHistory,
    tradeHistory,
    isLoadingOpenOrders,
    isLoadingOrderHistory,
    isLoadingTradeHistory,
  } = useSelector(stateSelector);
  const [ActivePage, setActivePage] = useState(OrderPages.OPEN_ORDER);

  const pageSelector = (params: { page: OrderPages; openOrders: Order[] }) => {
    switch (params.page) {
      case OrderPages.OPEN_ORDER:
        return (
          <OpenOrderPage
            isLoadingOpenOrders={isLoadingOpenOrders}
            openOrders={params.openOrders}
          />
        );
      case OrderPages.ORDER_HISTORY:
        return (
          <OrderHistoryPage
            orderHistory={orderHistory}
            isLoadingOrderHistory={isLoadingOrderHistory}
          />
        );
      case OrderPages.TRADE_HISTORY:
        return (
          <TradeHistoryPage
            tradeHistory={tradeHistory}
            isLoadingTradeHistory={isLoadingTradeHistory}
          />
        );

      default:
        return (
          <OpenOrderPage
            isLoadingOpenOrders={isLoadingOpenOrders}
            openOrders={params.openOrders}
          />
        );
    }
  };

  return (
    <div>
      <Helmet>
        <title>Orders</title>
        <meta name='description' content='Description of OrdersPage' />
      </Helmet>
      <MaxWidthWrapper1600>
        <SegmentSelector
          onChange={(page: OrderPages) => {
            setActivePage(page);
          }}
          options={[
            {
              title: <FormattedMessage {...translate.openOrder} />,
              page: OrderPages.OPEN_ORDER,
            },
            {
              title: <FormattedMessage {...translate.orderHistory} />,
              page: OrderPages.ORDER_HISTORY,
            },
            {
              title: <FormattedMessage {...translate.tradeHistory} />,
              page: OrderPages.TRADE_HISTORY,
            },
          ]}
        />
        <OverflowWrapper>
          <MainWrapper>
            <TitledComponent
              id='withHidableBorder'
              title={<FormattedMessage {...translate.history} />}
            >
              {pageSelector({ page: ActivePage, openOrders: openOrders })}
            </TitledComponent>
          </MainWrapper>
        </OverflowWrapper>
      </MaxWidthWrapper1600>
    </div>
  );
}

export default memo(OrdersPage);
const MainWrapper = styled(Card)`
  background: white;
  border-radius: 10px !important;
  box-shadow: none !important;
  height: calc(95vh - 115px);
  min-height: 585px;
  min-width: 1092px;
  margin-top: 2vh;
  display: flex;
  flex-direction: column;
  align-items: center;
`;
